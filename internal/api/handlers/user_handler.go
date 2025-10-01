package handlers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
	"github.com/MartinPatricio/GoGinAPISimple/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
		return
	}

	// 2. Imprimir el contenido del JSON en la terminal donde corre el servidor.
	// Usamos log.Printf para un formato más limpio.
	log.Printf("--- Raw JSON Received on /register ---\n%s\n--------------------------------------", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var req db.CreateUserParams
	//if err := c.ShouldBindJSON(&req); err != nil {
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
			NameFilter:  pgtype.Text{String: "%" + nameFilter + "%", Valid: nameFilter != ""},
			EmailFilter: pgtype.Text{String: "%" + emailFilter + "%", Valid: emailFilter != ""},
		}
		users, err = h.service.GetUsersWithFilters(c.Request.Context(), params)
	} else {
		params := db.GetAllUsersParams{
			Limit:  int32(pageSize),
			Offset: int32(offset),
		}
		users, err = h.service.GetAllUsers(c.Request.Context(), params)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// internal/api/handlers/user_handler.go

// ... (resto de tu código de handler) ...

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), int32(id))
	if err != nil {
		// Maneja el caso en que el usuario no se encuentra
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
