package routers

import (
	"github.com/gin-gonic/gin"
	"gorm-demo/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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
