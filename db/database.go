package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

var Con *sql.DB

func InitConnection(connectionString string) (err error) {
	Con, err = sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Unable to open database:", err)
		return
	}

	// Test the connection
	err = Con.Ping()
	if err != nil {
		fmt.Println("Database Connection Failed:", err)
		return
	}
	fmt.Println("Connected to MySQL successfully!")
	return
}
