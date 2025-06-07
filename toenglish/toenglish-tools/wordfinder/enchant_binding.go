package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/jdkato/prose/v2"
// )

// func main() {
// 	// Check if a word argument was provided
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run nlp_dictionary.go <word>")
// 		os.Exit(1)
// 	}

// 	// Get the word from command line arguments
// 	word := strings.TrimSpace(strings.ToLower(os.Args[1]))

// 	// Create a document with the word
// 	// The prose library will analyze the word
// 	doc, err := prose.NewDocument(word)
// 	if err != nil {
// 		fmt.Printf("Error analyzing word: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Get the tokens from the document
// 	tokens := doc.Tokens()
// 	if len(tokens) == 0 {
// 		fmt.Printf("The word '%s' could not be analyzed.\n", word)
// 		os.Exit(1)
// 	}

// 	// Check if the word is recognized and not marked as misspelled
// 	// Prose doesn't directly provide spell checking, but we can check
// 	// if the token is recognized as a valid word with a tag
// 	token := tokens[0]

// 	// Process the result
// 	// Note: This is a simplified check - prose is primarily a natural language processing
// 	// library, not a spelling checker, but it does recognize common words
// 	if token.Tag != "XX" && token.Tag != "" {
// 		fmt.Printf("The word '%s' appears to be valid (recognized as %s).\n", word, describeTag(token.Tag))
// 	} else {
// 		fmt.Printf("The word '%s' does not appear to be a recognized word.\n", word)
// 	}
// }

// // describeTag returns a description of the POS tag
// func describeTag(tag string) string {
// 	tags := map[string]string{
// 		"CC":   "coordinating conjunction",
// 		"CD":   "cardinal number",
// 		"DT":   "determiner",
// 		"EX":   "existential there",
// 		"FW":   "foreign word",
// 		"IN":   "preposition/subordinating conjunction",
// 		"JJ":   "adjective",
// 		"JJR":  "adjective, comparative",
// 		"JJS":  "adjective, superlative",
// 		"LS":   "list item marker",
// 		"MD":   "modal",
// 		"NN":   "noun, singular or mass",
// 		"NNS":  "noun, plural",
// 		"NNP":  "proper noun, singular",
// 		"NNPS": "proper noun, plural",
// 		"PDT":  "predeterminer",
// 		"POS":  "possessive ending",
// 		"PRP":  "personal pronoun",
// 		"PRP$": "possessive pronoun",
// 		"RB":   "adverb",
// 		"RBR":  "adverb, comparative",
// 		"RBS":  "adverb, superlative",
// 		"RP":   "particle",
// 		"SYM":  "symbol",
// 		"TO":   "to",
// 		"UH":   "interjection",
// 		"VB":   "verb, base form",
// 		"VBD":  "verb, past tense",
// 		"VBG":  "verb, gerund/present participle",
// 		"VBN":  "verb, past participle",
// 		"VBP":  "verb, non-3rd person singular present",
// 		"VBZ":  "verb, 3rd person singular present",
// 		"WDT":  "wh-determiner",
// 		"WP":   "wh-pronoun",
// 		"WP$":  "possessive wh-pronoun",
// 		"WRB":  "wh-adverb",
// 	}

// 	if desc, ok := tags[tag]; ok {
// 		return desc
// 	}
// 	return "unknown word type"
// }
