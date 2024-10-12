package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Line struct {
	ID    int
	Title string
}

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	if os.Getenv("DB_HOST") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open a DB connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	database := &Database{DB: db}
	if err := database.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return database, nil
}

func (d *Database) createTable() error {
	_, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS golang_project_table (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL
	);`)
	return err
}

func (d *Database) ClearTable() error {
	_, err := d.DB.Exec("TRUNCATE TABLE golang_project_table")
	return err
}

func (d *Database) GetData() ([]Line, error) {
	rows, err := d.DB.Query("SELECT id, title FROM golang_project_table")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var projects []Line
	for rows.Next() {
		var project Line
		if err := rows.Scan(&project.ID, &project.Title); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %v", err)
	}

	return projects, nil
}

func (d *Database) SetData(line string) error {
	_, err := d.DB.Exec("INSERT INTO golang_project_table (title) VALUES ($1)", line)
	return err
}

func (d *Database) DeleteData(id string) error {
	_, err := d.DB.Exec("DELETE FROM golang_project_table WHERE id = $1", id)
	return err
}

func (d *Database) Close() error {
	return d.DB.Close()
}
