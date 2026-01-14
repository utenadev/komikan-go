package manga

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
	// TODO: Add database connection
}

// NewManager creates a new manga manager
func NewManager() *Manager {
	return &Manager{}
}

// Add adds a manga to the collection
func (m *Manager) Add(manga Manga) error {
	// TODO: Implement database insert
	return nil
}

// GetLatestVolume returns the latest volume for a given series/title
func (m *Manager) GetLatestVolume(title string) (int, error) {
	// TODO: Implement database query
	return 0, nil
}

// List returns all manga in the collection
func (m *Manager) List() ([]Manga, error) {
	// TODO: Implement database query
	return nil, nil
}
