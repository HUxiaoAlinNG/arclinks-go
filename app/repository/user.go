/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 17:40:41
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 18:31:49
 * @FilePath: /arclinks-go/app/repository/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package repository

import (
	"context"
)

type UserRepositoryIface interface {
	Save(ctx context.Context, username string, password string) error
}
