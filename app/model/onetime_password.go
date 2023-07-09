/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 17:26:59
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:27:06
 * @FilePath: /arclinks-go/app/model/onetime_password.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package model

type OneTimePasswordModel struct {
	Id        int64  `gorm:"column:id"`
	OpenKeyId int64  `gorm:"column:open_key_id"`
	Password  string `gorm:"column:password"` // 一次性密码
	Verify    int64  `gorm:"column:verify"`   // 0 未被认证 1：已认证
}

func (m OneTimePasswordModel) TableName() string {
	return "onetime_password"
}
