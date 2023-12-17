package repository

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrEmailDuplicated = dao.ErrEmailDuplicated
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (u *UserRepository) Create(ctx context.Context, email string, password string) error {
	err := u.dao.Create(ctx, email, password)
	if errors.Is(err, ErrEmailDuplicated) {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.cache.Get(ctx, int(id))
	if err == nil {
		return user, err
	} else {
		user, err = u.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		_ = u.cache.Set(ctx, user)
	}
	return user, nil
}
