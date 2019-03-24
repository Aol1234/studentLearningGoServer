package devRoom

import (
	aly "github.com/Aol1234/studentLearningGoServer/analysis"
	g "github.com/Aol1234/studentLearningGoServer/groups"
	mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"time"
)

// ** Inside a http handle func** ///
// 	db, err := gorm.Open("mysql", "root:@/studentlearning?charset=utf8&parseTime=True&loc=Local")
//	if err != nil {
//		panic("failed to connect database")
//	}
//	defer db.Close()
// Migrate the schema
//db.AutoMigrate(&product.Product{})
// Create
//db.Create(&Product{Code: "L1212", Price: 1000})
// Read
//var product Product
//db.First(&product, 1) // find product with id 1
//db.First(&product, "code = ?", "L1212") // find product with code l1212
// Update - update product's price to 2000
//db.Model(&product).Update("Price", 2000)
// Delete - delete product
//db.Delete(&product)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Body struct {
	// Needs to Capitalize all JSON variables
	ObjectOne string   `json:"object_one"`
	ObjectTwo string   `json:"object_two"`
	ArrayOne  ArrayOne `json:"array_one"`
}

type ArrayOne struct {
	A string `json:"a"`
	B string `json:"b"`
}

type FirebaseToken struct {
	Idtoken string `json:"idtoken"`
}

type MCQResults map[string]string

// Dictionary{ array[dictionary {k:v k:v ...} dictionary {k:v k:v} ...] array[dictionary {k:v} ] ... }
type Questionnaire map[int]map[int]map[string]string
type Identifier struct {
	Id string `json:"id"`
}
type Q struct {
	Identifier `json:"identifier"`
	Questions  `json:"Questions"`
}
type Questions []Que
type Que struct {
	Question string
	Answers
}
type Answer struct {
	Text  string `json:"text, omitempty"`
	Value string `json:"value, omitempty"`
}

type Answers []Answer

var Test2 = Q{
	Identifier{"123456789"},
	Questions{
		Que{"Question 1", Answers{Answer{"Radio  One", "0"}, Answer{"Radio Two", "1"}}},
		Que{"Question 2", Answers{Answer{"Radio One", "0"}, Answer{"Radio Two", "2"}, Answer{"Radio Three", "3"}}},
	},
}
var Mcq = Questionnaire{
	0: {0: {"text": "Test First radio", "value": "0"}, 1: {"text": "Second radio", "value": "1"}, 2: {"text": "Third radio", "value": "2"}},
	1: {0: {"text": "aaa radio", "value": "1", "misc": "testing that other attributes can be added"}, 1: {"text": "bbb radio", "value": "0"}, 2: {"text": "ccc radio", "value": "2"}},
	2: {0: {"text": "111 radio", "value": "2"}, 1: {"text": "222 radio", "value": "1"}, 2: {"text": "333 radio", "value": "0"}}}

//c:= map[string]string{"text": "First radio", "value": "0" }
//m := body{"ResponseOne", "ResponseTwo", ArrayOne{"Aone", "Bone"}}
/*
var Mc = Q{
Que { q.Answer {"Test First radio", "0"}, q.Answer { "Second radio", "1" }, q.Answer { "Third radio", "2" }},
Que { q.Answer {"Test First radio", "0"}, q.Answer { "Second radio", "1" }, q.Answer { "Third radio", "2" }},
Que { q.Answer {"Test First radio", "0"}, q.Answer { "Second radio", "1" }, q.Answer { "Third radio", "2" }},}
*/
type QuestionsResults struct {
	Results `json:"result"`
}

type Results []Result

type Result struct {
	ID        uint   `gorm:"primary_key"`
	ResultID  int    `gorm:"index"`
	Value     string `json:"value"`
	TotalTime int    `json:"total_time"`
}
type Test4 struct {
	Results `json:"result"`
}

type Test6 struct {
	Id      uint `gorm:"primary_key;AUTO_INCREMENT"`
	Results []Result
}

var Test1 = QuestionsResults{
	Results{
		Result{Value: "e", TotalTime: 1},
	},
}
var Test8 = Test6{
	//	id,
	//		[]Result{{"e", 2}, {"3", 3}},
}

/*
var (
	TestQuestions = []McqQuestion{
	{Question:"What is X?"},
	{Question:"What is Y?"},
	}
)
*/

/*
	// Migrate the schema
	db.AutoMigrate(&mcq.MCQ{}, &mcq.McqQuestion{}, &mcq.McqAnswer{})
	// Create
	/*
	db.Create(&mcq.MCQ{
		McqQuestions:[]mcq.McqQuestion{
			{Question:"What is x",
				Answers:[]mcq.McqAnswer{{Text:"true", Value:0}, {Text:"false", Value:1}},
			},
			{Question:"What is y",
				Answers:[]mcq.McqAnswer{{Text:"true", Value:0}, {Text:"false", Value:1}},
			},
			{Question:"What is z",
				Answers:[]mcq.McqAnswer{{Text:"true", Value:0}, {Text:"false", Value:1}},
			},
		},
	})
	//
	// Get MCQ Where mcq_id = ?
*/
var resultTest map[string]string

var Test3 = mcq.McqResult{
	1, 1, 1, 0.55, time.Now(), []mcq.McqQuestionResult{
		{1, 1, 1, 4, 4, time.Duration(20), 0},
	},
}

func SetUp(db *gorm.DB) {
	db.AutoMigrate(aly.Topic{}, aly.TopicAnalysis{}, aly.WeeklyMcqAnalysis{}, aly.WeeklyMcqAnalysisResult{},
		aly.MonthlyMcqAnalysis{}, aly.MonthlyMcqAnalysisResult{}, aly.YearlyMcqAnalysis{}, aly.YearlyMcqAnalysisResult{},
		aly.TotalMcqAnalysis{}, aly.TotalMcqAnalysisResult{})
	db.AutoMigrate(g.Group{}, g.Member{}, g.GroupTopicAnalysis{})
}
