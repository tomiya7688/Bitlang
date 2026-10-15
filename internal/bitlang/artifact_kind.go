package bitlang

//	{
//	  責務: [
//	    ArtifactKind: パイプライン成果物の種類を表す
//	  ]
//	  補足: [
//	    値はステージ境界で使う正規の成果物種別を識別する
//	  ]
//	}
type ArtifactKind string

//	{
//	  責務: [
//	    ArtifactKind の定義値: 公開パイプライン成果物のcanonical identityを列挙する
//	  ]
//	  項目: [
//	    ArtifactSource: 入力ソース
//	    ArtifactExplicit: "explicit" を表す正規identity
//	    ArtifactLow: "low" を表す正規identity
//	    ArtifactPreprocessed: ArtifactExplicitを指す旧Go名のsource互換alias。最初のmajor API version v1で削除する。文字列"preprocessed"は受け入れない
//	    ArtifactCompiled: ArtifactLowを指す旧Go名のsource互換alias。最初のmajor API version v1で削除する。文字列"compiled"は受け入れない
//	    ArtifactTreeObject: VM Backend内部表現。公開stage境界には使わない
//	    ArtifactVMAssembly: "vm_assembly" を表す正規identity
//	  ]
//	}
const (
	ArtifactSource       ArtifactKind = "source"
	ArtifactExplicit     ArtifactKind = "explicit"
	ArtifactPreprocessed ArtifactKind = ArtifactExplicit
	ArtifactLow          ArtifactKind = "low"
	ArtifactCompiled     ArtifactKind = ArtifactLow
	ArtifactTreeObject   ArtifactKind = "tree_object"
	ArtifactVMAssembly   ArtifactKind = "vm_assembly"
)

//	{
//	  責務: [
//	    isPublicArtifactKind: 公開stage境界として有効なcanonical identityか判定する
//	  ]
//	  引数: [
//	    kind: 判定対象の成果物identity
//	  ]
//	  戻り値: [
//	    valid: 公開境界として使える場合はtrue
//	  ]
//	}
func isPublicArtifactKind(kind ArtifactKind) bool {
	switch kind {
	case ArtifactSource, ArtifactExplicit, ArtifactLow, ArtifactVMAssembly:
		return true
	default:
		return false
	}
}
