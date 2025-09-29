package models

import (
	"time"
)

/*Representa un usuario*/
type User struct {
	IdUser        int       `json:"iduser" db:"idUser" binding:"required,iduser"`
	IdRol         int16     `json:"idrol" db:"idRol" binding:"required,min=1,max=10"`
	NameUser      string    `json:"nameuser" db:"NameUser" binding:"required,nameuser"`
	Email         string    `json:"email" db:"Email" binding:"required,email"`
	LastName      string    `json:"lastname" db:"LastName" binding:"required,lastname"`
	LastActivitie time.Time `json:"lastactivitie" db:"LastActivitie" binding:"required,LastActivitie"`
	DateCreated   time.Time `json:"datecreated" db:"DateCreated" binding:"required,LastActivitie"`
	PasswordHash  string    `json:"passwordhash" db:"Password" binding:"required,passwordhash"`
}

/*Representa un Rol*/
type Rol struct {
	IdRol          int16  `json:"idrol" db:"idRol" binding:"required,idRol"`
	DescriptionRol string `json:"descriptionrol" db:"DescriptionRol" binding:"required,descriptionrol"`
	IsActive       bool   `json:"isactive" db:"IsActive" binding:"required,isActive"`
}

/*RESPUESTAS DE LAS OPERACIONES CON LOS MODELOS*/
type CreateUserRequest struct {
	IdRol         int16     `json:"idrol" binding:"required,min=1,max=10"`
	NameUser      string    `json:"nameuser" binding:"required,nameuser"`
	Email         string    `json:"email" db:"Email" binding:"required,email"`
	LastName      string    `json:"lastname" binding:"required,lastname"`
	LastActivitie time.Time `json:"lastactivitie" binding:"required,LastActivitie"`
	Password      string    `json:"password" binding:"required,password"`
}

type UpdateUserRequest struct {
	NameUser string `json:"nameuser" binding:"required,nameuser"`
	Email    string `json:"email" binding:"required,email"`
	LastName string `json:"lastname" binding:"required,lastname"`
}

type UserResponse struct {
	IdUser        int       `json:"iduser"`
	IdRol         int16     `json:"idrol"`
	NameUser      string    `json:"nameuser"`
	Email         string    `json:"email"`
	LastName      string    `json:"lastname"`
	LastActivitie time.Time `json:"lastactivitie"`
	DateCreated   time.Time `json:"datecreated"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required, min=6"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

/*METODOS DEL MODELO USER*/

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		IdUser:        u.IdUser,
		IdRol:         u.IdRol,
		NameUser:      u.NameUser,
		Email:         u.Email,
		LastName:      u.LastName,
		LastActivitie: u.LastActivitie,
		DateCreated:   u.DateCreated,
	}
}

func (u *User) GetFullNameUser() string {
	return u.NameUser + " " + u.LastName
}

func (u *User) SetPassword(password string) bool {
	u.Password = password
	return true
}

func NewUserFromRequest(req *CreateUserRequest) (*User, error) {
	return &User{
		IdRol:         req.IdRol,
		NameUser:      req.NameUser,
		Email:         req.Email,
		LastName:      req.LastName,
		LastActivitie: req.LastActivitie,
		Password:      req.Password,
	}, nil
}

func ValidateRol(idRol int) bool {
	totalrol := []int{1, 2, 3, 4, 5}
	for _, rol := range totalrol {
		if idRol == rol {
			return true
		}
	}
	return false
}

/*Paginacion*/

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Offset   int `json:"-"`
}

type UserFilter struct {
	IdRol  *int   `json:"idrol,omitempty" form:"id_rol"`
	Name   string `json:"name,omitempty" form:"name"`
	Search string `json:"search,omitempty" form:"search"` // Búsqueda general

}

func PaginationNew(page int, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

func (p *Pagination) CalculatePages(totalRecords int) int {
	if totalRecords == 0 || p.PageSize == 0 {
		return 0
	}
	return (totalRecords + p.PageSize - 1) / p.PageSize
}
