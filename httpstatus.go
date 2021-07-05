package bnp

import "net/http"

// IsHTTPInformational checks if the status code means informational
func IsHTTPInformational(r *http.Response) bool {
	return r != nil && r.StatusCode >= 100 && r.StatusCode <= 199
}

// IsHTTPSuccess checks if the status code means success
func IsHTTPSuccess(r *http.Response) bool {
	return r != nil && r.StatusCode >= 200 && r.StatusCode <= 299
}

// IsHTTPRedirection checks if the status code means redirection
func IsHTTPRedirection(r *http.Response) bool {
	return r != nil && r.StatusCode >= 300 && r.StatusCode <= 399
}

// IsHTTPClientError checks if the status code means client error
func IsHTTPClientError(r *http.Response) bool {
	return r != nil && r.StatusCode >= 400 && r.StatusCode <= 499
}

// IsHTTPError checks if the status code means client error
func IsHTTPError(r *http.Response) bool {
	return IsHTTPClientError(r) || IsHTTPServerError(r)
}

// IsHTTPServerError checks if the status code means client error
func IsHTTPServerError(r *http.Response) bool {
	return r != nil && r.StatusCode >= 500 && r.StatusCode <= 599
}
