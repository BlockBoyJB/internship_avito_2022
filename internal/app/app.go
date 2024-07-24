package app

import (
	"avito_intership/config"
	v1 "avito_intership/internal/api/v1"
	"avito_intership/internal/repo"
	"avito_intership/internal/service"
	"avito_intership/pkg/broker"
	"avito_intership/pkg/httpserver"
	"avito_intership/pkg/postgres"
	"avito_intership/pkg/redis"
	"avito_intership/pkg/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

//	@title			Api for account balance management
//	@version		1.0
//	@description	Api for balance management. Include operations, e.g. deposit, withdraw, transfer, reservation, etc

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT token

func Run() {
	// config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// set up json logger
	setLogger(cfg.Log.Level, cfg.Log.Output)

	// postgresql database
	pg, err := postgres.NewPG(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatalf("Initializing postgres error: %s", err)
	}
	defer pg.Close()

	// redis
	rdb := redis.NewRedis(cfg.Redis.Url, redis.SetPassword(cfg.Redis.Password))
	defer rdb.Close()

	// database repositories
	repos := repo.NewRepositories(pg, rdb)

	// kafka producer
	producer, err := broker.NewProducer(cfg.Kafka.Url)
	if err != nil {
		log.Fatalf("Initializing kafka error: %s", err)
	}
	defer producer.Close()

	d := &service.ServicesDependencies{
		Repos:      repos,
		Producer:   producer,
		PrivateKey: cfg.JWT.PrivateKey,
		PublicKey:  cfg.JWT.PublicKey,
	}
	services := service.NewServices(d)

	// validator for incoming messages
	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf("Initializing handler validator error: %s", err)
	}

	// handler
	handler := echo.New()
	handler.Validator = v
	v1.LoggingMiddleware(handler, cfg.Log.Output)
	v1.NewRouter(handler, services)

	// http server
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))

	log.Infof("App started! Listening port %s", cfg.HTTP.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run, signal " + s.String())

	case err = <-httpServer.Notify():
		log.Errorf("/app/run http server notify error: %s", err)
	}
	// graceful shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("/app/run http server shutdown error: %s", err)
	}

	log.Infof("App shutdown with exit code 0")
}

// loading environment params from .env
func init() {
	if _, ok := os.LookupEnv("HTTP_PORT"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("load env file error: %s", err)
		}
	}
}
