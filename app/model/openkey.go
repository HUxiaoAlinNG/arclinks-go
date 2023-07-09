/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 17:25:41
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 18:15:03
 * @FilePath: /arclinks-go/app/model/openkey.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package model

type OpenKeyModel struct {
	Id  int64  `gorm:"column:id"`
	Key string `gorm:"column:key"` // 钱包地址
}

func (m OpenKeyModel) TableName() string {
	return "open_key"
}
