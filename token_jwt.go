package ginz

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/auth.zgo"
)

// GetBearerToken 获取用户令牌
func GetBearerToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		prefix := "Bearer "
		if auth := c.GetHeader("Authorization"); auth != "" && strings.HasPrefix(auth, prefix) {
			return auth[len(prefix):], nil
		}
	}

	return "", auth.ErrNoneToken
}

// GetQueryToken 获取用户令牌
func GetQueryToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if auth, ok := c.GetQuery("token"); ok && auth != "" {
			return auth, nil
		}
	}

	return "", auth.ErrNoneToken
}

// GetCookieToken 获取用户令牌
func GetCookieToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if auth, err := c.Cookie("authorization"); err == nil && auth != "" {
			return auth, nil
		}
	}

	return "", auth.ErrNoneToken
}
