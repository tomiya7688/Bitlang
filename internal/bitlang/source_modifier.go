package bitlang

import "fmt"

//	{
//	  責務: [
//	    SourceModifierMapping: アダプターの修飾子から正規プロパティへの対応を定義する
//	  ]
//	  フィールド: [
//	    Modifier: ソース言語側の修飾子名
//	    Properties: 展開先のプロパティ状態
//	  ]
//	  補足: [
//	    修飾子の意味を共通プロパティ機構へ固定しない
//	  ]
//	}
type SourceModifierMapping struct {
	Modifier   string                 `json:"modifier"`
	Properties []PreprocessedProperty `json:"properties"`
}

//	{
//	  責務: [
//	    ExpandSourceModifiers: ソース修飾子を仕様上の正規プロパティへ変換する
//	  ]
//	  処理: [
//	    1: 修飾子マッピングの重複を拒否する
//	    2: 入力順に修飾子をプロパティ列へ展開する
//	    3: 対象へ適用可能な正規状態へ解決する
//	    4: 重複状態を除き、排他軸の競合を拒否する
//	  ]
//	  引数: [
//	    spec: 正規プロパティ仕様
//	    target: プロパティ適用対象
//	    modifiers: ソース宣言に現れた修飾子
//	    mappings: アダプターが定義する展開規則
//	  ]
//	  戻り値: [
//	    properties: 重複を除いた正規プロパティ列
//	    error: 未定義・重複・競合がある場合の失敗
//	  ]
//	}
func ExpandSourceModifiers(spec PropertySpecification, target string, modifiers []string, mappings []SourceModifierMapping) ([]PreprocessedProperty, error) {
	byModifier := make(map[string][]PreprocessedProperty, len(mappings))
	for _, mapping := range mappings {
		name, err := CanonicalizeIdentifier(mapping.Modifier)
		if err != nil {
			return nil, fmt.Errorf("invalid modifier mapping: %w", err)
		}
		if _, exists := byModifier[name]; exists {
			return nil, fmt.Errorf("duplicate modifier mapping %q", mapping.Modifier)
		}
		byModifier[name] = mapping.Properties
	}

	expanded := make([]PreprocessedProperty, 0, len(modifiers))
	for _, modifier := range modifiers {
		name, err := CanonicalizeIdentifier(modifier)
		if err != nil {
			return nil, err
		}
		properties, found := byModifier[name]
		if !found {
			return nil, fmt.Errorf("no property mapping for source modifier %q", modifier)
		}
		expanded = append(expanded, properties...)
	}

	resolved, err := ResolveProperties(spec, target, expanded)
	if err != nil {
		return nil, err
	}
	result := make([]PreprocessedProperty, 0, len(resolved))
	seen := make(map[PreprocessedProperty]bool, len(resolved))
	for _, property := range resolved {
		if !seen[property] {
			seen[property] = true
			result = append(result, property)
		}
	}
	for _, axis := range spec.Axes {
		if axis.Exclusive && propertyAxisApplies(axis, target) && countPropertyAxisStates(axis, result) > 1 {
			return nil, fmt.Errorf("conflicting property states for axis %q", axis.Name)
		}
	}
	return result, nil
}
