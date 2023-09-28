package main

import (
	"github.com/gin-gonic/gin"
	"ms_edge_tts/src"
	"net/http"
)

var tts = src.NewMsEdgeTTS(true)

func main() {
	r := setRouter()
	err := r.Run("0.0.0.0:2580")
	if err != nil {
		panic(err)
	}
}

func setRouter() *gin.Engine {
	r := gin.Default()

	r.POST("", receive)
	return r
}

type body struct {
	Text      string `json:"tex"`
	Speed     int    `json:"spd"`
	VoiceName string `json:"vn"`
}

func receive(c *gin.Context) {
	var form body
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]error{"error": err})
	}
	tts.SetMetaData(form.VoiceName, src.WEBM_24KHZ_16BIT_MONO_OPUS, 0, form.Speed, 0)
	speechCh := tts.TextToSpeech(form.Text)
	c.Header("Context-Type", "Content-Type: audio/webm")
	for ch := range speechCh {
		_, err := c.Writer.Write(ch)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
			break
		}
	}
}
