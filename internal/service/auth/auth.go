package auth

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	svcErr "github.com/passwordhash/jwt-test-task/internal/service/errors"
)

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) GetPair(ctx context.Context, id, ip, userAgent string) (access, refresh string, err error) {
	const op = "auth.service.GetPair"

	log := s.log.With("op", op, "id", id, "ip", ip, "userAgent", userAgent)

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid id", slog.String("id", id), slog.Any("error", err))

		return "", "", svcErr.ErrInvalidID
	}

	return "foo", "bar", nil
}
