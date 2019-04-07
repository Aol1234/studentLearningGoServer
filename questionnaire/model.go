package questionnaire

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Sql struct {
	db *gorm.DB
}

type MCQ struct { // struct to store details of Mcq
	McqId        uint          `gorm:"primary_key;AUTO_INCREMENT"`
	UserId       uint          `gorm:"foreignkey"` // Created by UserId
	Topic        string        // Topic associated with mcq
	Name         string        // Name of Mcq
	Desc         string        // Description for user
	CreatedAt    time.Time     // Created at
	LastUsed     time.Time     // Last user to use this mcq
	McqQuestions []McqQuestion `gorm:"foreignkey:McqId"`
}

type McqQuestion struct { // struct for MCQ questions
	QId      uint        `gorm:"primary_key;AUTO_INCREMENT"`
	McqId    uint        `gorm:"foreignkey:McqId"` // Mcq Id to associate with mcq test
	Question string      // Question text
	Answers  []McqAnswer `gorm:"foreignkey:QId"`
}

type McqAnswer struct { // struct for MCQ answer
	AId    uint   `gorm:"primary_key;AUTO_INCREMENT"` // Answer Id
	QId    uint   // Question Id to associate with Question
	Text   string `json:"text, omitempty"`  // Answer text
	Value  int    `json:"value, omitempty"` // Id of Answer // Possibility of multiple or close answers
	Result int    // result mark of Answer
}

type McqResult struct { // struct for MCQ results
	McqRId            uint                `gorm:"primary_key;AUTO_INCREMENT"` // Result Id
	McqId             uint                `gorm:"foreignkey:McqId"`           // Mcq Id to associate with mcq test
	UserId            uint                `gorm:"foreignkey:UserId"`          // User Id to associate with user profile
	AverageResult     float64             // Average result
	CreatedAt         time.Time           // Time MCQ Test finished
	McqQuestionResult []McqQuestionResult `gorm:"foreignkey:McqRId"`
}

type McqQuestionResult struct { // struct for MCQ Question result
	RId     uint          `gorm:"primary_key;AUTO_INCREMENT"` // Result Id
	QId     uint          `gorm:"foreignkey:QId"`             // Associate with question derived from
	McqRId  uint          // Mcq Result Id
	Result  int           `json:"result, omitempty"`  // Answer value
	Total   int           `json:"total, omitempty"`   // Max value
	Time    time.Duration `json:"time, omitempty"`    // Amount of time taken to answer question in seconds
	Changes int           `json:"changes, omitempty"` // Amount of time answer changed
}
