package questionnaire

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Sql struct {
	db *gorm.DB
}

type MCQ struct {
	McqId        uint          `gorm:"primary_key;AUTO_INCREMENT"`
	UserId       uint          `gorm:"foreignkey"` // Created by UserId
	Topic        string        // Topic associated with mcq
	Name         string        // Name of Mcq
	Desc         string        // Description for user
	CreatedAt    time.Time     // Created at
	LastUsed     time.Time     // Last user to use this mcq
	McqQuestions []McqQuestion `gorm:"foreignkey:McqId"`
}

type McqQuestion struct {
	QId      uint `gorm:"primary_key;AUTO_INCREMENT"`
	McqId    uint `gorm:"foreignkey:McqId"` // Mcq Id to associate with mcq test
	Question string
	Answers  []McqAnswer `gorm:"foreignkey:QId"`
}

type McqAnswer struct {
	AId    uint   `gorm:"primary_key;AUTO_INCREMENT"` // Answer Id
	QId    uint   // Question Id to associate with Question
	Text   string `json:"text, omitempty"`  // Answer text
	Value  int    `json:"value, omitempty"` // Id of Answer // Possibility of multiple or close answers
	Result int    // result mark of Answer
}

type McqResult struct {
	McqRId            uint                `gorm:"primary_key;AUTO_INCREMENT"` // Result Id
	McqId             uint                `gorm:"foreignkey:McqId"`           // Mcq Id to associate with mcq test
	UserId            uint                `gorm:"foreignkey:UserId"`          // User Id to associate with user profile
	CreatedAt         time.Time           `json:"created_at, omitempty"`      // Time MCQ Test finished
	McqQuestionResult []McqQuestionResult `gorm:"foreignkey:McqRId"`
}

type McqQuestionResult struct {
	RId     uint          `gorm:"primary_key;AUTO_INCREMENT"` // Result Id
	QId     uint          `gorm:"foreignkey:QId"`
	McqRId  uint          // Mcq Result Id
	Result  int           `json:"result, omitempty"`  // Answer value
	Total   int           `json:"total, omitempty"`   // Max value
	Time    time.Duration `json:"time, omitempty"`    // Amount of time taken to answer question in seconds
	Changes int           `json:"changes, omitempty"` // Amount of time answer changed
}
