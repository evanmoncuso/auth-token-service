package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Logger is just some fun logging middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("[ %s ] :: %s :: %v\n", req.Method, req.URL, time.Now())
		next.ServeHTTP(res, req)
	})
}
