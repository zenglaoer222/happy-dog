package model

import (
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Manager struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Password string `gorm:"type:varchar(20);not null" json:"password"`
}

func ManagerLogin(tx *gorm.DB, name string, password string) int {
	var manager Manager
	tx.Where("name = ?", name).First(&manager)
	if manager.Password == password {
		return errmsg.SUCCESS
	}
	return errmsg.ERROR_MANAGER_NOT_EXIST

}
