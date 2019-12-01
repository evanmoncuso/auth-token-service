package rest

import (
	"auth-token-service/pkg/auth"
	"auth-token-service/pkg/db"
	"auth-token-service/pkg/middleware"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleAuthenticate(res http.ResponseWriter, req *http.Request) {
	middleware.EnableCors(&res)

	var c = new(credentials)

	defer req.Body.Close()

	reader, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reader, &c)
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := datastore.FindUserByUsername(c.Username)
	// TODO consider the case where this username doesn't exist for some reason
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
	}

	storeHash := []byte(user.Password)
	input := []byte(c.Password)

	err = bcrypt.CompareHashAndPassword(storeHash, input)
	if err != nil {
		handleHTTPError(res, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	token, expiration, err := auth.GenerateToken()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json")
	encoder.Encode(map[string]interface{}{
		"token":      token,
		"expiration": expiration,
		"lifetime":   auth.TokenLifespan,
	})
}
