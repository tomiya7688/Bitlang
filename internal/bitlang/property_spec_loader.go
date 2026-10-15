package bitlang

import (
	"fmt"
)

//	{
//	  責務: [
//	    LoadPropertySpecification: 機械可読なプロパティ仕様を厳密に読み込み検証する
//	  ]
//	  処理: [
//	    1: 未知・重複フィールドと余分な JSON 値を拒否する
//	    2: 軸・状態・対象の仕様整合性を検証する
//	  ]
//	  引数: [
//	    data: プロパティ仕様の JSON バイト列
//	  ]
//	  戻り値: [
//	    spec: 検証済みのプロパティ仕様
//	    error: デコードまたは検証に失敗した場合の内容
//	  ]
//	}
func LoadPropertySpecification(data []byte) (PropertySpecification, error) {
	if err := validatePropertySpecificationJSONFields(data); err != nil {
		return PropertySpecification{}, fmt.Errorf("decode property specification: %w", err)
	}
	var spec PropertySpecification
	if err := decodeSingleJSONValue(data, &spec); err != nil {
		return PropertySpecification{}, fmt.Errorf("decode property specification: %w", err)
	}
	if err := validatePropertySpecification(spec); err != nil {
		return PropertySpecification{}, err
	}
	return spec, nil
}
