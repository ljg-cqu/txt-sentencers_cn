package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sqweek/dialog"
)

/*
Description:
<light>This program extracts, cleans, and processes sentences from text files, separating Chinese and English content and creating three output files.</light>

Features:
- Extracts Chinese sentences, English sentences, and combined sentences using regex.
- Cleans content by splitting after punctuation, removing empty lines, and discarding punctuation-only lines.
- Outputs cleaned data in dedicated files: Chinese, English, and combined sentences.
- User-friendly file selection through GUI dialog.

Workflow:
1. Select input file through a GUI.
2. Read and categorize Chinese and English sentences using regex patterns.
3. Clean content: split after punctuation, remove empty lines, discard punctuation-only lines.
4. Write cleaned and formatted data into three output files.
*/

// Main function
func main() {
	// File paths
	fmt.Println("Select the input file:")
	inputFile, err := dialog.File().
		Title("Select Input File").
		Filter("Text Files (*.txt)", "txt").
		Load()

	if err != nil {
		fmt.Printf("Error selecting input file: %v\n", err)
		return
	}
	if inputFile == "" {
		fmt.Println("No input file selected.")
		return
	}
	fmt.Printf("Selected input file: %s\n", inputFile)

	pureChineseSentencesFile := "pure_chinese_sentences.txt"
	pureEnglishSentencesFile := "pure_english_sentences.txt"
	combinedSentencesFile := "combined_sentences.txt"

	// Open the input file for reading
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	// `([︱|丨，,，。.?？/\\、：;；:——……！!])`
	// Regex patterns for filtering sentences
	chineseSentenceRegex := `[\p{Han}\d０-９。，！？：；（）【】《》“”‘’\-:.\s︱、\\]+` // Matches Chinese characters, Chinese/Arabic numbers, punctuation, and times
	englishSentenceRegex := `[a-zA-Z0-9.,!?;:'"()\-:\s|\\]+`            // Matches English sentences, numbers, and punctuation

	// Slices to store sentences
	chineseSentences := []string{}
	englishSentences := []string{}
	combinedSentences := []string{}

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Extract Chinese sentences
		chineseMatches := regexp.MustCompile(chineseSentenceRegex).FindAllString(line, -1)
		for _, sentence := range chineseMatches {
			chineseSentences = append(chineseSentences, sentence)
			combinedSentences = append(combinedSentences, sentence) // Include in combined output
		}

		// Extract English sentences
		englishMatches := regexp.MustCompile(englishSentenceRegex).FindAllString(line, -1)
		for _, sentence := range englishMatches {
			englishSentences = append(englishSentences, sentence)
			combinedSentences = append(combinedSentences, sentence) // Include in combined output
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Apply punctuation splitting, remove empty lines, and strip punctuation-only lines
	err = writeCleanedContent(removePunctuationOnlyLines(splitAfterPunctuation(joinLines(chineseSentences))), pureChineseSentencesFile)
	if err != nil {
		fmt.Printf("Error writing to Chinese sentences file: %v\n", err)
		return
	}

	err = writeCleanedContent(removePunctuationOnlyLines(splitAfterPunctuation(joinLines(englishSentences))), pureEnglishSentencesFile)
	if err != nil {
		fmt.Printf("Error writing to English sentences file: %v\n", err)
		return
	}

	err = writeCleanedContent(removePunctuationOnlyLines(splitAfterPunctuation(joinLines(combinedSentences))), combinedSentencesFile)
	if err != nil {
		fmt.Printf("Error writing to Combined sentences file: %v\n", err)
		return
	}

	fmt.Println("All output files written and cleaned successfully!")
}

// Function to write formatted and cleaned content to a file
func writeCleanedContent(content string, filePath string) error {
	// Clean content: remove empty lines and punctuation-only lines
	cleanContent := removeEmptyLines(content)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(cleanContent) // Write cleaned content
	if err != nil {
		return err
	}
	return writer.Flush()
}

// Function to split content after specific punctuation and insert newline
func splitAfterPunctuation(content string) string {
	// Define punctuation marks to split and insert a newline
	pattern := `([︱|丨，,，。.?？/\\、：;；:——……“"”！!])` // Matches a range of designated punctuation marks
	re := regexp.MustCompile(pattern)            // Compile regex pattern

	// Replace matched punctuation with itself followed by newline
	return re.ReplaceAllString(content, "$1\n")
}

// Helper function to join slices of strings into a single string
func joinLines(lines []string) string {
	return strings.Join(lines, "\n") // Join slices with a newline separator
}

// Function to remove empty lines from content
func removeEmptyLines(content string) string {
	lines := strings.Split(content, "\n") // Split content into lines
	var nonEmptyLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line) // Trim whitespace
		if trimmed != "" {                 // Keep non-empty lines
			nonEmptyLines = append(nonEmptyLines, trimmed)
		}
	}
	return strings.Join(nonEmptyLines, "\n") // Join cleaned lines back
}

// Additional function to remove lines containing only punctuation
func removePunctuationOnlyLines(content string) string {
	// Define regex for punctuation-only lines: both Chinese and English
	punctuationOnlyRegex := `^[.,!?;:'【】。、：；……——！丨︱-]+$`
	re := regexp.MustCompile(punctuationOnlyRegex)

	lines := strings.Split(content, "\n") // Split content into lines
	var nonPunctuationLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)             // Trim whitespace
		if trimmed != "" && !re.MatchString(trimmed) { // Remove punctuation-only lines
			nonPunctuationLines = append(nonPunctuationLines, trimmed)
		}
	}
	return strings.Join(nonPunctuationLines, "\n") // Join non-punctuation lines
}
