package bitlang

import "fmt"

// Transform converts one stage payload into the next stage payload.
type Transform func(any) (any, error)

// Stage represents one explicit Bitlang pipeline conversion.
type Stage struct {
	Name       string
	InputKind  ArtifactKind
	OutputKind ArtifactKind
	Transform  Transform
}

// Run executes one stage after validating its input contract.
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
