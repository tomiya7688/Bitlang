package main

import (
	"flag"
	"os"

	docpairchecker "github.com/tomiya7688/Bitlang/tools/doc-pair-checker"
)

//	{
//	  責務: [
//	    main: 日英ドキュメント対チェッカーを起動する
//	  ]
//	  処理: [
//	    1: コマンド引数から基準コミットと検査対象を取得する
//	    2: 検査結果をプロセス終了状態へ対応付ける
//	  ]
//	  引数: []
//	  戻り値: []
//	}
func main() {
	baseRef := flag.String("base", "", "git base ref used to require paired document updates")
	flag.Parse()
	os.Exit(docpairchecker.Run(".", *baseRef, os.Stdout, os.Stderr))
}
