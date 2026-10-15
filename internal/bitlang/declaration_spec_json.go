package bitlang

import "encoding/json"

//	{
//	  責務: [
//	    validateDeclarationSpecificationJSONFields: 宣言仕様 JSON のキー名を schema に照合する
//	  ]
//	  処理: [
//	    1: ルートの version と kinds を検査する
//	    2: 各宣言種別のキーを検査する
//	  ]
//	  引数: [
//	    data: 宣言仕様 JSON のバイト列
//	  ]
//	  戻り値: [
//	    error: 未知・誤表記・不正な JSON の失敗
//	  ]
//	}
func validateDeclarationSpecificationJSONFields(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if err := ValidateJSONFieldNames(fields, "version", "kinds"); err != nil {
		return err
	}
	var kinds []json.RawMessage
	if err := json.Unmarshal(fields["kinds"], &kinds); err != nil {
		return err
	}
	for _, kindData := range kinds {
		var kindFields map[string]json.RawMessage
		if err := json.Unmarshal(kindData, &kindFields); err != nil {
			return err
		}
		if err := ValidateJSONFieldNames(kindFields, "name", "property_target", "terminator", "layout", "contexts"); err != nil {
			return err
		}
	}
	return nil
}
