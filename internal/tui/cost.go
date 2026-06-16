package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chxmxii/3a/internal/storage"
)

type costView struct {
	costs     []storage.CostEstimate
	resources []storage.Resource
}

func (v *costView) render(width int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("💰 Cost Analysis"))
	b.WriteString("\n\n")

	if len(v.costs) == 0 {
		b.WriteString(normalStyle.Render("  No cost estimates available. Run an assessment first."))
		return b.String()
	}

	// Calculate totals.
	totalCost := 0.0
	byCategory := make(map[string]float64)
	var idleResources []storage.CostEstimate
	var oversizedResources []storage.CostEstimate

	for _, c := range v.costs {
		if c.MonthlyCost != nil {
			totalCost += *c.MonthlyCost
			byCategory[c.Category] += *c.MonthlyCost
		}
		if c.IdleFlag {
			idleResources = append(idleResources, c)
		}
		if c.OversizedFlag {
			oversizedResources = append(oversizedResources, c)
		}
	}

	// Total.
	b.WriteString(headerStyle.Render("  Estimated Monthly Total"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  $%.2f/month (~$%.2f/year)\n\n", totalCost, totalCost*12))

	// By category.
	b.WriteString(headerStyle.Render("  Cost by Category"))
	b.WriteString("\n")
	for cat, cost := range byCategory {
		pct := 0.0
		if totalCost > 0 {
			pct = (cost / totalCost) * 100
		}
		b.WriteString(fmt.Sprintf("    %-15s $%8.2f  (%4.1f%%)\n", cat, cost, pct))
	}
	b.WriteString("\n")

	// Top 5 cost drivers.
	b.WriteString(headerStyle.Render("  Top Cost Drivers"))
	b.WriteString("\n")

	type costItem struct {
		resourceID string
		cost       float64
		resType    string
	}
	var items []costItem
	for _, c := range v.costs {
		if c.MonthlyCost != nil && *c.MonthlyCost > 0 {
			items = append(items, costItem{
				resourceID: c.ResourceID,
				cost:       *c.MonthlyCost,
				resType:    c.ResourceType,
			})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].cost > items[j].cost })
	if len(items) > 5 {
		items = items[:5]
	}

	// Build name lookup.
	nameMap := make(map[string]string)
	for _, r := range v.resources {
		nameMap[r.ResourceID] = r.Name
	}

	for i, item := range items {
		name := nameMap[item.resourceID]
		if name == "" {
			name = item.resourceID
			if len(name) > 30 {
				name = "..." + name[len(name)-27:]
			}
		}
		b.WriteString(fmt.Sprintf("    %d. %-30s %-15s $%.2f/mo\n", i+1, name, item.resType, item.cost))
	}
	b.WriteString("\n")

	// Optimization hints.
	if len(idleResources) > 0 || len(oversizedResources) > 0 {
		b.WriteString(headerStyle.Render("  Optimization Opportunities"))
		b.WriteString("\n")
		if len(idleResources) > 0 {
			b.WriteString(warnStyle.Render(fmt.Sprintf("    ⚠ %d potentially idle resource(s)\n", len(idleResources))))
		}
		if len(oversizedResources) > 0 {
			b.WriteString(warnStyle.Render(fmt.Sprintf("    ⚠ %d potentially oversized resource(s)\n", len(oversizedResources))))
		}
	}

	return b.String()
}
