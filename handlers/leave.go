package handlers

import (
	"DynamicWebsiteProject/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LeaveRequest struct {
	ID        int    `json:"id,omitempty"`
	StudentID int    `json:"student_id"`
	Reason    string `json:"reason"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	Status    string `json:"status,omitempty"`
}

// POST /apply_leave
func ApplyLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "ajaxjwtdemo.com")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := db.FetchStudentProfile(username)
	if err != nil {
		http.Error(w, "Could not fetch profile", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	reason := r.FormValue("reason")
	fromDate := r.FormValue("from_date")
	toDate := r.FormValue("to_date")

	query := `INSERT INTO leave_requests (student_id, reason, from_date, to_date) VALUES (?, ?, ?, ?)`
	_, err = db.Con.Exec(query, profile.ID, reason, fromDate, toDate)
	if err != nil {
		http.Error(w, "Failed to apply for leave", http.StatusInternalServerError)
		fmt.Println("Leave insert error:", err)
		return
	}

	http.Redirect(w, r, "/student_dashboard", http.StatusSeeOther)
}

// GET /view_leave?status=Pending
func ViewLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	fmt.Println("Requested status:", status)

	rows, err := db.Con.Query("SELECT id, student_id, reason, from_date, to_date, status FROM leave_requests WHERE TRIM(status) = ?", status)
	if err != nil {
		http.Error(w, "Error retrieving leaves", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leaves []LeaveRequest
	for rows.Next() {
		var lr LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.StudentID, &lr.Reason, &lr.FromDate, &lr.ToDate, &lr.Status); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		lr.Status = strings.TrimSpace(lr.Status)
		leaves = append(leaves, lr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaves)
}

// PUT /update_leave_status?id=1&status=Approved
func UpdateLeaveStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	_, err := db.Con.Exec("UPDATE leave_requests SET status = ? WHERE id = ?", strings.TrimSpace(status), id)
	if err != nil {
		http.Error(w, "Error updating leave status", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Leave status updated"))
}

// GET /view_all_leaves
func ViewAllLeavesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Con.Query("SELECT id, student_id, reason, from_date, to_date, status FROM leave_requests")
	if err != nil {
		http.Error(w, "Error retrieving all leaves", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leaves []LeaveRequest
	for rows.Next() {
		var lr LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.StudentID, &lr.Reason, &lr.FromDate, &lr.ToDate, &lr.Status); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		lr.Status = strings.TrimSpace(lr.Status)
		leaves = append(leaves, lr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaves)
}
