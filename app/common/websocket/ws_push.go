/*
 * @Author: hiLin 123456
 * @Date: 2021-11-08 17:18:56
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 21:19:29
 * @FilePath: /./app/common/websocket/ws_push.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package websocket

import (
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"

	"github.com/gogf/gf/net/ghttp"
)

type Msg struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	// 消息类型
	OK            = "msg"           // "message"
	SvcErr        = "err"           // 业务错误
	InvalidParams = "invalidParams" // 参数有误
)

const (

	// 心跳间隔
	HEARTBEAT_DURATION = 5 * time.Minute
)

// 响应客户端心跳验证, 维持连接可用
func ResponseToClientsHeartBeat(ws *ghttp.WebSocket) {
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

// 关闭websocket连接，向参数中的close channel写入关闭消息
func closeConnection(ws *ghttp.WebSocket, closeSignal chan<- struct{}) {
	ws.Close()
	closeSignal <- struct{}{}
}

// 处理消息推送
func HandleMessagePush(ws *ghttp.WebSocket, dataCh chan interface{}, closeSignal chan<- struct{}) {
	ticker := time.NewTicker(HEARTBEAT_DURATION)
	defer ticker.Stop()
	defer func() {
		recover()
		return
	}()

ForEnd:
	for {

		select {
		// 从service层获得的消息
		case msgData := <-dataCh:

			data := msgData.(Msg)
			if data.Type != OK {
				_ = ws.WriteMessage(ghttp.WS_MSG_CLOSE, gorillaWebsocket.FormatCloseMessage(gorillaWebsocket.CloseInternalServerErr, data.Message))
				closeConnection(ws, closeSignal)
				break ForEnd
			}

			err := ws.WriteJSON(data)
			if err != nil {
				closeConnection(ws, closeSignal)
				break ForEnd
			}

		// 心跳监测
		case <-ticker.C:
			err := ws.WriteMessage(ghttp.WS_MSG_TEXT, []byte(""))

			if err != nil {
				closeConnection(ws, closeSignal)
				break ForEnd
			}
		}
	}
}
