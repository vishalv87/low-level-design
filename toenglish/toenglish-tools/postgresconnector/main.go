package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BaseModel contains common fields for all models with UUID primary key
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

// WordCategory represents a word_category with a primary name and alternate names
type WordCategory struct {
	BaseModel
	PrimaryName    string         `gorm:"type:text;not null;column:primary_name"`
	AlternateNames pq.StringArray `gorm:"type:text[];column:alternate_names"`
}

// TableName overrides the table name for WordCategory
func (WordCategory) TableName() string {
	return "word_categories"
}

// ClusterJSON represents the JSON structure for clusters
type ClusterJSON struct {
	ClusterID      int      `json:"cluster_id"`
	PrimaryName    string   `json:"primary_name"`
	AlternateNames []string `json:"alternate_names"`
}

// ClustersData represents the top-level JSON structure
type ClustersData struct {
	Clusters []ClusterJSON `json:"clusters"`
}

// Seeder interface defines the methods any seeder must implement
type Seeder interface {
	// GetTableName returns the name of the table
	GetTableName() string

	// GetData returns the seed data
	GetData() ([]interface{}, error)

	// ShouldSeed checks if seeding is necessary
	ShouldSeed(db *gorm.DB) bool
}

// WordCategorySeeder implements the Seeder interface for Cluster model
type WordCategorySeeder struct {
	JsonFilePath string
}

func (s WordCategorySeeder) GetTableName() string {
	return "word_categories"
}

func (s WordCategorySeeder) ShouldSeed(db *gorm.DB) bool {
	var count int64
	db.Table(s.GetTableName()).Count(&count)
	return count == 0
}

func (s WordCategorySeeder) GetData() ([]interface{}, error) {
	// Read JSON file
	fileData, err := ioutil.ReadFile(s.JsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	// Parse JSON data
	var clustersData ClustersData
	if err := json.Unmarshal(fileData, &clustersData); err != nil {
		return nil, fmt.Errorf("error parsing JSON data: %w", err)
	}

	// chatgpt:change - Change the return type to return pointers to structs
	wordCategories := make([]interface{}, len(clustersData.Clusters))
	for i, clusterJSON := range clustersData.Clusters {
		// chatgpt:change - Create struct pointers and properly convert string arrays
		category := &WordCategory{
			BaseModel:      BaseModel{ID: uuid.New()},
			PrimaryName:    clusterJSON.PrimaryName,
			AlternateNames: pq.StringArray(clusterJSON.AlternateNames),
		}
		wordCategories[i] = category
	}

	return wordCategories, nil
}

// DatabaseSeeder manages all seeders
type DatabaseSeeder struct {
	db      *gorm.DB
	seeders []Seeder
}

// NewDatabaseSeeder creates a new database seeder
func NewDatabaseSeeder(db *gorm.DB, clusterJsonPath string) *DatabaseSeeder {
	return &DatabaseSeeder{
		db: db,
		seeders: []Seeder{
			WordCategorySeeder{JsonFilePath: clusterJsonPath},
		},
	}
}

// Migrate performs schema migrations for all models
func (s *DatabaseSeeder) Migrate() error {
	fmt.Println("Starting database migration...")

	// Enable UUID extension in PostgreSQL if it's not already enabled
	s.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Add all model structs that need to be migrated
	err := s.db.AutoMigrate(&WordCategory{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("Database migration completed successfully")
	return nil
}

// Seed runs all registered seeders
func (s *DatabaseSeeder) Seed() error {
	fmt.Println("Starting database seeding...")

	for _, seeder := range s.seeders {
		tableName := seeder.GetTableName()

		if seeder.ShouldSeed(s.db) {
			fmt.Printf("Seeding %s table...\n", tableName)

			data, err := seeder.GetData()
			if err != nil {
				return fmt.Errorf("failed to get seed data for %s: %w", tableName, err)
			}

			if len(data) == 0 {
				fmt.Printf("No data to seed for %s\n", tableName)
				continue
			}

			// chatgpt:change - Modified batch creation to directly use pointers
			// Insert data in batches of 10
			batchSize := 10
			for i := 0; i < len(data); i += batchSize {
				end := i + batchSize
				if end > len(data) {
					end = len(data)
				}

				batch := data[i:end]
				for _, item := range batch {
					// The item is already a pointer, so we don't need the & operator
					if err := s.db.Create(item).Error; err != nil {
						return fmt.Errorf("failed to seed %s: %w", tableName, err)
					}
				}
			}

			fmt.Printf("Successfully seeded %d records into %s\n", len(data), tableName)
		} else {
			fmt.Printf("Skipping seed for %s (data already exists)\n", tableName)
		}
	}

	fmt.Println("Database seeding completed successfully")
	return nil
}

func main() {
	// Connection string for Neon PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DatabaseHost, DatabaseUsername, DatabasePassword, DatabaseName, DatabasePort, DatabaseSSLMode)

	// Path to JSON data file
	clusterJsonPath := "/Users/vishalvaibhav/Code/toenglish/toenglish-tools/wordcategorizer/cluster_data.json"

	// chatgpt:change - Add file existence check before proceeding
	if _, err := ioutil.ReadFile(clusterJsonPath); err != nil {
		log.Fatalf("Could not find or access the JSON file: %v", err)
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to Neon PostgreSQL database")

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Create and use the database seeder
	seeder := NewDatabaseSeeder(db, clusterJsonPath)

	// Run migrations
	if err := seeder.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Run seeders
	if err := seeder.Seed(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	// Print a summary of what was done
	fmt.Println("Database setup complete!")

	// Query and display a sample of data
	var wordCategories []WordCategory
	result := db.Limit(3).Find(&wordCategories)
	if result.Error != nil {
		log.Printf("Warning: Failed to query word categories: %v", result.Error)
	} else if len(wordCategories) > 0 {
		fmt.Println("\nSample word categories in database:")
		for _, category := range wordCategories {
			fmt.Printf("UUID: %s, Name: %s\n", category.ID, category.PrimaryName)
		}
	}
}
