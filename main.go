package main

import (
	"flag"
	"github.com/guquan-lengyue/ms_edge_tts/data"
	"github.com/guquan-lengyue/ms_edge_tts/service"
)

func main() {
	port := flag.Int("port", 2580, "listen port")
	host := flag.String("host", "0.0.0.0", "listen host")
	clientType := flag.String("ct", "baidu", "client type")
	dbHost := flag.String("dbhost", "", "db host")
	dbUsr := flag.String("dbusr", "", "db usr")
	dbPass := flag.String("dbpass", "", "db pass")
	dbName := flag.String("dbname", "", "db name")
	flag.Parse()

	data.InitDb(*dbHost, *dbUsr, *dbPass, *dbName)

	service.Service(*clientType, *host, *port)
}
