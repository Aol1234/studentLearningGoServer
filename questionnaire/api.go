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
