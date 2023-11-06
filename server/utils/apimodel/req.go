package apimodel

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// 封装一些常用请求参数作为上下文
//
// jason.liao 2023.08.01

// ReqWithGin 保留 gin 的环境
type ReqWithGin struct {
	c *gin.Context
}

func (R *ReqWithGin) SetGinContext(c *gin.Context) {
	R.c = c
}

func (R *ReqWithGin) GetGinContext() *gin.Context {
	return R.c
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithNow 统一业务时间点
type ReqWithNow struct {
	now time.Time
}

func (R *ReqWithNow) SetNow(now time.Time) {
	R.now = now
}
func (R *ReqWithNow) GetNow() time.Time {
	if R.now.IsZero() {
		R.now = time.Now()
	}
	return R.now
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithIP 客户端 ip
type ReqWithIP struct {
	ip string
}

func (R *ReqWithIP) SetIP(ip string) {
	R.ip = ip
}
func (R *ReqWithIP) GetIP() string {
	return R.ip
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithCtx 上下文环境
type ReqWithCtx struct {
	ctx context.Context
}

func (R *ReqWithCtx) SetCtx(ctx context.Context) {
	R.ctx = ctx
}
func (R *ReqWithCtx) GetCtx() context.Context {
	return R.ctx
}
func (R *ReqWithCtx) SetCtxValue(k, v interface{}) {
	R.ctx = context.WithValue(R.ctx, k, v)
}
func (R *ReqWithCtx) GetCtxValue(k interface{}) interface{} {
	return R.ctx.Value(k)
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

type IUserInfoProvider interface {
	GetUserID() uint
	GetUserAuthorityId() uint
	GetUserName() string
	GetNickName() string
}

// ReqWithAuthUser 用户信息
type ReqWithAuthUser struct {
	user IUserInfoProvider
}

func (R *ReqWithAuthUser) GetUserID() uint {
	if R.user == nil {
		return 0
	}
	return R.user.GetUserID()
}

func (R *ReqWithAuthUser) GetUserAuthorityId() uint {
	if R.user == nil {
		return 0
	}
	return R.user.GetUserAuthorityId()
}

func (R *ReqWithAuthUser) GetUserName() string {
	if R.user == nil {
		return ""
	}
	return R.user.GetUserName()
}

func (R *ReqWithAuthUser) GetNickName() string {
	if R.user == nil {
		return ""
	}
	return R.user.GetNickName()
}

func (R *ReqWithAuthUser) SetAuthUser(info IUserInfoProvider) {
	R.user = info
}

func (R *ReqWithAuthUser) GetAuthUser() IUserInfoProvider {
	return R.user
}
