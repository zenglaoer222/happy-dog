package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	MG "happy-dog/manage_socket_conn"
	"happy-dog/middleware"
	"happy-dog/utils/errmsg"
)

type Message struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Message string `json:"message"`
}

// 查看消息发送请求调用连接
func ConCreate(c *gin.Context) {
	// 初始一个连接
	var conn *websocket.Conn
	var err error

	if websocket.IsWebSocketUpgrade(c.Request) {
		// 升级为websocket连接
		conn, err = websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
		if err != nil {
			return
		}
		println("升级为websocket连接")
	} else {
		return
	}

	var uintUserId uint

	// todo:验证其token是否合法
	messageType, messageData, err := conn.ReadMessage()
	if err != nil || messageType != 1 {
		return
	}

	var data struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(messageData, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	claims, err := middleware.ParseToken(data.Token)
	if err != nil {
		return
	}
	uintUserId = claims.UserId

	total, code := MG.GetUserSet().ConnConnect(int(uintUserId), conn)
	if code == errmsg.ERROR {
		println("出错")
		return
	}
	println("正常，当前连接数：", total)
	// todo:返回历史信息
	//model.GetMessage()
	// 先采用redis

	go func() {
		defer func() {
			MG.GetUserSet().ConnDisconnect(int(uintUserId))
		}()

		for {
			messageType, messageData, err = conn.ReadMessage()
			println("收到消息", messageType)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					println("正常关闭")
				} else {
					println("异常断开")
				}
				break
			}
			if messageType == websocket.CloseMessage {
				println("收到关闭消息")
				break
			} else if messageType == websocket.TextMessage {
				println("收到文本消息", string(messageData))

				var msg Message
				err = json.Unmarshal(messageData, &msg)
				if err != nil {
					println("解析失败")
					return
				}
				fmt.Println(msg.Message, msg.From, msg.To)
				// todo:转发

				if msg.To == 0 {
					return
				}
				MG.GetUserSet().SendMsg(msg.Message, msg.From, msg.To)
			}

		}
	}()

	return

}

func ConDelete(c *gin.Context) {
	userId, _ := c.Get("user_id")
	uintUserId, _ := userId.(uint)
	// 关闭当前用户连接
	code := MG.GetUserSet().ConnDisconnect(int(uintUserId))
	if code == errmsg.ERROR {
		println("断开失败")
		return
	}
	println("断开成功")
	return
}
