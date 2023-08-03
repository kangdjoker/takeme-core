package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func JWTDecode(paramLog *basic.ParamLog, tokenString string) (domain.Claims, error) {

	claims := domain.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return domain.Claims{}, ErrorUnauthorized(paramLog)
		}
		return domain.Claims{}, ErrorUnauthorized(paramLog)
	}
	if !tkn.Valid {
		return domain.Claims{}, ErrorUnauthorized(paramLog)
	}

	return claims, nil
}

func JWTEncode(paramLog *basic.ParamLog, claimsAble domain.ClaimsAble, corporate domain.Corporate) (string, error) {
	expirationTime := time.Now().Add(time.Minute * time.Duration(corporate.TokenExpired))

	claims := &domain.Claims{
		SocketID:        claimsAble.GetID(),
		FullName:        claimsAble.GetFullName(),
		PhoneNumber:     claimsAble.GetFullName(),
		Verified:        claimsAble.GetVerified(),
		IsPinAlreadySet: claimsAble.GetIsPinAlreadySet(),
		CorporateID:     claimsAble.GetCorporateID(),
		CorporateName:   corporate.Name,
		AccessLevel:     claimsAble.GetAccessLevel(),
		Privileges:      claimsAble.GetPrivileges(),
		Resources:       corporate.Products,
		SAAS:            corporate.SAAS,
		CorporateURL:    corporate.DashboardURL,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", ErrorInternalServer(paramLog, DecodeTokenFailed, err.Error())
	}

	return tokenString, nil
}
