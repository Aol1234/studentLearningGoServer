package user

import "github.com/jinzhu/gorm"

type User struct { // User struct
	UserId uint   `gorm:"primary_key;AUTO_INCREMENT"`
	UID    string `json:"uid"`
	Role   string // ADMIN
}

type UserPreference struct { // User Preferences struct
	UserId         uint  `gorm:"primary_key"`
	GraphInfoModal int64 // 0 = False 1 = True
}
type Sql struct { // Sql Database
	db *gorm.DB
}
