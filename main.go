package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	bufReadCount    = 128
	defaultAlphabet = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
)

type CipherEncoder interface {
	Encode([]rune) []rune
}

type EncoderGronsfeld struct {
	alphabet map[rune]int
	keyWord  []int
}

func (encoder *EncoderGronsfeld) SetAlphabet(alphabet string) {
	encoder.alphabet = make(map[rune]int)
	runes := []rune(alphabet)
	for i, symbol := range runes {
		encoder.alphabet[symbol] = i
	}
}

func (encoder *EncoderGronsfeld) Encode(toEncode []rune) []rune {
	//encoded := make([]rune, len(toEncode))
	keyWordCount := 0

	for position, symbol := range toEncode {

		if keyWordCount == len(encoder.keyWord) {
			keyWordCount = 0
		}

		symbolPosition, ok := encoder.alphabet[symbol]

		if !ok {
			continue
		}

		newPosition := symbolPosition + encoder.keyWord[keyWordCount]

		if newPosition > (len(encoder.alphabet) - 1) {
			newPosition = newPosition % (len(encoder.alphabet))
		}

		var encodedSymbol rune

		for symb, symbNumber := range encoder.alphabet {
			if symbNumber == newPosition {
				encodedSymbol = symb
			}
		}

		toEncode[position] = encodedSymbol
		keyWordCount++
	}

	return toEncode
}

func getNextMessagePart(reader *bufio.Reader) ([]rune, int) {
	// `vsem privet, sevodnya ya nachinau svoi videoblog, znau tak ne interesno, y menya och vazhnaya informaciya!
	// Dorogoy dnevnic! Mne ne peredat slovami vse cho ya isputala v dannuyu secundu. Yuri ahuel nu vot chestnoe slovo`

	buf := make([]byte, bufReadCount)
	bytesReadCount, _ := reader.Read(buf)
	buf = bytes.Trim(buf, "\x00")
	runes := bytes.Runes(buf)

	return runes, bytesReadCount
}

func writeToResult(file *os.File, encodedString string) int {
	bytesWriteCount, err := fmt.Fprint(file, encodedString)
	if err != nil {
		log.Printf("Write encoded part error: %s", err.Error())
		return bytesWriteCount
	}

	return bytesWriteCount
}

func getKeyWord(sc *bufio.Scanner) ([]int, string) {
	sc.Scan()
	input := sc.Text()
	keyWords := make([]int, len(input))
	strNumbers := strings.Split(input, "")
	for i, symbol := range strNumbers {
		number, _ := strconv.Atoi(symbol)
		keyWords[i] = number
	}

	return keyWords, input
}

func getAlphabet(sc *bufio.Scanner) (string, bool) {
	sc.Scan()
	alphabet := sc.Text()

	if alphabet == "" {
		return alphabet, false
	}

	return alphabet, true
}

func main() {
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
	scanner := bufio.NewScanner(os.Stdin)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Введите ключевое слов: ")
	keyWord, keyWordString := getKeyWord(scanner)
	fmt.Println("Введите алфавит (по надобности): ")
	alphabet, exists := getAlphabet(scanner)

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
