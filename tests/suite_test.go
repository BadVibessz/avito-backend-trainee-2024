package tests

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	userbannerhandler "avito-backend-trainee-2024/internal/handler/banner/user"
	midlewares "avito-backend-trainee-2024/internal/handler/middleware"
	bannerrepo "avito-backend-trainee-2024/internal/repository/postgres/banner"
	featurerepo "avito-backend-trainee-2024/internal/repository/postgres/feature"
	tagrepo "avito-backend-trainee-2024/internal/repository/postgres/tag"
	userrepo "avito-backend-trainee-2024/internal/repository/postgres/user"
	bannerservice "avito-backend-trainee-2024/internal/service/banner"
	"avito-backend-trainee-2024/pkg/hasher"
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	authservice "avito-backend-trainee-2024/internal/service/auth"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type BannerService interface {
	GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error)
	CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error)
}

type BannerRepo interface {
	GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error)
	GetBannerByID(ctx context.Context, id int) (*entity.Banner, error)
	GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error)
	CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, id int, updateModel entity.Banner) error
	DeleteBanner(ctx context.Context, id int) (*entity.Banner, error)
}

type BannerHandler interface {
	GetBannerByFeatureAndTags(rw http.ResponseWriter, req *http.Request)
	Routes() *chi.Mux
}

var (
	dbConnectionStr string
	jwtSecret       string
)

func init() {
	dbConnectionStr = "postgresql://postgres:postgres@localhost:5433/test?sslmode=disable"
	jwtSecret = os.Getenv("TEST_JWT_SECRET")
}

type Suite struct {
	suite.Suite

	db *sqlx.DB

	bannerRepo    BannerRepo
	bannerService BannerService
	bannerHandler BannerHandler
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) setupDB() {
	conn, err := sql.Open("pgx", dbConnectionStr)
	if err != nil {
		s.FailNowf("cannot open database connection with connection string: %v, err: %v", dbConnectionStr, err)
	}

	s.db = sqlx.NewDb(conn, "postgres")
}

func (s *Suite) setupRepos() {
	s.bannerRepo = bannerrepo.New(s.db)
}

func (s *Suite) setupServices() {
	featureRepo := featurerepo.New(s.db)
	tagRepo := tagrepo.New(s.db)

	s.bannerService = bannerservice.New(s.bannerRepo, featureRepo, tagRepo)
}

func (s *Suite) setupHandlers() {
	logger := logrus.New()
	valid := validator.New(validator.WithRequiredStructEnabled())
	cache := gocache.New(5*time.Minute, 10*time.Minute)

	authMiddleware := midlewares.JWTAuthentication("token", jwtSecret, logger)
	cacheMiddleware := midlewares.InMemUserBannerCache(cache, logger)

	s.bannerHandler = userbannerhandler.New(s.bannerService, logger, valid, authMiddleware, cacheMiddleware)
}

func (s *Suite) SetupSuite() {
	s.setupDB()
	s.setupRepos()
	s.setupServices()
	s.setupHandlers()
	s.loadFixturesInDB()
}

func (s *Suite) TearDownSuite() {
	_ = s.db.Close() // close db connection
}

func (s *Suite) loadFixturesInDB() {
	hash := hasher.New() // TODO: mock!!
	ctx := context.Background()

	userRepo := userrepo.New(s.db)
	authService := authservice.New(userRepo, hash)

	for _, user := range users {
		_, _ = authService.RegisterUser(ctx, user) // todo:
	}

	for _, banner := range banners {
		_, _ = s.bannerService.CreateBanner(ctx, banner)
	}
}
