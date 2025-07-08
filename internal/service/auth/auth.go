package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	svcErr "github.com/passwordhash/jwt-test-task/internal/service/errors"
	"github.com/passwordhash/jwt-test-task/pkg/jwt"
)

type Service struct {
	log *slog.Logger

	accessTTL time.Duration
	secret    string
}

func NewService(
	log *slog.Logger,
	accessTTL time.Duration,
	secret string,
) *Service {
	return &Service{
		log:       log,
		accessTTL: accessTTL,
		secret:    secret,
	}
}

func (s *Service) GetPair(ctx context.Context, id, ip, userAgent string) (access, refresh string, err error) {
	const op = "auth.service.GetPair"

	log := s.log.With("op", op, "id", id, "ip", ip, "userAgent", userAgent)

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid id", slog.String("id", id), slog.Any("error", err))

		return "", "", svcErr.ErrInvalidID
	}

	access, err = jwt.NewToken("HS512", id, s.accessTTL, s.secret)
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err))

		return "", "", err
	}

	return access, "bar", nil
}
