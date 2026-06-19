package oci

import (
	"context"
	"fmt"

	"github.com/chxmxii/a3/internal/assessment"
	"github.com/chxmxii/a3/internal/provider"
	"github.com/chxmxii/a3/internal/storage"
)

// AllRules returns all OCI assessment rules.
func AllRules() []assessment.Rule {
	return []assessment.Rule{
		&PublicBucketRule{},
		&NSGOpenIngressRule{},
		&UnencryptedVolumeRule{},
		&DBPublicAccessRule{},
	}
}

// PublicBucketRule checks for OCI Object Storage buckets with public access.
type PublicBucketRule struct{}

func (r *PublicBucketRule) ID() string                            { return "oci-public-bucket" }
func (r *PublicBucketRule) Standard() string                      { return "3A Security Baseline" }
func (r *PublicBucketRule) ControlID() string                     { return "SEC-008" }
func (r *PublicBucketRule) Category() assessment.FindingCategory  { return assessment.CategorySecurity }
func (r *PublicBucketRule) AppliesTo() []provider.ResourceType    { return []provider.ResourceType{provider.ResourceTypeObjectStorage} }

func (r *PublicBucketRule) Evaluate(_ context.Context, resource storage.Resource) ([]assessment.Finding, error) {
	meta := resource.RawMetadata

	// Steampipe: public_access_type field — "NoPublicAccess" is the secure value.
	publicAccess, ok := meta["public_access_type"].(string)
	if ok && publicAccess != "NoPublicAccess" && publicAccess != "" {
		return []assessment.Finding{{
			Severity:       assessment.SeverityHigh,
			ResourceID:     resource.ResourceID,
			Description:    fmt.Sprintf("Object storage bucket %s has public access type: %s", resource.Name, publicAccess),
			Recommendation: "Set public_access_type to NoPublicAccess unless public access is explicitly required",
			StandardName:   r.Standard(),
			ControlID:      r.ControlID(),
			Category:       r.Category(),
		}}, nil
	}

	return nil, nil
}

// NSGOpenIngressRule checks for NSGs allowing all ingress traffic.
type NSGOpenIngressRule struct{}

func (r *NSGOpenIngressRule) ID() string                            { return "oci-nsg-open-ingress" }
func (r *NSGOpenIngressRule) Standard() string                      { return "3A Security Baseline" }
func (r *NSGOpenIngressRule) ControlID() string                     { return "SEC-009" }
func (r *NSGOpenIngressRule) Category() assessment.FindingCategory  { return assessment.CategorySecurity }
func (r *NSGOpenIngressRule) AppliesTo() []provider.ResourceType    { return []provider.ResourceType{provider.ResourceTypeNSG} }

func (r *NSGOpenIngressRule) Evaluate(_ context.Context, resource storage.Resource) ([]assessment.Finding, error) {
	meta := resource.RawMetadata
	var findings []assessment.Finding

	// Steampipe: rules column contains security rules.
	rules, ok := meta["rules"].([]any)
	if !ok {
		return nil, nil
	}

	for _, rule := range rules {
		ruleMap, ok := rule.(map[string]any)
		if !ok {
			continue
		}

		direction, _ := ruleMap["direction"].(string)
		if direction != "INGRESS" {
			continue
		}

		source, _ := ruleMap["source"].(string)
		if source == "0.0.0.0/0" {
			protocol, _ := ruleMap["protocol"].(string)
			if protocol == "all" || protocol == "6" { // all or TCP
				findings = append(findings, assessment.Finding{
					Severity:       assessment.SeverityHigh,
					ResourceID:     resource.ResourceID,
					Description:    fmt.Sprintf("NSG %s allows ingress from 0.0.0.0/0 (protocol: %s)", resource.Name, protocol),
					Recommendation: "Restrict ingress rules to specific source CIDRs",
					StandardName:   r.Standard(),
					ControlID:      r.ControlID(),
					Category:       r.Category(),
				})
			}
		}
	}

	return findings, nil
}

// UnencryptedVolumeRule checks for OCI block volumes without encryption.
type UnencryptedVolumeRule struct{}

func (r *UnencryptedVolumeRule) ID() string                            { return "oci-volume-unencrypted" }
func (r *UnencryptedVolumeRule) Standard() string                      { return "3A Security Baseline" }
func (r *UnencryptedVolumeRule) ControlID() string                     { return "SEC-010" }
func (r *UnencryptedVolumeRule) Category() assessment.FindingCategory  { return assessment.CategorySecurity }
func (r *UnencryptedVolumeRule) AppliesTo() []provider.ResourceType    { return []provider.ResourceType{provider.ResourceTypeBlockVolume} }

func (r *UnencryptedVolumeRule) Evaluate(_ context.Context, resource storage.Resource) ([]assessment.Finding, error) {
	meta := resource.RawMetadata

	// Steampipe: is_hydrated field and kms_key_id.
	kmsKeyID, _ := meta["kms_key_id"].(string)
	if kmsKeyID == "" {
		// OCI volumes are encrypted by default with Oracle-managed keys,
		// but no customer-managed key indicates lower security posture.
		return []assessment.Finding{{
			Severity:       assessment.SeverityLow,
			ResourceID:     resource.ResourceID,
			Description:    fmt.Sprintf("Block volume %s uses Oracle-managed encryption (no customer-managed key)", resource.Name),
			Recommendation: "Consider using a customer-managed encryption key (KMS) for sensitive data",
			StandardName:   r.Standard(),
			ControlID:      r.ControlID(),
			Category:       r.Category(),
		}}, nil
	}

	return nil, nil
}

// DBPublicAccessRule checks for publicly accessible OCI database systems.
type DBPublicAccessRule struct{}

func (r *DBPublicAccessRule) ID() string                            { return "oci-db-public" }
func (r *DBPublicAccessRule) Standard() string                      { return "3A Security Baseline" }
func (r *DBPublicAccessRule) ControlID() string                     { return "SEC-011" }
func (r *DBPublicAccessRule) Category() assessment.FindingCategory  { return assessment.CategorySecurity }
func (r *DBPublicAccessRule) AppliesTo() []provider.ResourceType    { return []provider.ResourceType{provider.ResourceTypeOCIDB} }

func (r *DBPublicAccessRule) Evaluate(_ context.Context, resource storage.Resource) ([]assessment.Finding, error) {
	meta := resource.RawMetadata

	// Check if the subnet is public (no prohibition on public IPs).
	// This is a heuristic — DB systems in public subnets are publicly accessible.
	if hostname, ok := meta["hostname"].(string); ok && hostname != "" {
		// If there's a hostname and the subnet doesn't block public IPs.
		if nsgIDs, ok := meta["nsg_ids"].([]any); ok && len(nsgIDs) == 0 {
			return []assessment.Finding{{
				Severity:       assessment.SeverityMedium,
				ResourceID:     resource.ResourceID,
				Description:    fmt.Sprintf("Database system %s has no NSG protection", resource.Name),
				Recommendation: "Apply NSG rules to restrict database access to authorized sources only",
				StandardName:   r.Standard(),
				ControlID:      r.ControlID(),
				Category:       r.Category(),
			}}, nil
		}
	}

	return nil, nil
}
