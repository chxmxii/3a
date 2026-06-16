package sizing

// SizingCategory groups resources by their function.
type SizingCategory string

const (
	CategoryCompute    SizingCategory = "compute"
	CategoryKubernetes SizingCategory = "kubernetes"
	CategoryDatabase   SizingCategory = "database"
	CategoryStorage    SizingCategory = "storage"
)

// InstanceSpec describes a compute instance's specifications.
type InstanceSpec struct {
	VCPUs  int
	MemGB  float64
	Family string
}

// SizingSummary aggregates sizing data across all resources.
type SizingSummary struct {
	TotalVCPUs    int
	TotalMemoryGB float64
	TotalStorageGB float64
	ByCategory    map[SizingCategory]CategorySizing
}

// CategorySizing holds sizing totals for a category.
type CategorySizing struct {
	Count     int
	VCPUs     int
	MemoryGB  float64
	StorageGB float64
}

// instanceCatalog maps instance types to their specs.
var instanceCatalog = map[string]InstanceSpec{
	// General Purpose
	"t3.nano":     {VCPUs: 2, MemGB: 0.5, Family: "general"},
	"t3.micro":    {VCPUs: 2, MemGB: 1, Family: "general"},
	"t3.small":    {VCPUs: 2, MemGB: 2, Family: "general"},
	"t3.medium":   {VCPUs: 2, MemGB: 4, Family: "general"},
	"t3.large":    {VCPUs: 2, MemGB: 8, Family: "general"},
	"t3.xlarge":   {VCPUs: 4, MemGB: 16, Family: "general"},
	"t3.2xlarge":  {VCPUs: 8, MemGB: 32, Family: "general"},
	"m5.large":    {VCPUs: 2, MemGB: 8, Family: "general"},
	"m5.xlarge":   {VCPUs: 4, MemGB: 16, Family: "general"},
	"m5.2xlarge":  {VCPUs: 8, MemGB: 32, Family: "general"},
	"m5.4xlarge":  {VCPUs: 16, MemGB: 64, Family: "general"},
	"m5.8xlarge":  {VCPUs: 32, MemGB: 128, Family: "general"},
	"m5.12xlarge": {VCPUs: 48, MemGB: 192, Family: "general"},
	"m5.16xlarge": {VCPUs: 64, MemGB: 256, Family: "general"},
	"m5.24xlarge": {VCPUs: 96, MemGB: 384, Family: "general"},
	"m6i.large":   {VCPUs: 2, MemGB: 8, Family: "general"},
	"m6i.xlarge":  {VCPUs: 4, MemGB: 16, Family: "general"},
	"m6i.2xlarge": {VCPUs: 8, MemGB: 32, Family: "general"},
	"m6i.4xlarge": {VCPUs: 16, MemGB: 64, Family: "general"},

	// Compute Optimized
	"c5.large":    {VCPUs: 2, MemGB: 4, Family: "compute"},
	"c5.xlarge":   {VCPUs: 4, MemGB: 8, Family: "compute"},
	"c5.2xlarge":  {VCPUs: 8, MemGB: 16, Family: "compute"},
	"c5.4xlarge":  {VCPUs: 16, MemGB: 32, Family: "compute"},
	"c5.9xlarge":  {VCPUs: 36, MemGB: 72, Family: "compute"},
	"c5.18xlarge": {VCPUs: 72, MemGB: 144, Family: "compute"},

	// Memory Optimized
	"r5.large":    {VCPUs: 2, MemGB: 16, Family: "memory"},
	"r5.xlarge":   {VCPUs: 4, MemGB: 32, Family: "memory"},
	"r5.2xlarge":  {VCPUs: 8, MemGB: 64, Family: "memory"},
	"r5.4xlarge":  {VCPUs: 16, MemGB: 128, Family: "memory"},
	"r5.8xlarge":  {VCPUs: 32, MemGB: 256, Family: "memory"},
	"r5.12xlarge": {VCPUs: 48, MemGB: 384, Family: "memory"},

	// RDS Instance Classes
	"db.t3.micro":   {VCPUs: 2, MemGB: 1, Family: "database"},
	"db.t3.small":   {VCPUs: 2, MemGB: 2, Family: "database"},
	"db.t3.medium":  {VCPUs: 2, MemGB: 4, Family: "database"},
	"db.t3.large":   {VCPUs: 2, MemGB: 8, Family: "database"},
	"db.m5.large":   {VCPUs: 2, MemGB: 8, Family: "database"},
	"db.m5.xlarge":  {VCPUs: 4, MemGB: 16, Family: "database"},
	"db.m5.2xlarge": {VCPUs: 8, MemGB: 32, Family: "database"},
	"db.m5.4xlarge": {VCPUs: 16, MemGB: 64, Family: "database"},
	"db.r5.large":   {VCPUs: 2, MemGB: 16, Family: "database"},
	"db.r5.xlarge":  {VCPUs: 4, MemGB: 32, Family: "database"},
	"db.r5.2xlarge": {VCPUs: 8, MemGB: 64, Family: "database"},
}

// GetInstanceSpec returns the spec for a given instance type, or nil if unknown.
func GetInstanceSpec(instanceType string) *InstanceSpec {
	spec, ok := instanceCatalog[instanceType]
	if !ok {
		return nil
	}
	return &spec
}
