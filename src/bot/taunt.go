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
	rules map[string][][]Token
}

type GrammarRules struct {
	rules map[string]Grammar
}

func (grammar Grammar) taunt(name string) string {
	paths := grammar.rules[name]

	rand_num := rand.Int()
	path := paths[rand_num % len(paths)]

	var result string
	for _, token := range path {
		if token.is_token {
			result += grammar.taunt(token.name)
		} else {
			result += token.name
		}
	}
	return result
}

func (grammars* GrammarRules) Taunt(lang string, input string) string {
	phrase := strings.ToLower(input)
	if (holy_grail_mentioned(phrase)) {
		return "(A childish hand gesture)."
	}
	return grammars.rules[lang].taunt("<taunt>")
}

func parse_grammar(filename string) *Grammar {
	var grammar Grammar
	grammar.rules = make(map[string][][]Token)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error parsing file '%s': %v\n", filename, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		splitted := strings.Split(line, " ::= ")
		from := splitted[0]
		raw_to := splitted[1]
		for _,val := range strings.Split(raw_to, " | ") {
			var istart,iend int
			istart = 0
			var to []Token
			for iend = range val {
				if val[iend] == '<' {
					if istart < iend {
						to = append(to, Token{val[istart:iend], false})
						istart = iend
					}
				} else if val[iend] == '>' {
					if istart < iend {
						to = append(to, Token{val[istart:iend+1], true})
						istart = iend+1
					}
				}
			}
			if istart <= iend {
				to = append(to, Token{val[istart:iend+1], false})
			}
			grammar.rules[from] = append(grammar.rules[from], to)
		}
	}
	return &grammar
}

func LoadLangs(folder string) GrammarRules {
	var grammars GrammarRules
	grammars.rules = make(map[string]Grammar)

	grammar_file_re := regexp.MustCompile(`rules\.([a-zA-Z]+)$`)
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		groups := grammar_file_re.FindStringSubmatch(file.Name())
		if len(groups) > 0 {
			lang := groups[1]
			new_grammar := parse_grammar(folder + "/" + file.Name())
			if new_grammar != nil {
				grammars.rules[lang] = *new_grammar
			}
		}
	}
	return grammars
}
