package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/spf13/viper/remote"
	"gorm-demo/init_pkg"
	"gorm-demo/routers"
	"net/http"
	"time"
)

func main() {
	init_pkg.Gorm()
	gin.SetMode(gin.DebugMode)

	routersInit := routers.InitRouter()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        routersInit,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

//
//func main() {
//	router := gin.Default()
//
//	​
//	//  group: v1
//	v1 := router.Group("/v1")
//	{
//		v1.POST("/login", v1Login)
//		v1.POST("/submit", v1Submit)
//		v1.POST("/add", v1Add)
//	}
//	​
//	// group: v2
//	v2 := router.Group("/v2")
//	{
//		v2.POST("/login", v2Login)
//		v2.POST("/submit", v2Submit)
//		v2.POST("/add", v2Add)
//	}
//	router.Run(":8080")
//}
