// Filename: cmd/api/main.go

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"flag"
)

// configuration holds application configuration settings.
type configuration struct {
	port int    // Port for the API server to listen on
	env  string // Environment (development|staging|production)
}

// application holds the dependencies for the HTTP server.
type application struct {
	config configuration // App configuration
	logger *slog.Logger  // Structured logger
}

func main() {
	// Initialize configuration from command line flags.
	cfg := loadConfig()
	// Initialize logger based on environment.
	logger := setupLogger(cfg.env)

	// Create application struct with config and logger.
	app := application{
		config: cfg,
		logger: logger,
	}

	// Start the HTTP server.
	err := app.serve()
	if err != nil {
		// Log error and exit if server fails to start.
		logger.Error(err.Error())
		os.Exit(1)
	}
} // end of main

// loadConfig reads configuration from command line flags.
func loadConfig() configuration {
	var cfg configuration

	// Define command line flags for port and environment.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	return cfg
}

// setupLogger initializes a structured logger for the application.
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	// Create a new text-based logger writing to stdout.
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

// serve starts the HTTP server with configured settings.
func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port), // Server address
		Handler:      app.routes(),                        // HTTP handler routes
		IdleTimeout:  time.Minute,                         // Idle connection timeout
		ReadTimeout:  5 * time.Second,                     // Read timeout
		WriteTimeout: 10 * time.Second,                    // Write timeout
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError), // Error logger
	}

	// Log server startup information.
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	return srv.ListenAndServe() // Start listening for HTTP requests.
}
