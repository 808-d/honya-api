package rbac

import (
	"context"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: sqlc.New()}
}

// Begin starts a transaction owned by the service.
func (r *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *Repository) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := r.q.ListRoles(ctx, r.db)
	if err != nil {
		return nil, err
	}
	roles := make([]Role, 0, len(rows))
	for _, row := range rows {
		roles = append(roles, Role{
			ID:          row.ID,
			PublicID:    row.PublicID,
			Name:        row.Name,
			Description: row.Description,
			Features:    row.Features,
			MemberCount: row.MemberCount,
			CreateTime:  row.CreateTime.Time,
			UpdateTime:  row.UpdateTime.Time,
		})
	}
	return roles, nil
}

func (r *Repository) GetRole(ctx context.Context, publicID string) (*Role, error) {
	row, err := r.q.GetRoleByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Name:        row.Name,
		Description: row.Description,
		Features:    row.Features,
		MemberCount: row.MemberCount,
		CreateTime:  row.CreateTime.Time,
		UpdateTime:  row.UpdateTime.Time,
	}, nil
}

func (r *Repository) CreateRole(ctx context.Context, db sqlc.DBTX, name, description string) (*Role, error) {
	row, err := r.q.CreateRole(ctx, db, sqlc.CreateRoleParams{Name: name, Description: description})
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Name:        row.Name,
		Description: row.Description,
		Features:    []string{},
		CreateTime:  row.CreateTime.Time,
		UpdateTime:  row.UpdateTime.Time,
	}, nil
}

func (r *Repository) UpdateRole(ctx context.Context, db sqlc.DBTX, publicID, name, description string) (*Role, error) {
	row, err := r.q.UpdateRole(ctx, db, sqlc.UpdateRoleParams{
		Name:        name,
		Description: description,
		PublicID:    publicID,
	})
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Name:        row.Name,
		Description: row.Description,
		Features:    []string{},
		CreateTime:  row.CreateTime.Time,
		UpdateTime:  row.UpdateTime.Time,
	}, nil
}

func (r *Repository) SoftDeleteRole(ctx context.Context, db sqlc.DBTX, publicID string) (int64, error) {
	return r.q.SoftDeleteRole(ctx, db, publicID)
}

func (r *Repository) ListFeatures(ctx context.Context) ([]Feature, error) {
	rows, err := r.q.ListFeatures(ctx, r.db)
	if err != nil {
		return nil, err
	}
	features := make([]Feature, 0, len(rows))
	for _, row := range rows {
		features = append(features, Feature{
			ID:          row.ID,
			PublicID:    row.PublicID,
			Code:        row.Code,
			Name:        row.Name,
			Description: row.Description,
		})
	}
	return features, nil
}

func (r *Repository) CountFeaturesByPublicIDs(ctx context.Context, db sqlc.DBTX, ids []string) (int64, error) {
	return r.q.CountFeaturesByPublicIDs(ctx, db, ids)
}

func (r *Repository) ReplaceRoleFeatures(ctx context.Context, db sqlc.DBTX, roleID int64, featureIDs []string) error {
	if err := r.q.DeleteRoleFeatures(ctx, db, roleID); err != nil {
		return err
	}
	for _, id := range featureIDs {
		affected, err := r.q.InsertRoleFeature(ctx, db, sqlc.InsertRoleFeatureParams{
			RoleID:          roleID,
			FeaturePublicID: id,
		})
		if err != nil {
			return err
		}
		if affected == 0 {
			continue
		}
	}
	return nil
}

func (r *Repository) CountRolesByPublicIDs(ctx context.Context, db sqlc.DBTX, ids []string) (int64, error) {
	return r.q.CountRolesByPublicIDs(ctx, db, ids)
}

func (r *Repository) ListUsers(ctx context.Context, search string, limit, offset int32) ([]User, error) {
	rows, err := r.q.ListUsers(ctx, r.db, sqlc.ListUsersParams{
		Search:     utils.TextFromString(search),
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toUser(row.ID, row.PublicID, row.Username, row.DisplayName, row.Email, row.Phone, row.IsVerified, row.CreateTime, row.Roles))
	}
	return users, nil
}

func (r *Repository) CountUsers(ctx context.Context, search string) (int64, error) {
	return r.q.CountUsers(ctx, r.db, utils.TextFromString(search))
}

func (r *Repository) GetUser(ctx context.Context, publicID string) (*User, error) {
	row, err := r.q.GetUserByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	u := toUser(row.ID, row.PublicID, row.Username, row.DisplayName, row.Email, row.Phone, row.IsVerified, row.CreateTime, row.Roles)
	return &u, nil
}

func (r *Repository) ReplaceUserRoles(ctx context.Context, db sqlc.DBTX, userID int64, roleIDs []string) error {
	if err := r.q.DeleteUserRoles(ctx, db, userID); err != nil {
		return err
	}
	for _, id := range roleIDs {
		affected, err := r.q.InsertUserRole(ctx, db, sqlc.InsertUserRoleParams{
			UserID:       userID,
			RolePublicID: id,
		})
		if err != nil {
			return err
		}
		if affected == 0 {
			continue
		}
	}
	return nil
}

func (r *Repository) UserHasRole(ctx context.Context, userPublicID, rolePublicID string) (bool, error) {
	return r.q.UserHasRole(ctx, r.db, sqlc.UserHasRoleParams{
		UserPublicID: userPublicID,
		RolePublicID: rolePublicID,
	})
}

func (r *Repository) FeatureCodesExcludingRole(ctx context.Context, db sqlc.DBTX, userPublicID string, roleID int64) ([]string, error) {
	return r.q.GetFeatureCodesByPublicIDExcludingRole(ctx, db, sqlc.GetFeatureCodesByPublicIDExcludingRoleParams{
		PublicID: userPublicID,
		RoleID:   roleID,
	})
}

func (r *Repository) FeatureCodes(ctx context.Context, db sqlc.DBTX, publicID string) ([]string, error) {
	return r.q.GetFeatureCodesByPublicID(ctx, db, publicID)
}

func toUser(id int64, publicID, username, displayName, email string, phone pgtype.Text, verified bool, created pgtype.Timestamptz, roles []string) User {
	return User{
		ID:          id,
		PublicID:    publicID,
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Phone:       utils.TextPtr(phone),
		IsVerified:  verified,
		Roles:       roles,
		CreateTime:  created.Time,
	}
}
