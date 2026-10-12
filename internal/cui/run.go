package cui

import (
	"fmt"
	"io"

	core "github.com/tomiya7688/Bitlang/internal/bitlang"
)

//	{
//	  責務: [
//	    Run: CUI の一回の呼び出しを処理し終了状態を返す
//	  ]
//	  処理: [
//	    1: 引数に応じてコンパイラー操作を選択する
//	    2: 結果または診断を指定された出力先へ渡す
//	    3: 成否をプロセス終了状態へ対応付ける
//	  ]
//	  引数: [
//	    args: CUI に渡されたコマンドライン引数
//	    out: 通常出力先
//	    errOut: エラー出力先
//	  ]
//	  戻り値: [
//	    exitCode: 呼び出し元へ返すプロセス終了状態
//	  ]
//	}
func Run(args []string, out io.Writer, errOut io.Writer) int {
	if len(args) < 1 {
		writeUsage(errOut)
		return 2
	}

	switch args[0] {
	case "canonicalize":
		if len(args) < 2 {
			writeUsage(errOut)
			return 2
		}
		for _, identifier := range args[1:] {
			canonical, err := core.CanonicalizeIdentifier(identifier)
			if err != nil {
				fmt.Fprintln(errOut, err)
				return 1
			}
			fmt.Fprintf(out, "%s\t%s\n", identifier, canonical)
		}
		return 0
	default:
		writeUsage(errOut)
		return 2
	}
}

func writeUsage(out io.Writer) {
	fmt.Fprintln(out, "usage:")
	fmt.Fprintln(out, "  bitlang canonicalize <identifier> [identifier ...]")
}
