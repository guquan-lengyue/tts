package main

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const params = `tex=%25E6%2589%258B%25EF%25BC%258C%25E7%25BA%25B7%25E7%25BA%25B7%25E8%25A2%25AB%25E6%2596%25A9%25E6%259D%2580%25E7%259A%2584%25E4%25B8%2580%25E5%25B9%25B2%25E4%25BA%258C%25E5%2587%2580%25E3%2580%2582&spd=9.5&vn=zh-CN-YunjianNeural`

func Test_TTS(t *testing.T) {
	r := setRouter()
	w := httptest.NewRecorder()
	reader := bytes.NewReader([]byte(params))
	req, _ := http.NewRequest("POST", "/", reader)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	open, _ := os.OpenFile("text.webm", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	_, err := open.Write(w.Body.Bytes())
	if err != nil {
		return
	}
}

func Test_UrlParams(t *testing.T) {

	parse, err := url.ParseQuery(params)
	if err != nil {
		t.Error(err)
	}
	t.Log(parse.Get("a"))

}
