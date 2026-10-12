package gorulechecker

import (
	"fmt"
	"io"
	"sort"
)

//	{
//	  責務: [
//	    Run: 指定された Go ソース範囲にルール検査を実行する
//	  ]
//	  処理: [
//	    1: 対象ファイルを収集し、無視設定を読み込む
//	    2: 有効な静的ルールを各ファイルへ適用する
//	    3: findings を出力して検査結果を返す
//	  ]
//	  引数: [
//	    paths: 検査するファイルまたはディレクトリ
//	    out: 検出結果の出力先
//	  ]
//	  戻り値: [
//	    exitCode: 検出結果を表す終了状態
//	    error: 検査処理自体を続行できない場合の失敗
//	  ]
//	}
func Run(paths []string, out io.Writer) (int, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	config, err := loadIgnoreConfig()
	if err != nil {
		return 0, err
	}
	files, err := collectGoFiles(paths)
	if err != nil {
		return 0, err
	}
	var findings []finding
	for _, path := range files {
		fileFindings, err := inspectFile(path)
		if err != nil {
			return 0, err
		}
		findings = append(findings, fileFindings...)
	}
	responsibilityFindings, err := checkResponsibilities(files, isFullRepositoryCheck(paths))
	if err != nil {
		return 0, err
	}
	findings = append(findings, responsibilityFindings...)

	filtered := findings[:0]
	for _, item := range findings {
		if !config.ignores(item) {
			filtered = append(filtered, item)
		}
	}
	findings = filtered
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].path != findings[j].path {
			return findings[i].path < findings[j].path
		}
		return findings[i].line < findings[j].line
	})
	for _, item := range findings {
		fmt.Fprintln(out, item.String())
	}
	if len(findings) == 0 {
		fmt.Fprintln(out, "OK go-rules")
		return 0, nil
	}
	fmt.Fprintf(out, "NG go-rules: %d\n", len(findings))
	return 1, nil
}
