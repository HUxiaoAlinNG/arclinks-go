# go-giro

## 框架结构  
```
arclinks-go
├── app                         // 应用程序代码
│   ├── auth                    // 鉴权模块
│   ├── common                  // 应用程序内部公用函数(系统相关但不涉及具体的业务逻辑)
│   │   ├── err_msg             // 错误统一处理
│   │   ├── response            // 返回值统一处理
│   │   └── util                // 其他公用模块
│   ├── controller              // 控制器
│   ├── model                   // model(只定义数据库表相关的结构体，具体操作放到repository里面)
│   ├── repository              // 仓储层
│   │   ├── impl                // 仓储层接口实现
│   │   ├── iteration.go        // 仓储层接口
│   │   └── mock                // mock
│   ├── service                 // 核心业务逻辑
│   │   ├── iteration.go        // 业务代码
│   │   └── iteration_test.go   // 单元测试
│   ├── ext_api                 // 第三方系统 api 接口
│   │   ├── impl                // 第三方api接口实现
│   │   ├── st2_client.go       // 第三方api接口
│   ├── vars                    // 通用结构体
│   │   ├── setting.go          // apollo 配置对应的结构体
│   │   └── vars.go             // 程序通用结构体
│   ├── dto                     // controller与service层数据交互结构体以及复杂请求参数校验
├── boot                        // 应用初始化
│   └── boot.go
├── config                      // 初始化配置
│   └── config.go
├── router                      // 路由规则
│   └── router.go
├── library                     // 公共库(封装第三方库)
│   ├── apollo_sdk              // apollo SDK
│   ├── queue                   // 消息队列
│   ├── log                     // log SDK
│   └── rabbitmq                // rabbitMQ
├── doc                         // 文档
│   └── workflow                // st2 工作流相关文档
├── main.go                     // 入口文件
├── go.mod
├── go.sum
└── README.MD
```

> * library:系统公用方法，拿到别的系统也可以直接使用的模块。
> * library/apollo_sdk:apollo SDK，使用详情见下面"如何使用apollo添加配置"。
> * library/log:log SDK，使用详情见下面"log使用"。
> * auth:只处理权限控制，后期考虑独立系统。用户项目权限不放到这个目录。
> * router:路由配置(路由组)&通过中间件鉴权控制。
> * controller:控制器层，只包含校验参数、调用service层处理业务逻辑、返回结果集。
> * service:处理核心业务逻辑。
> * repository:所有操作数据库的方法均封装在这一层，其它层不直接操作数据库。
> * model:数据库对应的结构体，一个表一个结构体。
> * ext_api:对接第三方api统一入口，不做业务逻辑处理

## GF框架教程
```
https://goframe.org/start/index
```

## model模块字段类型
| mysql字段类型 | model字段类型 |
| --------      | :----:        |
| tinyint       | uint8         |
| smallint      | uint32        |
| int           | int64         |
| bigint        | uint64        |
| *char	        | string        |
| *text	        | string        |
| *binary       | bytes         |
| *blob	        | bytes         |
| float         | float64       |
| double	    | float64       |
| decimal       | float64       |
| bool          | bool          |
| 其他          | string        |

## 如何使用apollo添加配置
1.设置连接参数&初始化apollo，应用初始化时已设置
```
// 设置环境变量: APOLLO_META_SERVER_URL

func loadApolloConfig() error {
	// 1.传入初始化apollo所需参数
	apollo_sdk.NewApolloServer(os.Getenv("APOLLO_META_SERVER_URL"))
	err := apollo_sdk.ApolloServer.Start(appName, "application")
	if err != nil {
		return err
	}

	// 2.获取apollo配置的参数
	err = vars.LoadApplicationConfig()
	if err != nil {
		return err
	}
	return nil
}
```

2.在app/vars/setting.go文件中添加对应的结构体，以及配置加载映射即可
```
// 后续新增配置，只需要在这个文件中创建对应的结构体，初始化，LoadApplicationConfig()配置加载映射，在其他模块直接通过vars.AppSetting.DefaultPageSize调用
type AppSettingS struct {
	DefaultPageSize int64
	MaxPageSize     int64
}

var (
	AppSetting = &AppSettingS{}
)

// 加载配置对象映射
func LoadApplicationConfig() error {
    var config *gcfg.Config
	if vars.EnvSetting.IsDevelop { //判断是否开启了本地配置
		config = g.Config()
	}
	err := loadConfig(config, "App", AppSetting)
	if err != nil {
		return err
	}

	return nil
}
```

3.在apollo配置对应的参数，并发布

4.设置本地配置(若开启了本地配置，Apollo配置会被覆盖)
```
(1)配置环境变量: export IS_DEVELOP=1;
(2)./config目录下增加config.toml文件，并添加相应的配置，eg:
[App]
    DefaultPageSize      = 12
    MaxPageSize          = 100
```

## 如何使用mysql
```
// 在应用程序中使用g.DB()对象，g.DB()对象是一个单例对象
func (iter *iteration) GetIterationList() ([]*model.IterationModel, error) {
	list := []*model.IterationModel{}

	err := g.DB().Table("iteration").Where("id>1").Structs(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
```

## 如何使用redis
```
import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/util/gconv"
)

value, err := g.Redis().Do("GET", "key_name")
if err != nil {
	fmt.Println("err: ", err)
}

fmt.Println(gconv.String(value))
```

## log使用
```
// 系统内置了三类log，分别为access、business、error三类，使用demo如下：
import "./library/log"

log.ErrLog.Errorf("参数异常, %s", param)

log.BusinessLog.Infof("请求参数为%s", param)

log.AccessLog.Info("access")
```

## app/service层业务代码，需提供单元测试
```
单元测试主要用的gomock、GoConvey
https://www.jianshu.com/p/598a11bbdafb
https://blog.csdn.net/zwqjoy/article/details/79474196
```

```
// 1.app/repository/iteration.go 定义接口

// 2.app/repository/impl/ 实现接口

// 3.在app/repository/目录生成mock文件
$ go generate .
或者 
$ make gen-mock

// 4.编辑单元测试文件。eg:app/service/iteration_test.go

// 5.运行单元测试方法
$ go test -v -run TestGetIterationList .
=== RUN   TestGetIterationList

  db err ✔

1 total assertion

--- PASS: TestGetIterationList (0.00s)
PASS
ok  	./app/service	0.232s

或者 
$ make test

// 6.查看测试代码覆盖率
go test -v -coverprofile=cover.out
go tool cover -html=cover.out -o coverage.html
最终生成一个html文件，浏览器直接打开即可 
```

## 安装说明
由于引入github.com/containers/image(cgo)库实现镜像同步功能，需要安装依赖程序gpgme, btrfs, devmapper:
```bash
#ubuntu
sudo apt-get install -y libgpgme11-dev
sudo apt-get install -y libbtrfs-dev
sudo apt-get install -y libdevmapper-dev

#macos
brew install gpgme
```
如果出现libdevmapper依赖版本冲突需要手动下载安装 [devmappper](https://packages.ubuntu.com/bionic-updates/libdevmapper-dev)
