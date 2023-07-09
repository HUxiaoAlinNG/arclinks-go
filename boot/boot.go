package boot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 用于应用初始化。
func init() {
	serverPort, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		serverPort = "9000"
	}
	port, err := strconv.Atoi(serverPort)
	if err != nil {
		fmt.Printf("strconv.Atoi(serverPort) err: %v", err)
		return
	}

	s := g.Server()
	// Web Server配置
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(false)
	s.SetPort(port)
}
