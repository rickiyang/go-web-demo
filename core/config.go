package core

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm-demo/constants"
	"gorm-demo/utils"
	"os"
)

var configCon string

func init() {
	flag.StringVar(&configCon, "c", "", "choose configCon file.")
	flag.Parse()
	// 优先级: 命令行 > 环境变量 > 默认值
	if configCon == "" {
		if configEnv := os.Getenv(utils.ConfigEnv); configEnv == "" {
			configCon = utils.ConfigFile
			fmt.Printf("您正在使用config的默认值,config的路径为%v\n", utils.ConfigFile)
		} else {
			configCon = configEnv
			fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", configCon)
		}
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", configCon)
	}
	v := viper.New()
	v.SetConfigFile(configCon)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configCon file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("configCon file changed:", e.Name)
		if err := v.Unmarshal(&constants.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&constants.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	constants.GVA_VP = v
}
