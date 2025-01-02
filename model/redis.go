package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"happy-dog/utils"
	"happy-dog/utils/errmsg"
	"strconv"
	"sync"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client
var lock sync.Mutex

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     utils.RedisAddr,
		Password: utils.RedisPwd, // no password set
		DB:       utils.RedisDb,
	})

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("无法连接到 Redis:", err)
		return
	}
	fmt.Println("连接成功:", pong)
	//GetTest()
}

func StoreMessage(from, to int, message string) {
	lock.Lock()
	defer lock.Unlock()

	timestamp := time.Now().Unix()
	t1, t2, flag := swapFromTo(from, to)
	// 储存到redis
	err := rdb.RPush(ctx, fmt.Sprintf("%d_%d", t1, t2), fmt.Sprintf("%d_%d_%s", flag, timestamp, message)).Err()
	if err != nil {
		fmt.Println("无法存储消息:", err)
		return
	}

	length, err := rdb.LLen(ctx, fmt.Sprintf("%d_%d", t1, t2)).Result()
	if err != nil {
		fmt.Println("无法获取长度:", err)
		return
	}

	if length > 5 {
		fmt.Println("长度大于5", length)
		// 处理后续操作
		// 将数据存放到mysql中
		var rList []string
		for i := 0; i < 5; i++ {
			rdata, err := rdb.LPop(ctx, fmt.Sprintf("%d_%d", t1, t2)).Result()
			if err != nil {
				fmt.Println("无法获取消息:", err)
			}
			rList = append(rList, rdata)
		}
		// 将数据存放到mysql中
		code := AddMessage(DB, from, to, rList)
		if code == errmsg.ERROR {
			return
		}

	}
	return
	// 检查条数是否大于15

	// 如果大于15 则将所有的记录存放到mysql

	// 删除redis中的记录

}

func GetTest() {
	res, err := rdb.LPop(ctx, "20_21").Result()
	if err != nil {
		fmt.Println("无法获取消息:", err)
		return
	}
	fmt.Println("获取消息:", res)
	return
}

func swapFromTo(from, to int) (int, int, int) {
	if from > to {
		return to, from, 1
	}
	return from, to, 0
}

func InitConn() {
	//manage_socket_conn.GetUserSet()
}

func GetRedisHistory(id1, id2 int) ([]HistoryMessage, int) {
	historyMessage, err := rdb.LRange(ctx, fmt.Sprintf("%d_%d", id1, id2), 0, -1).Result()
	if err != nil {
		fmt.Println("无法获取历史消息:", err)
		return nil, errmsg.ERROR
	}

	var historyMessageList []HistoryMessage
	for _, v := range historyMessage {
		strFlag := v[0:1]
		flag, _ := strconv.Atoi(strFlag)

		strTime := v[2:12]
		strContent := v[13:]
		intTime, _ := strconv.Atoi(strTime)
		parsedTime := time.Unix(int64(intTime), 0)
		historyMessageList = append(historyMessageList, HistoryMessage{
			Id1:     id1,
			Id2:     id2,
			Flag:    flag,
			Time:    parsedTime,
			Content: strContent,
		})
	}

	return historyMessageList, errmsg.SUCCESS
}
