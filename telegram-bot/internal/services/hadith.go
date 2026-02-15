package services

import (
	"hadith-bot/internal/data"
	"hadith-bot/internal/logger"
	"hadith-bot/internal/models"
	"sync"
	"time"
)

// HadithService handles hadith-related operations
type HadithService struct {
	data       *models.CollectionData
	log        *logger.Logger
	apiURL     string
	apiKey     string
	apiTimeout time.Duration
	mu         sync.RWMutex
}

// NewHadithService creates a new hadith service
func NewHadithService(dataDir string, apiURL, apiKey string, apiTimeout time.Duration, log *logger.Logger) *HadithService {
	s := &HadithService{
		log:        log,
		apiURL:     apiURL,
		apiKey:     apiKey,
		apiTimeout: apiTimeout,
	}

	// Try to load data from files first
	if dataDir != "" {
		if loadedData, err := data.LoadHadithData(dataDir); err == nil && loadedData != nil {
			s.data = loadedData
			s.log.Info("Loaded hadith data from files")
			return s
		}
	}

	// Fall back to default data
	s.data = data.GetDefaultCollectionData()
	s.log.Info("Using default hadith data")

	return s
}

// GetCollections returns all available collections
func (s *HadithService) GetCollections() []models.Collection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetCollections()
}

// GetCollection returns a specific collection
func (s *HadithService) GetCollection(name string) *models.Collection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetCollection(name)
}

// GetBooks returns all books in a collection
func (s *HadithService) GetBooks(collection string) []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetBooks(collection)
}

// GetBook returns a specific book
func (s *HadithService) GetBook(collection string, bookNumber int) *models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetBook(collection, bookNumber)
}

// GetHadiths returns hadiths with pagination
func (s *HadithService) GetHadiths(collection string, bookNumber int, page int, limit int) models.HadithResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetHadiths(collection, bookNumber, page, limit)
}

// SearchHadiths searches hadiths by keyword
func (s *HadithService) SearchHadiths(query string, page int, limit int) models.SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Validate input
	if query == "" {
		return models.SearchResult{
			Hadiths:    []models.Hadith{},
			Total:      0,
			Page:       page,
			TotalPages: 0,
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	return s.data.SearchHadiths(query, page, limit)
}

// GetRandomHadith returns a random hadith
func (s *HadithService) GetRandomHadith() models.RandomHadithResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.GetRandomHadith()
}

// CollectionNames returns display names for collections
var CollectionNames = map[string]string{
	"bukhari":   "Sahih al-Bukhari",
	"muslim":    "Sahih Muslim",
	"abudawud":  "Sunan Abu Dawood",
	"tirmidhi":  "Jami' at-Tirmidhi",
	"nasai":     "Sunan an-Nasa'i",
	"ibnmajah":  "Sunan Ibn Majah",
}

// GetCollectionDisplayName returns the display name for a collection
func GetCollectionDisplayName(collection string) string {
	if name, ok := CollectionNames[collection]; ok {
		return name
	}
	return collection
}

