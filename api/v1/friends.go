package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"strconv"
)

type FriendsWithData struct {
	Friends        model.Friends          `json:"friends"`
	HistoryMessage []model.HistoryMessage `json:"history_message"`
}

func CreateFriends(c *gin.Context) {
	userId1, err := strconv.Atoi(c.PostForm("id1"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "格式有误",
		})
		return
	}
	userId2, err := strconv.Atoi(c.PostForm("id2"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "格式有误",
		})
		return
	}
	//var friends model.Friends
	//_ = c.ShouldBindJSON(&friends)
	//userId1 := friends.ID1
	//userId2 := friends.ID2
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "创建好友失败",
		})
		return
	}
	//println(userId1, user_id.(uint))
	if userId1 != int(user_id.(uint)) {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "非本人操作好友",
		})
		return
	}

	// 判断是否已经是好友或者已经添加过
	mode, _ := model.ExistFriends(model.DB, uint(userId1), uint(userId2))
	if mode == 1 || mode == 2 {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "已经是好友或者已经添加过",
		})
		return
	}

	code := model.CreateFriends(model.DB, uint(userId1), uint(userId2))
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "创建失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(200),
	})
	return
}

func GetFriends(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取token失败",
		})
		return
	}
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取user_id失败",
		})
		return
	}
	if userId != int(user_id.(uint)) {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "非本人操作",
		})
		return
	}

	friends, code := model.GetFriends(model.DB, user_id.(uint), 2)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "查找好友失败",
		})
		return
	}

	// todo:redis返回每个好友的聊天记录
	// q:
	var friendsWithData []FriendsWithData

	for i := 0; i < len(friends); i++ {
		var friendWithData FriendsWithData
		friendWithData.Friends = friends[i]
		friendWithData.HistoryMessage, code = model.GetMessage(model.DB, int(friends[i].ID1), int(friends[i].ID2))
		if code == errmsg.ERROR {
			c.JSON(200, gin.H{
				"code":    errmsg.ERROR,
				"message": "聊天记录查找失败",
			})
			return
		}

		friendsWithData = append(friendsWithData, friendWithData)
	}

	c.JSON(200, gin.H{
		"code":    errmsg.SUCCESS,
		"data":    friendsWithData,
		"message": "查找好友成功",
	})
	return
}

func GetFriendsWait(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取token失败",
		})
		return
	}
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取user_id失败",
		})
		return
	}
	if userId != int(user_id.(uint)) {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "非本人操作",
		})
		return
	}

	friends, code := model.GetFriends(model.DB, user_id.(uint), 1)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "查找好友失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    errmsg.SUCCESS,
		"data":    friends,
		"message": "查找好友成功",
	})
	return
}

func AcceptFriends(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取token失败",
		})
		return
	}
	Id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取参数id失败",
		})
		return
	}
	code := model.BoolFriends(model.DB, uint(Id), user_id.(uint))
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "非本人待添加好友",
		})
		return
	}

	code = model.AcceptFriends(model.DB, uint(Id))
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "添加好友失败",
		})
		return

	}
	c.JSON(200, gin.H{
		"code":    errmsg.SUCCESS,
		"message": "添加好友成功",
	})
	return
}

func SearchFriends(c *gin.Context) {

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取token失败",
		})
		return
	}

	username := c.Query("username")

	customerAll, total, code := model.GetCustomer(model.DB, username, -1, 0)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "查找好友失败",
		})
		return
	}
	println(total, len(customerAll))
	var customerRes []model.Customer

	for _, item := range customerAll {
		id1 := item.ID
		id2 := user_id.(uint)
		println(id1, id2)
		if id1 != id2 {
			mode, _ := model.ExistFriends(model.DB, id1, id2)

			//println(id1, id2)

			if mode == 0 {

				customerRes = append(customerRes, item)
			}
		}

	}

	c.JSON(200, gin.H{
		"code":  errmsg.SUCCESS,
		"data":  customerRes,
		"total": len(customerRes),
	})
	return

}
