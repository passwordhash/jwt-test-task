package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"

	svcErr "github.com/passwordhash/jwt-test-task/internal/service/errors"
	"github.com/passwordhash/jwt-test-task/pkg/jwt"
)

const (
	refreshTokenLength = 32
)

type RefreshTokenSaver interface {
	Save(ctx context.Context, userID, token, userAgent, ip string) (string, error)
}

type RefreshTokenGenerator interface {
	Generate(length int) (string, error)
	Hash(token string) (string, error)
}

type Service struct {
	log *slog.Logger

	refreshSaver          RefreshTokenSaver
	refreshTokenGenerator RefreshTokenGenerator

	accessTTL time.Duration
	secret    string
}

func New(
	log *slog.Logger,
	refreshSaver RefreshTokenSaver,
	refreshTokenGenerator RefreshTokenGenerator,
	accessTTL time.Duration,
	secret string,
) *Service {
	return &Service{
		log:                   log,
		refreshSaver:          refreshSaver,
		refreshTokenGenerator: refreshTokenGenerator,
		accessTTL:             accessTTL,
		secret:                secret,
	}
}

func (s *Service) GetPair(ctx context.Context, id, remoteAddr, userAgent string) (access, refresh string, err error) {
	const op = "tokens.service.GetPair"

	log := s.log.With("op", op, "id", id, "remoteAddr", remoteAddr, "userAgent", userAgent)

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid id", slog.String("id", id), slog.Any("error", err))

		return "", "", svcErr.ErrInvalidID
	}

	access, err = jwt.NewToken("HS512", id, s.accessTTL, s.secret)
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err))

		return "", "", err
	}

	refresh, err = s.refreshTokenGenerator.Generate(refreshTokenLength)
	if err != nil {
		log.Error("failed to generate refresh token", slog.Any("error", err))

		return "", "", err
	}

	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Error("failed to split host and port from IP", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.refreshSaver.Save(ctx, id, refresh, userAgent, ip)
	// TODO: handle error properly
	if err != nil {
		log.Error("failed to save refresh token", slog.Any("error", err))

		return "", "", err
	}

	return access, refresh, nil
}
