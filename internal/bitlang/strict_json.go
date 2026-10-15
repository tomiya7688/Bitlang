package bitlang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

//	{
//	  責務: [
//	    DecodeStrictJSON: 重複フィールドや後続値を許さず単一JSON値をデコードする
//	  ]
//	  引数: [
//	    data: デコード対象のJSONバイト列
//	    output: 値を書き込む出力先
//	  ]
//	  戻り値: [
//	    error: 重複・未知フィールド・後続値・不正構文の失敗
//	  ]
//	}
func DecodeStrictJSON(data []byte, output any) error {
	return decodeSingleJSONValue(data, output)
}

//	{
//	  責務: [
//	    ValidateJSONFieldNames: JSON オブジェクトのキーを許可された表記に限定する
//	  ]
//	  引数: [
//	    fields: 検査対象の JSON オブジェクト
//	    allowed: 許可するキーの完全一致表記
//	  ]
//	  戻り値: [
//	    error: 未知または大文字小文字が異なるキーの失敗
//	  ]
//	}
func ValidateJSONFieldNames(fields map[string]json.RawMessage, allowed ...string) error {
	for field := range fields {
		known := false
		for _, candidate := range allowed {
			if field == candidate {
				known = true
				break
			}
		}
		if !known {
			return fmt.Errorf("unknown or incorrectly cased JSON field %q", field)
		}
	}
	return nil
}

//	{
//	  責務: [
//	    decodeSingleJSONValue: 未知フィールドや後続値を許さず JSON 値をデコードする
//	  ]
//	  処理: [
//	    1: 重複オブジェクトキーを事前検出する
//	    2: 未知フィールドを拒否して出力先へデコードする
//	    3: 末尾が空白以外を含まない単一値か確認する
//	  ]
//	  引数: [
//	    data: デコード対象の JSON バイト列
//	    output: 値を書き込む出力先
//	  ]
//	  戻り値: [
//	    error: 重複・未知フィールド・複数値・不正構文の失敗
//	  ]
//	}
func decodeSingleJSONValue(data []byte, output any) error {
	if err := rejectDuplicateJSONFields(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(output); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err != nil {
			return fmt.Errorf("decode trailing JSON data: %w", err)
		}
		return fmt.Errorf("multiple JSON values are not allowed")
	}
	return nil
}

func rejectDuplicateJSONFields(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := scanJSONValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err != nil {
			return fmt.Errorf("decode trailing JSON data: %w", err)
		}
		return fmt.Errorf("multiple JSON values are not allowed")
	}
	return nil
}

func scanJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return nil
	}
	switch delimiter {
	case '{':
		return scanJSONObject(decoder)
	case '[':
		return scanJSONArray(decoder)
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
	}
}

func scanJSONObject(decoder *json.Decoder) error {
	keys := make(map[string]bool)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("JSON object key must be a string")
		}
		canonical := strings.ToLower(key)
		if keys[canonical] {
			return fmt.Errorf("duplicate JSON object field %q", key)
		}
		keys[canonical] = true
		if err := scanJSONValue(decoder); err != nil {
			return err
		}
	}
	return consumeJSONDelimiter(decoder, '}')
}

func scanJSONArray(decoder *json.Decoder) error {
	for decoder.More() {
		if err := scanJSONValue(decoder); err != nil {
			return err
		}
	}
	return consumeJSONDelimiter(decoder, ']')
}

func consumeJSONDelimiter(decoder *json.Decoder, expected json.Delim) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if token != expected {
		return fmt.Errorf("expected JSON delimiter %q", expected)
	}
	return nil
}
