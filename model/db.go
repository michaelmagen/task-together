package model

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
)

var db *sql.DB

func InitDB() {
	// Read the connection string from Viper
	connectionString := viper.GetString("dbConnectionString")
	if connectionString == "" {
		log.Fatal("Connection string not found in configuration")
	}

	// Open a database connection
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	// Read the SQL schema file
	sqlFile, err := os.ReadFile("model/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL statements
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Database connected successfully!")
}

// GetDB returns the *sql.DB object for other packages to use
func GetDB() *sql.DB {
	return db
}
