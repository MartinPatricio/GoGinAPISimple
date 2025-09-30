// cmd/api/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MartinPatricio/GoGinAPISimple/internal/api"
	"github.com/MartinPatricio/GoGinAPISimple/internal/config"
	"github.com/MartinPatricio/GoGinAPISimple/internal/repository"
	"github.com/MartinPatricio/GoGinAPISimple/internal/service"

	// PASO 1: Importa el driver PGX para que se registre con database/sql.
	// El guion bajo es crucial.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	// PASO 2: Dile a sql.Open que use el driver "pgx" que acabamos de importar.
	// No uses "postgres" aquí.
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("could not open db connection: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	log.Println("Database connection successful!")

	// De aquí en adelante, todo funciona porque el *sql.DB es compatible con la interfaz de SQLC.
	userRepo := repository.NewSQLUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	router := api.SetupRouter(userService, cfg)

	log.Printf("Starting server on port %s", cfg.ApiPort)
	if err := router.Run(":" + cfg.ApiPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
