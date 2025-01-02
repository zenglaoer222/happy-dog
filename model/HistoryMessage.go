package model

import (
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
	"strconv"
	"time"
)

type HistoryMessage struct {
	gorm.Model
	Id1     int       `json:"id1"`
	Id2     int       `json:"id2"`
	Flag    int       `json:"flag"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func GetMessage(tx *gorm.DB, id1 int, id2 int) (historyMessage []HistoryMessage, code int) {
	t1 := id1
	t2 := id2
	if id1 > id2 {
		t1 = id2
		t2 = id1
	}
	err := tx.Where("id1 = ? AND id2 = ?", t1, t2).Find(&historyMessage).Error
	if err != nil {
		return nil, errmsg.ERROR
	}

	// todo:调用redis读取
	redisHistory, code := GetRedisHistory(t1, t2)
	if code == errmsg.ERROR {
		return nil, errmsg.ERROR
	}
	historyMessage = append(historyMessage, redisHistory...)

	return historyMessage, errmsg.SUCCESS
}

func AddMessage(tx *gorm.DB, id1 int, id2 int, content []string) int {
	for _, v := range content {
		// 解析时间戳和内容
		strFlag := v[0:1]
		flag, _ := strconv.Atoi(strFlag)

		strTime := v[2:12]
		println("strTime", strTime)
		strContent := v[13:]
		println("strContent", strContent)
		intTime, _ := strconv.Atoi(strTime)
		parsedTime := time.Unix(int64(intTime), 0)
		err := tx.Create(&HistoryMessage{
			Id1:     id1,
			Id2:     id2,
			Flag:    flag,
			Content: strContent,
			Time:    parsedTime,
		}).Error
		if err != nil {
			return errmsg.ERROR
		}
	}
	return errmsg.SUCCESS
}
