package main

import (
	"bloomfilter"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("Hello World\n")
	err := Execute("cmd/wiki-100k.txt", 1000000, 10)
	fmt.Println(err)
}

func Execute(path string, n uint64, k int) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	wordCount := 0
	testPropotion := 10
	postive := 0
	negative := 0
	comments := 0

	bf := bloomfilter.NewBloomFilter(n, k)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		record := scanner.Text()
		if strings.HasPrefix(record, "#") {
			comments++
			continue
		}
		wordCount++
		if wordCount%testPropotion != 0 {
			err := bf.Add(record)
			if err != nil {
				fmt.Println("error while adding record: ", err)
			}
		} else {
			contains, err := bf.Contains(record)
			if err != nil {
				fmt.Println("error while testing record", err)
			}
			if contains {
				negative++
			} else {
				postive++
			}
		}
	}

	fmt.Println("Total words: ", wordCount)
	fmt.Println("Non Presence correct prediction: ", postive)
	fmt.Println("Non Presence wrong prediction: ", negative)
	fmt.Println("Total Comments: ", comments)
	return nil
}
