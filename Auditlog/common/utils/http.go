package utils

import "net/http"

func IsHttpErr(r *http.Response) bool {
	return r.StatusCode < http.StatusOK || r.StatusCode >= http.StatusMultipleChoices
}
