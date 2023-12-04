package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
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
	//server := gin.Default()
	//server.GET("/", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello,k8s")
	//})
	//server.Run(":8180")
}

func InitServer() *gin.Engine {
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
	udd := dao.NewUserDetailDao(db)
	userRepo := repository.NewUserRepository(ud)
	userDetailRepo := repository.NewUserDetailRepository(udd)
	svc := service.NewUserService(userRepo, userDetailRepo)
	u := web.NewUserHandler(svc)
	return u
}
