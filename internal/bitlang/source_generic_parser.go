package bitlang

type genericListState uint8

const (
	genericExpectName genericListState = iota
	genericAfterName
	genericAfterType
)

//	{
//	  責務: [
//	    parseGenericTokenSequence: 型引数または型パラメーター列の形を検証する
//	  ]
//	  処理: [
//	    1: 反復スタックで入れ子の型引数を走査する
//	    2: 型の意味解決を行わず入力トークンを保持する
//	  ]
//	  引数: [
//	    p: トークン位置を持つクラス解析器
//	    parameters: 型パラメーター列として解析するか
//	  ]
//	  戻り値: [
//	    tokens: 入れ子と元表記を保つ型列トークン
//	    error: 区切りまたは要素形式が不正な場合の失敗
//	  ]
//	}
func (p *sourceClassParser) parseGenericTokenSequence(parameters bool) ([]PreprocessedToken, error) {
	start := p.index
	if _, err := p.expectSymbol("<"); err != nil {
		return nil, err
	}
	stack := []genericListState{genericExpectName}
	for len(stack) > 0 {
		last := len(stack) - 1
		switch stack[last] {
		case genericExpectName:
			token, err := p.expectIdentifier("generic parameter or type argument")
			if err != nil {
				return nil, err
			}
			if _, err := sourceName(token); err != nil {
				return nil, err
			}
			stack[last] = genericAfterName
		case genericAfterName:
			if !parameters && p.isSymbol(".") {
				p.advance()
				stack[last] = genericExpectName
				continue
			}
			stack[last] = genericAfterType
			if !parameters && p.isSymbol("<") {
				p.advance()
				stack = append(stack, genericExpectName)
			}
		case genericAfterType:
			if p.isSymbol(">") {
				p.advance()
				stack = stack[:last]
				continue
			}
			if _, err := p.expectSymbol(","); err != nil {
				return nil, err
			}
			stack[last] = genericExpectName
		}
	}
	return append([]PreprocessedToken(nil), p.tokens[start:p.index]...), nil
}
