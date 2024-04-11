package middleware

import (
	"avito-backend-trainee-2024/internal/handler/response"
	handlerutils "avito-backend-trainee-2024/pkg/utils/handler"
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MiddlewareData = map[string]any

func InMemUserBannerCache(cache *cache.Cache, logger *logrus.Logger) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Method != "GET" {
				next.ServeHTTP(rw, req) // cache only get requests
			}

			key := req.URL.RequestURI() // take requested uri as key
			req = req.WithContext(context.WithValue(req.Context(), req.URL.RequestURI(), make(MiddlewareData)))

			useLastRevision, err := handlerutils.GetStringParamFromQuery(req, "use_last_revision")
			if err != nil {
				msg := fmt.Sprintf("error occurred getting 'use_last_revision' param from query: %v", err)

				handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, msg, msg)
				return
			}

			var data MiddlewareData

			if val := req.Context().Value(key); val != nil {
				data = val.(MiddlewareData)
			}

			setToCacheFunc := func() {
				banner, exists := data["banner"]
				if exists {
					cache.Set(key, banner, 0)
				}

				isActive, exists := data["is_active"]
				if exists {
					cache.Set(key+"?is_active=", isActive, 0)
				}
			}

			if useLastRevision == "true" {
				next.ServeHTTP(rw, req)

				// set retrieved from db object to cache
				setToCacheFunc()
				return
			}

			if cached, found := cache.Get(key); found {
				banner, ok := cached.(response.GetUserBannerResponse)
				if !ok {
					msg := fmt.Sprintf("error occurred casting cached value to GetUserBannerResponse struct: %v", err)

					handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusInternalServerError, msg, "")
					return
				}

				isActive, exists := cache.Get(key + "?is_active=")
				if !exists {
					msg := "no key 'is_active' in MiddlewareData"

					handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusInternalServerError, msg, "")
					return
				}

				if isActive != "true" && req.Header.Get("is_admin") != "true" {
					msg := "banner is inactive"

					handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusNoContent, msg, msg)
					return
				}

				render.JSON(rw, req, banner)
				rw.WriteHeader(http.StatusOK)

				return
			}

			// serve main handler
			next.ServeHTTP(rw, req)

			// set retrieved from db object to cache
			setToCacheFunc()
		})
	}
}
