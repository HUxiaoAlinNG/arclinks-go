/*
 * @Author: hiLin 123456
 * @Date: 2021-11-09 16:41:18
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 22:12:37
 * @FilePath: /arclinks-go/app/common/util/pager.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package util

import (
	"strconv"
	"strings"

	"arclinks-go/app/vars"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

type Pager struct {
	Page     int `p:"page_index"`
	PageSize int `p:"page_size"`
}

func (p *Pager) Offset() int {
	return p.PageSize * (p.Page - 1)
}

func GetPageSize(pageSize int) int {
	if pageSize <= 0 {
		return vars.AppSetting.DefaultPageSize
	}
	if pageSize > vars.AppSetting.MaxPageSize {
		return vars.AppSetting.MaxPageSize
	}
	return pageSize
}

func GetPage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func GetPager(page, pageSize int) *Pager {
	return &Pager{
		Page:     GetPage(page),
		PageSize: GetPageSize(pageSize),
	}
}

func InitPager(r *ghttp.Request) *Pager {
	return GetPager(r.GetInt("page_index"), r.GetInt("page_size"))
}

func ParseIdsFromStr(idStr string) (ids []int64) {
	if idStr == "" {
		return
	}
	strArr := strings.Split(idStr, ",")
	return gconv.SliceInt64(strArr)
}

func ParseUint8FromStr(str string) (res []uint8) {
	if str == "" {
		return
	}
	strArr := strings.Split(str, ",")
	for _, s := range strArr {
		i, _ := strconv.ParseUint(s, 10, 8)
		res = append(res, uint8(i))
	}
	return res
}
