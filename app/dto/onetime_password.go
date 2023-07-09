/*
 * @Author: hilin hilin
 * @Date: 2023-07-09 18:31:26
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 18:38:31
 * @FilePath: /arclinks-go/app/dto/onetime_password.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package dto

type OneTimePasswordReq struct {
	Password string `json:"password"`
}
type GenPasswordReq struct {
	Key string `json:"key"`
}
