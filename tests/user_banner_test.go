package tests

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	router "avito-backend-trainee-2024/pkg/route"
	jwtutils "avito-backend-trainee-2024/pkg/utils/jwt"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
)

func (s *Suite) TestGetNotExistingBannerByUser() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       1,
		"username": "user",
		"is_admin": false,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "10")
	q.Set("tag_ids", "1,2,3")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusBadRequest, recorder.Result().StatusCode)

	respMsg := recorder.Body.String()

	assertions.Equal("error occurred fetching banner: no such banner", respMsg)
}

func (s *Suite) TestGetActiveBannerByUser() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       1,
		"username": "user",
		"is_admin": false,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "1")
	q.Set("tag_ids", "1,2")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusOK, recorder.Result().StatusCode)

	var banner entity.Banner

	s.NoError(json.NewDecoder(recorder.Body).Decode(&banner))

	assertions.Equal("title", banner.Content.Title)
	assertions.Equal("text", banner.Content.Text)
	assertions.Equal("http://url.com", banner.Content.Url)
}

func (s *Suite) TestGetNotExistingBannerByAdmin() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       2,
		"username": "admin",
		"is_admin": true,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "10")
	q.Set("tag_ids", "1,2,3")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusBadRequest, recorder.Result().StatusCode)

	respMsg := recorder.Body.String()

	assertions.Equal("error occurred fetching banner: no such banner", respMsg)
}

func (s *Suite) TestGetActiveBannerByAdmin() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       2,
		"username": "admin",
		"is_admin": true,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "1")
	q.Set("tag_ids", "1,2")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusOK, recorder.Result().StatusCode)

	var banner entity.Banner

	s.NoError(json.NewDecoder(recorder.Body).Decode(&banner))

	assertions.Equal("title", banner.Content.Title)
	assertions.Equal("text", banner.Content.Text)
	assertions.Equal("http://url.com", banner.Content.Url)
}

func (s *Suite) TestGetInactiveBannerByUser() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       1,
		"username": "user",
		"is_admin": false,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "2")
	q.Set("tag_ids", "1")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusForbidden, recorder.Result().StatusCode)

	respMsg := recorder.Body.String()

	assertions.Equal("banner is inactive", respMsg)
}

func (s *Suite) TestGetInactiveBannerByAdmin() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "/test/api/user_banner", nil)

	payload := map[string]any{ // this user should exist in db
		"id":       2,
		"username": "admin",
		"is_admin": true,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)
	s.NoError(err)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	q := req.URL.Query()

	q.Set("feature_id", "2")
	q.Set("tag_ids", "1")
	q.Set("use_last_revision", "true")

	req.URL.RawQuery = q.Encode()

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusOK, recorder.Result().StatusCode)

	var banner entity.Banner

	s.NoError(json.NewDecoder(recorder.Body).Decode(&banner))

	assertions.Equal("title2", banner.Content.Title)
	assertions.Equal("text2", banner.Content.Text)
	assertions.Equal("http://url2.com", banner.Content.Url)
}
