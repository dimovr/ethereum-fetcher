package auth

import (
	"ethereum_fetcher/api"
	"ethereum_fetcher/pkg/logging"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(req api.AuthRequest) (*string, error)
	GetUserId(tokenString api.AuthToken) (uint64, error)
}

type impl struct {
	repo   *UserRepo
	jm     JwtManager
	logger *logrus.Logger
}

func NewAuthService(db *gorm.DB, jwtSecretKey string) (AuthService, error) {
	logger := logging.New()

	repo := &UserRepo{db: db}
	JwtManager := NewJwtManager(jwtSecretKey)

	return &impl{
		repo:   repo,
		jm:     JwtManager,
		logger: logger,
	}, nil
}

func (s *impl) Authenticate(req api.AuthRequest) (*string, error) {
	user, err := s.repo.FindUser(req.Username, req.Password)
	if err != nil {
		s.logger.Debugf("failed to authenticate user '%s' : %s", req.Username, err)
		return nil, err
	}

	jwt, err := s.jm.GenerateJwt(user.ID)
	if err != nil {
		s.logger.Warnf("failed to generate JWT for user '%s' : %s", req.Username, err)
		return nil, err
	}

	return &jwt, nil
}

func (s *impl) GetUserId(authToken api.AuthToken) (uint64, error) {
	token, err := s.jm.ParseJwt(authToken)
	if err != nil || !token.Valid {
		s.logger.Warnf("failed to parse JWT : %s", err)
		return 0, InvalidToken
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		s.logger.Warnf("failed to get subject from token : %s", err)
		return 0, NoSubject
	}
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		s.logger.Warnf("failed to convert subject to user id : %s", err)
		return 0, NoSubject
	}

	return id, nil
}
