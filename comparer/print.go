package comparer

import (
	"strings"

	"bahmut.de/pdx-documentation-manager/parser"
)

func (compare *CompareResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString("# Script Documentation\n")
	builder.WriteString("## Table of Contents\n")
	builder.WriteString(" * [Scopes](#scopes)\n")
	builder.WriteString(" * [Effects](#effects)\n")
	builder.WriteString(" * [Triggers](#triggers)\n")
	builder.WriteString(" * [Event Targets](#event-targets)\n")
	builder.WriteString(" * [Iterators](#iterators)\n")
	builder.WriteString("## Notes\n")
	builder.WriteString(" * **Changed** means the description, scopes or anything related to the documentation for this element has changed\n")
	builder.WriteString(" * The list of iterators do **not** include generated geographic region based iterators\n")
	builder.WriteString("## Scopes\n")
	builder.WriteString(compare.ScopeResult.Print())
	builder.WriteString("\n")
	builder.WriteString("## Effects\n")
	builder.WriteString(compare.EffectResult.Print())
	builder.WriteString("\n")
	builder.WriteString("## Triggers\n")
	builder.WriteString(compare.TriggerResult.Print())
	builder.WriteString("\n")
	builder.WriteString("## Event Targets\n")
	builder.WriteString(compare.EventTargetResult.Print())
	builder.WriteString("\n")
	builder.WriteString("## Iterators\n")
	builder.WriteString(compare.IteratorResult.Print())
	builder.WriteString("\n")

	return builder.String()
}

func (compare *EffectResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Effect", "Description"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Changed {
		builder.WriteString(printTableLine(
			"Changed",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (compare *TriggerResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Trigger", "Trait", "Description"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printTriggerTraits(element),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Changed {
		builder.WriteString(printTableLine(
			"Changed",
			printInlineCode(element.Name),
			printTriggerTraits(element),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode(element.Name),
			printTriggerTraits(element),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (compare *EventTargetResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Event Target", "Description"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Changed {
		builder.WriteString(printTableLine(
			"Changed",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode(element.Name),
			printFirstLine(element.Description),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (compare *IteratorResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Iterator"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode("{any\\|every\\|ordered\\|random}_"+element.Name),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode("{any\\|every\\|ordered\\|random}_"+element.Name),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (compare *ScopeResult) Print() string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Scope", "Supports Variables", "Supports Effects", "Supports Triggers", "Save Game Identifier"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printBool(element.SupportsVariables),
			printBool(element.SupportsEffects),
			printBool(element.SupportsTriggers),
			printInlineCode(element.SaveIdentifier),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode(element.Name),
			printBool(element.SupportsVariables),
			printBool(element.SupportsEffects),
			printBool(element.SupportsTriggers),
			printInlineCode(element.SaveIdentifier),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func printTableLine(elements ...string) string {
	var builder = strings.Builder{}
	builder.WriteString("| ")
	for _, element := range elements {
		builder.WriteString(element + " | ")
	}
	return strings.TrimSpace(builder.String())
}

func printTableHeader(elements ...string) string {
	var builder = strings.Builder{}
	builder.WriteString(printTableLine(elements...))
	builder.WriteString("\n")
	builder.WriteString("|")
	for range elements {
		builder.WriteString("--|")
	}
	return builder.String()
}

func printFirstLine(block string) string {
	return strings.Split(block, "\n")[0]
}

func printInlineCode(code string) string {
	return "`" + code + "`"
}

func printTriggerTraits(element *parser.Trigger) string {
	if element.Value {
		return "Value"
	}
	if element.Boolean {
		return "Boolean"
	}
	return " - "
}

func printBool(element bool) string {
	if element {
		return "True"
	}
	return "False"
}

func isGeographicRegionIterator(name string) bool {
	return strings.HasPrefix(name, "country_in_") ||
		strings.HasPrefix(name, "province_in_") ||
		strings.HasPrefix(name, "state_in_") ||
		strings.HasPrefix(name, "state_region_in_") ||
		strings.HasPrefix(name, "strategic_region_in_")
}
