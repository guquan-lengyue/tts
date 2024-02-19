package baidu

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	urlUtil "net/url"
	"strings"
	"sync"
	"time"

	src "github.com/guquan-lengyue/ms_edge_tts/ttsclient"
)

var _ src.ITtsClient = (*BaiduTTS)(nil)

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

// NewClient 用于创建百度tts客户端实例
// clientName 客户端标识 主要用于打印日志时区分日志来源
// enableLogger 是否打印相识日志标识
func NewClient(clientName string, enableLogger bool) src.ITtsClient {
	return &BaiduTTS{
		enableLogger: enableLogger,
		clientName:   clientName,
	}
}

func (m *BaiduTTS) log(a ...any) {
	if m.enableLogger {
		log.Print(m.clientName + "----")
		log.Println(a...)
	}
}

func (t *BaiduTTS) SetClient(voiceName string, rate float32, volume float32) {
	// 百度tts朗读速度只有1-10 10级
	t.rate = int(rate)
	if t.rate > 10 {
		t.rate = 10
	}
	if t.rate < 0 {
		t.rate = 0
	}
	t.volume = volume
	t.voiceName = voiceName
}

func (t *BaiduTTS) TextToSpeech(input string) chan []byte {
	// 百度tts试用接口最大字数为200, 需要将input分段
	text := []rune(input)
	textLength := len(text)
	// ttsRst 是input分段后返回的结果, ttsRst[0]为input第一段的内音频内容
	// 主要是为了分段后异步请求tts且能保持返回结果为请求时的顺序, 保持朗读内容顺序正常
	ttsRst := make([]*[]byte, 0)
	// wg 用于等待所有input分段tts请求结束
	var wg sync.WaitGroup
	for idx := 0; idx < textLength; idx += 200 {
		end := idx + 200
		if end > textLength {
			end = textLength
		}
		begin := idx
		// region 异步请求input分段tts内容, 将内容按顺序存放到相应位置
		var buffer = new([]byte)
		ttsRst = append(ttsRst, buffer)
		wg.Add(1)
		go func() {
			defer wg.Done()
			*buffer = t.tts(string(text[begin:end]))
		}()
		// endregion
	}
	// 同步等待tts数据
	bch := make(chan []byte, 1)
	go func() {
		wg.Wait()
		for i := range ttsRst {
			tts := *ttsRst[i]
			// 过滤掉返回值为nil的tts数据
			if len(tts) > 0 {
				bch <- tts
			}
		}
		close(bch)
	}()
	return bch
}

// tts 请求百度tts试用接口将text转为音频
// 当发生错误, 或者接口返回音频内容为空时, 返回值为nil
func (t *BaiduTTS) tts(text string) []byte {
	// 根据百度tts接口文档, 将text内容进行2次url编码
	// 2次url编码为了将特殊字符能够正确传递
	s := urlUtil.QueryEscape(text)
	s = urlUtil.QueryEscape(s)

	form := fmt.Sprintf("type=tns&per=%s&spd=%d&pit=5&vol=5&aue=6&tex=%s", t.voiceName, t.rate, s)
	payload := strings.NewReader(form)

	client := &http.Client{}
	req, err := http.NewRequest(method, baiduHost, payload)

	if err != nil {
		t.log(err)
		return nil
	}
	timestamp := time.Now().UnixMicro()
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Add("Referer", fmt.Sprintf("https://ai.baidu.com/tech/speech/tts_online?_=%d", timestamp))
	req.Header.Add("Host", "ai.baidu.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		t.log(err)
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.log(err)
		return nil
	}
	b := &bodyContent{}
	err = json.Unmarshal(body, b)
	if err != nil {
		t.log(err)
		return nil
	}
	var b64 string
	// 删除base64前缀 data:audio/x-mpeg;base64,
	if len(b.Data) > 25 {
		b64 = b.Data[25:]
	} else {
		t.log(err)
		return nil
	}
	audio, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.log(err)
		return nil
	}
	return audio
}
