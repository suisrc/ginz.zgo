package ginz

import (
	"net/http"

	"github.com/suisrc/logger.zgo"
	"github.com/suisrc/res.zgo"
)

// ResSuccess 包装响应错误
// 禁止service层调用,请使用NewSuccess替换
func ResSuccess(ctx Context, v interface{}) error {
	res := res.NewSuccess(ctx, v)
	//ctx.JSON(http.StatusOK, res)
	//ctx.Abort()
	ResJSON(ctx, http.StatusOK, res)
	return res
}

// ResError 包装响应错误
// 禁止service层调用,请使用NewWarpError替换
func ResError(ctx Context, em *res.ErrorModel) error {
	res := res.NewWrapError(ctx, em)
	//ctx.JSON(http.StatusOK, res)
	//ctx.Abort()
	ResJSON(ctx, em.Status, res)
	return res
}

// ResJSON 响应JSON数据
// 禁止service层调用
func ResJSON(ctx Context, status int, v interface{}) {
	if ctx == nil {
		return
	}
	buf, err := res.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	if status == 0 {
		status = http.StatusOK
	}
	ctx.Data(status, res.ResponseTypeJSON, buf)
	ctx.Abort()

	// ctx.JSON(status, v)
	// ctx.PureJSON(status, v)
}

// FixResponseError 上级应用已经处理了返回值
func FixResponseError(c Context, err error) bool {
	switch err.(type) {
	case *res.Success, *res.ErrorInfo:
		ResJSON(c, http.StatusOK, err)
		return true
	case *res.ErrorRedirect:
		status := err.(*res.ErrorRedirect).Status
		if status <= 0 {
			status = http.StatusSeeOther
		}
		c.Redirect(status, err.(*res.ErrorRedirect).Location)
		return true
	case *res.ErrorNone:
		// do nothing
		return true
	case *res.ErrorModel:
		em := err.(*res.ErrorModel)
		ResJSON(c, em.Status, res.NewWrapError(c, em))
		return true
	default:
		// e := err.Error()
		return false
	}
}

// FixResponse400Error 修复返回的异常
func FixResponse400Error(c Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, res.Err400BadRequest)
}

// FixResponse401Error 修复返回的异常, 注意, 401异常会导致系统重定向到登陆页面
func FixResponse401Error(c Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, res.Err401Unauthorized)
}

// FixResponse403Error 修复返回的异常
func FixResponse403Error(c Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, res.Err403Forbidden)
}

// FixResponse406Error 修复返回的异常
func FixResponse406Error(c Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, res.Err406NotAcceptable)
}

// FixResponse500Error 修复返回的异常
func FixResponse500Error(c Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, res.Err500InternalServer)
}

// FixResponse500Error2Logger 修复返回的异常
func FixResponse500Error2Logger(c Context, err error) {
	FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) })
}

//=============================================================
//=============================================================
//=============================================================

//// ResErrorResBody 包装响应错误
//// 禁止service层调用
//func ResErrorResBody(ctx Context, em *res.ErrorModel) error {
//	res := res.NewWrapError(ctx, em)
//	ResJSONResBody(ctx, em.Status, res)
//	return res
//}
//
//// ResJSONResBody 响应JSON数据
//// 禁止service层调用
//func ResJSONResBody(ctx Context, status int, v interface{}) {
//	if ctx == nil {
//		return
//	}
//	buf, err := res.JSONMarshal(v)
//	if err != nil {
//		panic(err)
//	}
//	ctx.Set(res.ResBodyKey, buf)
//	if status == 0 {
//		status = http.StatusOK
//	}
//	ctx.Data(status, res.ResponseTypeJSON, buf)
//	ctx.Abort()
//
//	// ctx.JSON(status, v)
//	// ctx.PureJSON(status, v)
//}
