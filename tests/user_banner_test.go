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

func (s *Suite) TestGetUserBanner() {
	assertions := s.Require()

	req, _ := http.NewRequest("GET", "test/api/user_banner", nil)

	payload := map[string]any{ // todo: this user should exist in db
		"id":       1,
		"username": "user",
		"is_admin": false,
	}

	token, err := jwtutils.CreateJWT(payload, jwt.SigningMethodHS256, jwtSecret)

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", token)

	req.URL.Query().Set("feature_id", "1")
	req.URL.Query().Set("tag_ids", "1,2")
	req.URL.Query().Set("use_last_revision", "true")

	routers := make(map[string]chi.Router)

	routers["/user_banner"] = s.bannerHandler.Routes()

	r := router.MakeRoutes("/test/api", routers)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assertions.Equal(http.StatusOK, recorder.Result().StatusCode)

	var banner entity.Banner

	err = json.NewDecoder(recorder.Body).Decode(&banner)

	s.NoError(err)

	assertions.Equal("title", banner.Content.Title)
	assertions.Equal("text", banner.Content.Text)
	assertions.Equal("http://url.com", banner.Content.Url)
}
