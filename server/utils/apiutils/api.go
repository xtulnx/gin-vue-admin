package apiutils

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/apimodel"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

var ginBinding = map[string]binding.Binding{}

// RegBinding 注册请求与 model 的解析器
func RegBinding(mime string, b binding.Binding) {
	ginBinding[mime] = b
}

var bindingHandler []BindingHandler

type BindingHandler func(c *gin.Context, req any) error

func RegBindingHandler(h BindingHandler) {
	bindingHandler = append(bindingHandler, h)
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

type tWithNow interface {
	SetNow(now time.Time)
}
type tWithIP interface {
	SetIP(string)
}
type tWithCtx interface {
	SetCtx(ctx context.Context)
	SetCtxValue(k, v interface{})
}
type tWithAuthUser interface {
	SetAuthUser(info apimodel.IUserInfoProvider)
}
type tWithGinContext interface {
	SetGinContext(c *gin.Context)
}
type tWithBinder interface {
	BindGinContext(c *gin.Context) error
}

// GinMustBind 示例：解析请求参数
func GinMustBind(c *gin.Context, obj interface{}) error {
	reqMethod, reqContentType := c.Request.Method, c.ContentType()
	var b binding.Binding = nil
	if ginBinding != nil {
		if reqMethod == http.MethodGet {
			b, _ = ginBinding[binding.MIMEPOSTForm]
		} else {
			b, _ = ginBinding[reqContentType]
		}
	}
	if b == nil {
		b = binding.Default(reqMethod, reqContentType)
	}
	err := c.MustBindWith(obj, b)
	return GinBindAfter(c, obj, err)
}

// GinBindAfter 参数事后处理
func GinBindAfter(c *gin.Context, obj interface{}, err error) error {
	if err != nil {
		return err
	}
	if m, ok := obj.(tWithNow); ok {
		m.SetNow(time.Now())
	}
	if m, ok := obj.(tWithGinContext); ok {
		m.SetGinContext(c)
	}
	if m, ok := obj.(tWithCtx); ok {
		m.SetCtx(c.Request.Context())
	}
	if m, ok := obj.(tWithIP); ok {
		m.SetIP(c.ClientIP())
	}
	if m, ok := obj.(tWithAuthUser); ok {
		m.SetAuthUser(utils.GetUserInfo(c))
	}
	if m, ok := obj.(tWithBinder); ok {
		err = m.BindGinContext(c)
		if err != nil {
			return err
		}
	}
	if bindingHandler != nil {
		for _, h := range bindingHandler {
			err = h(c, obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type API struct {
}

func (API) GinBindAfter(c *gin.Context, obj interface{}, err error) error {
	return GinBindAfter(c, obj, err)
}

func (API) GinMustBind(c *gin.Context, obj interface{}) error {
	return GinMustBind(c, obj)
}
