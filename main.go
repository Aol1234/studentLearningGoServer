package main

import (
	"encoding/json"
	"fmt"
	aly "github.com/Aol1234/studentLearningGoServer/analysis"
	dev "github.com/Aol1234/studentLearningGoServer/devRoom"
	g "github.com/Aol1234/studentLearningGoServer/groups"
	mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	authApi "github.com/Aol1234/studentLearningGoServer/studentAuth"
	userApi "github.com/Aol1234/studentLearningGoServer/userModel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
	"net/http"
	"os"
)

func main() {
	db, err := gorm.Open("mysql", "o4fg1odluxwsx0cl:go6oeb46hzn6hccr@tcp(l6slz5o3eduzatkw.cbetxkdyhwsb.us-east-1.rds.amazonaws.com)/ltq4ywpuwsubopkz?charset=utf8&parseTime=True") //root:pass@/StudentLearning?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	http.HandleFunc("/publishMcq", func(w http.ResponseWriter, req *http.Request) {
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		_, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		var requestBody mcq.MCQ
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		fmt.Println(requestBody)
		mcq.CreateMcq(db, requestBody)
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
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		_, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		Options := mcq.GrabMcqs(db)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Options)
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
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		_, verify, err := verifyUser(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		var requestBody mcq.MCQ
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		Options := mcq.RetrieveMcq(db, requestBody.McqId)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Options)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/:ID/:mcqID/result", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		var requestBody mcq.McqResult
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		requestBody.UserId = userId
		mcq.StoreResult(db, requestBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(dev.Test3)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

	http.HandleFunc("/getProfile", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		user, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		mcqQuestions, week, month, year, mcqResults, topics := aly.GetProfile(db, user)
		profile := aly.Data{McqQuestions: mcqQuestions, Weekly: week, Monthly: month, Yearly: year, Results: mcqResults, Topics: topics}
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

	http.HandleFunc("/studentAuth/SignUp", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		var requestBody dev.FirebaseToken
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		token, err := authApi.VerifyUser(ctx, requestBody.Idtoken)
		if err != nil {
			panic(err)
		}
		userApi.CreateUser(db, token.UID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/studentAuth/Login", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		migrate(db)
		var requestBody dev.FirebaseToken
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		fmt.Println(requestBody.Idtoken)
		ctx := context.Background()
		token, err := authApi.VerifyUser(ctx, requestBody.Idtoken)
		if err != nil {
			panic(err)
		}
		user := userApi.LoginVerification(db, token.UID)
		verify := userApi.SetCookie(user, requestBody.Idtoken)
		if verify != true {
			panic(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/updateUserPreferences", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			//w.WriteHeader(http.StatusBadRequest)
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			panic(err)
			//w.WriteHeader(http.StatusForbidden)
			return
		}
		decoder := json.NewDecoder(req.Body)
		var requestBody userApi.UserPreference
		err := decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		requestBody.UserId = userId
		err = userApi.UpdateUserPreferences(db, requestBody)
		if err != nil {
			//panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/getUserPreferences", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		userId, verify := userApi.VerifyUserId(bearer)
		if verify != true {
			w.WriteHeader(http.StatusOK)
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		user := userApi.UserPreference{UserId: userId}
		userPreference, err := userApi.RetrieveUserPreferences(db, user)
		if err != nil {
			panic(err)
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
		setupResponse(&w, req)

		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(req.Body)
		var requestBody g.Group
		err = decoder.Decode(&requestBody)
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
		setupResponse(&w, req)

		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			panic(err)
		}
		decoder := json.NewDecoder(req.Body)
		var requestBody g.Group
		err = decoder.Decode(&requestBody)
		if err != nil {
			panic(err)
		}
		err = g.JoinGroup(db, requestBody.Code, userId)
		if err != nil {
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/viewUserGroups", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		bearer := req.Header.Get("Authorization")
		if bearer == "" {
			return
		}
		userId, _, err := verifyUser(bearer)
		if err != nil {
			panic(err)
		}
		group, groupAnalysis, err := g.GetGroups(db, userId)
		if err != nil {
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
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
		setupResponse(&w, req)

		// TODO: Add user authentication, select user
		var users []userApi.User
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			aly.CollectData(db, user.UserId, "Week") // Needs to specify user && mcq
		}
		aly.AnalyseGroups(db)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/admin/collectUserDataYear", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		// TODO: Add user authentication, select user
		var users []userApi.User
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			aly.CollectData(db, user.UserId, "Year") // Needs to specify user && mcq
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/admin/collectUserDataMonth", func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		// TODO: Add user authentication, select user
		var users []userApi.User
		db.Where("role <> ?", "ADMIN").Find(&users)
		for _, user := range users {
			aly.CollectData(db, user.UserId, "Month") // Needs to specify user && mcq
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
	if bearer == "" {
		return 0, false, nil
	}
	userId, verify = userApi.VerifyUserId(bearer)
	if verify != true {
		return 0, false, nil
	}
	return userId, true, nil
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&userApi.User{}, &userApi.UserPreference{})
	db.AutoMigrate(&aly.WeeklyMcqAnalysis{}, &aly.WeeklyMcqAnalysisResult{},
		&aly.MonthlyMcqAnalysis{}, &aly.MonthlyMcqAnalysisResult{},
		&aly.YearlyMcqAnalysis{}, &aly.YearlyMcqAnalysisResult{},
		&aly.TotalMcqAnalysis{}, &aly.TotalMcqAnalysisResult{})
	db.AutoMigrate(&mcq.MCQ{}, &mcq.McqQuestion{}, &mcq.McqAnswer{},
		&mcq.McqResult{}, &mcq.McqQuestionResult{})
	db.AutoMigrate(&g.Group{}, &g.GroupTopicAnalysis{}, &g.Member{},
		&aly.TopicAnalysis{}, &aly.TopicAnalysis{})
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
