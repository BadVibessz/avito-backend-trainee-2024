package main

import (
	"avito-backend-trainee-2024/internal/config"
	"avito-backend-trainee-2024/pkg/hasher"
	gocache "github.com/patrickmn/go-cache"
	"time"

	router "avito-backend-trainee-2024/pkg/route"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	chimiddlewares "github.com/go-chi/chi/v5/middleware"

	bannerrepo "avito-backend-trainee-2024/internal/repository/postgres/banner"
	featurerepo "avito-backend-trainee-2024/internal/repository/postgres/feature"
	tagrepo "avito-backend-trainee-2024/internal/repository/postgres/tag"
	userrepo "avito-backend-trainee-2024/internal/repository/postgres/user"

	authservice "avito-backend-trainee-2024/internal/service/auth"
	bannerservice "avito-backend-trainee-2024/internal/service/banner"

	midlewares "avito-backend-trainee-2024/internal/handler/middleware"

	adminbannerhandler "avito-backend-trainee-2024/internal/handler/banner/admin"
	userbannerhandler "avito-backend-trainee-2024/internal/handler/banner/user"

	authhandler "avito-backend-trainee-2024/internal/handler/auth"
	httpswagger "github.com/swaggo/http-swagger"

	_ "avito-backend-trainee-2024/docs"
)

const (
	configPath = "./config"
)

func initConfig() (*config.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// env variables
	if err := godotenv.Load(configPath + "/.env"); err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("avito_trainee")
	viper.AutomaticEnv()

	conf.Jwt.Secret = viper.GetString("JWT_SECRET")
	if conf.Jwt.Secret == "" {
		return nil, errors.New("JWT_SECRET env variable not set")
	}

	return &conf, nil
}

func main() {
	logger := logrus.New()
	valid := validator.New(validator.WithRequiredStructEnabled())
	cache := gocache.New(5*time.Minute, 10*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	conf, err := initConfig()
	if err != nil {
		logger.Fatalf("error occurred initializing config: %v", err)
	}

	conn, err := sql.Open("pgx", conf.Postgres.ConnectionURL())
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", conf.Postgres.ConnectionURL(), err)
	}

	db := sqlx.NewDb(conn, "postgres")

	if err = db.Ping(); err != nil {
		logger.Fatalf("can't ping database: %v", err)
	}

	userRepo := userrepo.New(db)
	bannerRepo := bannerrepo.New(db)
	featureRepo := featurerepo.New(db)
	tagRepo := tagrepo.New(db)

	bannerService := bannerservice.New(bannerRepo, featureRepo, tagRepo)
	authService := authservice.New(userRepo, hasher.New())

	authMiddleware := midlewares.JWTAuthentication("token", conf.Jwt.Secret, logger)
	adminAuthMiddleware := midlewares.AdminAuthorization(logger)
	cacheMiddleware := midlewares.InMemUserBannerCache(cache, logger)

	authHandler := authhandler.New(authService, conf.Jwt, logger, valid)
	userBannerHandler := userbannerhandler.New(bannerService, logger, valid, authMiddleware, cacheMiddleware)
	adminBannerHandler := adminbannerhandler.New(bannerService, logger, valid, authMiddleware, adminAuthMiddleware)

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = userBannerHandler.Routes()
	routers["/banner"] = adminBannerHandler.Routes()
	routers["/auth"] = authHandler.Routes()

	middlewares := []router.Middleware{
		chimiddlewares.Recoverer,
		chimiddlewares.Logger,
	}

	r := router.MakeRoutes("/avito-trainee/api/v1", routers, middlewares...)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Server.Port),
		Handler: r,
	}

	// add swagger middleware
	r.Get("/swagger/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", conf.Server.Port)), // The url pointing to API definition
	))

	logger.Infof("server started at port %v", server.Addr)

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.WithError(listenErr).Fatalf("server can't listen requests")
		}
	}()

	logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", conf.Server.Port)

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(interrupt, syscall.SIGINT)

	go func() {
		<-interrupt

		logger.Info("interrupt signal caught: shutting server down")

		if shutdownErr := server.Shutdown(ctx); err != nil {
			logger.WithError(shutdownErr).Fatalf("can't close server listening on '%s'", server.Addr)
		}

		cancel()
	}()

	<-ctx.Done()
}
