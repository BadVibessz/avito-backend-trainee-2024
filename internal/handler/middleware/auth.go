package middleware

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	handlerutils "avito-backend-trainee-2024/pkg/utils/handler"
	jwtutils "avito-backend-trainee-2024/pkg/utils/jwt"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	maputils "avito-backend-trainee-2024/pkg/utils/map"
)

type Handler = func(http.Handler) http.Handler

type AuthService interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

func JWTAuthentication(headerName, secret string, logger *logrus.Logger) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get(headerName)
			if authHeader == "" {
				msg := fmt.Sprintf("%v header is empty", headerName)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			token := authHeader

			payload, err := jwtutils.ValidateToken(token, secret)
			if err != nil {
				msg := fmt.Sprintf("error occurred validating token: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			id, err := maputils.GetIntFromAnyMap(payload, "id")
			if err != nil {
				msg := fmt.Sprintf("invalid payload: not contains id: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			username, err := maputils.GetStringFromAnyMap(payload, "username")
			if err != nil {
				msg := fmt.Sprintf("invalid payload: not contains username: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			isAdmin, err := maputils.GetBoolFromAnyMap(payload, "is_admin")
			if err != nil {
				msg := fmt.Sprintf("invalid payload: not contains is_admin: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusUnauthorized, msg, msg)
				return
			}

			req.Header.Set("id", strconv.Itoa(id))
			req.Header.Set("username", username)
			req.Header.Set("is_admin", fmt.Sprintf("%v", isAdmin))

			next.ServeHTTP(rw, req)
		})
	}
}

func AdminAuthorization(logger *logrus.Logger) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("is_admin") != "true" {
				msg := "only admin allowed to call this method"

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusForbidden, msg, msg)
				return
			}

			next.ServeHTTP(rw, req)
		})
	}
}
