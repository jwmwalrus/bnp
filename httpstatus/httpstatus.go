package httpstatus

import "net/http"

// IsInformational checks if the status code means informational
func IsInformational(r *http.Response) bool {
	return r != nil && r.StatusCode >= 100 && r.StatusCode <= 199
}

// IsSuccess checks if the status code means success
func IsSuccess(r *http.Response) bool {
	return r != nil && r.StatusCode >= 200 && r.StatusCode <= 299
}

// IsRedirection checks if the status code means redirection
func IsRedirection(r *http.Response) bool {
	return r != nil && r.StatusCode >= 300 && r.StatusCode <= 399
}

// IsClientError checks if the status code means client error
func IsClientError(r *http.Response) bool {
	return r != nil && r.StatusCode >= 400 && r.StatusCode <= 499
}

// IsError checks if the status code means client error
func IsError(r *http.Response) bool {
	return IsClientError(r) || IsServerError(r)
}

// IsServerError checks if the status code means client error
func IsServerError(r *http.Response) bool {
	return r != nil && r.StatusCode >= 500 && r.StatusCode <= 599
}
