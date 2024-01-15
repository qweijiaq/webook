package service

import (
	"context"

	"github.com/qweijiaq/webook/internal/domain"
	"github.com/qweijiaq/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里
	// 存储
	return svc.repo.Create(ctx, u)
}
