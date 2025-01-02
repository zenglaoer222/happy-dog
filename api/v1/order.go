package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"strconv"
)

type OrderDetails struct {
	Total     int               `json:"total"`
	Order     model.Order       `json:"order"`
	OrderItem []model.OrderItem `json:"order_item"`
}

// 新增订单
// 2024/12/2 更新为事务处理：扣除订单，生成订单详情，生成订单
func CreateOrder(c *gin.Context) {
	var code int
	var data OrderDetails
	_ = c.ShouldBindJSON(&data)
	fmt.Println(data)
	// 判断本人
	cid, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "获取user_id失败",
		})
	}

	if cid != data.Order.Cid {
		c.JSON(200, gin.H{
			"code":    errmsg.ERROR,
			"message": "id不匹配",
		})
		return
	}

	tx := model.DB.Begin()
	// 判断商家是否存在
	code = model.ExistShop(tx, data.Order.Sid)
	fmt.Println(code)
	if code == errmsg.ERROR_SHOP_NOT_EXIST {
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	// 判断用户是否存在
	code = model.ExistCustomer(tx, data.Order.Cid)
	if code == errmsg.ERROR_USER_NOT_EXIST {
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	// 判断余额
	//var balance float64
	//fmt.Println(balance)
	code, _ = model.InquireBalance(tx, data.Order.Cid)
	if code == errmsg.ERROR_WALLET_NOT_EXIST {
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 扣除余额
	code = model.DeductBalance(tx, data.Order.Cid, data.Order.TotalPrice)
	if code != errmsg.SUCCESS {
		tx.Rollback()
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	//  创建订单
	code = model.CreateOrder(tx, &data.Order)
	if code == errmsg.ERROR_ORDER_CREATE_FAIL {
		tx.Rollback()
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	//  创建订单详情
	fmt.Println(data.OrderItem)
	code = model.CreateItems(tx, data.OrderItem, data.Order.ID)
	if code == errmsg.ERROR {
		tx.Rollback()
		c.JSON(200, gin.H{
			"code":    code,
			"message": "创建订单详情失败",
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code":    code,
		"message": errmsg.GetErrMsg(code),
	})

}

// todo: 查询订单
func GetOrder(c *gin.Context) {
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

	var orderDetails []OrderDetails

	code, orders, total := model.InquireOrders(model.DB, uint(cid))
	if code != errmsg.SUCCESS {
		c.JSON(200, gin.H{
			"code":    code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	var ordersItem []model.OrderItem
	var itemCount int64
	for _, order := range orders {
		ordersItem, itemCount, code = model.GetItems(model.DB, order.ID)
		if code != errmsg.SUCCESS {
			c.JSON(200, gin.H{
				"code":    code,
				"message": "获取订单商品失败",
			})
			return
		}
		orderDetails = append(orderDetails, OrderDetails{
			Total:     int(itemCount),
			Order:     order,
			OrderItem: ordersItem,
		})

	}

	c.JSON(200, gin.H{
		"status":  code,
		"data":    orderDetails,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
	return

}
