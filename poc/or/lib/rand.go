package lib

import (
	"io"
	"math/rand"
	"os"
	"runtime"
	"strings"
)

var cachedWords []string

func RandomWord() string {
	if cachedWords == nil {
		if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
			panic("unsupported OS for this operation")
		}

		file, err := os.Open("/usr/share/dict/words")
		if err != nil {
			panic(err)
		}

		words, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		cachedWords = strings.Split(string(words), "\n")
	}
	return cachedWords[rand.Int()%len(cachedWords)]
}
