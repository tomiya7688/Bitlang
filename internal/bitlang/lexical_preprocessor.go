package bitlang

//	{
//	  責務: [
//	    LexicalPreprocess: ブートストラップ字句解析結果を明示的な中間値にする
//	  ]
//	  処理: [
//	    1: ソースをトークン列へ変換する
//	    2: 識別子とキーワードの比較表記を保持する
//	  ]
//	  引数: [
//	    source: 変更せずに扱う入力ソース
//	  ]
//	  戻り値: [
//	    lexicalSource: パス・元表記・トークンを保持する字句段階の成果物
//	    error: 字句解析に失敗した場合の内容
//	  ]
//	  補足: [
//	    構文糖衣展開、プロパティ補完、旧形式移行、収束までの反復は行わない
//	  ]
//	}
func LexicalPreprocess(source SourceText) (LexicalSource, error) {
	tokens, err := NewLexer(source).Lex()
	if err != nil {
		return LexicalSource{}, err
	}
	strict := make([]PreprocessedToken, 0, len(tokens))
	for _, token := range tokens {
		converted, err := newPreprocessedToken(token)
		if err != nil {
			return LexicalSource{}, err
		}
		strict = append(strict, converted)
	}
	return LexicalSource{
		Path: source.Path, OriginalText: source.Text, Tokens: strict,
	}, nil
}
