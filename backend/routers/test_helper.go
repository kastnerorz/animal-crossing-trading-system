package routers

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var router = SetupRouter()

func PerformRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	router.ServeHTTP(r, req)
	return r
}

func PerformRequestWithAuth(method, url string, body io.Reader, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r := httptest.NewRecorder()
	router.ServeHTTP(r, req)
	return r
}

var QuotationId string
var QuotationIdPassCode string
var ApplicationId string
var ReviewerToken string
var ApplicantToken string
