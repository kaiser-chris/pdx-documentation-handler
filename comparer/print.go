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
	builder.WriteString(" * [On Actions](#on-actions)\n")
	builder.WriteString("## Notes\n")
	builder.WriteString(" * **Changed** means the description, scopes or anything related to the documentation for this element has changed\n")
	builder.WriteString(" * The list of iterators do **not** include generated geographic region based iterators\n")
	builder.WriteString(" * The on action scope is based on the script documentation, for more information see the `common/on_actions` directory\n")
	builder.WriteString("## Scopes\n")
	builder.WriteString(printScopes(compare.ScopeResult))
	builder.WriteString("\n")
	builder.WriteString("## Effects\n")
	builder.WriteString(printEffects(compare.EffectResult))
	builder.WriteString("\n")
	builder.WriteString("## Triggers\n")
	builder.WriteString(printTriggers(compare.TriggerResult))
	builder.WriteString("\n")
	builder.WriteString("## Event Targets\n")
	builder.WriteString(printEventTargets(compare.EventTargetResult))
	builder.WriteString("\n")
	builder.WriteString("## Iterators\n")
	builder.WriteString(printIterators(compare.IteratorResult))
	builder.WriteString("## On Actions\n")
	builder.WriteString(printOnActions(compare.OnActionResult))
	builder.WriteString("\n")

	return builder.String()
}

func printEffects(compare *ElementResult[*parser.Effect]) string {
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

func printTriggers(compare *ElementResult[*parser.Trigger]) string {
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

func printEventTargets(compare *ElementResult[*parser.EventTarget]) string {
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

func printIterators(compare *ElementResult[*parser.Iterator]) string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Iterator"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode("{any|every|ordered|random}_"+element.Name),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		if isGeographicRegionIterator(element.Name) {
			continue
		}
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode("{any|every|ordered|random}_"+element.Name),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func printScopes(compare *ElementResult[*parser.Scope]) string {
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

func printOnActions(compare *ElementResult[*parser.OnAction]) string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "On Action", "Scope"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		if !element.FromCode {
			continue
		}
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printInlineCode(element.Scope),
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
			printInlineCode(element.Scope),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}

func printTableLine(elements ...string) string {
	var builder = strings.Builder{}
	builder.WriteString("| ")
	for _, element := range elements {
		builder.WriteString(clean(element) + " | ")
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

func clean(text string) string {
	result := text
	result = strings.ReplaceAll(result, "|", "\\|")
	result = strings.ReplaceAll(result, "[", "\\[")
	result = strings.ReplaceAll(result, "]", "\\]")
	result = strings.ReplaceAll(result, "<", "\\<")
	result = strings.ReplaceAll(result, ">", "\\>")
	return result
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
