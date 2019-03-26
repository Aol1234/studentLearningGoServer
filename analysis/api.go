package analysis

import (
	"fmt"
	g "github.com/Aol1234/studentLearningGoServer/groups"
	Mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	userApi "github.com/Aol1234/studentLearningGoServer/userModel"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"time"
)

func NewSql(db *gorm.DB) *Sql {
	return &Sql{db: db}
}

// Collect mcq results for a particular user
func CollectData(db *gorm.DB, userId uint, timePeriod string) {
	db.AutoMigrate(&WeeklyMcqAnalysis{}, &WeeklyMcqAnalysisResult{},
		&MonthlyMcqAnalysis{}, &MonthlyMcqAnalysisResult{},
		&YearlyMcqAnalysis{}, &YearlyMcqAnalysisResult{},
		&TotalMcqAnalysis{}, &TotalMcqAnalysisResult{})

	var user userApi.User
	db.Where("user_id = ?", userId).First(&user)

	var period time.Time
	if timePeriod == "Month" {
		period = time.Now().Add(-30 * (24 * time.Hour))
	} else if timePeriod == "Week" {
		period = time.Now().Add(-7 * (24 * time.Hour))
	} else if timePeriod == "Year" {
		period = time.Now().Add(-365 * (24 * time.Hour))
	}

	var collection []Mcq.MCQ
	collection = collectDistinctMcq(db, period)
	if timePeriod == "Week" {
		for _, mcq := range collection {
			var results []Mcq.McqResult // Collect all mcq results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).Preload("McqQuestionResult").Find(&results)
			if len(results) > 0 { // Check if results not empty
				err := checkTodayAnalysis(db, user, mcq)
				if err != nil {
					log.Println(err)
				}
				getWeeklyAnalysis(db, results)
				// getLastWeeksWorseQuestions(db, mcq.McqId)
			}
		}
		getTopicAnalysis(db, userId)
	} else if timePeriod == "Month" {
		for _, mcq := range collection {
			var analysis WeeklyMcqAnalysis // Weekly Analysis results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).Preload("WeeklyMcqAnalysisResults").First(&analysis)
			if len(analysis.WeeklyMcqAnalysisResults) > 0 { // Check if results not empty
				err := checkMonthlyAnalysis(db, user, mcq)
				if err != nil {
					log.Println(err)
				}
				getMonthlyAnalysis(db, analysis)
				// getLastWeeksWorseQuestions(db, mcq.McqId)
			}
		}
	} else if timePeriod == "Year" {
		for _, mcq := range collection {
			var analysis MonthlyMcqAnalysis // Weekly Analysis results relating to this mcq
			db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).Preload("MonthlyMcqAnalysisResults").First(&analysis)
			if len(analysis.MonthlyMcqAnalysisResults) > 0 { // Check if results not empty
				err := checkYearlyAnalysis(db, user, mcq)
				if err != nil {
					log.Println(err)
				}
				getYearlyAnalysis(db, analysis)
				// getLastWeeksWorseQuestions(db, mcq.McqId)
			}
		}
	}
}

func GetProfile(db *gorm.DB, userId uint) ([]Mcq.MCQ, []WeeklyMcqAnalysis, []MonthlyMcqAnalysis, []YearlyMcqAnalysis, [][]Mcq.McqResult, []TopicAnalysis) {
	var topics []TopicAnalysis
	db.Where("user_id = ?", userId).Find(&topics)
	var collectionQuestions []Mcq.MCQ
	db.Where("mcq_id IN (SELECT mcq_id FROM ltq4ywpuwsubopkz.weekly_mcq_analyses WHERE user_id = ?)", userId).
		Preload("McqQuestions").
		Find(&collectionQuestions)
	var collectionWeek []WeeklyMcqAnalysis
	db.Where("user_id = ?", userId).
		Preload("WeeklyMcqAnalysisResults").
		Find(&collectionWeek)
	var collectionMonth []MonthlyMcqAnalysis
	db.Where("user_id = ?", userId).
		Preload("MonthlyMcqAnalysisResults").
		Find(&collectionMonth)
	var collectionYearly []YearlyMcqAnalysis
	db.Where("user_id = ?", userId).
		Preload("YearlyMcqAnalysisResults").
		Find(&collectionYearly)
	var distinctMcq []Mcq.McqResult
	var mcqResults [][]Mcq.McqResult
	db.Select("DISTINCT(mcq_id)").Order("mcq_id").
		Where("created_at >= ? AND user_id = ?", time.Now().Add(-365*(24*time.Hour)), userId).
		Find(&distinctMcq)
	for _, result := range distinctMcq {
		var temp []Mcq.McqResult
		db.Where("user_id = ? AND created_at > ? AND mcq_id = ?", userId,
			time.Now().Add(-365*(24*time.Hour)), result.McqId).
			Order("created_at").
			Find(&temp)
		mcqResults = append(mcqResults, temp)
	}
	fmt.Println("Results", collectionWeek)
	return collectionQuestions, collectionWeek, collectionMonth, collectionYearly, mcqResults, topics
}

func AnalyseGroups(db *gorm.DB) { // Analyse all groups
	var groups []g.Group
	db.Find(&groups)
	fmt.Println("Groups", groups)
	for _, group := range groups { // For each Group
		var distinctTopicAnalysis []TopicAnalysis
		db.
			Order("topic_id").
			Where("user_id = ?", group.CreatorId).
			Find(&distinctTopicAnalysis) // Get all topics of creator
		fmt.Println("Topics", distinctTopicAnalysis)
		var members []g.Member
		db.Where("group_id = ?", group.GroupId).
			Find(&members) // Get all members of group
		fmt.Println("Members", members)
		for _, topic := range distinctTopicAnalysis {
			var cumulative float64
			var count float64
			for _, member := range members { // For each member
				var chosen TopicAnalysis
				db.Where("user_id = ? AND topic_id = ?", member.UserId, topic.TopicId).
					First(&chosen) // Get topic analysis for member for selected topic
				fmt.Println("Chosen", chosen)
				if chosen.TopicId != 0 {
					cumulative += chosen.AvgResult
					count += 1
				}
			}
			avgResult := cumulative / count
			var groupAnalysis g.GroupTopicAnalysis
			db.Where("group_id = ? AND topic_id = ?", group.GroupId, topic.TopicId).
				First(&groupAnalysis)
			fmt.Println("Analysis", groupAnalysis)
			if groupAnalysis.GTopId != 0 {
			} else {
				db.Create(&g.GroupTopicAnalysis{GroupId: group.GroupId, TopicId: topic.TopicId, TopicName: topic.TopicName, CreatedAt: time.Now(), AvgResult: avgResult})
			}
		}
	}
}

//////////  Private Methods  //////////

func collectDistinctMcq(db *gorm.DB, timePeriod time.Time) (collection []Mcq.MCQ) {
	// Collects all uniques ids associated with each mcq
	db.Select("DISTINCT(mcq_id)").Order("mcq_id").Where("created_at >= ? ", timePeriod).Find(&collection)
	return collection
}

func checkTodayAnalysis(db *gorm.DB, user userApi.User, mcq Mcq.MCQ) error {
	fmt.Print("HERE! ", user.UserId, mcq.Name)
	var mcqData Mcq.MCQ
	db.Where("mcq_id = ?", mcq.McqId).First(&mcqData)
	var topicInfo Topic
	db.Where("topic_name = ?", mcqData.Topic).First(&topicInfo)
	// Check user has weekly Analysis
	var QId []Mcq.McqQuestion
	db.Select("q_id").Where("mcq_id	= ?", mcq.McqId).Find(&QId)
	var WeeklyAnalysis WeeklyMcqAnalysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).First(&WeeklyAnalysis).Error
	if err != nil || WeeklyAnalysis.McqId == 0 {
		fmt.Printf("No Weekly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&WeeklyMcqAnalysis{UserId: user.UserId, McqId: mcq.McqId, WeeklyMcqAnalysisResults: []WeeklyMcqAnalysisResult{}, Topic: mcqData.Topic, TopicId: topicInfo.TopicId, AvgResult: 0})
	}
	// Create empty weekly Analysis
	var question []Mcq.McqQuestion
	db.Where("mcq_id = ?", mcq.McqId).Find(&question)
	numberOfQuestions := len(question)

	db.Where("user_id = ? AND mcq_id = ?", user.UserId, mcq.McqId).First(&WeeklyAnalysis)
	// Check user has weekly Analysis for this mcq question
	var questionAnalysis []WeeklyMcqAnalysisResult
	db.Where("mcq_id = ? AND weekly_r_ana = ?", mcq.McqId, WeeklyAnalysis.WeeklyRAna).Find(&questionAnalysis)
	fmt.Println("Length of Questions: ", len(questionAnalysis))
	if len(questionAnalysis) == 0 {
		fmt.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ {
			var question Mcq.McqQuestion
			db.Where("q_id = ?", QId[index].QId).First(&question)
			err = db.Create(&WeeklyMcqAnalysisResult{
				WeeklyRAna: WeeklyAnalysis.WeeklyRAna, McqId: mcq.McqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: QId[index].QId, Question: question.Question}).Error
			if err != nil {
				fmt.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}

	return nil
}

func checkMonthlyAnalysis(db *gorm.DB, user userApi.User, distinctMcq Mcq.MCQ) error {
	// Checked Weekly
	var mcq Mcq.MCQ
	db.Where("mcq_id = ?", distinctMcq.McqId).First(&mcq)

	// Check User has
	var MonthlyAnalysis MonthlyMcqAnalysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq.McqId).First(&MonthlyAnalysis).Error
	if err != nil || MonthlyAnalysis.McqId == 0 { // Create empty Monthly Analysis Results
		fmt.Printf("No Monthly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&MonthlyMcqAnalysis{UserId: user.UserId, McqId: mcq.McqId, MonthlyMcqAnalysisResults: []MonthlyMcqAnalysisResult{}, Topic: mcq.Topic, AvgResult: 0})
		db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq.McqId).First(&MonthlyAnalysis)
	}

	// Collect weekly analysis of this mcq
	var weeklyQuestionAnalysis []WeeklyMcqAnalysisResult
	db.Where("mcq_id = ? AND user_id = ?", distinctMcq.McqId, user.UserId).Find(&weeklyQuestionAnalysis)
	numberOfQuestions := len(weeklyQuestionAnalysis)

	// Check user has weekly Analysis for this mcq question
	var monthlyQuestionAnalysis []MonthlyMcqAnalysisResult
	err = db.Where("mcq_id = ? AND user_id = ?", distinctMcq.McqId, user.UserId).Find(&monthlyQuestionAnalysis).Error
	if len(monthlyQuestionAnalysis) == 0 || err != nil {
		fmt.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ {
			var question Mcq.McqQuestion
			db.Where("q_id = ?", weeklyQuestionAnalysis[index].QId).First(&question)
			err = db.Create(&MonthlyMcqAnalysisResult{
				MonthlyRAna: MonthlyAnalysis.MonthlyRAna, McqId: mcq.McqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: weeklyQuestionAnalysis[index].QId, Question: question.Question}).Error
			if err != nil {
				fmt.Printf("Failed to create Question Analysis: %s", err)
				return err
			}
		}
	}
	return nil
}

func checkYearlyAnalysis(db *gorm.DB, user userApi.User, distinctMcq Mcq.MCQ) error {
	// Checked Monthly
	var mcq Mcq.MCQ
	db.Where("mcq_id = ?", distinctMcq.McqId).First(&mcq)

	// FIXME: Check user has Monthly Analysis

	// Check User has
	var YearlyAnalysis YearlyMcqAnalysis
	err := db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq.McqId).First(&YearlyAnalysis).Error
	if err != nil || YearlyAnalysis.McqId == 0 {
		fmt.Printf("No Yearly Analysis created for this: %s  Creating new Analysis for Exam", err)
		db.Create(&YearlyMcqAnalysis{UserId: user.UserId, McqId: mcq.McqId, YearlyMcqAnalysisResults: []YearlyMcqAnalysisResult{}, Topic: mcq.Topic, AvgResult: 0})
		db.Where("user_id = ? AND mcq_id = ?", user.UserId, distinctMcq.McqId).First(&YearlyAnalysis)
	}

	// Create empty Monthly Analysis Results
	var monthlyQuestionAnalysis []MonthlyMcqAnalysisResult
	db.Where("mcq_id = ? AND user_id = ?", distinctMcq.McqId, user.UserId).Find(&monthlyQuestionAnalysis)
	numberOfQuestions := len(monthlyQuestionAnalysis)

	// Check user has weekly Analysis for this mcq question
	var yearlyQuestionAnalysis []YearlyMcqAnalysisResult
	err = db.Where("mcq_id = ? AND user_id = ?", distinctMcq.McqId, user.UserId).Find(&yearlyQuestionAnalysis).Error
	if len(yearlyQuestionAnalysis) == 0 || err != nil {
		fmt.Printf("Failed to find Question Analysis: %s", err)
		for index := 0; index < numberOfQuestions; index++ {
			var question Mcq.McqQuestion
			db.Where("q_id = ?", monthlyQuestionAnalysis[index].QId).First(&question)
			err = db.Create(&YearlyMcqAnalysisResult{
				YearlyRAna: YearlyAnalysis.YearlyRAna, McqId: mcq.McqId, NumberOfResults: 0,
				AvgTime: time.Duration(0), AvgResult: 0.00, QId: monthlyQuestionAnalysis[index].QId, Question: question.Question}).Error
			if err != nil {
				fmt.Printf("Failed to create Question Analysis: %s", err)
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
	var weekAvg WeeklyMcqAnalysis
	// Get Week Analysis
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
	db.Model(WeeklyMcqAnalysis{}).Updates(WeeklyMcqAnalysis{LastModified: weekAvg.LastModified, AvgResult: average}).Where("weekly_r_ana = ?", weekAvg.WeeklyRAna)
	return weekAvg
}

func getMonthlyAnalysis(db *gorm.DB, weekAnalysis WeeklyMcqAnalysis) {
	userId, mcqId := weekAnalysis.UserId, weekAnalysis.McqId
	var monthAvg MonthlyMcqAnalysis // Get Month Analysis
	db.Where("mcq_id = ? AND user_id = ?", mcqId, userId).
		Preload("MonthlyMcqAnalysisResults").
		First(&monthAvg)
	monthAvg.LastModified = time.Now()

	for index, weeklyAnalysisResult := range weekAnalysis.WeeklyMcqAnalysisResults {
		currentAvgTime := monthAvg.MonthlyMcqAnalysisResults[index].AvgTime
		currentAvgResult := monthAvg.MonthlyMcqAnalysisResults[index].AvgResult
		previousWeek := float64(weeklyAnalysisResult.AvgResult)
		if currentAvgTime == 0 {
			currentAvgTime = weeklyAnalysisResult.AvgTime // If no current average time, current average time duplicated
		}
		numberOfResults := float64(monthAvg.MonthlyMcqAnalysisResults[index].NumberOfResults)
		newAvgTime := (currentAvgTime + weeklyAnalysisResult.AvgTime) / time.Duration(2)
		newAvgResult := ((currentAvgResult * numberOfResults) + previousWeek) / (numberOfResults + 1)

		if numberOfResults < 4 {
			monthAvg.MonthlyMcqAnalysisResults[index].NumberOfResults += 1
		}
		monthAvg.MonthlyMcqAnalysisResults[index].AvgTime = newAvgTime
		monthAvg.MonthlyMcqAnalysisResults[index].AvgResult = newAvgResult
		monthAvg.MonthlyMcqAnalysisResults[index].AvgConfidenceString = weeklyAnalysisResult.AvgConfidenceString
		monthAvg.MonthlyMcqAnalysisResults[index].AvgConfidence = weeklyAnalysisResult.AvgConfidence
	}
	var cumulative float64
	for _, Question := range monthAvg.MonthlyMcqAnalysisResults {
		db.Model(&Question).
			Updates(MonthlyMcqAnalysisResult{AvgTime: Question.AvgTime, AvgResult: Question.AvgResult,
				NumberOfResults: Question.NumberOfResults, AvgConfidenceString: Question.AvgConfidenceString, AvgConfidence: Question.AvgConfidence})
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(monthAvg.MonthlyMcqAnalysisResults))
	db.Model(MonthlyMcqAnalysis{}).Updates(MonthlyMcqAnalysis{LastModified: monthAvg.LastModified, AvgResult: average}).Where("monthly_r_ana = ?", monthAvg.MonthlyRAna)
}

func getYearlyAnalysis(db *gorm.DB, monthlyAnalysis MonthlyMcqAnalysis) {
	userId, mcqId := monthlyAnalysis.UserId, monthlyAnalysis.McqId
	var yearlyAvg YearlyMcqAnalysis // Get Month Analysis
	db.Where("mcq_id = ? AND user_id = ?", mcqId, userId).
		Preload("YearlyMcqAnalysisResults").
		First(&yearlyAvg)
	yearlyAvg.LastModified = time.Now()

	for index, monthlyAnalysisResult := range monthlyAnalysis.MonthlyMcqAnalysisResults {
		currentAvgTime := yearlyAvg.YearlyMcqAnalysisResults[index].AvgTime
		currentAvgResult := yearlyAvg.YearlyMcqAnalysisResults[index].AvgResult
		previousWeek := float64(monthlyAnalysisResult.AvgResult)
		if currentAvgTime == 0 {
			currentAvgTime = monthlyAnalysisResult.AvgTime // If no current average time, current average time duplicated
		}
		numberOfResults := float64(yearlyAvg.YearlyMcqAnalysisResults[index].NumberOfResults)
		newAvgTime := (currentAvgTime + monthlyAnalysisResult.AvgTime) / time.Duration(2)
		newAvgResult := ((currentAvgResult * numberOfResults) + previousWeek) / (numberOfResults + 1)

		if numberOfResults < 12 {
			yearlyAvg.YearlyMcqAnalysisResults[index].NumberOfResults += 1
		}
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgTime = newAvgTime
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgResult = newAvgResult
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgConfidence = monthlyAnalysisResult.AvgConfidence
		yearlyAvg.YearlyMcqAnalysisResults[index].AvgConfidenceString = monthlyAnalysisResult.AvgConfidenceString

	}
	var cumulative float64
	for _, Question := range yearlyAvg.YearlyMcqAnalysisResults {
		db.Model(&Question).
			Updates(YearlyMcqAnalysisResult{AvgTime: Question.AvgTime, AvgResult: Question.AvgResult,
				NumberOfResults: Question.NumberOfResults, AvgConfidenceString: Question.AvgConfidenceString, AvgConfidence: Question.AvgConfidence})
		cumulative += Question.AvgResult
	}
	average := cumulative / float64(len(yearlyAvg.YearlyMcqAnalysisResults))
	db.Model(YearlyMcqAnalysis{}).Updates(YearlyMcqAnalysis{LastModified: yearlyAvg.LastModified, AvgResult: average}).Where("yearly_r_ana = ?", yearlyAvg.YearlyRAna)

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
func getTopicAnalysis(db *gorm.DB, userId uint) { // Get analysis for topic for SELECTED user
	db.AutoMigrate(Topic{}, TopicAnalysis{})
	var topics []Topic
	db.Find(&topics)               // Find all topics
	for _, topic := range topics { // For each topic
		var weekAvg []WeeklyMcqAnalysis
		db.Where("user_id = ? AND topic = ?", userId, topic.TopicName).
			Find(&weekAvg) // Find topics this user has tested
		var cumulative float64
		var count float64
		if len(weekAvg) != 0 { // User has at least one analysis in this topic
			for _, analysis := range weekAvg { // For each week analysis
				if analysis.AvgResult != 0 {
					cumulative += analysis.AvgResult
					count += 1
				}
			}
			avg := cumulative / count
			var topicAnalysis TopicAnalysis
			db.Where("topic_id = ? AND user_id = ?", topic.TopicId, userId).
				First(&topicAnalysis)
			if topicAnalysis.TopicId != 0 {
				db.Table("topic_analyses").Update("avg_result", avg).Where("user_id = ? AND topic_id = ?", userId, topic.TopicId)
			} else {
				db.Create(&TopicAnalysis{UserId: userId, TopicId: topic.TopicId, TopicName: topic.TopicName, AvgResult: avg})
			}

		}
	}
}

func getConfidence(numberOfAlterations int, timeUsed time.Duration) (confidenceLevel string, confidence float64) {
	VeryHigh := float64(0.85)
	High := float64(0.65)
	SomeWhat := float64(0.50)
	Low := float64(0.25)
	VeryLow := float64(0)
	// Calculate Confidence
	idealTime := time.Duration(10 * time.Second)
	idealNumberOfAlterations := 1.1
	if numberOfAlterations < 2 { // Provides  100% Confidence if first choice chosen
		numberOfAlterations = 0
	}
	confidenceNumberOfAlterations := math.Pow(idealNumberOfAlterations, -float64(numberOfAlterations))
	fmt.Println("Confidence Of Alterations", confidenceNumberOfAlterations, numberOfAlterations)
	timeUsed = time.Duration(timeUsed * 1000)
	var confidenceInTime float64
	if timeUsed < idealTime {
		confidenceInTime = math.Pow(float64(idealTime), -0)
	} else {
		// TODO: This does not correctly calculate confidence, y needs to divided by 10
		confidenceInTime = math.Pow(float64(idealTime), -(float64(timeUsed - idealTime)))
	}
	fmt.Println("Confidence of Time", confidenceInTime, timeUsed)
	confidence = (confidenceNumberOfAlterations + confidenceInTime) / 2
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
	fmt.Println("Confidence", confidenceLevel, confidence)
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
