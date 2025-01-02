package manage_socket_conn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"sync"
)

// 这是一个管理上线下线的类
// 用map实现，并且需要对此进行并发限制
type userSet struct {
	users map[int]*websocket.Conn
	lock  sync.Mutex
	once  sync.Once
}
type msg struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Message string `json:"message"`
}

var us = new(userSet)

// 单例模式
func GetUserSet() *userSet {
	us.once.Do(func() {
		us.users = make(map[int]*websocket.Conn)
		us.users[-1] = nil
		us.lock = sync.Mutex{}
	})
	return us
}

// 添加链接
func (u *userSet) ConnConnect(user_id int, conn *websocket.Conn) (total, code int) {
	u.lock.Lock()
	defer u.lock.Unlock()

	if _, ok := u.users[user_id]; ok {
		// 说明已经在线了
		return len(u.users) - 1, errmsg.ERROR
	}
	// 链接成功
	u.users[user_id] = conn
	return len(u.users) - 1, errmsg.SUCCESS
}

// 链接断开
func (u *userSet) ConnDisconnect(user_id int) int {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.users[user_id]; ok {
		delete(u.users, user_id)
		return errmsg.SUCCESS
	}
	return errmsg.ERROR
}

// 转发消息
func (u *userSet) SendMsg(MessageData string, From, To int) int {
	u.lock.Lock()
	defer u.lock.Unlock()

	// 判断是否在线
	if _, ok := u.users[To]; !ok {
		return errmsg.ERROR
	}
	message := msg{
		From:    From,
		To:      To,
		Message: MessageData,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	err := enc.Encode(message)
	if err != nil {
		fmt.Println("Error encoding person:", err)
		return errmsg.ERROR
	}

	// 将编码后的数据转换为[]byte
	binaryData := buf.Bytes()

	err = u.users[To].WriteMessage(websocket.TextMessage, binaryData)
	if err != nil {
		return errmsg.ERROR
	}
	go func() {
		// 并发的写到redis中
		model.StoreMessage(From, To, MessageData)
	}()
	//fmt.Println("发送成功")
	return errmsg.SUCCESS
}
