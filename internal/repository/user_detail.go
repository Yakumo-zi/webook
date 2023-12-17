package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotExsit   = dao.ErrUserNotExsit
	ErrDetailNotExist = dao.ErrDetailNotExist
)

type UserDetailRepository struct {
	dao   *dao.UserDetailDao
	cache *cache.UserCache
}

func NewUserDetailRepository(dao *dao.UserDetailDao, cache *cache.UserCache) *UserDetailRepository {
	return &UserDetailRepository{
		dao:   dao,
		cache: cache,
	}
}

func (u *UserDetailRepository) Create(ctx context.Context, ud domain.User) error {
	err := u.dao.Create(ctx, ud.ID, ud.NickName, ud.Introduction, ud.Birthday)
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			if me.Number == 1452 {
				return ErrUserNotExsit
			}
		}
		return err
	}
	return nil
}

func (u *UserDetailRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	detail, err := u.cache.Get(ctx, int(id))
	if err == nil {
		return detail, err
	} else {
		detail, err = u.dao.FindById(ctx, id)
		if err == ErrDetailNotExist {
			return domain.User{}, ErrDetailNotExist
		}
		if err != nil {
			return domain.User{}, err
		}
		_ = u.cache.Set(ctx, detail)
	}

	return detail, nil
}

func (u *UserDetailRepository) UpdateById(ctx context.Context, user domain.User) error {
	err := u.dao.UpdateById(ctx, user.ID, user.NickName, user.Introduction, user.Birthday, user.Avatar)
	if err != nil {
		return err
	}
	_ = u.cache.Set(ctx, user)
	return err
}
