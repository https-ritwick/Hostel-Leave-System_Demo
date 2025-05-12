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
