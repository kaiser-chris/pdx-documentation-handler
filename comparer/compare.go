package comparer

import (
	"slices"

	"bahmut.de/pdx-documentation-manager/parser"
)

type ScriptDocumentationCompareResult struct {
	EffectResult      *ElementResult[*parser.Effect]      `json:"effect-result"`
	TriggerResult     *ElementResult[*parser.Trigger]     `json:"trigger-result"`
	EventTargetResult *ElementResult[*parser.EventTarget] `json:"event-target-result"`
	IteratorResult    *ElementResult[*parser.Iterator]    `json:"iterator-result"`
	ScopeResult       *ElementResult[*parser.Scope]       `json:"scope-result"`
	OnActionResult    *ElementResult[*parser.OnAction]    `json:"on-action-result"`
}

type DataTypesCompareResult struct {
	DataTypeResult      *ElementResult[*parser.DataType]         `json:"data-type-result"`
	GlobalPromoteResult *ElementResult[*parser.DataTypeFunction] `json:"global-promote-result"`
}

type ElementResult[T parser.DocumentationElement] struct {
	Added    []T `json:"added"`
	Changed  []T `json:"changed"`
	Removed  []T `json:"removed"`
	Existing []T `json:"existing"`
}

func CompareScriptDocumentation(old *parser.ScriptDocumentation, new *parser.ScriptDocumentation) *ScriptDocumentationCompareResult {
	result := &ScriptDocumentationCompareResult{}

	result.EffectResult = compareEffects(old.EffectDocumentation.Elements, new.EffectDocumentation.Elements)
	result.TriggerResult = compareTriggers(old.TriggerDocumentation.Elements, new.TriggerDocumentation.Elements)
	result.EventTargetResult = compareBasic(old.EventTargetDocumentation.Elements, new.EventTargetDocumentation.Elements)
	result.IteratorResult = compareBasic(old.IteratorDocumentation.Elements, new.IteratorDocumentation.Elements)
	result.ScopeResult = compareBasic(old.ScopeDocumentation.Elements, new.ScopeDocumentation.Elements)
	result.OnActionResult = compareBasic(old.OnActionDocumentation.Elements, new.OnActionDocumentation.Elements)

	return result
}

func CompareDataTypes(old *parser.DataTypeDocumentation, new *parser.DataTypeDocumentation) *DataTypesCompareResult {
	result := &DataTypesCompareResult{}

	result.DataTypeResult = compareBasic(old.DataTypes, new.DataTypes)
	result.GlobalPromoteResult = compareBasic(old.GlobalPromotes, new.GlobalPromotes)

	return result
}

func compareEffects(old []*parser.Effect, new []*parser.Effect) *ElementResult[*parser.Effect] {
	result := &ElementResult[*parser.Effect]{
		Added:    make([]*parser.Effect, 0),
		Changed:  make([]*parser.Effect, 0),
		Removed:  make([]*parser.Effect, 0),
		Existing: make([]*parser.Effect, 0),
	}

	for _, effect := range new {
		compare, found := findElement(effect, old)
		if !found {
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
		_, found := findElement(effect, new)
		if !found {
			result.Removed = append(result.Removed, effect)
			continue
		}
	}

	return result
}

func compareTriggers(old []*parser.Trigger, new []*parser.Trigger) *ElementResult[*parser.Trigger] {
	result := &ElementResult[*parser.Trigger]{
		Added:    make([]*parser.Trigger, 0),
		Changed:  make([]*parser.Trigger, 0),
		Removed:  make([]*parser.Trigger, 0),
		Existing: make([]*parser.Trigger, 0),
	}

	for _, trigger := range new {
		compare, found := findElement(trigger, old)
		if !found {
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
		_, found := findElement(trigger, new)
		if !found {
			result.Removed = append(result.Removed, trigger)
			continue
		}
	}

	return result
}

func compareBasic[T parser.DocumentationElement](old []T, new []T) *ElementResult[T] {
	result := &ElementResult[T]{
		Added:    make([]T, 0),
		Changed:  make([]T, 0),
		Removed:  make([]T, 0),
		Existing: make([]T, 0),
	}

	for _, element := range new {
		_, found := findElement(element, old)
		if !found {
			result.Added = append(result.Added, element)
			continue
		}
		result.Existing = append(result.Existing, element)
	}

	for _, element := range old {
		_, found := findElement(element, new)
		if !found {
			result.Removed = append(result.Removed, element)
			continue
		}
	}

	return result
}

func findElement[T parser.DocumentationElement](base T, elements []T) (T, bool) {
	for _, element := range elements {
		if element.ElementName() == base.ElementName() {
			return element, true
		}
	}
	return elements[0], false
}
