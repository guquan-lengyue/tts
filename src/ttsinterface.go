package src

type ITts interface {
	TextToSpeech(input string) chan []byte 
	SetMetaData(voiceName string, outputFormat OutputFormat, pitch float32, rate float32, volume float32) 
}
