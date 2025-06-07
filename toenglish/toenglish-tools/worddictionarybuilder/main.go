package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// // WordAnalysisPrompt represents the structured prompt for word analysis
// type WordAnalysisPrompt struct {
// 	PromptTitle         string          `json:"prompt_title"`
// 	PromptDescription   string          `json:"prompt_description"`
// 	OutputFormat        string          `json:"output_format"`
// 	RequestedAttributes []AttributeSpec `json:"requested_attributes"`
// 	WordToAnalyze       string          `json:"word_to_analyze"`
// }

// // AttributeSpec represents an attribute specification in the prompt
// type AttributeSpec struct {
// 	Name          string                   `json:"name"`
// 	Description   string                   `json:"description"`
// 	ExampleValue  interface{}              `json:"example_value,omitempty"`
// 	SubAttributes []map[string]interface{} `json:"sub_attributes,omitempty"`
// }

// // WordAnalysis represents the complete analysis structure
// type WordAnalysis struct {
// 	Word             string     `json:"word"`
// 	PartOfSpeech     string     `json:"part_of_speech"`
// 	PronunciationIPA string     `json:"pronunciation_ipa"`
// 	Syllabification  string     `json:"syllabification"`
// 	Definition       string     `json:"definition"`
// 	ExampleSentences []string   `json:"example_sentences"`
// 	Synonyms         []string   `json:"synonyms"`
// 	Antonyms         []string   `json:"antonyms"`
// 	Etymology        string     `json:"etymology"`
// 	Tags             []string   `json:"tags"`
// 	UsageNotes       UsageNotes `json:"usage_notes"`
// 	Frequency        string     `json:"frequency"`
// }

// // UsageNotes represents usage context information
// type UsageNotes struct {
// 	Collocations         []string `json:"collocations"`
// 	CulturalSignificance string   `json:"cultural_significance"`
// 	Register             string   `json:"register"`
// }

// // GeminiRequest represents the request structure for Gemini API
// type GeminiRequest struct {
// 	Contents []Content `json:"contents"`
// }

// // Content represents the content structure in Gemini request
// type Content struct {
// 	Parts []Part `json:"parts"`
// }

// // Part represents a text part in the content
// type Part struct {
// 	Text string `json:"text"`
// }

// // GeminiResponse represents the response from Gemini API
// type GeminiResponse struct {
// 	Candidates []Candidate `json:"candidates"`
// }

// // Candidate represents a response candidate
// type Candidate struct {
// 	Content Content `json:"content"`
// }

// // GeminiClient handles communication with Gemini API
// type GeminiClient struct {
// 	APIKey  string
// 	BaseURL string
// }

// // Improved function to clean JSON response from markdown formatting
// // cleanJSONResponse removes markdown code blocks and extracts pure JSON
// func cleanJSONResponse(response string) string {
// 	// Remove leading/trailing whitespace
// 	response = strings.TrimSpace(response)

// 	// Handle multiple possible markdown formats
// 	// Case 1: ```json ... ```
// 	// Case 2: ``` ... ```
// 	// Case 3: Plain text with JSON

// 	if strings.Contains(response, "```") {
// 		lines := strings.Split(response, "\n")
// 		var jsonLines []string
// 		inCodeBlock := false

// 		for _, line := range lines {
// 			trimmed := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmed, "```") {
// 				inCodeBlock = !inCodeBlock
// 				continue
// 			}
// 			if inCodeBlock {
// 				jsonLines = append(jsonLines, line)
// 			}
// 		}

// 		if len(jsonLines) > 0 {
// 			response = strings.Join(jsonLines, "\n")
// 		}
// 	}

// 	// Remove any remaining non-JSON text before the first {
// 	startIdx := strings.Index(response, "{")
// 	if startIdx == -1 {
// 		return response
// 	}

// 	// Find the matching closing brace for the complete JSON object
// 	braceCount := 0
// 	endIdx := -1
// 	inString := false
// 	escaped := false

// 	for i := startIdx; i < len(response); i++ {
// 		char := response[i]

// 		if escaped {
// 			escaped = false
// 			continue
// 		}

// 		if char == '\\' {
// 			escaped = true
// 			continue
// 		}

// 		if char == '"' {
// 			inString = !inString
// 			continue
// 		}

// 		if !inString {
// 			if char == '{' {
// 				braceCount++
// 			} else if char == '}' {
// 				braceCount--
// 				if braceCount == 0 {
// 					endIdx = i
// 					break
// 				}
// 			}
// 		}
// 	}

// 	if endIdx != -1 {
// 		response = response[startIdx : endIdx+1]
// 	}

// 	return strings.TrimSpace(response)
// }

// // NewGeminiClient creates a new Gemini API client
// func NewGeminiClient(apiKey string) *GeminiClient {
// 	return &GeminiClient{
// 		APIKey:  apiKey,
// 		BaseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
// 	}
// }

// // AnalyzeWord sends a request to Gemini API to analyze the given word
// func (gc *GeminiClient) AnalyzeWord(word string) (*WordAnalysis, error) {
// 	prompt := fmt.Sprintf(`Analyze the English word "%s" and return ONLY a JSON response in this exact format:

// {
//   "word": "%s",
//   "part_of_speech": "Noun/Verb/Adjective/etc",
//   "pronunciation_ipa": "/phonetic notation/",
//   "syllabification": "word-broken-into-syllables",
//   "definition": "clear dictionary definition",
//   "example_sentences": [
//     "sentence 1",
//     "sentence 2",
//     "sentence 3",
//     "sentence 4",
//     "sentence 5"
//   ],
//   "synonyms": ["synonym1", "synonym2"],
//   "antonyms": ["antonym1", "antonym2"],
//   "etymology": "word origin explanation",
//   "tags": ["category1", "category2", "category3"],
//   "usage_notes": {
//     "collocations": ["common phrase 1", "common phrase 2"],
//     "cultural_significance": "cultural info or N/A",
//     "register": "Formal/Informal/Standard/etc"
//   },
//   "frequency": "Very Common/Common/Less Common/Rare"
// }

// Return ONLY the JSON, no other text, no code blocks, no explanations.`, word, word)

// 	// Create request payload
// 	reqPayload := GeminiRequest{
// 		Contents: []Content{
// 			{
// 				Parts: []Part{
// 					{Text: prompt},
// 				},
// 			},
// 		},
// 	}

// 	// Marshal request to JSON
// 	jsonData, err := json.Marshal(reqPayload)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling request: %v", err)
// 	}

// 	// Create HTTP request
// 	url := fmt.Sprintf("%s?key=%s", gc.BaseURL, gc.APIKey)
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	// Send request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("error sending request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read response
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response: %v", err)
// 	}

// 	// Check for HTTP errors
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
// 	}

// 	// Parse Gemini response
// 	var geminiResp GeminiResponse
// 	if err := json.Unmarshal(body, &geminiResp); err != nil {
// 		return nil, fmt.Errorf("error parsing Gemini response: %v", err)
// 	}

// 	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
// 		return nil, fmt.Errorf("no content in Gemini response")
// 	}

// 	// Extract the JSON response text
// 	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

// 	// Add debugging to see the raw response
// 	fmt.Printf("Raw Gemini Response:\n%s\n", responseText)
// 	fmt.Println(strings.Repeat("=", 50))

// 	// Clean the response text to extract pure JSON
// 	// Remove markdown code blocks and other formatting
// 	responseText = cleanJSONResponse(responseText)

// 	// Add debugging to see the cleaned response
// 	fmt.Printf("Cleaned JSON Response:\n%s\n", responseText)
// 	fmt.Println(strings.Repeat("=", 50))

// 	// Parse the word analysis JSON
// 	var analysis WordAnalysis
// 	if err := json.Unmarshal([]byte(responseText), &analysis); err != nil {
// 		return nil, fmt.Errorf("error parsing word analysis JSON: %v\nResponse text: %s", err, responseText)
// 	}

// 	return &analysis, nil
// }

// // PrettyPrint prints the word analysis in a formatted way
// func PrettyPrint(analysis *WordAnalysis) {
// 	fmt.Printf("=== Word Analysis for '%s' ===\n\n", analysis.Word)
// 	fmt.Printf("Part of Speech: %s\n", analysis.PartOfSpeech)
// 	fmt.Printf("Pronunciation (IPA): %s\n", analysis.PronunciationIPA)
// 	fmt.Printf("Syllabification: %s\n", analysis.Syllabification)
// 	fmt.Printf("Definition: %s\n", analysis.Definition)
// 	fmt.Printf("Frequency: %s\n\n", analysis.Frequency)

// 	fmt.Println("Example Sentences:")
// 	for i, sentence := range analysis.ExampleSentences {
// 		fmt.Printf("  %d. %s\n", i+1, sentence)
// 	}

// 	fmt.Printf("\nSynonyms: %v\n", analysis.Synonyms)
// 	fmt.Printf("Antonyms: %v\n", analysis.Antonyms)
// 	fmt.Printf("Tags: %v\n\n", analysis.Tags)

// 	fmt.Printf("Etymology: %s\n\n", analysis.Etymology)

// 	fmt.Println("Usage Notes:")
// 	fmt.Printf("  Collocations: %v\n", analysis.UsageNotes.Collocations)
// 	fmt.Printf("  Cultural Significance: %s\n", analysis.UsageNotes.CulturalSignificance)
// 	fmt.Printf("  Register: %s\n", analysis.UsageNotes.Register)
// }

// func main() {
// 	apiKey := ""

// 	// Get word to analyze from command line argument
// 	word := "April" // Default word
// 	if len(os.Args) > 1 {
// 		word = os.Args[1]
// 	}

// 	// Create Gemini client
// 	client := NewGeminiClient(apiKey)

// 	fmt.Printf("Analyzing word: %s\n", word)
// 	fmt.Println("Sending request to Gemini API...")

// 	// Analyze the word
// 	analysis, err := client.AnalyzeWord(word)
// 	if err != nil {
// 		fmt.Printf("Error analyzing word: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Print results
// 	fmt.Println("\n" + strings.Repeat("=", 50))
// 	PrettyPrint(analysis)

// 	// Also output raw JSON
// 	fmt.Println("\n" + strings.Repeat("=", 50))
// 	fmt.Println("Raw JSON Output:")
// 	jsonOutput, err := json.MarshalIndent(analysis, "", "  ")
// 	if err != nil {
// 		fmt.Printf("Error formatting JSON: %v\n", err)
// 	} else {
// 		fmt.Println(string(jsonOutput))
// 	}
// }
