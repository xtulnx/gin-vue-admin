package global

// 所有的事件

// 通用事件
const (
	EVENTROUTERRegisterInit    = "router.register.init"    // 路由注册开始, 参数: (app*gin.Engine)
	EVENTROUTERRegisterBegin   = "router.register.begin"   // 路由注册开始, 参数: (app *gin.Engine,root *gin.RouterGroup)
	EVENTROUTERRegisterSuccess = "router.register.success" // 路由注册成功, 参数: (app *gin.Engine,pri *gin.RouterGroup,pub *gin.RouterGroup)
)
