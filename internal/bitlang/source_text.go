package bitlang

//	{
//	  責務: [
//	    SourceText: 字句解析前のソース単位と入力内容を保持する
//	  ]
//	  フィールド: [
//	    Path: 診断表示用のパス情報
//	    Text: 受け取ったままのソース文字列
//	  ]
//	  補足: [
//	    後段の位置情報とリテラル内容を保つため、字句解析前に Text を正規化しない
//	  ]
//	}
type SourceText struct {
	Path string
	Text string
}

//	{
//	  責務: [
//	    NewSourceText: パスとソース文字列を保持する値を作る
//	  ]
//	  引数: [
//	    path: 診断情報に使うソースパス
//	    text: 変更せずに保持する入力ソース
//	  ]
//	  戻り値: [
//	    source: 指定内容をそのまま持つ SourceText
//	  ]
//	}
func NewSourceText(path string, text string) SourceText {
	return SourceText{Path: path, Text: text}
}
