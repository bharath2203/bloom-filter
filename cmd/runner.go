package main

import (
	"bloomfilter"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var path string
	var k int
	var n uint64

	rootCmd := &cobra.Command{
		Use:   "bloom",
		Short: "Analyze a file using Bloom Filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return fmt.Errorf("Please provide a file path using the --path flag")
			}
			return Execute(path, n, k)
		},
	}

	rootCmd.Flags().StringVarP(&path, "path", "p", "cmd/wiki-100k.txt", "Path to the file to analyze")
	rootCmd.Flags().IntVarP(&k, "hash-functions", "k", 10, "Number of hash functions for the Bloom Filter")
	rootCmd.Flags().Uint64VarP(&n, "bloom-filter-bits-counts", "n", 1000000, "Number of bits used for the Bloom Filter data structure")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

func Execute(path string, n uint64, k int) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	wordCount := 0
	testPropotion := 10
	correctPredictionCount := 0
	wrongPredictionCount := 0
	comments := 0

	bf := bloomfilter.NewBloomFilter(n, k)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		record := scanner.Text()
		// The dataset contains few comments and they start with #.
		// Ignore them.
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

			// This word is not added in the Bloom Filter. If the prediction comes to be true, its a false positive.
			if contains {
				wrongPredictionCount++
			} else {
				correctPredictionCount++
			}
		}
	}

	fmt.Println("Total words: ", wordCount)
	fmt.Println("Correct prediction: ", correctPredictionCount)
	fmt.Println("Wrong prediction: ", wrongPredictionCount)
	fmt.Println("Total Comments: ", comments)
	return nil
}
