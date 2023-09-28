package src

type OutputFormat string

const (
	RAW_16KHZ_16BIT_MONO_PCM         OutputFormat = "raw-16khz-16bit-mono-pcm"
	RAW_24KHZ_16BIT_MONO_PCM         OutputFormat = "raw-24khz-16bit-mono-pcm"
	RAW_48KHZ_16BIT_MONO_PCM         OutputFormat = "raw-48khz-16bit-mono-pcm"
	RAW_8KHZ_8BIT_MONO_MULAW         OutputFormat = "raw-8khz-8bit-mono-mulaw"
	RAW_8KHZ_8BIT_MONO_ALAW          OutputFormat = "raw-8khz-8bit-mono-alaw"
	RAW_16KHZ_16BIT_MONO_TRUESILK    OutputFormat = "raw-16khz-16bit-mono-truesilk"
	RAW_24KHZ_16BIT_MONO_TRUESILK    OutputFormat = "raw-24khz-16bit-mono-truesilk"
	RIFF_16KHZ_16BIT_MONO_PCM        OutputFormat = "riff-16khz-16bit-mono-pcm"
	RIFF_24KHZ_16BIT_MONO_PCM        OutputFormat = "riff-24khz-16bit-mono-pcm"
	RIFF_48KHZ_16BIT_MONO_PCM        OutputFormat = "riff-48khz-16bit-mono-pcm"
	RIFF_8KHZ_8BIT_MONO_MULAW        OutputFormat = "riff-8khz-8bit-mono-mulaw"
	RIFF_8KHZ_8BIT_MONO_ALAW         OutputFormat = "riff-8khz-8bit-mono-alaw"
	AUDIO_16KHZ_32KBITRATE_MONO_MP3  OutputFormat = "audio-16khz-32kbitrate-mono-mp3"
	AUDIO_16KHZ_64KBITRATE_MONO_MP3  OutputFormat = "audio-16khz-64kbitrate-mono-mp3"
	AUDIO_16KHZ_128KBITRATE_MONO_MP3 OutputFormat = "audio-16khz-128kbitrate-mono-mp3"
	AUDIO_24KHZ_48KBITRATE_MONO_MP3  OutputFormat = "audio-24khz-48kbitrate-mono-mp3"
	AUDIO_24KHZ_96KBITRATE_MONO_MP3  OutputFormat = "audio-24khz-96kbitrate-mono-mp3"
	AUDIO_24KHZ_160KBITRATE_MONO_MP3 OutputFormat = "audio-24khz-160kbitrate-mono-mp3"
	AUDIO_48KHZ_96KBITRATE_MONO_MP3  OutputFormat = "audio-48khz-96kbitrate-mono-mp3"
	AUDIO_48KHZ_192KBITRATE_MONO_MP3 OutputFormat = "audio-48khz-192kbitrate-mono-mp3"
	WEBM_16KHZ_16BIT_MONO_OPUS       OutputFormat = "webm-16khz-16bit-mono-opus"
	WEBM_24KHZ_16BIT_MONO_OPUS       OutputFormat = "webm-24khz-16bit-mono-opus"
	OGG_16KHZ_16BIT_MONO_OPUS        OutputFormat = "ogg-16khz-16bit-mono-opus"
	OGG_24KHZ_16BIT_MONO_OPUS        OutputFormat = "ogg-24khz-16bit-mono-opus"
	OGG_48KHZ_16BIT_MONO_OPUS        OutputFormat = "ogg-48khz-16bit-mono-opus"
)
