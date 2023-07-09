package log

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
)

const (
	LOGGER_ERR_TYPE      = "err"
	LOGGER_ACCESS_TYPE   = "access"
	LOGGER_BUSINESS_TYPE = "business"
)

var (
	globalRootPath string // 日志统一存储根目录
	globalLevel    uint8  // 配置日志级别
	globalAppName  string // app name
)

var (
	ErrLog      LoggerContextIface
	AccessLog   LoggerContextIface
	BusinessLog LoggerContextIface
)

var levelMap = map[string]uint8{
	"debug": LEVEL_DEBUG,
	"info":  LEVEL_INFO,
	"warn":  LEVEL_WARN,
	"error": LEVEL_ERROR,
}

func InitGlobalConfig(rootPath string, level string, appName string) error {
	if rootPath == "" {
		return errors.New("初始化全局配置失败,参数rootPath不能为空")
	}
	if _, ok := levelMap[level]; !ok {
		return errors.New("初始化全局配置失败,参数level错误")
	}
	if appName == "" {
		return errors.New("初始化全局配置失败,参数appName不能为空")
	}
	globalRootPath = rootPath
	globalLevel = levelMap[level]
	globalAppName = appName
	return nil
}

func GetAccessLogger(loggerTag string) error {
	var err error
	AccessLog, err = getLogger(LOGGER_ACCESS_TYPE, loggerTag, false)
	return err
}
func GetBusinessLogger(loggerTag string) error {
	var err error
	BusinessLog, err = getLogger(LOGGER_BUSINESS_TYPE, loggerTag, false)
	return err
}

func GetErrLogger(loggerTag string) error {
	var err error
	ErrLog, err = getLogger(LOGGER_ERR_TYPE, loggerTag, true)
	return err
}

func contextCall(l *LoggerContext) LoggerContextIface {
	return l
}

func getLogger(loggerType string, loggerTag string, caller bool) (*LoggerContext, error) {
	log, err := CreateLogger(LoggerConf{
		LogFilePath: fmt.Sprintf("%s/%s/%s/log", globalRootPath, loggerType, globalAppName),
		Level:       globalLevel,
		Caller:      caller,
		ContextCall: contextCall,
	})
	if err != nil {
		return nil, err
	}
	logger := log.Tag(loggerTag).WithCommonFields(Fields{
		"application": globalAppName,
	})
	return logger.(*LoggerContext), nil
}

// 创建日志对象
func CreateLogger(logConf LoggerConf) (LoggerContextIface, error) {
	if logConf.LogFilePath == "" {
		return nil, errors.New("Log Path is empty")
	}
	globalLogger := newLoggerContext()
	globalLogger.loggerConf = logConf
	fileWriter, err := newFileWriter(logConf.LogFilePath)
	if err != nil {
		return nil, err
	}
	ioLogger := zerolog.New(fileWriter)
	ioLogger.Level(zerolog.Level(globalLogger.loggerConf.Level))

	globalLogger.logger = &ioLogger
	return globalLogger, nil
}
