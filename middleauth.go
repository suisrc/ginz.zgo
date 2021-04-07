package ginz

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/suisrc/auth.zgo"
	"github.com/suisrc/config.zgo"
	"github.com/suisrc/res.zgo"
)

// package middleware

// UseAuthServerMiddleware 用户授权中间件
func UseAuthServerMiddleware(c *gin.Context, gf func(*gin.Context, auth.UserInfo) (aid, uid int64, err error)) {
	// 如果通过验证， 当前用户是一定登录的
	user, _ := GetUserInfo(c)

	h := c.Writer.Header()
	h.Set("X-Request-Z-Token-Kid", user.GetTokenID())
	h.Set("X-Request-Z-Token-Pid", user.GetTokenPID())

	if user.GetAccount() != "" && gf != nil {
		if acc, usr, err := gf(c, user); err == nil {
			h.Set("X-Request-Z-Account-Id", strconv.Itoa(int(acc)))
			if usr > 0 {
				h.Set("X-Request-Z-User-Id", strconv.Itoa(int(usr)))
			}
		}
	}

	h.Set("X-Request-Z-Account", user.GetAccount())   // 账户名
	h.Set("X-Request-Z-Account1", user.GetAccount1()) // 账户名1
	h.Set("X-Request-Z-Account2", user.GetAccount2()) // 账户名2
	h.Set("X-Request-Z-User-Kid", user.GetUserID())   // 用户唯一标识
	h.Set("X-Request-Z-User-Name", url.QueryEscape(user.GetUserName()))
	h.Set("X-Request-Z-User-Roles", strings.Join(user.GetUserRoles(), ";"))
	h.Set("X-Request-Z-Org-Code", user.GetOrgCode())
	h.Set("X-Request-Z-Org-Admin", user.GetOrgAdmin())
	h.Set("X-Request-Z-Org-Usrid", user.GetOrgUsrID())
	h.Set("X-Request-Z-Agent", user.GetAgent())
	h.Set("X-Request-Z-Scope", user.GetScope())
	h.Set("X-Request-Z-Domain", user.GetDomain())
	h.Set("X-Request-Z-Issuer", user.GetIssuer())
	h.Set("X-Request-Z-Audience", user.GetAudience())
	h.Set("X-Request-Z-Xip", GetHostIP(c))
}

// UseAuthzClientMiddleware 用户授权中间件, 只判定登录权限
func UseAuthzClientMiddleware(c *gin.Context) {
	conf := config.C.JWTAuth

	user := SigninUser{}
	if conf.AuthzServer == "" {
		c.BindHeader(user) // 执行验证
	} else if v := c.GetHeader("X-Request-Z-Xip"); v != "" {
		// XXX 这里需要验证请求的合法性才能进行下面的绑定操作
		c.BindHeader(user) // 执行验证
	} else if ok := UseRemoteAuth(c, &user, conf.AuthzServer); !ok {
		return // 无法获取认证， 结束处理
	}
	if user.TokenID == "" {
		// 令牌为空， 拒绝访问
		ResError(NewGinContext(c), res.Err401Unauthorized)
		return // 无有效登陆用户
	}
	// 为后端服务器提供服务
	SetUserInfo(c, &user)
}

// UseAuthxClientMiddleware 用户授权中间件, 只判定登录权限
func UseAuthxClientMiddleware(c *gin.Context) {
	conf := config.C.JWTAuth

	user := SigninUser{}
	if conf.AuthxServer == "" {
		c.BindHeader(user) // 执行验证
	} else if v := c.GetHeader("X-Request-Z-Xip"); v != "" {
		// XXX 这里需要验证请求的合法性才能进行下面的绑定操作
		c.BindHeader(user) // 执行验证
	} else if ok := UseRemoteAuth(c, &user, conf.AuthxServer); !ok {
		return // 无法获取认证， 结束处理
	}
	if user.TokenID == "" {
		// 令牌为空， 拒绝访问
		ResError(NewGinContext(c), res.Err401Unauthorized)
		return // 无有效登陆用户
	}
	// 为后端服务器提供服务
	SetUserInfo(c, &user)
}

// UseRemoteAuthz ...
// proxy_set_header X-Request-Id $req_id;
// proxy_set_header X-Request-Origin-Host $host;
// proxy_set_header X-Request-Origin-Path $request_uri;
// proxy_set_header X-Request-Origin-Method $request_method;
func UseRemoteAuth(c *gin.Context, user *SigninUser, authz string) bool {
	code := 0 // http code
	body := new([]byte)
	head := gout.H{
		"X-Request-Origin-Host":   c.Request.Host,
		"X-Request-Origin-Path":   c.Request.RequestURI,
		"X-Request-Origin-Method": c.Request.Method,
	}
	// 需要拷贝令牌： header['Authorization'] => api
	for k, v := range c.Request.Header {
		head[k] = strings.Join(v, "; ") // 拷贝请求头（拷贝过程中， 其中会带有cookie信息）
	}
	err := gout.GET(authz).
		SetHeader(head).
		BindHeader(user).
		BindBody(body).
		Code(&code).
		Do()
	if err != nil {
		// 远程认证访问发生异常
		FixResponse500Error2Logger(NewGinContext(c), err)
		return false
	} else if code >= 400 {
		// 直接返回上级服务的结果
		c.Data(code, res.ResponseTypeJSON, *body)
		c.Abort()
		return false
	}
	return true
}

//==================================================================================
//==================================================================================
//==================================================================================

// UserIdx ...
type UserIdx interface {
	GetAccountIdx() int64
	GetUserIdx() int64
}

var _ auth.UserInfo = &SigninUser{}
var _ UserIdx = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	TokenID   string `header:"X-Request-Z-Token-Kid"`
	TokenPID  string `header:"X-Request-Z-Token-Pid"`
	AccoIdx   int64  `header:"X-Request-Z-Account-Id"`
	UserIdx   int64  `header:"X-Request-Z-User-Id"`
	Account   string `header:"X-Request-Z-Account"`
	Account1  string `header:"X-Request-Z-Account1"`
	Account2  string `header:"X-Request-Z-Account2"`
	UserID    string `header:"X-Request-Z-User-Kid"`
	UserName  string `header:"X-Request-Z-User-Name"`
	UserRoles string `header:"X-Request-Z-User-Roles"`
	OrgCode   string `header:"X-Request-Z-Org-Code"`
	OrgAdmin  string `header:"X-Request-Z-Org-Admin"`
	OrgUsrID  string `header:"X-Request-Z-Org-Usrid"`
	Agent     string `header:"X-Request-Z-Agent"`
	Scope     string `header:"X-Request-Z-Scope"`
	Domain    string `header:"X-Request-Z-Domain"`
	Issuer    string `header:"X-Request-Z-Issuer"`
	Audience  string `header:"X-Request-Z-Audience"`
	ZgoXip    string `header:"X-Request-Z-Xip"`
}

// GetAccountIdx ...
func (u *SigninUser) GetAccountIdx() int64 {
	return u.AccoIdx
}

// GetUserIdx ...
func (u *SigninUser) GetUserIdx() int64 {
	return u.UserIdx
}

// GetTokenID xxx
func (u *SigninUser) GetTokenID() string {
	return u.TokenID
}

// GetTokenPID xxx
func (u *SigninUser) GetTokenPID() string {
	return u.TokenPID
}

// GetAccount xxx
func (u *SigninUser) GetAccount() string {
	return u.Account
}

// GetAccount1 xxx
func (u *SigninUser) GetAccount1() string {
	return u.Account1
}

// GetAccount2 xxx
func (u *SigninUser) GetAccount2() string {
	return u.Account2
}

// GetUserID xxx
func (u *SigninUser) GetUserID() string {
	return u.UserID
}

// GetUserName xxx
func (u *SigninUser) GetUserName() string {
	return u.UserName
}

// GetUserRoles xxx
func (u *SigninUser) GetUserRoles() []string {
	if u.UserRoles == "" {
		return nil
	}
	return strings.Split(u.UserRoles, ";")
}

// GetOrgCode xxx
func (u *SigninUser) GetOrgCode() string {
	return u.OrgCode
}

// GetOrgAdmin xxx
func (u *SigninUser) GetOrgAdmin() string {
	return u.OrgAdmin
}

// GetOrgUsrID xxx
func (u *SigninUser) GetOrgUsrID() string {
	return u.OrgUsrID
}

// GetAgent xxx
func (u *SigninUser) GetAgent() string {
	return u.Agent
}

// GetScope xxx
func (u *SigninUser) GetScope() string {
	return u.Scope
}

// GetDomain xxx
func (u *SigninUser) GetDomain() string {
	return u.Domain
}

// GetIssuer xxx
func (u *SigninUser) GetIssuer() string {
	return u.Issuer
}

// GetAudience xxx
func (u *SigninUser) GetAudience() string {
	return u.Audience
}

// GetUserSvcRoles xxx
func (u *SigninUser) GetUserSvcRoles(svc string) []string {
	roles := []string{}
	for _, role := range u.GetUserRoles() {
		if strings.HasPrefix(role, svc) {
			roles = append(roles, role)
		}
	}
	return roles
}

// SetUserRoles xxx
func (u *SigninUser) SetUserRoles(roles []string) {
	if roles == nil {
		u.UserRoles = ""
	} else {
		u.UserRoles = strings.Join(roles, ";")
	}
}
