package storage

import "fmt"

// Relationship represents a directed edge between two resources.
type Relationship struct {
	ID               int64
	AssessmentID     string
	SourceID         string
	TargetID         string
	RelationshipType string
	Status           string // "resolved" or "unresolved"
	UnresolvedReason string
	TargetRegion     string
	TargetAccount    string
}

// InsertRelationship inserts a relationship record into the database.
func (s *Store) InsertRelationship(rel *Relationship) error {
	result, err := s.DB.Exec(`
		INSERT INTO relationships (assessment_id, source_id, target_id, relationship_type, status, unresolved_reason, target_region, target_account)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rel.AssessmentID,
		rel.SourceID,
		rel.TargetID,
		rel.RelationshipType,
		rel.Status,
		rel.UnresolvedReason,
		rel.TargetRegion,
		rel.TargetAccount,
	)
	if err != nil {
		return fmt.Errorf("inserting relationship: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	rel.ID = id
	return nil
}

// GetRelationshipsByAssessment returns all relationships for a given assessment.
func (s *Store) GetRelationshipsByAssessment(assessmentID string) ([]Relationship, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, source_id, target_id, relationship_type, status, unresolved_reason, target_region, target_account
		FROM relationships
		WHERE assessment_id = ?`, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("querying relationships by assessment: %w", err)
	}
	defer rows.Close()

	return scanRelationships(rows)
}

// GetRelationshipsBySource returns all relationships for a given source resource within an assessment.
func (s *Store) GetRelationshipsBySource(assessmentID, sourceID string) ([]Relationship, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, source_id, target_id, relationship_type, status, unresolved_reason, target_region, target_account
		FROM relationships
		WHERE assessment_id = ? AND source_id = ?`, assessmentID, sourceID)
	if err != nil {
		return nil, fmt.Errorf("querying relationships by source: %w", err)
	}
	defer rows.Close()

	return scanRelationships(rows)
}

// scanRelationships scans rows into a slice of Relationship.
func scanRelationships(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]Relationship, error) {
	var rels []Relationship
	for rows.Next() {
		var r Relationship
		if err := rows.Scan(
			&r.ID,
			&r.AssessmentID,
			&r.SourceID,
			&r.TargetID,
			&r.RelationshipType,
			&r.Status,
			&r.UnresolvedReason,
			&r.TargetRegion,
			&r.TargetAccount,
		); err != nil {
			return nil, fmt.Errorf("scanning relationship row: %w", err)
		}
		rels = append(rels, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating relationship rows: %w", err)
	}
	return rels, nil
}
