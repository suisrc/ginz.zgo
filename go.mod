module github.com/suisrc/ginz.zgo

go 1.16

//replace (
//	github.com/suisrc/auth.zgo => ../auth
//	github.com/suisrc/buntdb.zgo => ../buntdb
//	github.com/suisrc/config.zgo => ../config
//	github.com/suisrc/crypto.zgo => ../crypto
//	github.com/suisrc/logger.zgo => ../logger
//	github.com/suisrc/res.zgo => ../res
//	github.com/suisrc/utils.zgo => ../utils
//)

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/google/uuid v1.2.0
	github.com/guonaihong/gout v0.1.8
	github.com/suisrc/auth.zgo v0.0.0-20210408060712-08bb878db327
	github.com/suisrc/config.zgo v0.0.0-20210407020836-a5a7e3c8595d
	github.com/suisrc/gin-i18n v0.1.1
	github.com/suisrc/logger.zgo v0.0.0-20210408054212-b4e804e2dc15
	github.com/suisrc/res.zgo v0.0.0-20210408020700-20221959252e
	github.com/suisrc/utils.zgo v0.0.0-20210402013101-2fffea7939ab
	github.com/urfave/cli/v2 v2.3.0
)
