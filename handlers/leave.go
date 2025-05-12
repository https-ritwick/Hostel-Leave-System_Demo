package handlers

import (
	"DynamicWebsiteProject/db"
	"encoding/json"
	"fmt"
	"net/http"
)

type LeaveRequest struct {
	ID        int    `json:"id,omitempty"`     // for fetching from DB
	StudentID int    `json:"student_id"`       // required in request
	Reason    string `json:"reason"`           // required in request
	FromDate  string `json:"from_date"`        // date in YYYY-MM-DD format
	ToDate    string `json:"to_date"`          // date in YYYY-MM-DD format
	Status    string `json:"status,omitempty"` // used in response/view
}

// GET: Serve HTML page
// POST: Handle form submission
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

	// Fetch student ID
	profile, err := db.FetchStudentProfile(username)
	if err != nil {
		http.Error(w, "Could not fetch profile", http.StatusInternalServerError)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	reason := r.FormValue("reason")
	fromDate := r.FormValue("from_date")
	toDate := r.FormValue("to_date")

	// Insert into DB
	query := `INSERT INTO leave_requests (student_id, reason, from_date, to_date) VALUES (?, ?, ?, ?)`
	_, err = db.Con.Exec(query, profile.ID, reason, fromDate, toDate)
	if err != nil {
		http.Error(w, "Failed to apply for leave", http.StatusInternalServerError)
		fmt.Println("Leave insert error:", err)
		return
	}

	// Redirect back to dashboard
	http.Redirect(w, r, "/student_dashboard", http.StatusSeeOther)
}

// GET /view_leave?status=Pending
func ViewLeaveHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	fmt.Println("Fetching leave requests with status:", status)

	rows, err := db.Con.Query("SELECT id, student_id, reason, from_date, to_date, status FROM leave_requests WHERE status = ?", status)
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
		leaves = append(leaves, lr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaves)
}

// PUT /update_leave_status?id=1&status=Approved
func UpdateLeaveStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	_, err := db.Con.Exec("UPDATE leave_requests SET status = ? WHERE id = ?", status, id)
	if err != nil {
		http.Error(w, "Error updating leave status", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Leave status updated"))
}
