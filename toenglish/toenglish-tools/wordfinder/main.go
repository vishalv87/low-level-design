package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// // Response structure for the dictionary API
// type DictionaryResponse struct {
// 	Word      string `json:"word"`
// 	Phonetic  string `json:"phonetic,omitempty"`
// 	Phonetics []struct {
// 		Text  string `json:"text,omitempty"`
// 		Audio string `json:"audio,omitempty"`
// 	} `json:"phonetics,omitempty"`
// 	Meanings []struct {
// 		PartOfSpeech string `json:"partOfSpeech,omitempty"`
// 		Definitions  []struct {
// 			Definition string   `json:"definition,omitempty"`
// 			Example    string   `json:"example,omitempty"`
// 			Synonyms   []string `json:"synonyms,omitempty"`
// 			Antonyms   []string `json:"antonyms,omitempty"`
// 		} `json:"definitions,omitempty"`
// 	} `json:"meanings,omitempty"`
// }

// // WordExists checks if a word exists in the English dictionary
// func WordExists(word string) (bool, error) {
// 	// Clean the word input
// 	word = strings.TrimSpace(strings.ToLower(word))

// 	// Create the API URL
// 	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)

// 	// Make the HTTP request
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return false, fmt.Errorf("error making HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	// 200 means the word exists, 404 means it doesn't
// 	if resp.StatusCode == http.StatusOK {
// 		return true, nil
// 	} else if resp.StatusCode == http.StatusNotFound {
// 		return false, nil
// 	} else {
// 		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}
// }

// func main() {
// 	// Check if a word argument was provided
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run dictionary.go <word>")
// 		os.Exit(1)
// 	}

// 	// Get the word from command line arguments
// 	word := os.Args[1]

// 	// Check if the word exists
// 	exists, err := WordExists(word)
// 	if err != nil {
// 		fmt.Printf("Error checking word: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Print the result
// 	if exists {
// 		fmt.Printf("The word '%s' exists in the English dictionary.\n", word)

// 		// Optionally get additional information about the word
// 		url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
// 		resp, err := http.Get(url)
// 		if err == nil && resp.StatusCode == http.StatusOK {
// 			defer resp.Body.Close()

// 			var dictResp []DictionaryResponse
// 			if err := json.NewDecoder(resp.Body).Decode(&dictResp); err == nil && len(dictResp) > 0 {
// 				fmt.Println("\nDefinitions:")
// 				for _, meaning := range dictResp[0].Meanings {
// 					fmt.Printf("  Part of speech: %s\n", meaning.PartOfSpeech)
// 					for i, def := range meaning.Definitions {
// 						if i >= 2 { // Limit to 2 definitions per part of speech
// 							break
// 						}
// 						fmt.Printf("    - %s\n", def.Definition)
// 						if def.Example != "" {
// 							fmt.Printf("      Example: %s\n", def.Example)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	} else {
// 		fmt.Printf("The word '%s' does not exist in the English dictionary.\n", word)
// 	}
// }
