package comparer

import (
	"slices"

	"bahmut.de/pdx-documentation-manager/parser"
)

type CompareResult struct {
	EffectResult      *EffectResult      `json:"effect-result"`
	TriggerResult     *TriggerResult     `json:"trigger-result"`
	EventTargetResult *EventTargetResult `json:"event-target-result"`
	IteratorResult    *IteratorResult    `json:"iterator-result"`
	ScopeResult       *ScopeResult       `json:"scope-result"`
}

type EffectResult struct {
	Added    []*parser.Effect `json:"added"`
	Changed  []*parser.Effect `json:"changed"`
	Removed  []*parser.Effect `json:"removed"`
	Existing []*parser.Effect `json:"existing"`
}

type TriggerResult struct {
	Added    []*parser.Trigger `json:"added"`
	Changed  []*parser.Trigger `json:"changed"`
	Removed  []*parser.Trigger `json:"removed"`
	Existing []*parser.Trigger `json:"existing"`
}

type EventTargetResult struct {
	Added    []*parser.EventTarget `json:"added"`
	Changed  []*parser.EventTarget `json:"changed"`
	Removed  []*parser.EventTarget `json:"removed"`
	Existing []*parser.EventTarget `json:"existing"`
}

type IteratorResult struct {
	Added    []*parser.Iterator `json:"added"`
	Removed  []*parser.Iterator `json:"removed"`
	Existing []*parser.Iterator `json:"existing"`
}

type ScopeResult struct {
	Added    []*parser.Scope `json:"added"`
	Removed  []*parser.Scope `json:"removed"`
	Existing []*parser.Scope `json:"existing"`
}

func Compare(old *parser.Documentation, new *parser.Documentation) *CompareResult {
	result := &CompareResult{}

	result.EffectResult = compareEffects(old.EffectDocumentation.Effects, new.EffectDocumentation.Effects)
	result.TriggerResult = compareTriggers(old.TriggerDocumentation.Triggers, new.TriggerDocumentation.Triggers)
	result.EventTargetResult = compareEventTargets(old.EventTargetDocumentation.EventTargets, new.EventTargetDocumentation.EventTargets)
	result.IteratorResult = compareIterators(old.IteratorDocumentation.Iterators, new.IteratorDocumentation.Iterators)
	result.ScopeResult = compareScopes(old.ScopeDocumentation.Scopes, new.ScopeDocumentation.Scopes)

	return result
}

func compareEffects(old []*parser.Effect, new []*parser.Effect) *EffectResult {
	result := &EffectResult{
		Added:    make([]*parser.Effect, 0),
		Changed:  make([]*parser.Effect, 0),
		Removed:  make([]*parser.Effect, 0),
		Existing: make([]*parser.Effect, 0),
	}

	for _, effect := range new {
		compare := findEffect(effect.Name, old)
		if compare == nil {
			result.Added = append(result.Added, effect)
			continue
		}
		if effect.Description != compare.Description {
			result.Changed = append(result.Changed, effect)
			continue
		}
		if slices.Compare(effect.SupportedScopes, compare.SupportedScopes) != 0 {
			result.Changed = append(result.Changed, effect)
			continue
		}
		if slices.Compare(effect.SupportedTargets, compare.SupportedTargets) != 0 {
			result.Changed = append(result.Changed, effect)
			continue
		}
		result.Existing = append(result.Existing, effect)
	}

	for _, effect := range old {
		compare := findEffect(effect.Name, new)
		if compare == nil {
			result.Removed = append(result.Removed, effect)
			continue
		}
	}

	return result
}

func compareTriggers(old []*parser.Trigger, new []*parser.Trigger) *TriggerResult {
	result := &TriggerResult{
		Added:    make([]*parser.Trigger, 0),
		Changed:  make([]*parser.Trigger, 0),
		Removed:  make([]*parser.Trigger, 0),
		Existing: make([]*parser.Trigger, 0),
	}

	for _, trigger := range new {
		compare := findTrigger(trigger.Name, old)
		if compare == nil {
			result.Added = append(result.Added, trigger)
			continue
		}
		if trigger.Description != compare.Description {
			result.Changed = append(result.Changed, trigger)
			continue
		}
		if trigger.Boolean != compare.Boolean {
			result.Changed = append(result.Changed, trigger)
			continue
		}
		if trigger.Value != compare.Value {
			result.Changed = append(result.Changed, trigger)
			continue
		}
		if slices.Compare(trigger.SupportedScopes, compare.SupportedScopes) != 0 {
			result.Changed = append(result.Changed, trigger)
			continue
		}
		if slices.Compare(trigger.SupportedTargets, compare.SupportedTargets) != 0 {
			result.Changed = append(result.Changed, trigger)
			continue
		}
		result.Existing = append(result.Existing, trigger)
	}

	for _, trigger := range old {
		compare := findTrigger(trigger.Name, new)
		if compare == nil {
			result.Removed = append(result.Removed, trigger)
			continue
		}
	}

	return result
}

func compareEventTargets(old []*parser.EventTarget, new []*parser.EventTarget) *EventTargetResult {
	result := &EventTargetResult{
		Added:    make([]*parser.EventTarget, 0),
		Changed:  make([]*parser.EventTarget, 0),
		Removed:  make([]*parser.EventTarget, 0),
		Existing: make([]*parser.EventTarget, 0),
	}

	for _, eventTarget := range new {
		compare := findEventTarget(eventTarget.Name, old)
		if compare == nil {
			result.Added = append(result.Added, eventTarget)
			continue
		}
		if eventTarget.Description != compare.Description {
			result.Changed = append(result.Changed, eventTarget)
			continue
		}
		if eventTarget.Parameterized != compare.Parameterized {
			result.Changed = append(result.Changed, eventTarget)
			continue
		}
		if slices.Compare(eventTarget.SupportedScopes, compare.SupportedScopes) != 0 {
			result.Changed = append(result.Changed, eventTarget)
			continue
		}
		if eventTarget.OutputScope != compare.OutputScope {
			result.Changed = append(result.Changed, eventTarget)
			continue
		}
		result.Existing = append(result.Existing, eventTarget)
	}

	for _, eventTarget := range old {
		compare := findEventTarget(eventTarget.Name, new)
		if compare == nil {
			result.Removed = append(result.Removed, eventTarget)
			continue
		}
	}

	return result
}

func compareIterators(old []*parser.Iterator, new []*parser.Iterator) *IteratorResult {
	result := &IteratorResult{
		Added:    make([]*parser.Iterator, 0),
		Removed:  make([]*parser.Iterator, 0),
		Existing: make([]*parser.Iterator, 0),
	}

	for _, iterator := range new {
		compare := findIterator(iterator.Name, old)
		if compare == nil {
			result.Added = append(result.Added, iterator)
			continue
		}
		result.Existing = append(result.Existing, iterator)
	}

	for _, iterator := range old {
		compare := findIterator(iterator.Name, new)
		if compare == nil {
			result.Removed = append(result.Removed, iterator)
			continue
		}
	}

	return result
}

func compareScopes(old []*parser.Scope, new []*parser.Scope) *ScopeResult {
	result := &ScopeResult{
		Added:    make([]*parser.Scope, 0),
		Removed:  make([]*parser.Scope, 0),
		Existing: make([]*parser.Scope, 0),
	}

	for _, scope := range new {
		compare := findScope(scope.Name, old)
		if compare == nil {
			result.Added = append(result.Added, scope)
			continue
		}
		result.Existing = append(result.Existing, scope)
	}

	for _, scope := range old {
		compare := findScope(scope.Name, new)
		if compare == nil {
			result.Removed = append(result.Removed, scope)
			continue
		}
	}

	return result
}

func findTrigger(name string, elements []*parser.Trigger) *parser.Trigger {
	for _, element := range elements {
		if element.Name == name {
			return element
		}
	}
	return nil
}

func findEffect(name string, elements []*parser.Effect) *parser.Effect {
	for _, element := range elements {
		if element.Name == name {
			return element
		}
	}
	return nil
}

func findEventTarget(name string, elements []*parser.EventTarget) *parser.EventTarget {
	for _, element := range elements {
		if element.Name == name {
			return element
		}
	}
	return nil
}

func findIterator(name string, elements []*parser.Iterator) *parser.Iterator {
	for _, element := range elements {
		if element.Name == name {
			return element
		}
	}
	return nil
}

func findScope(name string, elements []*parser.Scope) *parser.Scope {
	for _, element := range elements {
		if element.Name == name {
			return element
		}
	}
	return nil
}
