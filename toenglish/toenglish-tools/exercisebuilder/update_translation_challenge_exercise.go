package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"strings"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Old structure models
// type OldTranslation struct {
// 	Language string `bson:"language" json:"language"`
// 	Sentence string `bson:"sentence" json:"sentence"`
// }

// type OldExerciseData struct {
// 	Original           string           `bson:"original" json:"original"`
// 	Translations       []OldTranslation `bson:"translations" json:"translations"`
// 	CorrectTranslation OldTranslation   `bson:"correctTranslation" json:"correctTranslation"`
// }

// type OldExerciseDocument struct {
// 	ID              primitive.ObjectID `bson:"_id" json:"_id"`
// 	TopicID         primitive.ObjectID `bson:"topicID" json:"topicID"`
// 	DifficultyLevel string             `bson:"difficultyLevel" json:"difficultyLevel"`
// 	ExerciseType    string             `bson:"exerciseType" json:"exerciseType"`
// 	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
// 	ExerciseData    OldExerciseData    `bson:"exerciseData" json:"exerciseData"`
// }

// // New structure models
// type Question struct {
// 	Text     string `bson:"text" json:"text"`
// 	Language string `bson:"language" json:"language"`
// }

// type Option struct {
// 	Index    int    `bson:"index" json:"index"`
// 	Text     string `bson:"text" json:"text"`
// 	Language string `bson:"language" json:"language"`
// }

// type CorrectTranslation struct {
// 	Language string `bson:"language" json:"language"`
// 	Sentence string `bson:"sentence" json:"sentence"`
// }

// type NewExerciseData struct {
// 	Question           Question           `bson:"question" json:"question"`
// 	Options            []Option           `bson:"options" json:"options"`
// 	CorrectAnswerIndex int                `bson:"correctAnswerIndex" json:"correctAnswerIndex"`
// 	CorrectTranslation CorrectTranslation `bson:"correctTranslation" json:"correctTranslation"`
// }

// type NewExerciseDocument struct {
// 	ID              primitive.ObjectID `bson:"_id" json:"_id"`
// 	TopicID         primitive.ObjectID `bson:"topicID" json:"topicID"`
// 	DifficultyLevel string             `bson:"difficultyLevel" json:"difficultyLevel"`
// 	ExerciseType    string             `bson:"exerciseType" json:"exerciseType"`
// 	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
// 	ExerciseData    NewExerciseData    `bson:"exerciseData" json:"exerciseData"`
// }

// // Sample incorrect translations for generating multiple choice options
// var sampleIncorrectTranslations = map[string][]string{
// 	"English": {
// 		"Hello, how do you do?",
// 		"Hi, how are you doing?",
// 		"Hello, what are you doing?",
// 		"Hey, how have you been?",
// 		"Hello, where are you going?",
// 		"Hi, what's your name?",
// 		"Hello, nice to meet you",
// 		"How do you do today?",
// 		"What are you up to?",
// 		"How's everything going?",
// 	},
// 	"Hindi": {
// 		"आप क्या कर रहे हैं?",
// 		"आपका नाम क्या है?",
// 		"आप कहाँ जा रहे हैं?",
// 		"आप कैसे हो?",
// 		"क्या हाल है?",
// 		"आप कहाँ से हैं?",
// 		"आपको कैसा लग रहा है?",
// 		"आप क्या करते हैं?",
// 		"आज कैसा दिन है?",
// 		"सब कुछ कैसा चल रहा है?",
// 	},
// }

// // MigrateExerciseData migrates exercises from old TranslationChallenge format to new TranslationMultipleChoice format
// func MigrateExerciseData(client *mongo.Client, databaseName, collectionName string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30000*time.Second)
// 	defer cancel()

// 	collection := client.Database(databaseName).Collection(collectionName)

// 	// Find all documents with the old exerciseType
// 	filter := bson.M{
// 		"exerciseType": "TranslationChallenge",
// 	}

// 	cursor, err := collection.Find(ctx, filter)
// 	if err != nil {
// 		return fmt.Errorf("failed to find documents: %v", err)
// 	}
// 	defer cursor.Close(ctx)

// 	var documentsToUpdate []OldExerciseDocument
// 	if err := cursor.All(ctx, &documentsToUpdate); err != nil {
// 		return fmt.Errorf("failed to decode documents: %v", err)
// 	}

// 	log.Printf("Found %d documents to migrate", len(documentsToUpdate))

// 	// Process each document
// 	for i, oldDoc := range documentsToUpdate {
// 		log.Printf("Processing document %d/%d: %s", i+1, len(documentsToUpdate), oldDoc.ID.Hex())

// 		// Transform the data
// 		newExerciseData, err := transformExerciseData(oldDoc.ExerciseData)
// 		if err != nil {
// 			log.Printf("Error transforming document %s: %v", oldDoc.ID.Hex(), err)
// 			continue
// 		}

// 		// Update the document
// 		updateFilter := bson.M{"_id": oldDoc.ID}
// 		update := bson.M{
// 			"$set": bson.M{
// 				// chatgpt:change - Not updating exerciseType, keeping it as is
// 				"exerciseData": newExerciseData,
// 				"updatedAt":    time.Now(),
// 			},
// 		}

// 		result, err := collection.UpdateOne(ctx, updateFilter, update)
// 		if err != nil {
// 			log.Printf("Error updating document %s: %v", oldDoc.ID.Hex(), err)
// 			continue
// 		}

// 		if result.ModifiedCount == 0 {
// 			log.Printf("Warning: No modifications made to document %s", oldDoc.ID.Hex())
// 		} else {
// 			log.Printf("Successfully updated document %s", oldDoc.ID.Hex())
// 		}
// 	}

// 	log.Printf("Migration completed. Processed %d documents", len(documentsToUpdate))
// 	return nil
// }

// // transformExerciseData converts old exercise data format to new format
// func transformExerciseData(oldData OldExerciseData) (NewExerciseData, error) {
// 	// Find the non-English translation to use as the question
// 	var questionTranslation OldTranslation
// 	var correctOption Option

// 	// Find the question (non-English translation) and correct answer
// 	for _, translation := range oldData.Translations {
// 		if translation.Language != "English" {
// 			questionTranslation = translation
// 		} else {
// 			// chatgpt:change - Using index 0 instead of ID "A"
// 			correctOption = Option{
// 				Index:    0,
// 				Text:     translation.Sentence,
// 				Language: translation.Language,
// 			}
// 		}
// 	}

// 	// If no non-English translation found, use the original text
// 	if questionTranslation.Language == "" {
// 		questionTranslation = OldTranslation{
// 			Language: "Unknown", // You might want to detect this
// 			Sentence: oldData.Original,
// 		}
// 	}

// 	// If no English translation found in translations array, use the correct translation
// 	if correctOption.Text == "" {
// 		// chatgpt:change - Using index 0 instead of ID "A"
// 		correctOption = Option{
// 			Index:    0,
// 			Text:     oldData.CorrectTranslation.Sentence,
// 			Language: oldData.CorrectTranslation.Language,
// 		}
// 	}

// 	// Generate incorrect options
// 	incorrectOptions := generateIncorrectOptions(correctOption.Text, correctOption.Language)

// 	// Combine correct and incorrect options
// 	allOptions := []Option{correctOption}

// 	// chatgpt:change - Using sequential indices instead of letter IDs
// 	for i, incorrectText := range incorrectOptions {
// 		if i < 3 { // Only add 3 incorrect options
// 			allOptions = append(allOptions, Option{
// 				Index:    i + 1, // Start from 1 since 0 is the correct option
// 				Text:     incorrectText,
// 				Language: correctOption.Language,
// 			})
// 		}
// 	}

// 	// Shuffle options while keeping track of correct answer
// 	// chatgpt:change - Function now returns correctAnswerIndex instead of correctAnswerID
// 	shuffledOptions, correctAnswerIndex := shuffleOptionsAndGetCorrectIndex(allOptions)

// 	return NewExerciseData{
// 		Question: Question{
// 			Text:     questionTranslation.Sentence,
// 			Language: questionTranslation.Language,
// 		},
// 		Options:            shuffledOptions,
// 		CorrectAnswerIndex: correctAnswerIndex,
// 		CorrectTranslation: CorrectTranslation{
// 			Language: oldData.CorrectTranslation.Language,
// 			Sentence: oldData.CorrectTranslation.Sentence,
// 		},
// 	}, nil
// }

// // generateIncorrectOptions creates plausible but incorrect translation options
// func generateIncorrectOptions(correctText, language string) []string {
// 	incorrectOptions := []string{}

// 	// Get sample incorrect translations for the language
// 	if samples, exists := sampleIncorrectTranslations[language]; exists {
// 		// Filter out any that match the correct translation
// 		for _, sample := range samples {
// 			if !strings.EqualFold(sample, correctText) {
// 				incorrectOptions = append(incorrectOptions, sample)
// 			}
// 		}
// 	}

// 	// If we don't have enough samples, generate some variations
// 	if len(incorrectOptions) < 3 {
// 		// Add some generic incorrect options
// 		genericOptions := []string{
// 			"This is an incorrect translation",
// 			"Another wrong answer",
// 			"Not the right translation",
// 		}

// 		for _, generic := range genericOptions {
// 			if len(incorrectOptions) < 3 {
// 				incorrectOptions = append(incorrectOptions, generic)
// 			}
// 		}
// 	}

// 	// Return only the first 3 to ensure we have exactly 3 incorrect options
// 	if len(incorrectOptions) > 3 {
// 		incorrectOptions = incorrectOptions[:3]
// 	}

// 	return incorrectOptions
// }

// // chatgpt:change - Renamed function and changed return type from string to int
// // shuffleOptionsAndGetCorrectIndex shuffles the options and returns the new index of the correct answer
// func shuffleOptionsAndGetCorrectIndex(options []Option) ([]Option, int) {
// 	// Create a copy to avoid modifying the original
// 	shuffled := make([]Option, len(options))
// 	copy(shuffled, options)

// 	// Find the correct answer before shuffling
// 	correctText := ""
// 	for _, option := range shuffled {
// 		if option.Index == 0 { // chatgpt:change - 0 was the correct answer initially instead of "A"
// 			correctText = option.Text
// 			break
// 		}
// 	}

// 	// Shuffle the options
// 	rand.Seed(time.Now().UnixNano())
// 	for i := range shuffled {
// 		j := rand.Intn(i + 1)
// 		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
// 	}

// 	// chatgpt:change - Reassign indices and find the new correct answer index
// 	correctAnswerIndex := 0

// 	for i, option := range shuffled {
// 		shuffled[i].Index = i // chatgpt:change - Using sequential indices starting from 0
// 		if option.Text == correctText {
// 			correctAnswerIndex = i // chatgpt:change - Return the index instead of letter ID
// 		}
// 	}

// 	return shuffled, correctAnswerIndex // chatgpt:change - Return int instead of string
// }

// // Example usage function
// func main() {
// 	// MongoDB connection
// 	connectionString := ""
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
// 	if err != nil {
// 		log.Fatal("Failed to connect to MongoDB:", err)
// 	}
// 	defer client.Disconnect(context.TODO())

// 	// Test the connection
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal("Failed to ping MongoDB:", err)
// 	}

// 	log.Println("Connected to MongoDB")

// 	// Run the migration
// 	databaseName := "toenglish"
// 	collectionName := "updatedexercises"

// 	err = MigrateExerciseData(client, databaseName, collectionName)
// 	if err != nil {
// 		log.Fatal("Migration failed:", err)
// 	}

// 	log.Println("Migration completed successfully")
// }

// // Alternative function to migrate specific documents by IDs
// func MigrateSpecificDocuments(client *mongo.Client, databaseName, collectionName string, documentIDs []string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	collection := client.Database(databaseName).Collection(collectionName)

// 	// Convert string IDs to ObjectIDs
// 	var objectIDs []primitive.ObjectID
// 	for _, idStr := range documentIDs {
// 		objectID, err := primitive.ObjectIDFromHex(idStr)
// 		if err != nil {
// 			log.Printf("Invalid ObjectID: %s, error: %v", idStr, err)
// 			continue
// 		}
// 		objectIDs = append(objectIDs, objectID)
// 	}

// 	// Find specific documents
// 	filter := bson.M{
// 		"_id":          bson.M{"$in": objectIDs},
// 		"exerciseType": "TranslationChallenge",
// 	}

// 	cursor, err := collection.Find(ctx, filter)
// 	if err != nil {
// 		return fmt.Errorf("failed to find documents: %v", err)
// 	}
// 	defer cursor.Close(ctx)

// 	var documentsToUpdate []OldExerciseDocument
// 	if err := cursor.All(ctx, &documentsToUpdate); err != nil {
// 		return fmt.Errorf("failed to decode documents: %v", err)
// 	}

// 	log.Printf("Found %d specific documents to migrate", len(documentsToUpdate))

// 	// Process each document (same logic as MigrateExerciseData)
// 	for i, oldDoc := range documentsToUpdate {
// 		log.Printf("Processing document %d/%d: %s", i+1, len(documentsToUpdate), oldDoc.ID.Hex())

// 		newExerciseData, err := transformExerciseData(oldDoc.ExerciseData)
// 		if err != nil {
// 			log.Printf("Error transforming document %s: %v", oldDoc.ID.Hex(), err)
// 			continue
// 		}

// 		updateFilter := bson.M{"_id": oldDoc.ID}
// 		update := bson.M{
// 			"$set": bson.M{
// 				// chatgpt:change - Not updating exerciseType, keeping it as is
// 				"exerciseData": newExerciseData,
// 				"updatedAt":    time.Now(),
// 			},
// 		}

// 		result, err := collection.UpdateOne(ctx, updateFilter, update)
// 		if err != nil {
// 			log.Printf("Error updating document %s: %v", oldDoc.ID.Hex(), err)
// 			continue
// 		}

// 		if result.ModifiedCount == 0 {
// 			log.Printf("Warning: No modifications made to document %s", oldDoc.ID.Hex())
// 		} else {
// 			log.Printf("Successfully updated document %s", oldDoc.ID.Hex())
// 		}
// 	}

// 	return nil
// }
