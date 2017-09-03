package main

import (
	"fmt"
	"time"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"telegram"
	"strings"
)

type Config struct {
	TauntBotConf telegram.BotConfig
}

var AppConfig Config

type CommandType uint8

const (
	TauntCommand = iota + 1
	ShrugCommand
)

func log(format string, smth ...interface{}) {
	fmt.Printf("[%v] " + format + "\n", time.Now(), smth)
}

// TODO: don't answer to too old messages
func process_updates(grammars* GrammarRules, bot_state* telegram.BotState) {
	updates, err := telegram.GetUpdates(&AppConfig.TauntBotConf, bot_state)
	if err != nil {
		log("GetUpdates failed with error: %v", err)
		return
	}
	if len(updates) > 0 {
		log("Received %d updates", len(updates))
		for _, msg := range updates {
			printed, _ := json.Marshal(msg)
			fmt.Println(string(printed))
			// todo: change to regex match instead full string match
			var response string
			command := strings.Split(msg.Message.Text, "@")[0]
			switch command {
				case "/shrug":
					response = "¯\\_(ツ)_/¯"
				case "/taunt":
					response = Taunt(grammars, "eng", msg.Message.Text)
				default:
					// ignore
			}
			if response != "" {
				log("Sending response: %v", response)
				err := telegram.SendMessage(&AppConfig.TauntBotConf, telegram.OutgoingMessage{ChatId: msg.Message.Chat.Id, Text: response, DisableNotification: true})
				if err != nil {
					log("SendMessage failed with error: %v", err)
				}
			}
			bot_state.LastUpdateId = msg.Update_id
		}
	}
}

func main() {
	// init
	grammars := LoadLangs("etc/taunt")
	rand.Seed(time.Now().UnixNano())
	config_raw, _ := ioutil.ReadFile("etc/config.json")
	if err := json.Unmarshal(config_raw, &AppConfig); err != nil {
		panic(err)
	}
	restored_state_raw, _ := ioutil.ReadFile(AppConfig.TauntBotConf.StateFile)
	var bot_state telegram.BotState
	if err := json.Unmarshal(restored_state_raw, &bot_state); err != nil {
		bot_state.LastUpdateId = AppConfig.TauntBotConf.StartUpdateId
	}
	process_updates(&grammars, &bot_state)
	updated_state, _ := json.Marshal(bot_state)
	ioutil.WriteFile(AppConfig.TauntBotConf.StateFile, updated_state, 0644)
}
