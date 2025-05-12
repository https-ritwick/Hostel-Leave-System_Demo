/*
package db

import "fmt"

func CreateUser(name, email string) (id int64, err error) {
	query := "INSERT INTO user_details(name, email, validflag) VALUES(?, ?, 1)"
	result, err := Con.Exec(query, name, email)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err = result.LastInsertId()
	fmt.Println("Inserted record ID:", id)
	return
}
func CreateUserWithCredentials(name, email, username, password string) (id int64, err error) {
	id, err = CreateUser(name, email)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err = CreateLoginCredentials(username, password, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
func CreateLoginCredentials(username, password string, user_details_ref_id int64) (lastInsertId int64, err error) {
	query := "INSERT INTO user_login( username, userpassword, validflag,user_details_ref_id) VALUES(?, ?, 1, ?)"
	result, err := Con.Exec(query, username, password, user_details_ref_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	lastInsertId, err = result.LastInsertId()
	fmt.Println("Inserted record ID:", lastInsertId)
	return
}

func PrintUserDetails() {
	query := "SELECT id, name, email, validflag FROM user_details"
	rows, err := Con.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var id, validflag int32
	var name, email string
	for rows.Next() {
		rows.Scan(&id, &name, &email, &validflag)
		fmt.Println("Rows ID:", id, " - ", name, " - ", email, " - ", validflag)
	}
}

func PrintUserDetailsForSelectedUserName(userName string) {
	query := "SELECT id, name, email, validflag FROM user_details where name like ?"
	rows, err := Con.Query(query, "%"+userName+"%")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var id, validflag int32
	var name, email string
	for rows.Next() {
		rows.Scan(&id, &name, &email, &validflag)
		fmt.Println("Rows ID:", id, " - ", name, " - ", email, " - ", validflag)
	}
}
*/

package db

import "fmt"

// Inserts into user_details all user information including role and gender
func CreateUserWithAllDetails(name, email, username, password, room, address, contact, role, gender string) (id int64, err error) {
	query := `
		INSERT INTO user_details(name, email, room_number, address, contact_number, role, gender, validflag)
		VALUES(?, ?, ?, ?, ?, ?, ?, 1)`
	result, err := Con.Exec(query, name, email, room, address, contact, role, gender)
	if err != nil {
		fmt.Println("CreateUserWithAllDetails error:", err)
		return
	}
	id, err = result.LastInsertId()
	if err != nil {
		return
	}
	return CreateLoginCredentials(username, password, id)
}

// Adds user login credentials linked to user_details
func CreateLoginCredentials(username, password string, user_details_ref_id int64) (lastInsertId int64, err error) {
	query := "INSERT INTO user_login(username, userpassword, validflag, user_details_ref_id) VALUES(?, ?, 1, ?)"
	result, err := Con.Exec(query, username, password, user_details_ref_id)
	if err != nil {
		fmt.Println("CreateLoginCredentials error:", err)
		return
	}
	lastInsertId, err = result.LastInsertId()
	fmt.Println("Inserted record ID (user_login):", lastInsertId)
	return
}

// Prints basic user details
func PrintUserDetails() {
	query := "SELECT id, name, email, room_number, contact_number, role, gender, validflag FROM user_details"
	rows, err := Con.Query(query)
	if err != nil {
		fmt.Println("PrintUserDetails error:", err)
		return
	}
	defer rows.Close()

	var id, validflag int32
	var name, email, room, contact, role, gender string
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &room, &contact, &role, &gender, &validflag)
		if err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		fmt.Printf("ID: %d | Name: %s | Email: %s | Room: %s | Contact: %s | Role: %s | Gender: %s | Valid: %d\n",
			id, name, email, room, contact, role, gender, validflag)
	}
}

func FetchUserRoleByUsername(username string) (role string, err error) {
	query := `
		SELECT ud.role 
		FROM user_details ud 
		JOIN user_login ul ON ud.id = ul.user_details_ref_id 
		WHERE ul.username = ? AND ud.validflag = 1 AND ul.validflag = 1
		LIMIT 1`
	err = Con.QueryRow(query, username).Scan(&role)
	if err != nil {
		fmt.Println("FetchUserRoleByUsername error:", err)
	}
	return
}

type StudentProfile struct {
	ID         int
	Name       string
	RoomNumber string
	Contact    string
	Gender     string
}

func FetchStudentProfile(username string) (StudentProfile, error) {
	var profile StudentProfile

	query := `
		SELECT ud.id, ud.name, ud.room_number, ud.contact_number, ud.gender
		FROM user_details ud
		JOIN user_login ul ON ul.user_details_ref_id = ud.id
		WHERE ul.username = ? AND ud.validflag = 1 AND ul.validflag = 1
		LIMIT 1`

	err := Con.QueryRow(query, username).Scan(
		&profile.ID,
		&profile.Name,
		&profile.RoomNumber,
		&profile.Contact,
		&profile.Gender,
	)

	if err != nil {
		fmt.Println("FetchStudentProfile error:", err)
	}
	return profile, err
}

type LeaveEntry struct {
	Reason   string
	FromDate string
	ToDate   string
	Status   string
}

func FetchLeaveHistory(studentID int) ([]LeaveEntry, error) {
	query := `
		SELECT reason, from_date, to_date, status
		FROM leave_requests
		WHERE student_id = ?
		ORDER BY id DESC`

	rows, err := Con.Query(query, studentID)
	if err != nil {
		fmt.Println("FetchLeaveHistory error:", err)
		return nil, err
	}
	defer rows.Close()

	var history []LeaveEntry
	for rows.Next() {
		var entry LeaveEntry
		if err := rows.Scan(&entry.Reason, &entry.FromDate, &entry.ToDate, &entry.Status); err != nil {
			fmt.Println("Row scan error in FetchLeaveHistory:", err)
			continue
		}
		history = append(history, entry)
	}

	return history, nil
}
