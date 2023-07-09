/*
 * @Author: hiLin 123456
 * @Date: 2022-12-16 11:43:08
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:46:13
 * @FilePath: /arclinks-go/app/repository/impl/role.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package impl

import (
	"context"
	"errors"

	"arclinks-go/app/model"
	"arclinks-go/app/vars"

	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
)

type onetimePassword struct {
	tableName string
}

func NewOnetimePasswordRepository() *onetimePassword {
	return &onetimePassword{
		tableName: "onetime_password",
	}
}

func (r onetimePassword) FindOne(ctx context.Context, password string) (*model.OneTimePasswordModel, error) {
	db := apmgorm.WithContext(ctx, vars.DB)
	m := &model.OneTimePasswordModel{}

	err := db.First(
		m,
		"password = ?",
		password).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

// 更新verify
func (r onetimePassword) UpdateVerify(ctx context.Context, m *model.OneTimePasswordModel) error {
	db := apmgorm.WithContext(ctx, vars.DB)
	err := db.First(
		&model.OneTimePasswordModel{},
		"id = ?", m.Id).Update("verify", m.Verify).Error

	if err != nil {
		return err
	}

	return nil
}

func (r onetimePassword) Insert(ctx context.Context, m *model.OneTimePasswordModel) error {
	db := apmgorm.WithContext(ctx, vars.DB)
	err := db.Save(m).Error

	if err != nil {
		return err
	}

	return nil
}
