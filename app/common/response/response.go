/*
 * @Author: hiLin 123456
 * @Date: 2021-11-08 17:18:56
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 22:12:14
 * @FilePath: /arclinks-go/app/common/response/response.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package response

import (
	"arclinks-go/app/common/err_msg"

	"arclinks-go/library/log"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
)

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// code:  错误码(0:成功, 其它:失败);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func Json(r *ghttp.Request, code int, msg string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":    code,
		"message": msg,
		"data":    responseData,
	})
	if code == err_msg.SUCCESS {
		end_time := gtime.TimestampMicro()
		log.AccessLog.Infof("uri:%s, start_time:%v(Microsecond), end_time:%v(Microsecond), total_time:%v(s)", r.RequestURI, r.EnterTime, end_time, float64(end_time-r.EnterTime)*1e-06)
	}
	r.Exit()
}

func ReturnErr(r *ghttp.Request, err error, others ...map[string]interface{}) {
	log.ErrLog.Infof("params get:%s, post:%s", r.RequestURI, string(r.GetRaw()))
	r.Response.Header().Set("Content-Type", "application/json")
	res := g.Map{"code": err_msg.FAILURE, "message": err.Error()}

	if len(others) > 0 {
		for _, other := range others {
			for key, value := range other {
				res[key] = value
			}
		}
	}

	r.Response.WriteJson(res)
	r.Exit()
}

func Return500Err(r *ghttp.Request, err error, others ...map[string]interface{}) {
	log.ErrLog.Infof("params get:%s, post:%s", r.RequestURI, string(r.GetRaw()))
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.Status = 500
	res := g.Map{"code": err_msg.FAILURE, "message": err.Error()}

	if len(others) > 0 {
		for _, other := range others {
			for key, value := range other {
				res[key] = value
			}
		}
	}

	r.Response.WriteJson(res)
	r.Exit()
}
