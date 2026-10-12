package bitlangci

import (
	"encoding/json"

	"github.com/tomiya7688/Bitlang/internal/bitlang"
)

//	{
//	  責務: [
//	    validateFamilyPipelineManifestJSONFields: family manifest の全キーを schema 表記に照合する
//	  ]
//	  処理: [
//	    1: ルートのキー名を検証する
//	    2: 各配列要素のキー名を検証する
//	  ]
//	  引数: [
//	    data: family manifest JSON のバイト列
//	  ]
//	  戻り値: [
//	    error: 未知・誤表記・不正な JSON の失敗
//	  ]
//	}
func validateFamilyPipelineManifestJSONFields(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if err := bitlang.ValidateJSONFieldNames(fields,
		"schema_version", "contract_version", "family", "gate_status", "repositories",
		"public_stages", "transitions", "internal_representations", "legacy_aliases",
	); err != nil {
		return err
	}
	checks := []struct {
		field   string
		allowed []string
	}{
		{"repositories", []string{"id", "name", "github", "role"}},
		{"public_stages", []string{"id", "name", "repository"}},
		{"transitions", []string{"from", "to", "fixture_path", "fixture_status"}},
		{"internal_representations", []string{"id", "owner_repository", "optional", "public_stage"}},
		{"legacy_aliases", []string{"literal", "canonical", "go_identifier", "remove_identifier_at_major", "accept_literal"}},
	}
	for _, check := range checks {
		if err := validateFamilyPipelineObjectFields(fields[check.field], check.allowed...); err != nil {
			return err
		}
	}
	return nil
}

func validateFamilyPipelineObjectFields(data json.RawMessage, allowed ...string) error {
	var objects []json.RawMessage
	if err := json.Unmarshal(data, &objects); err != nil {
		return err
	}
	for _, object := range objects {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(object, &fields); err != nil {
			return err
		}
		if err := bitlang.ValidateJSONFieldNames(fields, allowed...); err != nil {
			return err
		}
	}
	return nil
}
