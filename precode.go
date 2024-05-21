package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	//req := ... // здесь нужно создать запрос к сервису
	req, err := http.NewRequest("GET", "/cafe?count=5&city=moscow", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("count", strconv.Itoa(totalCount)) // Используем totalCount
	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// здесь нужно добавить необходимые проверки
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "status code is not OK")
	assert.NotEmpty(t, responseRecorder.Body.String(), "response body is empty")

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "status code is not BadRequest")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "response body mismatch")

	expectedResponse := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	actualResponse := responseRecorder.Body.String()
	assert.Equal(t, expectedResponse, actualResponse, "response body mismatch")

	// Проверяем, что длина ответа соответствует ожидаемой
	expectedCafes := 4
	assert.Len(t, strings.Split(actualResponse, ","), expectedCafes, "number of cafes is incorrect")
}
