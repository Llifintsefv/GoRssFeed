package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Llifintsefv/GoRssFeed/internal/rss/repository"
	"github.com/Llifintsefv/GoRssFeed/internal/rss/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB,error) {
	err := godotenv.Load(".env") // Load environment variables from .env
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == ""{
		return nil,fmt.Errorf("DB_HOST,DB_PORT,DB_USER,DB_PASSWORD,DB_NAME,DB_SSL_MODE must be set")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db,err := sql.Open("postgres",connStr)
	if err != nil {
		return nil,fmt.Errorf("failed to open db connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil,fmt.Errorf("failed to ping db: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS news (id SERIAL PRIMARY KEY, title TEXT, Link TEXT, Published TIMESTAMP,description TEXT)")
	if err != nil {
		return nil,fmt.Errorf("failed to create table: %w", err)
	}
	return db,nil
}
func main() {
	db,err := InitDB()
	if err != nil {
		log.Fatal(err,"error init db")
	}
	defer db.Close()
	repo := repository.NewRepository(db)
	service := service.NewRssService(repo)
	allnews,err := service.GetAllNews()
	if err != nil {
		log.Fatal(err,"error get all news")
	}
	for _,item := range allnews {
		fmt.Println(item.Title,": ",item.Description,"\n------------------------")
	}

}