package main

import (
	"time"
	"math/rand"
)

type TauntGenerator struct {
	grammar GrammarRules
}

func newTauntGenerator(grammar* GrammarRules) *TauntGenerator {
	this := TauntGenerator{ grammar: *grammar }
	rand.Seed(time.Now().UnixNano())
	return &this
}

func (this *TauntGenerator)process(command Command) Action {
	return Action { action: ActionTextReply, text: this.grammar.Taunt("eng", command.text) }
}

func (this *TauntGenerator)processInline(command Command) []InlineQueryResponse {
	count := 3
	results := make([]InlineQueryResponse, count)
	for i := 0; i < count; i++ {
		nextTaunt := this.grammar.Taunt("eng", command.text)
		results[i].title = nextTaunt
		results[i].value = nextTaunt
		results[i].cacheTime = 3
	}
	results[count - 1].title = "random taunt"
	return results
}
