package msedge

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	src "github.com/guquan-lengyue/tts/ttsclient"
)

const trustedClientToken = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
const binaryDelim = "Path:audio\r\n"
const initMessage = `
Content-Type:application/json; charset=utf-8
Path:speech.config

{
	"context": {
		"synthesis": {
			"audio": {
				"metadataoptions": {
					"sentenceBoundaryEnabled": "false",
					"wordBoundaryEnabled": "false"
				},
				"outputFormat": "%s" 
			}
		}
	}
}
`
const ttsMessage = `
X-RequestId: %s
Content-Type:application/ssml+xml
Path:ssml

<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="https://www.w3.org/2001/mstts" xml:lang="%s">
	<voice name="%s">
		<prosody pitch='+%fHz' rate ='+%f%%' volume='+%f%%'>
			%s
		</prosody>
	</voice>
</speak>
`

var lock sync.Mutex

const heartBeatTime = 28 * time.Second
const overTime = 15 * time.Minute

var _ src.ITtsClient = (*MsEdgeTTS)(nil)

type MsEdgeTTS struct {
	// enableLogger 是否打印日志
	enableLogger bool
	// outputFormat 音频格式
	outputFormat OutputFormat
	// ws overtime
	overTime time.Duration
	// ws ws客户端
	ws *websocket.Conn
	// voiceName 朗读者名称
	voiceName string
	// pitch 声音码率
	pitch float32
	// rate 朗读速度
	rate float32
	// volume 音量
	volume float32
	timer  *time.Timer
	// 客户端名称
	clientName string
}

func NewClient(clientName string, enableLogger bool) src.ITtsClient {
	lock.Lock()
	defer lock.Unlock()
	return &MsEdgeTTS{
		enableLogger: enableLogger,
		clientName:   clientName,
	}
}
func (m *MsEdgeTTS) closeMs() {
	m.timer.Stop()
	m.ws.Close()
	m.ws = nil
}
func (m *MsEdgeTTS) ttl() {
	if m.timer != nil {
		m.timer.Stop()
	}
	m.log("overTime left", m.overTime)
	m.timer = time.AfterFunc(heartBeatTime, func() {
		m.overTime -= heartBeatTime
		if m.overTime <= 0 {
			m.closeMs()
			log.Println("静默超时,主动断开链接")
			return
		}
		go func() {
			defer func() {
				if recover() != nil {
					log.Println(m.clientName, "客户端错误, 正则重启客户端")
					m.initWsClient()
				}
			}()
			speech := m.sendSsmlTemplate(fmt.Sprintf("hello %d", m.overTime))
			l := 0
			for i := range speech {
				l += len(i)
			}
			if l <= 0 {
				m.closeMs()
				m.log("heart beat error, ws close")
				return
			}
			m.log("heart beat ...")
		}()
	})
}

func (m *MsEdgeTTS) log(a ...any) {
	if m.enableLogger {
		log.Println(m.clientName + "--------------------")
		log.Println(a...)
	}
}

// initWsClient 初始化ws客户端
func (m *MsEdgeTTS) initWsClient() error {
	header := http.Header{}
	u := url.URL{
		Scheme:   "wss",
		Host:     "speech.platform.bing.com",
		Path:     "/consumer/speech/synthesize/readaloud/edge/v1",
		RawQuery: fmt.Sprintf("TrustedClientToken=%s", trustedClientToken),
	}
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.36")
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	log.Println(m.clientName, "tts server connected....")
	if err != nil {
		m.log(err)
		return err
	}
	m.ws = ws
	// 发送初始化信息
	m.send(fmt.Sprintf(initMessage, m.outputFormat))
	m.ws.SetCloseHandler(func(code int, text string) error {
		m.log("code :", code, "text :", text)
		m.log("链接已经断开")
		return nil
	})
	return nil
}

func (m *MsEdgeTTS) SetClient(voiceName string, rate float32, volume float32) {
	outputFormat := WEBM_24KHZ_16BIT_MONO_OPUS
	var pitch float32 = 200
	oldVoiceName := m.voiceName
	oldOutputFormat := m.outputFormat
	oldPitch := m.pitch
	oldRate := m.rate
	oldVolume := m.volume

	m.voiceName = voiceName
	m.outputFormat = outputFormat
	m.pitch = pitch
	m.rate = rate
	m.volume = volume

	change := oldVoiceName != voiceName ||
		oldOutputFormat != outputFormat ||
		oldPitch != pitch ||
		oldRate != rate ||
		oldVolume != volume

	if change {
		if m.ws != nil {
			err := m.ws.Close()
			if err != nil {
				m.log(err)
			}
		}
		m.initWsClient()
	}
}

func (m *MsEdgeTTS) TextToSpeech(input string) []byte {
	m.overTime = overTime
	speech := m.sendSsmlTemplate(input)
	buffer := bytes.Buffer{}
	for i := range speech {
		buffer.Write(i)
	}
	return buffer.Bytes()
}

func (m *MsEdgeTTS) sendSsmlTemplate(input string) chan []byte {
	ssmlTemplate := m.ssmlTemplate(input)
	m.send(ssmlTemplate)
	buf := make(chan []byte)
	m.listen(buf)
	return buf
}

// send 发送消息
func (m *MsEdgeTTS) send(message string) {
	msg := strings.ReplaceAll(message, "\t", "")
	msg = strings.ReplaceAll(msg, "\n", "\r\n")
	msg = strings.Trim(msg, "\r\n")

	for m.ws == nil {
		m.initWsClient()
		time.Sleep(3 * time.Second)
	}
	if m.ws == nil {
		m.log(errors.New("无法建立链接"))
		return
	}

	err := m.ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		m.log(err)
		return
	}
	m.ttl()
}

// listen 接收数据
func (m *MsEdgeTTS) listen(out chan []byte) {
	go func() {
		defer func() {
			if recover() != nil {
				m.log(recover())
			}
			close(out)
		}()
		for {
			_, message, err := m.ws.ReadMessage()
			if err != nil {
				m.log(err)
				close(out)
				break
			}
			if strings.Contains(string(message), "Path:turn.end") {
				close(out)
				break
			}
			sep := []byte(binaryDelim)
			audioIndex := bytes.Index(message, sep)
			if audioIndex > 0 {
				audioData := message[audioIndex+len(sep):]
				out <- audioData
			}
		}
	}()
}

func (m *MsEdgeTTS) ssmlTemplate(input string) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	requestId := hex.EncodeToString(b)
	return fmt.Sprintf(ttsMessage,
		requestId,
		string([]rune(m.voiceName)[:5]),
		m.voiceName,
		m.pitch,
		m.rate,
		m.volume,
		input,
	)
}
