package gorulechecker

import "fmt"

type finding struct {
	level   string
	rule    string
	path    string
	line    int
	message string
}

//	{
//	  責務: [
//	    String: 一件のルール違反を簡潔な端末表示へ整形する
//	  ]
//	  引数: [
//	    f: ルール名・場所・詳細を持つ検出結果
//	  ]
//	  戻り値: [
//	    text: 端末へ出力する一行の表現
//	  ]
//	}
func (f finding) String() string {
	if f.line > 0 {
		return fmt.Sprintf("%s %s:%d %s %s", f.level, f.path, f.line, f.rule, f.message)
	}
	return fmt.Sprintf("%s %s %s %s", f.level, f.path, f.rule, f.message)
}
