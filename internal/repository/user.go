package repository

import (
	"context"

	"github.com/qweijiaq/webook/internal/domain"
	"github.com/qweijiaq/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindById(int642 int64) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了会写 cache
}
