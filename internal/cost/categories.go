package cost

// CostCategory groups resource costs by function.
type CostCategory string

const (
	CostCategoryCompute    CostCategory = "Compute"
	CostCategoryDatabase   CostCategory = "Database"
	CostCategoryStorage    CostCategory = "Storage"
	CostCategoryNetworking CostCategory = "Networking"
	CostCategoryKubernetes CostCategory = "Kubernetes"
	CostCategoryServerless CostCategory = "Serverless"
	CostCategoryOther      CostCategory = "Other"
)

// HoursPerMonth is the standard AWS billing hours per month.
const HoursPerMonth = 730.0

// pricingCatalog maps instance types to hourly on-demand prices (USD, us-east-1).
var pricingCatalog = map[string]float64{
	// General Purpose
	"t3.nano":     0.0052,
	"t3.micro":    0.0104,
	"t3.small":    0.0208,
	"t3.medium":   0.0416,
	"t3.large":    0.0832,
	"t3.xlarge":   0.1664,
	"t3.2xlarge":  0.3328,
	"m5.large":    0.096,
	"m5.xlarge":   0.192,
	"m5.2xlarge":  0.384,
	"m5.4xlarge":  0.768,
	"m5.8xlarge":  1.536,
	"m5.12xlarge": 2.304,
	"m5.16xlarge": 3.072,
	"m5.24xlarge": 4.608,
	"m6i.large":   0.096,
	"m6i.xlarge":  0.192,
	"m6i.2xlarge": 0.384,
	"m6i.4xlarge": 0.768,

	// Compute Optimized
	"c5.large":    0.085,
	"c5.xlarge":   0.170,
	"c5.2xlarge":  0.340,
	"c5.4xlarge":  0.680,
	"c5.9xlarge":  1.530,
	"c5.18xlarge": 3.060,

	// Memory Optimized
	"r5.large":    0.126,
	"r5.xlarge":   0.252,
	"r5.2xlarge":  0.504,
	"r5.4xlarge":  1.008,
	"r5.8xlarge":  2.016,
	"r5.12xlarge": 3.024,

	// RDS
	"db.t3.micro":   0.017,
	"db.t3.small":   0.034,
	"db.t3.medium":  0.068,
	"db.t3.large":   0.136,
	"db.m5.large":   0.171,
	"db.m5.xlarge":  0.342,
	"db.m5.2xlarge": 0.684,
	"db.m5.4xlarge": 1.368,
	"db.r5.large":   0.240,
	"db.r5.xlarge":  0.480,
	"db.r5.2xlarge": 0.960,
}

// storagePricing maps storage types to monthly per-GB prices.
var storagePricing = map[string]float64{
	"gp2":      0.10,
	"gp3":      0.08,
	"io1":      0.125,
	"io2":      0.125,
	"st1":      0.045,
	"sc1":      0.015,
	"standard": 0.05,
}

// natGatewayHourly is the hourly price for a NAT gateway.
const natGatewayHourly = 0.045

// albHourly is the hourly price for an ALB.
const albHourly = 0.0225

// nlbHourly is the hourly price for an NLB.
const nlbHourly = 0.0225

// resourceCategory maps resource types to cost categories.
var resourceCategory = map[string]CostCategory{
	"ec2_instance":      CostCategoryCompute,
	"compute_instance":  CostCategoryCompute,
	"rds_instance":      CostCategoryDatabase,
	"oci_database":      CostCategoryDatabase,
	"s3_bucket":         CostCategoryStorage,
	"object_storage":    CostCategoryStorage,
	"ebs_volume":        CostCategoryStorage,
	"block_volume":      CostCategoryStorage,
	"nat_gateway":       CostCategoryNetworking,
	"alb":              CostCategoryNetworking,
	"nlb":              CostCategoryNetworking,
	"oci_load_balancer": CostCategoryNetworking,
	"eks_cluster":       CostCategoryKubernetes,
	"oke_cluster":       CostCategoryKubernetes,
	"lambda_function":   CostCategoryServerless,
}

// GetCategory returns the cost category for a resource type.
func GetCategory(resourceType string) CostCategory {
	if cat, ok := resourceCategory[resourceType]; ok {
		return cat
	}
	return CostCategoryOther
}
