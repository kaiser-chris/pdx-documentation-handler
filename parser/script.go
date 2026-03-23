package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
)

const (
	scriptEffects             = "effects.log"
	scriptTriggers            = "triggers.log"
	scriptEventTargets        = "event_targets.log"
	scriptEventScopes         = "event_scopes.log"
	scriptOnActions           = "on_actions.log"
	scriptModifiers           = "modifiers.log"
	scriptCustomLocalizations = "custom_localization.log"
)

type Documentation struct {
	EffectDocumentation      *EffectDocumentation      `json:"effect-documentation"`
	TriggerDocumentation     *TriggerDocumentation     `json:"trigger-documentation"`
	EventTargetDocumentation *EventTargetDocumentation `json:"event-target-documentation"`
	IteratorDocumentation    *IteratorDocumentation    `json:"iterator-documentation"`
	ScopeDocumentation       *ScopeDocumentation       `json:"scope-documentation"`
}

type EffectDocumentation struct {
	File    string    `json:"file"`
	Effects []*Effect `json:"effects"`
}

type Effect struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	SupportedScopes  []string `json:"supported-scopes"`
	SupportedTargets []string `json:"supported-targets"`
}

type TriggerDocumentation struct {
	File     string     `json:"file"`
	Triggers []*Trigger `json:"triggers"`
}

type Trigger struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	SupportedScopes  []string `json:"supported-scopes"`
	SupportedTargets []string `json:"supported-targets"`
	Value            bool     `json:"is-value"`
	Boolean          bool     `json:"is-bool"`
}

type EventTargetDocumentation struct {
	File         string         `json:"file"`
	EventTargets []*EventTarget `json:"event-targets"`
}

type EventTarget struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	SupportedScopes []string `json:"supported-scopes"`
	OutputScope     string   `json:"output-scope"`
	Parameterized   bool     `json:"parameterized"`
}

type ScopeDocumentation struct {
	File   string   `json:"file"`
	Scopes []*Scope `json:"scopes"`
}

type Scope struct {
	Name              string `json:"name"`
	SaveIdentifier    string `json:"save-identifier"`
	SupportsTriggers  bool   `json:"supports-triggers"`
	SupportsEffects   bool   `json:"supports-effects"`
	SupportsVariables bool   `json:"supports-variables"`
	SupportsScopes    bool   `json:"supports-scopes"`
}

type IteratorDocumentation struct {
	Iterators []*Iterator `json:"iterators"`
}

type Iterator struct {
	Name     string   `json:"name"`
	Variants []string `json:"variants"`
}

func ParseScriptDocumentation(folder string) (*Documentation, error) {
	if !exists(folder) {
		return nil, fmt.Errorf("folder does not exist: %s", folder)
	}

	documentation := &Documentation{}

	effectFile := path.Join(folder, scriptEffects)
	if exists(effectFile) {
		effects, err := ParseEffectDocumentation(effectFile)
		if err != nil {
			return nil, err
		}
		documentation.EffectDocumentation = effects
	}

	triggerFile := path.Join(folder, scriptTriggers)
	if exists(triggerFile) {
		triggers, iterators, err := ParseTriggerDocumentation(triggerFile)
		if err != nil {
			return nil, err
		}
		documentation.TriggerDocumentation = triggers
		documentation.IteratorDocumentation = iterators
	}

	eventTargetFile := path.Join(folder, scriptEventTargets)
	if exists(eventTargetFile) {
		eventTargets, err := ParseEventTargetDocumentation(eventTargetFile)
		if err != nil {
			return nil, err
		}
		documentation.EventTargetDocumentation = eventTargets
	}

	scopeFile := path.Join(folder, scriptEventScopes)
	if exists(eventTargetFile) {
		scopes, err := ParseScopeDocumentation(scopeFile)
		if err != nil {
			return nil, err
		}
		documentation.ScopeDocumentation = scopes
	}

	return documentation, nil
}

func ParseEffectDocumentation(file string) (*EffectDocumentation, error) {
	documentation := &EffectDocumentation{
		File:    file,
		Effects: make([]*Effect, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var effect *Effect = nil
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixMedium) && effect != nil {
			return nil, fmt.Errorf("unterminated effect: %s, %s", effect.Name, line)
		}
		if strings.HasPrefix(line, prefixMedium) {
			if isIterator(cleanLine(line)) {
				continue
			}
			effect = &Effect{
				Name:             cleanLine(line),
				SupportedScopes:  make([]string, 0),
				SupportedTargets: make([]string, 0),
			}
			continue
		}
		if effect == nil {
			continue
		}
		if strings.TrimSpace(line) == terminator {
			// Finished Effect
			documentation.Effects = append(documentation.Effects, effect)
			effect = nil
			continue
		}
		if strings.HasPrefix(line, prefixSupportedScopes) {
			// Supported Scopes
			effect.SupportedScopes = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(effect.SupportedScopes)
			continue
		}
		if strings.HasPrefix(line, prefixSupportedTargets) {
			// Supported Targets
			effect.SupportedTargets = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(effect.SupportedTargets)
			continue
		}
		// Description
		effect.Description += "\n" + line
		effect.Description = strings.TrimPrefix(effect.Description, "\n")
	}

	return documentation, nil
}

func ParseTriggerDocumentation(file string) (*TriggerDocumentation, *IteratorDocumentation, error) {
	documentation := &TriggerDocumentation{
		File:     file,
		Triggers: make([]*Trigger, 0),
	}
	iterators := &IteratorDocumentation{
		Iterators: make([]*Iterator, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var trigger *Trigger = nil
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixMedium) && trigger != nil {
			// Unterminated Trigger
			return nil, nil, fmt.Errorf("unterminated effect: %s, %s", trigger.Name, line)
		}
		if strings.HasPrefix(line, prefixMedium) && isIterator(cleanLine(line)) {
			// Iterator
			name := strings.TrimPrefix(cleanLine(line), iteratorAny)
			iterators.Iterators = append(iterators.Iterators, createIterator(name))
			continue
		}
		if strings.HasPrefix(line, prefixMedium) {
			trigger = &Trigger{
				Name:             cleanLine(line),
				SupportedScopes:  make([]string, 0),
				SupportedTargets: make([]string, 0),
				Value:            false,
				Boolean:          false,
			}
			continue
		}
		if trigger == nil {
			continue
		}
		if cleanLine(line) == triggerTraitValue {
			trigger.Value = true
			continue
		}
		if cleanLine(line) == triggerTraitBoolean {
			trigger.Boolean = true
			continue
		}
		if strings.TrimSpace(line) == terminator {
			// Finished Trigger
			documentation.Triggers = append(documentation.Triggers, trigger)
			trigger = nil
			continue
		}
		if strings.HasPrefix(line, prefixSupportedScopes) {
			// Supported Scopes
			trigger.SupportedScopes = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(trigger.SupportedScopes)
			continue
		}
		if strings.HasPrefix(line, prefixSupportedTargets) {
			// Supported Targets
			trigger.SupportedTargets = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(trigger.SupportedTargets)
			continue
		}
		// Description
		trigger.Description += "\n" + line
		trigger.Description = strings.TrimPrefix(trigger.Description, "\n")
	}

	return documentation, iterators, nil
}

func ParseEventTargetDocumentation(file string) (*EventTargetDocumentation, error) {
	documentation := &EventTargetDocumentation{
		File:         file,
		EventTargets: make([]*EventTarget, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var eventTarget *EventTarget = nil
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixSmall) && eventTarget != nil {
			return nil, fmt.Errorf("unterminated effect: %s, %s", eventTarget.Name, line)
		}
		if strings.HasPrefix(line, prefixSmall) {
			eventTarget = &EventTarget{
				Name:            cleanLine(line),
				SupportedScopes: make([]string, 0),
				Parameterized:   false,
			}
			continue
		}
		if eventTarget == nil {
			continue
		}
		if cleanLine(line) == eventTargetTraitParameter {
			// Requires a parameter
			eventTarget.Parameterized = true
			continue
		}
		if strings.TrimSpace(line) == terminator {
			// Finished Event Target
			documentation.EventTargets = append(documentation.EventTargets, eventTarget)
			eventTarget = nil
			continue
		}
		if strings.HasPrefix(line, eventTargetInput) {
			// Supported Scopes
			eventTarget.SupportedScopes = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(eventTarget.SupportedScopes)
			continue
		}
		if strings.HasPrefix(line, eventTargetOutput) {
			// Output Scope
			eventTarget.OutputScope = cleanLine(line)
			continue
		}
		// Description
		eventTarget.Description += "\n" + line
		eventTarget.Description = strings.TrimPrefix(eventTarget.Description, "\n")
	}

	return documentation, nil
}

func ParseScopeDocumentation(file string) (*ScopeDocumentation, error) {
	documentation := &ScopeDocumentation{
		File:   file,
		Scopes: make([]*Scope, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var scope *Scope = nil
	for scanner.Scan() {
		line := scanner.Text()
		cleanLine := cleanLine(line)
		if isScopeName(line) && scope != nil {
			return nil, fmt.Errorf("unterminated scope: %s, %s", scope.Name, line)
		}
		if isScopeName(line) {
			scope = &Scope{
				Name:              strings.TrimSuffix(cleanLine, ":"),
				SupportsVariables: false,
				SupportsTriggers:  false,
				SupportsEffects:   false,
			}
			continue
		}
		if strings.TrimSpace(line) == terminator && scope != nil {
			// Finished Event Target
			documentation.Scopes = append(documentation.Scopes, scope)
			scope = nil
			continue
		}
		if scope == nil {
			continue
		}
		if strings.HasPrefix(line, scopeSupportTriggers) {
			scope.SupportsTriggers = parseScriptBool(cleanLine)
			continue
		}
		if strings.HasPrefix(line, scopeSupportEffects) {
			scope.SupportsEffects = parseScriptBool(cleanLine)
			continue
		}
		if strings.HasPrefix(line, scopeSupportScopes) {
			scope.SupportsScopes = parseScriptBool(cleanLine)
			continue
		}
		if strings.HasPrefix(line, scopeSupportVariables) {
			scope.SupportsVariables = parseScriptBool(cleanLine)
			continue
		}
		if strings.HasPrefix(line, scopeSaveGameIdentifier) {
			scope.SaveIdentifier = cleanLine
			continue
		}
	}

	return documentation, nil
}

func isIterator(name string) bool {
	return strings.HasPrefix(name, iteratorAny) ||
		strings.HasPrefix(name, iteratorEvery) ||
		strings.HasPrefix(name, iteratorOrdered) ||
		strings.HasPrefix(name, iteratorRandom)
}

func createIterator(name string) *Iterator {
	return &Iterator{
		Name: name,
		Variants: []string{
			iteratorAny + name,
			iteratorEvery + name,
			iteratorOrdered + name,
			iteratorRandom + name,
		},
	}
}

func isScopeName(line string) bool {
	return !strings.HasPrefix(line, scopeSupportTriggers) &&
		!strings.HasPrefix(line, scopeSupportEffects) &&
		!strings.HasPrefix(line, scopeSupportScopes) &&
		!strings.HasPrefix(line, scopeSaveGameIdentifier) &&
		!strings.HasPrefix(line, scopeSupportVariables) &&
		!(strings.TrimSpace(line) == terminator) &&
		!(line == "Scope Types:")
}
