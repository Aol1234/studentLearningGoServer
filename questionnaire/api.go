package questionnaire

import (
	"github.com/jinzhu/gorm"
	"time"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

func StoreResult(db *gorm.DB, results McqResult) {
	db.AutoMigrate(&McqResult{})
	db.AutoMigrate(&McqQuestionResult{})
	results.CreatedAt = time.Now()
	db.Create(&results)
}

func CreateMcq(db *gorm.DB, request MCQ) {
	db.AutoMigrate(&MCQ{}, &McqQuestion{}, &McqAnswer{})
	request.CreatedAt = time.Now()
	db.Create(&request)
}

func GrabMcqs(db *gorm.DB) []MCQ {
	db.AutoMigrate(&MCQ{}, &McqQuestion{}, &McqAnswer{})
	var Mcqs []MCQ
	db.Find(&Mcqs)
	return Mcqs
}

func RetrieveMcq(db *gorm.DB, mcqId uint) MCQ {
	var Mcq MCQ
	db.Where("mcq_id = ?", mcqId).Preload("McqQuestions").Preload("McqQuestions.Answers").First(&Mcq)
	// TODO: Update when mcq last used
	return Mcq
}
