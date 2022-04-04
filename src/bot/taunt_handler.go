package main

import (
	"time"
	"math/rand"
)

type TauntGenerator struct {
	grammar *Grammar
}

func newTauntGenerator(g* Grammar) *TauntGenerator {
	this := TauntGenerator{ grammar: g }
	rand.Seed(time.Now().UnixNano())
	return &this
}

func (this *TauntGenerator)process(command Command) Action {
	return Action { action: ActionTextReply, text: this.grammar.Taunt(command.text) }
}

func (this *TauntGenerator)processInline(command Command) []InlineQueryResponse {
	count := 3
	results := make([]InlineQueryResponse, count)
	for i := 0; i < count; i++ {
		nextTaunt := this.grammar.Taunt(command.text)
		results[i].title = nextTaunt
		results[i].value = nextTaunt
		results[i].cacheTime = 3
	}
	results[count - 1].title = "random taunt"
	return results
}
