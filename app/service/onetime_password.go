/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 14:56:06
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:45:08
 * @FilePath: /arclinks-go/app/service/onetime_password.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"arclinks-go/app/common/err_msg"
	"arclinks-go/app/model"
	"arclinks-go/library/log"
	"context"
	"errors"
	"fmt"
	"sync"

	"arclinks-go/app/repository"
	"arclinks-go/app/repository/impl"
	"time"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordLength = 10
)

type OnetimePassword struct {
	OneTimePasswordRepository repository.OneTimePasswordRepositoryIface
	OpenKeyRepository         repository.OpenKeyRepositoryIface
}

var onetimePasswordService *OnetimePassword
var onetimePasswordOnce sync.Once

func NewOnetimePasswordService() *OnetimePassword {
	onetimePasswordOnce.Do(func() {
		onetimePasswordService = &OnetimePassword{
			OneTimePasswordRepository: impl.NewOnetimePasswordRepository(),
			OpenKeyRepository:         impl.NewOpenKeyRepository(),
		}
	})

	return onetimePasswordService
}

// 验证一次性并返回对应 openKey
func (s *OnetimePassword) VerifyPassword(ctx context.Context, password string) (string, error) {
	item, err := s.OneTimePasswordRepository.FindOne(ctx, password)
	if err != nil {
		log.ErrLog.Errorf("VerifyPassword s.OneTimePasswordRepository.FindOnes err: %v", err)
		return "", err_msg.DBError
	}
	if item == nil {
		return "", errors.New("Invalid password")
	}
	if item.Verify == 1 {
		return "", errors.New("The password has been verified and cannot be used again.")
	}

	item.Verify = 1

	openKey, err := s.OpenKeyRepository.FindById(ctx, item.OpenKeyId)
	if err != nil {
		log.ErrLog.Errorf("VerifyPassword s.OpenKeyRepository.FindOne err: %v", err)
		return "", errors.New("Find OpenKey error")
	}

	if openKey == nil {
		return "", errors.New("No OpenKey")
	}

	err = s.OneTimePasswordRepository.UpdateVerify(ctx, item)
	if err != nil {
		log.ErrLog.Errorf("VerifyPassword s.OneTimePasswordRepository.UpdateVerify err: %v", err)
		return "", errors.New("Update Password error")
	}

	return openKey.Key, nil
}

// 生成一次性密码
func (s *OnetimePassword) GenPassword(ctx context.Context, openKey string) (string, error) {
	openKeyModel, err := s.OpenKeyRepository.FindByKey(ctx, openKey)
	if err != nil {
		log.ErrLog.Errorf("GenPassword s.OpenKeyRepository.FindByKey err: %v", err)
		return "", err_msg.DBError
	}

	if openKeyModel == nil {
		return "", errors.New("Can't find openKey Record")
	}

	password := fmt.Sprintf("arclinks-%v", time.Now().UnixMicro())
	oneTimePassword := &model.OneTimePasswordModel{
		Password:  password,
		OpenKeyId: openKeyModel.Id,
		Verify:    0,
	}

	err = s.OneTimePasswordRepository.Insert(ctx, oneTimePassword)
	if err != nil {
		log.ErrLog.Errorf("GenPassword s.OneTimePasswordRepository.Insert err: %v", err)
		return "", errors.New("Insert password error")
	}

	return password, nil
}
