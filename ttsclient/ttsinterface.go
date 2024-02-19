package ttsclient

// ITtsClient tts客户端接口
type ITtsClient interface {
	// SetClient 用于设置tts客户端的运行信息
	// voiceName 参数用于定义客户端使用的声音声色
	// rate 语速
	// volume 音量
	SetClient(voiceName string, rate float32, volume float32)
	// TextToSpeech 将调用tts客户端, 将input内容转换为音频, 音频内容将在返回值chan中异步返回, 音频结束后方法会主动关闭chan
	// input 需要转为音频的文字
	// return 异步返回音频内容
	TextToSpeech(input string) chan []byte
}
