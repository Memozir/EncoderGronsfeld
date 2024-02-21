package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestEncoding(t *testing.T) {
	inputFilePath := "input.txt"
	outputFilePath := "encoded.txt"

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		_, _ = os.Create(inputFilePath)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		_, err = os.Create(outputFilePath)
	}

	inputFile, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("Can`t open file: %s", inputFilePath)
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, 0777)

	if err != nil {
		log.Fatalf("Can`t open file: %s", outputFilePath)
	}

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

	reader := bufio.NewReader(inputFile)
	// writer := bufio.NewWriter(outputFile)
	// scanner := bufio.NewScanner(os.Stdin)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Введите ключевое слов: ")
	// keyWord, keyWordString := getKeyWord(scanner)
	keyWord, keyWordString := []int{1}, "1"
	fmt.Println("Введите алфавит (по надобности): ")
	// alphabet, exists := getAlphabet(scanner)
	alphabet, exists := defaultAlphabet, true

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
		fmt.Println(encoded)
		writeToResult(outputFile, string(encoded))
		// fmt.Print(string(encoded))
		input, count = getNextMessagePart(reader)
	}
}
