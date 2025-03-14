package helpers

import (
	"net/http"
)

func GetUserID(r *http.Request) string {
	return r.Context().Value("user_id").(string)
}
