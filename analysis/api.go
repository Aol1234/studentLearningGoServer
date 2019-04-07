package analysis

import (
	g "github.com/Aol1234/studentLearningGoServer/groups"
	Mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	userApi "github.com/Aol1234/studentLearningGoServer/user"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"time"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

func CollectData(db *gorm.DB, userId uint, timePeriod string) {
	// Collect mcq results for a particular user
	var user userApi.User
	db.Where("user_id = ?", userId).First(&user)

	var period time.Time // Select analysis
	if timePeriod == "Month" {
		period = time.Now().Add(-30 * (24 * time.Hour))
	} else if timePeriod == "Week" {
		period = time.Now().Add(-7 * (24 * time.Hour))
	} else if timePeriod == "Year" {
		period = time.Now().Add(-365 * (24 * time.Hour))
	}

	var collection []Mcq.McqResult // Collect all mcq ids
	collection = collectDistinctMcq(db, period)
	if timePeriod == "Week" { // Weekly analysis
		for _, mcq := range collection {
			var results []Mcq.McqResult // Collect all mcq results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ? AND created_at >= ?", user.UserId, mcq.McqId, period).Preload("McqQuestionResult").Find(&results)
			if len(results) > 0 { // Check if results not empty
				err := checkTodayAnalysis(db, user, mcq.McqId)
				if err != nil {
					log.Println(err)
				}
				getWeeklyAnalysis(db, results)
			}
		}
		getTopicAnalysis(db, userId) // Collect topic analysis
	} else if timePeriod == "Month" { // Monthly analysis
		for _, mcq := range collection {
			var analysis WeeklyMcqAnalysis // Weekly Analysis results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).Preload("WeeklyMcqAnalysisResults").First(&analysis)
			if len(analysis.WeeklyMcqAnalysisResults) > 0 { // Check if results not empty
				err := checkMonthlyAnalysis(db, user, mcq.McqId)
				if err != nil {
					log.Println(err)
				}
				getMonthlyAnalysis(db, analysis)
			}
		}
	} else if timePeriod == "Year" { // Yearly analysis
		for _, mcq := range collection {
			var analysis MonthlyMcqAnalysis // Weekly Analysis results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).Preload("MonthlyMcqAnalysisResults").First(&analysis)
			if len(analysis.MonthlyMcqAnalysisResults) > 0 { // Check if results not empty
				err := checkYearlyAnalysis(db, user, mcq.McqId)
				if err != nil {
					log.Println(err)
				}
				getYearlyAnalysis(db, analysis)
			}
		}
	}
}

func GetProfile(db *gorm.DB, userId uint) ([]Mcq.MCQ, []WeeklyMcqAnalysis, []MonthlyMcqAnalysis, []YearlyMcqAnalysis, [][]Mcq.McqResult, []TopicAnalysis) {
	// Collect personal information
	var topics []TopicAnalysis // Collect topic analysis
	db.Where("user_id = ?", userId).Find(&topics)
	database := "ltq4ywpuwsubopkz." // database name
	//local := ""
	var collectionQuestions []Mcq.MCQ // Collect questions
	db.Where("mcq_id IN (SELECT mcq_id FROM " + database + "mcqs)").
		Preload("McqQuestions").
		Find(&collectionQuestions)

	var collectionWeek []WeeklyMcqAnalysis // Collect weekly analysis
	db.Where("user_id = ?", userId).
		Preload("WeeklyMcqAnalysisResults").
		Find(&collectionWeek)

	var collectionMonth []MonthlyMcqAnalysis // Collect monthly analysis
	db.Where("user_id = ?", userId).
		Preload("MonthlyMcqAnalysisResults").
		Find(&collectionMonth)

	var collectionYearly []YearlyMcqAnalysis // Collect yearly analysis
	db.Where("user_id = ?", userId).
		Preload("YearlyMcqAnalysisResults").
		Find(&collectionYearly)

	var distinctMcq []Mcq.McqResult // Collect distinct mcq ids
	var mcqResults [][]Mcq.McqResult
	db.Select("DISTINCT(mcq_id)").Order("mcq_id").
		Where("created_at >= ? AND user_id = ?", time.Now().Add(-365*(24*time.Hour)), userId).
		Find(&distinctMcq)
	for _, result := range distinctMcq { // Collect user results relating to mcq id
		var temp []Mcq.McqResult
		db.Where("user_id = ? AND created_at > ? AND mcq_id = ?", userId,
			time.Now().Add(-365*(24*time.Hour)), result.McqId).
			Order("created_at").
			Find(&temp)
		mcqResults = append(mcqResults, temp)
	}
	return collectionQuestions, collectionWeek, collectionMonth, collectionYearly, mcqResults, topics
}

func AnalyseGroups(db *gorm.DB) { // Analyse all groups
	var groups []g.Group
	db.Find(&groups)               // Groups
	for _, group := range groups { // For each Group
		var distinctTopicAnalysis []TopicAnalysis
		db.
			Order("topic_id").
			Where("user_id = ?", group.CreatorId).
			Find(&distinctTopicAnalysis) // Get all topics of creator
		var members []g.Member
		db.Where("group_id = ?", group.GroupId).
			Find(&members) // Get all members of group
		for _, topic := range distinctTopicAnalysis { // For each topic
			var cumulative float64
			var count float64
			for _, member := range members { // For each member
				var chosen TopicAnalysis
				db.Where("user_id = ? AND topic_id = ?", member.UserId, topic.TopicId).
					First(&chosen) // Get topic analysis for member for selected topic
				if chosen.TopicId != 0 {
					cumulative += chosen.AvgResult
					count += 1
				}
			}
			avgResult := cumulative / count        // Get average
			var groupAnalysis g.GroupTopicAnalysis // Get analysis associated with group id and topic id
			db.Where("group_id = ? AND topic_id = ?", group.GroupId, topic.TopicId).
				First(&groupAnalysis)
			if groupAnalysis.GTopId != 0 { // if analysis for this group topic
				// update analysis
				db.Table("group_topic_analyses").Updates(&g.GroupTopicAnalysis{AvgResult: avgResult}).Where("g_top_id = ?", groupAnalysis.GTopId)
			} else {
				// create analysis
				db.Create(&g.GroupTopicAnalysis{GroupId: group.GroupId, TopicId: topic.TopicId, TopicName: topic.TopicName, CreatedAt: time.Now(), AvgResult: avgResult})
			}
		}
	}
}

//////////  Private Methods  //////////

func collectDistinctMcq(db *gorm.DB, timePeriod time.Time) (collection []Mcq.McqResult) {
	// Collects all uniques ids associated with each mcq
	db.Select("DISTINCT(mcq_id)").Order("mcq_id").Where("created_at >= ? ", timePeriod).Find(&collection)
	return collection
}

func checkTodayAnalysis(db *gorm.DB, user userApi.User, mcqId uint) error {
	// Check if user has topic analysis and create analysis if none exists
	var mcqData Mcq.MCQ // Get Mcq
	db.Where("mcq_id = ?", mcqId).First(&mcqData)
	var topicInfo Topic // Get topic of Mcq
	db.Where("topic_name = ?", mcqData.Topic).First(&topicInfo)
	// Check user has weekly Analysis
	var QId []Mcq.McqQuestion
	db.Select("q_id").Where("mcq_id = ?", mcqId).Find(&QId)
	var WeeklyAnalysis WeeklyMcqAnalysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcqId).First(&WeeklyAnalysis).Error
	if err != nil || WeeklyAnalysis.McqId == 0 { // Create analysis if none exist
		log.Printf("No Weekly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&WeeklyMcqAnalysis{UserId: user.UserId, McqId: mcqId, WeeklyMcqAnalysisResults: []WeeklyMcqAnalysisResult{}, Topic: mcqData.Topic, TopicId: topicInfo.TopicId, AvgResult: 0})
	}
	// Create empty weekly Analysis
	var question []Mcq.McqQuestion
	db.Where("mcq_id = ?", mcqId).Find(&question)
	numberOfQuestions := len(question)

	db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcqId).First(&WeeklyAnalysis)
	// Check user has weekly Analysis for this mcq question
	var questionAnalysis []WeeklyMcqAnalysisResult
	db.Where("mcq_id = ? AND weekly_r_ana = ?", mcqId, WeeklyAnalysis.WeeklyRAna).Find(&questionAnalysis)
	if len(questionAnalysis) == 0 { // Create analysis records
		log.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ {
			var question Mcq.McqQuestion
			db.Where("q_id = ?", QId[index].QId).First(&question)
			err = db.Create(&WeeklyMcqAnalysisResult{
				WeeklyRAna: WeeklyAnalysis.WeeklyRAna, McqId: mcqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: QId[index].QId, Question: question.Question}).Error
			if err != nil {
				log.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}
	return nil
}

func checkMonthlyAnalysis(db *gorm.DB, user userApi.User, distinctMcq uint) error {
	// Check monthly analysis
	var mcq Mcq.MCQ // Checked Weekly
	db.Where("mcq_id = ?", distinctMcq).First(&mcq)

	var MonthlyAnalysis MonthlyMcqAnalysis // Check User has monthly analysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq).First(&MonthlyAnalysis).Error
	if err != nil || MonthlyAnalysis.McqId == 0 { // Create empty Monthly Analysis Results
		log.Printf("No Monthly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&MonthlyMcqAnalysis{UserId: user.UserId, McqId: mcq.McqId, MonthlyMcqAnalysisResults: []MonthlyMcqAnalysisResult{}, Topic: mcq.Topic, AvgResult: 0})
		db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq).First(&MonthlyAnalysis)
	}
	var Week WeeklyMcqAnalysis // Collect weekly analysis of this mcq
	db.Where("mcq_id = ? AND user_id = ?", distinctMcq, user.UserId).First(&Week)
	// Collect weekly analysis of this mcq
	var weeklyQuestionAnalysis []WeeklyMcqAnalysisResult
	db.Where("weekly_r_ana = ?", Week.WeeklyRAna).Find(&weeklyQuestionAnalysis)
	numberOfQuestions := len(weeklyQuestionAnalysis)

	// Check user has monthly Analysis for this mcq question
	var monthlyQuestionAnalysis []MonthlyMcqAnalysisResult
	err = db.Where("mcq_id = ? AND monthly_r_ana = ?", distinctMcq, MonthlyAnalysis.MonthlyRAna).Find(&monthlyQuestionAnalysis).Error
	if len(monthlyQuestionAnalysis) == 0 || err != nil { // If no question analysis
		log.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ { // for each question
			var question Mcq.McqQuestion // Find question
			db.Where("q_id = ?", weeklyQuestionAnalysis[index].QId).First(&question)
			err = db.Create(&MonthlyMcqAnalysisResult{
				MonthlyRAna: MonthlyAnalysis.MonthlyRAna, McqId: mcq.McqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: weeklyQuestionAnalysis[index].QId, Question: question.Question}).Error
			if err != nil {
				log.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}
	return nil
}

func checkYearlyAnalysis(db *gorm.DB, user userApi.User, distinctMcq uint) error {
	// Check if user has yearly analysis
	var mcq Mcq.MCQ // Checked Monthly
	db.Where("mcq_id = ?", distinctMcq).First(&mcq)

	var YearlyAnalysis YearlyMcqAnalysis // Check User has yearly analysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq).First(&YearlyAnalysis).Error
	if err != nil || YearlyAnalysis.McqId == 0 { // If no yearly analysis
		log.Printf("No Yearly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&YearlyMcqAnalysis{UserId: user.UserId, McqId: mcq.McqId, YearlyMcqAnalysisResults: []YearlyMcqAnalysisResult{}, Topic: mcq.Topic, AvgResult: 0})
		db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq).First(&YearlyAnalysis)
	}
	var Monthly MonthlyMcqAnalysis // Find monthly analysis with Mcq id and User id
	db.Where("mcq_id = ? AND user_id = ?", distinctMcq, user.UserId).First(&Monthly)
	var monthlyQuestionAnalysis []MonthlyMcqAnalysisResult // Create empty Monthly Analysis Results
	db.Where("monthly_r_ana = ?", Monthly.MonthlyRAna).Find(&monthlyQuestionAnalysis)
	numberOfQuestions := len(monthlyQuestionAnalysis)

	var yearlyQuestionAnalysis []YearlyMcqAnalysisResult // Check user has weekly Analysis for this mcq question
	err = db.Where("mcq_id = ? AND yearly_r_ana = ?", distinctMcq, YearlyAnalysis.YearlyRAna).Find(&yearlyQuestionAnalysis).Error
	if len(yearlyQuestionAnalysis) == 0 || err != nil { // If no yearly question analysis
		log.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ { // For each question
			var question Mcq.McqQuestion // Creat yearly question analysis
			db.Where("q_id = ?", monthlyQuestionAnalysis[index].QId).First(&question)
			err = db.Create(&YearlyMcqAnalysisResult{
				YearlyRAna: YearlyAnalysis.YearlyRAna, McqId: mcq.McqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: monthlyQuestionAnalysis[index].QId, Question: question.Question}).Error
			if err != nil {
				log.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}
	return nil
}

/*  // Not utilised due to time constraints
func checkTotalAnalysis(db *gorm.DB, user userApi.User, distinctMcq Mcq.MCQ) error {
	return nil
}
*/

func getWeeklyAnalysis(db *gorm.DB, M []Mcq.McqResult) WeeklyMcqAnalysis {
	// Run analysis
	var weekAvg WeeklyMcqAnalysis // Get Week Analysis
	db.Where("mcq_id = ? AND user_id = ?", M[0].McqId, M[0].UserId).
		Preload("WeeklyMcqAnalysisResults").
		First(&weekAvg)
	weekAvg.LastModified = time.Now()
	for _, result := range M { // For each Question Result
		for i, answer := range result.McqQuestionResult { // Get weekly Analysis for Question Result
			numberOfResults := weekAvg.WeeklyMcqAnalysisResults[i].NumberOfResults
			numberOfResults += 1

			currentAvgTime := weekAvg.WeeklyMcqAnalysisResults[i].AvgTime
			currentAvgResult := weekAvg.WeeklyMcqAnalysisResults[i].AvgResult * float64(numberOfResults-1) // Normalise current average result

			thisResult := float64(result.McqQuestionResult[i].Result) / float64(answer.Total)
			if currentAvgTime == 0 { // If no current average time, current average time duplicated
				currentAvgTime = answer.Time
			}
			newAvgTime := (currentAvgTime + answer.Time) / time.Duration(2)
			newAvgResult := (currentAvgResult + thisResult) / float64(numberOfResults)

			weekAvg.WeeklyMcqAnalysisResults[i].AvgTime = newAvgTime
			weekAvg.WeeklyMcqAnalysisResults[i].AvgResult = newAvgResult
			weekAvg.WeeklyMcqAnalysisResults[i].NumberOfResults += 1
			ConfDesc, Confidence := getConfidence(result.McqQuestionResult[i].Changes, result.McqQuestionResult[i].Time)
			weekAvg.WeeklyMcqAnalysisResults[i].AvgConfidence = Confidence
			weekAvg.WeeklyMcqAnalysisResults[i].AvgConfidenceString = ConfDesc
		}
	}
	var cumulative float64
	for _, Question := range weekAvg.WeeklyMcqAnalysisResults {
		db.Model(&Question).
			Updates(WeeklyMcqAnalysisResult{AvgTime: Question.AvgTime,
				AvgResult: Question.AvgResult, NumberOfResults: Question.NumberOfResults,
				AvgConfidenceString: Question.AvgConfidenceString, AvgConfidence: Question.AvgConfidence})
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(weekAvg.WeeklyMcqAnalysisResults))
	db.Table("weekly_mcq_analyses").Updates(WeeklyMcqAnalysis{LastModified: weekAvg.LastModified, AvgResult: average}).Where("weekly_r_ana = ?", weekAvg.WeeklyRAna)
	return weekAvg
}

func getMonthlyAnalysis(db *gorm.DB, weekAnalysis WeeklyMcqAnalysis) {
	// Run analysis
	userId, mcqId := weekAnalysis.UserId, weekAnalysis.McqId
	var monthAvg MonthlyMcqAnalysis // Get Month Analysis
	db.Where("mcq_id = ? AND user_id = ?", mcqId, userId).
		Preload("MonthlyMcqAnalysisResults").
		First(&monthAvg)
	monthAvg.LastModified = time.Now()

	for index, weeklyAnalysisResult := range weekAnalysis.WeeklyMcqAnalysisResults { // For each weekly analysis
		currentAvgTime := monthAvg.MonthlyMcqAnalysisResults[index].AvgTime
		currentAvgResult := monthAvg.MonthlyMcqAnalysisResults[index].AvgResult
		previousWeek := float64(weeklyAnalysisResult.AvgResult)
		if currentAvgTime == 0 { // If no current average time, current average time duplicated
			currentAvgTime = weeklyAnalysisResult.AvgTime
		}
		numberOfResults := float64(monthAvg.MonthlyMcqAnalysisResults[index].NumberOfResults)
		newAvgTime := (currentAvgTime + weeklyAnalysisResult.AvgTime) / time.Duration(2)
		newAvgResult := ((currentAvgResult * numberOfResults) + previousWeek) / (numberOfResults + 1)

		if numberOfResults < 4 { // Prevent analysis of more than four weeks
			monthAvg.MonthlyMcqAnalysisResults[index].NumberOfResults += 1
		}
		monthAvg.MonthlyMcqAnalysisResults[index].AvgTime = newAvgTime
		monthAvg.MonthlyMcqAnalysisResults[index].AvgResult = newAvgResult
		monthAvg.MonthlyMcqAnalysisResults[index].AvgConfidenceString = weeklyAnalysisResult.AvgConfidenceString
		monthAvg.MonthlyMcqAnalysisResults[index].AvgConfidence = weeklyAnalysisResult.AvgConfidence
	}
	var cumulative float64
	for _, Question := range monthAvg.MonthlyMcqAnalysisResults { // For each question analysis
		// Update analysis for monthly question analysis with analysis id
		db.Model(&Question).
			Updates(MonthlyMcqAnalysisResult{AvgTime: Question.AvgTime, AvgResult: Question.AvgResult,
				NumberOfResults: Question.NumberOfResults, AvgConfidenceString: Question.AvgConfidenceString, AvgConfidence: Question.AvgConfidence}).
			Where("where monthly_r_ana = ?", monthAvg.MonthlyRAna)
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(monthAvg.MonthlyMcqAnalysisResults))
	// Update Monthly analysis
	db.Table("monthly_mcq_analyses").Updates(MonthlyMcqAnalysis{LastModified: monthAvg.LastModified, AvgResult: average}).Where("monthly_r_ana = ?", monthAvg.MonthlyRAna)
}

func getYearlyAnalysis(db *gorm.DB, monthlyAnalysis MonthlyMcqAnalysis) {
	// Run analysis
	userId, mcqId := monthlyAnalysis.UserId, monthlyAnalysis.McqId
	var yearlyAvg YearlyMcqAnalysis // Get Yearly Analysis
	db.Where("mcq_id = ? AND user_id = ?", mcqId, userId).
		Preload("YearlyMcqAnalysisResults").
		First(&yearlyAvg)
	yearlyAvg.LastModified = time.Now()

	for index, monthlyAnalysisResult := range monthlyAnalysis.MonthlyMcqAnalysisResults { // For each monthly analysis question
		currentAvgTime := yearlyAvg.YearlyMcqAnalysisResults[index].AvgTime
		currentAvgResult := yearlyAvg.YearlyMcqAnalysisResults[index].AvgResult
		previousWeek := float64(monthlyAnalysisResult.AvgResult)
		if currentAvgTime == 0 { // If no current average time, current average time duplicated
			currentAvgTime = monthlyAnalysisResult.AvgTime
		}
		numberOfResults := float64(yearlyAvg.YearlyMcqAnalysisResults[index].NumberOfResults)
		newAvgTime := (currentAvgTime + monthlyAnalysisResult.AvgTime) / time.Duration(2)
		newAvgResult := ((currentAvgResult * numberOfResults) + previousWeek) / (numberOfResults + 1)

		if numberOfResults < 12 { // Limit analysis to twelve months
			yearlyAvg.YearlyMcqAnalysisResults[index].NumberOfResults += 1
		}
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgTime = newAvgTime
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgResult = newAvgResult
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgConfidence = monthlyAnalysisResult.AvgConfidence
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgConfidenceString = monthlyAnalysisResult.AvgConfidenceString

	}
	var cumulative float64
	for _, Question := range yearlyAvg.YearlyMcqAnalysisResults { // Update each question analysis
		db.Model(&Question).
			Updates(YearlyMcqAnalysisResult{AvgTime: Question.AvgTime, AvgResult: Question.AvgResult,
				NumberOfResults: Question.NumberOfResults, AvgConfidenceString: Question.AvgConfidenceString, AvgConfidence: Question.AvgConfidence})
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(yearlyAvg.YearlyMcqAnalysisResults))
	// update yearly analysis
	db.Table("yearly_mcq_analyses").Updates(YearlyMcqAnalysis{LastModified: yearlyAvg.LastModified, AvgResult: average}).Where("yearly_r_ana = ?", yearlyAvg.YearlyRAna)

}

/*  // Not utilised due to time constraints
func getTotalAnalysis(db *gorm.DB, yearlyAnalysis YearlyMcqAnalysis) {
	userId, mcqId := yearlyAnalysis.UserId, yearlyAnalysis.McqId
	var totalAvg TotalMcqAnalysis // Get Month Analysis
	db.Where("mcq_id = ? AND user_id = ?", mcqId, userId).
		Preload("TotalMcqAnalysisResults").
		First(&totalAvg)
	totalAvg.LastModified = time.Now()

	for index, monthlyAnalysisResult := range yearlyAnalysis.YearlyMcqAnalysisResults {
		currentAvgTime := totalAvg.TotalMcqAnalysisResults[index].AvgTime
		currentAvgResult := totalAvg.TotalMcqAnalysisResults[index].AvgResult
		previousWeek := float64(monthlyAnalysisResult.AvgResult)
		if currentAvgTime == 0 {
			currentAvgTime = monthlyAnalysisResult.AvgTime // If no current average time, current average time duplicated
		}
		numberOfResults := float64(totalAvg.TotalMcqAnalysisResults[index].NumberOfResults)
		newAvgTime := (currentAvgTime + monthlyAnalysisResult.AvgTime) / time.Duration(2)
		newAvgResult := ((currentAvgResult * numberOfResults) + previousWeek) / (numberOfResults + 1)

		// FIXME: Does the number of results work here?
		totalAvg.TotalMcqAnalysisResults[index].NumberOfResults += 1
		totalAvg.TotalMcqAnalysisResults[index].AvgTime = newAvgTime
		totalAvg.TotalMcqAnalysisResults[index].AvgResult = newAvgResult
	}
	var cumulative float64
	for _, Question := range totalAvg.TotalMcqAnalysisResults {
		db.Model(&Question).
			Updates(TotalMcqAnalysisResult{AvgTime: Question.AvgTime, AvgResult: Question.AvgResult, NumberOfResults: Question.NumberOfResults})
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(totalAvg.TotalMcqAnalysisResults))
	db.Model(TotalMcqAnalysisResult{}).Updates(TotalMcqAnalysis{LastModified: totalAvg.LastModified, AvgResult: average}).Where("total_r_ana = ?", totalAvg.TotalRAna)

}
*/

func getTopicAnalysis(db *gorm.DB, userId uint) {
	// Get analysis for topic for user
	var topics []Topic
	db.Find(&topics)               // Find all topics
	for _, topic := range topics { // For each topic
		var weekAvg []WeeklyMcqAnalysis
		db.Where("user_id = ? AND topic_id = ?", userId, topic.TopicName).
			Find(&weekAvg) // Find topics this user has tested
		var cumulative float64
		var count float64
		if len(weekAvg) != 0 { // User has at least one analysis in this topic
			for _, analysis := range weekAvg { // For each week analysis
				if analysis.AvgResult != 0 { // If result
					cumulative += analysis.AvgResult
					count += 1
				}
			}
			avg := cumulative / count
			var topicAnalysis TopicAnalysis
			db.Where("topic_id = ? AND user_id = ?", topic.TopicId, userId).
				First(&topicAnalysis) // Get topic analysis
			if topicAnalysis.TopicId != 0 { // Update analysis
				db.Table("topic_analyses").Update("avg_result", avg).Where("user_id = ? AND topic_id = ?", userId, topic.TopicId).
					Where("topic_ana_id = ?", topic.TopicId)
			} else { // Create analysis
				db.Create(&TopicAnalysis{UserId: userId, TopicId: topic.TopicId, TopicName: topic.TopicName, AvgResult: avg})
			}

		}
	}
}

func getConfidence(numberOfAlterations int, timeUsed time.Duration) (confidenceLevel string, confidence float64) {
	// Calculate Confidence
	VeryHigh := float64(0.85)
	High := float64(0.65)
	SomeWhat := float64(0.50)
	Low := float64(0.25)
	VeryLow := float64(0)
	idealTime := time.Duration(10 * time.Second)
	idealNumberOfAlterations := 1.1
	if numberOfAlterations < 2 { // Provides  100% Confidence if result changed once
		numberOfAlterations = 0
	}
	confidenceNumberOfAlterations := math.Pow(idealNumberOfAlterations, -float64(numberOfAlterations))
	timeUsed = time.Duration(timeUsed * 1000000) //  time to seconds
	var confidenceInTime float64
	if timeUsed < idealTime { // Prevent greater than 100% confidence
		confidenceInTime = math.Pow(float64(idealTime), -0)
	} else {
		confidenceInTime = math.Pow(float64(idealTime), -(float64(timeUsed - idealTime)))
	}
	confidence = (confidenceNumberOfAlterations + confidenceInTime) / 2 // Get confidence
	// Get string value
	if confidence > VeryHigh {
		confidenceLevel = "Very High"
	} else if confidence > High {
		confidenceLevel = "High"
	} else if confidence > SomeWhat {
		confidenceLevel = "Some What"
	} else if confidence > Low {
		confidenceLevel = "Low"
	} else if confidence >= VeryLow {
		confidenceLevel = "Very Low"
	}
	return confidenceLevel, confidence
}

/*  // Not utilised due to time constraints
func Normalise(yAxis []float64) {
	min := yAxis[0]
	max := yAxis[0]
	for _, y := range yAxis {
		if y < min {
			min = y
		}
		if y > max {
			max = y
		}
	}
	var normalised []float64
	for _, y := range yAxis {
		normal := (y-min)/max - min
		normalised = append(normalised, normal)
	}

}
*/

/*  // Not utilised due to time constraints
func getLastWeeksWorseQuestions(db *gorm.DB, McqId uint) map[string]WeeklyMcqAnalysisResult {
	QuestionAndResult := make(map[string]WeeklyMcqAnalysisResult)
	var results []WeeklyMcqAnalysisResult
	db.
		Order("avg_result ASC").
		Where("mcq_id = ?", McqId).
		Limit(5).
		Find(&results)
	var questions []Mcq.McqQuestion
	for _, qid := range results {
		var question Mcq.McqQuestion
		db.Where("mcq_id = ? AND q_id = ?", McqId, qid.QId).First(&question)
		questions = append(questions, question)
		QuestionAndResult[question.Question] = qid
	}
	fmt.Println("Worse Results ", QuestionAndResult)
	return QuestionAndResult

}
*/
