package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sajari/word2vec"
)

func main() {
	reader, err := os.Open("model.bin")
	if err != nil {
		log.Fatal(err)
	}
	model, err := word2vec.FromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	vecs := model.Map([]string{"愛"})
	for _, vec := range vecs {
		fmt.Printf("%v\n", vec)
	}
}
