package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"happy-dog/utils/validator"
	"strconv"
)

// 查询用户是否存在
// todo 查询是否存在

// 新增商家
// todo 新增
func AddShop(c *gin.Context) {
	var data model.Shop
	var msg string
	var code int
	_ = c.ShouldBindJSON(&data)

	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(200, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}

	code, _ = model.ExistShopByName(model.DB, data.ShopName)
	if code != errmsg.ERROR_SHOP_NOT_EXIST {
		msg = errmsg.GetErrMsg(code)
	} else {
		code = model.CreateShop(model.DB, &data)
		msg = errmsg.GetErrMsg(code)
	}
	c.JSON(200, gin.H{
		"status":  code,
		"message": msg,
	})

}

// 商家列表
func GetShop(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	shop_name := c.Query("shop_name")

	pageNum = (pageNum - 1) * pageSize
	// 不分页
	if pageSize == 0 {
		pageSize = -1
	}
	shopList, total, code := model.GetShop(model.DB, shop_name, pageSize, pageNum)
	c.JSON(200, gin.H{
		"status":  code,
		"data":    shopList,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
	return

}

// 编辑商家
func EditShop(c *gin.Context) {
	var data model.Shop
	var id_now uint
	_ = c.ShouldBindJSON(&data)
	id, _ := strconv.Atoi(c.Param("id"))
	// 查询修改是否符合要求
	msg, code := validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(200, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}
	// 查询是否存在
	code, id_now = model.ExistShopByName(model.DB, data.ShopName)
	if code == errmsg.ERROR_SHOP_USED {
		// 该用户修改自身
		if id_now == uint(id) {
			code = model.EditShop(model.DB, id, &data)
		} //该用户修改为他人名称
		msg = errmsg.GetErrMsg(code)
	}
	// 修改数据
	if code == errmsg.ERROR_SHOP_NOT_EXIST {
		code = model.EditShop(model.DB, id, &data)
		msg = errmsg.GetErrMsg(code)
	}

	c.JSON(200, gin.H{
		"status":  code,
		"message": msg,
	})
	return

}

// 删除商家
// todo 删除商家
func DeleteShop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteShop(model.DB, id)
	c.JSON(200, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
