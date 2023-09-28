package main

import (
	"github.com/gin-gonic/gin"
	"ms_edge_tts/src"
	"net/http"
	"strconv"
)

var tts = src.NewMsEdgeTTS(true)

func main() {
	r := setRouter()
	err := r.Run("127.0.0.1:2580")
	if err != nil {
		panic(err)
	}
}

func setRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", receive)
	return r
}
func receive(c *gin.Context) {
	text := c.DefaultQuery("tex", "")                        // text
	speakSpeed := c.DefaultQuery("spd", "0")                 // speak speed
	voiceName := c.DefaultQuery("vn", "zh-CN-YunyangNeural") // voiceName
	spd, err := strconv.Atoi(speakSpeed)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]error{"error": err})
	}
	tts.SetMetaData(voiceName, src.WEBM_24KHZ_16BIT_MONO_OPUS, 0, spd, 0)
	speechCh := tts.TextToSpeech(text)
	c.Header("Context-Type", "Content-Type: audio/webm")
	for ch := range speechCh {
		_, err := c.Writer.Write(ch)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
			break
		}
	}
}
