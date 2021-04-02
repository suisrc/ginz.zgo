module github.com/suisrc/ginz.zgo

go 1.16

replace (
	github.com/suisrc/auth.zgo v0.0.0 => ../auth
	github.com/suisrc/buntdb.zgo v0.0.0 => ../buntdb
	github.com/suisrc/config.zgo v0.0.0 => ../config
	github.com/suisrc/crypto.zgo v0.0.0 => ../crypto
	github.com/suisrc/logger.zgo v0.0.0 => ../logger
	github.com/suisrc/res.zgo v0.0.0 => ../res
	github.com/suisrc/utils.zgo v0.0.0 => ../utils
)

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/google/uuid v1.2.0
	github.com/suisrc/auth.zgo v0.0.0
	github.com/suisrc/config.zgo v0.0.0
	github.com/suisrc/gin-i18n v0.1.1
	github.com/suisrc/logger.zgo v0.0.0
	github.com/suisrc/res.zgo v0.0.0
	github.com/suisrc/utils.zgo v0.0.0
	github.com/urfave/cli/v2 v2.3.0
)
