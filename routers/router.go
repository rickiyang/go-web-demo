package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"gorm-demo/config"
	"gorm-demo/routers/api"
	"net/http"
)

var (
	log    = config.GVA_LOG
	origin = "www.baidu.com"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(testGlobalMiddleWare())
	r.Use(cors())

	// 给表单限制上传大小 (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	apiUser := r.Group("/api/user")
	apiUser.Use()
	{
		//根据uid查询用户信息
		apiUser.GET("/getUserByUid/", api.GetUserByUid)

		//根据uid查询用户信息 - 参数在请求路径中
		apiUser.GET("/getUserByUid/:uid", api.GetUserByUidUseRouteParam)

		// 新增用户，表单请求
		apiUser.POST("/addUser", api.AddUser)
		//新增用户， json 格式请求
		apiUser.POST("/addUserUseJson", api.AddUserUseJson)
	}
	//单文件上传
	r.POST("/upload/singleFile", api.PostSingleFile)
	//多文件上传
	r.POST("upload/multiFile", api.PostMultiFiles)

	return r
}

// 全局中间件示例
func testGlobalMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("MiddleWare: 中间件开始执行")

		// 在gin.Context中设置一个值 演示中间件的能力
		traceId, _ := uuid.GenerateUUID()
		c.Set("trace_id", traceId)

		//todo 这里你可以执行你想做的任何事情

		// 执行完这里的逻辑之后别忘了 调用 Next 函数将请求交给下个 handler 处理
		c.Next()

		status := c.Writer.Status()
		log.Info("MiddleWare: 中间件执行结束, status: ", zap.Any("status", status))
	}
}

//添加跨域支持
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if origin != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
