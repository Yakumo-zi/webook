package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths map[string]bool
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		paths: make(map[string]bool, 4),
	}
}

func (l *LoginMiddlewareBuilder) IgnorePath(paths ...string) *LoginMiddlewareBuilder {
	for _, v := range paths {
		l.paths[v] = true
	}
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if _, ok := l.paths[ctx.Request.URL.String()]; ok {
			return
		}
		id := session.Get("userID")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每 10s 刷新一次 session，当首次登录时也会刷新 session
		updateTime := session.Get("updateTime")
		now := time.Now()
		if updateTime == nil || now.Second()-updateTime.(time.Time).Second() > 10 {
			println("刷新 session")
			session.Set("updateTime", now)
			session.Options(sessions.Options{
				MaxAge: 60,
			})
			session.Save()
		}
	}
}
