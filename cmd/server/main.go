package main

import (
	bannerrepo "avito-backend-trainee-2024/internal/repository/postgres/banner"
	featurerepo "avito-backend-trainee-2024/internal/repository/postgres/feature"
	tagrepo "avito-backend-trainee-2024/internal/repository/postgres/tag"
	router "avito-backend-trainee-2024/pkg/route"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	bannerservice "avito-backend-trainee-2024/internal/service/banner"

	bannerhandler "avito-backend-trainee-2024/internal/handler/banner"

	httpswagger "github.com/swaggo/http-swagger"

	_ "avito-backend-trainee-2024/docs"
)

func main() {
	logger := logrus.New()
	valid := validator.New(validator.WithRequiredStructEnabled())

	ctx, cancel := context.WithCancel(context.Background())

	// todo: config
	port := 5000
	connStr := "postgresql://postgres:postgres@localhost:5432/avito-trainee?sslmode=disable"

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", connStr, err)
	}

	db := sqlx.NewDb(conn, "postgres")

	bannerRepo := bannerrepo.New(db)

	//banner := entity.Banner{
	//	Name:      "test_banner",
	//	TagIDs:    []int{5, 6},
	//	FeatureID: 5,
	//	Content: entity.Content{
	//		Title: "test_content",
	//		Text:  "some_text",
	//		Url:   "google.com",
	//	},
	//	IsActive: true,
	//}
	//
	//_, err = bannerRepo.CreateBanner(ctx, banner)
	//if err != nil {
	//	if err != nil {
	//		logger.Fatalf("error occurred creating banner: %v", err)
	//	}
	//}

	//logger.Infof("created banner: %+v", createdBanner)
	//
	//banners, err := bannerRepo.GetAllBanners(ctx, 0, math.MaxInt64)
	//if err != nil {
	//	logger.Fatalf("error occurred fetching banners from db: %v", err)
	//}
	//
	//for _, banner = range banners {
	//	logger.Infof("banner: %+v", banner)
	//}
	//
	//updateModel := entity.Banner{
	//	TagIDs:    []int{6},
	//	FeatureID: 6,
	//	Content: entity.Content{
	//		Title: "new_title",
	//		Text:  "new_text",
	//		Url:   "new_url",
	//	},
	//	IsActive: true,
	//}
	//
	//err = bannerRepo.UpdateBanner(ctx, createdBanner.ID, updateModel)
	//if err != nil {
	//	logger.Fatalf("error occurred updating banner: %v", err)
	//}
	//
	//logger.Infof("updated banner")
	//
	//got, err := bannerRepo.GetBannerByFeatureAndTags(ctx, updateModel.FeatureID, updateModel.TagIDs)
	//if err != nil {
	//	logger.Fatalf("error occurred fetching banner from db: %v", err)
	//}
	//
	//logger.Infof("banner found: %+v", got)
	//
	//deleted, err := bannerRepo.DeleteBanner(ctx, createdBanner.ID)
	//if err != nil {
	//	logger.Fatalf("error occurred deleting banners from db: %v", err)
	//}
	//
	//logger.Infof("deleted: %+v", deleted)

	featureRepo := featurerepo.New(db)
	tagRepo := tagrepo.New(db)

	bannerService := bannerservice.New(bannerRepo, featureRepo, tagRepo)

	bannerHandler := bannerhandler.New(bannerService, logger, valid)

	routers := make(map[string]chi.Router)

	routers["/banner"] = bannerHandler.Routes()

	r := router.MakeRoutes("/avito-trainee/api/v1", routers)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	// add swagger middleware
	r.Get("/swagger/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port)), // The url pointing to API definition
	))

	logger.Infof("server started at port %v", server.Addr)

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.WithError(listenErr).Fatalf("server can't listen requests")
		}
	}()

	logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", port)

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
