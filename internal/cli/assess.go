package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/chxmxii/3a/internal/architecture"
	"github.com/chxmxii/3a/internal/assessment"
	awsrules "github.com/chxmxii/3a/internal/assessment/rules/aws"
	ocirules "github.com/chxmxii/3a/internal/assessment/rules/oci"
	"github.com/chxmxii/3a/internal/checklist"
	"github.com/chxmxii/3a/internal/config"
	"github.com/chxmxii/3a/internal/cost"
	"github.com/chxmxii/3a/internal/discovery"
	"github.com/chxmxii/3a/internal/provider/steampipe"
	"github.com/chxmxii/3a/internal/sizing"
	"github.com/chxmxii/3a/internal/storage"
	"github.com/chxmxii/3a/internal/tui"
)

func newAssessCmd() *cobra.Command {
	var connString string
	var noTUI bool

	cmd := &cobra.Command{
		Use:   "assess <profile>",
		Short: "Run a full assessment for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return runAssessment(profileName, connString, noTUI)
		},
	}

	cmd.Flags().StringVar(&connString, "steampipe-conn", "postgres://steampipe@localhost:9193/steampipe", "Steampipe connection string")
	cmd.Flags().BoolVar(&noTUI, "no-tui", false, "skip TUI and print summary to stdout")

	return cmd
}

func runAssessment(profileName, connString string, noTUI bool) error {
	ctx := context.Background()

	// Load config.
	cfgPath := config.DefaultConfigPath()
	cfg, err := config.Load(cfgPath)
	if err != nil {
		// If config doesn't exist, create minimal config.
		cfg = &config.Config{
			DBPath: resolveDBPath(getDBPath()),
			Profiles: []config.AccountProfile{
				{
					Name:     profileName,
					Provider: "aws",
					Regions:  []string{"us-east-1"},
				},
			},
		}
	}

	profile, err := config.GetProfile(cfg, profileName)
	if err != nil {
		return fmt.Errorf("profile error: %w", err)
	}

	// Open storage.
	dbFile := resolveDBPath(getDBPath())
	if cfg.DBPath != "" {
		dbFile = resolveDBPath(cfg.DBPath)
	}
	store, err := storage.Open(dbFile)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer store.Close()

	// Create assessment record.
	assessmentID := uuid.New().String()
	now := time.Now()
	a := &storage.Assessment{
		ID:       assessmentID,
		Profile:  profileName,
		Provider: profile.Provider,
		Status:   "in_progress",
		StartedAt: now,
		Regions:  profile.Regions,
	}
	if err := store.CreateAssessment(a); err != nil {
		return fmt.Errorf("creating assessment: %w", err)
	}

	fmt.Printf("🚀 Starting assessment %s for profile %s (%s)\n", assessmentID[:8], profileName, profile.Provider)

	// Step 1: Connect to Steampipe and discover resources.
	fmt.Println("📡 Connecting to Steampipe...")
	sp, err := steampipe.NewSteampipeProvider(connString, profile.Provider)
	if err != nil {
		return fmt.Errorf("creating steampipe provider: %w", err)
	}
	defer sp.Close()

	if err := sp.Authenticate(ctx); err != nil {
		return fmt.Errorf("connecting to steampipe: %w", err)
	}

	// Validate the profile can actually return data before running full discovery.
	fmt.Printf("🔑 Validating %s profile...\n", profile.Provider)
	if err := sp.ValidateProfile(ctx); err != nil {
		_ = store.UpdateAssessmentStatus(assessmentID, "failed", nil)
		return fmt.Errorf("profile validation failed:\n\n%w", err)
	}

	fmt.Println("🔍 Discovering resources...")
	engine := discovery.NewEngine(sp, store)
	summary, err := engine.Run(ctx, assessmentID, profile.Regions)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}
	fmt.Printf("   Found %d resources across %d regions\n", summary.TotalResources, len(summary.ByRegion))

	if summary.TotalResources == 0 {
		_ = store.UpdateAssessmentStatus(assessmentID, "failed", nil)
		errMsg := "discovery returned 0 resources\n\nPossible causes:\n"
		errMsg += "  • Steampipe plugin credentials may be expired or misconfigured\n"
		errMsg += "  • The account may have no resources in the queried regions\n"
		errMsg += "  • Steampipe connection may not have the right profile configured\n"
		errMsg += fmt.Sprintf("\nCheck: steampipe query \"SELECT count(*) FROM %s\"\n", canaryTableForProvider(profile.Provider))
		if len(summary.Errors) > 0 {
			errMsg += "\nDiscovery errors:\n"
			for _, e := range summary.Errors {
				errMsg += fmt.Sprintf("  • [%s/%s] %v\n", e.Service, e.Region, e.Err)
			}
		}
		return fmt.Errorf(errMsg)
	}

	// Step 2: Reconstruct architecture.
	fmt.Println("🏗️  Reconstructing architecture...")
	reconstructor := architecture.NewReconstructor(store, profile.Provider)
	if err := reconstructor.Reconstruct(assessmentID); err != nil {
		fmt.Printf("   ⚠ Architecture reconstruction: %v\n", err)
	}

	// Step 3: Run assessment rules.
	fmt.Println("🔐 Running security assessment...")
	var rules []assessment.Rule
	switch profile.Provider {
	case "aws":
		rules = awsrules.AllRules()
	case "oci":
		rules = ocirules.AllRules()
	}
	assessEngine := assessment.NewEngine(store, rules)
	if err := assessEngine.Run(ctx, assessmentID); err != nil {
		fmt.Printf("   ⚠ Assessment: %v\n", err)
	}

	findings, _ := store.GetFindingsByAssessment(assessmentID)
	fmt.Printf("   Found %d findings\n", len(findings))

	// Step 4: Sizing analysis.
	fmt.Println("📐 Analyzing sizing...")
	sizingAnalyzer := sizing.NewAnalyzer(store)
	sizingSummary, err := sizingAnalyzer.Analyze(assessmentID)
	if err != nil {
		fmt.Printf("   ⚠ Sizing: %v\n", err)
	} else {
		fmt.Printf("   %d vCPUs, %.1f GB memory\n", sizingSummary.TotalVCPUs, sizingSummary.TotalMemoryGB)
	}

	// Step 5: Cost estimation.
	fmt.Println("💰 Estimating costs...")
	costEstimator := cost.NewEstimator(store)
	costSummary, err := costEstimator.Estimate(assessmentID)
	if err != nil {
		fmt.Printf("   ⚠ Cost estimation: %v\n", err)
	} else {
		fmt.Printf("   Estimated: $%.2f/month\n", costSummary.TotalMonthlyCost)
	}

	// Step 6: Generate checklist.
	fmt.Println("✅ Generating checklist...")
	checkEngine := checklist.NewEngine(store)
	checkSummary, err := checkEngine.Generate(assessmentID)
	if err != nil {
		fmt.Printf("   ⚠ Checklist: %v\n", err)
	} else {
		fmt.Printf("   %d pass, %d fail, %d warn\n", checkSummary.PassCount, checkSummary.FailCount, checkSummary.WarnCount)
	}

	// Mark assessment as completed.
	completedAt := time.Now()
	_ = store.UpdateAssessmentStatus(assessmentID, "completed", &completedAt)

	fmt.Println("\n✅ Assessment complete!")

	if noTUI {
		return nil
	}

	// Launch TUI.
	fmt.Println("\nLaunching interactive view...")
	model := tui.NewModel(store, assessmentID)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

func resolveDBPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func canaryTableForProvider(providerType string) string {
	switch providerType {
	case "aws":
		return "aws_account"
	case "oci":
		return "oci_identity_compartment"
	default:
		return "unknown"
	}
}
