package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	bufReadCount    = 128
	defaultAlphabet = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ "
	fileInputPath   = "input1.txt"
	fileOutputPath  = "encoded.txt"
)

type CipherEncoder interface {
	Encode([]rune) []rune
}

type GronsfeldEncoder struct {
	alphabet map[rune]int
	keyWord  []int
}

func (encoder *GronsfeldEncoder) SetAlphabet(alphabet string) {
	encoder.alphabet = make(map[rune]int)
	runes := []rune(alphabet)
	for i, symbol := range runes {
		encoder.alphabet[symbol] = i
	}
}

func (encoder *GronsfeldEncoder) Encode(toEncode []rune) []rune {
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

func NewGronsfeldEncoder(keyword []int, alphabet string) *GronsfeldEncoder {
	encoder := &GronsfeldEncoder{keyWord: keyword}
	encoder.SetAlphabet(alphabet)
	return encoder
}

func getNextMessagePart(reader *bufio.Reader) ([]rune, int, error) {
	// `vsem privet, sevodnya ya nachinau svoi videoblog, znau tak ne interesno, y menya och vazhnaya informaciya!
	// Dorogoy dnevnic! Mne ne peredat slovami vse cho ya isputala v dannuyu secundu. Yuri ahuel nu vot chestnoe slovo`

	buf := make([]rune, bufReadCount)
	//bytesReadCount, err := reader.Read(buf)

	var readedRune rune
	var size int
	var err error

	for i := 0; i < bufReadCount; i++ {
		readedRune, size, err = reader.ReadRune()
		if err == io.EOF {
			buf = buf[:i]
			return buf, len(buf), err
		} else if err != nil {
			return nil, size, err
		}

		buf[i] = readedRune
	}

	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	//tempString := string(buf)
	//fmt.Println(tempString)
	//tempString = strings.Replace(tempString, "\r", "", -1)
	//tempString = strings.Replace(tempString, "\r", "", -1)
	////tempString = strings.Trim(tempString, string(utf8.RuneError))
	//buf = []byte(tempString)
	//buf = bytes.Trim(buf, "\x00")
	//runes := bytes.Runes(buf)

	return buf, size, err
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
	fmt.Println("Введите ключевое слов: ")
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
	fmt.Println("Введите алфавит (по надобности): ")
	sc.Scan()
	alphabet := sc.Text()

	if alphabet == "" {
		return alphabet, false
	}

	return alphabet, true
}

func GetInputFileReader(filePath string) *bufio.Reader {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, _ = os.Create(filePath)
	}

	inputFile, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Can`t open file: %s", filePath)
	}

	reader := bufio.NewReader(inputFile)

	return reader
}

func GetOutputFileWriter(filePath string) *os.File {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err = os.Create(filePath)
	}

	outputFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0777)

	if err != nil {
		log.Fatalf("Can`t open file: %s", filePath)
	}

	return outputFile
}

func startEncoding(encoder CipherEncoder, reader *bufio.Reader, writer *os.File) {
	input, count, err := getNextMessagePart(reader)
	for count > 0 {
		if err == io.EOF {
			encoded := encoder.Encode(input)
			writeToResult(writer, string(encoded))
			break
		} else if err != nil {
			break
		}
		encoded := encoder.Encode(input)
		writeToResult(writer, string(encoded))
		input, count, err = getNextMessagePart(reader)
	}
}

func main() {
	fileInputReader := GetInputFileReader(fileInputPath)
	fileOutputWriter := GetOutputFileWriter(fileOutputPath)

	defer func() {
		err := fileOutputWriter.Close()

		if err != nil {
			log.Printf("File <%s> close error: %s", fileOutputPath, err.Error())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	keyWord, keyWordString := getKeyWord(scanner)
	alphabet, exists := getAlphabet(scanner)

	if !exists {
		alphabet = defaultAlphabet
	}

	writeToResult(fileOutputWriter, keyWordString+"\n")
	writeToResult(fileOutputWriter, alphabet+"\n")
	encoder := NewGronsfeldEncoder(keyWord, alphabet)
	startEncoding(encoder, fileInputReader, fileOutputWriter)
}
