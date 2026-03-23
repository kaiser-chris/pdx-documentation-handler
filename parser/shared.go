package parser

import (
	"os"
	"strings"

	"bahmut.de/pdx-documentation-manager/logging"
)

const (
	terminator                = ""
	listSeparator             = ", "
	prefixSmall               = "### "
	prefixMedium              = "## "
	prefixSupportedScopes     = "**Supported Scopes**: "
	prefixSupportedTargets    = "**Supported Targets**: "
	triggerTraitValue         = "Traits: <, <=, =, !=, >, >="
	triggerTraitBoolean       = "Traits: yes/no"
	eventTargetInput          = "Input Scopes: "
	eventTargetOutput         = "Output Scopes: "
	eventTargetTraitParameter = "Requires Data: yes"
	iteratorAny               = "any_"
	iteratorEvery             = "every_"
	iteratorOrdered           = "ordered_"
	iteratorRandom            = "random_"
	scopeSupportTriggers      = "Evaluate Triggers: "
	scopeSupportEffects       = "Execute Effects: "
	scopeSupportScopes        = "Change Scopes: "
	scopeSaveGameIdentifier   = "Save Token: "
	scopeSupportVariables     = "Stores Variables: "
	scriptBoolTrue            = "yes"
	scriptBoolFalse           = "no"
)

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
	return output
}

func parseScriptBool(text string) bool {
	if text == scriptBoolTrue {
		return true
	}
	if text == scriptBoolFalse {
		return false
	}
	logging.Warnf("Unrecognized bool: %s", text)
	return false
}

func exists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
