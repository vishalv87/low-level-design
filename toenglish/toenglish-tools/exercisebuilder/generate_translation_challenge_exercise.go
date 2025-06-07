package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

// TempExercise represents the source document structure
type TempExercise struct {
	ID                          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Level                       string             `bson:"level" json:"level"`
	Translations                []Translation      `bson:"translations" json:"translations"`
	TextHash                    string             `bson:"textHash" json:"textHash"`
	ChapterDescription          string             `bson:"chapterDescription" json:"chapterDescription"`
	ChapterDescriptionTitleHash string             `bson:"chapterDescriptionTitleHash" json:"chapterDescriptionTitleHash"`
	ChapterTitle                string             `bson:"chapterTitle" json:"chapterTitle"`
	ChapterTitleHash            string             `bson:"chapterTitleHash" json:"chapterTitleHash"`
	CourseTitle                 string             `bson:"courseTitle" json:"courseTitle"`
	CourseTitleHash             string             `bson:"courseTitleHash" json:"courseTitleHash"`
	TopicTitle                  string             `bson:"topicTitle" json:"topicTitle"`
	TopicTitleHash              string             `bson:"topicTitleHash" json:"topicTitleHash"`
}

// Translation represents a translation pair
type Translation struct {
	Language string `bson:"language" json:"language"`
	Sentence string `bson:"sentence" json:"sentence"`
}

// TranslationChallenge represents the target document structure
type TranslationChallenge struct {
	ID                          primitive.ObjectID      `bson:"_id,omitempty" json:"_id,omitempty"`
	ChapterDescription          string                  `bson:"chapterDescription" json:"chapterDescription"`
	ChapterDescriptionTitleHash string                  `bson:"chapterDescriptionTitleHash" json:"chapterDescriptionTitleHash"`
	ChapterTitle                string                  `bson:"chapterTitle" json:"chapterTitle"`
	ChapterTitleHash            string                  `bson:"chapterTitleHash" json:"chapterTitleHash"`
	CourseTitle                 string                  `bson:"courseTitle" json:"courseTitle"`
	CourseTitleHash             string                  `bson:"courseTitleHash" json:"courseTitleHash"`
	TopicTitle                  string                  `bson:"topicTitle" json:"topicTitle"`
	TopicTitleHash              string                  `bson:"topicTitleHash" json:"topicTitleHash"`
	DifficultyLevel             string                  `bson:"difficultyLevel" json:"difficultyLevel"`
	ExerciseType                string                  `bson:"exerciseType" json:"exerciseType"`
	CreatedAt                   primitive.DateTime      `bson:"createdAt" json:"createdAt"`
	UpdatedAt                   primitive.DateTime      `bson:"updatedAt" json:"updatedAt"`
	ExerciseData                TranslationExerciseData `bson:"exerciseData" json:"exerciseData"`
}

type TranslationQuestion struct {
	Text     string `bson:"text" json:"text"`
	Language string `bson:"language" json:"language"`
}

type TranslationOption struct {
	Index    int    `bson:"index" json:"index"`
	Text     string `bson:"text" json:"text"`
	Language string `bson:"language" json:"language"`
}

type CorrectTranslation struct {
	Language string `bson:"language" json:"language"`
	Sentence string `bson:"sentence" json:"sentence"`
}

type TranslationExerciseData struct {
	Question           TranslationQuestion `bson:"question" json:"question"`
	Options            []TranslationOption `bson:"options" json:"options"`
	CorrectAnswerIndex int                 `bson:"correctAnswerIndex" json:"correctAnswerIndex"`
	CorrectTranslation CorrectTranslation  `bson:"correctTranslation" json:"correctTranslation"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Sentences []string `json:"sentences"`
}

func main() {
	// Environment variables
	mongoURI := ""
	geminiAPIKey := ""

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Initialize Gemini client
	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		log.Fatal("Failed to create Gemini client:", err)
	}
	defer geminiClient.Close()

	// Get collections
	db := client.Database("toenglish")
	sourceCollection := db.Collection("tempexercises")
	targetCollection := db.Collection("exercise-translationchallenge")

	// Process documents
	err = processDocuments(ctx, sourceCollection, targetCollection, geminiClient)
	if err != nil {
		log.Fatal("Error processing documents:", err)
	}

	fmt.Println("Successfully processed all documents!")
}

func processDocuments(ctx context.Context, sourceCollection, targetCollection *mongo.Collection, geminiClient *genai.Client) error {
	// Find all documents in source collection
	cursor, err := sourceCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	// Process each document
	for cursor.Next(ctx) {
		var tempExercise TempExercise
		if err := cursor.Decode(&tempExercise); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}

		// Extract English and Hindi translations
		var englishTranslation, hindiTranslation string
		for _, translation := range tempExercise.Translations {
			if strings.ToLower(translation.Language) == "english" {
				englishTranslation = translation.Sentence
			} else if strings.ToLower(translation.Language) == "hindi" {
				hindiTranslation = translation.Sentence
			}
		}

		if englishTranslation == "" || hindiTranslation == "" {
			log.Printf("Missing translations for document %s", tempExercise.ID.Hex())
			continue
		}

		// Generate incorrect options using Gemini
		incorrectOptions, err := generateIncorrectOptions(ctx, geminiClient, englishTranslation)
		if err != nil {
			log.Printf("Failed to generate options for document %s: %v", tempExercise.ID.Hex(), err)
			continue
		}

		// Create translation challenge document
		challenge := createTranslationChallenge(tempExercise, hindiTranslation, englishTranslation, incorrectOptions)

		// Insert into target collection
		_, err = targetCollection.InsertOne(ctx, challenge)
		if err != nil {
			log.Printf("Failed to insert challenge document: %v", err)
			continue
		}

		fmt.Printf("Processed document %s successfully\n", tempExercise.ID.Hex())
	}

	return cursor.Err()
}

func generateIncorrectOptions(ctx context.Context, client *genai.Client, correctTranslation string) ([]string, error) {
	model := client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`Given this English sentence: "%s"

Generate exactly 3 English sentences that are similar to the given sentence but are NOT correct translations of any Hindi text. The sentences should:
1. Be grammatically correct English
2. Have similar structure or words to the original sentence
3. Be plausible but incorrect translation options
4. Each sentence should be different from the others
5. If the sentence has names, then use same names in the generated sentences

Return only the 3 sentences, one per line, with no numbering or additional text.`, correctTranslation)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Extract text from response
	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	lines := strings.Split(strings.TrimSpace(responseText), "\n")

	var options []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			options = append(options, line)
		}
	}

	// Ensure we have exactly 3 options
	if len(options) < 3 {
		return nil, fmt.Errorf("insufficient options generated: got %d, need 3", len(options))
	}

	return options[:3], nil
}

func createTranslationChallenge(tempExercise TempExercise, hindiText, correctEnglishText string, incorrectOptions []string) TranslationChallenge {
	now := primitive.NewDateTimeFromTime(time.Now())

	// Create options array with correct answer and incorrect ones
	options := []TranslationOption{
		{
			Index:    0,
			Text:     correctEnglishText,
			Language: "English",
		},
	}

	// Add incorrect options
	optionIDs := []int{1, 2, 3}
	for i, incorrectOption := range incorrectOptions {
		if i < 3 { // Ensure we don't exceed array bounds
			options = append(options, TranslationOption{
				Index:    optionIDs[i],
				Text:     incorrectOption,
				Language: "English",
			})
		}
	}

	// Shuffle options to randomize correct answer position
	shuffleOptions(options)

	// Find the correct answer after shuffle
	var correctAnswerIndex int
	for _, option := range options {
		if option.Text == correctEnglishText {
			correctAnswerIndex = option.Index
			break
		}
	}

	return TranslationChallenge{
		ChapterDescription:          tempExercise.ChapterDescription,
		ChapterDescriptionTitleHash: tempExercise.ChapterDescriptionTitleHash,
		ChapterTitle:                tempExercise.ChapterTitle,
		ChapterTitleHash:            tempExercise.ChapterTitleHash,
		CourseTitle:                 tempExercise.CourseTitle,
		CourseTitleHash:             tempExercise.CourseTitleHash,
		TopicTitle:                  tempExercise.TopicTitle,
		TopicTitleHash:              tempExercise.TopicTitleHash,
		DifficultyLevel:             tempExercise.Level,
		ExerciseType:                "TranslationChallenge",
		CreatedAt:                   now,
		UpdatedAt:                   now,
		ExerciseData: TranslationExerciseData{
			Question: TranslationQuestion{
				Text:     hindiText,
				Language: "Hindi",
			},
			Options:            options,
			CorrectAnswerIndex: correctAnswerIndex,
			CorrectTranslation: CorrectTranslation{
				Language: "English",
				Sentence: correctEnglishText,
			},
		},
	}
}

func shuffleOptions(options []TranslationOption) {
	rand.Seed(time.Now().UnixNano())

	for i := len(options) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	// Reassign IDs after shuffle
	optionIDs := []int{0, 1, 2, 3}
	for i := range options {
		options[i].Index = optionIDs[i]
	}
}
