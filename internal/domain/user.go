package domain

import "time"

type User struct {
	ID           int64
	Email        string
	Password     string
	NickName     string
	Avatar       string
	Introduction string
	Birthday     string
	Ctime        time.Time
}
