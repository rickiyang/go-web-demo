package initialize

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func init() {
	SqlDB, _ = sql.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/test_db")
	//设置数据库最大连接数
	SqlDB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	SqlDB.SetMaxIdleConns(10)
	//验证连接
	if err := SqlDB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success")

}
