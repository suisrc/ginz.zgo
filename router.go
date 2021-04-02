package ginz

// 系统默认根路由
import (
	"github.com/suisrc/config.zgo"

	"github.com/gin-gonic/gin"
)

// Router 根路由 Register???
type Router gin.IRouter

// NewRouter 初始化根路由
func NewRouter(app *gin.Engine) Router {
	var router Router
	if v := config.C.HTTP.ContextPath; v != "" {
		router = app.Group(v)
	} else {
		router = app
	}

	return router
}
