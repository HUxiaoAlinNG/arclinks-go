/*
 * @Author: hiLin 123456
 * @Date: 2022-12-19 18:52:35
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-19 20:24:54
 * @FilePath: /arclinks-go/app/ext_api/token.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ext_api

type TokenImageUrlReqItem struct {
	ContractAddress string `json:"contract_address"`
	TokenId         string `json:"token_id"`
	NetID           string `json:"net_id"`
}

type OpenApiGetAssetResponse struct {
	ImageUrl string `json:"image_url"`
}

type TokenImageUrlResItem struct {
	TokenImageUrlReqItem
	ImageUrl string `json:"image_url"`
}

type NFTScanTokenImageUrlResItem struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ImageUri string `json:"image_uri"`
	} `json:"data"`
}
