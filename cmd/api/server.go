// Filename: cmd/api/server.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *applicationDependencies) serve() error {
	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.port),
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}
	// create a channel to keep track of any errors during the shutdown process

	shutdownError := make(chan error)
	// run the server in a goroutine so that it doesn't block the graceful shutdown handling below
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		a.logger.Info("shutting down server", "signal", sig.String())
		// create a context to attempt a graceful 5 second shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// call the server's Shutdown() method which is what will trigger all of our
		shutdownError <- apiServer.Shutdown(ctx)
	}()

	a.logger.Info("starting server", "addr", apiServer.Addr, "env", a.config.environment)
	err := apiServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error("server error", "err", err)
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	a.logger.Info("stopped server", "addr", apiServer.Addr)
	return nil
}
