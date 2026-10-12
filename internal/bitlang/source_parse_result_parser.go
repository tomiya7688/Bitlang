package bitlang

//	{
//	  責務: [
//	    ParseDeclarationsInSourceResult: 指定された宣言種別で解析し、既存の結果へ追加する
//	  ]
//	  処理: [
//	    1: 関連する字句トークン列を宣言列として解析する
//	    2: 他のソース・クラス・スコープ情報を保ち Declarations を設定する
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    kindName: 適用する宣言種別名
//	    source: 字句ソースと既存の解析情報
//	  ]
//	  戻り値: [
//	    result: 宣言列を追加した解析結果
//	    error: 宣言列の解析に失敗した場合の内容
//	  ]
//	}
func ParseDeclarationsInSourceResult(specs SpecificationSet, kindName string, source SourceParseResult) (SourceParseResult, error) {
	declarations, err := ParsePreprocessedDeclarations(specs, kindName, source.Source.Tokens)
	if err != nil {
		return SourceParseResult{}, err
	}
	source.Declarations = declarations
	return source, nil
}

//	{
//	  責務: [
//	    ParseMixedDeclarationsInSourceResult: 既存スコープを使って混在宣言列を解析する
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    source: 字句ソース・スコープ・既存解析情報
//	  ]
//	  戻り値: [
//	    result: 宣言列を追加した解析結果
//	    error: スコープまたは宣言解析の失敗
//	  ]
//	}
func ParseMixedDeclarationsInSourceResult(specs SpecificationSet, source SourceParseResult) (SourceParseResult, error) {
	return ParseMixedDeclarationsInSourceResultInScope(specs, source.Scope, source)
}

//	{
//	  責務: [
//	    ParseMixedDeclarationsInSourceResultInContext: 文字列文脈の入口から型付きスコープで混在宣言を解析する
//	  ]
//	  処理: [
//	    1: 空でない文脈を既存スコープの子として構成する
//	    2: 構成したスコープで宣言列を解析し、結果へ関連付ける
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    context: 現在の文脈名。空の場合は既存スコープを使う
//	    source: 字句ソースと既存の解析情報
//	  ]
//	  戻り値: [
//	    result: 文脈と宣言列を追加した解析結果
//	    error: スコープ生成または宣言解析の失敗
//	  ]
//	}
func ParseMixedDeclarationsInSourceResultInContext(specs SpecificationSet, context string, source SourceParseResult) (SourceParseResult, error) {
	if context == "" {
		return ParseMixedDeclarationsInSourceResultInScope(specs, nil, source)
	}
	scope, err := NewPreprocessedScope(context, source.Scope)
	if err != nil {
		return SourceParseResult{}, err
	}
	return ParseMixedDeclarationsInSourceResultInScope(specs, &scope, source)
}

//	{
//	  責務: [
//	    ParseMixedDeclarationsInSourceResultInScope: 明示スコープで混在宣言を解析し結果へ保存する
//	  ]
//	  処理: [
//	    1: スコープに適合する宣言列を字句トークンから解析する
//	    2: 指定スコープと宣言列を解析結果へ設定する
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    scope: 宣言種別を絞り込む明示スコープ
//	    source: 字句ソースと既存の解析情報
//	  ]
//	  戻り値: [
//	    result: スコープと宣言列を追加した解析結果
//	    error: 宣言解析の失敗
//	  ]
//	}
func ParseMixedDeclarationsInSourceResultInScope(specs SpecificationSet, scope *PreprocessedScope, source SourceParseResult) (SourceParseResult, error) {
	declarations, err := ParseMixedPreprocessedDeclarationsInScope(specs, scope, source.Source.Tokens)
	if err != nil {
		return SourceParseResult{}, err
	}
	source.Scope = scope
	source.Declarations = declarations
	return source, nil
}
