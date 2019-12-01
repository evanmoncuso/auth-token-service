package middleware

// EnableCors lets cors work because it's stupid
import (
	"net/http"
)

// EnableCors allows cors to be set
func EnableCors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
}
