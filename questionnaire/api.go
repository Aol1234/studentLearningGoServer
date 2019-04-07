package questionnaire

import (
	"github.com/jinzhu/gorm"
	"time"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

func StoreResult(db *gorm.DB, results McqResult) {
	// Store user result
	var questions []McqQuestion // Find questionnaire
	db.Where("mcq_id = ?", results.McqId).
		Find(&questions)
	for index, question := range questions { // Add question id to each question result
		results.McqQuestionResult[index].QId = question.QId
	}
	average := getAverageScore(results) // Get simple average
	results.AverageResult = average
	results.CreatedAt = time.Now() // Get time record stored
	db.Create(&results)            // Create record
}

func CreateMcq(db *gorm.DB, mcq MCQ) {
	// Create questionnaire
	mcq.CreatedAt = time.Now()
	db.Create(&mcq) // Create MCQ
}

func GetMcqs(db *gorm.DB) []MCQ {
	// Get all MCQs
	var Mcqs []MCQ
	db.Find(&Mcqs) // Collect MCQs
	return Mcqs
}

func RetrieveMcq(db *gorm.DB, mcqId uint) MCQ {
	// Retrieve Mcq
	var Mcq MCQ // Collect MCQ with this ID
	db.Where("mcq_id = ?", mcqId).Preload("McqQuestions").Preload("McqQuestions.Answers").First(&Mcq)
	db.Table("mcqs").Update("last_used", time.Now()).Where("mcq_id = ?", mcqId) // Update MCQ last_used to now
	return Mcq
}

func getAverageScore(results McqResult) float64 {
	// Get average of result
	var cumulative float64
	for _, result := range results.McqQuestionResult {
		cumulative += float64(result.Result) / float64(result.Total)
	}
	average := cumulative / float64(len(results.McqQuestionResult))
	return average
}
