package model

import (
	"fmt"
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Wallet struct {
	ID       uint    `gorm:"primarykey"`
	Cid      uint    `gorm:"type:int;" json:"cid"`
	Balance  float64 `gorm:"type:decimal(10,2);" json:"balance"`
	Password string  `gorm:"type:varchar(6);not null" json:"password" validate:"required,min=6,max=6" label:"支付密码"`
}

func CreateWallet(tx *gorm.DB, id uint, password string) int {
	var wallet Wallet
	wallet.Cid = id
	wallet.Balance = 1000
	wallet.Password = password
	err := tx.Model(&Wallet{}).Create(&wallet).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func InquireBalance(tx *gorm.DB, cid uint) (int, float64) {
	var wallet Wallet
	fmt.Println(cid)
	err := tx.Select("cid,balance").Where("cid = ?", cid).First(&wallet).Error
	if wallet.Cid > 0 && err == nil {
		return errmsg.SUCCESS, wallet.Balance
	}
	return errmsg.ERROR_WALLET_NOT_EXIST, float64(0)
}

func DeductBalance(tx *gorm.DB, cid uint, price float64) int {
	var maps = make(map[string]interface{})
	code, balance := InquireBalance(tx, cid)
	if code == errmsg.ERROR_WALLET_NOT_EXIST {
		return code
	}
	if balance < price {
		return errmsg.ERROR_WALLET_BALANCE_NOT_ENOUGH
	}
	maps["balance"] = balance - price
	err := tx.Model(&Wallet{}).Where("cid = ?", cid).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

func CheckPassword(tx *gorm.DB, cid uint, password string) (code int, ok bool) {
	var wallet Wallet
	err := tx.Where("cid = ?", cid).First(&wallet).Error
	if err != nil {
		return errmsg.ERROR_WALLET_NOT_EXIST, false
	}
	if wallet.Password != password {
		return errmsg.ERROR_WALLET_PASSWORD_WRONG, false
	}
	return errmsg.SUCCESS, true
}
