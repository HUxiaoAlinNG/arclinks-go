/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 23:47:17
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-19 20:41:20
 * @FilePath: /arclinks-go/app/ext_api/impl/nftScan.go
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

var (
	X_API_KEY = "9twmdU9j6brxZckEr8pdFFop"
)

type NFTScanApiInstance struct {
}

var nftScanApiInstanceService *NFTScanApiInstance
var nftScanApiInstanceOnce sync.Once

func NewNFTScanApiInstanceService() *NFTScanApiInstance {
	nftScanApiInstanceOnce.Do(func() {
		nftScanApiInstanceService = &NFTScanApiInstance{}
	})

	return nftScanApiInstanceService
}

func (i *NFTScanApiInstance) GetTokenImageUrls(ctx context.Context, ops []*ext_api.TokenImageUrlReqItem) ([]*ext_api.TokenImageUrlResItem, error) {

	wg := sync.WaitGroup{}
	wg.Add(len(ops))
	results := make([]*ext_api.TokenImageUrlResItem, 0)
	errs := make([]error, 0)
	for _, item := range ops {
		go func(item *ext_api.TokenImageUrlReqItem) {
			defer wg.Done()
			host, err := i.getHost(item.NetID)
			if err != nil {
				log.ErrLog.Errorf("GetTokenImageUrls i.getHost %s err: %v", item.NetID, err)
				errs = append(errs, err)
				return
			}
			url := fmt.Sprintf("%s/api/v2/assets/%s/%s?show_attribute=true", host, item.ContractAddress, item.TokenId)

			resp, err := httpctx.Get(
				ctx,
				url,
				map[string]string{
					"X-API-KEY": X_API_KEY,
				},
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

			body := &ext_api.NFTScanTokenImageUrlResItem{}
			err = json.Unmarshal(bodyString, body)

			if err != nil {
				log.ErrLog.Errorf("GetTokenImageUrls httpctx.Get %s json.Unmarshal err: %v", url, err)
				errs = append(errs, err)
				return
			}

			res := &ext_api.TokenImageUrlResItem{
				ImageUrl: fmt.Sprintf("https://ipfs.io/ipfs/%s", body.Data.ImageUri),
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

func (i *NFTScanApiInstance) getHost(netID string) (string, error) {
	switch netID {
	case "1": // ETH
		return "https://restapi.nftscan.com", nil
	case "137": // BNB
		return "https://bnbapi.nftscan.com", nil
	case "56": // Polygon
		return "https://polygonapi.nftscan.com", nil
	case "42161": // Arbitrum One
		return "https://arbitrumapi.nftscan.com", nil
	case "10": // Optimism
		return "https://optimismapi.nftscan.com", nil
	case "43114": // Avalanche
		return "https://avaxapi.nftscan.com", nil
	default:
		return "", fmt.Errorf("不支持NetID为%s查询token信息", netID)
	}
}
