package user

import (
	"fmt"
	"practise/database"
	"practise/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

func Login(username, password string) (string,error) {
	//先匹配db
	// collection:= GetCollection("users")
	user,err:=database.FindUserByUsername(username)
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
	token,err := jwt.GenerateToken(username)
	if err != nil {
		return "", fmt.Errorf("generate Token fail: %v", err)
	}

	// models.UserArticle.CreatedBy = username
	// models.UserInfo = models.User{
	// 	Username: username,
	// }

	

	return token,nil
}