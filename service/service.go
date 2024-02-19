package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync/atomic"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/guquan-lengyue/ms_edge_tts/ttsclient"
	"github.com/guquan-lengyue/ms_edge_tts/ttsclient/baidu"
	"github.com/guquan-lengyue/ms_edge_tts/ttsclient/msedge"
)

var ttsClients []ttsclient.ITtsClient

const clientNum = 10

func initClients(clientType string) {
	for i := 0; i < clientNum; i++ {
		if clientType == "baidu" {
			ttsClients = append(ttsClients, baidu.NewBaiduTTS(fmt.Sprintf("客户端[%d]", i), gin.Mode() == gin.DebugMode))
		}
		if clientType == "mstts" {
			ttsClients = append(ttsClients, msedge.NewMsEdgeTTS(fmt.Sprintf("客户端[%d]", i), gin.Mode() == gin.DebugMode))
		}
	}
}

func Service(clientType, host string, port int) {
	r := setRouter(clientType)
	r.Use(gzip.Gzip(gzip.BestCompression))
	err := r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
}

func setRouter(clientType string) *gin.Engine {
	initClients(clientType)

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

func getTts(form *body, c *gin.Context, tts ttsclient.ITtsClient) error {
	tts.SetClient(form.VoiceName, form.Speed, 0)
	log.Printf("request: %#v \n", form)
	size := 0
	audio := bytes.Buffer{}
	for i := 0; i < 3 && size == 0; i++ {
		audio.Reset()
		speechCh := tts.TextToSpeech(form.Text)
		c.Header("Context-Type", "Content-Type: audio/webm")
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

var chineseReg = regexp.MustCompile(`\p{Han}`)

func parseBody(b []byte) (*body, error) {
	query, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	spd := query.Get("spd")
	speed, err := strconv.ParseFloat(spd, 32)
	rst := &body{
		Speed:     float32(speed),
		VoiceName: query.Get("vn"),
	}
	if err != nil {
		return nil, err
	}
	text := query.Get("tex")
	text, err = url.QueryUnescape(text)
	if err != nil {
		return nil, err
	}
	if chineseReg.MatchString(text) {
		rst.Text = text
		return rst, nil
	}
	text, err = url.QueryUnescape(text)
	if err != nil {
		return nil, err
	}
	rst.Text = text
	return rst, nil
}
