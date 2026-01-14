package manga

import (
	"encoding/json"
	"fmt"

	"github.com/kench/komikan-go/internal/db"
)

// Manga represents a manga entry
type Manga struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Series      string   `json:"series,omitempty"` // Series name if part of one
	Volume      int      `json:"volume"`           // Volume number
	ISBN        string   `json:"isbn"`
	Publisher   string   `json:"publisher"`
	PublishDate string   `json:"publish_date"`
	URL         string   `json:"url"` // Purchase URL
	Tags        []string `json:"tags,omitempty"`
}

// Manager manages manga collection
type Manager struct {
	db *db.DB
}

// NewManager creates a new manga manager
func NewManager(database *db.DB) *Manager {
	return &Manager{db: database}
}

// Add adds a manga to the collection
func (m *Manager) Add(manga Manga) error {
	if manga.ID == "" {
		manga.ID = manga.ISBN
	}

	key := fmt.Sprintf("manga:isbn:%s", manga.ISBN)
	return m.db.SetJSON(key, manga)
}

// GetByISBN retrieves a manga by ISBN
func (m *Manager) GetByISBN(isbn string) (*Manga, error) {
	key := fmt.Sprintf("manga:isbn:%s", isbn)
	var manga Manga
	if err := m.db.GetJSON(key, &manga); err != nil {
		return nil, err
	}
	return &manga, nil
}

// GetBySeries returns all manga in a series
func (m *Manager) GetBySeries(series string) ([]Manga, error) {
	key := fmt.Sprintf("manga:series:%s", series)
	var mangaList []Manga
	if err := m.db.GetJSON(key, &mangaList); err != nil {
		return nil, err
	}
	return mangaList, nil
}

// GetLatestVolume returns the latest volume for a given series/title
func (m *Manager) GetLatestVolume(series string) (int, error) {
	mangaList, err := m.GetBySeries(series)
	if err != nil {
		return 0, err
	}

	maxVol := 0
	for _, m := range mangaList {
		if m.Volume > maxVol {
			maxVol = m.Volume
		}
	}

	return maxVol, nil
}

// AddToSeries adds a manga to a series index
func (m *Manager) AddToSeries(manga Manga) error {
	if manga.Series == "" {
		return nil // No series to update
	}

	// Get existing series list
	key := fmt.Sprintf("manga:series:%s", manga.Series)
	var mangaList []Manga
	if err := m.db.GetJSON(key, &mangaList); err != nil {
		// First entry in series
		mangaList = []Manga{manga}
		return m.db.SetJSON(key, mangaList)
	}

	// Check if already exists
	for i, existing := range mangaList {
		if existing.ISBN == manga.ISBN {
			// Update existing entry
			mangaList[i] = manga
			return m.db.SetJSON(key, mangaList)
		}
	}

	// Add new entry
	mangaList = append(mangaList, manga)
	return m.db.SetJSON(key, mangaList)
}

// List returns all manga in the collection
func (m *Manager) List() ([]Manga, error) {
	values, err := m.db.ListPrefixJSON("manga:isbn:")
	if err != nil {
		return nil, err
	}

	mangaList := make([]Manga, 0, len(values))
	for _, v := range values {
		var m Manga
		if err := json.Unmarshal(v, &m); err != nil {
			continue // Skip invalid entries
		}
		mangaList = append(mangaList, m)
	}

	return mangaList, nil
}

// ListSeries returns all series names
func (m *Manager) ListSeries() ([]string, error) {
	keys, err := m.db.ListPrefix("manga:series:")
	if err != nil {
		return nil, err
	}

	// Extract series names from keys
	series := make([]string, 0, len(keys))
	prefixLen := len("manga:series:")
	for _, key := range keys {
		if len(key) > prefixLen {
			series = append(series, string(key[prefixLen:]))
		}
	}

	return series, nil
}

// Update updates a manga entry
func (m *Manager) Update(manga Manga) error {
	return m.Add(manga)
}

// Delete removes a manga from the collection
func (m *Manager) Delete(isbn string) error {
	key := fmt.Sprintf("manga:isbn:%s", isbn)
	return m.db.Delete([]byte(key))
}
