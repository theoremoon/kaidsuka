package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ikawaha/kagome/tokenizer"
)

const (
	TYPE      = 0
	BASE_FORM = 6
)
const FEATURE_SIZE = 9

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <input text file>\n", os.Args[0])
		return
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	unsupported_types := map[string]bool{
		"助詞":  true,
		"記号":  true,
		"助動詞": true,
	}

	t := tokenizer.New()
	tokens := t.Tokenize(string(content))
	for _, token := range tokens {
		features := token.Features()
		if len(features) < FEATURE_SIZE {
			continue
		}
		if _, ok := unsupported_types[features[TYPE]]; ok {
			continue
		}

		fmt.Printf("%s %s\n", features[TYPE], features[BASE_FORM])
	}

}
