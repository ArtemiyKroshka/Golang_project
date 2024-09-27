package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Line struct {
	ID    int
	Title string
}

var DB *sql.DB

func createTable() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS golang_project_table (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func ConnectDatabase() {
	if os.Getenv("DB_HOST") == "" {
		// createTable()
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

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

	createTable()
}

func ExitDatabase() {
	ClearTable()

	defer DB.Close()
}

func ClearTable() {
	_, err := DB.Exec("TRUNCATE TABLE golang_project_table")
	if err != nil {
		log.Fatal("Failed to truncate data:", err)
	}
}

func GetData() ([]Line, error) {
	// Adjust query to select both id and title columns
	rows, err := DB.Query("SELECT id, title FROM golang_project_table")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var projects []Line
	for rows.Next() {
		var project Line
		// Scan both the id and title into the Project struct
		if err := rows.Scan(&project.ID, &project.Title); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		projects = append(projects, project)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %v", err)
	}

	return projects, nil
}

func SetData(line string) {
	_, err := DB.Exec("INSERT INTO golang_project_table (title) VALUES ($1)", line)
	if err != nil {
		log.Fatal("Failed to insert data:", err)
	}
}

func DeleteData(id string) {
	_, err := DB.Exec("DELETE FROM golang_project_table WHERE id = $1", id)
	if err != nil {
		log.Fatal("Failed to delete data:", err)
	}
}
