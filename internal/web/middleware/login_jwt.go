package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
	paths map[string]bool
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		paths: make(map[string]bool, 4),
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(paths ...string) *LoginJWTMiddlewareBuilder {
	for _, v := range paths {
		l.paths[v] = true
	}
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := l.paths[ctx.Request.URL.String()]; ok {
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		seg := strings.Split(authCode, " ")
		if len(seg) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := seg[1]
		uc := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte("xbcbtlzWUNZfXzmXvnLdpQnoIFRegaUK"), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 增强系统安全性，当 token 不是同一设备使用时 认为没有登录
		if uc.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每 10s 刷新一次
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			newToken, err := token.SignedString([]byte("xbcbtlzWUNZfXzmXvnLdpQnoIFRegaUK"))
			if err != nil {
				fmt.Printf("%v", err)
			} else {
				print("刷新jwt")
				ctx.Header("x-jwt-token", newToken)
			}
		}
		ctx.Set("userID", uc.Uid)
	}
}
