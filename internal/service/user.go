package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailDuplicated = repository.ErrEmailDuplicated
	ErrEmailOrPassword = errors.New("邮箱或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}

}
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, ErrEmailOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrEmailOrPassword
	}
	return domain.User{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = svc.repo.Create(ctx, user.Email, string(encrypted))
	if err == repository.ErrEmailDuplicated {
		return ErrEmailDuplicated
	}
	if err != nil {
		return err
	}
	return nil
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	panic("todo!")
}

func (svc *UserService) Profile(ctx context.Context, u domain.User) error {
	panic("todo!")
}
