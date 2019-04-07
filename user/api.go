package user

import (
	"github.com/Aol1234/studentLearningGoServer/sessions"
	"github.com/jinzhu/gorm"
	"log"
)

func NewSql(db *gorm.DB) *Sql { // Return database
	return &Sql{db: db}
}

func CreateUser(db *gorm.DB, UID string) { // Create User
	db.Create(&User{UID: UID})
	var user User // Return newly created user
	db.Where("uid = ?", UID).First(&user)
	db.Create(&UserPreference{UserId: user.UserId, GraphInfoModal: 0}) // Create record of user preference for user
}

func LoginVerification(db *gorm.DB, UID string) User { // Verify user Firebase UID
	if UID == "" { // Missing UID
		log.Println("Missing UID")
		return User{}
	}
	var user User // Find User
	db.Where("UID = ?", UID).First(&user)
	return user
}

func UpdateUserPreferences(db *gorm.DB, UserPreferences UserPreference) error {
	// Update user preferences
	var userP UserPreference
	userP.UserId = UserPreferences.UserId
	err := db.Model(&userP).Update("graph_info_modal", UserPreferences.GraphInfoModal).Error // Update preference of user
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func RetrieveUserPreferences(db *gorm.DB, UserPreferences UserPreference) (UserPreference, error) {
	// Retrieve user preferences
	var user UserPreference
	err := db.Where("user_id = ?", UserPreferences.UserId).First(&user).Error
	user.UserId = 0 // Remove identifying info
	if err != nil {
		log.Println("Failed to retrieve preferences")
		return user, err
	}
	return user, nil
}

func VerifyUserId(bearerToken string) (uint, bool) {
	// Verify user has session
	UserId, found := sessions.Get(bearerToken)
	if found == false {
		log.Println(bearerToken + ": Invalid token")
		return 0, false
	}
	return UserId.(uint), true
}
