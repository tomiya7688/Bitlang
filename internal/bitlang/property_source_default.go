package bitlang

import "fmt"

//	{
//	  責務: [
//	    ApplySourcePropertyDefaults: 明示値と上位規則を保ち未解決軸へ組み込み既定値を補う
//	  ]
//	  処理: [
//	    1: property仕様と対象を検証する
//	    2: 指定状態を仕様上の表記へ解決する
//	    3: 条件なし既定値を未指定の適用軸へ追加する
//	    4: 条件軸の既定値を解決できるまで条件付き既定値を追加する
//	  ]
//	  引数: [
//	    spec: 軸・状態・適用対象・組み込み既定値を定める仕様
//	    target: propertyを適用する宣言対象
//	    properties: 明示値および適用済み上位規則の状態列
//	  ]
//	  戻り値: [
//	    resolved: 入力状態と不足軸の組み込み既定値を正規表記で保持する列
//	    error: 仕様・対象・明示状態が不正な場合の失敗
//	  ]
//	  補足: [
//	    必須軸が全て埋まったことは保証しない。最終境界ではValidatePropertiesを実行する
//	    条件軸が未解決の場合は推測せず、条件付き既定値も未解決のまま残す
//	  ]
//	}
func ApplySourcePropertyDefaults(spec PropertySpecification, target string, properties []PreprocessedProperty) ([]PreprocessedProperty, error) {
	if err := validatePropertySpecification(spec); err != nil {
		return nil, err
	}
	if !propertyTargetKnown(spec, target) {
		return nil, fmt.Errorf("unknown property target %q", target)
	}
	resolved, err := ResolveProperties(spec, target, properties)
	if err != nil {
		return nil, err
	}
	resolved, err = applyUnconditionalSourceDefaults(spec, target, resolved)
	if err != nil {
		return nil, err
	}
	return applyConditionalSourceDefaults(spec, target, resolved)
}

func applyUnconditionalSourceDefaults(spec PropertySpecification, target string, properties []PreprocessedProperty) ([]PreprocessedProperty, error) {
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) || axis.SourceDefault == "" {
			continue
		}
		states := countPropertyAxisStates(axis, properties)
		if axis.Exclusive && states > 1 {
			return nil, fmt.Errorf("conflicting property states for axis %q", axis.Name)
		}
		if states == 0 {
			properties = append(properties, PreprocessedProperty(axis.SourceDefault))
		}
	}
	return properties, nil
}

func applyConditionalSourceDefaults(spec PropertySpecification, target string, properties []PreprocessedProperty) ([]PreprocessedProperty, error) {
	axisByName := make(map[string]PropertyAxisSpec, len(spec.Axes))
	for _, axis := range spec.Axes {
		canonicalName, err := CanonicalizeIdentifier(axis.Name)
		if err != nil {
			return nil, err
		}
		axisByName[canonicalName] = axis
	}

	dependents := make(map[string][]PropertyAxisSpec)
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) || axis.SourceDefaultWhen == nil {
			continue
		}
		states := countPropertyAxisStates(axis, properties)
		if axis.Exclusive && states > 1 {
			return nil, fmt.Errorf("conflicting property states for axis %q", axis.Name)
		}
		if states > 0 {
			continue
		}
		conditionName, err := CanonicalizeIdentifier(axis.SourceDefaultWhen.Axis)
		if err != nil {
			return nil, err
		}
		conditionAxis, found := axisByName[conditionName]
		if !found {
			return nil, fmt.Errorf("unknown condition axis %q for property axis %q", axis.SourceDefaultWhen.Axis, axis.Name)
		}
		if !propertyAxisApplies(conditionAxis, target) {
			return nil, fmt.Errorf("condition axis %q does not apply to %s", conditionAxis.Name, target)
		}
		dependents[conditionName] = append(dependents[conditionName], axis)
	}

	queue := make([]string, 0, len(spec.Axes))
	queued := make(map[string]bool, len(spec.Axes))
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) || countPropertyAxisStates(axis, properties) == 0 {
			continue
		}
		canonicalName, err := CanonicalizeIdentifier(axis.Name)
		if err != nil {
			return nil, err
		}
		queue = append(queue, canonicalName)
		queued[canonicalName] = true
	}

	for next := 0; next < len(queue); next++ {
		conditionAxis := axisByName[queue[next]]
		for _, outputAxis := range dependents[queue[next]] {
			defaultState, found, err := resolveConditionalSourceDefault(conditionAxis, target, outputAxis, properties)
			if err != nil {
				return nil, err
			}
			if !found {
				continue
			}
			properties = append(properties, defaultState)
			outputName, err := CanonicalizeIdentifier(outputAxis.Name)
			if err != nil {
				return nil, err
			}
			if !queued[outputName] {
				queue = append(queue, outputName)
				queued[outputName] = true
			}
		}
	}
	return properties, nil
}

func resolveConditionalSourceDefault(conditionAxis PropertyAxisSpec, target string, outputAxis PropertyAxisSpec, properties []PreprocessedProperty) (PreprocessedProperty, bool, error) {
	if !propertyAxisApplies(conditionAxis, target) {
		return "", false, fmt.Errorf("condition axis %q does not apply to %s", conditionAxis.Name, target)
	}
	conditionCount := countPropertyAxisStates(conditionAxis, properties)
	if conditionAxis.Exclusive && conditionCount > 1 {
		return "", false, fmt.Errorf("conflicting condition states for property axis %q", conditionAxis.Name)
	}
	return matchingConditionalDefault(conditionAxis, outputAxis, outputAxis.SourceDefaultWhen.States, properties)
}

func matchingConditionalDefault(conditionAxis PropertyAxisSpec, outputAxis PropertyAxisSpec, defaults map[string]string, properties []PreprocessedProperty) (PreprocessedProperty, bool, error) {
	var selected PreprocessedProperty
	for _, conditionState := range conditionAxis.States {
		if !hasPropertyState(properties, conditionState) {
			continue
		}
		mappedState, found := mappedDefaultState(defaults, conditionState)
		if !found {
			continue
		}
		canonicalState, found := canonicalPropertyState(outputAxis, mappedState)
		if !found {
			return "", false, fmt.Errorf("unknown conditional default %q in property axis %q", mappedState, outputAxis.Name)
		}
		if selected != "" && !propertyNamesEqual(string(selected), string(canonicalState)) {
			return "", false, fmt.Errorf("conflicting conditional defaults for property axis %q", outputAxis.Name)
		}
		selected = canonicalState
	}
	return selected, selected != "", nil
}

func hasPropertyState(properties []PreprocessedProperty, state string) bool {
	for _, property := range properties {
		if propertyNamesEqual(string(property), state) {
			return true
		}
	}
	return false
}

func mappedDefaultState(defaults map[string]string, conditionState string) (string, bool) {
	for state, defaultValue := range defaults {
		if propertyNamesEqual(state, conditionState) {
			return defaultValue, true
		}
	}
	return "", false
}

func canonicalPropertyState(axis PropertyAxisSpec, state string) (PreprocessedProperty, bool) {
	for _, candidate := range axis.States {
		if propertyNamesEqual(candidate, state) {
			return PreprocessedProperty(candidate), true
		}
	}
	return "", false
}
