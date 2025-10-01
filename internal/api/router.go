package api

import (
	"github.com/MartinPatricio/GoGinAPISimple/internal/api/handlers"
	"github.com/MartinPatricio/GoGinAPISimple/internal/api/middleware"
	"github.com/MartinPatricio/GoGinAPISimple/internal/config"
	"github.com/MartinPatricio/GoGinAPISimple/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userService *service.UserService, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Middlewares globales
	r.Use(middleware.SecurityHeaders())

	// Instanciar Handlers
	userHandler := handlers.NewUserHandler(userService)

	// Rutas públicas
	r.POST("/register", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	// Agrupar rutas que requieren autenticación
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)   // Esta línea ahora funcionará
			users.DELETE("/:id", userHandler.DeleteUser) // Y esta también
		}
	}

	return r
}
