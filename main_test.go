package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestEncoding(t *testing.T) {
	inputFile, err := os.Open("input.txt")
	outputFile, err := os.OpenFile("encoded.txt", os.O_APPEND, 0644)
	defer func() {
		err = inputFile.Close()
		if err != nil {
			log.Fatalf("File close error: %s", err)
		}

		err = outputFile.Close()
		if err != nil {
			log.Fatalf("File close error: %s", err)
		}
	}()

	if err != nil {
		log.Fatal("Can`t open file")
	}

	reader := bufio.NewReader(inputFile)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Введите ключевое слов: ")
	keyWord, keyWordString := []int{1, 2, 3}, "123"
	fmt.Println("Введите алфавит (по надобности): ")
	alphabet, exists := "", false

	writeToResult(outputFile, keyWordString+"\n")

	if !exists {
		alphabet = defaultAlphabet
	}

	writeToResult(outputFile, alphabet+"\n")

	encoder := EncoderGronsfeld{keyWord: keyWord}
	encoder.SetAlphabet(alphabet)

	input, count := getNextMessagePart(reader)
	for count > 0 {
		encoded := encoder.Encode(input)
		writeToResult(outputFile, string(encoded))
		input, count = getNextMessagePart(reader)
	}
}
