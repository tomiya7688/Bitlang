package bitlang

//	{
//	  責務: [
//	    LexicalSource: 一回の字句解析結果を意味正規化前の成果物として保持する
//	  ]
//	  フィールド: [
//	    Path: 診断用の元ソースパス
//	    OriginalText: 空白・コメントを含む変更前のソース全文
//	    Tokens: 字句解析結果とソース表記
//	  ]
//	  補足: [
//	    正規化済み Preprocessed 境界への到達を意味しない
//	  ]
//	}
type LexicalSource struct {
	Path         string
	OriginalText string
	Tokens       []PreprocessedToken
}
