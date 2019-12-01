package main

import (
	"auth-token-service/pkg/auth"
	"auth-token-service/pkg/db"
	"auth-token-service/pkg/http/rest"
	"auth-token-service/pkg/middleware"
	"auth-validation/pkg/validate"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	if port == "" {
		panic("No PORT specified for app")
	} else if tokenSecret == "" {
		panic("No TOKEN_SECRET specified for app")
	}

	// connect to USERS DB
	dbConnectionURL := os.Getenv("DATABASE_URL")
	err := datastore.InitializePostgres(dbConnectionURL)
	if err != nil {
		panic(err)
	}

	// connect to INVALID TOKEN store
	redisConnectionURL := os.Getenv("REDIS_URL")
	err = datastore.InitializeRedis(redisConnectionURL)
	if err != nil {
		panic(err)
	}

	// setup validation
	err = validate.Initialize(redisConnectionURL, tokenSecret)
	if err != nil {
		panic(err)
	}

	// set TokenSecret in auth
	auth.TokenSecret = []byte(tokenSecret)

	connection := fmt.Sprintf("127.0.0.1:%s", port)

	router := mux.NewRouter()
	router.Use(middleware.Logger)

	protected := router.PathPrefix("").Subrouter()
	protected.Use(validate.AuthenticateRoute)

	router.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}).Methods("GET")

	router.HandleFunc("/authenticate", rest.HandleAuthenticate).Methods(http.MethodPost)

	protected.HandleFunc("/invalidate", rest.HandleInvalidateToken).Methods(http.MethodGet, http.MethodPost)

	srv := &http.Server{
		Handler:      router,
		Addr:         connection,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
