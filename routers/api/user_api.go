package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm-demo/config"
	"gorm-demo/models"
	"gorm-demo/pkg"
	"gorm-demo/service"
	"net/http"
	"strconv"
)

var (
	log = config.GVA_LOG
)

//get请求，param 参数
func GetUserByUid(c *gin.Context) {
	newGin := pkg.Gin{C: c}
	uid, err := strconv.ParseInt(newGin.C.Query("uid"), 10, 64)
	if err != nil {
		log.Error("uid parse err", zap.Any("err", err))
		newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
		return
	}
	user := service.QueryUserByUid(uid)
	newGin.Response(http.StatusOK, pkg.SUCCESS, user)
}

//get请求 获取请求路径中的 参数
func GetUserByUidUseRouteParam(c *gin.Context) {
	newGin := pkg.Gin{C: c}
	uid, err := strconv.ParseInt(newGin.C.Param("uid"), 10, 64)
	if err != nil {
		log.Error("uid parse err", zap.Any("err", err))
		newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
		return
	}
	user := service.QueryUserByUid(uid)
	newGin.Response(http.StatusOK, pkg.SUCCESS, user)
}

// post 请求， 普通 form 表单获取参数
func AddUser(c *gin.Context) {
	newGin := pkg.Gin{C: c}

	name := newGin.C.PostForm("name")
	age, _ := strconv.ParseInt(newGin.C.PostForm("age"), 10, 32)
	sex, _ := strconv.ParseInt(newGin.C.PostForm("sex"), 10, 8)
	phone := newGin.C.PostForm("phone")

	user := models.User{0, name, int32(age), int8(sex), phone}
	err := service.AddNewUser(user)
	if err != nil {
		log.Error("AddUser err", zap.Any("err", err))
		newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
		return
	}
	newGin.Response(http.StatusOK, pkg.SUCCESS, user)
}

// post 请求， json 格式参数
func AddUserUseJson(c *gin.Context) {
	newGin := pkg.Gin{C: c}

	//第一种方式是榜单一个结构体
	var user models.User
	newGin.C.BindJSON(&user)

	//第二种方式可以绑定一个 map, 使用之前需要将第一次的注释掉，参数只能读取一次
	var user1 map[string]interface{}
	newGin.C.BindJSON(&user1)
	bytes, _ := json.Marshal(user1)
	json.Unmarshal(bytes, &user)

	err := service.AddNewUser(user)
	if err != nil {
		log.Error("AddUserUseJson err", zap.Any("err", err))
		newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
		return
	}
	newGin.Response(http.StatusOK, pkg.SUCCESS, user)
}

// post 请求， 单文件上传
// 注意文件路径这里要是具体的路径 + 文件名
func PostSingleFile(c *gin.Context) {
	newGin := pkg.Gin{C: c}
	file, _ := c.FormFile("file")
	// 上传文件到指定的路径
	filePath := "/Users/xiaoming/"
	err := c.SaveUploadedFile(file, filePath+file.Filename)
	if err != nil {
		log.Error("PostSingleFile err", zap.Any("err", err))
		newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
		return
	}
	newGin.Response(http.StatusOK, pkg.SUCCESS, "success")
}

// post 请求， 多文件上传
// 注意文件路径这里要是具体的路径 + 文件名
func PostMultiFiles(c *gin.Context) {
	newGin := pkg.Gin{C: c}
	// 上传文件到指定的路径
	filePath := "/Users/xiaoming/"
	form, _ := c.MultipartForm()
	files := form.File["upload"]
	for _, file := range files {
		//上传文件到指定的路径
		err := c.SaveUploadedFile(file, filePath+file.Filename)
		if err != nil {
			log.Error("PostMultiFiles err", zap.Any("err", err))
			newGin.Response(http.StatusBadRequest, pkg.INVALID_PARAMS, nil)
			return
		}
	}
	newGin.Response(http.StatusOK, pkg.SUCCESS, "success")
}
