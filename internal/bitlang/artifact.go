package bitlang

//	{
//	  責務: [
//	    Artifact: パイプライン段階間で成果物とその種別を一緒に渡す
//	  ]
//	  フィールド: [
//	    Kind: Payload の成果物種別
//	    Payload: 次段階へ渡す成果物本体
//	  ]
//	}
type Artifact struct {
	Kind    ArtifactKind
	Payload any
}
