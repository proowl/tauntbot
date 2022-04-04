package main

import (
	"log"
	"time"
	"io/ioutil"
	"encoding/json"
	"telegram"
	"strings"
	"regexp"
)

type Config struct {
	BotConf telegram.BotConfig
	SilentProcessing bool
	GrammarsPath string
	Debug bool
	SmileysConf string
}

type CommandType uint8

const (
	TauntCommand = iota + 1
	ShrugCommand
)

type Command struct {
	command string
	text string
}

type ActionType int

const (
	ActionIgnore ActionType = iota
	ActionTextReply
)

type Action struct {
	action ActionType
	text string
}

type InlineQueryResponse struct {
	title string
	cacheTime int
	value string
}

type HandlerInterface interface {
	process(command Command) Action
	processInline(command Command) []InlineQueryResponse
}

type Bot struct {
	handlers map[string]HandlerInterface
	state telegram.BotState
	config *Config
	commandRegexp *regexp.Regexp
}

func newBot(appConfig *Config) *Bot {
	var this Bot
	this.config = appConfig
	this.commandRegexp = regexp.MustCompile(`/(\w+)[^\w]*`)

	grammars, err := LoadLangs(appConfig.GrammarsPath)
	if err != nil {
		panic(err)
	}
	this.handlers = make(map[string]HandlerInterface)
	this.handlers["taunt"] = newTauntGenerator(grammars.FindGrammar("eng"))
	this.handlers["tauntru"] = newTauntGenerator(grammars.FindGrammar("ru"))
	this.handlers["shrug"] = newAsciiSmileys(appConfig.SmileysConf)
	// todo: add more ascii handlers

	this.state = telegram.RestoreBotState(appConfig.BotConf.StateFile, appConfig.BotConf.StartUpdateId)
	if appConfig.Debug {
		log.Printf("%+v", grammars.FindGrammar("ru"))
		log.Printf("'''%s'''", grammars.FindGrammar("ru").Taunt("<taunt>"))
	}
	return &this
}

func (this *Bot)start() {
	ticker := time.NewTicker(1 * time.Second)
	log.Print("Listening...")
	for {
		<- ticker.C
		this.iter()
		this.dumpState()
	}
}

func (this *Bot)dumpState() {
	updatedState, _ := json.Marshal(this.state)
	if err := ioutil.WriteFile(this.config.BotConf.StateFile, updatedState, 0644); err != nil {
		log.Printf("Error writing state: %v", err)
	}
}

func (this *Bot)iter() {
	updates, err := telegram.GetUpdates(&this.config.BotConf, &this.state)
	if err != nil {
		log.Printf("GetUpdates failed with error: %v", err)
		return
	}
	for _, msg := range updates {
		if this.config.Debug {
			printed, _ := json.Marshal(msg)
			log.Printf("Processing message: %s", string(printed))
		}
		if msg.Message.Message_id > 0 {
			this.processAsyncCommand(&msg)
		} else if msg.InlineQuery.Query_id != "" {
			this.processInlineQuery(&msg)
		}
		this.state.LastUpdateId = msg.Update_id
	}
}

func (this *Bot)processAsyncCommand(update *telegram.Update) {
	for _, matched := range this.commandRegexp.FindAllStringSubmatch(update.Message.Text, -1) {
		command := matched[1]
		if handler,found := this.handlers[command]; found {
			this.doAction( handler.process( Command{ command: command, text: update.Message.Text } ), update )
		}
	}
}

func (this *Bot)processInlineQuery(update *telegram.Update) {
	var results []telegram.InlineQueryResultArticle
	query := strings.TrimSpace(strings.ToLower(update.InlineQuery.Query))
	cache_time := 300
	id := 3000
	for _, handler := range this.handlers {
		for _, inlineResponse := range handler.processInline(Command{command: "taunt", text: query}) {
			id = id + 1
			results = append(results, telegram.InlineQueryResultArticle {
				Id: id,
				Type: "article",
				Title: inlineResponse.title,
				HideUrl: true,
				InputMessageContent: telegram.InputMessageContent { MessageText: inlineResponse.value },
			})
			if inlineResponse.cacheTime < cache_time {
				cache_time = inlineResponse.cacheTime
			}
		}
	}

	if this.config.SilentProcessing {
		return
	}
	if this.config.Debug {
		log.Print("Sending inline query response: %v", results)
	}
	err := telegram.SendInlineQueryResults(&this.config.BotConf,
		telegram.OutgoingInlineQuery {
			QueryId: update.InlineQuery.Query_id,
			CacheTime: cache_time,
			Results: results,
		})
	if err != nil {
		log.Printf("SendInlineQueryResults failed with error: %v", err)
	}
}

func (this *Bot)doAction(action Action, update *telegram.Update) {
	switch action.action {
	case ActionIgnore:
		// do nothing
	case ActionTextReply:
		if this.config.Debug || this.config.SilentProcessing {
			log.Printf("Sending response: %v to %v", action.text, update.Message.Chat.Id)
		}
		if this.config.SilentProcessing {
			break;
		}
		err := telegram.SendMessage(&this.config.BotConf, telegram.OutgoingMessage{
			ChatId: update.Message.Chat.Id,
			Text: action.text,
			DisableNotification: true,
		})
		if err != nil {
			log.Printf("SendMessage failed with error: %v", err)
		}
	}
}
