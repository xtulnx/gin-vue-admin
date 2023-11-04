package plugin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sort"
	"sync"
)

// 扩展插件注册模式
//
// jason.liao 2023.11.04

var (
	// ErrExist 已经存在
	ErrExist = errors.New("plugin exist")
	// ErrSame 有同名插件
	ErrSame = errors.New("plugin same name")
)

// Manger 插件管理器
type Manger struct {
	mc             sync.RWMutex
	setPlugin      map[Base]int
	mapName2Plugin map[string]Base
	lstPlugin      []Base
}

func (P *Manger) checkInitInner() {
	if P.setPlugin == nil {
		P.setPlugin = make(map[Base]int)
	}
	if P.mapName2Plugin == nil {
		P.mapName2Plugin = make(map[string]Base)
	}
}

// Register 注册插件
func (P *Manger) Register(p Base) error {
	P.mc.Lock()
	defer P.mc.Unlock()
	P.checkInitInner()
	code, priority := p.PluginCode(), 0
	if _, ok := P.setPlugin[p]; ok {
		return ErrExist
	}
	if _, ok := P.mapName2Plugin[code]; ok {
		return ErrSame
	}
	if p1, ok := p.(Priority); ok {
		priority = p1.PluginPriority()
	}
	P.setPlugin[p] = priority
	P.mapName2Plugin[code] = p
	idx := sort.Search(len(P.lstPlugin), func(i int) bool {
		return P.setPlugin[P.lstPlugin[i]] < priority
	})
	P.lstPlugin = append(P.lstPlugin, nil)
	copy(P.lstPlugin[idx+1:], P.lstPlugin[idx:])
	P.lstPlugin[idx] = p
	return nil
}

// Foreach 遍历所有插件
func (P *Manger) Foreach(f func(p Base) error) error {
	P.mc.RLock()
	defer P.mc.RUnlock()
	for _, p := range P.lstPlugin {
		if err := f(p); err != nil {
			return err
		}
	}
	return nil
}

// ForeachUnsafe 遍历所有插件，不加锁
func (P *Manger) ForeachUnsafe(f func(p Base) error) error {
	for _, p := range P.lstPlugin {
		if err := f(p); err != nil {
			return err
		}
	}
	return nil
}

// Base 阶段0 基础
type Base interface {
	PluginName() string // 插件名称
	PluginCode() string // 插件编码，必须唯一
}

type Priority interface {
	PluginPriority() int // 插件优先级，默认是0，越大越优先
}

// WithConfig 阶段1 配置
type WithConfig interface {
	PluginConfig() interface{} // 插件的配置项，返回对象指针，复用
}

type WithConfigReset interface {
	PluginConfigReset() // 插件的配置项重置
}

// WithDbInit 阶段2 数据库
type WithDbInit interface {
	// PluginInitDb 插件数据库初始化，可选
	//
	//   如果插件不是必须的，不要返回错误，否则会中止服务启动
	PluginInitDb(db *gorm.DB) error
}

// WithInit 阶段3 服务初始化
type WithInit interface {
	PluginInit() error // 插件初始化，如果插件不是必须的，不要返回错误，否则会中止服务启动
}

// WithRouter 阶段4 路由
type WithRouter interface {
	// PluginRouter 注册路由。如果插件不是必须的，不要返回错误，否则会中止服务启动
	//  pri 私有路由组
	//  pub 公共路由组
	PluginRouter(pri *gin.RouterGroup, pub *gin.RouterGroup) error
}

// WithRelease 阶段5 资源释放
type WithRelease interface {
	PluginRelease()
}
