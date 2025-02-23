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
	Picture  string  `gorm:"type:varchar(100)" json:"picture_url"`
}

func GetLastId(tx *gorm.DB) (int, uint) {
	var product Product
	err := tx.Model(&Product{}).Last(&product).Error
	if err != nil {
		return errmsg.ERROR, 0
	}
	return errmsg.SUCCESS, product.ID
}

func ExistProduct(tx *gorm.DB, id uint) (int, uint, int) {
	var product Product
	tx.Where("id = ?", id).First(&product)
	if product.ID > 0 {
		return errmsg.SUCCESS, product.ID, product.ShopId
	}
	return errmsg.ERROR_PRODUCT_NOT_EXIST, 0, 0
}

func CreateProduct(tx *gorm.DB, product *Product) int {
	//创建商品
	err := tx.Model(&Product{}).Create(&product).Error
	if err != nil {
		return errmsg.ERROR_PRODUCT_CREATE_FAIL
	}
	return errmsg.SUCCESS
}

func GetProductList(tx *gorm.DB, id uint, pageNum int) (code int, product []Product, total int64) {
	err := tx.Model(&Product{}).Where("shop_id = ?", id).Limit(8).Offset(8 * (pageNum - 1)).Find(&product).Error
	if err != nil {
		return errmsg.ERROR, nil, 0
	}
	err = tx.Model(&Product{}).Where("shop_id = ?", id).Count(&total).Error

	return errmsg.SUCCESS, product, total
}

func DeleteProduct(tx *gorm.DB, pid int) int {
	err := tx.Model(&Product{}).Where("id = ?", pid).Delete(&Product{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func uploadProductUrl(tx *gorm.DB, pid int, url string) int {
	err = tx.Model(&Product{}).Where("id = ?", pid).Update("picture", url).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
