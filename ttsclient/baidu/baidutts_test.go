package baidu

import (
	"os"
	"testing"
)

const test = "“赠尹兆先与君结识于谷雨之后，暂别于芒种之前，余深居小阁，县内友人唯君一人尔；忆往昔，摊桌初遇尚觉浅，笑言尹兄故孤高；然，君虽仅一县夫子，无愧圣贤之书，知理而善学，善学而擅改，学而时习，自勉自强；君子有欲明晰取之有道，小民常乐不扰他人一分，何人？宁安尹兆先也；只惜，天无皓月常清，地无宴席不散，星斗挂天余自去，君莫怪；夜走不辞别，临行赠一贴，对坐再弈棋，相逢会有期；望君，教书育人作于细，功参社稷勿须臾，持心如初，从始至终；他日著书立传，惠得百家子弟，教化天下万民，一代大儒皆可期；当是时，可游山川，踏天地，惊涛骇浪不改色，凌波微步亦自若，腹墨千千万，胸中有正气！”"

func TestBaiduTTS_TextToSpeech(t *testing.T) {
	client := NewClient("1", true)
	client.SetClient("5003", 6, 2)
	a := client.TextToSpeech(test)
	file, err := os.OpenFile("a.mp3", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	for bytes := range a {
		_, err = file.Write(bytes)
		if err != nil {
			panic(err)
		}
	}
}
