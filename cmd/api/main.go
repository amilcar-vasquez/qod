// Filename: cmd/api/main.go

package main

import (
	"fmt"
	"log/slog"
)

type configuration struct {
	port string
	env string
}

type application struct {
	config configuration
	logger *slog.Logger
}

func main {
	//initialize configuration
	cfg := loadConfig()
	//initalize logger
	logger := setupLogger()
	app := application{
		config: cfg,
		logger: logger,
	}
	app.serve()
if err != nil {
	logger.Error(err.Error())
	os.Exit (1)
}
}  //end of main

//loadConfig reads configuration from command line flags
func loadConfig() configuration {
    var cfg configuration
	
flag.IntVar(&cfg.port, "port", 4000, "API server port")
flag.StringVar(&cfg.env,"env","development","Environment(development
                   |staging|production)")
flag.Parse()
	
return cfg
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		
		
	return logger
}

// serve starts the HTTP server (server.go)
func (app *application) serve() error {
       srv := &http.Server {
       Addr:         fmt.Sprintf(":%d", app.config.port),
       Handler:      app.routes(),
       IdleTimeout:  time.Minute,
       ReadTimeout:  5 * time.Second,
       WriteTimeout: 10 * time.Second,
       ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
    }

 app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
 return srv.ListenAndServe()
}
