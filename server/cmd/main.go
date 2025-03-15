package main

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		if rand.Intn(100) < 3 { // 3% chance to simulate an error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if rand.Intn(100) < 3 { // 3% chance to simulate an error
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr: ":8080",
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Could not start the server", "error", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	slog.Info("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Could not shutdown the server", "error", err)
	}
	slog.Info("Server stopped")
	os.Exit(0)

}
