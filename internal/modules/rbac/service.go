package rbac

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
)

const (
	defaultPerPage = 20
	maxPerPage     = 100

	// DashboardAccessFeature is the feature that lets a user enter /admin.
	DashboardAccessFeature = "dashboard.access"

	// AdminRoleName is the bootstrap role that must keep dashboard access and
	// cannot be deleted.
	AdminRoleName = "admin"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListRoles(ctx context.Context) []RoleResponse {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		panic(err)
	}
	items := make([]RoleResponse, 0, len(roles))
	for i := range roles {
		items = append(items, toRoleResponse(&roles[i]))
	}
	return items
}

func (s *Service) GetRole(ctx context.Context, publicID string) RoleResponse {
	role, err := s.repo.GetRole(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("role not found"))
		}
		panic(err)
	}
	return toRoleResponse(role)
}

func (s *Service) CreateRole(ctx context.Context, in RoleInput) RoleResponse {
	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	if name == "" {
		panic(apperror.BadRequest("name is required"))
	}
	featureIDs := dedupe(in.FeatureIDs)

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	s.validateFeatures(ctx, tx, featureIDs)

	created, err := s.repo.CreateRole(ctx, tx, name, description)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			panic(apperror.Conflict("a role with that name already exists"))
		}
		panic(err)
	}
	if err := s.repo.ReplaceRoleFeatures(ctx, tx, created.ID, featureIDs); err != nil {
		panic(err)
	}
	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}

	role, err := s.repo.GetRole(ctx, created.PublicID)
	if err != nil {
		panic(err)
	}
	return toRoleResponse(role)
}

func (s *Service) UpdateRole(ctx context.Context, actorPublicID, publicID string, in RoleInput) RoleResponse {
	role, err := s.repo.GetRole(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("role not found"))
		}
		panic(err)
	}

	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	if name == "" {
		panic(apperror.BadRequest("name is required"))
	}
	featureIDs := dedupe(in.FeatureIDs)

	if role.Name == AdminRoleName {
		if name != AdminRoleName {
			panic(apperror.Conflict("the admin role cannot be renamed"))
		}
		if !contains(featureIDs, DashboardAccessFeature) {
			panic(apperror.Conflict("the admin role must keep dashboard.access"))
		}
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	s.validateFeatures(ctx, tx, featureIDs)

	if contains(role.Features, DashboardAccessFeature) && !contains(featureIDs, DashboardAccessFeature) {
		s.guardSelfLockout(ctx, tx, actorPublicID, role)
	}

	if _, err := s.repo.UpdateRole(ctx, tx, publicID, name, description); err != nil {
		if utils.IsUniqueViolation(err) {
			panic(apperror.Conflict("a role with that name already exists"))
		}
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("role not found"))
		}
		panic(err)
	}
	if err := s.repo.ReplaceRoleFeatures(ctx, tx, role.ID, featureIDs); err != nil {
		panic(err)
	}
	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}

	updated, err := s.repo.GetRole(ctx, publicID)
	if err != nil {
		panic(err)
	}
	return toRoleResponse(updated)
}

func (s *Service) DeleteRole(ctx context.Context, actorPublicID, publicID string) {
	role, err := s.repo.GetRole(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("role not found"))
		}
		panic(err)
	}
	if role.Name == AdminRoleName {
		panic(apperror.Conflict("the admin role cannot be deleted"))
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if contains(role.Features, DashboardAccessFeature) {
		s.guardSelfLockout(ctx, tx, actorPublicID, role)
	}

	affected, err := s.repo.SoftDeleteRole(ctx, tx, publicID)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.NotFound("role not found"))
	}
	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
}

func (s *Service) ListFeatures(ctx context.Context) []FeatureResponse {
	features, err := s.repo.ListFeatures(ctx)
	if err != nil {
		panic(err)
	}
	items := make([]FeatureResponse, 0, len(features))
	for i := range features {
		items = append(items, FeatureResponse{
			ID:          features[i].PublicID,
			Code:        features[i].Code,
			Name:        features[i].Name,
			Description: features[i].Description,
		})
	}
	return items
}

func (s *Service) ListUsers(ctx context.Context, search string, page, perPage int) *UserListResult {
	page, perPage = normalizePage(page, perPage)
	search = strings.TrimSpace(search)

	users, err := s.repo.ListUsers(ctx, search, int32(perPage), int32((page-1)*perPage))
	if err != nil {
		panic(err)
	}
	total, err := s.repo.CountUsers(ctx, search)
	if err != nil {
		panic(err)
	}

	items := make([]UserResponse, 0, len(users))
	for i := range users {
		items = append(items, toUserResponse(&users[i]))
	}
	return &UserListResult{Items: items, Total: total, Page: page, PerPage: perPage}
}

func (s *Service) UpdateUserRoles(ctx context.Context, actorPublicID, userPublicID string, roleIDs []string) UserResponse {
	user, err := s.repo.GetUser(ctx, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("user not found"))
		}
		panic(err)
	}

	roleIDs = dedupe(roleIDs)

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if len(roleIDs) > 0 {
		count, err := s.repo.CountRolesByPublicIDs(ctx, tx, roleIDs)
		if err != nil {
			panic(err)
		}
		if count != int64(len(roleIDs)) {
			panic(apperror.BadRequest("unknown role"))
		}
	}

	if err := s.repo.ReplaceUserRoles(ctx, tx, user.ID, roleIDs); err != nil {
		panic(err)
	}

	if actorPublicID == userPublicID {
		features, err := s.repo.FeatureCodes(ctx, tx, userPublicID)
		if err != nil {
			panic(err)
		}
		if !contains(features, DashboardAccessFeature) {
			panic(apperror.Conflict("cannot remove your own dashboard access"))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}

	updated, err := s.repo.GetUser(ctx, userPublicID)
	if err != nil {
		panic(err)
	}
	return toUserResponse(updated)
}

func (s *Service) validateFeatures(ctx context.Context, db sqlc.DBTX, featureIDs []string) {
	if len(featureIDs) == 0 {
		return
	}
	count, err := s.repo.CountFeaturesByPublicIDs(ctx, db, featureIDs)
	if err != nil {
		panic(err)
	}
	if count != int64(len(featureIDs)) {
		panic(apperror.BadRequest("unknown feature"))
	}
}

// guardSelfLockout refuses to remove dashboard.access from the actor's only
// role that grants it.
func (s *Service) guardSelfLockout(ctx context.Context, tx pgx.Tx, actorPublicID string, role *Role) {
	if actorPublicID == "" {
		return
	}
	has, err := s.repo.UserHasRole(ctx, actorPublicID, role.PublicID)
	if err != nil {
		panic(err)
	}
	if !has {
		return
	}
	other, err := s.repo.FeatureCodesExcludingRole(ctx, tx, actorPublicID, role.ID)
	if err != nil {
		panic(err)
	}
	if !contains(other, DashboardAccessFeature) {
		panic(apperror.Conflict("this change would remove your own dashboard access"))
	}
}

func normalizePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	return page, perPage
}

func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func toRoleResponse(role *Role) RoleResponse {
	return RoleResponse{
		ID:          role.PublicID,
		Name:        role.Name,
		Description: role.Description,
		Features:    nonNil(role.Features),
		MemberCount: role.MemberCount,
		CreateTime:  role.CreateTime.Format(time.RFC3339),
		UpdateTime:  role.UpdateTime.Format(time.RFC3339),
	}
}

func toUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:          user.PublicID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Phone:       user.Phone,
		IsVerified:  user.IsVerified,
		Roles:       nonNil(user.Roles),
		CreateTime:  user.CreateTime.Format(time.RFC3339),
	}
}

func nonNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
