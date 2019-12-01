package rest

import (
	"auth-token-service/pkg/auth"
	"auth-token-service/pkg/db"

	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var InvalidatedTokens = make(map[string]int64)

func HandleInvalidateToken(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		encoder := json.NewEncoder(res)
		res.Header().Set("Content-Type", "application/json")
		encoder.Encode(map[string]interface{}{
			"invalidatedTokens": InvalidatedTokens,
		})

	case http.MethodPost:
		invalidateToken(res, req)

	default:
		res.WriteHeader(http.StatusNotImplemented)
	}
}

func invalidateToken(res http.ResponseWriter, req *http.Request) {
	bearerToken := req.Header.Get("Authorization")

	if len(bearerToken) <= 7 {
		handleHTTPError(res, "No Token included", http.StatusInternalServerError)
		return
	}

	tokenString := bearerToken[7:]

	_, err := auth.ValidateToken(tokenString)
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.ParseToken(tokenString)
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	expireTime := token.Claims.(*jwt.StandardClaims).ExpiresAt

	err = datastore.SetInvalidToken(tokenString, expireTime)
	if err != nil {
		handleHTTPError(res, "Unable to invalidate token", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
