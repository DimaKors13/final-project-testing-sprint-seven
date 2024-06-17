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

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	params := url.Values{}
	params.Set("city", "moscow")
	params.Set("count", strconv.Itoa(totalCount))

	req := httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)

	handler := http.HandlerFunc(mainHandle)

	//Запрос корректный
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusOK, responseRecorder.Code, "The cod is not 200")
	assert.NotEmpty(t, responseRecorder.Body, "Responded body is empty")

	//Город, который передаётся в параметре city, не поддерживается
	responseRecorder = httptest.NewRecorder()
	params.Set("city", "rostov-on-don")
	req = httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "The cod is not 400")
	body := responseRecorder.Body.String()
	assert.Equal(t, body, "wrong city value", "Wrong body content when the city is incorrect")

	//В параметре count указано больше, чем есть всего
	responseRecorder = httptest.NewRecorder()
	params.Set("city", "moscow")
	params.Set("count", strconv.Itoa(totalCount+1))
	req = httptest.NewRequest("GET", "/cafe?"+params.Encode(), nil)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusOK, responseRecorder.Code, "The cod is not 200 when the count is more than total")
	body = responseRecorder.Body.String()
	cafes := strings.Split(body, ",")
	assert.Equal(t, totalCount, len(cafes), "Wrong body content when the count is more than total")
}
