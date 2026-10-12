package bitlang

//	{
//	  責務: [
//	    SpecificationSet: 個別・横断検証が済んだコンパイラー仕様をまとめる
//	  ]
//	  フィールド: [
//	    Properties: 宣言プロパティの軸・状態・対象仕様
//	    Declarations: 宣言種別と構文形式の仕様
//	  ]
//	}
type SpecificationSet struct {
	Properties   PropertySpecification
	Declarations DeclarationSpecification
}
