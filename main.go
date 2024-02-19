package main

import (
	"flag"

	"github.com/guquan-lengyue/ms_edge_tts/service"
)

func main() {
	port := flag.Int("port", 2580, "listen port")
	host := flag.String("host", "0.0.0.0", "listen host")
	clientType := flag.String("ct", "baidu", "client type, [mstts,baidu]")
	flag.Parse()
	service.Service(*clientType, *host, *port)
}
