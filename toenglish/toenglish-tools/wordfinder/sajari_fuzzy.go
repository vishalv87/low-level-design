package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/sajari/fuzzy"
// )

// func main() {
// 	// Check if a word argument was provided
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run dictionary_checker.go <word>")
// 		os.Exit(1)
// 	}

// 	// Get the word from command line arguments
// 	word := strings.TrimSpace(strings.ToLower(os.Args[1]))

// 	// Create a new fuzzy model
// 	model := fuzzy.NewModel()

// 	// Configure the model - you can adjust these parameters
// 	model.SetThreshold(1)
// 	model.SetDepth(2)

// 	// Train the model with a dictionary
// 	// The library comes with a small built-in dictionary
// 	// You can also load a custom dictionary
// 	model.Train(fuzzy.SampleEnglish())

// 	// Check if the word exists in the dictionary
// 	exists := model.SpellCheck(word)

// 	// Print the result
// 	if exists != "" {
// 		fmt.Printf("The word '%s' exists in the dictionary.\n", word)
// 	} else {
// 		fmt.Printf("The word '%s' does not exist in the dictionary.\n", word)

// 		// Get spelling suggestions
// 		suggestions := model.Suggestions(word, false)
// 		if len(suggestions) > 0 {
// 			fmt.Println("\nDid you mean:")
// 			for i, suggestion := range suggestions {
// 				if i >= 5 { // Limit to 5 suggestions
// 					break
// 				}
// 				fmt.Printf("  - %s\n", suggestion)
// 			}
// 		}
// 	}
// }
