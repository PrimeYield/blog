package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//POST localhost:8080/login
// the string in return is a token
func Login(username, password string) (string,error) {
	//先匹配db
	// collection:= GetCollection("users")
	user,err:=FindUsersByUsername(username)
	if err != nil {
		return "",err
	}
	if user == nil {
		return "",fmt.Errorf("%s is not exists",username)
		//也可以改成跳轉到CreatedUser之類的行為
	}
	// password驗證
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("wrong password")
	}

	return token,nil
}