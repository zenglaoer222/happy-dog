package model

import (
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
	"log"
)

type Customer struct {
	gorm.Model
	Customer  string `gorm:"type:varchar(20);not null" json:"Customer" validate:"required,min=4,max=12" label:"用户名"`
	Password  string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Address   string `gorm:"type:varchar(20);not null" json:"address" validate:"required,min=6,max=20" label:"地址"`
	PhoneNum  string `gorm:"type:varchar(20);not null" json:"phone_num" validate:"required,min=6,max=20" label:"电话"`
	AvatarUrl string `gorm:"type:varchar(100)" json:"avatar_url"`
}

func ExistCustomer(tx *gorm.DB, id uint) int {
	var customer Customer
	tx.Select("id").Where("id = ?", id).First(&customer)
	if customer.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //用户名已存在
	}
	return errmsg.ERROR_USER_NOT_EXIST
}

func ExistByNameCustomer(tx *gorm.DB, CName string) (code int, id uint) {
	var customer Customer
	tx.Select("id").Where("customer = ?", CName).First(&customer)
	if customer.ID > 0 {
		return errmsg.ERROR_USERNAME_USED, customer.ID //用户名已存在
	}
	return errmsg.ERROR_USER_NOT_EXIST, 0
}

func CreateCustomer(tx *gorm.DB, customer *Customer) (code int) {
	customer.Password = ScriptPw(customer.Password)
	err := tx.Create(&customer).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

func GetCustomer(tx *gorm.DB, customer string, pageSize int, pageNum int) ([]Customer, int64, int) {
	var customers []Customer
	var total int64

	if customer == "" {
		err = tx.Limit(pageSize).Offset(pageNum).Find(&customers).Error
		tx.Model(&Customer{}).Count(&total)
	} else {
		println("pageSize", pageSize, "pageNum", pageNum)
		err = tx.Where("customer LIKE ?", customer+"%").Limit(pageSize).Offset(pageNum).Find(&customers).Error
		if err != nil {
			return nil, 0, errmsg.ERROR
		}
		tx.Model(&Customer{}).Where("customer LIKE ?", customer+"%").Count(&total)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}

	for i := range customers {
		customers[i].Password = ""
	}
	return customers, total, errmsg.SUCCESS

}

func EditCustomer(tx *gorm.DB, id int, customer *Customer) int {
	var maps = make(map[string]interface{})
	maps["customer"] = customer.Customer
	maps["address"] = customer.Address
	maps["phone_num"] = customer.PhoneNum

	err = tx.Model(&Customer{}).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func DeleteCustomer(tx *gorm.DB, id int) int {
	err = tx.Where("id = ?", id).Delete(&Customer{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func uploadAvtar(tx *gorm.DB, id int, avatarUrl string) int {
	err = tx.Model(&Customer{}).Where("id = ?", id).Update("avatar_url", avatarUrl).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func ScriptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 34, 24, 53, 6, 75, 3, 56}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}

	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// 登录验证
func CheckLogin(tx *gorm.DB, username string, password string) (id uint, code int) {
	var customer Customer
	tx.Where("customer = ?", username).First(&customer)
	if customer.ID == 0 {
		return customer.ID, errmsg.ERROR_USER_NOT_EXIST
	}
	if ScriptPw(password) != customer.Password {
		return 0, errmsg.ERROR_PASSWORD_WRONG
	}
	return customer.ID, errmsg.SUCCESS

}
