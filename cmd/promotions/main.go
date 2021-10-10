package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/torwig/promotions/internal/http/server"
	"github.com/torwig/promotions/internal/promotions/app"
)

func main() {
	ctx := context.Background()

	httpPort := os.Getenv("HTTP_PORT")

	application := app.NewApplication(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	s := server.NewServer(application)
	s.RunHTTPServer(httpPort)

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	s.Shutdown(ctx)
}
