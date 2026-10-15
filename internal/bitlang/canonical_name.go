package bitlang

import (
	"fmt"
	"strings"
)

//	{
//	  責務: [
//	    CanonicalizeIdentifier: 識別子の比較用表記を作る
//	  ]
//	  処理: [
//	    1: 空の識別子を拒否する
//	    2: 比較用に大文字小文字を正規化する
//	  ]
//	  引数: [
//	    name: 識別子の元表記。文字列・文字リテラルには適用しない
//	  ]
//	  戻り値: [
//	    canonical: 比較用表記
//	    error: 空識別子などの失敗
//	  ]
//	}
func CanonicalizeIdentifier(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("identifier must not be empty")
	}
	// TODO: Replace Go lowercase mapping when Bitlang Unicode identifier semantics are fixed.
	return strings.ToLower(name), nil
}

//	{
//	  責務: [
//	    CanonicalName: 識別子の表示表記と比較表記を保持する
//	  ]
//	  フィールド: [
//	    Spelling: 診断やソース再表示に使う元表記
//	    Canonical: 識別子照合に使う比較表記
//	  ]
//	}
type CanonicalName struct {
	Spelling  string
	Canonical string
}

//	{
//	  責務: [
//	    NewCanonicalName: 識別子の表示表記と比較表記を組み立てる
//	  ]
//	  処理: [
//	    1: 元表記を正規化して比較表記を得る
//	    2: 両方の表記を保持する値を返す
//	  ]
//	  引数: [
//	    spelling: 識別子の元表記
//	  ]
//	  戻り値: [
//	    name: 元表記と比較表記を持つ値
//	    error: 識別子が無効な場合の失敗
//	  ]
//	}
func NewCanonicalName(spelling string) (CanonicalName, error) {
	canonical, err := CanonicalizeIdentifier(spelling)
	if err != nil {
		return CanonicalName{}, err
	}
	return CanonicalName{Spelling: spelling, Canonical: canonical}, nil
}
