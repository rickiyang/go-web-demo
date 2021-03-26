package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper
	GVA_CONFIG Server
)
