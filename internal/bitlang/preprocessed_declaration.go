package bitlang

//	{
//	  責務: [
//	    PreprocessedDeclaration: 正規化済み宣言の型・名前・位置・最終プロパティを保持する
//	  ]
//	  フィールド: [
//	    Kind: 仕様で定義された宣言種別
//	    Name: 宣言名の元表記と比較表記
//	    Type: 型名の元表記と比較表記
//	    Properties: 適用軸ごとに明示された最終状態
//	    Line: 宣言開始行
//	    Column: 宣言開始列
//	  ]
//	  補足: [
//	    適用される全プロパティ軸を保持し、後段で既定値を再構成させない
//	  ]
//	}
type PreprocessedDeclaration struct {
	Kind       string
	Name       CanonicalName
	Type       CanonicalName
	Properties []PreprocessedProperty
	Line       int
	Column     int
}
