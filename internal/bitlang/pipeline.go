package bitlang

import "fmt"

type ArtifactKind string

const (
	ArtifactSource       ArtifactKind = "source"
	ArtifactPreprocessed ArtifactKind = "preprocessed"
	ArtifactCompiled     ArtifactKind = "compiled"
	ArtifactTreeObject   ArtifactKind = "tree_object"
	ArtifactVMAssembly   ArtifactKind = "vm_assembly"
)

type Artifact struct {
	Kind    ArtifactKind
	Payload any
}

type Transform func(any) (any, error)

type Stage struct {
	Name       string
	InputKind  ArtifactKind
	OutputKind ArtifactKind
	Transform  Transform
}

func (s Stage) Run(input Artifact) (Artifact, error) {
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

type Pipeline struct {
	stages []Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func allowedTransition(from, to ArtifactKind) bool {
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

func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}

func (p *Pipeline) Run(input Artifact) (Artifact, error) {
	current := input
	var err error
	for _, stage := range p.stages {
		current, err = stage.Run(current)
		if err != nil {
			return Artifact{}, err
		}
	}
	return current, nil
}
