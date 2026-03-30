package parser

import (
	"strings"

	"bahmut.de/pdx-documentation-manager/logging"
)

const (
	terminator = ""
)

type DocumentationElement interface {
	ElementName() string
}

func cleanLine(line string) string {
	output := strings.TrimSpace(line)
	output = strings.TrimPrefix(output, prefixSmall)
	output = strings.TrimPrefix(output, prefixMedium)
	output = strings.TrimPrefix(output, prefixSupportedScopes)
	output = strings.TrimPrefix(output, prefixSupportedTargets)
	output = strings.TrimPrefix(output, eventTargetInput)
	output = strings.TrimPrefix(output, eventTargetOutput)
	output = strings.TrimPrefix(output, scopeSupportTriggers)
	output = strings.TrimPrefix(output, scopeSupportEffects)
	output = strings.TrimPrefix(output, scopeSupportScopes)
	output = strings.TrimPrefix(output, scopeSaveGameIdentifier)
	output = strings.TrimPrefix(output, scopeSupportVariables)
	output = strings.TrimPrefix(output, onActionFromCode)
	output = strings.TrimPrefix(output, onActionExpectedScope)
	output = strings.TrimPrefix(output, onActionSeparator)
	return output
}

func parseScriptBool(text string) bool {
	lower := strings.ToLower(text)
	if lower == scriptBoolTrue {
		return true
	}
	if lower == scriptBoolFalse {
		return false
	}
	logging.Warnf("Unrecognized bool: %s", text)
	return false
}
