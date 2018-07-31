package main

import (
)

type AsciiSmileys struct {
	smileys map[string]string
}

func newAsciiSmileys(configFile string) *AsciiSmileys {
	var this AsciiSmileys
	this.smileys = make(map[string]string)
	return &this
}

func (this *AsciiSmileys)process(command Command) Action {
	if value,found := this.smileys[command.command]; found {
		return Action{action: ActionTextReply, text: value}
	}
	return Action{action: ActionTextReply, text: "¯\\_(ツ)_/¯"}
}

func (this *AsciiSmileys)processInline(command Command) []InlineQueryResponse {
	return []InlineQueryResponse{
		InlineQueryResponse{title: "¯\\_(ツ)_/¯", value:"¯\\_(ツ)_/¯", cacheTime: 300},
		InlineQueryResponse{title: "( ͡° ͜ʖ ͡°)", value:"( ͡° ͜ʖ ͡°)", cacheTime: 300},
	}
}
