package main

import (
	"encoding/json"
	"fmt"
	alyApi "github.com/Aol1234/studentLearningGoServer/analysis"
	g "github.com/Aol1234/studentLearningGoServer/groups"
	mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	sesApi "github.com/Aol1234/studentLearningGoServer/sessions"
	authApi "github.com/Aol1234/studentLearningGoServer/studentAuth"
	userApi "github.com/Aol1234/studentLearningGoServer/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
)

func main() {
	// localDatabase := "root:pass@/StudentLearning?charset=utf8&parseTime=True&loc=Local"
	herokuDatabase := "o4fg1odluxwsx0cl:go6oeb46hzn6hccr@tcp(l6slz5o3eduzatkw.cbetxkdyhwsb.us-east-1.rds.amazonaws.com)/ltq4ywpuwsubopkz?charset=utf8&parseTime=True"
	db, err := gorm.Open("mysql", herokuDatabase)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	http.HandleFunc("/studentAuth/SignUp", func(w http.ResponseWriter, req *http.Request) {
		// Sign-up new user
		setupResponse(&w, req)
		migrate(db)                           // Create Database tables if not already created
		var requestBody authApi.FirebaseToken // Firebase Token
		decoder := json.NewDecoder(req.Body)  // Decode Token
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()                                // Create Context
		token, err := authApi.VerifyUser(ctx, requestBody.Idtoken) // Verify token
		if err != nil {
			panic(err)
		}
		userApi.CreateUser(db, token.UID) // Create user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/studentAuth/Login", func(w http.ResponseWriter, req *http.Request) {
		// Login User
		setupResponse(&w, req)
		migrate(db)
		var requestBody authApi.FirebaseToken // Firebase Token
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()                                // Create Context
		token, err := authApi.VerifyUser(ctx, requestBody.Idtoken) // Verify token
		if err != nil {
			panic(err)
		}
		user := userApi.LoginVerification(db, token.UID) // Verify user
		if user.UserId == 0 {
			log.Println("Failed to identify user")
			return
		}
		successful := sesApi.SetSession(user.UserId, requestBody.Idtoken) // Set session
		if successful != true {
			log.Println("Failed to set Session")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/publishMcq", func(w http.ResponseWriter, req *http.Request) {
		// Handles the storing of a newly created multiple choice questionnaire
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		_, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		var requestBody mcq.MCQ // Decode questionnaire
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		mcq.CreateMcq(db, requestBody) //Store questionnaire
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(requestBody)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/getMcqs", func(w http.ResponseWriter, req *http.Request) {
		// Collect all multiple choice questionnaires
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		_, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		collection := mcq.GetMcqs(db) // Collection all questionnaires

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(collection)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/getSelectedMcq", func(w http.ResponseWriter, req *http.Request) {
		// Retrieve questionnaire and all its questions and answers
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		_, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		var requestBody mcq.MCQ // Decode requested questionnaire
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		questionnaire := mcq.RetrieveMcq(db, requestBody.McqId) // Retrieve questions and answers of selected questionnaire
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(questionnaire)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/storeResult", func(w http.ResponseWriter, req *http.Request) {
		// Store results made by user
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		var requestBody mcq.McqResult // Decode result
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		requestBody.UserId = userId      // Identify who's result is being stored
		mcq.StoreResult(db, requestBody) // store result
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/getProfile", func(w http.ResponseWriter, req *http.Request) {
		// Get personal analysis data
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		mcqQuestions, week, month,
			year, mcqResults, topics := alyApi.GetProfile(db, userId) // Collect data relating to user

		profile := alyApi.Data{McqQuestions: mcqQuestions, Weekly: week,
			Monthly: month, Yearly: year, Results: mcqResults, Topics: topics} // structure data

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(profile)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/updateUserPreferences", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(verify)
			err = req.Body.Close()
			return
		}
		decoder := json.NewDecoder(req.Body) // Decode Preference
		var requestBody userApi.UserPreference
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		requestBody.UserId = userId
		err = userApi.UpdateUserPreferences(db, requestBody) // Update Preference
		if err != nil {
			log.Println("Failed to update Preference: ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/getUserPreferences", func(w http.ResponseWriter, req *http.Request) {
		// Retrieve Preference
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		user := userApi.UserPreference{UserId: userId}
		userPreference, err := userApi.RetrieveUserPreferences(db, user) // Retrieve Preference
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(userPreference)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/createUserGroup", func(w http.ResponseWriter, req *http.Request) {
		// Create Group
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		var requestBody g.Group
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&requestBody) // Decode group data
		if err != nil {
			panic(err)
		}
		code := g.CreateGroup(db, userId, requestBody.Name, requestBody.Desc)
		if code.IsNil() {
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(code)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/joinUserGroup", func(w http.ResponseWriter, req *http.Request) {
		// Allow user to join group
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		var requestBody g.Group
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		err = g.JoinGroup(db, requestBody.Code, userId) // Make user member of Group
		if err != nil {
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/viewUserGroups", func(w http.ResponseWriter, req *http.Request) {
		// Retrieve users groups
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		group, groupAnalysis, err := g.GetGroups(db, userId) // Retrieve groups
		if err != nil {
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(g.Data{Groups: group, GroupTopicAnalysis: groupAnalysis})
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	////////////////// ADMIN METHODS ///////////////////////
	http.HandleFunc("/admin/collectUserData", func(w http.ResponseWriter, req *http.Request) {
		// Analyse Weekly Data
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		var admin userApi.User // Get user's role
		db.Where("user_id = ? and role = ?", userId, "ADMIN").First(&admin)
		if admin.UserId == 0 { // Is user a administrator
			log.Println("User:", userId, " is not Administrator")
			err = req.Body.Close()
			return
		}
		var users []userApi.User // Collect all users not administrator
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			alyApi.CollectData(db, user.UserId, "Week") // Collect user results for week
		}
		alyApi.AnalyseGroups(db)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/admin/collectUserDataYear", func(w http.ResponseWriter, req *http.Request) {
		// Analyse Monthly Data
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		var admin userApi.User // Get user's role
		db.Where("user_id = ? and role = ?", userId, "ADMIN").First(&admin)
		if admin.UserId == 0 {
			log.Println("User:", userId, " is not Administrator")
			err = req.Body.Close()
			return
		}
		var users []userApi.User // Collect all users not administrator
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			alyApi.CollectData(db, user.UserId, "Year") // Collect user results for month
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/admin/collectUserDataMonth", func(w http.ResponseWriter, req *http.Request) {
		// Analyse Monthly Data
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			log.Println(err)
			err = req.Body.Close()
			return
		}
		var admin userApi.User // Get user's role
		db.Where("user_id = ? and role = ?", userId, "ADMIN").First(&admin)
		if admin.UserId == 0 {
			log.Println("User:", userId, " is not Administrator")
			err = req.Body.Close()
			return
		}
		var users []userApi.User // Collect all users not administrator
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			alyApi.CollectData(db, user.UserId, "Month") // Collect user results for month
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(getPort(), nil); err != nil {
		panic(err)
	}
}
func getPort() string {
	p := os.Getenv("PORT")
	fmt.Println("port", p)
	if p != "" {
		return ":" + p
	}
	return ":8000"
}
func verifyUser(bearer string) (userId uint, verify bool, err error) {
	if bearer == "" { // No Bearer Header
		return 0, false, nil
	}
	userId, verify = userApi.VerifyUserId(bearer) // Verify Header
	if verify != true {
		return 0, false, nil
	}
	return userId, true, nil
}

func migrate(db *gorm.DB) { // Create Tables in Database if it doesn't exist
	db.AutoMigrate(&userApi.User{}, &userApi.UserPreference{})
	db.AutoMigrate(&alyApi.WeeklyMcqAnalysis{}, &alyApi.WeeklyMcqAnalysisResult{},
		&alyApi.MonthlyMcqAnalysis{}, &alyApi.MonthlyMcqAnalysisResult{},
		&alyApi.YearlyMcqAnalysis{}, &alyApi.YearlyMcqAnalysisResult{},
		&alyApi.TotalMcqAnalysis{}, &alyApi.TotalMcqAnalysisResult{})
	db.AutoMigrate(&mcq.MCQ{}, &mcq.McqQuestion{}, &mcq.McqAnswer{},
		&mcq.McqResult{}, &mcq.McqQuestionResult{})
	db.AutoMigrate(&g.Group{}, &g.GroupTopicAnalysis{}, &g.Member{},
		&alyApi.Topic{}, &alyApi.TopicAnalysis{})
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
