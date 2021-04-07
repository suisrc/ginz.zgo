package ginz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/auth.zgo"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/logger.zgo"
	"github.com/suisrc/res.zgo"
)

type Context interface {
	context.Context
	res.Context
	res.ReqContext

	ReqContext
	ResContext

	logger.ContextTrace

	// Set(key string, value interface{})
	// Get(key string)
}

var _ Context = &ginContext{}

type ginContext struct {
	*gin.Context
}

func NewContext(c *gin.Context) Context {
	return &ginContext{Context: c}
}

//==========================================

// GetTraceID ...
func (a *ginContext) GetTraceID() string {
	return GetTraceID(a.Context)
}

// GetTraceCIP ...
func (a *ginContext) GetTraceCIP() string {
	return GetClientIP(a.Context)
}

// GetTraceUID ...
func (a *ginContext) GetTraceUID() string {
	if usr, ok := GetUserInfo(a.Context); ok {
		return fmt.Sprintf("[%s]->%s", usr.GetAccount1(), usr.GetUserID())
	}
	return ""
}

// FormatMessage ...
func (a *ginContext) FormatMessage(emsg *i18n.Message, args map[string]interface{}) string {
	return i18n.FormatMessage(a.Context, emsg, args)
}

// GetRequest ...
func (a *ginContext) GetRequest() *http.Request {
	return a.Context.Request
}

// GetUserInfo ...
func (a *ginContext) GetUserInfo() auth.UserInfo {
	if usr, ok := GetUserInfo(a.Context); ok {
		return usr
	}
	return nil
}
