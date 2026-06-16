package tui

import (
	"fmt"
	"strings"

	"github.com/chxmxii/3a/internal/storage"
)

type architectureView struct {
	resources     []storage.Resource
	relationships []storage.Relationship
}

func (v *architectureView) render(width int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🏗️  Architecture View"))
	b.WriteString("\n\n")

	if len(v.relationships) == 0 {
		b.WriteString(normalStyle.Render("  No relationships found. Run an assessment first."))
		return b.String()
	}

	b.WriteString(fmt.Sprintf("  %d relationships discovered\n\n", len(v.relationships)))

	// Build a resource name lookup.
	nameMap := make(map[string]string)
	typeMap := make(map[string]string)
	for _, r := range v.resources {
		display := r.Name
		if display == "" {
			display = r.ResourceID
		}
		nameMap[r.ResourceID] = display
		typeMap[r.ResourceID] = r.ResourceType
	}

	// Group relationships by source.
	bySource := make(map[string][]storage.Relationship)
	for _, rel := range v.relationships {
		bySource[rel.SourceID] = append(bySource[rel.SourceID], rel)
	}

	// Find root resources (those that are never targets).
	targetSet := make(map[string]bool)
	for _, rel := range v.relationships {
		targetSet[rel.TargetID] = true
	}

	var roots []string
	for source := range bySource {
		if !targetSet[source] {
			roots = append(roots, source)
		}
	}

	// If no clear roots, use all sources.
	if len(roots) == 0 {
		for source := range bySource {
			roots = append(roots, source)
		}
	}

	// Render tree (limit depth to avoid cycles).
	rendered := make(map[string]bool)
	for _, root := range roots {
		if rendered[root] {
			continue
		}
		v.renderTree(&b, root, "", true, bySource, nameMap, typeMap, rendered, 0)
	}

	return b.String()
}

func (v *architectureView) renderTree(b *strings.Builder, resourceID, prefix string, isLast bool, bySource map[string][]storage.Relationship, nameMap, typeMap map[string]string, rendered map[string]bool, depth int) {
	if depth > 4 || rendered[resourceID] {
		return
	}
	rendered[resourceID] = true

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	name := nameMap[resourceID]
	if name == "" {
		name = resourceID
		if len(name) > 30 {
			name = "..." + name[len(name)-27:]
		}
	}
	rType := typeMap[resourceID]

	line := fmt.Sprintf("%s%s[%s] %s", prefix, connector, rType, name)
	b.WriteString(normalStyle.Render(line))
	b.WriteString("\n")

	children := bySource[resourceID]
	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, rel := range children {
		isChildLast := i == len(children)-1
		v.renderTree(b, rel.TargetID, childPrefix, isChildLast, bySource, nameMap, typeMap, rendered, depth+1)
	}
}
