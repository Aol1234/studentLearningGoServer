package analysis

import (
	Mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	"github.com/jinzhu/gorm"
	"time"
)

// Single one for each MQC for each User
type WeeklyMcqAnalysis struct {
	WeeklyRAna               uint                      `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                   uint                      `gorm:"foreignkey"` // User Id
	McqId                    uint                      `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                    string                    // Topic Associated with MCQ
	LastModified             time.Time                 // Time user last updated result
	WeeklyMcqAnalysisResults []WeeklyMcqAnalysisResult `gorm:"foreignkey:WeeklyRAna"`
}

// Single one for each Question for each MQC
type WeeklyMcqAnalysisResult struct {
	McqAnaId        uint          `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	WeeklyRAna      uint          `json:"weekly_r_ana"`
	QId             uint          `gorm:"foreignkey"`        // Associated Question Id
	Question        string        `gorm:"foreignkey"`        // Question Name
	McqId           uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults int           `json:"number_of_results"` // Number of results over past week
	AvgTime         time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult       float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence   float64       // Avg level of confidence 1 change + 5sec  = v.h >5 = v.l
}

// Single one for each MQC for each User
type MonthlyMcqAnalysis struct {
	MonthlyRAna               uint                       `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                    uint                       `gorm:"foreignkey"` // User Id
	McqId                     uint                       `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                     string                     // Topic Associated with MCQ
	LastModified              time.Time                  // Time user last updated result
	MonthlyMcqAnalysisResults []MonthlyMcqAnalysisResult `gorm:"foreignkey:WeeklyRAna"`
}

// Single one for each Question for each MQC
type MonthlyMcqAnalysisResult struct {
	McqAnaId        uint `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	MonthlyRAna     uint
	QId             uint          `gorm:"foreignkey"`        // Associated Question Id
	Question        string        `gorm:"foreignkey"`        // Question Name
	McqId           uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults int           `json:"number_of_results"` // Number of results over past week
	AvgTime         time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult       float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence   float64       // Avg level of confidence 1 change + 5sec  = v.h >5 = v.l
}

// Single one for each MQC for each User
type YearlyMcqAnalysis struct {
	YearlyRAna               uint                      `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                   uint                      `gorm:"foreignkey"` // User Id
	McqId                    uint                      `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                    string                    // Topic Associated with MCQ
	LastModified             time.Time                 // Time user last updated result
	YearlyMcqAnalysisResults []YearlyMcqAnalysisResult `gorm:"foreignkey:YearlyRAna"`
}

// Single one for each Question for each MQC
type YearlyMcqAnalysisResult struct {
	McqAnaId        uint `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	YearlyRAna      uint
	QId             uint          `gorm:"foreignkey"`        // Associated Question Id
	Question        string        `gorm:"foreignkey"`        // Question Name
	McqId           uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults int           `json:"number_of_results"` // Number of results over past week
	AvgTime         time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult       float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence   float64       // Avg level of confidence 1 change + 5sec  = v.h >5 = v.l
}

// Single one for each MQC for each User
type TotalMcqAnalysis struct {
	TotalRAna               uint                     `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                  uint                     `gorm:"foreignkey"` // User Id
	McqId                   uint                     `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                   string                   // Topic Associated with MCQ
	LastModified            time.Time                // Time user last updated result
	TotalMcqAnalysisResults []TotalMcqAnalysisResult `gorm:"foreignkey:TotalRAna"`
}

// Single one for each Question for each MQC
type TotalMcqAnalysisResult struct {
	McqAnaId        uint `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	TotalRAna       uint
	QId             uint          `gorm:"foreignkey"`        // Associated Question Id
	Question        string        `gorm:"foreignkey"`        // Question Name
	McqId           uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults int           `json:"number_of_results"` // Number of results over past week
	AvgTime         time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult       float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence   float64       // Avg level of confidence 1 change + 5sec  = v.h >5 = v.l
}

type Data struct {
	Data []Mcq.McqResult
}

type Sql struct {
	db *gorm.DB
}
