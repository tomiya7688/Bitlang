package bitlang

func allowedTransition(from ArtifactKind, to ArtifactKind) bool {
	switch {
	case from == ArtifactSource && to == ArtifactPreprocessed:
		return true
	case from == ArtifactPreprocessed && to == ArtifactCompiled:
		return true
	case from == ArtifactCompiled && to == ArtifactTreeObject:
		return true
	case from == ArtifactTreeObject && to == ArtifactVMAssembly:
		return true
	default:
		return false
	}
}
