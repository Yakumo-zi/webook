package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}
	return db
}
