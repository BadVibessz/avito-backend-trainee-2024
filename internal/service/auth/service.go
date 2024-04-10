package auth

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	CheckUniqueConstraints(ctx context.Context, username string) error
}

type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type Service struct {
	UserRepo UserRepo
	Hasher   Hasher
}

func New(userRepo UserRepo, hasher Hasher) *Service {
	return &Service{
		UserRepo: userRepo,
		Hasher:   hasher,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// ensure that user with this email and username does not exist
	err := s.UserRepo.CheckUniqueConstraints(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	// user model sent with plain password
	hash, err := s.Hasher.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = string(hash)

	return s.UserRepo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = s.Hasher.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
