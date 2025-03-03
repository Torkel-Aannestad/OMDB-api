package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Torkel-Aannestad/OMDB-api/internal/database"
	"github.com/Torkel-Aannestad/OMDB-api/internal/mailer"
	"github.com/Torkel-Aannestad/OMDB-api/internal/vcs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type Config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	authLimiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config Config
	logger *slog.Logger
	models *database.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg Config

	godotenv.Load()
	dsn := os.Getenv("OMDB_API_DB_DSN_DEV")
	mailtrapUsername := os.Getenv("MAILTRAP_USERNAME_DEV")
	mailtrapPassword := os.Getenv("MAILTRAP_PASSWORD_DEV")

	flag.IntVar(&cfg.port, "port", 4000, "port to listen for request")
	flag.StringVar(&cfg.env, "env", "development", "development | staging | production")

	//DB flags
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "dsn for PG instance")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	//Rate limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	//Auth Rate limiter
	flag.Float64Var(&cfg.authLimiter.rps, "auth-limiter-rps", 0.0027778, "Rate limiter maximum requests per second for auth endpoints 10 req / hour / sec")
	flag.IntVar(&cfg.authLimiter.burst, "auth-limiter-burst", 10, "Rate limiter maximum burst for auth endpoints")
	flag.BoolVar(&cfg.authLimiter.enabled, "auth-limiter-enabled", true, "Enable auth rate limiter")

	//Mailer sandbox.smtp.mailtrap.io live.smtp.mailtrap.io
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", mailtrapUsername, "username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", mailtrapPassword, "password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "OMDB <no-reply@torkelaannestad.com>", "sender")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	loggerOptons := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	if cfg.env == "production" {
		loggerOptons.AddSource = true
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptons))

	db, err := database.OpenDB(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// //metrics
	// expvar.NewString("version").Set(version)
	// // Publish the number of active goroutines.
	// expvar.Publish("goroutines", expvar.Func(func() any {
	// 	return runtime.NumGoroutine()
	// }))
	// // Publish the database connection pool statistics.
	// expvar.Publish("database", expvar.Func(func() any {
	// 	return db.Stats()
	// }))
	// // Publish the current Unix timestamp.
	// expvar.Publish("timestamp", expvar.Func(func() any {
	// 	return time.Now().Unix()
	// }))

	app := &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
