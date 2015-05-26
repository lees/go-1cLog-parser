package main

import (
	"io/ioutil"
	"strings"
	//"fmt"
)

func readStr(text string) (string, string) {

	pos := strings.IndexAny(text, string('"'))

	if pos == -1 {
		panic("not ending string")
	}

	if text[pos+1] == '"' {
		subrez, rest := readStr(text[pos+2:])
		return text[:pos] + subrez, rest
	}

	return text[:pos], text[pos+1:]

}

func Tokenize(text string, tokens chan string) {

	seps := "{}," + string('"')

	for {

		pos := strings.IndexAny(text, seps)
		if pos == -1 {
			close(tokens)
			break
		}

		if pos > 0 {
			str := strings.TrimSpace(text[:pos])
			if len(str) > 0 {
				tokens <- str
			}
		}

		switch text[pos] {
		case ',':
			text = text[pos+1:]
		case '"':
			str, rest := readStr(text[pos+1:])

			//fmt.Printf("\nstr:%s\nrest:%s",str,rest[:5])

			tokens <- string('"') + str + string('"')
			//tokens <- str
			text = rest
		default:
			tokens <- string(text[pos])
			text = text[pos+1:]
		}

	}

}

func Parse(tokens chan string, events chan []string) {

	event := make([]string, 0, 20)
	stringBuffer := ""

	brackets := 0

	for token := range tokens {

		switch token {
		case "{":
			switch {
			case brackets <= 0:
				event = make([]string, 0, 20)
				brackets = 0
			case brackets == 1:
				stringBuffer = "{"
			default:
				stringBuffer = stringBuffer + "{"
			}
			brackets++
		case "}":
			switch {
			case brackets == 1:
				events <- event
			case brackets == 2:
				stringBuffer = stringBuffer + "}"
				event = append(event, stringBuffer)
			default:
				stringBuffer = stringBuffer + "}"
			}
			brackets--
		default:
			if brackets > 1 {
				if len(stringBuffer) > 1 {
					stringBuffer = stringBuffer + ","
				}
				stringBuffer = stringBuffer + token
			} else {
				event = append(event, token)
			}

		}

	}
	close(events)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ConvertFile(filefrom string) chan []string {

	dat, err := ioutil.ReadFile(filefrom)
	check(err)
	tokens := make(chan string, 200)
	go Tokenize(string(dat), tokens)

	events := make(chan []string, 20)
	go Parse(tokens, events)

	return events
}
