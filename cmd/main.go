package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logger := logrus.New()
	port := 5000

	r := chi.NewRouter()

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	// add swagger middleware
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port)), // The url pointing to API definition
	))

	logger.Infof("server started at port %v", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatalf("server can't listen requests")
		}
	}()

	logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", port)
}
