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

type UserRepository interface {
	Create(ctx context.Context, email string, password string) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type UserRepositoryI struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &UserRepositoryI{
		dao:   dao,
		cache: cache,
	}
}
func (u *UserRepositoryI) Create(ctx context.Context, email string, password string) error {
	err := u.dao.Create(ctx, email, password)
	if errors.Is(err, ErrEmailDuplicated) {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepositoryI) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepositoryI) FindById(ctx context.Context, id int64) (domain.User, error) {
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
