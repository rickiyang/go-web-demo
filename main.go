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
