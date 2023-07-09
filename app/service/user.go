/*
 * @Author: hiLin 123456
 * @Date: 2022-12-16 11:43:08
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 17:32:08
 * @FilePath: /arclinks-go/app/service/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"arclinks-go/app/common/err_msg"
	"arclinks-go/app/dto"
	"arclinks-go/library/log"
	"context"
	"sync"

	"arclinks-go/app/repository"
	"arclinks-go/app/repository/impl"
)

type User struct {
	UserRepository repository.UserRepositoryIface
}

var userService *User
var userOnce sync.Once

func NewUserService() *User {
	userOnce.Do(func() {
		userService = &User{
			UserRepository: impl.NewUserRepository(),
		}
	})

	return userService
}

func (s *User) Login(ctx context.Context, ops *dto.LoginReq) error {
	err := s.UserRepository.Save(ctx, ops.Username, ops.Password)

	if err != nil {
		log.ErrLog.Errorf("Login s.UserRepository.UserInfo err: %v", err)
		return err_msg.DBError
	}

	return nil
}
