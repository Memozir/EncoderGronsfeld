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

	if err != nil {
		fmt.Println(err.Error())
	}

	keyWord, keyWordString := []int{1, 2, 3, 4, 6}, "123456"

	writeToResult(outputFile, keyWordString+"\n")

	alphabet := defaultAlphabet

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
