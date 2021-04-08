package ginz

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/suisrc/config.zgo"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/res.zgo"
)

// HandlerFunc -> res.HandlerFunc -> gin.HandlerFunc
func HandlerFunc(fx res.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		fx(NewContext(c))
	}
}

// HandlerFunc2 -> res.HandlerFunc -> gin.HandlerFunc
func HandlerFunc2(fx res.HandlerFunc, nc func(*gin.Context) res.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		fx(nc(c))
	}
}

// NoMethodHandler 未找到请求方法的处理函数
func NoMethodHandler(c *gin.Context) {
	ResError(NewContext(c), res.Err405MethodNotAllowed)
	// Abort, 终止
}

// NoRouteHandler 未找到请求路由的处理函数
func NoRouteHandler(c *gin.Context) {
	ResError(NewContext(c), res.Err404NotFound)
	// Abort, 终止
}

//===================================================================
//===================================================================
//===================================================================

// GizMiddleware Giz, 主要部署前端时候(www中间件)对静态资源进行压缩
func GizMiddleware() gin.HandlerFunc {
	conf := config.C.GZIP
	return gzip.Gzip(gzip.BestCompression,
		gzip.WithExcludedExtensions(conf.ExcludedExtentions),
		gzip.WithExcludedPaths(conf.ExcludedPaths),
	)
}

// CORSMiddleware 跨域
func CORSMiddleware() gin.HandlerFunc {
	conf := config.C.CORS
	return cors.New(cors.Config{
		AllowOrigins:     conf.AllowOrigins,
		AllowMethods:     conf.AllowMethods,
		AllowHeaders:     conf.AllowHeaders,
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           time.Second * time.Duration(conf.MaxAge),
	})
}

//===================================================================
//===================================================================
//===================================================================

// I18nMiddleware 国际化
func I18nMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	// bundle := i18n.NewBundle(language.Chinese)
	// bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// bundle.LoadMessageFile("locales/active.zh-CN.toml")
	// bundle.LoadMessageFile("locales/active.en-US.toml")
	return i18n.Serve(bundle)
}

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, skippers ...res.SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := NewContext(c)
		if res.SkipHandler(cc, skippers...) {
			c.Next()
			return
		}
		if root == "" {
			root = config.C.WWW.RootDir
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, config.C.WWW.Index)
		}

		c.File(fpath)
		c.Abort()
	}
}

//===================================================================
//===================================================================
//===================================================================

// CORSMiddleware2 跨域
// 不推荐使用,可以使用gin中的跨域处理
func CORSMiddleware2(skippers ...res.SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := NewContext(c)
		if res.SkipHandler(cc, skippers...) {
			c.Next()
			return
		}

		origin := c.GetHeader("Origin") //请求头部
		if origin != "" {
			//c.Writer.Header().Set("Access-Control-Allow-Origin","*")
			// 这是允许访问所有域
			c.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
			// header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Length,X-CSRF-Token,Token,session,X_Requested_With,Accept,Origin,Host,Connection,Accept-Encoding,Accept-Language,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Pragma")
			// 允许跨域设置 | 可以返回其他子段
			// 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			// 缓存请求信息 单位为秒
			c.Header("Access-Control-Max-Age", "172800")
			// 跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Credentials", "false")
			// 设置返回格式是json
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()
	}
}
