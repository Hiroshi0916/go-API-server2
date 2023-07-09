package main

import (
	"net/http"

	"go_API_server/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//環境変数で設定。
	// .env ファイルをロード
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数を読み取り
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	// dsn := "user=api password=api dbname=api search_path=public sslmode=disable"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error DB connect")
	}

	// 環境変数を使用（ここではコンソールに表示）
	fmt.Println("DB Host:", dbHost)
	fmt.Println("DB User:", dbUser)
	fmt.Println("DB Pass:", dbPass)
	fmt.Printf("DB Host: %s\n", dbHost)
	fmt.Printf("DB Port: %s\n", dbPort)
	fmt.Printf("DB User: %s\n", dbUser)
	fmt.Printf("DB Name: %s\n", dbName)
	fmt.Printf("DB Pass: %s\n", dbPass)

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	router := routes.SetupRouter(db)
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
} //end of main
