package tokens

import (
	"context"
	"fmt"

	repoErr "github.com/passwordhash/jwt-test-task/internal/storage/errors"
	"github.com/passwordhash/jwt-test-task/pkg/postgres"
)

type Storage struct {
	db postgres.DB
}

func New(db postgres.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Save(ctx context.Context, userID, tokenID, tokenHash, userAgent, ip string) (string, error) {
	const op = "storage.tokens.Save"

	query := `
	INSERT INTO refresh_tokens (user_id, token_id, token_hash, user_agent, ip_address)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (user_id, user_agent) 
	DO UPDATE SET 
    		token_hash = EXCLUDED.token_hash,
		ip_address = EXCLUDED.ip_address,
    		updated_at = NOW(),
		created_at = NOW(),
    		is_revoked = FALSE
	RETURNING id;
	`

	var id string
	row := s.db.QueryRow(ctx, query, userID, tokenID, tokenHash, userAgent, ip)
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
