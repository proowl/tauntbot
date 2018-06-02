package main

import (
	"io/ioutil"
	"encoding/json"
	"os"
)

// ./bot etc/config.json
func main() {
	// init
	config_raw, _ := ioutil.ReadFile(os.Args[1])
	var appConfig Config
	if err := json.Unmarshal(config_raw, &appConfig); err != nil {
		panic(err)
	}

	bot := newBot(&appConfig)
	bot.start()
}
