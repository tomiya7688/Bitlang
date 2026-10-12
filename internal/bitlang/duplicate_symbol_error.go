package bitlang

import "fmt"

//	{
//	  責務: [
//	    DuplicateSymbolError: 同じ識別子へ解決された二つの表記を診断情報として保持する
//	  ]
//	  フィールド: [
//	    Existing: 先に定義された表記
//	    Incoming: 重複として拒否された表記
//	  ]
//	}
type DuplicateSymbolError struct {
	Existing string
	Incoming string
}

//	{
//	  責務: [
//	    Error: 重複シンボルの診断文を生成する
//	  ]
//	  処理: [
//	    1: 既存表記と入力表記を区別してメッセージへ含める
//	  ]
//	  引数: [
//	    e: 重複シンボルの両表記
//	  ]
//	  戻り値: [
//	    message: エラー表示用の診断文
//	  ]
//	}
func (e DuplicateSymbolError) Error() string {
	return fmt.Sprintf("duplicate symbol: %q conflicts with existing %q", e.Incoming, e.Existing)
}
