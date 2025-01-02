package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/middleware"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"net/http"
)

func Login(c *gin.Context) {
	var data model.Customer
	_ = c.ShouldBindJSON(&data)
	var token string

	cid, code := model.CheckLogin(model.DB, data.Customer, data.Password)

	if code == errmsg.SUCCESS {
		token, _ = middleware.SetToken(data.Customer, "customer", cid)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      cid,
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}

func ManagerLogin(c *gin.Context) {
	var data model.Manager
	_ = c.ShouldBindJSON(&data)
	var token string
	code := model.ManagerLogin(model.DB, data.Name, data.Password)

	if code == errmsg.SUCCESS {
		token, _ = middleware.SetToken(data.Name, "manager", 0)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}
