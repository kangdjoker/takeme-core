package utils

import (
	"encoding/json"
	"net/http"

	"github.com/kangdjoker/takeme-core/domain"
	log "github.com/sirupsen/logrus"
)

// ContextValue is a context key
type ContextValue map[string]interface{}

func CorporateContext(request *http.Request) domain.Corporate {
	data := request.Context().Value("data").(ContextValue)["corporate"].(domain.Corporate)
	return data
}

func AccessLevelByContext(request *http.Request) string {
	claims := request.Context().Value("data").(ContextValue)["claims"].(domain.Claims)

	return claims.AccessLevel
}

func UserContext(request *http.Request) domain.User {
	data := request.Context().Value("data").(ContextValue)["user"].(domain.User)
	return data
}

func LoadPayload(r *http.Request, payload interface{}) error {
	log.Info("------------------------ Client payload ------------------------")
	log.Info(
		string(r.Context().Value("payload").([]byte)[:]),
	)
	err := json.Unmarshal(r.Context().Value("payload").([]byte), &payload)
	if err != nil {
		return ErrorBadRequest(InvalidRequestPayload, "Invalid request payload")
	}
	log.Info("------------------------ Client payload ------------------------")

	return nil
}
