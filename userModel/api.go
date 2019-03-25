package userModel

import (
	"fmt"
	"github.com/Aol1234/studentLearningGoServer/sessions"
	"github.com/jinzhu/gorm"
	"log"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

func CreateUser(db *gorm.DB, UID string) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserPreference{})
	db.AutoMigrate(&UserScoreTest{})
	fmt.Println("here")
	db.Create(User{UID: UID})
	var user User
	fmt.Println("test", UID)

	db.Where("UID = ?", UID).First(&user)
	fmt.Println(user)
}

func LoginVerification(db *gorm.DB, UID string) User {
	if UID == "" {
		log.Panicf("Missing UID")
		return User{}
	}
	db.AutoMigrate(&User{})
	var user User
	db.Where("UID = ?", UID).First(&user)
	return user
}

func UpdateUserPreferences(db *gorm.DB, UserPreferences UserPreference) error {
	db.AutoMigrate(&UserPreference{})
	var userP UserPreference
	userP.UserId = UserPreferences.UserId
	db.Model(&userP).Update("graph_info_modal", UserPreferences.GraphInfoModal)
	return nil
}

func RetrieveUserPreferences(db *gorm.DB, UserPreferences UserPreference) (UserPreference, error) {
	var user UserPreference
	db.Where("user_id = ?", UserPreferences.UserId).First(&user)
	user.UserId = 0 // Remove identifying info
	return user, nil
}

func VerifyUserId(bearerToken string) (uint, bool) {
	UserId, found := sessions.Get(bearerToken)
	if found == false {
		panic(bearerToken + ": Invalid token")
		return 0, false
	}
	return UserId.(uint), true
}

func SetCookie(user User, bearerToken string) bool {
	if bearerToken == "" {
		return false
	}
	sessions.Set("Bearer  "+bearerToken, user.UserId, 0)
	return true
}
