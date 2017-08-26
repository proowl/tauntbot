package main

import (
	"fmt"
	"time"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"telegram"
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

// TODO: don't answer to too old messages
func process_updates(bot_state* telegram.BotState) {
	updates, err := telegram.GetUpdates(&AppConfig.TauntBotConf, bot_state)
	if err != nil {
		fmt.Println("GetUpdates failed with error: %v", err)
		return
	}
	printed, _ := json.Marshal(updates)
	fmt.Println(string(printed))
	if len(updates) > 0 {
		fmt.Printf("Received %d updates\n", len(updates))
		for _, msg := range updates {
			printed, _ := json.Marshal(msg)
			fmt.Println(string(printed))
			// todo: change to regex match instead full string match
			var response string
			switch msg.Message.Text {
				case "/shrug":
					response = "¯\\_(ツ)_/¯"
				// case "/taunt":
				default:
					response = Taunt(msg.Message.Text)
				// default:
				// 	fmt.Println(":: ignore")
			}
			if response != "" {
				fmt.Println("Sending response: ", response)
				err := telegram.SendMessage(&AppConfig.TauntBotConf, telegram.OutgoingMessage{ChatId: msg.Message.Chat.Id, Text: response, DisableNotification: true})
				if err != nil {
					fmt.Println("SendMessage failed with error: %v", err)
					return
				} else {
					bot_state.LastUpdateId = msg.Update_id
				}
			} else {
				bot_state.LastUpdateId = msg.Update_id
			}
		}
	}
}

func main() {
	// init
	rand.Seed(time.Now().UnixNano())
	fmt.Println(Taunt("hi bot"))
	fmt.Println(Taunt("the holy grail"))

	config_raw, _ := ioutil.ReadFile("etc/config.json")
	if err := json.Unmarshal(config_raw, &AppConfig); err != nil {
		panic(err)
	}
	restored_state_raw, _ := ioutil.ReadFile(AppConfig.TauntBotConf.StateFile)
	var bot_state telegram.BotState
	if err := json.Unmarshal(restored_state_raw, &bot_state); err != nil {
		bot_state.LastUpdateId = AppConfig.TauntBotConf.StartUpdateId
	}
	// process_updates();
	updated_state, _ := json.Marshal(bot_state)
	ioutil.WriteFile(AppConfig.TauntBotConf.StateFile, updated_state, 0644)
}
