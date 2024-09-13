package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UsersClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserInfoLoginRespone struct {
	ID           int    `json:"id"`
	UserName     string `json:"username"`
	LoginSession string `json:"login_session"`
	Token        string `json:"token"`
}

func GenerateJwtToken(userInfo UsersClaims, session string) (UserInfoLoginRespone, error) {

	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            userInfo.ID,
		"username":      userInfo.Username,
		"login_session": session,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	var userRes = UserInfoLoginRespone{}
	if err != nil {
		fmt.Println(err.Error())
		return userRes, err
	}

	userRes = UserInfoLoginRespone{
		ID:           userInfo.ID,
		UserName:     userInfo.Username,
		LoginSession: session,
		Token:        tokenString,
	}

	return userRes, err
}
