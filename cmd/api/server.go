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
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Conectar a la base de datos usando PGX
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	// Usamos el adaptador stdlib para obtener un *sql.DB compatible con la interfaz de pgx
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	log.Println("Database connection successful!")

	// 3. Inyectar dependencias (Repositorio -> Servicio)
	userRepo := repository.NewSQLUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)

	// 4. Configurar el router de Gin
	router := api.SetupRouter(userService, cfg)

	// 5. Iniciar el servidor
	log.Printf("Starting server on port %s", cfg.ApiPort)
	if err := router.Run(":" + cfg.ApiPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
