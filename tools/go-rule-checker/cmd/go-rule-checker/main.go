package main

import (
	"fmt"
	"os"

	gorulechecker "github.com/tomiya7688/Bitlang/tools/go-rule-checker"
)

//	{
//	  責務: [
//	    main: Go 実装ルールチェッカーを起動する
//	  ]
//	  処理: [
//	    1: コマンドラインの検査対象を取得する
//	    2: チェッカーの結果をプロセス終了状態へ対応付ける
//	  ]
//	  引数: []
//	  戻り値: []
//	}
func main() {
	code, err := gorulechecker.Run(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(code)
}
