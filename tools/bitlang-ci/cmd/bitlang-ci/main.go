package main

import (
	"os"

	bitlangci "github.com/tomiya7688/Bitlang/tools/bitlang-ci"
)

//	{
//	  責務: [
//	    main: 共有 Bitlang CI ゲートを起動して終了状態をプロセスへ返す
//	  ]
//	  処理: [
//	    1: 現在のリポジトリと標準入出力をゲートへ渡す
//	    2: ゲート結果をプロセス終了状態へ反映する
//	  ]
//	  引数: []
//	  戻り値: []
//	}
func main() {
	os.Exit(bitlangci.Run(".", os.Stdout, os.Stderr))
}
