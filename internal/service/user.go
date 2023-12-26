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
	ErrDetailNotExist  = repository.ErrDetailNotExist
)

type UserService interface {
	Login(ctx context.Context, email string, password string) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}
type userService struct {
	repo       repository.UserRepository
	detailRepo repository.UserDetailRepository
}

func NewUserService(userRepo repository.UserRepository, userDetailRepo repository.UserDetailRepository) UserService {
	return &userService{
		repo:       userRepo,
		detailRepo: userDetailRepo,
	}

}
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, ErrEmailOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrEmailOrPassword
	}
	return domain.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = svc.repo.Create(ctx, user.Email, string(encrypted))
	if errors.Is(err, repository.ErrEmailDuplicated) {
		return ErrEmailDuplicated
	}
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) Edit(ctx context.Context, u domain.User) error {
	_, err := svc.detailRepo.FindById(ctx, u.ID)
	if errors.Is(err, ErrDetailNotExist) {
		err = svc.detailRepo.Create(ctx, u)
		return err
	}
	err = svc.detailRepo.UpdateById(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.detailRepo.FindById(ctx, id)
	if errors.Is(err, ErrDetailNotExist) {
		user, err = svc.repo.FindById(ctx, id)
		if err == nil {
			_ = svc.detailRepo.Create(ctx, user)
		}
		return user, err
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}
