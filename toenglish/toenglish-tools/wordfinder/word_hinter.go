package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// DataMuseWord represents a single word response from DataMuse API
type DataMuseWord struct {
	Word  string   `json:"word"`
	Score int      `json:"score"`
	Tags  []string `json:"tags,omitempty"`
	Defs  []string `json:"defs,omitempty"`
}

// HintGenerator provides methods to generate hints for words
type HintGenerator struct {
	httpClient *http.Client
}

// NewHintGenerator creates a new hint generator with proper timeout
func NewHintGenerator() *HintGenerator {
	return &HintGenerator{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GenerateHints creates a list of hints for the given word
func (hg *HintGenerator) GenerateHints(word string, maxHints int) ([]string, error) {
	hints := []string{}
	word = strings.ToLower(word)

	// 1. Get definition hints
	defHints, err := hg.getDefinitionHints(word)
	if err == nil {
		hints = append(hints, defHints...)
	}

	// 2. Get synonym hints
	if len(hints) < maxHints {
		synHints, err := hg.getSynonymHints(word)
		if err == nil {
			hints = append(hints, synHints...)
		}
	}

	// 3. Get "sounds like" hints
	if len(hints) < maxHints {
		soundHints, err := hg.getSoundsLikeHints(word)
		if err == nil {
			hints = append(hints, soundHints...)
		}
	}

	// 4. Get "related words" hints
	if len(hints) < maxHints {
		relatedHints, err := hg.getRelatedWordHints(word)
		if err == nil {
			hints = append(hints, relatedHints...)
		}
	}

	// Add basic hints about word length and first/last letter for very difficult words
	if len(hints) < 2 {
		hints = append(hints, fmt.Sprintf("This word has %d letters", len(word)))
		if len(word) > 0 {
			hints = append(hints, fmt.Sprintf("The word starts with '%s'", string(word[0])))
		}
		if len(word) > 1 {
			hints = append(hints, fmt.Sprintf("The word ends with '%s'", string(word[len(word)-1])))
		}
	}

	// Limit the number of hints
	if len(hints) > maxHints {
		hints = hints[:maxHints]
	}

	return hints, nil
}

// getDefinitionHints gets definition-based hints
func (hg *HintGenerator) getDefinitionHints(word string) ([]string, error) {
	url := fmt.Sprintf("https://api.datamuse.com/words?sp=%s&md=d", word)

	resp, err := hg.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var words []DataMuseWord
	if err := json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, err
	}

	hints := []string{}
	wordLower := strings.ToLower(word)

	// Process definition hints
	for _, w := range words {
		if strings.ToLower(w.Word) == wordLower && len(w.Defs) > 0 {
			for _, def := range w.Defs {
				// Clean the definition (remove part of speech prefix like "n\t")
				parts := strings.SplitN(def, "\t", 2)
				if len(parts) == 2 {
					cleanDef := parts[1]
					// Check that the definition doesn't contain the original word
					if !strings.Contains(strings.ToLower(cleanDef), wordLower) {
						hints = append(hints, fmt.Sprintf("Definition: %s", cleanDef))
						if len(hints) >= 2 {
							break
						}
					}
				}
			}
			break
		}
	}

	return hints, nil
}

// getSynonymHints gets synonym-based hints
func (hg *HintGenerator) getSynonymHints(word string) ([]string, error) {
	url := fmt.Sprintf("https://api.datamuse.com/words?rel_syn=%s", word)

	resp, err := hg.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var words []DataMuseWord
	if err := json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, err
	}

	hints := []string{}

	// Get top 2 synonyms as hints
	for i, w := range words {
		if i >= 2 || len(hints) >= 2 {
			break
		}
		hints = append(hints, fmt.Sprintf("Similar to: %s", w.Word))
	}

	return hints, nil
}

// getSoundsLikeHints gets "sounds like" hints
func (hg *HintGenerator) getSoundsLikeHints(word string) ([]string, error) {
	url := fmt.Sprintf("https://api.datamuse.com/words?sl=%s", word)

	resp, err := hg.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var words []DataMuseWord
	if err := json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, err
	}

	hints := []string{}

	// Only add "sounds like" if it's not too similar
	for i, w := range words {
		if i >= 2 || len(hints) >= 1 {
			break
		}
		// Avoid returning words that are very close to original
		if len(w.Word) > 3 && w.Word != word {
			hints = append(hints, fmt.Sprintf("Sounds like: %s", w.Word))
		}
	}

	return hints, nil
}

// getRelatedWordHints gets related word hints
func (hg *HintGenerator) getRelatedWordHints(word string) ([]string, error) {
	url := fmt.Sprintf("https://api.datamuse.com/words?rel_trg=%s", word)

	resp, err := hg.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var words []DataMuseWord
	if err := json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, err
	}

	hints := []string{}

	// Get a couple related words
	for i, w := range words {
		if i >= 3 || len(hints) >= 2 {
			break
		}
		hints = append(hints, fmt.Sprintf("Think of: %s", w.Word))
	}

	return hints, nil
}

func main() {
	hintGen := NewHintGenerator()

	// Example word to generate hints for
	word := "contractor"

	// How many hints you want at max
	maxHints := 5

	hints, err := hintGen.GenerateHints(word, maxHints)
	if err != nil {
		log.Fatalf("Error generating hints: %v", err)
	}

	fmt.Printf("Hints for the word '%s':\n", word)
	for i, hint := range hints {
		fmt.Printf("%d. %s\n", i+1, hint)
	}
}
