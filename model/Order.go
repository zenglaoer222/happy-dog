package model

import (
	"fmt"
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Order struct {
	gorm.Model
	Sid        uint     `json:"sid"` //商店id
	Cid        uint     `json:"cid"` //顾客id
	Shop       Shop     `gorm:"foreignKey:Sid"`
	Customer   Customer `gorm:"foreignKey:Cid"`
	TotalPrice float64  `gorm:"type:decimal(10,2);" json:"total_price"` //总价
	State      int8     `json:"state"`                                  //订单状态
}

func CreateOrder(tx *gorm.DB, order *Order) int {

	err := tx.Create(&order).Error
	if err != nil {
		return errmsg.ERROR_ORDER_CREATE_FAIL
	}
	return errmsg.SUCCESS
}

func InquireOrders(tx *gorm.DB, cid uint) (code int, orders []Order, total int64) {
	err := tx.Joins("Customer").Joins("Shop").Where("cid = ?", cid).Find(&orders).Count(&total).Error
	if err != nil {
		return errmsg.ERROR_ORDER_INQUIRE_FAIL, nil, 0
	}
	//清除密码
	for i := range orders {
		orders[i].Customer.Password = ""
		orders[i].Shop.Password = ""
	}
	return errmsg.SUCCESS, orders, total

}
func InquireOrdersForShop(tx *gorm.DB, sid uint) (code int, orders []Order, total int64) {
	err := tx.Joins("Customer").Joins("Shop").Where("sid = ?", sid).Find(&orders).Count(&total).Error
	if err != nil {
		return errmsg.ERROR_ORDER_INQUIRE_FAIL, nil, 0
	}
	//清除密码
	for i := range orders {
		orders[i].Customer.Password = ""
		orders[i].Shop.Password = ""
	}
	return errmsg.SUCCESS, orders, total

}

func CheckOrderShop(tx *gorm.DB, oid uint, sid uint) (code int) {
	var order Order
	err := tx.Where("id = ?", oid).First(&order).Error
	if err != nil {
		return errmsg.ERROR
	}
	fmt.Println(order.Sid, sid)
	if order.Sid != sid {
		return errmsg.ERROR_ORDER_SHOP_NOT_SAME
	}
	return errmsg.SUCCESS
}

func FinishOrder(tx *gorm.DB, oid uint) (code int) {
	err := tx.Model(&Order{}).Where("id = ?", oid).Update("state", 1).Error //更新订单状态为已完成
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
