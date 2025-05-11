package db

import (
	"crypto/sha512"
	"fmt"
)

func ValidateUserCredentials(userName, password, nonce string) (validateflag bool, err error) {
	validateflag = false
	query := "SELECT ul.userpassword FROM user_details ud JOIN user_login ul ON ud.id=ul.user_details_ref_id where ul.username=? AND ud.validflag=1 AND ul.validflag=1 LIMIT 1"
	var dbHashedPassword string
	err = Con.QueryRow(query, userName).Scan(&dbHashedPassword)
	if err != nil {
		return
	}
	randomDbHashPassword := sha512Hash(dbHashedPassword + nonce)
	if password == randomDbHashPassword {
		validateflag = true
	}
	return
}
func sha512Hash(input string) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil)) //%x is for Hexa decimal string
}
