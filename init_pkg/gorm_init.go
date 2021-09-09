package init_pkg

import (
	"go.uber.org/zap"
	"gorm-demo/config"
	_ "gorm-demo/init_pkg/first"
	"gorm-demo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() {
	GormMysql()
}

var err error

// Gorm Mysql 初始化Mysql数据库
func GormMysql() {
	m := config.GVA_CONFIG.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	gormConfig := logConfig(m.LogMode)
	if config.GVA_DB, err = gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		config.GVA_LOG.Error("MySQL start fail", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(config.GVA_DB)
		sqlDB, _ := config.GVA_DB.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	}
}

// logConfig 根据配置决定是否开启日志
func logConfig(mod bool) (c *gorm.Config) {
	if mod {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
			//SkipDefaultTransaction: true,                         //禁用事务操作
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
				//NameReplacer:  strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
			},
		}

	} else {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
			},
		}
	}
	return
}

// GormDBTables 注册数据库表专用
func GormDBTables(db *gorm.DB) {
	err := db.AutoMigrate(
		models.User{},
	)
	if err != nil {
		config.GVA_LOG.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	config.GVA_LOG.Info("register table success")
}
