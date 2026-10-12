package bitlang

//	{
//	  責務: [
//	    Symbol: 名前付き意味値とその元表記を一組で保持する
//	  ]
//	  フィールド: [
//	    Name: 識別子の表示表記と比較表記
//	    Value: 名前に関連付けられた意味値
//	  ]
//	  補足: [
//	    他言語実装でも再現しやすいよう Go ジェネリクスに依存しない
//	  ]
//	}
type Symbol struct {
	Name  CanonicalName
	Value any
}
