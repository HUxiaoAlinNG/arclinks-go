/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 17:36:23
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-18 17:58:23
 * @FilePath: /arclinks-go/app/model/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package model

type UserModel struct {
	Id            int64  `gorm:"column:id"`
	WalletAddress string `gorm:"column:wallet_address"` // 钱包地址
	NetID         string `gorm:"column:net_id"`         // 网络协议
	LoginTimes    int64  `gorm:"column:login_times"`    // 登陆次数
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`
}

func (m UserModel) TableName() string {
	return "user"
}
