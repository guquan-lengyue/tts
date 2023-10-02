package src

import (
	"log"
	"os"
	"testing"
)

func Test_MsEdgeTTS(t *testing.T) {
	tts := NewMsEdgeTTS(false)
	tts.SetMetaData(
		"zh-CN-XiaoxiaoNeural",
		WEBM_24KHZ_16BIT_MONO_OPUS,
		0,
		0,
		0,
	)
	speech := tts.TextToSpeech("tts错误")
	open, err := os.OpenFile("error.webm", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for bytes := range speech {
		log.Println(string(bytes))
		_, err = open.Write(bytes)
	}
	if err != nil {
		panic(err)
	}
	open.Close()
}
