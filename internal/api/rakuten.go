package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RakutenClient represents a Rakuten Books API client
type RakutenClient struct {
	ApplicationID string
	HTTPClient    *http.Client
}

// BookInfo represents book information from Rakuten API
type BookInfo struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisherName"`
	Isbn        string `json:"isbn"`
	SalesDate   string `json:"salesDate"`
	ItemURL     string `json:"itemUrl"`
	MediumImage string `json:"mediumImageUrl"`
	Volume      string `json:"volume"`
}

// RakutenBooksResponse represents the API response
type RakutenBooksResponse struct {
	Items     []BookInfo `json:"Items"`
	PageCount int        `json:"pageCount"`
}

// NewRakutenClient creates a new Rakuten API client
func NewRakutenClient(appID string) *RakutenClient {
	return &RakutenClient{
		ApplicationID: appID,
		HTTPClient:    &http.Client{},
	}
}

// SearchByISBN searches for a book by ISBN
func (r *RakutenClient) SearchByISBN(isbn string) (*BookInfo, error) {
	baseURL := "https://app.rakuten.co.jp/services/api/BooksBook/Search/20170404"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("applicationId", r.ApplicationID)
	q.Set("isbnjan", isbn)
	q.Set("formatVersion", "2")
	u.RawQuery = q.Encode()

	resp, err := r.HTTPClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result RakutenBooksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("book not found")
	}

	return &result.Items[0], nil
}

// SearchByTitle searches for books by title
func (r *RakutenClient) SearchByTitle(title string) ([]BookInfo, error) {
	baseURL := "https://app.rakuten.co.jp/services/api/BooksBook/Search/20170404"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("applicationId", r.ApplicationID)
	q.Set("title", title)
	q.Set("formatVersion", "2")
	u.RawQuery = q.Encode()

	resp, err := r.HTTPClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result RakutenBooksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	books := make([]BookInfo, len(result.Items))
	for i, item := range result.Items {
		books[i] = item
	}

	return books, nil
}

// SearchByTitleSorted searches for books with sorting
func (r *RakutenClient) SearchByTitleSorted(title string, sort string, hits int) ([]BookInfo, error) {
	baseURL := "https://app.rakuten.co.jp/services/api/BooksBook/Search/20170404"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("applicationId", r.ApplicationID)
	q.Set("title", title)
	q.Set("formatVersion", "2")
	q.Set("sort", sort)
	if hits > 0 {
		q.Set("hits", fmt.Sprintf("%d", hits))
	}
	u.RawQuery = q.Encode()

	resp, err := r.HTTPClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result RakutenBooksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	books := make([]BookInfo, len(result.Items))
	for i, item := range result.Items {
		books[i] = item
	}

	return books, nil
}
