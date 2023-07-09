/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 17:42:06
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-19 00:50:18
 * @FilePath: /arclinks-go/app/repository/impl/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package impl

import (
	"arclinks-go/app/model"
	"arclinks-go/app/vars"
	"context"
	"errors"

	"github.com/gogf/gf/frame/g"
	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
)

type user struct {
	tableName string
}

func NewUserRepository() *user {
	return &user{
		tableName: "user",
	}
}

// 只需记录一次，有就保持
func (r user) Save(ctx context.Context, walletAddress string, netProtocol string) error {
	db := apmgorm.WithContext(ctx, vars.DB)
	m := &model.UserModel{}
	err := db.First(
		m,
		"wallet_address = ? AND net_id = ?",
		walletAddress,
		netProtocol).Error

	// 无则保持，有就更新
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.Create(&model.UserModel{
			WalletAddress: walletAddress,
			NetID:         netProtocol,
			LoginTimes:    1,
		}).Error

		if err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return err
	}

	err = db.Model(m).Update(g.Map{
		"login_times": m.LoginTimes + 1,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
