/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 18:02:50
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:16:13
 * @FilePath: /arclinks-go/app/repository/impl/openkey.go
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

type openKey struct {
	tableName string
}

func NewOpenKeyRepository() *openKey {
	return &openKey{
		tableName: "open_key",
	}
}

func (r openKey) FindById(ctx context.Context, id int64) (*model.OpenKeyModel, error) {
	db := apmgorm.WithContext(ctx, vars.DB)
	m := &model.OpenKeyModel{}

	err := db.First(
		m,
		"id = ?",
		id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r openKey) FindByKey(ctx context.Context, openKey string) (*model.OpenKeyModel, error) {
	db := apmgorm.WithContext(ctx, vars.DB)
	m := &model.OpenKeyModel{}

	err := db.First(
		m,
		"`key` = ?",
		openKey).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r openKey) Insert(ctx context.Context, m *model.OpenKeyModel) error {
	db := apmgorm.WithContext(ctx, vars.DB)
	err := db.Save(m).Error

	if err != nil {
		return err
	}

	return nil
}
