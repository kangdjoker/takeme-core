package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	log "github.com/sirupsen/logrus"
)

func Middleware(h http.HandlerFunc, secure bool) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		basic.LogInformation(basic.ParamLog{Tag: r.Header.Get("requestID")}, "----------------------------- REQUEST START -----------------------------")

		var ctx context.Context
		var corporate domain.Corporate
		var claims domain.Claims
		var user domain.User

		err := setupContextHeader(r)
		if err != nil {
			utils.ResponseError(err, w, r)
			return
		}

		corporate, err = service.CorporateByRequest(r)
		if err != nil {
			utils.ResponseError(err, w, r)
			return
		}

		// If signature invalid reduce access_attempt
		err = validateSignature(r, corporate)
		if err != nil {
			go InvalidCorporateAuth(corporate)
			utils.ResponseError(err, w, r)
			return
		}

		if secure == true {
			claims, err = validateJWT(r)
			if err != nil {
				utils.ResponseError(err, w, r)
				return
			}

			user, err = service.UserByIDWithValidation(claims.SocketID, []func(domain.User) error{
				service.ValidateUserExist,
				service.ValidateUserLocked,
			})

			if err != nil {
				utils.ResponseError(err, w, r)
				return
			}
		}

		data := utils.ContextValue{
			"claims":    claims,
			"user":      user,
			"corporate": corporate,
		}

		ctx = context.WithValue(r.Context(), "data", data)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setupContextHeader(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	language := r.Header.Get(" Accept-Language")
	corporate := r.Header.Get("corporate")
	requestID := r.Header.Get("requestID")
	signature := r.Header.Get("signature")

	// Set default value
	if language == "" {
		language = "en"
	}

	fmt.Println()
	log.Info(fmt.Sprintf("Header : {contentType : %v, language: %v, corporate: %v, requestID: %v }",
		contentType, language, corporate, requestID))

	if corporate == "" || requestID == "" || language == "" || signature == "" {
		return utils.ErrorBadRequest(utils.InvalidHeader, "Invalid Header")
	}

	return nil
}

func logPayloadBaseonLength(payload []byte, requestID string, signature string, result string) {
	if len(payload) <= 4000 {
		log.Info(fmt.Sprintf("Original Signature for requestID %v : (%v)", requestID, signature))
		log.Info(fmt.Sprintf("Should be Signature for requestID %v : (%v)", requestID, result))
	} else {
		log.Info(fmt.Sprintf("Original Signature for requestID %v : (%v)", requestID, "..."))
		log.Info(fmt.Sprintf("Should be Signature for requestID %v : (%v)", requestID, "..."))
	}
}

func validateSignature(r *http.Request, corporate domain.Corporate) error {

	// Get secret by corporateID and signature
	signature := r.Header.Get("signature")
	requestID := r.Header.Get("requestID")
	payload := r.Context().Value("payload").([]byte)

	secretKey := corporate.Secret
	result := hmacSHA512(payload, []byte(secretKey))

	logPayloadBaseonLength(payload, requestID, signature, result)

	if result == signature {
		return nil
	}

	return utils.ErrorBadRequest(utils.InvalidCorporateKey, "Invalid secret")
}

func validateJWT(r *http.Request) (domain.Claims, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return domain.Claims{}, utils.ErrorUnauthorized()
	}

	tokenString := authorization[7:]

	claims, err := utils.JWTDecode(tokenString)
	if err != nil {
		return domain.Claims{}, utils.ErrorUnauthorized()
	}

	return claims, nil
}

func hmacSHA512(data, secret []byte) string {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha512.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha
}

func MiddlewareWithoutSignature(h http.HandlerFunc, secure bool) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx context.Context
		var corporate domain.Corporate
		var claims domain.Claims
		var err error

		if secure == true {
			claims, err = validateJWT(r)
			if err != nil {
				utils.ResponseError(err, w, r)
				return
			}

		}

		corporate, err = service.CorporateByRequest(r)
		if err != nil {
			utils.ResponseError(err, w, r)
			return
		}

		data := utils.ContextValue{
			"claims":    claims,
			"userID":    claims.SocketID,
			"corporate": corporate,
		}

		ctx = context.WithValue(r.Context(), "data", data)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AddPayloadContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO REMOVE THIS TEMPORARY SOLUTION
		_, _, err := r.FormFile("file")
		if err != nil {
			payload, err := ioutil.ReadAll(r.Body)
			if err != nil {
				utils.ResponseError(utils.ErrorBadRequest(utils.InvalidRequestPayload, "Invalid payload request"), w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "payload", payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}

	})
}
