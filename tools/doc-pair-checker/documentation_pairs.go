package docpairchecker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type documentationPair struct {
	englishPath  string
	japanesePath string
}

var requiredDocumentationPairs = []documentationPair{
	{englishPath: "README.md", japanesePath: "README.ja.md"},
	{englishPath: "CONTRIBUTING.md", japanesePath: "CONTRIBUTING.ja.md"},
	{englishPath: "CURRENT_STATE.md", japanesePath: "CURRENT_STATE.ja.md"},
	{englishPath: "SECURITY.md", japanesePath: "SECURITY.ja.md"},
	{englishPath: ".github/SECURITY_SETTINGS.md", japanesePath: ".github/SECURITY_SETTINGS.ja.md"},
	{englishPath: "tools/bitlang-ci/README.md", japanesePath: "tools/bitlang-ci/README.ja.md"},
	{englishPath: "tools/doc-pair-checker/README.md", japanesePath: "tools/doc-pair-checker/README.ja.md"},
}

//	{
//	  責務: [
//	    Run: 登録済みの英日ドキュメント対の整合性を検査する
//	  ]
//	  処理: [
//	    1: 指定された変更範囲を取得する
//	    2: 変更された登録済みドキュメント対を検証する
//	    3: findings と成否を出力する
//	  ]
//	  引数: [
//	    root: 検査対象リポジトリのルート
//	    baseRef: 差分比較に使う基準コミット
//	    out: 通常出力先
//	    errOut: エラー出力先
//	  ]
//	  戻り値: [
//	    exitCode: 検査結果を表す終了状態
//	  ]
//	}
func Run(root string, baseRef string, out io.Writer, errOut io.Writer) int {
	changedFiles, err := loadChangedFiles(root, baseRef)
	if err != nil {
		fmt.Fprintf(errOut, "FAIL doc-pairs: %v\n", err)
		return 1
	}

	failures := 0
	for _, pair := range requiredDocumentationPairs {
		if !validateDocumentationPair(root, pair, changedFiles, errOut) {
			failures++
		}
	}

	if failures > 0 {
		fmt.Fprintf(errOut, "NG doc-pairs: %d pair(s) failed\n", failures)
		return 1
	}

	fmt.Fprintln(out, "OK doc-pairs")
	return 0
}

func validateDocumentationPair(root string, pair documentationPair, changedFiles map[string]struct{}, errOut io.Writer) bool {
	english, englishErr := readRegisteredDocumentation(root, pair.englishPath)
	japanese, japaneseErr := readRegisteredDocumentation(root, pair.japanesePath)
	if englishErr != nil || japaneseErr != nil {
		fmt.Fprintf(errOut, "FAIL doc-pairs: invalid pair %s <-> %s: English=%v Japanese=%v\n", pair.englishPath, pair.japanesePath, englishErr, japaneseErr)
		return false
	}

	valid := true
	if !bytes.Contains(english, []byte(path.Base(pair.japanesePath))) {
		fmt.Fprintf(errOut, "FAIL doc-pairs: %s does not link to %s\n", pair.englishPath, pair.japanesePath)
		valid = false
	}
	if !bytes.Contains(japanese, []byte(path.Base(pair.englishPath))) {
		fmt.Fprintf(errOut, "FAIL doc-pairs: %s does not link to %s\n", pair.japanesePath, pair.englishPath)
		valid = false
	}

	if changedFiles != nil {
		_, englishChanged := changedFiles[pair.englishPath]
		_, japaneseChanged := changedFiles[pair.japanesePath]
		if englishChanged != japaneseChanged {
			fmt.Fprintf(errOut, "FAIL doc-pairs: update both %s and %s in the same change\n", pair.englishPath, pair.japanesePath)
			valid = false
		}
	}

	return valid
}

func loadChangedFiles(root string, baseRef string) (map[string]struct{}, error) {
	if strings.TrimSpace(baseRef) == "" {
		return nil, nil
	}
	if err := validateBaseCommit(baseRef); err != nil {
		return nil, err
	}

	// #nosec G204 -- baseRef is validated as a full hexadecimal commit object ID and no shell is involved.
	command := exec.Command("git", "diff", "--name-only", baseRef+"...HEAD", "--")
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git diff from %s: %v: %s", baseRef, err, strings.TrimSpace(string(output)))
	}

	changedFiles := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		filePath := filepath.ToSlash(strings.TrimSpace(scanner.Text()))
		if filePath != "" {
			changedFiles[filePath] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read git diff output: %w", err)
	}
	return changedFiles, nil
}
