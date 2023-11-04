package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
// Author SliverHorn
func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(
		// 系统模块表
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},
		system.SysChatGptOption{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")

	// 插件的数据库表
	if err = global.GVA_PLUGIN.ForeachUnsafe(func(p plugin.Base) error {
		if m, ok := p.(plugin.WithDbInit); ok {
			if e1 := m.PluginInitDb(db); e1 != nil {
				global.GVA_LOG.Error("插件「"+p.PluginName()+"」初始化数据库失败！", zap.Error(e1))
				return e1
			}
		}
		if m, ok := p.(WithTableInit); ok {
			if tables := m.PluginInitTables(); len(tables) > 0 {
				if e1 := db.AutoMigrate(tables...); e1 != nil {
					global.GVA_LOG.Error("插件「"+p.PluginName()+"」注册数据库表失败！", zap.Error(e1))
					return e1
				}
			}
		}
		if m, ok := p.(WithMenuInit); ok {
			if menus := m.PluginInitMenus(); len(menus) > 0 {
				utils.RegisterMenus(menus...)
			}
		}
		if m, ok := p.(WithApiInit); ok {
			if apis := m.PluginInitApis(); len(apis) > 0 {
				utils.RegisterApis(apis...)
			}
		}
		return nil
	}); err != nil {
		global.GVA_LOG.Error("插件注册数据库表失败，停止启动！")
		os.Exit(1)
	}
}
