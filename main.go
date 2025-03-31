package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sqweek/dialog" // Import sqweek/dialog for file selection
)

/*
Description:
This program processes text files by inserting newlines after Chinese punctuation and removing empty lines. The output is saved as a new file with a "_sc" suffix added to the original file name.

Features:
- Adds newlines after Chinese punctuation marks using regex.
- Removes empty lines from the content for cleanliness.
- Allows file selection via a GUI and writes processed content to an output file.

Workflow:
1. User selects an input file through a GUI using sqweek/dialog.
2. The program reads the file content and processes it by inserting newlines after Chinese punctuation.
3. Removes unnecessary empty lines from the processed content.
4. Saves cleaned content in a new file with a "_sc" suffix in the same directory.
5. Notifies the user after successful processing.
*/

func main() {
	// Step 1: Use sqweek/dialog to let the user select the input file
	inputFilePath, err := dialog.File().
		Filter("Text Files", "txt").
		Title("Select Input File").
		Load()
	if err != nil {
		if err == dialog.Cancelled {
			fmt.Println("File selection was cancelled.")
		} else {
			fmt.Println("Error selecting input file:", err)
		}
		return
	}

	// Display selected input file path
	fmt.Println("Selected input file:", inputFilePath)

	// Step 2: Construct output file path by appending the suffix '_sc' to the input file base name
	fileDir := filepath.Dir(inputFilePath)
	fileName := strings.TrimSuffix(filepath.Base(inputFilePath), filepath.Ext(inputFilePath))
	outputFilePath := filepath.Join(fileDir, fileName+"_sc"+filepath.Ext(inputFilePath))

	// Step 3: Open the input file for reading
	inputFileContent, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Step 4: Define regex to match Chinese punctuation marks
	punctuationRegex := regexp.MustCompile(`([，。？：！；、……——])`) // Captures Chinese punctuation marks

	// Step 5: Process the entire file content to insert a newline after each punctuation mark
	processedContent := punctuationRegex.ReplaceAllString(string(inputFileContent), "$1\n") // Replace punctuation with itself followed by a newline (actual newline)

	// Step 6: Remove empty lines from the processed content
	var cleanedLines []string
	scanner := bufio.NewScanner(strings.NewReader(processedContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Trim whitespace around each line
		if line != "" {                           // Exclude empty lines
			cleanedLines = append(cleanedLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error while removing empty lines:", err)
		return
	}

	// Combine cleaned lines into the final output
	cleanedContent := strings.Join(cleanedLines, "\n")

	// Step 7: Write cleaned content to the output file
	err = os.WriteFile(outputFilePath, []byte(cleanedContent), 0644)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	// Notify the user of successful processing
	fmt.Printf("Processed file with empty lines removed has been saved to: %s\n", outputFilePath)
}
