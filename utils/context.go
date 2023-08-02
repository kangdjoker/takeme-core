package utils

import (
	"encoding/json"
	"net/http"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/basic"
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
	trClose, span, tag := basic.RequestToTracing(r)
	paramLog := basic.ParamLog{TrCloser: trClose, Span: span, Tag: tag}
	basic.LogInformation(&paramLog, "------------------------ Client payload ------------------------")
	basic.LogInformation(&paramLog,
		string(r.Context().Value("payload").([]byte)[:]),
	)
	err := json.Unmarshal(r.Context().Value("payload").([]byte), &payload)
	if err != nil {
		return ErrorBadRequest(InvalidRequestPayload, "Invalid request payload")
	}
	basic.LogInformation(&paramLog, "------------------------ Client payload ------------------------")

	return nil
}
