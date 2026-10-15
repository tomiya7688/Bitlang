package bitlang

//	{
//	  責務: [
//	    PreprocessedScope: 宣言の意味解釈に使う文脈と親スコープを保持する
//	  ]
//	  フィールド: [
//	    Context: 仕様で定義された文脈名
//	    Parent: 外側のスコープ。最上位では nil
//	  ]
//	  補足: [
//	    文脈名は Go 実装へ固定せず、仕様データから取得する
//	  ]
//	}
type PreprocessedScope struct {
	Context CanonicalName
	Parent  *PreprocessedScope
}

//	{
//	  責務: [
//	    NewPreprocessedScope: 文脈名と親スコープを持つスコープ値を作る
//	  ]
//	  処理: [
//	    1: 文脈名の元表記と比較表記を作る
//	    2: 親スコープと関連付ける
//	  ]
//	  引数: [
//	    context: 現在のスコープ文脈名
//	    parent: 外側のスコープ。最上位では nil
//	  ]
//	  戻り値: [
//	    scope: 作成した型付きスコープ
//	    error: 文脈名が無効な場合の失敗
//	  ]
//	}
func NewPreprocessedScope(context string, parent *PreprocessedScope) (PreprocessedScope, error) {
	name, err := NewCanonicalName(context)
	if err != nil {
		return PreprocessedScope{}, err
	}
	return PreprocessedScope{Context: name, Parent: parent}, nil
}
