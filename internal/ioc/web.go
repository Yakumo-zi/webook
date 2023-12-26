package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/ratelimit"
)

func InitGin(mdls []gin.HandlerFunc, userHandler *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHandler.RegisterRoutes(server)
	return server
}
func InitMiddleware(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
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
		}),
		middleware.NewLoginJWTMiddlewareBuilder().IgnorePath("/users/login", "/users/signup").Build(),
		ratelimit.NewBuilder(cmd, time.Minute, 100).Build(),
	}
}
