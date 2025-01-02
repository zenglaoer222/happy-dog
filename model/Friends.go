package model

import (
	"gorm.io/gorm"
	"happy-dog/utils/errmsg"
)

type Friends struct {
	gorm.Model
	ID1 uint `json:"uid1"`
	ID2 uint `json:"uid2"`

	User1 Customer `gorm:"foreignKey:ID1" json:"user1"`
	User2 Customer `gorm:"foreignKey:ID2" json:"user2"`

	Accept bool `json:"accept"`

	historyMessage []HistoryMessage
}

func ExistFriends(tx *gorm.DB, id1, id2 uint) (mode, code int) {
	var friends Friends
	err := tx.Model(&Friends{}).Where("(id1 = ? and id2 = ?) or (id2 = ? and id1 = ?)", id1, id2, id1, id2).First(&friends).Error
	if err != nil {
		return 0, errmsg.ERROR
	}
	// mode == 1表示已经是好友
	if friends.Accept {
		return 1, errmsg.SUCCESS
	}
	// mode == 2表示已经发送过好友请求
	return 2, errmsg.SUCCESS

}

func BoolFriends(tx *gorm.DB, id uint, id2 uint) (code int) {
	var friends Friends
	err := tx.Model(&Friends{}).Where("id = ? and id2 = ?", id, id2).First(&friends).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func CreateFriends(tx *gorm.DB, userId1, userId2 uint) int {
	err := tx.Create(&Friends{
		ID1:    userId1,
		ID2:    userId2,
		Accept: false,
	}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

func GetFriends(tx *gorm.DB, id uint, mode int) (friends []Friends, code int) {
	var err error
	// mode == 1 得到待添加的好友
	if mode == 1 {
		err = tx.Preload("User1").Where("id2 = ? and accept = false", id).Find(&friends).Error
		if err != nil {
			return nil, errmsg.ERROR
		}
	} else if mode == 2 {
		// mode == 2 得到所有的好友
		err = tx.Preload("User1").Preload("User2").Where("(id1 = ? and accept = true) or ( id2 = ? and accept = true) ", id, id).Find(&friends).Error
		if err != nil {
			return nil, errmsg.ERROR
		}
	}
	for i := range friends {
		friends[i].User1.Password = ""
		friends[i].User2.Password = ""

	}
	return friends, errmsg.SUCCESS

}

func AcceptFriends(tx *gorm.DB, id uint) int {

	err := tx.Model(&Friends{}).Where("id = ?", id).Update("accept", true).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
