package first

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm-demo/config"
	"gorm-demo/utils"
	"log"
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
			log.Printf("use local config, path=%v", utils.ConfigFile)
		} else {
			configCon = configEnv
			log.Printf("use remote config, path=%v", configCon)
		}
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
		if err := v.Unmarshal(&config.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&config.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	config.GVA_VP = v
}
