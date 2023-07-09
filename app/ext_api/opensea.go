/*
 * @Author: hiLin 123456
 * @Date: 2022-12-19 00:12:12
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-19 18:52:58
 * @FilePath: /arclinks-go/app/ext_api/opensea.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ext_api

import "context"

type OpenseaInstanceIface interface {
	GetTokenImageUrls(ctx context.Context, ops []*TokenImageUrlReqItem) ([]*TokenImageUrlResItem, error)
}
