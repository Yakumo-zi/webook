package main

import (
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := InitServer()
	db := InitDatabase()
	dao.InitTables(db)
	u := InitUser(db)

	u.RegisterRoutes(server)
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitServer() *gin.Engine {
	server := gin.Default()
	// 使用中间允许跨域
	server.Use(cors.New(cors.Config{
		// 不主动设置表示允许简单方法
		// AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	}))
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.
		NewLoginMiddlewareBuilder().
		IgnorePath("/users/login", "/users/signup").
		Build())
	return server
}
func InitDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}
	return db.Debug()
}
func InitUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
