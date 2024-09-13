package users

import (
	"api_v2/pkg"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	UserService *Service
}

func NewController(s *Service) *Controller {
	return &Controller{
		UserService: s,
	}
}

func (s *Controller) GetUsers(c *fiber.Ctx) error {

	perPageParam := c.Query("perPage")
	perPage, errVal := strconv.Atoi(perPageParam)
	if errVal != nil {
		fmt.Println("error", errVal)
		perPage = 20
	}

	currentPageParam := c.Query("currentPage")
	currentPage, errVal := strconv.Atoi(currentPageParam)
	if errVal != nil {
		fmt.Println("error", errVal)
		currentPage = 1
	}

	params := UserParamsConfig{
		CurrentPage: currentPage,
		PerPage:     perPage,
	}
	user, err := s.UserService.GetUser(params)

	if err != nil {
		var successRespone = pkg.ResBuilder(fiber.StatusBadRequest, "Failed to get users", pkg.Null())
		return c.Status(fiber.StatusBadRequest).JSON(successRespone)
	}
	var successRespone = pkg.ResBuilder(fiber.StatusOK, pkg.SUCCESS, user)
	return c.Status(fiber.StatusOK).JSON(successRespone)
}

func (s *Controller) Register(c *fiber.Ctx) error {

	regisRequest := new(RegisterUser)

	if err := c.BodyParser(regisRequest); err != nil {
		fmt.Println("Error get params", err)
		return err
	}

	validErrors, mes := pkg.ValidateStuct(regisRequest)

	if len(validErrors) > 0 {
		var successRespone = pkg.ResBuilder(fiber.StatusBadRequest, mes, validErrors)
		return c.Status(fiber.StatusBadRequest).JSON(successRespone)
	}

	regis, err := s.UserService.Register(*regisRequest)
	if err != nil {
		var successRespone = pkg.ResBuilder(fiber.StatusBadRequest, "Registeration failed!", pkg.Null())
		return c.Status(fiber.StatusBadRequest).JSON(successRespone)
	}

	var successRespone = pkg.ResBuilder(fiber.StatusOK, pkg.SUCCESS, regis)
	return c.Status(fiber.StatusOK).JSON(successRespone)
}

func (s *Controller) Login(c *fiber.Ctx) error {

	params := new(UserLogin)
	if err := c.BodyParser(params); err != nil {
		fmt.Println("Error Login: get params", err)
		return err
	}

	//Validate params
	validErrors, mes := pkg.ValidateStuct(params)
	if len(validErrors) > 0 {
		var successRespone = pkg.ResBuilder(fiber.StatusBadRequest, mes, validErrors)
		return c.Status(fiber.StatusBadRequest).JSON(successRespone)
	}

	//Check is user exist
	userInfo, err := s.UserService.Login(*params)
	if err != nil || userInfo.Password != params.Password {
		var successRespone = pkg.ResBuilder(fiber.StatusUnauthorized, "Login failed!", pkg.Null())
		return c.Status(fiber.StatusUnauthorized).JSON(successRespone)
	}

	//Set new session for user
	newSession := uuid.New().String()
	sess, err := s.UserService.SetUserSession(userInfo.Username, newSession)
	if err != nil || !sess {
		var successRespone = pkg.ResBuilder(fiber.StatusUnauthorized, "Login failed!", pkg.Null())
		return c.Status(fiber.StatusUnauthorized).JSON(successRespone)
	}

	//Generate User Token
	userClaimsInfo := pkg.UsersClaims{
		ID:       userInfo.ID,
		Username: userInfo.Username,
	}

	resData, err := pkg.GenerateJwtToken(userClaimsInfo, newSession)
	if err != nil {
		var successRespone = pkg.ResBuilder(fiber.StatusUnauthorized, "Login failed!", pkg.Null())
		return c.Status(fiber.StatusUnauthorized).JSON(successRespone)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    resData.Token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   false,
		HTTPOnly: true,
	})

	var successRespone = pkg.ResBuilder(fiber.StatusOK, pkg.SUCCESS, resData)
	return c.Status(fiber.StatusOK).JSON(successRespone)
}
