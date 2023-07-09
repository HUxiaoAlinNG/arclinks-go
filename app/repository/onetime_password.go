/*
 * @Author: hiLin 123456
 * @Date: 2021-11-09 16:41:18
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 17:39:32
 * @FilePath: /arclinks-go/app/repository/deployment.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package repository

import (
	"context"

	"arclinks-go/app/model"
)

type OneTimePasswordRepositoryIface interface {
	UpdateVerify(ctx context.Context, m *model.OneTimePasswordModel) error
	Insert(ctx context.Context, m *model.OneTimePasswordModel) error
	FindOne(ctx context.Context, password string) (*model.OneTimePasswordModel, error)
}
