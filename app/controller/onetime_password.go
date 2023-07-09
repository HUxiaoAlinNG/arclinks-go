/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 18:29:42
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 21:58:02
 * @FilePath: /arclinks-go/app/controller/onetime_password.go
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

type OnetimePassword struct {
}

func (c *OnetimePassword) VerifyPassword(r *ghttp.Request) {
	// 1.校验参数
	params := r.GetMap()
	rules := map[string]string{
		"password": "required",
	}
	msgs := map[string]interface{}{
		"password": "password is required",
	}
	if e := gvalid.CheckMap(params, rules, msgs); e != nil {
		response.ReturnErr(r, e)
	}

	ops := &dto.OneTimePasswordReq{}
	err := r.Parse(ops)
	if err != nil {
		response.ReturnErr(r, err)
	}

	openKey, err := service.NewOnetimePasswordService().VerifyPassword(r.Context(), ops.Password)

	if err != nil {
		response.Json(r, err_msg.FAILURE, err.Error())
		return
	}
	response.Json(r, err_msg.SUCCESS, openKey)
}

func (c *OnetimePassword) GenPassword(r *ghttp.Request) {
	// 1.校验参数
	params := r.GetMap()
	rules := map[string]string{
		"key": "required",
	}
	msgs := map[string]interface{}{
		"key": "key is required",
	}
	if e := gvalid.CheckMap(params, rules, msgs); e != nil {
		response.ReturnErr(r, e)
	}

	ops := &dto.GenPasswordReq{}
	err := r.Parse(ops)
	if err != nil {
		response.ReturnErr(r, err)
	}

	password, err := service.NewOnetimePasswordService().GenPassword(r.Context(), ops.Key)
	if err != nil {
		response.Json(r, err_msg.FAILURE, err.Error())
		return
	}
	response.Json(r, err_msg.SUCCESS, password)
}
