package main

import (
	"math/rand"
	"strings"
	"regexp"
)

func holy_grail_mentioned(input string) bool {
	split_re := regexp.MustCompile("[\\s\\.,!?]")
	splitted := split_re.Split(input, -1)
	var result []string
	for i := range splitted {
		if splitted[i] != "" {
			result = append(result, splitted[i])
		}
	}
	return strings.Contains(strings.Join(splitted, ""), "theholygrail");
}

func sentence() string {
	rand_num := rand.Int()
	switch (rand_num % 3) {
	case 0:
		return past_rel() + " " + noun_phrase()
	case 1:
		return present_rel() + " " + noun_phrase()
	case 2:
		return past_rel() + " " + article() + " " + noun()
	}
	return ""
}
func noun_phrase() string {
	return article() + " " + modified_noun()
}
func modified_noun() string {
	rand_num := rand.Int()
	switch (rand_num % 2) {
	case 0:
		return noun()
	case 1:
		return modifier() + " " + noun()
	}
	return ""
}
func modifier() string {
	rand_num := rand.Int()
	switch (rand_num % 2) {
	case 0:
		return adjective()
	case 1:
		return adverb() + " " + adjective()
	}
	return ""
}
func past_rel() string {
	return "your " + past_person() + " " + past_verb()
}
func present_rel() string {
	return "your " + present_person() + " " + present_verb()
}
func present_person() string {
	rand_num := rand.Int()
	arr := []string{"steed", "king", "first-born"}
	return arr[rand_num % len(arr)]
}
func past_person() string {
	rand_num := rand.Int()
	arr := []string{"mother", "father", "grandmother", "grandfather", "godfather"}
	return arr[rand_num % len(arr)]
}
func present_verb() string {
	rand_num := rand.Int()
	arr := []string{"is", "masquerades as"}
	return arr[rand_num % len(arr)]
}
func past_verb() string {
	rand_num := rand.Int()
	arr := []string{"was", "personified"}
	return arr[rand_num % len(arr)]
}
func noun() string {
	rand_num := rand.Int()
	arr := []string{"hamster", "coconut", "duck", "herring", "newt", "peril", "chicken", "vole", "parrot", "mouse", "twit"}
	return arr[rand_num % len(arr)]
}
func article() string {
	return "a"
}
func adjective() string {
	rand_num := rand.Int()
	arr := []string{"silly", "wicked", "sordid", "naughty", "repulsive", "malodorous", "ill-tempered"}
	return arr[rand_num % len(arr)]
}
func adverb() string {
	rand_num := rand.Int()
	arr := []string{"conspicuously", "categorically", "positively", "cruelly", "incontrovertibly"}
	return arr[rand_num % len(arr)]
}
func taunt(depth int) string {
	rand_num := rand.Int()
	divide_by := 3
	if depth > 0 {
		divide_by = 2
	}
	switch (rand_num % divide_by) {
	case 0:
		return sentence()
	case 1:
		return noun() + "!"
	case 2:
		return taunt(depth + 1) + " and " + sentence()
	}
	return ""
}

func Taunt(input string) string {
	phrase := strings.ToLower(input)
	if (holy_grail_mentioned(phrase)) {
		return "(A childish hand gesture)."
	}
	return taunt(0);
}
