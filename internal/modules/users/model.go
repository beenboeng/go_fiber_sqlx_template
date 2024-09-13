package users

type Users struct {
	ID        int     `db:"id" json:"id"`
	FirstName *string `db:"first_name" json:"first_name"`
	LastName  *string `db:"last_name" json:"last_name"`
	Username  string  `db:"username" json:"username"`
	Password  string  `db:"password" json:"password"`
	Email     string  `db:"email" json:"email"`
	Phone     *string `db:"phone" json:"phone"`
}

type RegisterUser struct {
	FirstName *string `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name" validate:"required"`
	Username  string  `json:"username" validate:"required"`
	Password  string  `json:"password" validate:"required,min=6,max=20"`
	Email     string  `json:"email" validate:"required"`
	Phone     *string `json:"phone" validate:"required"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserParamsConfig struct {
	PerPage     int `json:"perPage" validate:"required"`
	CurrentPage int `json:"currentPage" validate:"required"`
}

type PaginationConfig struct {
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	Total       int `json:"total"`
}

type UserInfoDataRespone struct {
	User       []Users          `json:"user"`
	Pagination PaginationConfig `json:"pagination"`
}
