/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 22:12:29
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 22:13:38
 * @FilePath: /arclinks-go/app/dto/open_key.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package dto

type AddOpenKeyReq struct {
	Key string `json:"password"`
}
