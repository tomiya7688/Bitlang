package bitlang

// ArtifactKind identifies one canonical Bitlang pipeline representation.
type ArtifactKind string

// Canonical Bitlang pipeline artifact kinds.
const (
	ArtifactSource       ArtifactKind = "source"
	ArtifactPreprocessed ArtifactKind = "preprocessed"
	ArtifactCompiled     ArtifactKind = "compiled"
	ArtifactTreeObject   ArtifactKind = "tree_object"
	ArtifactVMAssembly   ArtifactKind = "vm_assembly"
)
