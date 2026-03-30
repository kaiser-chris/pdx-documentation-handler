package comparer

import (
	"strings"

	"bahmut.de/pdx-documentation-manager/parser"
)

func (compare *DataTypesCompareResult) Print(version string) string {
	var builder = strings.Builder{}

	builder.WriteString("# Data Type Documentation " + version + "\n")
	builder.WriteString("## Table of Contents\n")
	builder.WriteString(" * [Types](#types)\n")
	builder.WriteString(" * [Global Promotes](#global-promotes)\n")
	builder.WriteString("## Notes\n")
	builder.WriteString("This is just a very basic overview of added and removed data types.\n\nChanged elements are **not** mentioned here.\n")
	builder.WriteString("## Types\n")
	builder.WriteString(printTypes(compare.DataTypeResult))
	builder.WriteString("\n")
	builder.WriteString("## Global Promotes\n")
	builder.WriteString(printPromotes(compare.GlobalPromoteResult))
	builder.WriteString("\n")

	return builder.String()
}

func printTypes(compare *ElementResult[*parser.DataType]) string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Data Type"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
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

func printPromotes(compare *ElementResult[*parser.DataTypeFunction]) string {
	var builder = strings.Builder{}

	builder.WriteString(printTableHeader("Type", "Promote", "Return Type"))
	builder.WriteString("\n")
	for _, element := range compare.Added {
		builder.WriteString(printTableLine(
			"Added",
			printInlineCode(element.Name),
			printInlineCode(element.ReturnType),
		))
		builder.WriteString("\n")
	}
	for _, element := range compare.Removed {
		builder.WriteString(printTableLine(
			"Removed",
			printInlineCode(element.Name),
			printInlineCode(element.ReturnType),
		))
		builder.WriteString("\n")
	}

	return builder.String()
}
