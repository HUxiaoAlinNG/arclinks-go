/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 23:47:17
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-19 17:25:01
 * @FilePath: /arclinks-go/app/ext_api/impl/opensea.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package impl

import (
	"arclinks-go/app/ext_api"
	"arclinks-go/library/httpctx"
	"arclinks-go/library/log"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var OpenseaApiUrl = "https://api.opensea.io"

type OpenseaApiInstance struct {
}

var openseaApiInstanceService *OpenseaApiInstance
var openseaApiInstanceOnce sync.Once

func NewOpenseaApiInstanceService() *OpenseaApiInstance {
	openseaApiInstanceOnce.Do(func() {
		openseaApiInstanceService = &OpenseaApiInstance{}
	})

	return openseaApiInstanceService
}

func (i *OpenseaApiInstance) GetTokenImageUrls(ctx context.Context, ops []*ext_api.TokenImageUrlReqItem) ([]*ext_api.TokenImageUrlResItem, error) {

	wg := sync.WaitGroup{}
	wg.Add(len(ops))
	results := make([]*ext_api.TokenImageUrlResItem, 0)
	errs := make([]error, 0)
	for _, item := range ops {
		go func(item *ext_api.TokenImageUrlReqItem) {
			defer wg.Done()

			url := fmt.Sprintf("%s/api/v1/asset/%s/%s/?include_orders=false", OpenseaApiUrl, item.ContractAddress, item.TokenId)

			resp, err := httpctx.Get(
				ctx,
				url,
				map[string]string{},
			)

			if err != nil {
				log.ErrLog.Errorf("GetTokenImageUrls httpctx.Get %s err: %v", url, err)
				errs = append(errs, err)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.ErrLog.Errorf("GetTokenImageUrls httpctx.Get %s request failed. err: %v", url, resp.Status)
				errs = append(errs, fmt.Errorf("GetTokenImageUrls httpctx.Get %s request failed. err: %v", url, resp.Status))
				return
			}

			bodyString, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.ErrLog.Errorf("GetTokenImageUrls httpctx.Get %s ioutil.ReadAll err: %v", url, err)
				errs = append(errs, err)
				return
			}

			body := &ext_api.OpenApiGetAssetResponse{}
			err = json.Unmarshal(bodyString, body)

			if err != nil {
				log.ErrLog.Errorf("GetTokenImageUrls httpctx.Get %s json.Unmarshal err: %v", url, err)
				errs = append(errs, err)
				return
			}

			res := &ext_api.TokenImageUrlResItem{
				ImageUrl: body.ImageUrl,
			}
			res.ContractAddress = item.ContractAddress
			res.NetID = item.NetID
			res.TokenId = item.TokenId

			results = append(results, res)

		}(item)
	}
	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return results, nil
}
