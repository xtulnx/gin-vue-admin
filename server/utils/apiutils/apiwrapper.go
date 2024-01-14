package apiutils

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"reflect"
)

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
	ginType = reflect.TypeOf((*gin.Context)(nil)).Elem()
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

// WrapperApi0 包装器, 用于将 gin 的请求参数绑定到结构体上，
//
// 特殊参数: context.Context, *gin.Context
//
// 不列出需要解析的参数的话，则只有第一个非 context 的参数用于接收请求参数
//
// jason.liao 2023.12.12
func WrapperApi0(callbackFunc any, reqSample ...any) gin.HandlerFunc {
	fnType := reflect.TypeOf(callbackFunc)
	if fnType.Kind() != reflect.Func {
		panic("callbackFunc must be func")
	}
	fnVal := reflect.ValueOf(callbackFunc)
	reqTypes := make(map[reflect.Type]reflect.Value)
	for _, v := range reqSample {
		rt1, rv1 := reflect.TypeOf(v), reflect.ValueOf(v)
		if rt1.Kind() == reflect.Ptr {
			rt1 = rt1.Elem()
			rv1 = rv1.Elem()
		}
		reqTypes[rt1] = rv1
	}
	skipReqType := len(reqTypes) > 0
	return func(c *gin.Context) {
		var args = make([]reflect.Value, fnType.NumIn())
		for i := 0; i < fnType.NumIn(); i++ {
			var rt1 = fnType.In(i)
			if rt1 == ctxType {
				args[i] = reflect.ValueOf(c.Request.Context())
				continue
			}
			if rt1 == ginType {
				args[i] = reflect.ValueOf(c)
				continue
			}
			var rt2 = rt1
			if rt1.Kind() == reflect.Ptr {
				rt2 = rt1.Elem()
				args[i] = reflect.New(rt1.Elem())
			} else {
				args[i] = reflect.New(rt1).Elem()
			}
			if rt2.Kind() != reflect.Struct {
				continue
			}
			if skipReqType {
				//if _, ok := reqTypes[rt2]; !ok {
				//	continue
				//}
			}
			if err := GinMustBind(c, args[i].Interface()); err != nil {
				response.Err(err, c)
				return
			}
		}
		var resp = fnVal.Call(args)
		var data interface{}
		var err error
		if len(resp) > 0 {
			for _, v := range resp {
				if v.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					if !v.IsNil() {
						err = v.Interface().(error)
					}
				} else {
					data = v.Interface()
				}
			}
		}
		response.ResultErr(data, err, c)
	}
}
