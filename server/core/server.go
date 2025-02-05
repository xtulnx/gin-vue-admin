package core

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	if global.GVA_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	// 初始化插件
	if err := global.GVA_PLUGIN.ForeachUnsafe(func(p plugin.Base) error {
		if m1, ok := p.(plugin.WithInit); ok {
			if e1 := m1.PluginInit(); e1 != nil {
				global.GVA_LOG.Error("插件「"+p.PluginName()+"」初始化失败!", zap.Any("err", e1))
				return e1
			}
		}
		return nil
	}); err != nil {
		global.GVA_LOG.Error("插件初始化失败，停止启动！")
		return
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	defer func() {
		_ = global.GVA_PLUGIN.ForeachUnsafe(func(p plugin.Base) error {
			if m, ok := p.(plugin.WithRelease); ok {
				m.PluginRelease()
			}
			return nil
		})
	}()

	fmt.Printf(`
	欢迎使用 gin-vue-admin
	当前版本:v2.5.7
    加群方式:微信号：shouzi_1994 QQ群：622360840
	插件市场:https://plugin.gin-vue-admin.com
	GVA讨论社区:https://support.qq.com/products/371961
	默认自动化文档地址:http://127.0.0.1%s%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
	如果项目让您获得了收益，希望您能请团队喝杯可乐:https://www.gin-vue-admin.com/coffee/index.html
`, address, global.GVA_CONFIG.System.RouterPrefix)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
