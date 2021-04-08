package ginz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/suisrc/auth.zgo"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/logger.zgo"
	"github.com/suisrc/res.zgo"
)

type Context interface {
	logger.ContextTrace
	context.Context
	res.Context

	GetUserInfo() auth.UserInfo                              // 获取登陆用户信息
	ShouldBindWith(obj interface{}, b binding.Binding) error // bind
	Data(int, string, []byte)                                // 写入数据
	Redirect(int, string)                                    // 重定向RX

}

var _ Context = &GinContext{}

type GinContext struct {
	*gin.Context
}

func NewContext(c *gin.Context) Context {
	return &GinContext{Context: c}
}

//==========================================

// GetTraceID ...
func (a *GinContext) GetTraceID() string {
	return GetTraceID(a.Context)
}

// GetTraceCIP ...
func (a *GinContext) GetTraceCIP() string {
	return GetClientIP(a.Context)
}

// GetTraceUID ...
func (a *GinContext) GetTraceUID() string {
	if usr, ok := GetUserInfo(a.Context); ok {
		return fmt.Sprintf("[%s]->%s", usr.GetAccount1(), usr.GetUserID())
	}
	return ""
}

// FormatMessage ...
func (a *GinContext) FormatMessage(emsg *i18n.Message, args map[string]interface{}) string {
	return i18n.FormatMessage(a.Context, emsg, args)
}

// GetRequest ...
func (a *GinContext) GetRequest() *http.Request {
	return a.Context.Request
}

// GetUserInfo ...
func (a *GinContext) GetUserInfo() auth.UserInfo {
	if usr, ok := GetUserInfo(a.Context); ok {
		return usr
	}
	return nil
}
