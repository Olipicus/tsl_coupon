package main

import (
	"flag"
	"log"
	"net/http"

	"code.olipicus.com/trueselect_coupon/config"
	"code.olipicus.com/trueselect_coupon/line"
)

var configFile string

func initFlag() {
	flag.StringVar(&configFile, "config", "./config.json", "Configuration file location")
	flag.Parse()
}

func main() {

	initFlag()
	config := config.LoadConfig(configFile)

	channelSecret, _ := config.GetString("channel_secret")
	channelToken, _ := config.GetString("channel_token")

	app, err := line.NewLineApp(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	var port string
	if port, err = config.GetString("http_port"); err != nil {
		port = "8088" //Default Port
	}

	http.HandleFunc("/line", app.CallbackHandler)
	http.HandleFunc("/verify", verify)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func verify(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
