package bitlang

func allowedTransition(from ArtifactKind, to ArtifactKind) bool {
	switch {
	case from == ArtifactSource && to == ArtifactExplicit:
		return true
	case from == ArtifactExplicit && to == ArtifactLow:
		return true
	case from == ArtifactLow && to == ArtifactVMAssembly:
		return true
	default:
		return false
	}
}
