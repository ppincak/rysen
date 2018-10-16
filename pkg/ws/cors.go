package ws

import (
	"net/http"
)

// Allow any origin
func CheckOriginAny(request *http.Request) bool {
	return true
}
