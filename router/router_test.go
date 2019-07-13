package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/unknown/route", nil)

	rr := httptest.NewRecorder()
	router := GetRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRecover(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/500", nil)

	rr := httptest.NewRecorder()
	router := GetRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/test/500", nil)

	rr := httptest.NewRecorder()
	router := GetRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
