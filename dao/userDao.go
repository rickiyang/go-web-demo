package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm-demo/initialize"
	"gorm-demo/model"
	"strconv"
)

//多行数据查询操作
func FindUserByIds(uids []int64) map[int64]model.User {
	result := make(map[int64]model.User)
	if len(uids) <= 0 {
		return result
	}
	var searchParam string
	for idx, uid := range uids {
		if idx == 0 {
			searchParam = strconv.FormatInt(uid, 10)
		} else {
			searchParam = fmt.Sprintf("%s,%v", searchParam, uid)
		}
	}
	sql := "select * from user where id in (" + searchParam + ")"
	rows, e := initialize.SqlDB.Query(sql)
	if e == nil {
		errors.New("query incur error")
	}
	for rows.Next() {
		var user model.User
		e := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Phone)
		if e == nil {
			bytes, _ := json.Marshal(user)
			fmt.Println(string(bytes))
			result[user.Id] = user
		}
	}
	defer rows.Close()
	return result
}

//单行查询操作
func FindUserById(uid int64) model.User {
	var user model.User
	searchSql := fmt.Sprintf("%s%v", "select * from user where id=", uid)
	initialize.SqlDB.QueryRow(searchSql).Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Phone)
	return user
}

//预处理使用
func PrepareUse(uid int64) model.User {
	var user model.User

	searchSql := fmt.Sprintf("%s%v", "select * from user where id=", uid)
	searchSql = "select * from user where id = ?"

	stmt, e := initialize.SqlDB.Prepare(searchSql)
	if e != nil {
		fmt.Println(e)
		return model.User{}
	}
	defer stmt.Close()
	rows, e := stmt.Query(uid)
	if e != nil {
		fmt.Println(e)
	}
	rows.Scan(&user.Id, &user.Name, &user.Age, &user.Sex, &user.Phone)
	return user
}

//插入新数据
func Insert(user model.User) (err error) {
	sqlStr := "insert into user(name,age,sex,phone) values(?,?,?,?)"
	result, err := initialize.SqlDB.Exec(sqlStr, user.Name, user.Age, user.Sex, user.Phone)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()

	fmt.Println("id is ", id)
	return nil
}

//更新数据
func UpdateById(user model.User) (err error) {
	sqlStr := "update user set name=?,age=?,sex=?,phone=? where id=?"
	result, err := initialize.SqlDB.Exec(sqlStr, user.Name, user.Age, user.Sex, user.Phone, user.Id)
	if err != nil {
		return err
	}
	affectRow, _ := result.RowsAffected()
	fmt.Println("affectRow is ", affectRow)
	return nil
}

//删除数据
func DeleteById(id int64) (err error) {
	sqlStr := "delete from user where id=?"
	result, err := initialize.SqlDB.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	affectRow, _ := result.RowsAffected()
	fmt.Println("affectRow is ", affectRow)
	return nil
}

func TestTx() {

	//开启事务
	tx, err := initialize.SqlDB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	//准备sql语句
	stmt, err := tx.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		fmt.Println("Prepare fail")
		return
	}
	//设置参数以及执行sql语句
	res, err := stmt.Exec(2)
	fmt.Println(res)
	if err != nil {
		fmt.Println("DELETE Exec fail")
		return
	}

	panic(errors.New(" special err in the cur throw").Error())

	sqlStr := "insert into user(name,age,sex,phone) values(?,?,?,?)"
	result, err := tx.Exec(sqlStr, "方德峰", 23, 1, "13321123543")
	if err != nil {
		fmt.Println("insert Exec fail")
		return
	}
	affectRow, _ := result.RowsAffected()
	fmt.Println("affectRow is ", affectRow)
	//提交事务
	tx.Commit()

}
