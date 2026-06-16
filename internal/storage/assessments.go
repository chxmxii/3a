package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Assessment represents an assessment run stored in SQLite.
type Assessment struct {
	ID          string
	Profile     string
	Provider    string
	Status      string // "in_progress", "completed", "failed", "partial"
	StartedAt   time.Time
	CompletedAt *time.Time
	Regions     []string
}

// CreateAssessment inserts a new assessment record. The Regions slice is
// serialized as a JSON array and StartedAt/CompletedAt are stored as ISO 8601 text.
func (s *Store) CreateAssessment(a *Assessment) error {
	regions := a.Regions
	if regions == nil {
		regions = []string{}
	}

	regionsJSON, err := json.Marshal(regions)
	if err != nil {
		return fmt.Errorf("marshaling regions: %w", err)
	}

	var completedAt *string
	if a.CompletedAt != nil {
		t := a.CompletedAt.UTC().Format(time.RFC3339)
		completedAt = &t
	}

	_, err = s.DB.Exec(`
		INSERT INTO assessments (id, profile, provider, status, started_at, completed_at, regions)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.ID,
		a.Profile,
		a.Provider,
		a.Status,
		a.StartedAt.UTC().Format(time.RFC3339),
		completedAt,
		string(regionsJSON),
	)
	if err != nil {
		return fmt.Errorf("inserting assessment: %w", err)
	}

	return nil
}

// UpdateAssessmentStatus updates the status and optional completion time of an assessment.
func (s *Store) UpdateAssessmentStatus(id, status string, completedAt *time.Time) error {
	var completedAtStr *string
	if completedAt != nil {
		t := completedAt.UTC().Format(time.RFC3339)
		completedAtStr = &t
	}

	result, err := s.DB.Exec(`
		UPDATE assessments SET status = ?, completed_at = ? WHERE id = ?`,
		status, completedAtStr, id,
	)
	if err != nil {
		return fmt.Errorf("updating assessment status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("assessment not found: %s", id)
	}

	return nil
}

// GetAssessment retrieves a single assessment by ID. Returns nil if not found.
func (s *Store) GetAssessment(id string) (*Assessment, error) {
	row := s.DB.QueryRow(`
		SELECT id, profile, provider, status, started_at, completed_at, regions
		FROM assessments
		WHERE id = ?`, id)

	a, err := scanAssessment(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying assessment by id: %w", err)
	}
	return a, nil
}

// GetLatestAssessment returns the most recent assessment for a given profile.
// Returns nil if no assessments exist for the profile.
func (s *Store) GetLatestAssessment(profile string) (*Assessment, error) {
	row := s.DB.QueryRow(`
		SELECT id, profile, provider, status, started_at, completed_at, regions
		FROM assessments
		WHERE profile = ?
		ORDER BY started_at DESC
		LIMIT 1`, profile)

	a, err := scanAssessment(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying latest assessment: %w", err)
	}
	return a, nil
}

// ListAssessments returns all assessments ordered by started_at descending.
func (s *Store) ListAssessments() ([]Assessment, error) {
	rows, err := s.DB.Query(`
		SELECT id, profile, provider, status, started_at, completed_at, regions
		FROM assessments
		ORDER BY started_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying assessments: %w", err)
	}
	defer rows.Close()

	return scanAssessments(rows)
}

// scanAssessment scans a single row into an Assessment struct.
func scanAssessment(row *sql.Row) (*Assessment, error) {
	var a Assessment
	var startedAtStr string
	var completedAtStr *string
	var regionsJSON string

	err := row.Scan(
		&a.ID,
		&a.Profile,
		&a.Provider,
		&a.Status,
		&startedAtStr,
		&completedAtStr,
		&regionsJSON,
	)
	if err != nil {
		return nil, err
	}

	startedAt, err := time.Parse(time.RFC3339, startedAtStr)
	if err != nil {
		return nil, fmt.Errorf("parsing started_at: %w", err)
	}
	a.StartedAt = startedAt

	if completedAtStr != nil {
		t, err := time.Parse(time.RFC3339, *completedAtStr)
		if err != nil {
			return nil, fmt.Errorf("parsing completed_at: %w", err)
		}
		a.CompletedAt = &t
	}

	if err := json.Unmarshal([]byte(regionsJSON), &a.Regions); err != nil {
		return nil, fmt.Errorf("unmarshaling regions: %w", err)
	}

	return &a, nil
}

// scanAssessments scans multiple rows into a slice of Assessment structs.
func scanAssessments(rows *sql.Rows) ([]Assessment, error) {
	var assessments []Assessment

	for rows.Next() {
		var a Assessment
		var startedAtStr string
		var completedAtStr *string
		var regionsJSON string

		err := rows.Scan(
			&a.ID,
			&a.Profile,
			&a.Provider,
			&a.Status,
			&startedAtStr,
			&completedAtStr,
			&regionsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning assessment row: %w", err)
		}

		startedAt, err := time.Parse(time.RFC3339, startedAtStr)
		if err != nil {
			return nil, fmt.Errorf("parsing started_at: %w", err)
		}
		a.StartedAt = startedAt

		if completedAtStr != nil {
			t, err := time.Parse(time.RFC3339, *completedAtStr)
			if err != nil {
				return nil, fmt.Errorf("parsing completed_at: %w", err)
			}
			a.CompletedAt = &t
		}

		if err := json.Unmarshal([]byte(regionsJSON), &a.Regions); err != nil {
			return nil, fmt.Errorf("unmarshaling regions: %w", err)
		}

		assessments = append(assessments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating assessment rows: %w", err)
	}

	return assessments, nil
}
