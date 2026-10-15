package bitlang

//	{
//	  責務: [
//	    PropertySpecification: プリプロセス済み成果物のプロパティ軸仕様を保持する
//	  ]
//	  フィールド: [
//	    Version: 仕様形式のバージョン
//	    Axes: 軸ごとの状態・適用対象・制約
//	  ]
//	  補足: [
//	    Bitlang 固有の状態名を Go コードへ重複定義しない
//	  ]
//	}
type PropertySpecification struct {
	Version int                `json:"version"`
	Axes    []PropertyAxisSpec `json:"axes"`
}

//	{
//	  責務: [
//	    PropertyAxisSpec: 一つのプロパティ軸の有効状態と適用規則を定義する
//	  ]
//	  フィールド: [
//	    Name: 軸名
//	    States: 軸で許可する最終状態
//	    Exclusive: 同時に複数状態を許可しないか
//	    Required: 対象で状態指定を必須にするか
//	    AppliesTo: 軸が適用される対象名
//	    SourceDefault: sourceで省略され、上位規則にも値がない場合の組み込み状態
//	    SourceDefaultWhen: 条件軸の状態に応じて組み込み状態を選ぶ規則
//	  ]
//	}
type PropertyAxisSpec struct {
	Name              string                          `json:"name"`
	States            []string                        `json:"states"`
	Exclusive         bool                            `json:"exclusive"`
	Required          bool                            `json:"required"`
	AppliesTo         []string                        `json:"applies_to"`
	SourceDefault     string                          `json:"source_default,omitempty"`
	SourceDefaultWhen *PropertySourceDefaultCondition `json:"source_default_when,omitempty"`
}

//	{
//	  責務: [
//	    PropertySourceDefaultCondition: 一つの条件軸からsource既定値を選ぶ規則を保持する
//	  ]
//	  フィールド: [
//	    Axis: 既定値選択の入力にするproperty軸名
//	    States: 入力軸の状態ごとに選択する出力状態
//	  ]
//	}
type PropertySourceDefaultCondition struct {
	Axis   string            `json:"axis"`
	States map[string]string `json:"states"`
}
