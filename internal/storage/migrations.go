package storage

// migrate runs all schema migrations. Tables and indexes use IF NOT EXISTS
// so this is safe to call on every Open.
func (s *Store) migrate() error {
	_, err := s.DB.Exec(schema)
	return err
}

const schema = `
CREATE TABLE IF NOT EXISTS assessments (
    id           TEXT PRIMARY KEY,
    profile      TEXT NOT NULL,
    provider     TEXT NOT NULL,
    status       TEXT NOT NULL,
    started_at   TEXT NOT NULL,
    completed_at TEXT,
    regions      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS resources (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id TEXT NOT NULL REFERENCES assessments(id),
    provider_type TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    region        TEXT NOT NULL,
    name          TEXT NOT NULL DEFAULT '',
    tags          TEXT NOT NULL DEFAULT '{}',
    raw_metadata  TEXT NOT NULL DEFAULT '{}',
    UNIQUE(assessment_id, resource_id)
);

CREATE INDEX IF NOT EXISTS idx_resources_assessment ON resources(assessment_id);
CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(assessment_id, resource_type);
CREATE INDEX IF NOT EXISTS idx_resources_region ON resources(assessment_id, region);

CREATE TABLE IF NOT EXISTS relationships (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id     TEXT NOT NULL REFERENCES assessments(id),
    source_id         TEXT NOT NULL,
    target_id         TEXT NOT NULL,
    relationship_type TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'resolved',
    unresolved_reason TEXT DEFAULT '',
    target_region     TEXT DEFAULT '',
    target_account    TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_relationships_assessment ON relationships(assessment_id);
CREATE INDEX IF NOT EXISTS idx_relationships_source ON relationships(assessment_id, source_id);

CREATE TABLE IF NOT EXISTS findings (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id  TEXT NOT NULL REFERENCES assessments(id),
    severity       TEXT NOT NULL,
    resource_id    TEXT NOT NULL,
    description    TEXT NOT NULL,
    recommendation TEXT NOT NULL,
    standard_name  TEXT NOT NULL,
    control_id     TEXT NOT NULL,
    category       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_findings_assessment ON findings(assessment_id);
CREATE INDEX IF NOT EXISTS idx_findings_severity ON findings(assessment_id, severity);
CREATE INDEX IF NOT EXISTS idx_findings_category ON findings(assessment_id, category);

CREATE TABLE IF NOT EXISTS cost_estimates (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id  TEXT NOT NULL REFERENCES assessments(id),
    resource_id    TEXT NOT NULL,
    resource_type  TEXT NOT NULL,
    monthly_cost   REAL,
    confidence     TEXT,
    category       TEXT NOT NULL,
    idle_flag      INTEGER NOT NULL DEFAULT 0,
    oversized_flag INTEGER NOT NULL DEFAULT 0,
    unestimable    INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_costs_assessment ON cost_estimates(assessment_id);

CREATE TABLE IF NOT EXISTS sizing (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id TEXT NOT NULL REFERENCES assessments(id),
    category      TEXT NOT NULL,
    resource_id   TEXT NOT NULL,
    data          TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sizing_assessment ON sizing(assessment_id);
CREATE INDEX IF NOT EXISTS idx_sizing_category ON sizing(assessment_id, category);
`
