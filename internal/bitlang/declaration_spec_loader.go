package bitlang

import "fmt"

//	{
//	  責務: [
//	    LoadDeclarationSpecification: 宣言文法仕様を読み込み、使用前に検証する
//	  ]
//	  処理: [
//	    1: 厳密な JSON 値としてデコードする
//	    2: 仕様の必須項目と種別定義を検証する
//	  ]
//	  引数: [
//	    data: 宣言仕様の JSON バイト列
//	  ]
//	  戻り値: [
//	    specification: 検証済み宣言仕様
//	    error: デコードまたは仕様検証の失敗
//	  ]
//	}
func LoadDeclarationSpecification(data []byte) (DeclarationSpecification, error) {
	if err := validateDeclarationSpecificationJSONFields(data); err != nil {
		return DeclarationSpecification{}, fmt.Errorf("decode declaration specification: %w", err)
	}
	var spec DeclarationSpecification
	if err := decodeSingleJSONValue(data, &spec); err != nil {
		return DeclarationSpecification{}, fmt.Errorf("decode declaration specification: %w", err)
	}
	if spec.Version <= 0 || len(spec.Kinds) == 0 {
		return DeclarationSpecification{}, fmt.Errorf("declaration specification requires version and kinds")
	}
	seen := map[string]bool{}
	for _, kind := range spec.Kinds {
		if kind.Name == "" || kind.PropertyTarget == "" || kind.Terminator == "" {
			return DeclarationSpecification{}, fmt.Errorf("declaration kind requires name, property target, and terminator")
		}
		if seen[kind.Name] {
			return DeclarationSpecification{}, fmt.Errorf("duplicate declaration kind %q", kind.Name)
		}
		seen[kind.Name] = true
		if !validDeclarationLayout(kind.Layout) {
			return DeclarationSpecification{}, fmt.Errorf("unsupported declaration layout for %q", kind.Name)
		}
		if err := validateDeclarationContexts(kind); err != nil {
			return DeclarationSpecification{}, err
		}
	}
	return spec, nil
}

//	{
//	  責務: [
//	    DeclarationKind: 仕様内の宣言種別を名前で解決する
//	  ]
//	  処理: [
//	    1: 種別名が一致する定義を返す
//	    2: 未定義名の場合は失敗を返す
//	  ]
//	  引数: [
//	    name: 仕様から解決する宣言種別名
//	  ]
//	  戻り値: [
//	    kind: 対応する宣言種別定義
//	    error: 未定義または重複時の失敗
//	  ]
//	}
func (s DeclarationSpecification) DeclarationKind(name string) (DeclarationKindSpec, error) {
	for _, kind := range s.Kinds {
		if kind.Name == name {
			return kind, nil
		}
	}
	return DeclarationKindSpec{}, fmt.Errorf("unknown declaration kind %q", name)
}
