package telegram

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

type User struct {
	Id int64
	First_name string
	Username string
	Language_code string
}

type Room struct {
	Id int64
	First_name string
	Username string
	Type string
}

type Message struct {
	Message_id int64
	From User
	Chat Room
	Fate int64
	Text string
}

type Update struct {
	Update_id int64
	Message Message
}

type UpdateResponse struct {
	Ok bool
	Result []Update
}

type OutgoingMessage struct {
	ChatId int64 `json:"chat_id"`
	Text string `json:"text"`
	DisableNotification bool `json:"disable_notification"`
}

type BotConfig struct {
	ApiToken string
	Host string
	StartUpdateId int64
	StateFile string
}

type BotState struct {
	LastUpdateId int64
}

func GetUpdates(config* BotConfig, state* BotState) ([]Update, error) {
	resp, err := http.Get(fmt.Sprintf("%s/bot%s/getUpdates?offset=%d", config.Host, config.ApiToken, state.LastUpdateId + 1))
	if err != nil {
		return nil, fmt.Errorf("GET /getUpdates: %v", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var parsed_response UpdateResponse
		if err := json.Unmarshal(body, &parsed_response); err != nil {
			return nil, fmt.Errorf("unmarshal response: %v", err)
		}
		return parsed_response.Result, nil
	}
}

func SendMessage(config* BotConfig, message OutgoingMessage) error {
	data, _ := json.Marshal(message)
	resp, err := http.Post(fmt.Sprintf("%s/bot%s/sendMessage", config.Host, config.ApiToken), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("POST /sendMessage: %v", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response: " + string(body))
	}
	return nil
}
