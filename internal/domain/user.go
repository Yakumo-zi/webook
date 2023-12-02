package domain

import "time"

type User struct {
	Id           int64
	Email        string
	Password     string
	NickName     string
	Avatar       string
	Introduction string
	Birthday     time.Time
	Ctime        time.Time
}
