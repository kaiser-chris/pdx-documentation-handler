package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"bahmut.de/pdx-documentation-manager/util"
)

const (
	scriptEffects      = "effects.log"
	scriptTriggers     = "triggers.log"
	scriptEventTargets = "event_targets.log"
	scriptEventScopes  = "event_scopes.log"
	scriptOnActions    = "on_actions.log"
)

const (
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
	onActionSeparator         = "--------------------"
	onActionFromCode          = "From Code: "
	onActionExpectedScope     = "Expected Scope: "
)

type ScriptDocumentation struct {
	EffectDocumentation      *ElementDocumentation[Effect]      `json:"effect-documentation"`
	TriggerDocumentation     *ElementDocumentation[Trigger]     `json:"trigger-documentation"`
	EventTargetDocumentation *ElementDocumentation[EventTarget] `json:"event-target-documentation"`
	IteratorDocumentation    *ElementDocumentation[Iterator]    `json:"iterator-documentation"`
	ScopeDocumentation       *ElementDocumentation[Scope]       `json:"scope-documentation"`
	OnActionDocumentation    *ElementDocumentation[OnAction]    `json:"on-action-documentation"`
}

type ScriptElement interface {
	ElementName() string
}

type ScriptElements interface {
	Effect | Trigger | EventTarget | Iterator | Scope | OnAction
}

type ElementDocumentation[T ScriptElements] struct {
	File     string `json:"file"`
	Elements []*T   `json:"elements"`
}

type Effect struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	SupportedScopes  []string `json:"supported-scopes"`
	SupportedTargets []string `json:"supported-targets"`
}

func (e *Effect) ElementName() string {
	return e.Name
}

type Trigger struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	SupportedScopes  []string `json:"supported-scopes"`
	SupportedTargets []string `json:"supported-targets"`
	Value            bool     `json:"is-value"`
	Boolean          bool     `json:"is-bool"`
}

func (t *Trigger) ElementName() string {
	return t.Name
}

type Iterator struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Variants         []string `json:"variants"`
	SupportedScopes  []string `json:"supported-scopes"`
	SupportedTargets []string `json:"supported-targets"`
}

func (i *Iterator) ElementName() string {
	return i.Name
}

type EventTarget struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	SupportedScopes []string `json:"supported-scopes"`
	OutputScope     string   `json:"output-scope"`
	Parameterized   bool     `json:"parameterized"`
}

func (e *EventTarget) ElementName() string {
	return e.Name
}

type Scope struct {
	Name              string `json:"name"`
	SaveIdentifier    string `json:"save-identifier"`
	SupportsTriggers  bool   `json:"supports-triggers"`
	SupportsEffects   bool   `json:"supports-effects"`
	SupportsVariables bool   `json:"supports-variables"`
	SupportsScopes    bool   `json:"supports-scopes"`
}

func (s *Scope) ElementName() string {
	return s.Name
}

type OnAction struct {
	Name     string `json:"name"`
	FromCode bool   `json:"from-code"`
	Scope    string `json:"scope"`
}

func (o *OnAction) ElementName() string {
	return o.Name
}

func ParseScriptDocumentation(folder string) (*ScriptDocumentation, error) {
	if !util.Exists(folder) {
		return nil, fmt.Errorf("script documentation folder does not exist: %s", folder)
	}

	documentation := &ScriptDocumentation{}

	effectFile := path.Join(folder, scriptEffects)
	if util.Exists(effectFile) {
		effects, err := ParseEffectDocumentation(effectFile)
		if err != nil {
			return nil, err
		}
		documentation.EffectDocumentation = effects
	} else {
		return nil, fmt.Errorf("effect documentation does not exist: %s", effectFile)
	}

	triggerFile := path.Join(folder, scriptTriggers)
	if util.Exists(triggerFile) {
		triggers, err := ParseTriggerDocumentation(triggerFile)
		if err != nil {
			return nil, err
		}
		documentation.TriggerDocumentation = triggers
		iterators, err := ParseIteratorDocumentation(triggerFile)
		if err != nil {
			return nil, err
		}
		documentation.IteratorDocumentation = iterators
	} else {
		return nil, fmt.Errorf("trigger documentation does not exist: %s", effectFile)
	}

	eventTargetFile := path.Join(folder, scriptEventTargets)
	if util.Exists(eventTargetFile) {
		elements, err := ParseEventTargetDocumentation(eventTargetFile)
		if err != nil {
			return nil, err
		}
		documentation.EventTargetDocumentation = elements
	} else {
		return nil, fmt.Errorf("event target documentation does not exist: %s", effectFile)
	}

	scopeFile := path.Join(folder, scriptEventScopes)
	if util.Exists(eventTargetFile) {
		elements, err := ParseScopeDocumentation(scopeFile)
		if err != nil {
			return nil, err
		}
		documentation.ScopeDocumentation = elements
	} else {
		return nil, fmt.Errorf("scope documentation does not exist: %s", effectFile)
	}

	onActionFile := path.Join(folder, scriptOnActions)
	if util.Exists(eventTargetFile) {
		elements, err := ParseOnActionDocumentation(onActionFile)
		if err != nil {
			return nil, err
		}
		documentation.OnActionDocumentation = elements
	} else {
		return nil, fmt.Errorf("on action documentation does not exist: %s", effectFile)
	}

	return documentation, nil
}

func ParseEffectDocumentation(file string) (*ElementDocumentation[Effect], error) {
	documentation := &ElementDocumentation[Effect]{
		File:     file,
		Elements: make([]*Effect, 0),
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
			documentation.Elements = append(documentation.Elements, effect)
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

func ParseTriggerDocumentation(file string) (*ElementDocumentation[Trigger], error) {
	documentation := &ElementDocumentation[Trigger]{
		File:     file,
		Elements: make([]*Trigger, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var trigger *Trigger = nil
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixMedium) && trigger != nil {
			// Unterminated Trigger
			return nil, fmt.Errorf("unterminated trigger: %s, %s", trigger.Name, line)
		}
		if strings.HasPrefix(line, prefixMedium) && isIterator(cleanLine(line)) {
			// Iterator
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
			documentation.Elements = append(documentation.Elements, trigger)
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

	return documentation, nil
}

func ParseIteratorDocumentation(file string) (*ElementDocumentation[Iterator], error) {
	iterators := &ElementDocumentation[Iterator]{
		File:     file,
		Elements: make([]*Iterator, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var iterator *Iterator = nil
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixMedium) && iterator != nil {
			// Unterminated Trigger
			return nil, fmt.Errorf("unterminated iterator: %s, %s", iterator.Name, line)
		}
		if strings.HasPrefix(line, prefixMedium) && !isIterator(cleanLine(line)) {
			// Normal Trigger
			continue
		}
		if strings.HasPrefix(line, prefixMedium) {
			iterator = &Iterator{
				Name:             strings.TrimPrefix(cleanLine(line), "any_"),
				SupportedScopes:  make([]string, 0),
				SupportedTargets: make([]string, 0),
			}
			continue
		}
		if iterator == nil {
			continue
		}
		if strings.TrimSpace(line) == terminator {
			// Finished Trigger
			iterators.Elements = append(iterators.Elements, iterator)
			iterator = nil
			continue
		}
		if strings.HasPrefix(line, prefixSupportedScopes) {
			// Supported Scopes
			iterator.SupportedScopes = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(iterator.SupportedScopes)
			continue
		}
		if strings.HasPrefix(line, prefixSupportedTargets) {
			// Supported Targets
			iterator.SupportedTargets = strings.Split(cleanLine(line), listSeparator)
			slices.Sort(iterator.SupportedTargets)
			continue
		}
		// Description
		iterator.Description += "\n" + line
		iterator.Description = strings.TrimPrefix(iterator.Description, "\n")
	}

	return iterators, nil
}

func ParseEventTargetDocumentation(file string) (*ElementDocumentation[EventTarget], error) {
	documentation := &ElementDocumentation[EventTarget]{
		File:     file,
		Elements: make([]*EventTarget, 0),
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
			documentation.Elements = append(documentation.Elements, eventTarget)
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

func ParseScopeDocumentation(file string) (*ElementDocumentation[Scope], error) {
	documentation := &ElementDocumentation[Scope]{
		File:     file,
		Elements: make([]*Scope, 0),
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
			documentation.Elements = append(documentation.Elements, scope)
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

func ParseOnActionDocumentation(file string) (*ElementDocumentation[OnAction], error) {
	documentation := &ElementDocumentation[OnAction]{
		File:     file,
		Elements: make([]*OnAction, 0),
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var element *OnAction = nil
	for scanner.Scan() {
		line := scanner.Text()
		cleanLine := cleanLine(line)
		if isOnActionName(line) && element != nil {
			return nil, fmt.Errorf("unterminated on action: %s, %s", element.Name, line)
		}
		if isOnActionName(line) {
			element = &OnAction{
				Name:     strings.TrimSuffix(cleanLine, ":"),
				FromCode: false,
			}
			continue
		}
		if strings.TrimSpace(line) == terminator && element != nil {
			// Finished On Action
			documentation.Elements = append(documentation.Elements, element)
			element = nil
			continue
		}
		if element == nil {
			continue
		}
		if strings.HasPrefix(line, onActionFromCode) {
			element.FromCode = parseScriptBool(cleanLine)
			continue
		}
		if strings.HasPrefix(line, onActionExpectedScope) {
			element.Scope = cleanLine
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

func isScopeName(line string) bool {
	return !strings.HasPrefix(line, scopeSupportTriggers) &&
		!strings.HasPrefix(line, scopeSupportEffects) &&
		!strings.HasPrefix(line, scopeSupportScopes) &&
		!strings.HasPrefix(line, scopeSaveGameIdentifier) &&
		!strings.HasPrefix(line, scopeSupportVariables) &&
		!(strings.TrimSpace(line) == terminator) &&
		!(line == "Scope Types:")
}

func isOnActionName(line string) bool {
	return !strings.HasPrefix(line, onActionFromCode) &&
		!strings.HasPrefix(line, onActionExpectedScope) &&
		!strings.HasPrefix(line, onActionSeparator) &&
		!(strings.TrimSpace(line) == terminator) &&
		!(line == "On Action Documentation:")
}
