package bitlang

//	{
//	  責務: [
//	    propertyNamesEqual: プロパティ名を識別子と同じ規則で比較する
//	  ]
//	  処理: [
//	    1: 両方の名前を比較用に正規化する
//	    2: 正規表記の一致を返す
//	  ]
//	  引数: [
//	    left: 比較する一方の名前
//	    right: 比較するもう一方の名前
//	  ]
//	  戻り値: [
//	    equal: 正規表記が一致するか
//	  ]
//	  補足: [
//	    診断向けの元表記や仕様文字列は書き換えない
//	  ]
//	}
func propertyNamesEqual(left, right string) bool {
	leftName, err := CanonicalizeIdentifier(left)
	if err != nil {
		return false
	}
	rightName, err := CanonicalizeIdentifier(right)
	return err == nil && leftName == rightName
}

func propertyAxisApplies(axis PropertyAxisSpec, target string) bool {
	for _, candidate := range axis.AppliesTo {
		if propertyNamesEqual(candidate, target) {
			return true
		}
	}
	return false
}

func propertyTargetKnown(spec PropertySpecification, target string) bool {
	for _, axis := range spec.Axes {
		if propertyAxisApplies(axis, target) {
			return true
		}
	}
	return false
}

func countPropertyAxisStates(axis PropertyAxisSpec, properties []PreprocessedProperty) int {
	count := 0
	for _, property := range properties {
		for _, state := range axis.States {
			if propertyNamesEqual(string(property), state) {
				count++
			}
		}
	}
	return count
}

func propertyInAxis(axis PropertyAxisSpec, property PreprocessedProperty) bool {
	for _, state := range axis.States {
		if propertyNamesEqual(state, string(property)) {
			return true
		}
	}
	return false
}
