# v0.1.0 — Initial Release

First public release of 3A (Agnostic Account Assessment).

## What is 3A

A terminal-first tool that assesses cloud accounts (AWS, OCI) using Steampipe. One command gives you a full picture: resource inventory, architecture map, security findings, cost analysis, and actionable reports.

## Features

### Discovery
- 25+ AWS resource types via Steampipe SQL queries
- Three-tier fallback for restricted IAM roles (SELECT * → specific columns → minimal)
- Concurrent multi-region discovery

### Security Assessment
- 26 rules covering EC2, RDS, S3, Lambda, VPC, ALB, EBS, IAM, EKS
- Standards: CIS Benchmarks, AWS Well-Architected Framework
- Categories: Security, Reliability, Cost Optimization, Operational Excellence

### Architecture
- Split views: Network (VPC/Subnet/SG/Gateways) and Resource (ALB/EKS/RDS relationships)
- Color-coded resource types in tree visualization

### Cost Analysis
- Primary: real billing data from AWS Cost Explorer via Steampipe
- Fallback: static pricing catalog (90+ instance types)
- Idle and oversized resource detection

### TUI
- 5 interactive views: Overview, Inventory, Architecture, Findings, Cost
- Region and type filtering (r/R, t/T)
- Resource detail panels (Enter) with type-aware rendering for SGs and Route Tables
- Scrollable views with keyboard navigation

### Reports
- Markdown and JSON export
- Excel workbook with 5 sheets (Summary, Inventory, Findings, Cost, Architecture)
- Auto-filters and formatted columns

### CLI
- `3a assess <profile>` — full pipeline with animated progress
- `3a configure` — interactive setup wizard (credentials, Steampipe config, regions)
- `3a report <profile> --format excel` — generate reports without TUI
- `3a profiles list/add` — profile management

## Install

```bash
curl -fsSL https://chxmxii.github.io/3a/install.sh | sh
```

Or build from source:
```bash
go install github.com/chxmxii/3a/cmd/3a@v0.1.0
```

## Requirements

- Steampipe with the AWS plugin installed and configured
- Read-only AWS credentials (ReadOnlyAccess or custom policy)

## Known Limitations

- OCI discovery is defined but untested (awaiting OCI Steampipe plugin setup)
- Cost estimation accuracy depends on `ce:GetCostAndUsage` permission
- Some Steampipe hydrate columns require permissions beyond basic ReadOnly
- Property-based tests not yet implemented

## What's Next (v0.2)

- Well-Architected pillar scoring
- 50+ assessment rules
- Tag compliance checks
- Assessment diff (compare runs)
