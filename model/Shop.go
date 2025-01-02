package model

import (
	"fmt"
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Shop struct {
	gorm.Model
	ShopName string `gorm:"type:varchar(20);not null" json:"shop_name" validate:"required,min=4,max=12" label:"商家名称"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Address  string `gorm:"type:varchar(20);not null" json:"address" validate:"required,min=6,max=20" label:"地址"`
}

func ExistShop(tx *gorm.DB, id uint) (code int) {
	var shop Shop
	tx.Select("id").Where("id = ?", id).First(&shop)
	fmt.Println(shop)
	if shop.ID > 0 {
		return errmsg.ERROR_SHOP_USED //商家已存在
	}
	return errmsg.ERROR_SHOP_NOT_EXIST //商家不存在
}

func ExistShopByName(tx *gorm.DB, CName string) (code int, id uint) {
	var shop Shop
	tx.Select("id").Where("shop_name = ?", CName).First(&shop)
	if shop.ID > 0 {
		return errmsg.ERROR_SHOP_USED, shop.ID //商家已存在
	}
	return errmsg.ERROR_SHOP_NOT_EXIST, 0 //商家不存在
}

func CreateShop(tx *gorm.DB, shop *Shop) (code int) {
	shop.Password = ScriptPw(shop.Password)
	err := tx.Create(&shop).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

func GetShop(tx *gorm.DB, shop_name string, pageSize int, pageNum int) ([]Shop, int64, int) {
	var shops []Shop
	var total int64

	if shop_name == "" {
		err = tx.Limit(pageSize).Offset(pageNum).Find(&shops).Error
		tx.Model(&Shop{}).Count(&total)
	} else {
		err = tx.Where("shop_name LIKE ?", shop_name+"%").Limit(pageSize).Offset(pageNum).Find(&shops).Error
		tx.Model(&Shop{}).Where("shop_name LIKE ?", shop_name+"%").Count(&total)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	// 屏蔽密码
	for i, _ := range shops {
		shops[i].Password = ""
	}

	return shops, total, errmsg.SUCCESS

}

func EditShop(tx *gorm.DB, id int, shop *Shop) int {
	var maps = make(map[string]interface{})
	maps["shop_name"] = shop.ShopName
	maps["address"] = shop.Address

	err = tx.Model(&Shop{}).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func DeleteShop(tx *gorm.DB, id int) int {
	err = tx.Where("id = ?", id).Delete(&Shop{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
