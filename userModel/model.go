package userModel

import "github.com/jinzhu/gorm"

type User struct {
	UserId uint   `gorm:"primary_key;AUTO_INCREMENT"`
	UID    string `json:"uid"`
}

type UserScoreTest struct {
	UserScoreId uint   `gorm:"primary_key"`
	UserId      uint   `gorm:"foreignkey"`
	Value1      string `json:"value_1"`
	Value2      string `json:"value_2"`
}

type UserPreference struct {
	UserId         uint `gorm:"primary_key"`
	GraphInfoModal int64
}
type Sql struct {
	db *gorm.DB
}
