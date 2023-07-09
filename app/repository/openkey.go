/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 18:02:01
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:16:18
 * @FilePath: /arclinks-go/app/repository/openkey.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package repository

import (
	"arclinks-go/app/model"
	"context"
)

type OpenKeyRepositoryIface interface {
	FindById(ctx context.Context, id int64) (*model.OpenKeyModel, error)
	FindByKey(ctx context.Context, openKey string) (*model.OpenKeyModel, error)
	Insert(ctx context.Context, m *model.OpenKeyModel) error
}
