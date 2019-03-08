package analysis

import (
	"fmt"
	"github.com/jinzhu/gorm"
	Mcq "studentLearningGoServer/questionnaire"
	userApi "studentLearningGoServer/userModel"
	"time"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

func CollectData(db *gorm.DB) []Mcq.McqResult {
	var results []Mcq.McqResult
	db.Where("user_id = ?", 3).Preload("McqQuestionResult").Find(&results)
	return results

}

func CheckUsersAnalysis(db *gorm.DB, user userApi.User, mcqID uint) error {
	// Check user has used this mcq

	db.AutoMigrate(WeeklyMcqAnalysis{})
	db.AutoMigrate(&WeeklyMcqAnalysisResult{})
	// Check user has weekly Analysis

	var WeeklyAnalysis WeeklyMcqAnalysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcqID).First(&WeeklyAnalysis).Error
	// FIXME: Make condition to checks if WeeklyAnalysis is not null/empty
	if err != nil || WeeklyAnalysis.McqId == 0 {
		fmt.Printf("No Weekly Analysis created for this: %s", err)
		db.Create(&WeeklyMcqAnalysis{UserId: user.UserId, McqId: mcqID, WeeklyMcqAnalysisResults: []WeeklyMcqAnalysisResult{}})
	}
	// Create empty weekly Analysis
	var question []Mcq.McqQuestion
	err = db.Where("mcq_id = ?", mcqID).Find(&question).Error
	if err != nil || len(question) == 0 {
		fmt.Printf("No MCQ with this id: %s", err)
		return err
	}
	numberOfQuestions := len(question)
	db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcqID).First(&WeeklyAnalysis)
	// Check user has weekly Analysis for this mcq question
	var questionAnalysis []WeeklyMcqAnalysisResult
	db.Where("mcq_id = ?", mcqID).Find(&questionAnalysis)
	if len(questionAnalysis) == 0 {
		fmt.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ {
			err = db.Create(&WeeklyMcqAnalysisResult{WeeklyRAna: WeeklyAnalysis.WeeklyRAna, McqId: mcqID, NumberOfResults: 0, AvgTime: time.Duration(0), AvgResult: 0.00}).Error
			if err != nil {
				fmt.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}

	return nil
}

// Analysis SINGLE MCQ EXAM
func GetNewAvg(db *gorm.DB, M []Mcq.McqResult) WeeklyMcqAnalysis {
	db.AutoMigrate(&WeeklyMcqAnalysis{})
	db.AutoMigrate(&WeeklyMcqAnalysisResult{})
	weekAvg := WeeklyMcqAnalysis{}
	// Get Week Analysis
	db.Where("mcq_id = ? AND user_id = ?", M[0].McqId, M[0].UserId).
		Preload("WeeklyMcqAnalysisResults").
		First(&weekAvg)
	// For each Result
	// FIXME: JS cant sent negative values
	for _, result := range M {
		// For each Question Result
		fmt.Print("Range ", result.McqQuestionResult)
		for i, answer := range result.McqQuestionResult {
			// Get weekly Analysis for Question Result
			numberOfResults := weekAvg.WeeklyMcqAnalysisResults[i].NumberOfResults
			fmt.Println("Number of Results", numberOfResults)
			// Include this result in numberOfResults
			numberOfResults += 1
			currentAvgTime := weekAvg.WeeklyMcqAnalysisResults[i].AvgTime
			// Normalise current average result
			currentAvgResult := weekAvg.WeeklyMcqAnalysisResults[i].AvgResult * float64(numberOfResults-1)
			// TODO: ADD new attribute called result value
			thisResult := float64(result.McqQuestionResult[i].Result)
			if currentAvgTime == 0 {
				// If no current average time, current average time duplicated
				currentAvgTime = answer.Time
			}
			newAvgTime := (currentAvgTime + answer.Time) / time.Duration(2)
			newAvgResult := (currentAvgResult + thisResult) / float64(numberOfResults)

			weekAvg.WeeklyMcqAnalysisResults[i].AvgTime = newAvgTime
			weekAvg.WeeklyMcqAnalysisResults[i].AvgResult = newAvgResult
			weekAvg.WeeklyMcqAnalysisResults[i].NumberOfResults += 1
		}
	}

	fmt.Println("Week Avg: ", weekAvg)
	return weekAvg
	//db.Updates(weekAvg)
}
