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
	WeeklyMcqAnalysisResults []WeeklyMcqAnalysisResult `gorm:"foreignkey:WeeklyRAna"`
}

// Single one for each Question for each MQC
type WeeklyMcqAnalysisResult struct {
	McqAnaId        uint          `gorm:"primary_key;AUTO_INCREMENT"` // Analysis Id
	WeeklyRAna      uint          `json:"weekly_r_ana"`
	McqId           uint          `json:"mcq_id"`            // Mcq Id to associate with mcq test
	NumberOfResults int           `json:"number_of_results"` // Number of results over past week
	AvgTime         time.Duration `json:"avg_time"`          // Avg time taken to answer this question
	AvgResult       float64       `json:"avg_result"`        // Avg chance answer is correct
}

type Data struct {
	Data []Mcq.McqResult
}

type Sql struct {
	db *gorm.DB
}
