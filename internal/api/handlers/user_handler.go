package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
	"github.com/MartinPatricio/GoGinAPISimple/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser godoc
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req db.CreateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserByID godoc
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// ... Implementación para obtener usuario por ID
}

// GetAllUsers godoc
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	nameFilter := c.Query("name")
	emailFilter := c.Query("email")

	offset := (page - 1) * pageSize

	var users []db.Tbluser
	var err error

	if nameFilter != "" || emailFilter != "" {
		params := db.GetUsersWithFiltersParams{
			Limit:       int32(pageSize),
			Offset:      int32(offset),
			NameFilter:  sql.NullString{String: "%" + nameFilter + "%", Valid: nameFilter != ""},
			EmailFilter: sql.NullString{String: "%" + emailFilter + "%", Valid: emailFilter != ""},
		}
		// users, err = h.service.GetUsersWithFilters(c.Request.Context(), params)
	} else {
		params := db.GetAllUsersParams{
			Limit:  int32(pageSize),
			Offset: int32(offset),
		}
		// users, err = h.service.GetAllUsers(c.Request.Context(), params)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// ... otros handlers (DeleteUser)
