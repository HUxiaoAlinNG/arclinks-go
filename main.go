/*
 * @Author: hiLin 123456
 * @Date: 2021-11-08 17:18:56
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 23:19:55
 * @FilePath: /./main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"

	_ "arclinks-go/boot"
	"arclinks-go/config"
	_ "arclinks-go/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Println("config.InitConfig err: ", err)
		return
	}

	// 启动MQ的消费端
	// go rabbitmq.Server.RunReceiver()

	g.Server().EnablePProf()
	g.Server().SetLogStdout(true)
	g.Server().SetAccessLogEnabled(true)
	g.Server().Run()
}
