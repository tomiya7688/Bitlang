package bitlang

import "testing"

func TestPipelineRunsInOrder(t *testing.T) {
	pipeline := NewPipeline()

	stages := []Stage{
		{
			Name:       "preprocessor",
			InputKind:  ArtifactSource,
			OutputKind: ArtifactPreprocessed,
			Transform: func(v any) (any, error) { return v.(string) + "|pre", nil },
		},
		{
			Name:       "compiler",
			InputKind:  ArtifactPreprocessed,
			OutputKind: ArtifactCompiled,
			Transform: func(v any) (any, error) { return v.(string) + "|compiled", nil },
		},
	}

	for _, stage := range stages {
		if err := pipeline.Add(stage); err != nil {
			t.Fatal(err)
		}
	}

	out, err := pipeline.Run(Artifact{Kind: ArtifactSource, Payload: "source"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Kind != ArtifactCompiled {
		t.Fatalf("got kind %q, want %q", out.Kind, ArtifactCompiled)
	}
	if out.Payload != "source|pre|compiled" {
		t.Fatalf("got payload %q", out.Payload)
	}
}

func TestPipelineRejectsSkippedStage(t *testing.T) {
	pipeline := NewPipeline()
	err := pipeline.Add(Stage{
		Name:       "invalid",
		InputKind:  ArtifactSource,
		OutputKind: ArtifactCompiled,
		Transform:  func(v any) (any, error) { return v, nil },
	})
	if err == nil {
		t.Fatal("expected invalid transition error")
	}
}
