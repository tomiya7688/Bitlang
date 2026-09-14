package bitlang

// Artifact carries one value between canonical Bitlang pipeline stages.
type Artifact struct {
	Kind    ArtifactKind
	Payload any
}
