package model

import (
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
