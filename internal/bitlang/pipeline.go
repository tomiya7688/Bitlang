package bitlang

import "fmt"

// Pipeline stores and executes an ordered list of Bitlang stages.
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates an empty pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends one stage after validating canonical order and adjacency.
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

// Stages returns a copy of the configured stage sequence.
func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}

// Run executes the configured stage sequence from the supplied artifact.
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
