package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"code.sajari.com/word2vec"
	"github.com/ikawaha/kagome/tokenizer"
)

type Song struct {
	Title  string    `json:"title"`
	Vector []float32 `json:"vector"`
}

func getBestMatchSong(vec word2vec.Vector, songs []Song) *Song {
	max := float32(0.0)
	maxIndex := 0
	for i, song := range songs {
		cos := vec.Dot(song.Vector)
		if cos > max {
			max = cos
			maxIndex = i
		}
	}
	return &(songs[maxIndex])
}

func aggregateVectors(vecs []word2vec.Vector) word2vec.Vector {
	vec := word2vec.Vector(make([]float32, len(vecs[0])))
	for _, v := range vecs {
		vec.Add(1.0, v)
	}
	vec.Normalise()
	return vec
}

func getWordVectors(model *word2vec.Model, words []string) ([]word2vec.Vector, error) {
	vecmap := model.Map(words)
	vecs := make([]word2vec.Vector, 0, len(vecmap))
	for _, v := range vecmap {
		vecs = append(vecs, v)
	}
	return vecs, nil
}

func wordRegularize(text string) []string {
	const (
		TYPE      = 0
		BASE_FORM = 6
	)
	const FEATURE_SIZE = 9

	unsupported_types := map[string]bool{
		"助詞":  true,
		"記号":  true,
		"助動詞": true,
	}

	t := tokenizer.New()
	tokens := t.Tokenize(text)
	words := make([]string, 0, 100)
	for _, token := range tokens {
		features := token.Features()
		if len(features) < FEATURE_SIZE {
			continue
		}
		if _, ok := unsupported_types[features[TYPE]]; ok {
			continue
		}

		words = append(words, features[BASE_FORM])
	}
	return words
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <input text file>\n", os.Args[0])
		return
	}

	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	song_data, err := ioutil.ReadFile("song_vec.json")
	if err != nil {
		log.Fatal(err)
	}
	var songs []Song
	if err := json.Unmarshal(song_data, &songs); err != nil {
		log.Fatal(err)
	}

	reader, err := os.Open("model.bin")
	if err != nil {
		return nil, err
	}
	defer func() { os.Close(reader) }()
	model, err := word2vec.FromReader(reader)
	if err != nil {
		return nil, err
	}

	words := wordRegularize(string(text))
	vecs, err := getWordVectors(model, words)
	if err != nil {
		log.Fatal(err)
	}
	vec := aggregateVectors(vecs)

	fmt.Printf("%#v\n", len(vec))
	fmt.Printf("%#v\n", []float32(vec))

	best_song := getBestMatchSong(vec, songs)
	fmt.Printf("%s\n", best_song.Title)
}
