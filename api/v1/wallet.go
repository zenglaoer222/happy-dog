package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"happy-dog/utils/validator"
	"strconv"
)

// 创建钱包
func CreateWallet(c *gin.Context) {
	var data model.Wallet
	_ = c.ShouldBindJSON(&data)
	//cid, _ := strconv.Atoi(c.Query("cid"))

	msg, code := validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(200, gin.H{
			"status": code,
			"msg":    msg,
		})
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    "user_id获取失败",
		})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    "role获取失败",
		})
		return
	}

	fmt.Println("确认用户身份")
	fmt.Println(userId, data.Cid)
	if data.Cid != userId.(uint) && role != "manager" { // 检查是否为当前用户
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}

	fmt.Println("查询钱包")
	code, _ = model.InquireBalance(model.DB, data.Cid)
	if code == errmsg.SUCCESS {
		code = errmsg.ERROR_WALLET_EXIST
		c.JSON(200, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	fmt.Println("创建钱包")

	code = model.CreateWallet(model.DB, data.Cid, data.Password)
	c.JSON(200, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
	return
}

// 查询余额
func InquireBalance(c *gin.Context) {
	var msg string
	cid, _ := strconv.Atoi(c.Query("cid"))

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    "user_id获取失败",
		})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    "role获取失败",
		})
		return
	}

	if uint(cid) != userId.(uint) && role != "manager" { // 检查是否为当前用户
		c.JSON(200, gin.H{
			"status": errmsg.ERROR,
			"msg":    errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}

	code, balance := model.InquireBalance(model.DB, uint(cid))
	msg = errmsg.GetErrMsg(code)
	if code == errmsg.ERROR_WALLET_NOT_EXIST {
		balance = 0
	}
	c.JSON(200, gin.H{
		"status":  code,
		"msg":     msg,
		"balance": balance,
	})
}
