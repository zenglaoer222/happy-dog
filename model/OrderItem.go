package model

import (
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type OrderItem struct {
	gorm.Model

	GoodsId uint    `json:"goods_id"` // 商品id
	OrderId uint    `json:"order_id"` // 订单id
	Order   Order   `gorm:"foreignKey:OrderId"`
	Product Product `gorm:"foreignKey:GoodsId"`
	Count   int     `json:"count"` // 商品数量
}

func AddItem(tx *gorm.DB, item OrderItem) int {

	err := tx.Create(&item).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func CreateItems(tx *gorm.DB, items []OrderItem, Oid uint) int {
	for _, item := range items {
		// 将每一项都放进表中
		item.OrderId = Oid
		code := AddItem(tx, item)
		if code != errmsg.SUCCESS {
			return errmsg.ERROR
		}
	}
	return errmsg.SUCCESS
}

func GetItems(tx *gorm.DB, Oid uint) (orderItem []OrderItem, total int64, code int) {
	var items []OrderItem
	err := tx.Joins("Product").Joins("Order").Where("order_id = ?", Oid).Find(&items).Count(&total).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return items, total, errmsg.SUCCESS
}
