package config

import (
	"fmt"
	"os"
	"time"

	"arclinks-go/app/common/util"

	"github.com/yunkeCN/ali-logger-golang/logger"

	"github.com/gogf/gf/database/gredis"
	"github.com/pkg/errors"

	"arclinks-go/app/vars"
	"arclinks-go/library/log"
	"arclinks-go/library/rabbitmq"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.elastic.co/apm/module/apmgorm"
)

var appName = "arclinks-go"

func InitConfig() error {
	// 统一读取环境变量
	loadEnvConfig()

	// 1.初始化log配置
	err := initLog()
	if err != nil {
		log.ErrLog.Errorf("init log err: %v", err)
		return err
	}

	// 2.加载配置, 配置读取顺序: apollo配置>本地配置(配置export IS_DEVELOP=1, 本地配置会覆盖Apollo配置)
	err = loadApplicationConfig()
	if err != nil {
		log.ErrLog.Errorf("InitConfig loadApplicationConfig err: %v", err)
		return err
	}

	// 3.apollo热更新
	// apollo_sdk.TriggerApolloHotUpdateListen(vars.EnvSetting.ApolloMetaServerUrl, appName, "application")

	// 4.初始化db配置
	err = initDbConfig()
	if err != nil {
		log.ErrLog.Errorf("InitConfig initDbConfig err: %v", err)
		return err
	}

	// // 5.初始化redis配置
	// err = initRedis()
	// if err != nil {
	// 	log.ErrLog.Errorf("InitConfig initRedis err: %v", err)
	// 	return err
	// }

	// // 6.初始化RabbitMQ配置
	// err = initRabbitMQ()
	// if err != nil {
	// 	log.ErrLog.Errorf("InitConfig initRabbitMQ err: %v", err)
	// 	return err
	// }

	// 7.初始化mysql连接池 gorm
	err = initMysql()
	if err != nil {
		log.ErrLog.Errorf("InitConfig initMysql err: %v", err)
		return err
	}

	return nil
}

func loadEnvConfig() {
	// vars.EnvSetting.TenantNamespace = os.Getenv("TENANT_NAMESPACE")
	// vars.EnvSetting.ApolloMetaServerUrl = os.Getenv("APOLLO_META_SERVER_URL")
	// vars.EnvSetting.ElasticApmServerUrl = os.Getenv("ELASTIC_APM_SERVER_URL")
	// vars.EnvSetting.EtcdV3ServerUrls = os.Getenv("ETCDV3_SERVER_URLS")
	// vars.EnvSetting.HostName = os.Getenv("HOSTNAME")

	_, ok := os.LookupEnv("IS_DEVELOP")
	vars.EnvSetting.IsDevelop = ok
}

// 初始化并加载apollo配置
func loadApplicationConfig() error {
	// 1.传入初始化apollo所需参数
	// url := vars.EnvSetting.ApolloMetaServerUrl
	// if url == "" {
	// 	url = "http://localhost:18011"
	// }
	// apollo_sdk.NewApolloServer(url)
	// apolloLog := &apollo_sdk.ApolloLogger{Logger: log.BusinessLog}
	// err := apollo_sdk.ApolloServer.SetLog(apolloLog).Start(appName, "application") //"10.10.3.61:18011"
	// if err != nil {
	// 	return err
	// }

	// 2.获取apollo配置、本地配置
	err := vars.LoadApplicationConfig()
	if err != nil {
		return err
	}

	return nil
}

// 初始化log
func initLog() error {
	isDev := false
	if os.Getenv("GO_ENV") != "" {
		isDev = true
	}
	logger.Init(logger.Options{ProjectName: appName, IsDev: isDev})

	logPath := "/var/log/service"
	if vars.EnvSetting.IsDevelop {
		logPath = "./_log"
	}
	err := log.InitGlobalConfig(logPath, "error", appName)

	if err != nil {
		return err
	}

	err = log.GetAccessLogger("access")
	if err != nil {
		return err
	}

	err = log.GetBusinessLogger("business")
	if err != nil {
		return err
	}

	err = log.GetErrLogger("err")
	if err != nil {
		return err
	}

	return nil
}

// 初始化db配置
func initDbConfig() error {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:             vars.MysqlSetting.Host,
				Port:             vars.MysqlSetting.Port,
				User:             vars.MysqlSetting.UserName,
				Pass:             vars.MysqlSetting.Password,
				Name:             vars.MysqlSetting.DBName,
				Charset:          vars.MysqlSetting.Charset,
				MaxConnLifetime:  time.Duration(vars.MysqlSetting.MaxConnLifetime),
				MaxIdleConnCount: vars.MysqlSetting.MaxIdleConnCount,
				MaxOpenConnCount: vars.MysqlSetting.MaxOpenConnCount,
				Type:             "mysql",
			},
		},
		// "moon": gdb.ConfigGroup{
		// 	gdb.ConfigNode{
		// 		Host:             vars.MysqlSetting.Host,
		// 		Port:             vars.MysqlSetting.Port,
		// 		User:             vars.MysqlSetting.UserName,
		// 		Pass:             vars.MysqlSetting.Password,
		// 		Name:             vars.MysqlSetting.MoonDB,
		// 		Charset:          vars.MysqlSetting.Charset,
		// 		MaxConnLifetime:  time.Duration(vars.MysqlSetting.MaxConnLifetime),
		// 		MaxIdleConnCount: vars.MysqlSetting.MaxIdleConnCount,
		// 		MaxOpenConnCount: vars.MysqlSetting.MaxOpenConnCount,
		// 		Type:             "mysql",
		// 	},
		// },
	})
	if vars.EnvSetting.IsDevelop {
		g.DB().SetDebug(true)
	}

	return g.DB().PingMaster()
}

// 初始化mysql连接池 gorm
func initMysql() error {
	var err error
	args := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=true", vars.MysqlSetting.UserName, vars.MysqlSetting.Password, vars.MysqlSetting.Host, vars.MysqlSetting.Port, vars.MysqlSetting.DBName, vars.MysqlSetting.Charset)
	fmt.Println("mysql: ", args)
	vars.DB, err = apmgorm.Open("mysql", args)
	if err != nil {
		return err
	}

	if _, ok := os.LookupEnv("IS_DEVELOP"); ok {
		vars.DB.LogMode(true)
	}

	vars.DB.DB().SetMaxIdleConns(vars.MysqlSetting.MaxIdleConnCount)
	vars.DB.DB().SetMaxOpenConns(vars.MysqlSetting.MaxOpenConnCount)
	vars.DB.DB().SetConnMaxLifetime(time.Duration(vars.MysqlSetting.MaxConnLifetime))

	vars.DB.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	vars.DB.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	return nil
}

// 初始化redis配置
func initRedis() error {
	setting := vars.RedisSetting
	switch {
	case setting.MaxIdle == 0:
		setting.MaxIdle = 30
	case setting.MaxActive == 0:
		setting.MaxActive = 30
	case setting.IdleTimeout == 0:
		setting.IdleTimeout = 10
	}
	//host:port[,db,pass?maxIdle=x&maxActive=x&idleTimeout=x&maxConnLifetime=x]
	confString := fmt.Sprintf("%s:%s,%d,%s?maxIdle=%d&maxActive=%d&idleTimeout=%d",
		setting.Host,
		setting.Port,
		setting.Database,
		setting.Password,
		setting.MaxIdle,
		setting.MaxActive,
		setting.IdleTimeout,
	)
	err := gredis.SetConfigByStr(confString, "default")
	if err != nil {
		return err
	}
	_, err = g.Redis().Do("PING")
	if err != nil {
		return err
	}

	return nil
}

func initRabbitMQ() error {
	vhost := util.GetVhost()
	setting := vars.RabbitMQSetting

	// 初始化API配置
	rabbitmq.NewRabbitMQApi(rabbitmq.RbApiConfig{
		Host:     setting.Host,
		Port:     setting.Port,
		User:     setting.User,
		Password: setting.Password,
		ApiPort:  setting.ApiPort,
		Vhost:    vhost,
	})
	// 创建vhost 并授权
	_, err := rabbitmq.Api.AddVhost(vhost)
	if err != nil {
		return errors.Wrap(err, "vhost create fail")
	}
	_, err = rabbitmq.Api.SetPermissions(vhost, setting.User)
	if err != nil {
		return errors.Wrap(err, "set permissions fail")
	}

	// 初始rabbitmq服务
	err = rabbitmq.NewServer(rabbitmq.Config{
		Host:     setting.Host,
		User:     setting.User,
		Password: setting.Password,
		Port:     setting.Port,
		ApiPort:  setting.ApiPort,
		Vhost:    vhost,
		BackOff:  vars.RabbitMQSetting.BackOff,
	})
	if err != nil {
		return err
	}
	return nil
}
