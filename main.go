package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"ms_edge_tts/assets"
	"ms_edge_tts/src"
	"net/http"
	"net/url"
	"strconv"
)

var tts = src.NewMsEdgeTTS(gin.Mode() == gin.DebugMode)

func main() {
	r := setRouter()
	r.Use(gzip.Gzip(gzip.BestCompression))
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
	Text      string  `json:"tex"`
	Speed     float32 `json:"spd"`
	VoiceName string  `json:"vn"`
}

func receive(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]error{"error": err})
		return
	}
	form, err := parseBody(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]error{"error": err})
		return
	}
	tts.SetMetaData(form.VoiceName, src.WEBM_24KHZ_16BIT_MONO_OPUS, 0, form.Speed, 0)
	log.Printf("request: %#v \n", form)
	size := 0
	for i := 0; i < 3 && size == 0; i++ {
		speechCh := tts.TextToSpeech(form.Text)
		c.Header("Context-Type", "Content-Type: audio/webm")
		for ch := range speechCh {
			size += len(ch)
			_, err := c.Writer.Write(ch)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
				break
			}
		}
	}
	log.Println("response size: ", size)
	if size == 0 {
		_, err = c.Writer.Write(assets.ErrorTttWebm)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
		}
	}
}

func parseBody(b []byte) (*body, error) {
	query, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	spd := query.Get("spd")
	speed, err := strconv.ParseFloat(spd, 10)
	if err != nil {
		return nil, err
	}
	text := query.Get("tex")
	text, err = url.QueryUnescape(text)
	if err != nil {
		return nil, err
	}
	text, err = url.QueryUnescape(text)
	if err != nil {
		return nil, err
	}
	return &body{
		Text:      text,
		Speed:     float32(speed),
		VoiceName: query.Get("vn"),
	}, nil
}
