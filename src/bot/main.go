package main

import (
	"fmt"
	"time"
	"math/rand"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"telegram"
	"bytes"
)

type Config struct {
	BotAPIToken string
	LastUpdateId int64
}

var AppConfig Config

type CommandType uint8

const (
	TauntCommand = iota + 1
	ShrugCommand
)

func Telegram_GetUpdates_HTTP() []telegram.Update {
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", AppConfig.BotAPIToken, AppConfig.LastUpdateId + 1))
	if err != nil {
		fmt.Println("Error on /getUpdates request:", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		var parsed_response telegram.UpdateResponse
		if err := json.Unmarshal(body, &parsed_response); err != nil {
			panic(err)
		}
		printed, _ := json.Marshal(parsed_response)
		fmt.Println(string(printed))
		fmt.Println("parsed_response: ", parsed_response)
		return parsed_response.Result
	}
	var updates []telegram.Update
	return updates
}

func Telegram_SendMessage(message telegram.OutgoingMessage) bool { 
	data, _ := json.Marshal(message)
	resp, err := http.Post("https://api.telegram.org/bot" + AppConfig.BotAPIToken + "/sendMessage", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error on /sendMessage request:", err)		
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response: " + string(body))
		return true
	}
	return false
}

func Telegram_GetUpdates_Test() []telegram.Update {
	var updates []telegram.Update
	var update telegram.Update
	update.Update_id = 1
	var msg telegram.Message
	msg.Message_id = 1
	msg.Text = "/shrug"
	update.Message = msg
	updates = append(updates, update)
	return updates
}

func Telegram_GetUpdates() []telegram.Update {
	// return Telegram_GetUpdates_Test()
	return Telegram_GetUpdates_HTTP()
}

func process_updates() {
	updates := Telegram_GetUpdates()
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
				ok := Telegram_SendMessage(telegram.OutgoingMessage{ChatId: msg.Message.Chat.Id, Text: response, DisableNotification: true})
				if !ok {
					return
				} else {
					AppConfig.LastUpdateId = msg.Update_id
				}
			}
		}
	}

}

func main() {
	// init
	rand.Seed(time.Now().UnixNano())
	config_raw, _ := ioutil.ReadFile("config.json")
	if err := json.Unmarshal(config_raw, &AppConfig); err != nil {
		panic(err)
	}
	fmt.Println(AppConfig)
	process_updates();
	updated_config, _ := json.Marshal(AppConfig)
	ioutil.WriteFile("config.json", updated_config, 0644)
}
