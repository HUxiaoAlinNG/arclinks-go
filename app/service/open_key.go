/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 22:14:55
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:26:12
 * @FilePath: /arclinks-go/app/service/open_key.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"arclinks-go/app/model"
	"arclinks-go/library/log"
	"context"
	"errors"
	"sync"

	"arclinks-go/app/repository"
	"arclinks-go/app/repository/impl"
)

type OpenKey struct {
	OpenKeyRepository repository.OpenKeyRepositoryIface
}

var openKeyService *OpenKey
var openKeyOnce sync.Once

func NewOpenKeyService() *OpenKey {
	openKeyOnce.Do(func() {
		openKeyService = &OpenKey{
			OpenKeyRepository: impl.NewOpenKeyRepository(),
		}
	})

	return openKeyService
}

// 新增openKey
func (s *OpenKey) AddOpenKey(ctx context.Context, key string) error {
	item, _ := s.OpenKeyRepository.FindByKey(ctx, key)

	if item != nil {
		return errors.New("Same key already exists")
	}
	openKey := &model.OpenKeyModel{
		Key: key,
	}
	err := s.OpenKeyRepository.Insert(ctx, openKey)
	if err != nil {
		log.ErrLog.Errorf("AddOpenKey s.OpenKeyRepository.Insert err: %v", err)
		return errors.New("Insert openKey error")
	}

	return nil
}
