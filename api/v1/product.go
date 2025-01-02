package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
)

func CreateProduct(c *gin.Context) {
	// todo创建商品
	var product model.Product
	_ = c.ShouldBindJSON(&product)
	code := model.CreateProduct(model.DB, &product)

	msg := errmsg.GetErrMsg(code)
	c.JSON(200, gin.H{
		"status": code,
		"msg":    msg,
	})
}
