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
	var questions []McqQuestion
	db.Where("mcq_id = ?", results.McqId).
		Find(&questions)
	for index, question := range questions {
		results.McqQuestionResult[index].QId = question.QId
	}
	average := getAverageScore(results)
	results.AverageResult = average
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
	db.AutoMigrate(&MCQ{})
	var Mcq MCQ
	db.Where("mcq_id = ?", mcqId).Preload("McqQuestions").Preload("McqQuestions.Answers").First(&Mcq)
	db.Table("mcqs").Update("last_used", time.Now()) // TODO: Update when mcq last used
	return Mcq
}

func getAverageScore(results McqResult) float64 {
	var cumulative float64
	for _, result := range results.McqQuestionResult {
		cumulative += float64(result.Result) / float64(result.Total)
	}
	average := cumulative / float64(len(results.McqQuestionResult))
	return average
}
