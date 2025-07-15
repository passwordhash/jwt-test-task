package tokens

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	repoErr "github.com/passwordhash/jwt-test-task/internal/storage/errors"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Storage struct {
	db DB
}

func New(db DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Save(
	ctx context.Context,
	userID, token, userAgent, ip string,
) (string, error) {
	const op = "storage.tokens.Save"

	query := `
	INSERT INTO refresh_tokens (user_id, token, user_agent, ip_address)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id, user_agent) 
	DO UPDATE SET 
    		token = EXCLUDED.token,
		ip_address = EXCLUDED.ip_address,
    		updated_at = NOW(),
		created_at = NOW(),
    		is_revoked = FALSE
	RETURNING id;
	`

	var id string
	row := s.db.QueryRow(ctx, query, userID, token, userAgent, ip)
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Revoke(
	ctx context.Context,
	userID, userAgent string,
) error {
	const op = "storage.tokens.Revoke"

	query := `
	UPDATE refresh_tokens
	SET is_revoked = TRUE, updated_at = NOW()
	WHERE user_id = $1 AND user_agent = $2 AND is_revoked = FALSE;
	`

	res, err := s.db.Exec(ctx, query, userID, userAgent)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repoErr.ErrRefreshTokenNotFound)
	}

	return nil
}
