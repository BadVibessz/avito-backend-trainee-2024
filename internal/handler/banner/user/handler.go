package user

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/handler/mapper"
	"avito-backend-trainee-2024/internal/handler/middleware"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/go-playground/validator/v10"

	handlerutils "avito-backend-trainee-2024/pkg/utils/handler"
	urlutils "avito-backend-trainee-2024/pkg/utils/url"
)

type Service interface {
	GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error)
}

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	Service     Service
	Middlewares []Middleware

	logger    *logrus.Logger
	validator *validator.Validate
}

func New(service Service, logger *logrus.Logger, validator *validator.Validate, middlewares ...Middleware) *Handler {
	return &Handler{
		Service:     service,
		Middlewares: middlewares,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)

		r.Get("/", h.GetBannerByFeatureAndTags)
	})

	return router
}

// GetBannerByFeatureAndTags godoc
//
//	@Summary		Get banner with feature and tags
//	@Description	Get banner with feature and tags
//	@Security		JWT
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param token 	header string true "user auth token"
//	@Param			feature_id	query		string	true	"id of the feature"
//	@Param			tag_ids		query		[]int	true	"ids of the tags"
//	@Param			use_last_revision		query		bool	true	"use last revision?"
//	@Success		200			{object}	response.GetUserBannerResponse
//	@Failure		401			{string}	Unauthorized
//	@Failure		400			{string}	invalid		request
//	@Failure		403			{string}	invalid		request
//	@Failure		500			{string}	internal	error
//	@Router			/avito-trainee/api/v1/user_banner [get]
func (h *Handler) GetBannerByFeatureAndTags(rw http.ResponseWriter, req *http.Request) {
	featureID, err := handlerutils.GetIntParamFromQuery(req, "feature_id")
	if err != nil {
		msg := fmt.Sprintf("error occurred getting 'feature_id' query param: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	tagIDs, err := handlerutils.GetIntArrayParamFromQuery(req, "tag_ids")
	if err != nil {
		msg := fmt.Sprintf("error occurred getting 'tag_ids' query param: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	banner, err := h.Service.GetBannerByFeatureAndTags(req.Context(), featureID, tagIDs)
	if err != nil {
		msg := fmt.Sprintf("error occurred fetching banner: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	// return to users only active banners, if user = admin, then return anyway
	if !banner.IsActive && req.Header.Get("is_admin") != "true" {
		msg := "banner is inactive"

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusForbidden, msg, msg)
		return
	}

	resp := mapper.MapBannerToUserBannerResponse(banner)

	if middlewareData, ok := req.Context().
		Value(urlutils.RemoveQueryParamByKey(*req.URL, "use_last_revision").
			RequestURI()).(middleware.MiddlewareData); ok {
		middlewareData["banner"] = resp
		middlewareData["is_active"] = banner.IsActive
	}

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
