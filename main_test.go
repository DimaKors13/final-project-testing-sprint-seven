package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var totalCount = 4
var handler = http.HandlerFunc(mainHandle)

// Запрос корректный
func TestMainHandlerCorrectRequest(t *testing.T) {

	params := requestParams("moscow", totalCount)
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusOK, responseRecorder.Code, "The cod is not 200")
	assert.NotEmpty(t, responseRecorder.Body, "Responded body is empty")

}

// Город, который передаётся в параметре city, не поддерживается
func TestMainHandlerWhenWrongCity(t *testing.T) {

	params := requestParams("wrong-city", totalCount)
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "The cod is not 400")
	body := responseRecorder.Body.String()
	assert.Equal(t, body, "wrong city value", "Wrong body content when the city is incorrect")

}

// В параметре count указано больше, чем есть всего
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	params := requestParams("moscow", totalCount+1)
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusOK, responseRecorder.Code, "The cod is not 200 when the count is more than total")
	body := responseRecorder.Body.String()
	cafes := strings.Split(body, ",")
	assert.Len(t, cafes, totalCount, "Wrong body content when the count is more than total")
}

func requestParams(city string, count int) url.Values {
	result := url.Values{}
	result.Set("city", city)
	result.Set("count", strconv.Itoa(count))
	return result
}
