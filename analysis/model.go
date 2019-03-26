package analysis

import (
	mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	"github.com/jinzhu/gorm"
	"time"
)

// Single one for each MQC for each User
type WeeklyMcqAnalysis struct {
	WeeklyRAna               uint   `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                   uint   `gorm:"foreignkey"` // User Id
	McqId                    uint   `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                    string // Topic Associated with MCQ
	TopicId                  uint
	AvgResult                float64
	LastModified             time.Time                 // Time user last updated result
	WeeklyMcqAnalysisResults []WeeklyMcqAnalysisResult `gorm:"foreignkey:WeeklyRAna"`
}

// Single one for each Question for each MQC
type WeeklyMcqAnalysisResult struct {
	McqAnaId            uint          `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	WeeklyRAna          uint          `json:"weekly_r_ana"`
	QId                 uint          `gorm:"foreignkey"`        // Associated Question Id
	Question            string        `gorm:"foreignkey"`        // Question Name
	McqId               uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults     int           `json:"number_of_results"` // Number of results over past week
	AvgTime             time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult           float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidenceString string
	AvgConfidence       float64 // Avg level of confidence 1 change + 5sec  = v.h >5 = v.l
}

// Single one for each MQC for each User
type MonthlyMcqAnalysis struct {
	MonthlyRAna               uint   `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                    uint   `gorm:"foreignkey"` // User Id
	McqId                     uint   `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                     string // Topic Associated with MCQ
	AvgResult                 float64
	LastModified              time.Time                  // Time user last updated result
	MonthlyMcqAnalysisResults []MonthlyMcqAnalysisResult `gorm:"foreignkey:MonthlyRAna"`
}

// Single one for each Question for each MQC
type MonthlyMcqAnalysisResult struct {
	McqAnaId            uint `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	MonthlyRAna         uint
	QId                 uint          `gorm:"foreignkey"`        // Associated Question Id
	Question            string        `gorm:"foreignkey"`        // Question Name
	McqId               uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults     int           `json:"number_of_results"` // Number of results over past month // Shouldn't go higher than 4
	AvgTime             time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult           float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence       float64
	AvgConfidenceString string
}

// Single one for each MQC for each User
type YearlyMcqAnalysis struct {
	YearlyRAna               uint   `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                   uint   `gorm:"foreignkey"` // User Id
	McqId                    uint   `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                    string // Topic Associated with MCQ
	AvgResult                float64
	LastModified             time.Time                 // Time user last updated result
	YearlyMcqAnalysisResults []YearlyMcqAnalysisResult `gorm:"foreignkey:YearlyRAna"`
}

// Single one for each Question for each MQC
type YearlyMcqAnalysisResult struct {
	McqAnaId            uint `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	YearlyRAna          uint
	QId                 uint          `gorm:"foreignkey"`        // Associated Question Id
	Question            string        `gorm:"foreignkey"`        // Question Name
	McqId               uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults     int           `json:"number_of_results"` // Number of results over past week
	AvgTime             time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult           float64       `json:"avg_result"`        // Avg chance answer is correct
	AvgConfidence       float64
	AvgConfidenceString string
}

// Single one for each MQC for each User
type TotalMcqAnalysis struct {
	TotalRAna               uint   `gorm:"primary_key; AUTO_INCREMENT"`
	UserId                  uint   `gorm:"foreignkey"` // User Id
	McqId                   uint   `gorm:"foreignkey"` // Mcq Id to associate with mcq test
	Topic                   string // Topic Associated with MCQ
	AvgResult               float64
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

// Struct used to return profile data to user
type Data struct {
	McqQuestions []mcq.MCQ
	Weekly       []WeeklyMcqAnalysis
	Monthly      []MonthlyMcqAnalysis
	Yearly       []YearlyMcqAnalysis
	Results      [][]mcq.McqResult
	Topics       []TopicAnalysis
}

// Single topic associated with topic, currently only possible to create new topics through sql
type Topic struct {
	TopicId   uint `gorm:"primary_key;AUTO_INCREMENT"` // Topic Id
	TopicName string
}

type TopicAnalysis struct {
	TopicAnaId uint `gorm:"primary_key;AUTO_INCREMENT"` // Topic Id
	UserId     uint `gorm:"foreignkey"`                 // Associated User Id
	TopicId    uint `gorm:"foreignkey"`                 // Associated Topic Id
	TopicName  string
	AvgResult  float64
}

type Sql struct {
	db *gorm.DB
}
