package lib

import (
	"io"
	"math/rand"
	"os"
	"strings"
)

func RandomWord() string {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	words := strings.Split(string(bytes), "\n")
	return words[rand.Int()%len(words)]
}
