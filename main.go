package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"ms_edge_tts/src"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

var ttsClients []*src.MsEdgeTTS

const clientNum = 10

func init() {
	for i := 0; i < clientNum; i++ {
		ttsClients = append(ttsClients, src.NewMsEdgeTTS(fmt.Sprintf("客户端[%d]", i), gin.Mode() == gin.DebugMode))
	}
}

func main() {
	port := flag.Int("port", 2580, "listen port")
	host := flag.String("host", "0.0.0.0", "listen host")
	flag.Parse()

	r := setRouter()
	r.Use(gzip.Gzip(gzip.BestCompression))
	err := r.Run(fmt.Sprintf("%s:%d", *host, *port))
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

var ban int32 = 0

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
	for {
		atomic.StoreInt32(&ban, (ban+1)%clientNum)
		err = getTts(form, c, ttsClients[ban])
		if err == nil {
			break
		}
	}
}

func getTts(form *body, c *gin.Context, tts *src.MsEdgeTTS) error {
	tts.SetMetaData(form.VoiceName, src.WEBM_24KHZ_16BIT_MONO_OPUS, 0, form.Speed, 0)
	log.Printf("request: %#v \n", form)
	size := 0
	for i := 0; i < 3 && size == 0; i++ {
		speechCh := tts.TextToSpeech(form.Text)
		c.Header("Context-Type", "Content-Type: audio/webm")
		audio := bytes.Buffer{}
		for ch := range speechCh {
			size += len(ch)
			audio.Write(ch)
		}
		_, err := c.Writer.Write(audio.Bytes())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
			break
		}
	}
	log.Println("response size: ", size)
	if size == 0 {
		return errors.New("tts错错误")
	}

	return nil
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
