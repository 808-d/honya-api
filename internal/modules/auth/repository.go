package auth

import (
	"context"

	"github.com/808-d/honya-api/internal/database/sqlc"
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

// FindByIdentifier looks up an active user by email or username.
func (r *Repository) FindByIdentifier(ctx context.Context, identifier string) (*User, error) {
	row, err := r.q.FindUserByIdentifier(ctx, r.db, identifier)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID,
		PublicID:     row.PublicID,
		Username:     row.Username,
		DisplayName:  row.DisplayName,
		HashPassword: row.HashPassword,
		Email:        row.Email,
		Phone:        textPtr(row.Phone),
		IsVerified:   row.IsVerified,
	}, nil
}

// FindByEmail looks up an active user by email.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.q.FindUserByEmail(ctx, r.db, email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID,
		PublicID:     row.PublicID,
		Username:     row.Username,
		DisplayName:  row.DisplayName,
		HashPassword: row.HashPassword,
		Email:        row.Email,
		Phone:        textPtr(row.Phone),
		IsVerified:   row.IsVerified,
	}, nil
}

// EmailExists reports whether an active user already uses the email.
func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	return r.q.EmailExists(ctx, r.db, email)
}

// UsernameExists reports whether an active user already uses the username.
func (r *Repository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return r.q.UsernameExists(ctx, r.db, username)
}

// Create inserts a new unverified user.
func (r *Repository) Create(ctx context.Context, username, displayName, hashPassword, email string) (*User, error) {
	row, err := r.q.CreateUser(ctx, r.db, sqlc.CreateUserParams{
		Username:     username,
		DisplayName:  displayName,
		HashPassword: hashPassword,
		Email:        email,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Email:       row.Email,
	}, nil
}

// Verify marks the user as verified and returns the updated row.
func (r *Repository) Verify(ctx context.Context, publicID string) (*User, error) {
	row, err := r.q.VerifyUser(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Email:       row.Email,
		Phone:       textPtr(row.Phone),
		IsVerified:  true,
	}, nil
}

// AssignRoleByName grants the named role to a user if the role exists.
func (r *Repository) AssignRoleByName(ctx context.Context, userID int64, name string) error {
	return r.q.AssignRoleByName(ctx, r.db, sqlc.AssignRoleByNameParams{
		UserID: userID,
		Name:   name,
	})
}

// UpdatePassword replaces the user's password hash.
func (r *Repository) UpdatePassword(ctx context.Context, publicID, hashPassword string) error {
	return r.q.UpdateUserPassword(ctx, r.db, sqlc.UpdateUserPasswordParams{
		PublicID:     publicID,
		HashPassword: hashPassword,
	})
}

// RolesByUserID returns the role names assigned to a user.
func (r *Repository) RolesByUserID(ctx context.Context, userID int64) ([]string, error) {
	return r.q.GetRoleNamesByUserID(ctx, r.db, userID)
}

func textPtr(t pgtype.Text) *string {
	if t.Valid {
		s := t.String
		return &s
	}
	return nil
}
