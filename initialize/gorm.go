package initialize

import (
	"go.uber.org/zap"
	"gorm-demo/constants"
	_ "gorm-demo/core"
	"gorm-demo/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
)

var err error

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() {
	switch constants.GVA_CONFIG.System.DbType {
	case "mysql":
		GormMysql()
	case "postgresql":
		GormPostgreSql()
	//case "sqlite": // sqlite需要gcc支持 windows用户需要自行安装gcc 如需使用打开注释即可
	//	GormSqlite()
	case "sqlserver":
		GormSqlServer()
	default:
		GormMysql()
	}
}

// Gorm Mysql 初始化Mysql数据库
func GormMysql() {
	m := constants.GVA_CONFIG.Mysql
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
	if constants.GVA_DB, err = gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		constants.GVA_LOG.Error("MySQL start fail", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(constants.GVA_DB)
		sqlDB, _ := constants.GVA_DB.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	}
}

// GormPostgreSql 初始化PostgreSql数据库
func GormPostgreSql() {
	p := constants.GVA_CONFIG.Postgresql
	dsn := "host=" + p.Host + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
	postgresConfig := postgres.Config{
		DSN:                  dsn,                    // DSN data source name
		PreferSimpleProtocol: p.PreferSimpleProtocol, // 禁用隐式 prepared statement
	}
	gormConfig := logConfig(p.Logger)
	if constants.GVA_DB, err = gorm.Open(postgres.New(postgresConfig), gormConfig); err != nil {
		constants.GVA_LOG.Error("PostgreSql启动异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(constants.GVA_DB)
		sqlDB, _ := constants.GVA_DB.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	}
}

// GormSqlite 初始化Sqlite数据库 sqlite需要gcc支持 windows用户需要自行安装gcc 如需使用打开注释即可
//func GormSqlite() {
//	s := global.GVA_CONFIG.Sqlite
//	gormConfig := logConfig(s.Logger)
//	if global.GVA_DB, err = gorm.Open(sqlite.Open(s.Path), gormConfig); err != nil {
//		global.GVA_LOG.Error("Sqlite启动异常", zap.Any("err", err))
//		os.Exit(0)
//	} else {
//		sqlDB, _ := global.GVA_DB.DB()
//		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
//		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
//	}
//}

// GormSqlServer 初始化SqlServer数据库
func GormSqlServer() {
	ss := constants.GVA_CONFIG.Sqlserver
	dsn := "sqlserver://" + ss.Username + ":" + ss.Password + "@" + ss.Path + "?database=" + ss.Dbname
	if constants.GVA_DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{}); err != nil {
		constants.GVA_LOG.Error("SqlServer start fail", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(constants.GVA_DB)
		sqlDB, _ := constants.GVA_DB.DB()
		sqlDB.SetMaxIdleConns(ss.MaxIdleConns)
		sqlDB.SetMaxOpenConns(ss.MaxOpenConns)
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
		model.User{},
	)
	if err != nil {
		constants.GVA_LOG.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	constants.GVA_LOG.Info("register table success")
}
