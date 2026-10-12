package bitlang

import "fmt"

//	{
//	  責務: [
//	    ValidateSpecificationConsistency: 複数の機械可読仕様をまたぐ参照整合性を検証する
//	  ]
//	  処理: [
//	    1: プロパティ仕様そのものを検証する
//	    2: 各宣言種別の対象名が定義済みか確認する
//	  ]
//	  引数: [
//	    properties: プロパティ軸と対象の仕様
//	    declarations: 宣言種別と対象参照の仕様
//	  ]
//	  戻り値: [
//	    error: 個別仕様または相互参照が不整合な場合の失敗
//	  ]
//	}
func ValidateSpecificationConsistency(properties PropertySpecification, declarations DeclarationSpecification) error {
	if err := validatePropertySpecification(properties); err != nil {
		return err
	}
	for _, kind := range declarations.Kinds {
		if !propertyTargetKnown(properties, kind.PropertyTarget) {
			return fmt.Errorf("declaration kind %q references unknown property target %q", kind.Name, kind.PropertyTarget)
		}
	}
	return nil
}

func validateUniquePropertyStates(spec PropertySpecification) error {
	owners := map[string]string{}
	for _, axis := range spec.Axes {
		for _, state := range axis.States {
			if state == "" {
				return fmt.Errorf("property axis %q contains an empty state", axis.Name)
			}
			canonical, err := CanonicalizeIdentifier(state)
			if err != nil {
				return err
			}
			if owner, exists := owners[canonical]; exists {
				return fmt.Errorf("property state %q collides case-insensitively between axes %q and %q", state, owner, axis.Name)
			}
			owners[canonical] = axis.Name
		}
	}
	return nil
}
