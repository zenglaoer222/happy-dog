package model

import (
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Product struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(20);not null" json:"name"`
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	ShopId   int     `gorm:"type:int;not null" json:"sid"`
	StockNum int     `gorm:"type:int;not null" json:"stock_num"`
}

func ExistProduct(tx *gorm.DB, id uint) (int, uint) {
	var product Product
	tx.Select("id").Where("id = ?", id).First(&product)
	if product.ID > 0 {
		return errmsg.SUCCESS, product.ID
	}
	return errmsg.ERROR_PRODUCT_NOT_EXIST, 0
}

func CreateProduct(tx *gorm.DB, product *Product) int {
	//创建商品
	err := tx.Model(&Product{}).Create(&product).Error
	if err != nil {
		return errmsg.ERROR_PRODUCT_CREATE_FAIL
	}
	return errmsg.SUCCESS
}
