/*
 * @Author: hiLin 123456
 * @Date: 2022-12-18 23:23:26
 * @LastEditors: hilin hilin
 * @LastEditTime: 2023-07-09 17:22:04
 * @FilePath: /arclinks-go/app/dto/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package dto

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
