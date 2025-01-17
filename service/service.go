package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync/atomic"

	"github.com/guquan-lengyue/tts/assets"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/guquan-lengyue/tts/ttsclient"
	"github.com/guquan-lengyue/tts/ttsclient/baidu"
	"github.com/guquan-lengyue/tts/ttsclient/msedge"
)

var ttsClients []ttsclient.ITtsClient

const clientNum = 10

func initClients(clientType string) {
	for i := 0; i < clientNum; i++ {
		if clientType == "baidu" {
			ttsClients = append(ttsClients, baidu.NewClient(fmt.Sprintf("客户端[%d]", i), gin.Mode() == gin.DebugMode))
		}
		if clientType == "mstts" {
			ttsClients = append(ttsClients, msedge.NewClient(fmt.Sprintf("客户端[%d]", i), gin.Mode() == gin.DebugMode))
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
	Volume    float32 `json:"v"`
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
	for i := 3; i > 0; i-- {
		atomic.StoreInt32(&ban, (ban+1)%clientNum)
		err = getTts(form, c, ttsClients[ban])
		if err == nil {
			return
		}
	}
	c.Header("Context-Type", "Content-Type: audio/webm")
	_, _ = c.Writer.Write(assets.ErrorTttWebm)
}

func getTts(form *body, c *gin.Context, tts ttsclient.ITtsClient) error {
	tts.SetClient(form.VoiceName, form.Speed, form.Volume)
	log.Printf("request: %#v \n", form)
	size := 0
	for i := 0; i < 3; i++ {
		speechCh := tts.TextToSpeech(form.Text)
		size = len(speechCh)
		if size == 0 {
			continue
		}
		c.Header("Context-Type", "Content-Type: audio/webm")
		_, err := c.Writer.Write(speechCh)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]error{"error": err})
			break
		}
		break
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
	if err != nil {
		return nil, err
	}
	volume := query.Get("v")
	if err != nil {
		return nil, err
	}
	v, err := strconv.ParseFloat(volume, 32)
	if err != nil {
		return nil, err
	}
	rst := &body{
		Speed:     float32(speed),
		VoiceName: query.Get("vn"),
		Volume:    float32(v),
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
