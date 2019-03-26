package groups

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"time"
)

func CreateGroup(db *gorm.DB, userId uint, name string, desc string) xid.ID {
	db.AutoMigrate(&Group{}, &Member{})
	code := xid.New()
	fmt.Printf("%s\n", code.String())
	timeNow := time.Now()
	fmt.Print(userId)
	var find Group
	db.Where("code = ? AND created_at = ?", code, timeNow).First(&find)
	if find.GroupId != 0 {
		return xid.ID{}
	}
	db.Create(&Group{Name: name, Desc: desc, Code: code.String(), CreatorId: userId, CreatedAt: timeNow, Members: []Member{{UserId: userId, JoinedAt: timeNow}}})
	var member Member
	db.Where("user_id = ? AND joined_at = ?", userId, timeNow).First(&member)
	return code
}

func JoinGroup(db *gorm.DB, code string, userID uint) error {
	var group Group
	joinDate := time.Now()
	db.Where("code = ? AND creator_id <> ?", code, userID).First(&group)
	if group.GroupId == 0 {
		return nil
	}
	db.Create(&Member{UserId: 4, JoinedAt: joinDate, GroupId: group.GroupId})
	return nil
}

func GetGroups(db *gorm.DB, userID uint) (groups []Group, groupsAnalysis []GroupTopicAnalysis, err error) {
	fmt.Println("Groups", userID)
	var membership []Member
	err = db.Where("user_id = ?", userID).Find(&membership).Error
	if err != nil {
		return []Group{}, []GroupTopicAnalysis{}, err
	}
	for _, chosenGroup := range membership {
		var groupAnalysis GroupTopicAnalysis
		db.Where("group_id = ?", chosenGroup.GroupId).First(&groupAnalysis)
		groupsAnalysis = append(groupsAnalysis, groupAnalysis)
		var group Group
		db.Where("group_id = ?", chosenGroup.GroupId).First(&group)
		groups = append(groups, group)
	}
	return groups, groupsAnalysis, nil
}
