package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotExsit   = dao.ErrUserNotExsit
	ErrDetailNotExist = dao.ErrDetailNotExist
)

type UserDetailRepository struct {
	dao *dao.UserDetailDao
}

func NewUserDetailRepository(dao *dao.UserDetailDao) *UserDetailRepository {
	return &UserDetailRepository{
		dao: dao,
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
	detail, err := u.dao.FindById(ctx, id)
	if err != nil {
		if err == ErrDetailNotExist {
			return domain.User{}, ErrDetailNotExist
		}
		return domain.User{}, err
	}
	return domain.User{
		ID:           detail.UserID,
		Avatar:       detail.Avatar,
		Introduction: detail.Introduction,
		Birthday:     detail.Birthday,
		NickName:     detail.NickName,
		Email:        detail.User.Email,
		Password:     detail.User.Password,
	}, nil
}

func (u *UserDetailRepository) UpdateById(ctx context.Context, user domain.User) error {
	err := u.dao.UpdateById(ctx, user.ID, user.NickName, user.Introduction, user.Birthday, user.Avatar)
	if err != nil {
		return err
	}
	return err
}
