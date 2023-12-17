package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/ratelimit"
)

func main() {

	db, cmd := InitDatabase()
	server := InitServer(cmd)
	dao.InitTables(db)
	u := InitUser(db, cmd)

	u.RegisterRoutes(server)
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitServer(cmd redis.Cmdable) *gin.Engine {
	server := gin.Default()
	// 使用中间允许跨域
	server.Use(cors.New(cors.Config{
		// 不主动设置表示允许简单方法
		// AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 设置允许前端读那些header,没有在这里面声明的 header 前端是无法拿到的
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	}))
	// 使用 cookie 存放 session
	//store := cookie.NewStore([]byte("secret"))
	// 使用 redis 存放 session
	//store, _ := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("xbcbtlzWUNZfXzmXvnLdpQnoIFRegaUK"), []byte("asUfpYzKvPCEDJFEPWPTcXoaaFhMVKMy"))
	//server.Use(sessions.Sessions("mysession", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/login", "/users/signup").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePath("/users/login", "/users/signup").Build())

	// 接入限流中间件
	server.Use(ratelimit.NewBuilder(cmd, time.Minute, 100).Build())
	return server
}
func InitDatabase() (*gorm.DB, redis.Cmdable) {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	return db.Debug(), cmd
}
func InitUser(db *gorm.DB, cmd redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDao(db)
	udd := dao.NewUserDetailDao(db)
	userCache := cache.NewUserCache(cmd)
	userRepo := repository.NewUserRepository(ud, userCache)
	userDetailRepo := repository.NewUserDetailRepository(udd, userCache)
	svc := service.NewUserService(userRepo, userDetailRepo)
	u := web.NewUserHandler(svc)
	return u
}
