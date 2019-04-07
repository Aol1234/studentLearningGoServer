package groups

import (
	"time"
)

type Group struct { // struct for groups
	GroupId            uint                 `gorm:"primary_key;AUTO_INCREMENT"`
	Name               string               // Name of group
	Desc               string               // Description of group
	Code               string               // Unique joining group
	CreatorId          uint                 // Identifier of Creator
	CreatedAt          time.Time            // time created
	GroupTopicAnalyses []GroupTopicAnalysis `gorm:"foreignkey:GroupId"`
	Members            []Member             `gorm:"foreignkey:GroupId"`
}

type Member struct { // struct for members
	MemberId uint      `gorm:"primary_key;AUTO_INCREMENT"`
	UserId   uint      // Identifier of user
	GroupId  uint      // Identifier of group
	JoinedAt time.Time // time of membership
	Banned   int
}

type GroupTopicAnalysis struct { // struct for group analysis
	GTopId    uint      `gorm:"primary_key;AUTO_INCREMENT"`
	GroupId   uint      // Identifier of group
	TopicId   uint      // Identifier of topic
	TopicName string    // Name of topic
	AvgResult float64   // Average result
	CreatedAt time.Time // Time Created
}

type Data struct { // struct to temporarily hold data
	Groups             []Group
	GroupTopicAnalysis []GroupTopicAnalysis
}
