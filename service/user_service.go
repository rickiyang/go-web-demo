package service

import (
	"go.uber.org/zap"
	"gorm-demo/config"
	"gorm-demo/models"
)

var log = config.GVA_LOG

//新增用户
func AddNewUser(user models.User) (err error) {
	id, err := models.InsertOneUser(user)
	if err != nil {
		return err
	}
	log.Info("add new user ,return id={}", zap.Any("id", id))
	return nil
}

//根据uid 查询用户
func QueryUserByUid(uid int64) models.User {
	return models.QueryUserByUid(uid)
}
