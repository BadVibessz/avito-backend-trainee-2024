package admin

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/handler/mapper"
	"avito-backend-trainee-2024/internal/handler/request"
	sliceutils "avito-backend-trainee-2024/pkg/utils/slice"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/go-playground/validator/v10"

	handlerinternalutils "avito-backend-trainee-2024/internal/pkg/utils/handler"
	handlerutils "avito-backend-trainee-2024/pkg/utils/handler"
)

type Service interface {
	GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error)
	GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error)
	CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, id int, updateModel entity.Banner) error
	DeleteBanner(ctx context.Context, id int) (*entity.Banner, error)
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

		r.Get("/", h.GetAllBanners)
		r.Post("/", h.CreateBanner)
		r.Patch("/{id}", h.UpdateBanner)
		r.Delete("/{id}", h.DeleteBanner)
	})

	return router
}

// GetAllBanners godoc
//
//	@Summary		Get all banners
//	@Description	Get all banners sorting by featureID
//	@Security		JWT
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param token 	header string true "admin auth token"
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Success		200		{object}	[]response.GetAdminBannerResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		403		{string}	Forbidden
//	@Failure		400		{string}	invalid		request
//	@Failure		500		{string}	internal	error
//	@Router			/avito-trainee/api/v1/banner [get]
func (h *Handler) GetAllBanners(rw http.ResponseWriter, req *http.Request) {
	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, DefaultOffset, DefaultLimit)

	if err := paginationOpts.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("invalid pagination options provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	banners, err := h.Service.GetAllBanners(req.Context(), paginationOpts.Offset, paginationOpts.Limit)
	if err != nil {
		msg := fmt.Sprintf("error occurred fetching banners: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	render.JSON(rw, req, sliceutils.Map(banners, mapper.MapBannerToAdminBannerResponse))
	rw.WriteHeader(http.StatusOK)
}

// CreateBanner godoc
//
//	@Summary		Create new banner
//	@Description	Create new banner
//	@Security		JWT
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param token 	header string true "admin auth token"
//	@Param			input	body		request.CreateBannerRequest	true	"create banner schema"
//	@Success		200		{object}	response.CreateBannerResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		403		{string}	Forbidden
//	@Failure		400		{string}	invalid		request
//	@Failure		500		{string}	internal	error
//	@Router			/avito-trainee/api/v1/banner [post]
func (h *Handler) CreateBanner(rw http.ResponseWriter, req *http.Request) {
	var bannerReq request.CreateBannerRequest

	if err := render.DecodeJSON(req.Body, &bannerReq); err != nil {
		msg := fmt.Sprintf("error occurred decoding request body to CreateBannerRequest srtuct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err := bannerReq.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("error occurred validating CreateBannerRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	created, err := h.Service.CreateBanner(req.Context(), mapper.MapCreateBannerRequestToEntity(&bannerReq))
	if err != nil {
		msg := fmt.Sprintf("error occurred creating banner: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	render.JSON(rw, req, mapper.MapBannerToCreateBannerResponse(created))
	rw.WriteHeader(http.StatusCreated)
}

// UpdateBanner godoc
//
//	@Summary		Update existing banner
//	@Description	Update existing banner
//	@Security		JWT
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param token 	header string true "admin auth token"
//	@Param			input	body	request.UpdateBannerRequest	true	"update banner schema"
//	@Param			id		path	int							true	"id of the updating banner"
//	@Success		200
//	@Failure		401	{string}	Unauthorized
//	@Failure		403	{string}	Forbidden
//	@Failure		400	{string}	invalid		request
//	@Failure		500	{string}	internal	error
//	@Router			/avito-trainee/api/v1/banner/{id} [patch]
func (h *Handler) UpdateBanner(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		msg := fmt.Sprintf("inavlid url param for id provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	var updateReq request.UpdateBannerRequest

	if err = render.DecodeJSON(req.Body, &updateReq); err != nil {
		msg := fmt.Sprintf("error occurred decoding request body to UpdateBannerRequest srtuct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err = updateReq.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("error occurred validating UpdateBannerRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err = h.Service.UpdateBanner(req.Context(), id, mapper.MapUpdateBannerRequestToEntity(&updateReq)); err != nil {
		msg := fmt.Sprintf("error occurred updating banner: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// DeleteBanner godoc
//
//	@Summary		Delete banner
//	@Description	Delete banner
//	@Security		JWT
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param token 	header string true "admin auth token"
//	@Param			id	path	int	true	"id of the banner"
//	@Success		200
//	@Failure		401	{string}	Unauthorized
//	@Failure		403	{string}	Forbidden
//	@Failure		400	{string}	invalid		request
//	@Failure		500	{string}	internal	error
//	@Router			/avito-trainee/api/v1/banner/{id} [delete]
func (h *Handler) DeleteBanner(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		msg := fmt.Sprintf("inavlid url param for id provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	_, err = h.Service.DeleteBanner(req.Context(), id)
	if err != nil {
		msg := fmt.Sprintf("error occurred deleting banner: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
