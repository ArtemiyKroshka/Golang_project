package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	database := os.Getenv("DATABASE")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping the DB:", err)
	}

	fmt.Println("Successfully connected to the database!")
}

func ClearTable() {
	_, err := DB.Query("TRUNCATE TABLE golang_project")
	if err != nil {
		log.Fatal("Failed to truncate data:", err)
	}
}

func GetData() ([]string, error) {
	rows, err := DB.Query("SELECT line FROM golang_project")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		lines = append(lines, line)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %v", err)
	}

	return lines, nil
}

func SetData(line string) {
	_, err := DB.Exec("INSERT INTO golang_project (line) VALUES ($1)", line)
	if err != nil {
		log.Fatal("Failed to insert data:", err)
	}
	fmt.Println("Inserted line:", line)
}

func DeleteData(line string) {
	_, err := DB.Exec("DELETE FROM golang_project WHERE line = $1", line)
	if err != nil {
		log.Fatal("Failed to delete data:", err)
	}
	fmt.Println("Deleted line:", line)
}
