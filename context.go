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

var _ Context = &GinContext{}

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) *GinContext {
	return &GinContext{Context: c}
}

//==========================================

func (a *GinContext) FormatMessage(emsg *i18n.Message, args map[string]interface{}) string {
	return i18n.FormatMessage(a.Context, emsg, args)
}
func (a *GinContext) GetRequest() *http.Request {
	return a.Context.Request
}
func (a *GinContext) GetTraceID() string {
	return GetTraceID(a.Context)
}
func (a *GinContext) GetTraceCIP() string {
	return GetClientIP(a.Context)
}
func (a *GinContext) GetTraceUID() string {
	if usr, ok := GetUserInfo(a.Context); ok {
		return fmt.Sprintf("[%s]->%s", usr.GetAccount1(), usr.GetUserID())
	}
	return ""
}
func (a *GinContext) GetUserInfo() auth.UserInfo {
	if usr, ok := GetUserInfo(a.Context); ok {
		return usr
	}
	return nil
}