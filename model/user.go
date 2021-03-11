package model

type User struct {
	Id    int64  `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Age   int32  `json:"age"`
	Sex   int8   `json:"sex"`
	Phone string `json:"phone"`
}
