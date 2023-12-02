package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrEmailDuplicated = errors.New("邮箱重复")
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}
func (u *UserDao) Create(ctx context.Context, email string, password string) error {
	now := time.Now().UnixMilli()
	err := u.db.WithContext(ctx).Create(&User{
		Email:    email,
		Password: password,
		CTime:    now,
		UTime:    now,
	}).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const errUniqueConflict = 1062
		if me.Number == errUniqueConflict {
			return ErrEmailDuplicated
		}
	}
	if err != nil {
		return err
	}
	return nil
}
func (u *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	CTime    int64
	UTime    int64
}

type UserDetail struct {
	ID           int64 `gorm:"primaryKey"`
	NickName     string
	Avatar       string
	Introduction string
	Birthday     int64
}