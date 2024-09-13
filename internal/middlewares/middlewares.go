package middlewares

import (
	"api_v2/database"
	"api_v2/pkg"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isCheckReqAuthentication(c) {
			return c.Next()
		}
		return abortRequest(c)
	}
}

func abortRequest(c *fiber.Ctx) error {
	var failedRespone = pkg.ResBuilder(fiber.StatusUnauthorized, pkg.FAIL, pkg.Null())
	return c.Status(fiber.StatusUnauthorized).JSON(failedRespone)
}

func isCheckReqAuthentication(c *fiber.Ctx) bool {

	fullToken := c.Get("Authorization")
	tokenCookies := c.Cookies("token")

	if tokenCookies != "" {
		tokenValidate, err := validateToken(c, tokenCookies)
		if err != nil {
			fmt.Println("Error in validate")
			return false
		}
		return tokenValidate
	} else {
		if strings.TrimSpace(fullToken) != "" {
			trimedToken := strings.TrimSpace(fullToken)
			if len(trimedToken) > 10 {
				userToken := strings.Split(trimedToken, " ")[1]

				tokenValidate, err := validateToken(c, userToken)

				if err != nil {
					fmt.Println("Error in validate")
					return false
				}
				return tokenValidate

			} else {
				return false
			}

		} else {
			fmt.Println("Authorization is empty")
			return false
		}
	}

}

type ClaimsObject struct {
	Exp          int
	ID           int    `json:"id"`
	UserName     string `json:"username"`
	LoginSession string `json:"login_session"`
	jwt.RegisteredClaims
}

func validateToken(c *fiber.Ctx, tokenStr string) (val bool, err error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsObject{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("Error in ValidateToken:", err)
		return false, err
	} else if claimsData, ok := token.Claims.(*ClaimsObject); ok {

		var userCtx = ClaimsObjectData{
			ID:       claimsData.ID,
			UserName: claimsData.UserName,
			Session:  claimsData.LoginSession,
			Exp:      claimsData.Exp,
		}

		//Set new context every time when request
		c.Locals("userCtx", userCtx)

		validateSession := ValidateUserSession(userCtx)
		if !validateSession {
			return false, err
		}
	} else {
		fmt.Println("unknown claims type, cannot proceed")
		return false, err
	}
	return true, err
}

type ClaimsObjectData struct {
	ID       int
	UserName string
	Session  string
	Exp      int
}

type Users struct {
	ID           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	LoginSession string `db:"login_session" json:"login_session"`
}

func ValidateUserSession(userReq ClaimsObjectData) bool {

	user := Users{}
	err := database.DB.Get(&user, `
		SELECT u.id,
			u.username,
			u.login_session
		FROM   tbl_users AS u
		WHERE  u.username = $1
	`, userReq.UserName)

	if err != nil {
		fmt.Println("Error in ValidateUserSession", err)
		return false
	}

	return strings.Contains(user.LoginSession, userReq.Session)

}
