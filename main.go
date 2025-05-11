package main

import (
	"DynamicWebsiteProject/db"
	"DynamicWebsiteProject/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var connectionString string

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handlers.JWTKey = []byte(os.Getenv("JWT_SECRET"))
	if len(handlers.JWTKey) == 0 {
		log.Fatal("JWT_SECRET not set in .env file")
	}
	connectionString = os.Getenv("CONN")
	if connectionString == "" {
		log.Fatal("DB Connection not set in .env file")
	}
}
func main() {
	err := db.InitConnection(connectionString)
	if err != nil {
		fmt.Println("Error conecting DB server:", err)
		return
	}
	defer db.Con.Close()

	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("resource"))))

	http.HandleFunc("/", handlers.AuthMiddlewareCookie(handlers.ApplyLeaveHandler))
	//http.HandleFunc("/quiz", handlers.AuthMiddlewareCookie(handlers.QuizPageHandler))
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/login_check", handlers.LoginCheckHandler)
	http.HandleFunc("/apply_leave", handlers.AuthMiddlewareCookie(handlers.ApplyLeaveHandler))
	http.HandleFunc("/view_leave", handlers.AuthMiddlewareCookie(handlers.ViewLeaveHandler))
	http.HandleFunc("/update_leave_status", handlers.AuthMiddlewareCookie(handlers.UpdateLeaveStatusHandler))
	http.HandleFunc("/leave_handle", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/leave_handle.html")
	})

	fmt.Println("Server started at :8000")
	err = http.ListenAndServe(":8000", nil) // Start the WEB Server AT PORT 8000
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
