package bitlang

//	{
//	  責務: [
//	    SourceMemberKind: ソースメンバーが表す宣言形式を区別する
//	  ]
//	  項目: [
//	    SourceFieldMember: フィールド宣言
//	    SourcePropertyMember: プロパティ宣言
//	    SourceConstructorMember: コンストラクター宣言
//	    SourceMethodMember: メソッド宣言
//	  ]
//	}
type SourceMemberKind string

const (
	// {
	//   責務: [
	//     SourceFieldMember: フィールド宣言を識別する
	//   ]
	// }
	SourceFieldMember SourceMemberKind = "field"
	// {
	//   責務: [
	//     SourcePropertyMember: プロパティ宣言を識別する
	//   ]
	// }
	SourcePropertyMember SourceMemberKind = "property"
	// {
	//   責務: [
	//     SourceConstructorMember: コンストラクター宣言を識別する
	//   ]
	// }
	SourceConstructorMember SourceMemberKind = "constructor"
	// {
	//   責務: [
	//     SourceMethodMember: メソッド宣言を識別する
	//   ]
	// }
	SourceMethodMember SourceMemberKind = "method"
)

//	{
//	  責務: [
//	    SourceClass: 意味変換前のクラス宣言とメンバー順を保持する
//	  ]
//	  フィールド: [
//	    Name: クラス名の元表記と比較表記
//	    TypeParameters: 型パラメーターの元トークン
//	    Modifiers: 修飾子の文字列表記
//	    ModifierTokens: 修飾子の元表記と位置
//	    Members: 元の順序を保ったメンバー列
//	    Line: クラス名の行番号
//	    Column: クラス名の列番号
//	  ]
//	  補足: [
//	    修飾子展開や文本体解析より前の情報を保持する
//	  ]
//	}
type SourceClass struct {
	Name           CanonicalName
	TypeParameters []PreprocessedToken
	Modifiers      []string
	ModifierTokens []PreprocessedToken
	Members        []SourceMember
	Line           int
	Column         int
}

//	{
//	  責務: [
//	    SourceMember: フィールド・プロパティ・コンストラクター・メソッドを保持する
//	  ]
//	  フィールド: [
//	    Kind: メンバーの宣言形式
//	    Name: 名前の元表記と比較表記
//	    Type: 型名の元表記と比較表記
//	    TypeArguments: 型引数の元トークン
//	    MethodTypeParameters: メソッド型パラメーターの元トークン
//	    Modifiers: 修飾子の文字列表記
//	    ModifierTokens: 修飾子の元表記と位置
//	    Parameters: パラメーター宣言
//	    Accessors: プロパティアクセサー宣言
//	    Initializer: 初期化式の元トークン
//	    Body: メソッド本体の未解析トークン
//	    Line: メンバー名の行番号
//	    Column: メンバー名の列番号
//	  ]
//	}
type SourceMember struct {
	Kind                 SourceMemberKind
	Name                 CanonicalName
	Type                 CanonicalName
	TypeArguments        []PreprocessedToken
	MethodTypeParameters []PreprocessedToken
	Modifiers            []string
	ModifierTokens       []PreprocessedToken
	Parameters           []SourceParameter
	Accessors            []SourceAccessor
	Initializer          []PreprocessedToken
	Body                 []PreprocessedToken
	Line                 int
	Column               int
}

//	{
//	  責務: [
//	    SourceParameter: パラメーター宣言と元の修飾子情報を保持する
//	  ]
//	  フィールド: [
//	    Name: 名前の元表記と比較表記
//	    Type: 型名の元表記と比較表記
//	    TypeArguments: 型引数の元トークン
//	    Modifiers: 修飾子の文字列表記
//	    ModifierTokens: 修飾子の元表記と位置
//	    Line: パラメーター名の行番号
//	    Column: パラメーター名の列番号
//	  ]
//	}
type SourceParameter struct {
	Name           CanonicalName
	Type           CanonicalName
	TypeArguments  []PreprocessedToken
	Modifiers      []string
	ModifierTokens []PreprocessedToken
	Line           int
	Column         int
}

//	{
//	  責務: [
//	    SourceAccessor: プロパティアクセサーの種類・表記・修飾子を保持する
//	  ]
//	  フィールド: [
//	    Kind: アクセサー種別
//	    Spelling: ソース上の元表記
//	    Modifiers: 修飾子の文字列表記
//	    ModifierTokens: 修飾子の元表記と位置
//	    Line: アクセサーの行番号
//	    Column: アクセサーの列番号
//	  ]
//	}
type SourceAccessor struct {
	Kind           string
	Spelling       string
	Modifiers      []string
	ModifierTokens []PreprocessedToken
	Line           int
	Column         int
}
