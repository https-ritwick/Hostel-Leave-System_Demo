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
		fmt.Println("Error connecting DB server:", err)
		return
	}
	defer db.Con.Close()

	// Static file handler
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("resource"))))

	// Public routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/login_check", handlers.LoginCheckHandler)

	// Authenticated routes
	http.HandleFunc("/", handlers.AuthMiddlewareCookie(handlers.ApplyLeaveHandler))
	http.HandleFunc("/apply_leave", handlers.AuthMiddlewareCookie(handlers.ApplyLeaveHandler))
	http.HandleFunc("/update_leave_status", handlers.AuthMiddlewareCookie(handlers.UpdateLeaveStatusHandler))
	http.HandleFunc("/student_dashboard", handlers.AuthMiddlewareCookie(handlers.StudentDashboardHandler))
	http.HandleFunc("/warden_dashboard", handlers.AuthMiddlewareCookie(handlers.WardenDashboardHandler))
	http.HandleFunc("/view_leave", handlers.AuthMiddlewareCookie(handlers.ViewLeaveHandler))
	http.HandleFunc("/view_all_leaves", handlers.AuthMiddlewareCookie(handlers.ViewAllLeavesHandler))

	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Warden-only protected route
	http.HandleFunc("/leave_handle", handlers.LeaveHandleHandler)

	fmt.Println("Server started at :localhost:8000/login")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
