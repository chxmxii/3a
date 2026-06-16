package storage

import (
	"encoding/json"
	"fmt"
)

// SizingEntry represents a sizing data record stored in SQLite.
type SizingEntry struct {
	ID           int64
	AssessmentID string
	Category     string         // "compute", "kubernetes", "database", "storage"
	ResourceID   string
	Data         map[string]any // Category-specific fields, serialized as JSON
}

// InsertSizing inserts a sizing entry into the database, JSON-serializing the Data field.
func (s *Store) InsertSizing(entry *SizingEntry) error {
	data := entry.Data
	if data == nil {
		data = map[string]any{}
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling sizing data: %w", err)
	}

	result, err := s.DB.Exec(`
		INSERT INTO sizing (assessment_id, category, resource_id, data)
		VALUES (?, ?, ?, ?)`,
		entry.AssessmentID,
		entry.Category,
		entry.ResourceID,
		string(dataJSON),
	)
	if err != nil {
		return fmt.Errorf("inserting sizing entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	entry.ID = id

	return nil
}

// GetSizingByAssessment returns all sizing entries for the given assessment ID.
func (s *Store) GetSizingByAssessment(assessmentID string) ([]SizingEntry, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, category, resource_id, data
		FROM sizing
		WHERE assessment_id = ?`, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("querying sizing by assessment: %w", err)
	}
	defer rows.Close()

	return scanSizingEntries(rows)
}

// GetSizingByCategory returns sizing entries for a specific category within an assessment.
func (s *Store) GetSizingByCategory(assessmentID, category string) ([]SizingEntry, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, category, resource_id, data
		FROM sizing
		WHERE assessment_id = ? AND category = ?`, assessmentID, category)
	if err != nil {
		return nil, fmt.Errorf("querying sizing by category: %w", err)
	}
	defer rows.Close()

	return scanSizingEntries(rows)
}

// scanSizingEntries scans multiple rows into a slice of SizingEntry structs.
func scanSizingEntries(rows interface{ Next() bool; Scan(...any) error; Err() error }) ([]SizingEntry, error) {
	var entries []SizingEntry

	for rows.Next() {
		var entry SizingEntry
		var dataJSON string

		err := rows.Scan(
			&entry.ID,
			&entry.AssessmentID,
			&entry.Category,
			&entry.ResourceID,
			&dataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning sizing entry row: %w", err)
		}

		if err := json.Unmarshal([]byte(dataJSON), &entry.Data); err != nil {
			return nil, fmt.Errorf("unmarshaling sizing data: %w", err)
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating sizing entry rows: %w", err)
	}

	return entries, nil
}
