package src

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlUtil "net/url"
	"strings"
	"sync"
	"time"
)

type BaiduTTS struct {
	// enableLogger 是否打印日志
	enableLogger bool
	// 客户端名称
	clientName string
	// rate 朗读速度
	rate int
	// volume 音量
	volume float32
	// 声音名称
	voiceName string
}

const baiduHost = "https://ai.baidu.com/aidemo"
const method = "POST"

type bodyContent struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Data  string `json:"data"`
}

func NewBaiduTTS(clientName string, enableLogger bool) ITts {
	lock.Lock()
	defer lock.Unlock()
	var m ITts = &BaiduTTS{
		enableLogger: enableLogger,
		clientName:   clientName,
	}
	return m
}

func (t *BaiduTTS) SetMetaData(voiceName string, _ OutputFormat, _ float32, rate float32, volume float32) {
	t.rate = int(rate)
	t.volume = volume
	t.voiceName = voiceName
}

func (t *BaiduTTS) TextToSpeech(input string) chan []byte {
	bch := make(chan []byte, 1)
	text := []rune(input)
	textLength := len(text)
	rst := make([]*[]byte, 0)
	var wg sync.WaitGroup
	for idx := 0; idx < textLength; idx += 200 {
		end := idx + 200
		if end > textLength {
			end = textLength
		}
		begin := idx
		wg.Add(1)
		var buffer *[]byte = new([]byte)
		rst = append(rst, buffer)
		go func() {
			defer wg.Done()
			*buffer = t.getRst(string(text[begin:end]))
		}()
	}
	go func() {
		wg.Wait()
		for i := range rst {
			bch <- *rst[i]
		}
		close(bch)
	}()
	return bch
}

func (t *BaiduTTS) getRst(text string) []byte {
	s := urlUtil.QueryEscape(text)
	s = urlUtil.QueryEscape(s)

	form := fmt.Sprintf("type=tns&per=%s&spd=%d&pit=5&vol=5&aue=6&tex=%s", t.voiceName, t.rate, s)
	payload := strings.NewReader(form)

	client := &http.Client{}
	req, err := http.NewRequest(method, baiduHost, payload)

	if err != nil {
		fmt.Println(err)
	}
	timestamp := time.Now().UnixMicro()
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Add("Referer", fmt.Sprintf("https://ai.baidu.com/tech/speech/tts_online?_=%d", timestamp))
	req.Header.Add("Host", "ai.baidu.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	b := &bodyContent{}
	err = json.Unmarshal(body, b)
	if err != nil {
		fmt.Println(err)
	}
	b64 := strings.ReplaceAll(b.Data, "data:audio/x-mpeg;base64,", "")
	audio, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		fmt.Println(err)
	}

	return audio
}
