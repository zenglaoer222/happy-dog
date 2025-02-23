package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"happy-dog/utils/validator"
	"strconv"
)

// 新增顾客
func AddCustomer(c *gin.Context) {
	var data model.Customer
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

	code, _ = model.ExistByNameCustomer(model.DB, data.Customer)
	if code != errmsg.ERROR_USER_NOT_EXIST {
		msg = errmsg.GetErrMsg(code)
	} else {
		code = model.CreateCustomer(model.DB, &data)
		msg = errmsg.GetErrMsg(code)
	}
	c.JSON(200, gin.H{
		"status":  code,
		"message": msg,
	})

}

// 顾客列表
func GetCustomer(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	customer := c.Query("customer")

	pageNum = (pageNum - 1) * pageSize
	// 不分页
	if pageSize == 0 {
		pageSize = -1
	}
	customerList, total, code := model.GetCustomer(model.DB, customer, pageSize, pageNum)
	c.JSON(200, gin.H{
		"status":  code,
		"data":    customerList,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
	return

}

// 获取用户个人信息
func GetCustomerInfo(c *gin.Context) {
	id, _ := c.Get("user_id")

	uid, ok := id.(uint)
	if !ok {
		c.JSON(200, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}
	customer, code := model.GetCustomerInfo(model.DB, uid)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  code,
		"data":    customer,
		"message": errmsg.GetErrMsg(code),
	})

}

// 编辑用户
func EditCustomer(c *gin.Context) {
	var data model.Customer
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
	code, id_now = model.ExistByNameCustomer(model.DB, data.Customer)
	if code == errmsg.ERROR_USERNAME_USED {
		// 该用户修改自身
		if id_now == uint(id) {
			code = model.EditCustomer(model.DB, id, &data)
		} //该用户修改为他人名称
		msg = errmsg.GetErrMsg(code)
	}
	// 修改数据
	if code == errmsg.ERROR_USER_NOT_EXIST {
		code = model.EditCustomer(model.DB, id, &data)
		msg = errmsg.GetErrMsg(code)
	}

	c.JSON(200, gin.H{
		"status":  code,
		"message": msg,
	})
	return

}

// 删除用户
// todo 删除用户
func DeleteCustomer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteCustomer(model.DB, id)
	c.JSON(200, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
