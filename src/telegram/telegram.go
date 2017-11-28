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
	Is_bot bool
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

type InlineQuery struct {
	Query_id string `json:"id"`
	From User
	Query string
	Offset string
}

type Update struct {
	Update_id int64
	Message Message
	InlineQuery InlineQuery `json:"inline_query"`
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

type InputMessageContent struct {
	MessageText string `json:"message_text"`
}

type InlineQueryResultArticle struct {
	Type string `json:"type"`
	Id int `json:"id"`
	Title string `json:"title"`
	HideUrl bool `json:"hide_url"`
	InputMessageContent InputMessageContent `json:"input_message_content"`
}

type OutgoingInlineQuery struct {
	QueryId string `json:"inline_query_id"`
	Results []InlineQueryResultArticle `json:"results"`
	CacheTime int `json:"cache_time"`
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

func SendInlineQueryResults(config* BotConfig, message OutgoingInlineQuery) error {
	data, _ := json.Marshal(message)
	fmt.Println(string(data))
	resp, err := http.Post(fmt.Sprintf("%s/bot%s/answerInlineQuery", config.Host, config.ApiToken), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("POST /answerInlineQuery: %v", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response: " + string(body))
	}
	return nil
}

func RestoreBotState(filename string, default_start_id int64) BotState {
	restored_state_raw, _ := ioutil.ReadFile(filename)
	var bot_state BotState
	if err := json.Unmarshal(restored_state_raw, &bot_state); err != nil {
		bot_state.LastUpdateId = default_start_id
	}
	return bot_state
}