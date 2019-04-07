package groups

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"log"
	"time"
)

func CreateGroup(db *gorm.DB, userId uint, name string, desc string) xid.ID {
	// Create group
	code := xid.New() // Create unique code
	timeNow := time.Now()
	var find Group // Check if group has code already
	db.Where("code = ? AND created_at = ?", code, timeNow).First(&find)
	if find.GroupId != 0 {
		log.Println("Group: ", find.GroupId, " already has this code")
		return xid.ID{}
	}
	// Create group
	db.Create(&Group{Name: name, Desc: desc, Code: code.String(), CreatorId: userId, CreatedAt: timeNow, Members: []Member{{UserId: userId, JoinedAt: timeNow}}})
	var member Member // Add creator as member
	db.Where("user_id = ? AND joined_at = ?", userId, timeNow).First(&member)
	return code
}

func JoinGroup(db *gorm.DB, code string, userID uint) error {
	// Join group
	var group Group
	joinDate := time.Now() // Fid group with code and user does not have creator's UserId
	db.Where("code = ? AND creator_id <> ?", code, userID).First(&group)
	if group.GroupId == 0 { // Check if group exists
		log.Println("Could not find group with this code")
		return nil
	}
	// Create new member of group
	db.Create(&Member{UserId: 4, JoinedAt: joinDate, GroupId: group.GroupId})
	return nil
}

func GetGroups(db *gorm.DB, userID uint) (groups []Group, groupsAnalysis []GroupTopicAnalysis, err error) {
	// Return all group details which user member of.
	var membership []Member // Find all groups
	err = db.Where("user_id = ?", userID).Find(&membership).Error
	if err != nil {
		return []Group{}, []GroupTopicAnalysis{}, err
	}
	for _, chosenGroup := range membership { // Iterate through all groups user is a member
		var groupAnalysis []GroupTopicAnalysis // Get analysis for each group
		db.Where("group_id = ?", chosenGroup.GroupId).Find(&groupAnalysis)
		for _, analysis := range groupAnalysis {
			groupsAnalysis = append(groupsAnalysis, analysis)
		}
		var group Group
		db.Where("group_id = ?", chosenGroup.GroupId).First(&group)
		groups = append(groups, group)
	}
	return groups, groupsAnalysis, nil
}
