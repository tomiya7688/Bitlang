package bitlang

//	{
//	  責務: [
//	    SourceParseResult: 字句ソースと各解析結果を結び付けて中間状態として保持する
//	  ]
//	  フィールド: [
//	    Source: 字句解析時のパス・元ソース・トークン
//	    Classes: 元情報を保ったクラス解析結果
//	    Declarations: プリプロセス済み宣言解析結果
//	    Scope: 宣言解析に使った意味スコープ
//	  ]
//	  補足: [
//	    正規化済み Preprocessed 成果物とは区別する
//	  ]
//	}
type SourceParseResult struct {
	Source       LexicalSource
	Classes      []SourceClass
	Declarations []PreprocessedDeclaration
	Scope        *PreprocessedScope
}

//	{
//	  責務: [
//	    NewSourceParseResult: 字句ソースに結び付いた空の解析結果を作る
//	  ]
//	  引数: [
//	    source: 字句解析済みの中間ソース
//	  ]
//	  戻り値: [
//	    result: ソースだけを関連付け、解析結果を未設定にした値
//	  ]
//	  補足: [
//	    構文解析と正規プリプロセスは別の段階として維持する
//	  ]
//	}
func NewSourceParseResult(source LexicalSource) SourceParseResult {
	return SourceParseResult{Source: source}
}
