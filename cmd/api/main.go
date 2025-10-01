package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MartinPatricio/GoGinAPISimple/internal/api"
	"github.com/MartinPatricio/GoGinAPISimple/internal/config"
	"github.com/MartinPatricio/GoGinAPISimple/internal/repository"
	"github.com/MartinPatricio/GoGinAPISimple/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListarRutas(router *gin.Engine) {
	fmt.Println("===================================================")
	fmt.Println("    RUTAS DE API DISPONIBLES (GIN FRAMEWORK)       ")
	fmt.Println("===================================================")

	// router.Routes() devuelve un slice de gin.RouteInfo
	for _, route := range router.Routes() {
		// Imprime el método HTTP y la ruta (Path)
		fmt.Printf(" [GIN] Método: %-7s | URL: %s\n", route.Method, route.Path)
	}
	fmt.Println("===================================================")
}
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)
	fmt.Println(connStr)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer dbpool.Close()

	if err = dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	log.Println("Database connection successful!")

	//userRepo := repository.NewSQLUserRepository(dbpool)
	sqlUserRepo := repository.NewSQLUserRepository(dbpool)
	userRepo := repository.NewLoggedUserRepository(sqlUserRepo)

	userService := service.NewUserService(userRepo, cfg)
	router := api.SetupRouter(userService, cfg)

	log.Printf("Starting server on port %s", cfg.ApiPort)
	if err := router.Run(":" + cfg.ApiPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
