package main

import (
	"encoding/json"
	"fmt"
	aly "github.com/Aol1234/studentLearningGoServer/analysis"
	dev "github.com/Aol1234/studentLearningGoServer/devRoom"
	mcq "github.com/Aol1234/studentLearningGoServer/questionnaire"
	authApi "github.com/Aol1234/studentLearningGoServer/studentAuth"
	userApi "github.com/Aol1234/studentLearningGoServer/userModel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
	"net/http"
)

func main() {
	db, err := gorm.Open("mysql", "root:pass@/StudentLearning?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	http.HandleFunc("/post", func(w http.ResponseWriter, req *http.Request) {
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
		var testMCQ mcq.MCQ
		test := db.Where("mcq_id = ?", 1).Preload("McqQuestions").Preload("McqQuestions.Answers").Find(&testMCQ) //db.Where("mcq_id = ?", 1).Find(&testMCQ)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(test)
		if err != nil {
			panic(err)
		}
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
		return
	})

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
		fmt.Println("Caught")
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
		profileInfo := aly.GetProfile(db, userId)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(profileInfo)
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
		req.Body.Close()
	})

	http.HandleFunc("/studentAuth/Login", func(w http.ResponseWriter, req *http.Request) {
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

	http.HandleFunc("/admin/collectUserData", func(w http.ResponseWriter, req *http.Request) {

		// TODO: Add user authentication, select user
		aly.CollectData(db) // Needs to specify user && mcq
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = req.Body.Close()
		if err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
