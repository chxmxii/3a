# Implementation Plan: Agnostic Account Assessment (3A)

## Overview

Implementation follows the pipeline architecture: Foundation (project setup, config, storage) → Core Engines (discovery, architecture, assessment, sizing, cost, checklist) → UI/Reporting (TUI, reports, CLI). Each task produces working, testable code that builds on previous tasks.

## Tasks

- [x] 1. Project foundation and storage layer
  - [x] 1.1 Initialize Go project structure and dependencies
    - Create `go.mod` with module `github.com/chxmxii/3a`
    - Create directory structure matching design (`cmd/3a/`, `internal/cli/`, `internal/config/`, `internal/provider/`, `internal/discovery/`, `internal/architecture/`, `internal/assessment/`, `internal/sizing/`, `internal/cost/`, `internal/checklist/`, `internal/report/`, `internal/storage/`, `internal/tui/`)
    - Add dependencies: `github.com/mattn/go-sqlite3`, `github.com/spf13/cobra`, `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, `gopkg.in/yaml.v3`, `github.com/flynaio/rapid` (test), `github.com/google/uuid`
    - Create minimal `cmd/3a/main.go` entry point
    - Create `Makefile` with `build`, `test`, `lint` targets
    - _Requirements: 9.1, 10.1_

  - [x] 1.2 Implement configuration types and loader
    - Implement `internal/config/config.go` with `Config`, `AccountProfile`, `Credential` types with YAML struct tags
    - Implement `Load(path string) (*Config, error)` that reads and parses `~/.3a/config.yaml`
    - Implement `DefaultConfigPath()` returning `~/.3a/config.yaml`
    - Return `ConfigError` for malformed or unreadable files
    - _Requirements: 9.2, 9.8, 9.9_

  - [x] 1.3 Implement profile management
    - Implement `internal/config/profile.go` with `GetProfile(name string) (*AccountProfile, error)`
    - Return `ProfileNotFoundError` with available profile names when not found
    - Implement `ListProfiles() []AccountProfile`
    - Implement `AddProfile(profile AccountProfile) error` for appending and saving
    - _Requirements: 9.3, 9.5, 9.6_

  - [ ]* 1.4 Write property test for configuration round-trip
    - **Property 20: Configuration Round-Trip**
    - Generate random valid Config structs with varied profiles, providers, credentials, and regions
    - Assert serialize-to-YAML then deserialize produces equivalent struct
    - **Validates: Requirements 9.2, 9.8**

  - [x] 1.5 Implement SQLite storage layer - schema and connection
    - Implement `internal/storage/db.go` with `Open(path string) (*Store, error)`
    - Auto-create database file and parent directories if not exist
    - Implement `internal/storage/migrations.go` with full schema creation (assessments, resources, relationships, findings, cost_estimates, sizing tables with all indexes)
    - Run migrations on Open
    - _Requirements: 10.1, 10.4_

  - [x] 1.6 Implement resource CRUD operations
    - Implement `internal/storage/resources.go` with `InsertResource`, `GetResourcesByAssessment`, `GetResourceByID`, `GetResourcesByType`, `GetResourcesByRegion`
    - Handle JSON serialization for tags and raw_metadata fields
    - Enforce UNIQUE constraint on (assessment_id, resource_id)
    - _Requirements: 1.3, 10.5_

  - [ ]* 1.7 Write property test for resource storage round-trip
    - **Property 1: Resource Storage Round-Trip**
    - Generate random DiscoveredResource with varied provider types, resource types, IDs, regions, names, tags, and metadata
    - Assert store-then-retrieve produces identical field values
    - **Validates: Requirements 1.3**

  - [x] 1.8 Implement relationship CRUD operations
    - Implement `internal/storage/relationships.go` with `InsertRelationship`, `GetRelationshipsByAssessment`, `GetRelationshipsBySource`
    - Store status, unresolved_reason, target_region, target_account fields
    - _Requirements: 2.3, 2.4_

  - [ ]* 1.9 Write property test for relationship storage round-trip
    - **Property 4: Relationship Storage Round-Trip**
    - Generate random Relationship structs with varied source/target IDs, types, and statuses
    - Assert store-then-retrieve produces identical field values
    - **Validates: Requirements 2.3**

  - [x] 1.10 Implement findings CRUD operations
    - Implement `internal/storage/findings.go` with `InsertFinding`, `GetFindingsByAssessment`, `GetFindingsBySeverity`, `GetFindingsByCategory`
    - _Requirements: 3.9_

  - [x] 1.11 Implement cost estimates and sizing CRUD operations
    - Implement `internal/storage/costs.go` with `InsertCostEstimate`, `GetCostsByAssessment`, `GetCostsByCategory`
    - Implement `internal/storage/assessments.go` with `CreateAssessment`, `UpdateAssessmentStatus`, `GetAssessment`, `GetLatestAssessment`, `ListAssessments`
    - Implement `internal/storage/sizing.go` (if not in costs.go) with `InsertSizing`, `GetSizingByAssessment`, `GetSizingByCategory`
    - _Requirements: 10.2, 10.3, 10.5, 10.6_

  - [ ]* 1.12 Write property test for assessment isolation
    - **Property 21: Assessment Persistence and Isolation**
    - Generate two assessments for the same profile, insert resources/findings into each
    - Assert querying by assessment ID returns only that assessment's data
    - **Validates: Requirements 10.2, 10.3, 10.5**

- [x] 2. Checkpoint - Foundation validation
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 3. Provider interface layer
  - [x] 3.1 Define provider interfaces and types
    - Implement `internal/provider/provider.go` with `Provider`, `Discoverer`, `MetricsClient`, `PricingClient` interfaces
    - Define `DiscoveredResource`, `ResourceType`, `PricingRequest`, `PricingResponse`, `PricingConfidence` types
    - _Requirements: 1.1, 1.2, 5.1_

  - [x] 3.2 Implement provider registry
    - Implement `internal/provider/registry.go` with `Registry`, `ProviderFactory`, `ProviderConfig`, `CredentialSource` types
    - Implement `NewRegistry()`, `Register(name string, factory ProviderFactory)`, `Get(cfg ProviderConfig) (Provider, error)`
    - _Requirements: 9.1, 9.10, 9.11_

  - [x] 3.3 Implement custom error types
    - Implement `internal/errors/errors.go` with `AuthError`, `DiscoveryError`, `RuleEvaluationError`, `ConfigError`, `ProfileNotFoundError`
    - Each error type implements the `error` interface with descriptive messages
    - _Requirements: 1.4, 1.6, 3.8, 9.3, 9.9_

  - [x] 3.4 Implement AWS provider skeleton with authentication
    - Implement `internal/provider/aws/provider.go` with `AWSProvider` struct implementing `Provider` interface
    - Implement `Authenticate(ctx)` using AWS SDK credential chain (named profile from `~/.aws/credentials` or environment variables)
    - Return `Discoverer()`, `MetricsClient()`, `PricingClient()` stubs
    - _Requirements: 1.5, 1.6, 9.10_

  - [x] 3.5 Implement OCI provider skeleton with authentication
    - Implement `internal/provider/oci/provider.go` with `OCIProvider` struct implementing `Provider` interface
    - Implement `Authenticate(ctx)` using OCI SDK config file (`~/.oci/config`)
    - Return `Discoverer()`, `MetricsClient()`, `PricingClient()` stubs
    - _Requirements: 1.5, 1.6, 9.11_

- [ ] 4. Discovery engine
  - [ ] 4.1 Implement discovery engine orchestration
    - Implement `internal/discovery/engine.go` with `Engine` struct, `NewEngine(provider, store, opts)`, `Run(ctx, assessmentID, regions) (DiscoverySummary, error)`
    - Implement bounded concurrency using semaphore pattern (max 10 goroutines)
    - Collect results from provider's Discoverer channel and persist to storage
    - Build and return `DiscoverySummary` with counts by type and region
    - _Requirements: 1.3, 1.7_

  - [ ] 4.2 Implement concurrent region scanner with retry logic
    - Implement `internal/discovery/scanner.go` with `retryWithBackoff(ctx, cfg, fn)` using exponential backoff (1s, 2s, 4s)
    - Implement `scanRegion(ctx, region, discoverer, results)` that wraps discovery calls with retry
    - Record `DiscoveryError` for exhausted retries, continue with remaining regions
    - _Requirements: 1.4, 1.7_

  - [ ]* 4.3 Write property test for retry logic
    - **Property 2: Retry Logic Invariants**
    - Generate random sequences of success/failure results
    - Assert max 3 retries, exponential backoff timing, failure recorded on exhaustion, continuation after failure
    - **Validates: Requirements 1.4**

  - [ ]* 4.4 Write property test for concurrency bound
    - **Property 3: Discovery Concurrency Bound**
    - Generate random region lists (1-50 regions)
    - Assert never more than 10 concurrent scans active at any point
    - **Validates: Requirements 1.7**

  - [ ] 4.5 Implement AWS resource discovery
    - Implement `internal/provider/aws/discovery.go` with `AWSDiscoverer` implementing `Discoverer` interface
    - Implement `DiscoverResources` enumerating: Organizations, Accounts, VPCs, Subnets, Route Tables, Internet Gateways, NAT Gateways, Transit Gateways, Security Groups, EC2, EKS, ECS, Lambda, RDS, S3, IAM, Load Balancers, Route53, CloudWatch, KMS, Secrets Manager
    - Stream results to channel as `DiscoveredResource` with full metadata
    - _Requirements: 1.1_

  - [ ] 4.6 Implement OCI resource discovery
    - Implement `internal/provider/oci/discovery.go` with `OCIDiscoverer` implementing `Discoverer` interface
    - Implement `DiscoverResources` enumerating: Compartments, VCNs, Subnets, Route Tables, Security Lists, NSGs, DRGs, Internet Gateways, NAT Gateways, Service Gateways, Compute Instances, OKE, Load Balancers, Databases, Object Storage, IAM, Vault, Logging, Monitoring
    - Stream results to channel as `DiscoveredResource` with full metadata
    - _Requirements: 1.2_

- [ ] 5. Checkpoint - Discovery validation
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 6. Architecture reconstructor
  - [ ] 6.1 Implement architecture reconstructor framework
    - Implement `internal/architecture/reconstructor.go` with `Reconstructor` struct, `NewReconstructor(store, rules)`, `Run(ctx, assessmentID) error`
    - Load resources from storage, apply each rule, persist inferred relationships
    - Handle unresolved relationships by recording partial edges with reason
    - _Requirements: 2.3, 2.4, 2.5_

  - [ ] 6.2 Implement AWS relationship rules
    - Implement `internal/architecture/rules.go` with AWS rules: VPC-to-Subnet, Subnet-to-Route Table, EC2-to-Security Group, EC2-to-EBS, ALB-to-Target Group-to-EC2, EKS-to-Node Group, Lambda-to-VPC, Transit Gateway-to-VPC
    - Each rule implements `RelationshipRule` interface, infers edges from resource metadata
    - _Requirements: 2.1_

  - [ ] 6.3 Implement OCI relationship rules
    - Implement OCI rules: VCN-to-Subnet, Subnet-to-Route Table, Compute-to-NSG, OKE-to-Node Pool, Load Balancer-to-Backend Set, VCN-to-DRG, Compute-to-Block Volume
    - Each rule implements `RelationshipRule` interface
    - _Requirements: 2.2_

  - [ ]* 6.4 Write property test for unresolved relationship handling
    - **Property 5: Unresolved Relationship Handling**
    - Generate resource sets with dangling references (missing targets, cross-region, cross-account)
    - Assert unresolved status recorded with non-empty reason and target annotations
    - **Validates: Requirements 2.4, 2.5**

- [ ] 7. Assessment engine
  - [ ] 7.1 Implement assessment engine framework
    - Implement `internal/assessment/engine.go` with `Engine` struct, `NewEngine(store, rules)`, `Run(ctx, assessmentID) error`
    - Implement `internal/assessment/rule.go` with `Rule` interface, `Finding`, `Severity`, `FindingCategory` types
    - Load resources, evaluate applicable rules per resource type, persist findings
    - On rule evaluation error: log, skip rule for that resource, continue
    - _Requirements: 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9_

  - [ ] 7.2 Implement assessment standards registry
    - Implement `internal/assessment/standards.go` with standard/control registry
    - Define mapping structures for AWS Well-Architected, CIS AWS, CIS OCI, OCI Well-Architected, NIST CSF
    - _Requirements: 3.1, 3.2_

  - [ ] 7.3 Implement AWS assessment rules (initial set)
    - Implement `internal/assessment/rules/aws/` with rules covering: S3 public access, security group open ports, unencrypted EBS, RDS public access, IAM overly permissive policies, EKS public endpoint, Lambda without VPC, missing encryption at rest
    - Each rule implements `Rule` interface with proper standard/control references
    - _Requirements: 3.1, 3.3, 3.4, 3.5, 3.6_

  - [ ] 7.4 Implement OCI assessment rules (initial set)
    - Implement `internal/assessment/rules/oci/` with rules covering: public buckets, NSG open ports, unencrypted block volumes, DB public access, overly permissive policies, OKE public endpoint, missing vault encryption
    - Each rule implements `Rule` interface with proper standard/control references
    - _Requirements: 3.2, 3.3, 3.4, 3.5, 3.6_

  - [ ]* 7.5 Write property tests for finding validity and engine resilience
    - **Property 6: Finding Completeness and Validity**
    - Generate random resources and evaluate against rules
    - Assert each finding has valid severity, category, non-empty resource ID, description ≤500 chars, recommendation ≤1000 chars, non-empty standard/control
    - **Property 7: Assessment Engine Resilience**
    - Generate rules that randomly error for specific resources
    - Assert engine skips failed rule, logs error, continues evaluating remaining rules
    - **Validates: Requirements 3.3, 3.4, 3.5, 3.6, 3.8**

- [ ] 8. Sizing analyzer
  - [ ] 8.1 Implement sizing analyzer
    - Implement `internal/sizing/analyzer.go` with `Analyzer` struct, `NewAnalyzer(store, metricsClient)`, `Run(ctx, assessmentID) (*SizingSummary, error)`
    - Implement `internal/sizing/types.go` with `SizingSummary`, `ComputeSummary`, `ComputeInstance`, `KubernetesSummary`, `KubernetesCluster`, `DatabaseSummary`, `DatabaseInstance`, `StorageSummary`, `StorageResource`
    - Extract sizing from resource metadata (instance type → vCPU/memory lookup)
    - Optionally enrich with metrics (CPU/memory utilization) when MetricsClient available
    - Aggregate totals (total vCPUs, total memory)
    - Omit categories with zero resources
    - Persist sizing entries to storage
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

  - [ ]* 8.2 Write property tests for sizing aggregation and category presence
    - **Property 8: Sizing Aggregation Correctness**
    - Generate random compute instances and kubernetes clusters with varied vCPU/memory
    - Assert total vCPUs = sum of individual, total memory = sum of individual
    - **Property 9: Sizing Category Presence and Omission**
    - Generate random resource sets with some empty categories
    - Assert category present iff at least one resource of that type discovered
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5**

- [ ] 9. Cost estimator
  - [ ] 9.1 Implement cost estimator engine
    - Implement `internal/cost/estimator.go` with `Estimator` struct, `NewEstimator(store, pricingClient, metricsClient)`, `Run(ctx, assessmentID) (*CostSummary, error)`
    - Implement `internal/cost/categories.go` with `CostCategory` constants and resource-type-to-category mapping
    - Calculate monthly cost as hourly_price × 730
    - Assign confidence levels based on pricing match quality
    - Identify top 5 cost drivers
    - Flag idle (CPU < 5% or zero traffic+connections) and oversized (peak < 20% capacity) resources
    - Mark unestimable resources when pricing unavailable
    - Skip idle/oversized evaluation when metrics unavailable
    - Persist cost estimates to storage
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_

  - [ ] 9.2 Implement AWS pricing client
    - Implement `internal/provider/aws/pricing.go` with `AWSPricingClient` implementing `PricingClient` interface
    - Query AWS Pricing API for on-demand prices by region, instance type, and attributes
    - Implement `internal/provider/aws/metrics.go` with `AWSMetricsClient` implementing `MetricsClient` for CloudWatch queries
    - _Requirements: 5.1, 5.5_

  - [ ] 9.3 Implement OCI pricing client
    - Implement `internal/provider/oci/pricing.go` with `OCIPricingClient` implementing `PricingClient` interface
    - Implement `internal/provider/oci/metrics.go` with `OCIMetricsClient` implementing `MetricsClient` for OCI Monitoring queries
    - _Requirements: 5.1, 5.5_

  - [ ]* 9.4 Write property tests for cost calculation and flagging
    - **Property 10: Monthly Cost Calculation**
    - Generate random hourly prices, assert monthly = hourly × 730
    - **Property 11: Cost Category Aggregation**
    - Generate random estimates with categories, assert category sums equal individual sums and total equals sum of categories
    - **Property 12: Top Cost Drivers Ordering**
    - Generate ≥5 cost estimates, assert top 5 are the highest sorted descending
    - **Property 13: Idle and Oversized Flagging**
    - Generate random utilization values around thresholds
    - Assert idle flag when CPU < 5% or zero traffic+connections, oversized when peak < 20%, no flags when metrics unavailable
    - **Validates: Requirements 5.1, 5.2, 5.4, 5.5, 5.7**

- [ ] 10. Checklist engine
  - [ ] 10.1 Implement checklist engine
    - Implement `internal/checklist/engine.go` with `Engine` struct, `NewEngine(store)`, `Generate(ctx, assessmentID) ([]Item, error)`
    - Implement `internal/checklist/types.go` with `Item`, `Status` types
    - Generate items only for rules applicable to discovered resource types
    - Derive status: FAIL when Critical/High findings, WARN when Medium/Low/Informational only, PASS when no findings
    - Group items by category
    - Set affected count = 0 for PASS, > 0 for FAIL/WARN
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [ ]* 10.2 Write property tests for checklist relevance and completeness
    - **Property 14: Checklist Relevance and Status Derivation**
    - Generate random resource types and findings, assert items only for discovered types, correct status derivation
    - **Property 15: Checklist Item Completeness**
    - Generate random checklist inputs, assert non-empty name, valid status/category, correct affected count
    - **Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5**

- [ ] 11. Checkpoint - Core engines validation
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 12. TUI application
  - [ ] 12.1 Implement TUI app shell and navigation
    - Implement `internal/tui/app.go` with root `App` model, view switching via numbered keys (1-5), persistent help bar
    - Implement `internal/tui/styles.go` with Lip Gloss style definitions
    - Implement `internal/tui/components/loading.go` with loading indicator component
    - Handle "no assessment data" state with prompt to run assessment
    - _Requirements: 7.1, 7.2, 7.8, 7.9, 7.11_

  - [ ] 12.2 Implement Overview view
    - Implement `internal/tui/overview.go` displaying: total resource count, findings by severity (Critical/High/Medium/Low/Informational), estimated monthly cost, assessment status
    - Load data from storage on view activation
    - _Requirements: 7.3_

  - [ ] 12.3 Implement Inventory view with pagination and filtering
    - Implement `internal/tui/inventory.go` with scrollable resource list, max 50 items per page
    - Implement `internal/tui/components/table.go` reusable table component
    - Implement `internal/tui/components/filter.go` filter input supporting resource type, region, provider filtering
    - Update displayed items within 200ms on filter application
    - _Requirements: 7.4, 7.10_

  - [ ] 12.4 Implement Architecture view with tree navigation
    - Implement `internal/tui/architecture.go` with tree format, collapsible parent-child nodes
    - Implement `internal/tui/components/tree.go` reusable tree view component with expand/collapse
    - _Requirements: 7.5_

  - [ ] 12.5 Implement Findings view with pagination and filtering
    - Implement `internal/tui/findings.go` with scrollable findings list, max 50 per page
    - Support filtering by severity, category, and standard
    - _Requirements: 7.6, 7.10_

  - [ ] 12.6 Implement Cost view
    - Implement `internal/tui/cost.go` displaying: cost breakdown by category, top 5 cost drivers, idle/oversized flagged resources
    - _Requirements: 7.7_

  - [ ]* 12.7 Write property tests for TUI pagination and filtering
    - **Property 16: TUI Pagination Invariant**
    - Generate random item lists of varying length, assert max 50 per page, pages = ceil(total/50)
    - **Property 17: TUI Filter Correctness**
    - Generate random items and filter criteria, assert all displayed items match filter, no matching items excluded
    - **Validates: Requirements 7.4, 7.6, 7.10**

- [ ] 13. Report generator
  - [ ] 13.1 Implement report data assembly and generator framework
    - Implement `internal/report/generator.go` with `Generator` struct, `NewGenerator(store, outputDir)`, `GenerateAll(ctx, assessmentID, profileName) error`
    - Implement `ReportData` aggregation from storage (resources, relationships, findings, costs, sizing, checklist)
    - Create output directory and parents if not exist
    - Return error if no completed assessment available
    - Prevent partial file output on write failure
    - _Requirements: 8.5, 8.6, 8.7, 8.8_

  - [ ] 13.2 Implement executive summary report
    - Implement `internal/report/executive.go` building executive summary: total resource count, top 5 findings per severity, estimated monthly cost, up to 10 prioritized recommendations
    - _Requirements: 8.1_

  - [ ] 13.3 Implement technical report
    - Implement `internal/report/technical.go` building technical report: full resource inventory, all relationships, all findings with details, complete cost breakdown, sizing summary
    - _Requirements: 8.2_

  - [ ] 13.4 Implement Markdown and JSON formatters
    - Implement `internal/report/markdown.go` rendering report data as Markdown
    - Implement `internal/report/json.go` rendering report data as JSON
    - File naming: `<profile-name>-<report-type>-<YYYYMMDD-HHMMSS>.<format-extension>`
    - Produce exactly 4 files per generation (executive.md, executive.json, technical.md, technical.json)
    - _Requirements: 8.3, 8.4_

  - [ ]* 13.5 Write property tests for report completeness and file output
    - **Property 18: Report Data Completeness**
    - Generate random assessment data, assert executive summary contains required fields, technical report contains all entries
    - **Property 19: Report File Output Correctness**
    - Generate random profile names and timestamps, assert exactly 4 files with correct naming pattern and valid format
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4**

- [ ] 14. Checkpoint - TUI and reports validation
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 15. CLI layer and integration wiring
  - [ ] 15.1 Implement Cobra CLI commands
    - Implement `internal/cli/root.go` with root command, `--db` flag for database path
    - Implement `internal/cli/assess.go` with `3a assess <profile-name>` command that loads profile, creates assessment, runs full pipeline (discovery → architecture → assessment → sizing → cost → checklist), updates assessment status
    - Implement `internal/cli/profiles.go` with `3a profiles list` and `3a profiles add` commands
    - Implement `internal/cli/report.go` with `3a report <profile-name>` command generating reports for latest completed assessment
    - Wire provider registry, engines, and storage together in command handlers
    - _Requirements: 9.1, 9.3, 9.4, 9.5, 9.6, 9.7_

  - [ ] 15.2 Wire TUI launch into CLI
    - Add TUI launch command (default when running `3a` with no subcommand or `3a tui`)
    - Connect TUI app to storage layer with assessment selection
    - Handle no-assessment-data state gracefully
    - _Requirements: 7.1, 7.9, 10.3_

  - [ ] 15.3 Implement end-to-end pipeline integration
    - Wire `cmd/3a/main.go` to CLI root command
    - Ensure assessment pipeline handles partial failures (mark status "partial", continue)
    - Handle authentication failure with non-zero exit and clear error message
    - Handle missing config file with clear error message
    - _Requirements: 1.6, 9.4, 9.9, 10.6_

- [ ] 16. Final checkpoint - Full integration
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation of each architectural layer
- Property tests validate universal correctness properties defined in the design document
- Unit tests validate specific examples and edge cases
- AWS provider is implemented before OCI per constraint, but interfaces are defined first
- All engines depend on the storage layer (tasks 1.5-1.11) being complete
- TUI depends on all storage CRUD operations
- Reports depend on all engines being complete to assemble full data

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["1.2", "1.5"] },
    { "id": 2, "tasks": ["1.3", "1.4", "1.6", "1.8", "1.10"] },
    { "id": 3, "tasks": ["1.7", "1.9", "1.11"] },
    { "id": 4, "tasks": ["1.12", "3.1", "3.3"] },
    { "id": 5, "tasks": ["3.2", "3.4", "3.5"] },
    { "id": 6, "tasks": ["4.1", "4.2"] },
    { "id": 7, "tasks": ["4.3", "4.4", "4.5", "4.6"] },
    { "id": 8, "tasks": ["6.1"] },
    { "id": 9, "tasks": ["6.2", "6.3"] },
    { "id": 10, "tasks": ["6.4", "7.1"] },
    { "id": 11, "tasks": ["7.2"] },
    { "id": 12, "tasks": ["7.3", "7.4"] },
    { "id": 13, "tasks": ["7.5", "8.1"] },
    { "id": 14, "tasks": ["8.2", "9.1"] },
    { "id": 15, "tasks": ["9.2", "9.3"] },
    { "id": 16, "tasks": ["9.4", "10.1"] },
    { "id": 17, "tasks": ["10.2"] },
    { "id": 18, "tasks": ["12.1", "13.1"] },
    { "id": 19, "tasks": ["12.2", "12.3", "12.4", "12.5", "12.6", "13.2", "13.3"] },
    { "id": 20, "tasks": ["12.7", "13.4"] },
    { "id": 21, "tasks": ["13.5", "15.1"] },
    { "id": 22, "tasks": ["15.2", "15.3"] }
  ]
}
```
