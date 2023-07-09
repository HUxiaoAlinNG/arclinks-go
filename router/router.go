/*
 * @Author: hiLin 123456
 * @Date: 2021-11-09 16:41:18
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:23:45
 * @FilePath: /arclinks-go/router/router.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package router

import (
	"arclinks-go/app/controller"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 统一路由注册.
func init() {
	s := g.Server()

	s.Group("/api", func(group *ghttp.RouterGroup) {

		// 一次性密码
		onetime := new(controller.OnetimePassword)
		group.Group("/onetime", func(group *ghttp.RouterGroup) {
			group.GET("/getPsw", onetime.GenPassword)
			group.POST("/verify", onetime.VerifyPassword)
		})
		openkey := new(controller.OpenKey)
		group.Group("/openKey", func(group *ghttp.RouterGroup) {
			group.GET("/add", openkey.AddOpenKey)
		})
	})

}
