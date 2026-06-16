package tui

import (
	"fmt"
	"strings"

	"github.com/chxmxii/3a/internal/storage"
)

type inventoryView struct {
	resources []storage.Resource
	filter    string
	cursor    int
	offset    int
}

func (v *inventoryView) render(width, height int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📦 Resource Inventory"))
	b.WriteString("\n")

	if v.filter != "" {
		b.WriteString(fmt.Sprintf("  Filter: %s\n", v.filter))
	}
	b.WriteString("\n")

	filtered := v.filteredResources()
	b.WriteString(fmt.Sprintf("  Showing %d of %d resources\n\n", len(filtered), len(v.resources)))

	// Table header.
	header := fmt.Sprintf("  %-20s %-40s %-15s %-15s", "TYPE", "NAME", "REGION", "ID (short)")
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString("  " + strings.Repeat("─", 90))
	b.WriteString("\n")

	// Calculate visible rows.
	maxRows := height - 10
	if maxRows < 5 {
		maxRows = 5
	}

	if v.offset > len(filtered)-maxRows {
		v.offset = max(0, len(filtered)-maxRows)
	}

	end := v.offset + maxRows
	if end > len(filtered) {
		end = len(filtered)
	}

	for i := v.offset; i < end; i++ {
		r := filtered[i]
		shortID := r.ResourceID
		if len(shortID) > 15 {
			shortID = "..." + shortID[len(shortID)-12:]
		}
		name := r.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}

		line := fmt.Sprintf("  %-20s %-40s %-15s %-15s", r.ResourceType, name, r.Region, shortID)
		if i == v.cursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(normalStyle.Render(line))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (v *inventoryView) filteredResources() []storage.Resource {
	if v.filter == "" {
		return v.resources
	}
	var filtered []storage.Resource
	lower := strings.ToLower(v.filter)
	for _, r := range v.resources {
		if strings.Contains(strings.ToLower(r.ResourceType), lower) ||
			strings.Contains(strings.ToLower(r.Name), lower) ||
			strings.Contains(strings.ToLower(r.Region), lower) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
