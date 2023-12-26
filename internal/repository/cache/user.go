package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

type UserCache interface {
	Get(ctx context.Context, id int) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
	Del(ctx context.Context, id int) error
}

type userCache struct {
	cmd redis.Cmdable
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &userCache{
		cmd: cmd,
	}
}
func (u *userCache) generateKey(id int) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (u *userCache) Get(ctx context.Context, id int) (domain.User, error) {
	var user domain.User
	res := u.cmd.Get(ctx, u.generateKey(id))
	data, err := res.Result()
	if err != nil {
		return domain.User{}, err
	}
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}
func (u *userCache) Set(ctx context.Context, du domain.User) error {
	ru, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return u.cmd.Set(ctx, u.generateKey(int(du.ID)), ru, time.Hour*24).Err()
}

func (u *userCache) Del(ctx context.Context, id int) error {
	return u.cmd.Del(ctx, u.generateKey(id)).Err()
}
