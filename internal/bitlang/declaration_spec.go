package bitlang

//	{
//	  責務: [
//	    DeclarationSpecification: 宣言種別に依存しない厳密な宣言形式仕様を保持する
//	  ]
//	  フィールド: [
//	    Version: 仕様形式のバージョン
//	    Kinds: 宣言種別ごとの構文定義
//	  ]
//	  補足: [
//	    宣言種別名は Go 実装へ固定せず、入力仕様から取得する
//	  ]
//	}
type DeclarationSpecification struct {
	Version int                   `json:"version"`
	Kinds   []DeclarationKindSpec `json:"kinds"`
}

//	{
//	  責務: [
//	    DeclarationKindSpec: 一つの宣言種別の解析規則を定義する
//	  ]
//	  フィールド: [
//	    Name: 仕様内で参照する種別名
//	    PropertyTarget: 適用するプロパティ対象名
//	    Terminator: 宣言列を区切る終端記号
//	    Layout: 構成要素の並び
//	    Contexts: この種別を許可する宣言コンテキスト
//	  ]
//	}
type DeclarationKindSpec struct {
	Name           string   `json:"name"`
	PropertyTarget string   `json:"property_target"`
	Terminator     string   `json:"terminator"`
	Layout         []string `json:"layout"`
	Contexts       []string `json:"contexts,omitempty"`
}
