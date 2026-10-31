package audit

import (
	"context"
	"fmt"

	"github.com/808-d/honya-api/internal/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const actorSetting = "app.actor_public_id"

// Begin starts a transaction and attributes writes in it to the authenticated
// user in the request context. The actor setting is transaction-local, so it
// cannot leak to another request through a pooled connection.
func Begin(ctx context.Context, db *pgxpool.Pool) (pgx.Tx, error) {
	return BeginAs(ctx, db, middleware.UserIDFromContext(ctx))
}

// BeginAs starts a transaction with an explicit actor. Use it for operations
// where the actor is not the authenticated caller, such as account verification
// or password reset. An empty actor records writes with a NULL actor_user_id.
func BeginAs(ctx context.Context, db *pgxpool.Pool, actorPublicID string) (pgx.Tx, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, "SELECT set_config($1, $2, true)", actorSetting, actorPublicID); err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("set audit actor: %w", err)
	}

	return tx, nil
}
