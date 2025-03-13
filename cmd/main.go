package main

import (
	"fmt"
	"log"
	"net/http"
	"test_data_flow/configs"
	"test_data_flow/internal/asset"
	"test_data_flow/internal/auth"
	"test_data_flow/internal/session"
	"test_data_flow/internal/user"
	"test_data_flow/pkg/middleware"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func App(cfg *configs.Config, db *sqlx.DB) http.Handler {
	router := http.NewServeMux()

	sessionRepository := session.NewSessionRepository(db)
	userRepository := user.NewUserRepository(db)
	assetRepository := asset.NewAssetRepository(db)

	sessionService := session.NewSessionService(sessionRepository)
	authService := auth.NewAuthService(userRepository, sessionService)
	assetService := asset.NewAssetService(assetRepository)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      cfg,
		AuthService: authService,
	})
	asset.NewAssetHandler(router, asset.AssetHandlerDeps{
		AssetService: assetService,
		Config:       cfg,
	})

	middlewareStack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return middlewareStack(router)
}

func main() {
	cfg := configs.LoadConfig()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.DB_USER,
		cfg.DB.DB_PASSWORD,
		cfg.DB.DB_HOST,
		cfg.DB.DB_PORT,
		cfg.DB.DB_NAME,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	app := App(cfg, db)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.SERVER_PORT),
		Handler: app,
	}

	log.Printf("Server is listening on port: %s\n", cfg.SERVER_PORT)
	server.ListenAndServe()
}
