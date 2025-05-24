# 🏨 Hostel Management System

A web-based **Hostel Management System** built using **Golang** and **MySQL**, designed to streamline hostel operations such as student registration, leave requests, room allocations, and complaint handling. The system supports role-based logins for students and wardens ensuring each role accesses the relevant functionalities.

---

## 🚀 Features

### 👨‍🎓 Student
- Secure login and session handling
- View profile (Name, Room No., Gender, Contact Info)
- Apply for leave with date range & reason
- View leave request history with status
- Logout functionality

### 🧑‍🏫 Warden
- View list of student leave requests
- Approve/Reject leave requests with one click
- Visual status indicators (Green: Approved, Red: Rejected)

### 🧑‍💼 UWD *(Planned for future scope)*
- Dashboard to view overall hostel stats
- Manage wardens and rooms

---

## ⚙️ Tech Stack

| Layer        | Technology Used          |
|--------------|---------------------------|
| **Frontend** | HTML, CSS (Bootstrap), JavaScript |
| **Backend**  | Go (Golang)               |
| **Database** | MySQL                     |
| **Security** | JWT Authentication, Session Cookies |

---

## 🗂️ Folder Structure

```plaintext
hostel-management-system/
│
├── main.go                  # Entry point of the application
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies lock file
│
├── handlers/               # Backend handlers
│   ├── auth.go             # Handles login and registration
│   ├── dashboard.go        # Renders student and warden dashboards
│   └── leave.go            # Handles leave request logic
│
├── db/                     # Database access logic
│   ├── connection.go       # DB connection setup
│   ├── user.go             # User CRUD operations
│   └── leave.go            # Leave CRUD operations
│
├── templates/              # HTML templates
│   ├── login.html
│   ├── register.html
│   ├── student_dashboard.html
│   ├── warden_dashboard.html
│   ├── apply_leave.html
│   └── view_leave_handle.html
│
├── resource/               # Static assets
│   ├── css/
│   │   └── styles.css      # Custom styles
│   ├── js/
│   │   └── security.js      
│   └── images/
│       └── GGSIU_logo.svg.png
│
├── screenshots/            # UI screenshots for README
│   ├── login.png
│   ├── student.png
│   ├── leave.png
│   └── warden.png
│
├── db/
│   └── schema.sql          # Database schema
│
└── README.md               # Project documentation

---
```

## 🧑‍💻 Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/https-ritwick/Hostel-Leave-System_Demo
cd hostel-management-system
```
### 2. Set Up MySQL Database
*Import the SQL schema provided in /db/schema.sql
*Update your DB credentials in db/connection.go

### 3. Run the server using 
```bash
go run main.go
```
### 4. Access the system via: http://localhost:8000

## 👥 Authors

### Ritwick Johari  
[![LinkedIn - Ritwick](https://img.shields.io/badge/LinkedIn-blue?style=flat&logo=linkedin&label=Connect)](https://www.linkedin.com/in/ritwick-johari)

### Khushi Thakur  
[![LinkedIn - Khushi](https://img.shields.io/badge/LinkedIn-blue?style=flat&logo=linkedin&label=Connect)](https://www.linkedin.com/in/khushi-thakur-a031a9298/)

Developed as part of the B.Tech (AI & DS) at **GGSIPU – USAR** 🎓
