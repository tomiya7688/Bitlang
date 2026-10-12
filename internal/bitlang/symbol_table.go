package bitlang

//	{
//	  責務: [
//	    SymbolTable: 正規化した識別子をキーに意味値を保持する
//	  ]
//	  フィールド: [
//	    symbols: 比較表記から元表記付きシンボルへの対応
//	  ]
//	}
type SymbolTable struct {
	symbols map[string]Symbol
}

//	{
//	  責務: [
//	    NewSymbolTable: シンボル未登録の表を作る
//	  ]
//	  戻り値: [
//	    table: 識別子を登録・検索できる空の表
//	  ]
//	}
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]Symbol)}
}

//	{
//	  責務: [
//	    Define: 元表記を保持したシンボルを登録し、正規名の衝突を拒否する
//	  ]
//	  処理: [
//	    1: 名前の表示表記と比較表記を作る
//	    2: 既存キーとの衝突を検査する
//	    3: 未登録の場合に値と元表記を保存する
//	  ]
//	  引数: [
//	    t: 登録先のシンボル表
//	    name: シンボル名の元表記
//	    value: 名前に関連付ける意味値
//	  ]
//	  戻り値: [
//	    symbol: 登録したシンボル
//	    error: 無効名または正規名の重複
//	  ]
//	}
func (t *SymbolTable) Define(name string, value any) (Symbol, error) {
	parsed, err := NewCanonicalName(name)
	if err != nil {
		return Symbol{}, err
	}
	if previous, ok := t.symbols[parsed.Canonical]; ok {
		return Symbol{}, DuplicateSymbolError{Existing: previous.Name.Spelling, Incoming: parsed.Spelling}
	}
	symbol := Symbol{Name: parsed, Value: value}
	t.symbols[parsed.Canonical] = symbol
	return symbol, nil
}

//	{
//	  責務: [
//	    Find: 大文字小文字を区別せず識別子名からシンボルを検索する
//	  ]
//	  処理: [
//	    1: 検索名を比較表記へ正規化する
//	    2: 登録済みシンボルと存在有無を返す
//	  ]
//	  引数: [
//	    t: 検索対象のシンボル表
//	    name: 検索する識別子名
//	  ]
//	  戻り値: [
//	    symbol: 見つかった場合の元表記付きシンボル
//	    found: 対応するシンボルがあるか
//	  ]
//	}
func (t *SymbolTable) Find(name string) (Symbol, bool) {
	canonical, err := CanonicalizeIdentifier(name)
	if err != nil {
		return Symbol{}, false
	}
	symbol, ok := t.symbols[canonical]
	return symbol, ok
}

//	{
//	  責務: [
//	    Get: 識別子名から登録済みの意味値を取得する
//	  ]
//	  引数: [
//	    t: 検索対象のシンボル表
//	    name: 検索する識別子名
//	  ]
//	  戻り値: [
//	    value: 登録値。未登録時は nil
//	    found: 対応するシンボルがあるか
//	  ]
//	}
func (t *SymbolTable) Get(name string) (any, bool) {
	symbol, ok := t.Find(name)
	if !ok {
		return nil, false
	}
	return symbol.Value, true
}

//	{
//	  責務: [
//	    Len: シンボル表の登録件数を返す
//	  ]
//	  引数: [
//	    t: 件数を調べるシンボル表
//	  ]
//	  戻り値: [
//	    count: 登録されているシンボル数
//	  ]
//	}
func (t *SymbolTable) Len() int {
	return len(t.symbols)
}
