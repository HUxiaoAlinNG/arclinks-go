/*
 * @Author: hiLin 123456
 * @Date: 2021-11-09 16:41:18
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 18:29:53
 * @FilePath: /arclinks-go/app/controller/deployment.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"arclinks-go/app/common/err_msg"
	"arclinks-go/app/common/response"
	"arclinks-go/app/dto"
	"arclinks-go/app/service"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

type User struct {
}

func (c *User) Login(r *ghttp.Request) {
	// 1.校验参数
	params := r.GetMap()
	rules := map[string]string{
		"username": "required",
		"password": "required",
	}
	msgs := map[string]interface{}{
		"username": "username is required",
		"password": "password is required",
	}
	if e := gvalid.CheckMap(params, rules, msgs); e != nil {
		response.ReturnErr(r, e)
	}

	ops := &dto.LoginReq{}
	err := r.Parse(ops)
	if err != nil {
		response.ReturnErr(r, err)
	}

	err = service.NewUserService().Login(r.Context(), ops)

	// 3.返回结果集
	response.Json(r, err_msg.SUCCESS, "请求成功", err)
}
