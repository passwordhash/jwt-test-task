package tokens

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
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
		RETURNING id
	`

	var id string
	row := s.db.QueryRow(ctx, query, userID, token, userAgent, ip)
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
