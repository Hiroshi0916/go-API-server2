package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ログイン情報
type User struct {
	ID       string `json:"login_id"`
	Password string `json:"password"`
}

// ログインIDとパスワード
type LoginRequest struct {
	LoginID  string `json:"loginID"`
	Password string `json:"password"`
}

// ユーザー認証用のハンドラ
func UserAuthHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Failed to decode user:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticated, err := authenticateUser(db, user)
	if err != nil {
		log.Println("Authentication error:", err)
		http.Error(w, "Authentication error", http.StatusInternalServerError)
		return
	}

	if !authenticated {
		log.Println("User not authenticated:", user.ID)
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// 認証が成功した場合の処理
	log.Println("User authenticated:", user.ID)
}

func LoginHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("Received a login request")

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Println("Invalid request body:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	result := db.Where("login_id = ?", loginRequest.LoginID).First(&user)
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		user := User{
			ID:       loginRequest.LoginID,
			Password: string(hashedPassword),
		}

		result = db.Create(&user)
		if result.Error != nil {
			log.Println("Error storing new user:", result.Error)
			http.Error(w, "error storing new user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully. Please login."))
	} else if result.Error != nil {
		log.Println("Error retrieving user:", result.Error)
		http.Error(w, "error retrieving user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			log.Println("Login error:", err)
			http.Error(w, "login error", http.StatusUnauthorized)
		} else {
			w.Write([]byte("Login successful"))

		}
	}
}

// ユーザーの認証(認証成功時:true, 認証失敗時:fail)
func authenticateUser(db *gorm.DB, user User) (bool, error) {
	// // データベースからユーザー情報を取得する
	var storedUser User

	// データベースからユーザー情報を取得する
	result := db.Select("password").Where("login_id = ?", user.ID).First(&storedUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// ユーザーが見つからない場合
		return false, nil
	} else if result.Error != nil {
		// 他のエラーが発生した場合
		return false, result.Error
	}

	// パスワードの検証
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		// パスワードが一致しない場合
		return false, nil
	}
	// パスワードが一致した場合
	return true, nil
}
