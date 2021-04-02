package ginz

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/suisrc/res.zgo"
)

type ReqContext interface {
	res.Context

	ShouldBindWith(obj interface{}, b binding.Binding) error // bind
}

// ParseJSON 解析请求JSON, 注意,解析失败后需要直接返回
func ParseJSON(c ReqContext, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.JSON); err != nil {
		return res.Wrap400Response(c, err)
	}
	return nil
}

// BindXML 解析请求JSON, 注意,解析失败后需要直接返回
func BindXML(c ReqContext, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.XML); err != nil {
		return res.Wrap400Response(c, err)
	}
	return nil
}

// ParseQuery 解析Query参数, 注意,解析失败后需要直接返回
func ParseQuery(c ReqContext, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Query); err != nil {
		return res.Wrap400Response(c, err)
	}
	return nil
}

// ParseForm 解析Form请求, 注意,解析失败后需要直接返回
func ParseForm(c ReqContext, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return res.Wrap400Response(c, err)
	}
	return nil
}
