package main

// Document is Document struct
type Document struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Deleted   bool   `json:"deleted"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ChangeResultsChangesDocument is Change results document struct
type ChangeResultsChangesDocument struct {
	Rev string `json:"rev"`
}

// ChangeResultsDocument is Change results document struct
type ChangeResultsDocument struct {
	Changes []ChangeResultsChangesDocument `json:"changes"`
	ID      string                         `json:"id"`
	Seq     int                            `json:"seq"`
	Deleted bool                           `json:"deleted"`
}

// ChangeDocument is Change document struct
type ChangeDocument struct {
	LastSeq int                     `json:"last_seq"`
	Pending int                     `json:"pending"`
	Results []ChangeResultsDocument `json:"results"`
}

// ReplCheckpoint is Checkoint to latest replication
type ReplCheckpoint struct {
	ID  string `json:"id"`
	Ok  bool   `json:"ok"`
	Rev string `json:"rev"`
}
