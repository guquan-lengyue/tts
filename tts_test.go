package main

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_TTS(t *testing.T) {
	r := setRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?tex=你好世界&spd=20", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	open, _ := os.OpenFile("text.webm", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	_, err := open.Write(w.Body.Bytes())
	if err != nil {
		return
	}
}
