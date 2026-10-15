package bitlang

import "fmt"

//	{
//	  責務: [
//	    ValidateProperties: 対象宣言のプロパティ状態を外部仕様と照合する
//	  ]
//	  処理: [
//	    1: 仕様全体と対象名を先に検証する
//	    2: 適用軸ごとの必須状態と排他制約を確認する
//	    3: 未知・対象外状態を拒否する
//	  ]
//	  引数: [
//	    spec: 検証対象のプロパティ仕様
//	    target: 宣言種別などのプロパティ対象
//	    properties: 宣言に明示された最終状態
//	  ]
//	  戻り値: [
//	    error: 仕様・対象・状態のいずれかが不正な場合の失敗
//	  ]
//	  補足: [
//	    状態名と対象名は大文字小文字を区別せず比較し、入力表記は保持する
//	  ]
//	}
func ValidateProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) error {
	if err := validatePropertySpecification(spec); err != nil {
		return err
	}
	if !propertyTargetKnown(spec, target) {
		return fmt.Errorf("unknown property target %q", target)
	}
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) {
			continue
		}
		matches := countPropertyAxisStates(axis, properties)
		if axis.Required && matches == 0 {
			return fmt.Errorf("missing required property axis %q for %s", axis.Name, target)
		}
		if axis.Exclusive && matches > 1 {
			return fmt.Errorf("conflicting property states for axis %q", axis.Name)
		}
	}
	return validateKnownProperties(spec, target, properties)
}

func validateKnownProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) error {
	for _, property := range properties {
		known := false
		for _, axis := range spec.Axes {
			if propertyAxisApplies(axis, target) && propertyInAxis(axis, property) {
				known = true
				break
			}
		}
		if !known {
			return fmt.Errorf("property %q is not applicable to %s", property, target)
		}
	}
	return nil
}
