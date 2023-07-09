package vars

import (
	"errors"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/util/gconv"
)

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize     int
}

/**
 * apollo配置对应的结构体
 */

type EnvSettingS struct {
	TenantNamespace     string
	ApolloMetaServerUrl string
	ElasticApmServerUrl string
	EtcdV3ServerUrls    string
	IsDevelop           bool
	HostName            string
}

type MysqlSettingS struct {
	Host             string
	Port             string
	UserName         string
	Password         string
	DBName           string
	MoonDB           string
	Charset          string
	MaxConnLifetime  int
	MaxIdleConnCount int
	MaxOpenConnCount int
}

type RedisSettingS struct {
	RedisConnLink string
	Host          string
	Port          string
	Password      string
	Database      int
	MaxIdle       int
	MaxActive     int
	IdleTimeout   int
}

type RabbitMQSettingS struct {
	Host     string
	Port     string
	User     string
	Password string
	Vhost    string
	BackOff  string
	ApiPort  string
}

//OssSettings 对象存储配置信息
type OssSettings struct {
	AccessKeyId     string
	AccessKeySecret string
	RoleArn         string
	Bucket          string
	Endpoint        string
	SchemaForUI     string
}

var (
	AppSetting      = &AppSettingS{}
	EnvSetting      = &EnvSettingS{}
	MysqlSetting    = &MysqlSettingS{}
	RedisSetting    = &RedisSettingS{}
	RabbitMQSetting = &RabbitMQSettingS{}
)

// LoadApplicationConfig 加载配置对象映射
func LoadApplicationConfig() error {
	var config *gcfg.Config
	// if EnvSetting.IsDevelop { //判断是否开启了本地配置
	// 	config = g.Config()
	// }
	// 直接读取本地配置
	config = g.Config()

	err := loadConfig(config, "App", AppSetting)
	if err != nil {
		return err
	}

	err = loadConfig(config, "Mysql", MysqlSetting)
	if err != nil {
		return err
	}

	// err = loadConfig(config, "Redis", RedisSetting)
	// if err != nil {
	// 	return err
	// }

	// err = loadConfig(config, "RabbitMQ", RabbitMQSetting)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// 先加载apollo配置, 如果开启了本地配置(配置export IS_DEVELOP=1), apollo配置会被本地配置覆盖
func loadConfig(config *gcfg.Config, section string, structure interface{}) error {
	var apolloConfigErr error
	var fileConfigErr = errors.New(section + "初始化")
	// apolloConfigErr = apollo_sdk.MapApolloConfig(section, structure)
	if config != nil {
		fileConfigErr = gconv.Struct(config.Get(section), structure)
	}

	if apolloConfigErr != nil && fileConfigErr != nil {
		return errors.New(section + "加载配置失败")
	}

	return nil
}

// ReplaceSensitiveWords 替换敏感字符
func ReplaceSensitiveWords(raw string) string {
	sensitiveWords := []string{MysqlSetting.Password,
		RabbitMQSetting.Password, RedisSetting.Password, ApiAuthKey}
	for _, word := range sensitiveWords {
		if word == "" {
			continue
		}
		raw = strings.ReplaceAll(raw, word, "******")
	}
	return raw
}
