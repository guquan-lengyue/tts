package main

import (
	"bytes"
	"fmt"
	"ms_edge_tts/src"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

const params = `tex=%25E7%2599%25BD%25E9%2593%25B6%25E4%25B9%258B%25E7%2588%25B9%25EF%25BC%259A%25E2%2580%259C%2525%25E2%2580%2594%25E2%2580%2594%25EF%25BC%2581%25E2%2580%259D&spd=9.5&vn=5003`

func Test_TTS(t *testing.T) {
	r := setRouter("baidu")
	request(t, r)
	time.Sleep(10 * time.Minute)
}

func request(t *testing.T, r *gin.Engine) {
	w := httptest.NewRecorder()
	reader := bytes.NewReader([]byte(params))
	req, _ := http.NewRequest("POST", "/", reader)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	size := len(w.Body.Bytes())
	assert.Equal(t, size > 0, true)
	t.Log("接收到文件大小: ", size)
}

func Test_BaiduTts(t *testing.T) {
	tts := src.BaiduTTS{}
	audio := tts.TextToSpeech("你好世界")
	fmt.Println(<-audio)
}
