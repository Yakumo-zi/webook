package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrEmailDuplicated = dao.ErrEmailDuplicated
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (u *UserRepository) Create(ctx context.Context, email string, password string) error {
	err := u.dao.Create(ctx, email, password)
	if err == ErrEmailDuplicated {
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
	user, err := u.dao.FindById(ctx, id)
	if err != nil {
		if err == ErrDetailNotExist {
			return domain.User{}, ErrDetailNotExist
		}
		return domain.User{}, err
	}
	return domain.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
