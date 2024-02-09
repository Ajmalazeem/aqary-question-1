package main

import (
    "database/sql"
	"fmt"
    "log"
	"os"
    // "net/http"
	"question-1/route"
	"github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	loadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}
	defer db.Close()

// Create a Gin router
    router := gin.Default()

	controller := route.NewController(db)
	controller.RegisterRoutes(router)

//start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
