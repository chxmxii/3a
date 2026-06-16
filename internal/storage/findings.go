package storage

import "fmt"

// Finding represents a standards-based assessment result.
type Finding struct {
	ID             int64
	AssessmentID   string
	Severity       string
	ResourceID     string
	Description    string
	Recommendation string
	StandardName   string
	ControlID      string
	Category       string
}

// InsertFinding inserts a new finding into the database.
func (s *Store) InsertFinding(finding *Finding) error {
	result, err := s.DB.Exec(`
		INSERT INTO findings (assessment_id, severity, resource_id, description, recommendation, standard_name, control_id, category)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		finding.AssessmentID,
		finding.Severity,
		finding.ResourceID,
		finding.Description,
		finding.Recommendation,
		finding.StandardName,
		finding.ControlID,
		finding.Category,
	)
	if err != nil {
		return fmt.Errorf("inserting finding: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	finding.ID = id
	return nil
}

// GetFindingsByAssessment returns all findings for a given assessment ID.
func (s *Store) GetFindingsByAssessment(assessmentID string) ([]Finding, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, severity, resource_id, description, recommendation, standard_name, control_id, category
		FROM findings
		WHERE assessment_id = ?`, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("querying findings by assessment: %w", err)
	}
	defer rows.Close()

	var findings []Finding
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ID, &f.AssessmentID, &f.Severity, &f.ResourceID, &f.Description, &f.Recommendation, &f.StandardName, &f.ControlID, &f.Category); err != nil {
			return nil, fmt.Errorf("scanning finding: %w", err)
		}
		findings = append(findings, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating findings: %w", err)
	}
	return findings, nil
}

// GetFindingsBySeverity returns findings for a given assessment filtered by severity.
func (s *Store) GetFindingsBySeverity(assessmentID, severity string) ([]Finding, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, severity, resource_id, description, recommendation, standard_name, control_id, category
		FROM findings
		WHERE assessment_id = ? AND severity = ?`, assessmentID, severity)
	if err != nil {
		return nil, fmt.Errorf("querying findings by severity: %w", err)
	}
	defer rows.Close()

	var findings []Finding
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ID, &f.AssessmentID, &f.Severity, &f.ResourceID, &f.Description, &f.Recommendation, &f.StandardName, &f.ControlID, &f.Category); err != nil {
			return nil, fmt.Errorf("scanning finding: %w", err)
		}
		findings = append(findings, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating findings: %w", err)
	}
	return findings, nil
}

// GetFindingsByCategory returns findings for a given assessment filtered by category.
func (s *Store) GetFindingsByCategory(assessmentID, category string) ([]Finding, error) {
	rows, err := s.DB.Query(`
		SELECT id, assessment_id, severity, resource_id, description, recommendation, standard_name, control_id, category
		FROM findings
		WHERE assessment_id = ? AND category = ?`, assessmentID, category)
	if err != nil {
		return nil, fmt.Errorf("querying findings by category: %w", err)
	}
	defer rows.Close()

	var findings []Finding
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ID, &f.AssessmentID, &f.Severity, &f.ResourceID, &f.Description, &f.Recommendation, &f.StandardName, &f.ControlID, &f.Category); err != nil {
			return nil, fmt.Errorf("scanning finding: %w", err)
		}
		findings = append(findings, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating findings: %w", err)
	}
	return findings, nil
}
