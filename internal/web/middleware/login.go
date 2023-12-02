package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if _, ok := l.paths[ctx.Request.URL.String()]; ok {
			return
		}
		id := session.Get("id")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
