//go:build debug
// +build debug

package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/ebus"
	"github.com/gin-gonic/gin"
)

func init() {
	ebus.Subscribe(global.EVENTROUTERRegisterInit, func(r *gin.Engine) {
		r.Use(middleware.Cors())
	})
}
