package main

import (
	"fmt"
	"math/rand"
	"strings"
	"regexp"
	"io/ioutil"
	"bufio"
	"os"
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


type Token struct {
	name string
	is_token bool
}

type Grammar struct {
	lang string
	rules map[string][][]Token
}

type GrammarRules struct {
	rules map[string]Grammar
}

func (gr *GrammarRules) FindGrammar(lang string) *Grammar {
	if g, found := gr.rules[lang]; !found {
		panic(fmt.Sprintf("Unknown lang %v", lang))
	} else {
		return &g
	}
}

func (g *Grammar) taunt(name string) string {
	paths := g.rules[name]
	if len(paths) == 0 {
		return ""
	}

	rand_num := rand.Int()
	path := paths[rand_num % len(paths)]

	var result string
	for _, token := range path {
		if token.is_token {
			result += g.taunt(token.name)
		} else {
			result += token.name
		}
	}
	return result
}

func (g* Grammar) Taunt(input string) string {
	phrase := strings.ToLower(input)
	if (holy_grail_mentioned(phrase)) {
		return g.holy_grail()
	}
	return g.taunt("<taunt>")
}

func (g* Grammar) holy_grail() string {
	if g.lang == "eng" {
		return "(A childish hand gesture)."
	}
	return "!!23li7s8i32u!12!!!1"
}

func parse_grammar(filename string) (*Grammar, error) {
	var grammar Grammar
	grammar.rules = make(map[string][][]Token)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error parsing file '%s': %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		splitted := strings.Split(line, " ::= ")
		from := splitted[0]
		raw_to := splitted[1]
		for _, val := range strings.Split(raw_to, " | ") {
			var to []Token
			tmp := ""
			started := false
			for _, c := range val {
				if c == '<' {
					if started {
						to = append(to, Token{tmp, false})
					}
					tmp = ""
					started = true
				} else if c == '>' {
					to = append(to, Token{fmt.Sprintf("<%s>", tmp), true})
					tmp = ""
					started = true
				} else if c == ' ' && started {
					tmp = tmp + string(c)
				} else {
					started = true
					tmp = tmp + string(c)
				}
			}
			if tmp != "" {
				to = append(to, Token{tmp, false})
			}
			grammar.rules[from] = append(grammar.rules[from], to)
		}
	}
	return &grammar, nil
}

func LoadLangs(folder string) (*GrammarRules, error) {
	var grammars GrammarRules
	grammars.rules = make(map[string]Grammar)

	grammar_file_re := regexp.MustCompile(`rules\.([a-zA-Z]+)$`)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		groups := grammar_file_re.FindStringSubmatch(file.Name())
		if len(groups) > 0 {
			lang := groups[1]
			new_grammar, err := parse_grammar(folder + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			new_grammar.lang = lang
			grammars.rules[lang] = *new_grammar
		}
	}
	return &grammars, nil
}
