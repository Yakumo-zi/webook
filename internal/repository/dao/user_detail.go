package dao

import (
	"context"
	"errors"
	"time"
	"webook/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrUserNotExsit   = errors.New("用户不存在")
	ErrDetailNotExist = errors.New("详情不存在")
)

type UserDetailDao struct {
	db *gorm.DB
}

func NewUserDetailDao(db *gorm.DB) *UserDetailDao {
	return &UserDetailDao{
		db: db,
	}
}

func (u *UserDetailDao) Create(ctx context.Context, id int64, nick string, intro string, birthday string) error {
	now := time.Now().UnixMilli()
	err := u.db.WithContext(ctx).Create(&UserDetail{
		UserID:       id,
		NickName:     nick,
		Introduction: intro,
		Birthday:     birthday,
		CTime:        now,
		UTime:        now,
	}).Error
	if err != nil {
		return ErrUserNotExsit
	}
	return nil

}
func (u *UserDetailDao) FindById(ctx context.Context, id int64) (domain.User, error) {
	var ud UserDetail
	row := u.db.WithContext(ctx).Where("user_id = ?", id).Preload("User").Find(&ud).RowsAffected
	if row == 0 {
		return domain.User{}, ErrDetailNotExist
	}
	return domain.User{
		ID:           ud.UserID,
		NickName:     ud.NickName,
		Avatar:       ud.Avatar,
		Introduction: ud.Introduction,
		Birthday:     ud.Birthday,
		Email:        ud.User.Email,
		Password:     ud.User.Password,
	}, nil
}

func (u *UserDetailDao) UpdateById(ctx context.Context, id int64, nick string, intro string, birthday string, avatar string) error {
	now := time.Now().UnixMilli()
	detail := UserDetail{
		NickName:     nick,
		UTime:        now,
		Introduction: intro,
		Birthday:     birthday,
		Avatar:       avatar,
	}
	err := u.db.WithContext(ctx).Model(&UserDetail{}).Where("user_id = ?", id).Updates(detail).Error
	if err != nil {
		return err
	}
	return nil
}

type UserDetail struct {
	ID           int64 `gorm:"primaryKey"`
	NickName     string
	Avatar       string
	Introduction string
	Birthday     string
	CTime        int64
	UTime        int64
	UserID       int64 `gorm:"unique"`
	User         User
}
