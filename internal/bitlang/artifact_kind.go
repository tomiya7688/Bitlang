package bitlang

// ArtifactKind identifies one canonical Bitlang pipeline representation.
type ArtifactKind string

const (
	ArtifactSource       ArtifactKind = "source"
	ArtifactPreprocessed ArtifactKind = "preprocessed"
	ArtifactCompiled     ArtifactKind = "compiled"
	ArtifactTreeObject   ArtifactKind = "tree_object"
	ArtifactVMAssembly   ArtifactKind = "vm_assembly"
)
