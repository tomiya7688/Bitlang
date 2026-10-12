package bitlang

import (
	"fmt"
	"sort"
)

//	{
//	  責務: [
//	    validatePropertySpecification: 照合時に検査が抜けないようプロパティ仕様全体を検証する
//	  ]
//	  処理: [
//	    1: バージョンと軸の存在を確認する
//	    2: 軸名・状態・適用対象の空値と重複を拒否する
//	    3: 全軸をまたぐ状態名の曖昧さを拒否する
//	  ]
//	  引数: [
//	    spec: 検証対象のプロパティ仕様
//	  ]
//	  戻り値: [
//	    error: 仕様が不完全または曖昧な場合の失敗
//	  ]
//	}
func validatePropertySpecification(spec PropertySpecification) error {
	if spec.Version <= 0 {
		return fmt.Errorf("property specification version must be positive")
	}
	if len(spec.Axes) == 0 {
		return fmt.Errorf("property specification must define at least one axis")
	}
	axisNames := make(map[string]bool, len(spec.Axes))
	for _, axis := range spec.Axes {
		if axis.Name == "" || len(axis.States) == 0 || len(axis.AppliesTo) == 0 {
			return fmt.Errorf("property axis requires name, states, and targets")
		}
		axisName, err := CanonicalizeIdentifier(axis.Name)
		if err != nil {
			return fmt.Errorf("invalid property axis name: %w", err)
		}
		if axisNames[axisName] {
			return fmt.Errorf("duplicate property axis name %q", axis.Name)
		}
		axisNames[axisName] = true
		for _, state := range axis.States {
			if _, err := CanonicalizeIdentifier(state); err != nil {
				return fmt.Errorf("invalid property state in axis %q: %w", axis.Name, err)
			}
		}
		if axis.SourceDefault != "" && !propertyInAxis(axis, PreprocessedProperty(axis.SourceDefault)) {
			return fmt.Errorf("source default %q is not a state in property axis %q", axis.SourceDefault, axis.Name)
		}
		for _, target := range axis.AppliesTo {
			if _, err := CanonicalizeIdentifier(target); err != nil {
				return fmt.Errorf("invalid target in property axis %q: %w", axis.Name, err)
			}
		}
	}
	if err := validatePropertySourceDefaults(spec); err != nil {
		return err
	}
	return validateUniquePropertyStates(spec)
}

func validatePropertySourceDefaults(spec PropertySpecification) error {
	for _, outputAxis := range spec.Axes {
		if outputAxis.SourceDefaultWhen == nil {
			continue
		}
		if outputAxis.SourceDefault != "" {
			return fmt.Errorf("property axis %q cannot define both source default forms", outputAxis.Name)
		}
		if err := validatePropertySourceDefaultCondition(spec, outputAxis); err != nil {
			return err
		}
	}
	return nil
}

func validatePropertySourceDefaultCondition(spec PropertySpecification, outputAxis PropertyAxisSpec) error {
	condition := outputAxis.SourceDefaultWhen
	conditionAxis, found := propertyAxisByName(spec, condition.Axis)
	if !found {
		return fmt.Errorf("unknown condition axis %q for property axis %q", condition.Axis, outputAxis.Name)
	}
	if propertyNamesEqual(conditionAxis.Name, outputAxis.Name) {
		return fmt.Errorf("property axis %q cannot depend on itself", outputAxis.Name)
	}
	if !propertyAxisCoversTargets(conditionAxis, outputAxis) {
		return fmt.Errorf("condition axis %q does not apply to every target of property axis %q", conditionAxis.Name, outputAxis.Name)
	}
	return validatePropertySourceDefaultStates(conditionAxis, outputAxis, condition.States)
}

func propertyAxisByName(spec PropertySpecification, name string) (PropertyAxisSpec, bool) {
	for _, axis := range spec.Axes {
		if propertyNamesEqual(axis.Name, name) {
			return axis, true
		}
	}
	return PropertyAxisSpec{}, false
}

func propertyAxisCoversTargets(conditionAxis PropertyAxisSpec, outputAxis PropertyAxisSpec) bool {
	for _, target := range outputAxis.AppliesTo {
		if !propertyAxisApplies(conditionAxis, target) {
			return false
		}
	}
	return len(outputAxis.AppliesTo) > 0
}

func validatePropertySourceDefaultStates(conditionAxis PropertyAxisSpec, outputAxis PropertyAxisSpec, defaults map[string]string) error {
	if len(defaults) == 0 {
		return fmt.Errorf("conditional source default for property axis %q must define at least one state", outputAxis.Name)
	}
	conditionStates := make([]string, 0, len(defaults))
	for conditionState := range defaults {
		conditionStates = append(conditionStates, conditionState)
	}
	sort.Strings(conditionStates)
	seen := make(map[string]bool, len(conditionStates))
	for _, conditionState := range conditionStates {
		canonical, err := CanonicalizeIdentifier(conditionState)
		if err != nil || !propertyInAxis(conditionAxis, PreprocessedProperty(conditionState)) {
			return fmt.Errorf("unknown condition state %q in property axis %q", conditionState, conditionAxis.Name)
		}
		if seen[canonical] {
			return fmt.Errorf("duplicate condition state %q in property axis %q", conditionState, conditionAxis.Name)
		}
		seen[canonical] = true
		defaultState := defaults[conditionState]
		if !propertyInAxis(outputAxis, PreprocessedProperty(defaultState)) {
			return fmt.Errorf("unknown default state %q in property axis %q", defaultState, outputAxis.Name)
		}
	}
	return nil
}
