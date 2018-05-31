package main

import (
	"fmt"
	"log"
	"time"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"telegram"
	"strings"
	"regexp"
	"os"
)

type Config struct {
	TauntBotConf telegram.BotConfig
	ReminderBotConf telegram.BotConfig
	SilentProcessing bool
}

var AppConfig Config

type CommandType uint8

const (
	TauntCommand = iota + 1
	ShrugCommand
)

func taunt_bot_iter(grammars* GrammarRules, bot_state* telegram.BotState) {
	updates, err := telegram.GetUpdates(&AppConfig.TauntBotConf, bot_state)
	if err != nil {
		log.Printf("GetUpdates failed with error: %v", err)
		return
	}
	if len(updates) > 0 {
		log.Printf("Received %d updates", len(updates))
		for _, msg := range updates {
			printed, _ := json.Marshal(msg)
			fmt.Println(string(printed))
			if msg.Message.Message_id > 0 {
				var response string
				split_re := regexp.MustCompile("[\\s\\.,!?@]")
				splitted := split_re.Split(msg.Message.Text, -1)
				for _, command := range splitted {
					switch command {
						case "/shrug":
							response = "¯\\_(ツ)_/¯"
							break;
						case "/taunt":
							response = grammars.Taunt("eng", msg.Message.Text)
							break;
						// default:
						// 	response = grammars.Taunt("eng", msg.Message.Text)
					}
				}
				if response != "" {
					log.Printf("Sending response: %v to %v", response, msg.Message.Chat.Id)
					if !AppConfig.SilentProcessing {
						err := telegram.SendMessage(&AppConfig.TauntBotConf, telegram.OutgoingMessage{ChatId: msg.Message.Chat.Id, Text: response, DisableNotification: true})
						if err != nil {
							log.Printf("SendMessage failed with error: %v", err)
						}
					}
				}
			} else if msg.InlineQuery.Query_id != "" {
				var results []telegram.InlineQueryResultArticle
				query := strings.TrimSpace(strings.ToLower(msg.InlineQuery.Query))
				cache_time := 300
				if strings.Contains(query, "taunt") {
					for i := 0; i < 5; i++ {
						taunt := grammars.Taunt("eng", msg.Message.Text)
						results = append(results, telegram.InlineQueryResultArticle {
							Id: 3000 + i,
							Type: "article",
							Title: taunt,
							HideUrl: true,
							InputMessageContent: telegram.InputMessageContent { MessageText: taunt },
						})
					}
					cache_time = 5
				} else if strings.Contains(query, "rand") {
					taunt := grammars.Taunt("eng", msg.Message.Text)
					results = append(results, telegram.InlineQueryResultArticle {
						Id: 2000,
						Type: "article",
						Title: "random taunt",
						HideUrl: true,
						InputMessageContent: telegram.InputMessageContent { MessageText: taunt },
					})
					cache_time = 0
				} else {
					smile := "¯\\_(ツ)_/¯"
					results = append(results, telegram.InlineQueryResultArticle {
						Id: 1000,
						Type: "article",
						Title: smile,
						HideUrl: true,
						InputMessageContent: telegram.InputMessageContent { MessageText: smile },
					})
				}
				if !AppConfig.SilentProcessing {
					err := telegram.SendInlineQueryResults(&AppConfig.TauntBotConf,
						telegram.OutgoingInlineQuery {
							QueryId: msg.InlineQuery.Query_id,
							CacheTime: cache_time,
							Results: results,
						})
					if err != nil {
						log.Printf("SendInlineQueryResults failed with error: %v", err)
					}
				}
			}
			bot_state.LastUpdateId = msg.Update_id
		}
	}
}

func run_taunt_bot(grammars* GrammarRules, bot_state* telegram.BotState) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<- ticker.C
		taunt_bot_iter(grammars, bot_state)
		taunt_updated_state, _ := json.Marshal(bot_state)
		ioutil.WriteFile(AppConfig.TauntBotConf.StateFile, taunt_updated_state, 0644)
	}
}

func addNewReminder(input string) string {
	return ""
}

func listAllReminders(input string) string {
	return ""
}

func removeReminder(input string) string {
	return ""
}

// ./bot etc/config.json
func main() {
	// init
	grammars := LoadLangs("etc/taunt")
	rand.Seed(time.Now().UnixNano())
	config_raw, _ := ioutil.ReadFile(os.Args[1])
	if err := json.Unmarshal(config_raw, &AppConfig); err != nil {
		panic(err)
	}
	taunt_bot_state := telegram.RestoreBotState(AppConfig.TauntBotConf.StateFile, AppConfig.TauntBotConf.StartUpdateId)

	run_taunt_bot(&grammars, &taunt_bot_state)
}
