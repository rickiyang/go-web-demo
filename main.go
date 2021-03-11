package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm-demo/dao"
	"gorm-demo/initialize"
	"gorm-demo/model"
)

func main() {
	initialize.Gorm()
	insertTest()
}

func insertTest() {
	user := model.User{
		Name:  "发森",
		Age:   12,
		Sex:   0,
		Phone: "13321234564",
	}
	//user1 := model.User{
	//	Name:  "凤岗",
	//	Age:   33,
	//	Sex:   0,
	//	Phone: "13214325689",
	//}
	//users :=  []model.User{user, user1}
	//
	//dao.BatchInsert(users)
	dao.InsertOne(user)

}
