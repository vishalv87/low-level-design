package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	// chatgpt:change - Add MongoDB driver imports
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// chatgpt:change - Add MongoDB document structures
// VocabularyWord represents a word document from MongoDB
type VocabularyWord struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Word             string             `bson:"word" json:"word"`
	WordCategoryID   primitive.ObjectID `bson:"wordCategoryID" json:"wordCategoryID"`
	WordCategoryName string             `bson:"wordCategoryName" json:"wordCategoryName"`
	Difficulty       string             `bson:"difficulty" json:"difficulty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// VocabularyWordDetail represents the detailed analysis with original word reference
type VocabularyWordDetail struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	VocabularyWordID primitive.ObjectID `bson:"vocabularyWordID" json:"vocabularyWordID"` // chatgpt:change - Reference to original word document
	Word             string             `bson:"word" json:"word"`
	PartOfSpeech     string             `bson:"part_of_speech" json:"part_of_speech"`
	PronunciationIPA string             `bson:"pronunciation_ipa" json:"pronunciation_ipa"`
	Syllabification  string             `bson:"syllabification" json:"syllabification"`
	Definition       string             `bson:"definition" json:"definition"`
	ExampleSentences []string           `bson:"example_sentences" json:"example_sentences"`
	Synonyms         []string           `bson:"synonyms" json:"synonyms"`
	Antonyms         []string           `bson:"antonyms" json:"antonyms"`
	Etymology        string             `bson:"etymology" json:"etymology"`
	Tags             []string           `bson:"tags" json:"tags"`
	UsageNotes       UsageNotes         `bson:"usage_notes" json:"usage_notes"`
	Frequency        string             `bson:"frequency" json:"frequency"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// WordAnalysisPrompt represents the structured prompt for word analysis
type WordAnalysisPrompt struct {
	PromptTitle         string          `json:"prompt_title"`
	PromptDescription   string          `json:"prompt_description"`
	OutputFormat        string          `json:"output_format"`
	RequestedAttributes []AttributeSpec `json:"requested_attributes"`
	WordToAnalyze       string          `json:"word_to_analyze"`
}

// AttributeSpec represents an attribute specification in the prompt
type AttributeSpec struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	ExampleValue  interface{}              `json:"example_value,omitempty"`
	SubAttributes []map[string]interface{} `json:"sub_attributes,omitempty"`
}

// WordAnalysis represents the complete analysis structure
type WordAnalysis struct {
	Word             string     `json:"word"`
	PartOfSpeech     string     `json:"part_of_speech"`
	PronunciationIPA string     `json:"pronunciation_ipa"`
	Syllabification  string     `json:"syllabification"`
	Definition       string     `json:"definition"`
	ExampleSentences []string   `json:"example_sentences"`
	Synonyms         []string   `json:"synonyms"`
	Antonyms         []string   `json:"antonyms"`
	Etymology        string     `json:"etymology"`
	Tags             []string   `json:"tags"`
	UsageNotes       UsageNotes `json:"usage_notes"`
	Frequency        string     `json:"frequency"`
}

// UsageNotes represents usage context information
type UsageNotes struct {
	Collocations         []string `json:"collocations"`
	CulturalSignificance string   `json:"cultural_significance"`
	Register             string   `json:"register"`
}

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents the content structure in Gemini request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents a text part in the content
type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// Candidate represents a response candidate
type Candidate struct {
	Content Content `json:"content"`
}

// chatgpt:change - Add MongoDB client structure
// MongoDBClient handles MongoDB operations
type MongoDBClient struct {
	Client     *mongo.Client
	Database   *mongo.Database
	WordsCol   *mongo.Collection
	DetailsCol *mongo.Collection
}

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient(connectionString, dbName string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	database := client.Database(dbName)
	wordsCol := database.Collection("vocabularywords")
	detailsCol := database.Collection("vocabularyworddetails")

	return &MongoDBClient{
		Client:     client,
		Database:   database,
		WordsCol:   wordsCol,
		DetailsCol: detailsCol,
	}, nil
}

// GetAllVocabularyWords retrieves all vocabulary words from MongoDB
func (mc *MongoDBClient) GetAllVocabularyWords(ctx context.Context) ([]VocabularyWord, error) {
	cursor, err := mc.WordsCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find vocabulary words: %v", err)
	}
	defer cursor.Close(ctx)

	var words []VocabularyWord
	if err := cursor.All(ctx, &words); err != nil {
		return nil, fmt.Errorf("failed to decode vocabulary words: %v", err)
	}

	return words, nil
}

// SaveWordDetail saves a detailed word analysis to MongoDB
func (mc *MongoDBClient) SaveWordDetail(ctx context.Context, detail *VocabularyWordDetail) error {
	detail.CreatedAt = time.Now()
	detail.UpdatedAt = time.Now()

	_, err := mc.DetailsCol.InsertOne(ctx, detail)
	if err != nil {
		return fmt.Errorf("failed to save word detail: %v", err)
	}

	return nil
}

// CheckIfDetailExists checks if a word detail already exists
func (mc *MongoDBClient) CheckIfDetailExists(ctx context.Context, vocabularyWordID primitive.ObjectID) (bool, error) {
	count, err := mc.DetailsCol.CountDocuments(ctx, bson.M{"vocabularyWordID": vocabularyWordID})
	if err != nil {
		return false, fmt.Errorf("failed to check if detail exists: %v", err)
	}
	return count > 0, nil
}

// Close closes the MongoDB connection
func (mc *MongoDBClient) Close(ctx context.Context) error {
	return mc.Client.Disconnect(ctx)
}

// GeminiClient handles communication with Gemini API
type GeminiClient struct {
	APIKey  string
	BaseURL string
}

// chatgpt:change - Improved function to clean JSON response from markdown formatting
// cleanJSONResponse removes markdown code blocks and extracts pure JSON
func cleanJSONResponse(response string) string {
	// Remove leading/trailing whitespace
	response = strings.TrimSpace(response)

	// Handle multiple possible markdown formats
	// Case 1: ```json ... ```
	// Case 2: ``` ... ```
	// Case 3: Plain text with JSON

	if strings.Contains(response, "```") {
		lines := strings.Split(response, "\n")
		var jsonLines []string
		inCodeBlock := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}
			if inCodeBlock {
				jsonLines = append(jsonLines, line)
			}
		}

		if len(jsonLines) > 0 {
			response = strings.Join(jsonLines, "\n")
		}
	}

	// Remove any remaining non-JSON text before the first {
	startIdx := strings.Index(response, "{")
	if startIdx == -1 {
		return response
	}

	// Find the matching closing brace for the complete JSON object
	braceCount := 0
	endIdx := -1
	inString := false
	escaped := false

	for i := startIdx; i < len(response); i++ {
		char := response[i]

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if !inString {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					endIdx = i
					break
				}
			}
		}
	}

	if endIdx != -1 {
		response = response[startIdx : endIdx+1]
	}

	return strings.TrimSpace(response)
}

// NewGeminiClient creates a new Gemini API client
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		APIKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
	}
}

// AnalyzeWord sends a request to Gemini API to analyze the given word
func (gc *GeminiClient) AnalyzeWord(word string) (*VocabularyWordDetail, error) {
	// chatgpt:change - Simplify the prompt to be more direct and clear
	prompt := fmt.Sprintf(`Analyze the English word "%s" and return ONLY a JSON response in this exact format:

{
  "word": "%s",
  "part_of_speech": "Noun/Verb/Adjective/etc",
  "pronunciation_ipa": "/phonetic notation/",
  "syllabification": "word-broken-into-syllables",
  "definition": "clear dictionary definition",
  "example_sentences": [
    "sentence 1",
    "sentence 2", 
    "sentence 3",
    "sentence 4",
    "sentence 5"
  ],
  "synonyms": ["synonym1", "synonym2"],
  "antonyms": ["antonym1", "antonym2"],
  "etymology": "word origin explanation",
  "tags": ["category1", "category2", "category3"],
  "usage_notes": {
    "collocations": ["common phrase 1", "common phrase 2"],
    "cultural_significance": "cultural info or N/A",
    "register": "Formal/Informal/Standard/etc"
  },
  "frequency": "Very Common/Common/Less Common/Rare"
}

Return ONLY the JSON, no other text, no code blocks, no explanations.`, word, word)

	// Create request payload
	reqPayload := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s?key=%s", gc.BaseURL, gc.APIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Gemini response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("error parsing Gemini response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in Gemini response")
	}

	// Extract the JSON response text
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// chatgpt:change - Add debugging to see the raw response (can be removed in production)
	fmt.Printf("Raw Gemini Response for '%s':\n%s\n", word, responseText)
	fmt.Println(strings.Repeat("=", 50))

	// chatgpt:change - Clean the response text to extract pure JSON
	// Remove markdown code blocks and other formatting
	responseText = cleanJSONResponse(responseText)

	// chatgpt:change - Add debugging to see the cleaned response (can be removed in production)
	fmt.Printf("Cleaned JSON Response for '%s':\n%s\n", word, responseText)
	fmt.Println(strings.Repeat("=", 50))

	// Parse the word analysis JSON
	var analysis VocabularyWordDetail
	if err := json.Unmarshal([]byte(responseText), &analysis); err != nil {
		return nil, fmt.Errorf("error parsing word analysis JSON for '%s': %v\nResponse text: %s", word, err, responseText)
	}

	return &analysis, nil
}

// chatgpt:change - Add method to get all existing vocabulary word IDs
// GetExistingVocabularyWordIDs retrieves all vocabulary word IDs that already have details
func (mc *MongoDBClient) GetExistingVocabularyWordIDs(ctx context.Context) (map[primitive.ObjectID]bool, error) {
	cursor, err := mc.DetailsCol.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"vocabularyWordID": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find existing vocabulary word IDs: %v", err)
	}
	defer cursor.Close(ctx)

	existingIDs := make(map[primitive.ObjectID]bool)
	for cursor.Next(ctx) {
		var doc struct {
			VocabularyWordID primitive.ObjectID `bson:"vocabularyWordID"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode vocabulary word ID: %v", err)
		}
		existingIDs[doc.VocabularyWordID] = true
	}

	return existingIDs, nil
}

// chatgpt:change - Optimize ProcessAllWords to batch check existing documents
// ProcessAllWords processes all vocabulary words and saves detailed analysis
func ProcessAllWords(geminiClient *GeminiClient, mongoClient *MongoDBClient) error {
	ctx := context.Background()

	// Get all vocabulary words
	fmt.Println("Fetching all vocabulary words from MongoDB...")
	words, err := mongoClient.GetAllVocabularyWords(ctx)
	if err != nil {
		return fmt.Errorf("failed to get vocabulary words: %v", err)
	}

	fmt.Printf("Found %d vocabulary words to process\n", len(words))

	// chatgpt:change - Fetch all existing vocabulary word IDs upfront for efficiency
	fmt.Println("Checking for existing word details...")
	existingIDs, err := mongoClient.GetExistingVocabularyWordIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing vocabulary word IDs: %v", err)
	}
	fmt.Printf("Found %d existing word details\n", len(existingIDs))

	successCount := 0
	errorCount := 0
	skippedCount := 0

	for i, word := range words {
		fmt.Printf("\nProcessing word %d/%d: '%s'\n", i+1, len(words), word.Word)

		// chatgpt:change - Check if detail already exists using pre-fetched map
		if existingIDs[word.ID] {
			fmt.Printf("Detail already exists for '%s', skipping...\n", word.Word)
			skippedCount++
			continue // chatgpt:change - Skip to next iteration if document already exists
		}

		// Analyze the word with Gemini
		detail, err := geminiClient.AnalyzeWord(word.Word)
		if err != nil {
			fmt.Printf("Error analyzing word '%s': %v\n", word.Word, err)
			errorCount++
			continue
		}

		// Set the vocabulary word ID reference
		detail.VocabularyWordID = word.ID

		// Save to MongoDB
		if err := mongoClient.SaveWordDetail(ctx, detail); err != nil {
			fmt.Printf("Error saving detail for '%s': %v\n", word.Word, err)
			errorCount++
			continue
		}

		fmt.Printf("Successfully processed and saved detail for '%s'\n", word.Word)
		successCount++

		// chatgpt:change - Add the newly processed ID to the existing map to avoid duplicates in case of retry
		existingIDs[word.ID] = true

		// Add a small delay to avoid hitting rate limits
		time.Sleep(4 * time.Second)
	}

	fmt.Printf("\n=== Processing Summary ===\n")
	fmt.Printf("Total words: %d\n", len(words))
	fmt.Printf("Successfully processed: %d\n", successCount)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Skipped (already exists): %d\n", skippedCount)

	return nil
}

// PrettyPrint prints the word analysis in a formatted way
func PrettyPrint(analysis *VocabularyWordDetail) {
	fmt.Printf("=== Word Analysis for '%s' ===\n\n", analysis.Word)
	fmt.Printf("Vocabulary Word ID: %s\n", analysis.VocabularyWordID.Hex())
	fmt.Printf("Part of Speech: %s\n", analysis.PartOfSpeech)
	fmt.Printf("Pronunciation (IPA): %s\n", analysis.PronunciationIPA)
	fmt.Printf("Syllabification: %s\n", analysis.Syllabification)
	fmt.Printf("Definition: %s\n", analysis.Definition)
	fmt.Printf("Frequency: %s\n\n", analysis.Frequency)

	fmt.Println("Example Sentences:")
	for i, sentence := range analysis.ExampleSentences {
		fmt.Printf("  %d. %s\n", i+1, sentence)
	}

	fmt.Printf("\nSynonyms: %v\n", analysis.Synonyms)
	fmt.Printf("Antonyms: %v\n", analysis.Antonyms)
	fmt.Printf("Tags: %v\n\n", analysis.Tags)

	fmt.Printf("Etymology: %s\n\n", analysis.Etymology)

	fmt.Println("Usage Notes:")
	fmt.Printf("  Collocations: %v\n", analysis.UsageNotes.Collocations)
	fmt.Printf("  Cultural Significance: %s\n", analysis.UsageNotes.CulturalSignificance)
	fmt.Printf("  Register: %s\n", analysis.UsageNotes.Register)
}

func main() {
	// chatgpt:change - Get required environment variables
	apiKey := ""

	mongoURI := ""

	// chatgpt:change - Initialize MongoDB client
	fmt.Println("Connecting to MongoDB...")
	mongoClient, err := NewMongoDBClient(mongoURI, "toenglish")
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient.Close(ctx)
	}()

	fmt.Println("Successfully connected to MongoDB!")

	// Create Gemini client
	geminiClient := NewGeminiClient(apiKey)

	// Process all words from MongoDB
	fmt.Println("Starting batch processing of all vocabulary words...")
	if err := ProcessAllWords(geminiClient, mongoClient); err != nil {
		fmt.Printf("Error processing words: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Batch processing completed!")
}
