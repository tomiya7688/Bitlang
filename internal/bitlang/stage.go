package bitlang

import "fmt"

//	{
//	  責務: [
//	    Transform: 一段階分の成果物変換処理を表す
//	  ]
//	  引数: [
//	    payload: 入力成果物の本体
//	  ]
//	  戻り値: [
//	    output: 変換後の成果物本体
//	    error: 変換処理の失敗
//	  ]
//	}
type Transform func(any) (any, error)

//	{
//	  責務: [
//	    Stage: 入力種別・出力種別・変換処理を一段階として定義する
//	  ]
//	  フィールド: [
//	    Name: 診断に使う段階名
//	    InputKind: 受け入れる成果物種別
//	    OutputKind: 生成する成果物種別
//	    Transform: 成果物本体の変換関数
//	  ]
//	}
type Stage struct {
	Name       string
	InputKind  ArtifactKind
	OutputKind ArtifactKind
	Transform  Transform
}

//	{
//	  責務: [
//	    Run: 入力境界を検証して一段階の変換を実行する
//	  ]
//	  処理: [
//	    1: 入力種別と変換関数の存在を確認する
//	    2: 入力本体を変換関数へ渡す
//	    3: 出力種別を付与した成果物を返す
//	  ]
//	  引数: [
//	    s: 入出力契約と変換関数を持つ段階
//	    input: 現段階に渡す成果物
//	  ]
//	  戻り値: [
//	    output: 変換後の出力種別付き成果物
//	    error: 入力契約または変換処理の失敗
//	  ]
//	}
func (s Stage) Run(input Artifact) (Artifact, error) {
	if !isPublicArtifactKind(s.InputKind) || !isPublicArtifactKind(s.OutputKind) {
		return Artifact{}, fmt.Errorf(
			"stage %q uses a non-canonical public artifact kind: %q -> %q",
			s.Name,
			s.InputKind,
			s.OutputKind,
		)
	}
	if !allowedTransition(s.InputKind, s.OutputKind) {
		return Artifact{}, fmt.Errorf("invalid Bitlang transition: %s -> %s", s.InputKind, s.OutputKind)
	}
	if input.Kind != s.InputKind {
		return Artifact{}, fmt.Errorf("stage %q expects %q, got %q", s.Name, s.InputKind, input.Kind)
	}
	if s.Transform == nil {
		return Artifact{}, fmt.Errorf("stage %q has no transform", s.Name)
	}
	payload, err := s.Transform(input.Payload)
	if err != nil {
		return Artifact{}, fmt.Errorf("stage %q: %w", s.Name, err)
	}
	return Artifact{Kind: s.OutputKind, Payload: payload}, nil
}
