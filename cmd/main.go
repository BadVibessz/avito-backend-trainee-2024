package main

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/repository/postgres/banner"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math"
)

func main() {
	logger := logrus.New()
	ctx := context.Background()

	//port := 5000
	//
	//r := chi.NewRouter()
	//
	//server := http.Server{
	//	Addr:    fmt.Sprintf(":%v", port),
	//	Handler: r,
	//}
	//
	//// add swagger middleware
	//r.Get("/swagger/*", httpSwagger.Handler(
	//	httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port)), // The url pointing to API definition
	//))
	//
	//logger.Infof("server started at port %v", server.Addr)
	//
	//go func() {
	//	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		logger.WithError(err).Fatalf("server can't listen requests")
	//	}
	//}()
	//
	//logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", port)

	// todo: config
	connStr := "postgresql://postgres:postgres@localhost:5432/avito-trainee?sslmode=disable"

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", connStr, err)
	}

	db := sqlx.NewDb(conn, "postgres")

	bannerRepo := banner.New(db)

	banner := entity.Banner{
		Name:      "test_banner",
		TagIDs:    []int{5, 6},
		FeatureID: 5,
		Content: entity.Content{
			Title: "test_content",
			Text:  "some_text",
			Url:   "google.com",
		},
		IsActive: true,
	}

	createdBanner, err := bannerRepo.CreateBanner(ctx, banner)
	if err != nil {
		if err != nil {
			logger.Fatalf("error occurred creating banner: %v", err)
		}
	}

	logger.Infof("created banner: %+v", createdBanner)

	banners, err := bannerRepo.GetAllBanners(ctx, 0, math.MaxInt64)
	if err != nil {
		logger.Fatalf("error occurred fetching banners from db: %v", err)
	}

	for _, banner := range banners {
		logger.Infof("banner: %+v", banner)
	}

	updateModel := entity.Banner{
		TagIDs:    []int{6},
		FeatureID: 6,
		Content: entity.Content{
			Title: "new_title",
			Text:  "new_text",
			Url:   "new_url",
		},
		IsActive: true,
	}

	err = bannerRepo.UpdateBanner(ctx, createdBanner.ID, updateModel)
	if err != nil {
		logger.Fatalf("error occurred updating banner: %v", err)
	}

	got, err := bannerRepo.GetBannerByFeatureAndTags(ctx, updateModel.FeatureID, updateModel.TagIDs)
	if err != nil {
		logger.Fatalf("error occurred fetching banner from db: %v", err)
	}

	logger.Infof("banner found: %+v", got)

	deleted, err := bannerRepo.DeleteBanner(ctx, createdBanner.ID)
	if err != nil {
		logger.Fatalf("error occurred deleting banners from db: %v", err)
	}

	logger.Infof("deleted: %+v", deleted)

}
