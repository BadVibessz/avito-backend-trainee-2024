package auth

import (
	"avito-backend-trainee-2024/internal/config"
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/handler/mapper"
	"avito-backend-trainee-2024/internal/handler/request"
	"avito-backend-trainee-2024/internal/handler/response"
	handlerutils "avito-backend-trainee-2024/pkg/utils/handler"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"

	jwtutils "avito-backend-trainee-2024/pkg/utils/jwt"
)

type Service interface {
	RegisterUser(ctx context.Context, user entity.User) (*entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	Service     Service
	Middlewares []Middleware

	jwtConfig config.Jwt
	logger    *logrus.Logger
	validator *validator.Validate
}

func New(service Service, jwtConfig config.Jwt, logger *logrus.Logger, validator *validator.Validate, middlewares ...Middleware) *Handler {
	return &Handler{
		Service:     service,
		Middlewares: middlewares,
		jwtConfig:   jwtConfig,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)

		r.Post("/user_register", h.RegisterUser)
		r.Post("/admin_register", h.RegisterAdmin)
		r.Post("/login", h.Login)
	})

	return router
}

// RegisterUser godoc
//
//	@Summary		Register new user
//	@Description	Register new user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.RegisterRequest	true	"register user schema"
//	@Success		200		{object}	response.RegisterUserResponse
//	@Failure		400		{string}	invalid		request
//	@Failure		500		{string}	internal	error
//	@Router			/avito-trainee/api/v1/auth/user_register [post]
func (h *Handler) RegisterUser(rw http.ResponseWriter, req *http.Request) {
	var registerReq request.RegisterRequest

	if err := render.DecodeJSON(req.Body, &registerReq); err != nil {
		msg := fmt.Sprintf("error occurred decoding request body to RegisterRequest srtuct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err := registerReq.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("error occurred validating RegisterRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	user, err := h.Service.RegisterUser(req.Context(), mapper.MapRegisterRequestToEntity(&registerReq, false))
	if err != nil {
		msg := fmt.Sprintf("error occurred registering user: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	render.JSON(rw, req, mapper.MapUserToRegisterUserResponse(user))
	rw.WriteHeader(http.StatusCreated)
}

// RegisterAdmin godoc
//
//	@Summary		Register new user
//	@Description	Register new user
//	@Security		JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.RegisterRequest	true	"register user schema"
//	@Success		200		{object}	response.RegisterUserResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		403		{string}	Forbidden
//	@Failure		400		{string}	invalid		request
//	@Failure		500		{string}	internal	error
//	@Router			/avito-trainee/api/v1/auth/admin_register [post]
func (h *Handler) RegisterAdmin(rw http.ResponseWriter, req *http.Request) {
	// todo: admin auth (only admin can register new admin)

	var registerReq request.RegisterRequest

	if err := render.DecodeJSON(req.Body, &registerReq); err != nil {
		msg := fmt.Sprintf("error occurred decoding request body to RegisterRequest srtuct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err := registerReq.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("error occurred validating RegisterRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	user, err := h.Service.RegisterUser(req.Context(), mapper.MapRegisterRequestToEntity(&registerReq, true))
	if err != nil {
		msg := fmt.Sprintf("error occurred registering admin: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	render.JSON(rw, req, mapper.MapUserToRegisterAdminResponse(user))
	rw.WriteHeader(http.StatusCreated)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	login user via JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		request.LoginRequest	true	"login info"
//	@Success		200		{object}	response.LoginResponse
//	@Failure		400		{string}	invalid		login	data	provided
//	@Failure		500		{string}	internal	error
//	@Router			/avito-trainee/api/v1/auth/login [post]
func (h *Handler) Login(rw http.ResponseWriter, req *http.Request) {
	var loginReq request.LoginRequest

	if err := render.DecodeJSON(req.Body, &loginReq); err != nil {
		msg := fmt.Sprintf("error occurred decoding request body to LoginRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	if err := loginReq.Validate(h.validator); err != nil {
		msg := fmt.Sprintf("error occurred validating LoginRequest struct: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	user, err := h.Service.Login(req.Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		msg := fmt.Sprintf("error occurred while user login: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	// construct jwt token
	payload := map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, h.jwtConfig.Secret)
	if err != nil {
		msg := fmt.Sprintf("error occurred signing jwt token: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, msg, msg)
		return
	}

	resp := response.LoginResponse{
		Token: token,
	}

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
