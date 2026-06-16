package checklist

// Status represents the outcome of a checklist item.
type Status string

const (
	StatusPass Status = "PASS"
	StatusFail Status = "FAIL"
	StatusWarn Status = "WARN"
)

// CheckItem represents a single checklist item.
type CheckItem struct {
	Name        string
	Description string
	Category    string
	Status      Status
	ResourceIDs []string
	Details     string
}

// ChecklistSummary holds the grouped checklist results.
type ChecklistSummary struct {
	Items      []CheckItem
	ByCategory map[string][]CheckItem
	PassCount  int
	FailCount  int
	WarnCount  int
}
