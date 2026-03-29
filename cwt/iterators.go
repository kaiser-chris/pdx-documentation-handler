package cwt

import (
	"strings"

	"bahmut.de/pdx-documentation-manager/parser"
)

func PrintIterators(docs *parser.ScriptDocumentation) string {
	var builder = strings.Builder{}

	for _, iterator := range docs.IteratorDocumentation.Elements {
		builder.WriteString(printIterator("every", "effect", "effect_every_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("random", "effect", "effect_random_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("ordered", "effect", "effect_ordered_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("every", "arithmetic_operation", "formula_every_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("random", "arithmetic_operation", "formula_random_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("ordered", "arithmetic_operation", "formula_ordered_list_clause", iterator))
		builder.WriteString("\n")
		builder.WriteString(printIterator("any", "trigger", "trigger_any_list_clause", iterator))
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func printIterator(prefix, listType, clause string, iterator *parser.Iterator) string {
	var builder = strings.Builder{}
	builder.WriteString("### ")
	builder.WriteString(strings.Split(iterator.Description, "\n")[0])
	builder.WriteString("\n")
	builder.WriteString("## scopes = { ")
	builder.WriteString(printScopes(iterator.SupportedScopes))
	builder.WriteString(" }\n")
	if len(iterator.SupportedTargets) > 0 {
		builder.WriteString("## push_scope = ")
		builder.WriteString(iterator.SupportedTargets[0])
		builder.WriteString("\n")
	}
	builder.WriteString("alias[")
	builder.WriteString(listType)
	builder.WriteString(":")
	builder.WriteString(prefix)
	builder.WriteString("_")
	builder.WriteString(iterator.Name)
	builder.WriteString("] = single_alias_right[")
	builder.WriteString(clause)
	builder.WriteString("]")
	return builder.String()
}

func printScopes(scopes []string) string {
	var builder = strings.Builder{}
	for _, scope := range scopes {
		if scope == "none" {
			builder.WriteString("any")
		} else {
			builder.WriteString(scope)
		}
		builder.WriteString(" ")
	}
	return strings.TrimSpace(builder.String())
}
