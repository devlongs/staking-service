package v1

import (
	"net/http"
)

// HealthHandler returns a handler that responds with a JSON health status
// TODO (LONGS): expandthe HealthHandler function
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"alive": true}`))
	}
}
