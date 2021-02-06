package util

// Document struct
type Document struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Deleted   bool   `json:"deleted"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
