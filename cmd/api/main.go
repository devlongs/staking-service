package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/devlongs/staking-service/internal/database"
	"github.com/devlongs/staking-service/internal/service"
	"github.com/devlongs/staking-service/internal/store"
)

type config struct {
	addr   string
	dbPath string
}

type application struct {
	config config
	store  store.Storage
	svc    service.StakeService
	logger *log.Logger
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// API v1 endpoints
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", v1.HealthHandler())
		r.Post("/stake", v1.StakeHandler(app.svc))
		r.Get("/rewards/{wallet_address}", v1.RewardsHandler(app.svc))
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine.
	serverErrors := make(chan error, 1)
	go func() {
		app.logger.Printf("Server starting on %s", app.config.addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err
	case sig := <-quit:
		app.logger.Printf("Shutting down server... (signal: %v)", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			app.logger.Printf("Could not gracefully shutdown the server: %v", err)
			return err
		}
		app.logger.Printf("Server exiting")
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, now using environment variables")
	}

	logger := log.New(os.Stdout, "stakeway: ", log.Ldate|log.Ltime)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "stakeway.db"
	}

	cfg := config{
		addr:   addr,
		dbPath: dbPath,
	}

	// Initialize SQLite database.
	db, err := database.NewSQLiteDB(cfg.dbPath)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize storage and service layers
	st := store.NewStorage(db)
	svc := service.NewStakeService(st)

	app := application{
		config: cfg,
		store:  st,
		svc:    svc,
		logger: logger,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		logger.Fatal(err)
	}
}
