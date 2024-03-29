package first

import (
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	globalConfig "gorm-demo/config"
	"gorm-demo/utils"
	"log"
	"os"
	"time"
)

var (
	err    error
	level  zapcore.Level
	writer zapcore.WriteSyncer
)

func init() {
	// 判断是否有Director文件夹
	if ok, _ := utils.PathExists(globalConfig.GVA_CONFIG.Zap.Director); !ok {
		log.Printf("create %v directory\n", globalConfig.GVA_CONFIG.Zap.Director)
		_ = os.Mkdir(globalConfig.GVA_CONFIG.Zap.Director, os.ModePerm)
	}

	// 初始化配置文件的Level
	switch globalConfig.GVA_CONFIG.Zap.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	writer, err = getWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		globalConfig.GVA_LOG = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		globalConfig.GVA_LOG = zap.New(getEncoderCore())
	}
	if globalConfig.GVA_CONFIG.Zap.ShowLine {
		globalConfig.GVA_LOG.WithOptions(zap.AddCaller())
	}
}

// getWriteSyncer zap logger中加入file-rotatelogs
func getWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		globalConfig.GVA_CONFIG.Zap.Director+string(os.PathSeparator)+"%Y-%m-%d.log",
		zaprotatelogs.WithLinkName(globalConfig.GVA_CONFIG.Zap.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if globalConfig.GVA_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  globalConfig.GVA_CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case globalConfig.GVA_CONFIG.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case globalConfig.GVA_CONFIG.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case globalConfig.GVA_CONFIG.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case globalConfig.GVA_CONFIG.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if globalConfig.GVA_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	return zapcore.NewCore(getEncoder(), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(globalConfig.GVA_CONFIG.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}
