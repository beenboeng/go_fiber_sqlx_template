package users

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		db,
	}
}

func (s *Service) GetUser(params UserParamsConfig) (UserInfoDataRespone, error) {

	var offset = (params.CurrentPage - 1) * params.PerPage
	users := []Users{}
	err := s.db.Select(&users, `
		SELECT 
			u.id, 
			u.first_name, 
			u.last_name, 
			u.username,
			u.password,
			u.email,
			u.phone
		FROM tbl_users AS u
		ORDER BY u.id
		LIMIT $1
		OFFSET $2    
	`, params.PerPage, offset)

	if err != nil {
		fmt.Println("Error GetUser(): Get users", err.Error())
		return UserInfoDataRespone{}, err
	}

	total := 0
	errTotal := s.db.QueryRowx(`
		SELECT 
			count(*) AS total
		FROM tbl_users AS u
	`).Scan(&total)

	if errTotal != nil {
		fmt.Println("Error GetUser(): Get users", errTotal.Error())
		return UserInfoDataRespone{}, errTotal
	}

	return UserInfoDataRespone{
		User: users,
		Pagination: PaginationConfig{
			PerPage:     params.PerPage,
			CurrentPage: params.CurrentPage,
			Total:       total,
		},
	}, nil
}

func (s *Service) Register(regis RegisterUser) (RegisterUser, error) {

	_, err := s.db.NamedExec(`
		INSERT INTO tbl_users
		(
			first_name,
			last_name,
			username,
			"password",
			email,
			phone
		)
		VALUES
		(:first_name, :last_name, :username, :password, :email, :phone);`,
		map[string]interface{}{
			"first_name": regis.FirstName,
			"last_name":  regis.LastName,
			"username":   regis.Username,
			"password":   regis.Password,
			"email":      regis.Email,
			"phone":      regis.Phone,
		},
	)

	if err != nil {
		return regis, err
	}

	return regis, nil
}

func (s *Service) Login(params UserLogin) (Users, error) {

	user := Users{}
	sqlErr := s.db.Get(&user, `
		SELECT u.id,
			u.first_name,
			u.last_name,
			u.username,
			u.password,
			u.email,
			u.phone
		FROM   tbl_users AS u
		WHERE  u.username = $1
	`, params.Username)

	if sqlErr != nil {
		fmt.Println("Error GetUser(): Get users", sqlErr.Error())
		return user, sqlErr
	}

	return user, nil
}

func (s *Service) SetUserSession(username, session string) (bool, error) {

	var rawSql = `UPDATE tbl_users 
			SET login_session = $1
			where tbl_users.username  = $2`

	_, err := s.db.Exec(rawSql, session, username)

	if err != nil {
		fmt.Println("Error in set user session", err)
		return false, err
	}

	return true, nil
}
