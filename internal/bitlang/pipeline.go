package bitlang

import "fmt"

//	{
//	  責務: [
//	    Pipeline: 隣接順序を検証した変換ステージ列を保持して実行する
//	  ]
//	  フィールド: [
//	    stages: 登録順のパイプラインステージ
//	  ]
//	}
type Pipeline struct {
	stages []Stage
}

//	{
//	  責務: [
//	    NewPipeline: ステージ未登録のパイプラインを作る
//	  ]
//	  戻り値: [
//	    pipeline: ステージを追加できる空のパイプライン
//	  ]
//	}
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

//	{
//	  責務: [
//	    Add: 入出力種別が隣接する次段階をパイプラインへ追加する
//	  ]
//	  処理: [
//	    1: ステージの変換関数と入出力種別を検証する
//	    2: 直前段階の出力との接続と正規順序を検証する
//	    3: 条件を満たす場合だけ末尾へ追加する
//	  ]
//	  引数: [
//	    p: 追加先のパイプライン
//	    stage: 追加する変換段階
//	  ]
//	  戻り値: [
//	    error: ステージ構成が無効な場合の失敗
//	  ]
//	}
func (p *Pipeline) Add(stage Stage) error {
	if !allowedTransition(stage.InputKind, stage.OutputKind) {
		return fmt.Errorf("invalid Bitlang transition: %s -> %s", stage.InputKind, stage.OutputKind)
	}
	if len(p.stages) > 0 {
		previous := p.stages[len(p.stages)-1]
		if previous.OutputKind != stage.InputKind {
			return fmt.Errorf(
				"stage %q cannot follow %q: %s != %s",
				stage.Name,
				previous.Name,
				previous.OutputKind,
				stage.InputKind,
			)
		}
	}
	p.stages = append(p.stages, stage)
	return nil
}

//	{
//	  責務: [
//	    Stages: 内部状態を公開せず登録済みステージ列を返す
//	  ]
//	  引数: [
//	    p: 読み取り対象のパイプライン
//	  ]
//	  戻り値: [
//	    stages: 内部スライスと独立したステージ列の複製
//	  ]
//	}
func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}

//	{
//	  責務: [
//	    Run: 入力成果物から登録済みの変換段階を順に実行する
//	  ]
//	  処理: [
//	    1: 各ステージへ直前の成果物を渡す
//	    2: ステージ失敗時は後続を実行せず返す
//	    3: 最終成果物を返す
//	  ]
//	  引数: [
//	    p: 実行するステージ列
//	    input: 最初のステージへ渡す成果物
//	  ]
//	  戻り値: [
//	    output: 全段階を通過した成果物
//	    error: 入力境界または変換処理の失敗
//	  ]
//	}
func (p *Pipeline) Run(input Artifact) (Artifact, error) {
	current := input
	for _, stage := range p.stages {
		next, err := stage.Run(current)
		if err != nil {
			return Artifact{}, err
		}
		current = next
	}
	return current, nil
}
