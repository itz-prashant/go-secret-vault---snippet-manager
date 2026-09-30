package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itz-prashant/secret-vault-api/internal/config"
	"github.com/itz-prashant/secret-vault-api/internal/db"
	"github.com/itz-prashant/secret-vault-api/internal/utils/response"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	sqlDb, err := db.New(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize database %v", err)
	}

	defer sqlDb.Db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJson(w, http.StatusOK, map[string]string{
			"service": "secret-vault-api",
			"status": "ok",
		})
	})

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("[INFO] server listening on *%s", cfg.Address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listening failed %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server gracefully....")

	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cencel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown %v", err)
	}
}
