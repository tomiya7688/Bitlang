package bitlang

import "fmt"

//	{
//	  責務: [
//	    ResolveProperties: 入力プロパティ表記を仕様上の正規状態へ対応付ける
//	  ]
//	  処理: [
//	    1: 入力順に各プロパティを対象と仕様に照合する
//	    2: 仕様に定義された元表記で結果を返す
//	  ]
//	  引数: [
//	    spec: 検証済みのプロパティ仕様
//	    target: プロパティを適用する対象
//	    properties: ソース側のプロパティ状態列
//	  ]
//	  戻り値: [
//	    resolved: 仕様から取得した正規状態列
//	    error: 未定義または対象外の状態がある場合の失敗
//	  ]
//	}
func ResolveProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) ([]PreprocessedProperty, error) {
	resolved := make([]PreprocessedProperty, 0, len(properties))
	for _, property := range properties {
		state, err := resolveProperty(spec, target, property)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, state)
	}
	return resolved, nil
}

func resolveProperty(spec PropertySpecification, target string, property PreprocessedProperty) (PreprocessedProperty, error) {
	propertyName, err := CanonicalizeIdentifier(string(property))
	if err != nil {
		return "", err
	}
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) {
			continue
		}
		for _, state := range axis.States {
			stateName, err := CanonicalizeIdentifier(state)
			if err != nil {
				return "", err
			}
			if stateName == propertyName {
				return PreprocessedProperty(state), nil
			}
		}
	}
	return "", fmt.Errorf("property %q is not applicable to %s", property, target)
}
