package manga

import (
	"fmt"
	"log"

	"github.com/kench/komikan-go/internal/api"
)

// NewReleaseCheckResult represents the result of a new release check
type NewReleaseCheckResult struct {
	SeriesTitle   string
	LatestVolume  int
	PreviousVolume int
	NewVolume     int
	Author        string
	ISBN          string
	URL           string
	SalesDate     string
}

// CheckNewReleases checks for new releases for registered manga
func (m *Manager) CheckNewReleases(rakutenAPIKey string) ([]NewReleaseCheckResult, error) {
	// Get all registered manga
	allManga, err := m.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list manga: %w", err)
	}

	// Group by series title
	seriesMap := make(map[string][]Manga)
	for _, mg := range allManga {
		if mg.Series == "" {
			continue // Skip non-series manga
		}
		seriesMap[mg.Series] = append(seriesMap[mg.Series], mg)
	}

	client := api.NewRakutenClient(rakutenAPIKey)
	var newReleases []NewReleaseCheckResult

	// Check each series for new releases
	for seriesTitle, manga := range seriesMap {
		// Get current latest volume from local database
		currentLatest := 0
		for _, mg := range manga {
			if mg.Volume > currentLatest {
				currentLatest = mg.Volume
			}
		}

		if currentLatest == 0 {
			continue // Skip if no volume info
		}

		// Search Rakuten for latest volume
		books, err := client.SearchByTitleSorted(seriesTitle, "-releaseDate", 30)
		if err != nil {
			log.Printf("Failed to search for %s: %v", seriesTitle, err)
			continue
		}

		// Find latest numbered volume from API
		var latestVolume int
		var latestBook *api.BookInfo
		for _, book := range books {
			info := ExtractVolumeInfo(book.Title)
			if info.HasVolume && info.Volume > latestVolume {
				latestVolume = info.Volume
				latestBook = &book
			}
		}

		// Check if new volume is available
		if latestBook != nil && latestVolume > currentLatest {
			result := NewReleaseCheckResult{
				SeriesTitle:    seriesTitle,
				LatestVolume:   latestVolume,
				PreviousVolume: currentLatest,
				NewVolume:      latestVolume,
				Author:         latestBook.Author,
				ISBN:           latestBook.Isbn,
				URL:            latestBook.ItemURL,
				SalesDate:      latestBook.SalesDate,
			}
			newReleases = append(newReleases, result)
		}
	}

	return newReleases, nil
}
