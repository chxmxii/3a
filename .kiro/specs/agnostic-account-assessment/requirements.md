# Requirements Document

## Introduction

3A (Agnostic Account Assessment) is a terminal-first TUI application written in Go that performs comprehensive cloud account assessments. The application discovers resources, reconstructs architecture relationships, evaluates standards compliance, estimates costs, and generates reports. Initial provider support covers AWS and OCI. The application uses Bubble Tea for the TUI, cloud-native SDKs for discovery, and SQLite for local storage.

## Glossary

- **3A**: The Agnostic Account Assessment application
- **Discovery_Engine**: The component responsible for enumerating all active cloud resources within a target account
- **Architecture_Reconstructor**: The component that infers and maps relationships between discovered resources
- **Assessment_Engine**: The component that evaluates discovered resources against compliance standards and best practices
- **Cost_Estimator**: The component that calculates estimated monthly costs for discovered resources
- **Sizing_Analyzer**: The component that consolidates infrastructure sizing information across compute, Kubernetes, databases, and storage
- **Checklist_Engine**: The component that generates adaptive assessment checklists based on discovered resource types
- **Report_Generator**: The component that produces executive and technical reports in Markdown and JSON formats
- **TUI**: Terminal User Interface built with the Bubble Tea framework in Go
- **Provider**: A cloud platform (AWS or OCI) against which assessments are performed
- **Resource**: A discrete cloud infrastructure entity (e.g., EC2 instance, VCN, S3 bucket)
- **Finding**: A standards-based assessment result containing severity, resource reference, description, recommendation, and standard reference
- **Account_Profile**: A named configuration referencing a cloud account with authentication credentials and provider type

## Requirements

### Requirement 1: Cloud Resource Discovery

**User Story:** As a cloud engineer, I want to discover all active resources within a cloud account, so that I have a complete inventory without manually enumerating services.

#### Acceptance Criteria

1. WHEN an assessment is initiated for an AWS account, THE Discovery_Engine SHALL enumerate resources across all regions specified in the Account_Profile for Organizations, Accounts, VPCs, Subnets, Route Tables, Internet Gateways, NAT Gateways, Transit Gateways, Security Groups, EC2, EKS, ECS, Lambda, RDS, S3, IAM, Load Balancers, Route53, CloudWatch, KMS, and Secrets Manager services.
2. WHEN an assessment is initiated for an OCI account, THE Discovery_Engine SHALL enumerate resources across all compartments accessible to the configured credentials for Compartments, VCNs, Subnets, Route Tables, Security Lists, NSGs, DRGs, Internet Gateways, NAT Gateways, Service Gateways, Compute Instances, OKE, Load Balancers, Databases, Object Storage, IAM, Vault, Logging, and Monitoring services.
3. WHEN discovery completes, THE Discovery_Engine SHALL store each discovered resource with its provider type, resource type, resource identifier, region, name, tags, and raw metadata in the SQLite database.
4. IF an API call fails during discovery, THEN THE Discovery_Engine SHALL retry the call up to 3 times with exponential backoff starting at 1 second, and if all retries fail, log the failure with service name, region, and error details, and continue discovering remaining services.
5. WHEN discovery begins, THE Discovery_Engine SHALL authenticate using the provider credentials configured in the Account_Profile.
6. IF authentication fails during discovery initiation, THEN THE Discovery_Engine SHALL report the authentication error with the provider name and credential source, and terminate the assessment with a non-zero exit code.
7. WHEN multiple regions are specified in the Account_Profile, THE Discovery_Engine SHALL enumerate resources in each region concurrently using parallel API calls bounded by a maximum of 10 concurrent region scans.

### Requirement 2: Architecture Relationship Reconstruction

**User Story:** As a cloud engineer, I want to see inferred relationships between resources, so that I understand the architecture without drawing diagrams manually.

#### Acceptance Criteria

1. WHEN discovery completes for an AWS account, THE Architecture_Reconstructor SHALL infer, at minimum, the following relationship types: VPC-to-Subnet, Subnet-to-Route Table, EC2-to-Security Group, EC2-to-EBS, ALB-to-Target Group-to-EC2, EKS-to-Node Group, Lambda-to-VPC, and Transit Gateway-to-VPC.
2. WHEN discovery completes for an OCI account, THE Architecture_Reconstructor SHALL infer, at minimum, the following relationship types: VCN-to-Subnet, Subnet-to-Route Table, Compute-to-NSG, OKE-to-Node Pool, Load Balancer-to-Backend Set, VCN-to-DRG, and Compute-to-Block Volume.
3. THE Architecture_Reconstructor SHALL store each relationship as a directed edge with source resource identifier (parent or containing resource), target resource identifier (child or contained resource), and relationship type in the SQLite database.
4. IF a relationship cannot be resolved because the target resource is not present in the discovered inventory or belongs to a different region or account than the source resource, THEN THE Architecture_Reconstructor SHALL record a partial relationship with an unresolved status, the known resource identifier, and the reason for non-resolution.
5. IF a referenced resource resides in a different region or account than the source resource, THEN THE Architecture_Reconstructor SHALL record the relationship with the available resource identifiers and annotate it with the target region or account when determinable from resource metadata.

### Requirement 3: Standards-Based Assessment

**User Story:** As a cloud engineer, I want resources assessed against industry standards, so that I can identify compliance gaps and risks without memorizing every benchmark.

#### Acceptance Criteria

1. WHEN discovery completes for an AWS account, THE Assessment_Engine SHALL evaluate all discovered resources against the AWS Well-Architected Framework and CIS Benchmarks for AWS and NIST CSF controls applicable to AWS.
2. WHEN discovery completes for an OCI account, THE Assessment_Engine SHALL evaluate all discovered resources against the OCI Well-Architected Framework and CIS Benchmarks for OCI and NIST CSF controls applicable to OCI.
3. THE Assessment_Engine SHALL generate one Finding per violated control per resource, such that a single resource may have multiple Findings across different standards and controls.
4. THE Assessment_Engine SHALL categorize each Finding into one of these categories: Security, Reliability, Performance, Cost Optimization, or Operational Excellence.
5. THE Assessment_Engine SHALL assign a severity level to each Finding based on the potential impact: Critical for issues enabling unauthorized access or data loss, High for issues degrading security posture or availability, Medium for issues representing non-compliance with best practices that carry moderate risk, Low for minor deviations from best practices with limited risk, and Informational for observations with no immediate risk.
6. THE Assessment_Engine SHALL include in each Finding: the assigned severity level, the affected resource identifier, a description of the non-compliant condition limited to 500 characters, a remediation recommendation limited to 1000 characters, and a reference to the specific standard name and control identifier.
7. WHEN no findings exist for a given category, THE Assessment_Engine SHALL report that category as compliant with zero findings.
8. IF a rule evaluation fails for a specific resource, THEN THE Assessment_Engine SHALL log the failure with the rule identifier, resource identifier, and error details, skip that rule for that resource, and continue evaluating remaining rules.
9. THE Assessment_Engine SHALL store all generated Findings in the SQLite database associated with the current assessment identifier.

### Requirement 4: Infrastructure Sizing Summary

**User Story:** As a cloud engineer, I want a consolidated sizing report, so that I can understand the infrastructure footprint at a glance.

#### Acceptance Criteria

1. WHEN discovery completes, THE Sizing_Analyzer SHALL produce a sizing summary for Compute resources including instance type, vCPU count, memory, aggregated total vCPUs, aggregated total memory, and current CPU and memory utilization percentage for each instance where the provider's monitoring service returns utilization metrics.
2. WHEN discovery completes, THE Sizing_Analyzer SHALL produce a sizing summary for Kubernetes clusters including node count, node instance types, total vCPUs, and total memory.
3. WHEN discovery completes, THE Sizing_Analyzer SHALL produce a sizing summary for Database resources including engine type, instance class, storage allocated, and multi-AZ or replication status.
4. WHEN discovery completes, THE Sizing_Analyzer SHALL produce a sizing summary for Storage resources including storage type, total capacity, and usage in bytes consumed for each storage resource where the provider's API returns consumption metrics.
5. IF discovery returns zero resources for a sizing category (Compute, Kubernetes, Database, or Storage), THEN THE Sizing_Analyzer SHALL omit that category from the sizing summary.
6. IF the provider's monitoring service does not return utilization metrics for a Compute or Storage resource, THEN THE Sizing_Analyzer SHALL include that resource in the sizing summary without utilization data and mark the utilization fields as unavailable.

### Requirement 5: Monthly Cost Estimation

**User Story:** As a cloud engineer, I want an estimated monthly cost for the account, so that I can understand spending without navigating billing consoles.

#### Acceptance Criteria

1. WHEN discovery completes, THE Cost_Estimator SHALL calculate an estimated monthly cost in USD for each discovered resource using on-demand pricing data for the resource's region, instance type or SKU, and provisioned capacity, with pricing data no older than 7 days.
2. THE Cost_Estimator SHALL provide a cost breakdown grouped by category (Compute, Storage, Database, Networking, Other), showing the sum of estimated costs for all resources within each category.
3. THE Cost_Estimator SHALL assign a confidence level to each cost estimate: High when pricing data exactly matches the resource's region and configuration, Medium when pricing is derived from a similar region or configuration fallback, and Low when pricing is partially interpolated or based on incomplete resource attributes.
4. THE Cost_Estimator SHALL identify and highlight the top five individual resources with the highest estimated monthly cost in the account.
5. THE Cost_Estimator SHALL flag a resource as idle if its average CPU utilization is below 5% over the available monitoring period or if it has zero network traffic and zero active connections, and flag a resource as oversized if its peak utilization across available metrics remains below 20% of provisioned capacity.
6. IF pricing data is unavailable for a resource type, THEN THE Cost_Estimator SHALL mark that resource's cost as unestimable and include it in a separate list.
7. IF the Cost_Estimator cannot retrieve utilization data for a resource, THEN THE Cost_Estimator SHALL skip idle and oversized evaluation for that resource and indicate that utilization data was unavailable.

### Requirement 6: Dynamic Assessment Checklist

**User Story:** As a cloud engineer, I want an adaptive checklist based on what is actually deployed, so that I only see relevant checks for my environment.

#### Acceptance Criteria

1. WHEN discovery completes, THE Checklist_Engine SHALL generate a checklist containing only checks relevant to the discovered resource types, where each check maps to one or more assessment rules evaluated by the Assessment_Engine.
2. THE Checklist_Engine SHALL assign each checklist item a status of PASS when all mapped assessment rules pass for all applicable resources, FAIL when at least one mapped rule produces a Critical or High severity finding, or WARN when mapped rules produce only Medium, Low, or Informational findings.
3. WHEN a resource type is not present in the account, THE Checklist_Engine SHALL omit checks related to that resource type from the checklist.
4. THE Checklist_Engine SHALL group checklist items by assessment category (Security, Reliability, Performance, Cost Optimization, Operational Excellence).
5. THE Checklist_Engine SHALL include in each checklist item: the check name, status, category, and a count of affected resources when the status is FAIL or WARN.

### Requirement 7: TUI Interface

**User Story:** As a cloud engineer, I want a terminal interface to navigate assessment results, so that I can review findings without leaving the terminal.

#### Acceptance Criteria

1. THE TUI SHALL be built using the Bubble Tea framework in Go.
2. THE TUI SHALL provide five navigable views: Overview, Inventory, Architecture, Findings, and Cost.
3. WHEN the Overview view is active, THE TUI SHALL display a summary including total resource count, total findings count grouped by severity (Critical, High, Medium, Low, Informational), estimated monthly cost, and assessment status (In Progress, Completed, or Failed).
4. WHEN the Inventory view is active, THE TUI SHALL display discovered resources in a scrollable list with a maximum of 50 items per page, supporting filtering by resource type, region, and provider.
5. WHEN the Architecture view is active, THE TUI SHALL display resource relationships in a tree format with collapsible parent-child nodes supporting expand and collapse navigation.
6. WHEN the Findings view is active, THE TUI SHALL display assessment findings in a scrollable list with a maximum of 50 items per page, supporting filtering by severity, category, and standard.
7. WHEN the Cost view is active, THE TUI SHALL display the cost breakdown by category, top five cost drivers, and flagged idle or oversized resources.
8. THE TUI SHALL support keyboard navigation between views using numbered keys (1-5) corresponding to each view, and SHALL display available key bindings in a persistent help bar.
9. IF no assessment data exists in the database, THEN THE TUI SHALL display a message indicating no assessment results are available and prompt the user to run an assessment.
10. WHEN a filter is applied in the Inventory or Findings view, THE TUI SHALL update the displayed list to show only matching items within 200 milliseconds.
11. THE TUI SHALL display a loading indicator while retrieving data from the SQLite database for any view transition.

### Requirement 8: Report Generation

**User Story:** As a cloud engineer, I want exportable reports, so that I can share assessment results with stakeholders who do not use the terminal.

#### Acceptance Criteria

1. WHEN the user requests a report, THE Report_Generator SHALL produce an Executive Summary containing total resource count, the top 5 findings per severity level, estimated monthly cost, and up to 10 recommendations prioritized by finding severity from Critical to Low.
2. WHEN the user requests a report, THE Report_Generator SHALL produce a Technical Report containing full resource inventory, all relationships, all findings with details, complete cost breakdown, and sizing summary.
3. THE Report_Generator SHALL output both the Executive Summary and the Technical Report in Markdown format and in JSON format, producing four files per report generation.
4. THE Report_Generator SHALL name each output file using the pattern `<profile-name>-<report-type>-<YYYYMMDD-HHMMSS>.<format-extension>` where report-type is "executive" or "technical" and format-extension is "md" or "json".
5. THE Report_Generator SHALL write report files to the local filesystem in a user-specified output directory, or to a default directory of `./3a-reports/` relative to the current working directory when no output directory is specified.
6. IF the target output directory does not exist, THEN THE Report_Generator SHALL create the directory and any required parent directories before writing report files.
7. IF writing a report file fails due to filesystem errors, THEN THE Report_Generator SHALL display an error message indicating the file path and failure reason, and SHALL not produce partial report files for that format.
8. IF a completed assessment is not available for the selected profile, THEN THE Report_Generator SHALL display an error message indicating that no assessment data exists and SHALL not generate report files.

### Requirement 9: Account Profile and CLI Interface

**User Story:** As a cloud engineer, I want to run assessments via a simple CLI command, so that I can integrate assessments into my workflow with minimal setup.

#### Acceptance Criteria

1. WHEN the user invokes `3a assess <profile-name>`, THE 3A SHALL load the Account_Profile matching the provided profile name and initiate a full assessment.
2. THE 3A SHALL support Account_Profile configuration specifying provider type (aws or oci), authentication credentials or credential source, target regions, and an optional display name.
3. IF the specified profile name does not exist, THEN THE 3A SHALL display an error message indicating the profile was not found and list available profiles.
4. IF authentication fails for the configured credentials, THEN THE 3A SHALL display an error message indicating the authentication failure reason and provider name, and exit with a non-zero status code.
5. THE 3A SHALL provide a `3a profiles list` command that displays all configured Account_Profiles with their name, provider type, configured regions, and display name.
6. THE 3A SHALL provide a `3a profiles add` command that guides the user through creating a new Account_Profile by prompting for provider type, profile name, credential source, and target regions.
7. THE 3A SHALL provide a `3a report <profile-name>` command that generates both Executive Summary and Technical Report for the most recent completed assessment of the specified profile without launching the TUI.
8. THE 3A SHALL store Account_Profiles in a YAML configuration file located at `~/.3a/config.yaml`.
9. IF the configuration file at `~/.3a/config.yaml` is malformed or unreadable, THEN THE 3A SHALL display an error message indicating the configuration file path and the parse error, and exit with a non-zero status code.
10. FOR AWS profiles, THE 3A SHALL support referencing named profiles from `~/.aws/credentials` or using environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN) as credential sources.
11. FOR OCI profiles, THE 3A SHALL support referencing configuration profiles from `~/.oci/config` as credential sources.

### Requirement 10: Data Persistence

**User Story:** As a cloud engineer, I want assessment results stored locally, so that I can review past assessments without re-running discovery.

#### Acceptance Criteria

1. THE 3A SHALL store all assessment data (resources, relationships, findings, cost estimates, sizing summaries) in a local SQLite database located at `~/.3a/assessments.db` by default, or at a path specified via the `--db` CLI flag or `db_path` configuration key.
2. WHEN an assessment completes, THE 3A SHALL record the assessment with a timestamp, profile name, and completion status (completed, failed, or partial).
3. THE 3A SHALL support viewing results from previous assessments by timestamp or assessment identifier via the TUI or CLI.
4. IF the SQLite database does not exist on first run, THEN THE 3A SHALL create the database file and parent directories as needed, and initialize the required schema automatically.
5. THE 3A SHALL store resources from each assessment in isolation, associated with a unique assessment identifier, so that multiple assessments of the same profile maintain independent historical records.
6. IF a database write operation fails during assessment, THEN THE 3A SHALL log the error with the operation type and affected record, mark the assessment status as partial, and continue processing remaining data.
