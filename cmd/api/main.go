// Filename: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/amilcar-vasquez/qod/internal/data"
	_ "github.com/lib/pq"
)

const appVersion = "2.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
}

type applicationDependencies struct {
	config     serverConfig
	logger     *slog.Logger
	quoteModel *data.QuoteModel
	userModel  *data.UserModel
}

func main() {
	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development",
		"Environment(development|staging|production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://qod:password@localhost/qod?sslmode=disable",
		"PostgreSQL DSN")
	flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&settings.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)",
		func(val string) error {
			settings.cors.trustedOrigins = strings.Fields(val)
			return nil
		})
	flag.Parse()
	//print out flags values
	fmt.Printf(`Starting server with config:
	port: %d
	environment: %s
	db-dsn: %s
	limiter-rps: %.2f
	limiter-burst: %d
	limiter-enabled: %t
	cors-trusted-origins: %v
	`, settings.port, settings.environment, settings.db.dsn, settings.limiter.rps, settings.limiter.burst, settings.limiter.enabled, settings.cors.trustedOrigins)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// the call to openDB() sets up our connection pool
	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")

	appInstance := &applicationDependencies{
		config:     settings,
		logger:     logger,
		quoteModel: &data.QuoteModel{DB: db},
		userModel:  &data.UserModel{DB: db},
	}

	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(settings serverConfig) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second)
	defer cancel()
	// let's test if the connection pool was created
	// we trying pinging it with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}
