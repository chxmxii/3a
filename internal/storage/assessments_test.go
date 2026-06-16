package storage

import (
	"testing"
	"time"
)

func TestCreateAssessment(t *testing.T) {
	store := testStore(t)

	// Insert a second assessment (testStore already creates one).
	a := &Assessment{
		ID:        "assess-new-1",
		Profile:   "production",
		Provider:  "aws",
		Status:    "in_progress",
		StartedAt: time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		Regions:   []string{"us-east-1", "eu-west-1"},
	}

	err := store.CreateAssessment(a)
	if err != nil {
		t.Fatalf("CreateAssessment() error: %v", err)
	}

	// Verify it was stored.
	got, err := store.GetAssessment("assess-new-1")
	if err != nil {
		t.Fatalf("GetAssessment() error: %v", err)
	}
	if got == nil {
		t.Fatal("expected assessment, got nil")
	}
	if got.ID != "assess-new-1" {
		t.Errorf("ID = %q, want %q", got.ID, "assess-new-1")
	}
	if got.Profile != "production" {
		t.Errorf("Profile = %q, want %q", got.Profile, "production")
	}
	if got.Provider != "aws" {
		t.Errorf("Provider = %q, want %q", got.Provider, "aws")
	}
	if got.Status != "in_progress" {
		t.Errorf("Status = %q, want %q", got.Status, "in_progress")
	}
	if !got.StartedAt.Equal(time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)) {
		t.Errorf("StartedAt = %v, want 2024-06-15T10:30:00Z", got.StartedAt)
	}
	if got.CompletedAt != nil {
		t.Errorf("CompletedAt = %v, want nil", got.CompletedAt)
	}
	if len(got.Regions) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(got.Regions))
	}
	if got.Regions[0] != "us-east-1" || got.Regions[1] != "eu-west-1" {
		t.Errorf("Regions = %v, want [us-east-1, eu-west-1]", got.Regions)
	}
}

func TestCreateAssessmentWithCompletedAt(t *testing.T) {
	store := testStore(t)

	completedAt := time.Date(2024, 6, 15, 11, 0, 0, 0, time.UTC)
	a := &Assessment{
		ID:          "assess-completed-1",
		Profile:     "staging",
		Provider:    "oci",
		Status:      "completed",
		StartedAt:   time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		CompletedAt: &completedAt,
		Regions:     []string{"us-ashburn-1"},
	}

	err := store.CreateAssessment(a)
	if err != nil {
		t.Fatalf("CreateAssessment() error: %v", err)
	}

	got, err := store.GetAssessment("assess-completed-1")
	if err != nil {
		t.Fatalf("GetAssessment() error: %v", err)
	}
	if got.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}
	if !got.CompletedAt.Equal(completedAt) {
		t.Errorf("CompletedAt = %v, want %v", got.CompletedAt, completedAt)
	}
}

func TestCreateAssessmentNilRegions(t *testing.T) {
	store := testStore(t)

	a := &Assessment{
		ID:        "assess-nil-regions",
		Profile:   "default",
		Provider:  "aws",
		Status:    "in_progress",
		StartedAt: time.Now().UTC(),
		Regions:   nil,
	}

	err := store.CreateAssessment(a)
	if err != nil {
		t.Fatalf("CreateAssessment() with nil regions error: %v", err)
	}

	got, err := store.GetAssessment("assess-nil-regions")
	if err != nil {
		t.Fatalf("GetAssessment() error: %v", err)
	}
	if got.Regions == nil {
		t.Error("expected Regions to be non-nil after round-trip (should be empty slice)")
	}
	if len(got.Regions) != 0 {
		t.Errorf("expected 0 regions, got %d", len(got.Regions))
	}
}

func TestUpdateAssessmentStatus(t *testing.T) {
	store := testStore(t)

	a := &Assessment{
		ID:        "assess-update-1",
		Profile:   "production",
		Provider:  "aws",
		Status:    "in_progress",
		StartedAt: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		Regions:   []string{"us-east-1"},
	}
	if err := store.CreateAssessment(a); err != nil {
		t.Fatalf("CreateAssessment() error: %v", err)
	}

	completedAt := time.Date(2024, 6, 15, 10, 45, 0, 0, time.UTC)
	err := store.UpdateAssessmentStatus("assess-update-1", "completed", &completedAt)
	if err != nil {
		t.Fatalf("UpdateAssessmentStatus() error: %v", err)
	}

	got, err := store.GetAssessment("assess-update-1")
	if err != nil {
		t.Fatalf("GetAssessment() error: %v", err)
	}
	if got.Status != "completed" {
		t.Errorf("Status = %q, want %q", got.Status, "completed")
	}
	if got.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}
	if !got.CompletedAt.Equal(completedAt) {
		t.Errorf("CompletedAt = %v, want %v", got.CompletedAt, completedAt)
	}
}

func TestUpdateAssessmentStatusNotFound(t *testing.T) {
	store := testStore(t)

	err := store.UpdateAssessmentStatus("nonexistent-id", "completed", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent assessment, got nil")
	}
}

func TestGetAssessmentNotFound(t *testing.T) {
	store := testStore(t)

	got, err := store.GetAssessment("nonexistent-id")
	if err != nil {
		t.Fatalf("GetAssessment() error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for nonexistent assessment, got %+v", got)
	}
}

func TestGetLatestAssessment(t *testing.T) {
	store := testStore(t)

	assessments := []*Assessment{
		{
			ID:        "assess-old",
			Profile:   "myprofile",
			Provider:  "aws",
			Status:    "completed",
			StartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Regions:   []string{"us-east-1"},
		},
		{
			ID:        "assess-new",
			Profile:   "myprofile",
			Provider:  "aws",
			Status:    "in_progress",
			StartedAt: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
			Regions:   []string{"us-east-1", "us-west-2"},
		},
	}

	for _, a := range assessments {
		if err := store.CreateAssessment(a); err != nil {
			t.Fatalf("CreateAssessment() error: %v", err)
		}
	}

	got, err := store.GetLatestAssessment("myprofile")
	if err != nil {
		t.Fatalf("GetLatestAssessment() error: %v", err)
	}
	if got == nil {
		t.Fatal("expected assessment, got nil")
	}
	if got.ID != "assess-new" {
		t.Errorf("ID = %q, want %q (latest)", got.ID, "assess-new")
	}
}

func TestGetLatestAssessmentNotFound(t *testing.T) {
	store := testStore(t)

	got, err := store.GetLatestAssessment("nonexistent-profile")
	if err != nil {
		t.Fatalf("GetLatestAssessment() error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for nonexistent profile, got %+v", got)
	}
}

func TestListAssessments(t *testing.T) {
	store := testStore(t)

	assessments := []*Assessment{
		{
			ID:        "assess-list-1",
			Profile:   "profile-a",
			Provider:  "aws",
			Status:    "completed",
			StartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Regions:   []string{"us-east-1"},
		},
		{
			ID:        "assess-list-2",
			Profile:   "profile-b",
			Provider:  "oci",
			Status:    "in_progress",
			StartedAt: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
			Regions:   []string{"us-ashburn-1"},
		},
	}

	for _, a := range assessments {
		if err := store.CreateAssessment(a); err != nil {
			t.Fatalf("CreateAssessment() error: %v", err)
		}
	}

	got, err := store.ListAssessments()
	if err != nil {
		t.Fatalf("ListAssessments() error: %v", err)
	}

	// testStore inserts one assessment already, plus our 2.
	if len(got) != 3 {
		t.Fatalf("expected 3 assessments, got %d", len(got))
	}

	// Should be ordered by started_at DESC, so most recent first.
	if got[0].ID != "assess-list-2" {
		t.Errorf("first assessment ID = %q, want %q (most recent)", got[0].ID, "assess-list-2")
	}
}
