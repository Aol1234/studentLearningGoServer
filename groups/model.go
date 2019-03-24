package groups

import (
	"time"
)

type Group struct {
	GroupId            uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name               string
	Desc               string
	Code               string
	CreatorId          uint
	CreatedAt          time.Time
	GroupTopicAnalyses []GroupTopicAnalysis `gorm:"foreignkey:GroupId"`
	Members            []Member             `gorm:"foreignkey:GroupId"`
}

type Member struct {
	MemberId uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserId   uint
	GroupId  uint
	JoinedAt time.Time
	Banned   int
}

type GroupTopicAnalysis struct {
	GTopId    uint `gorm:"primary_key;AUTO_INCREMENT"`
	GroupId   uint
	TopicId   uint
	TopicName string
	AvgResult float64
	CreatedAt time.Time
}

type Data struct {
	Groups             []Group
	GroupTopicAnalysis []GroupTopicAnalysis
}
