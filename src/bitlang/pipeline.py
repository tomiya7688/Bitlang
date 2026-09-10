"""Explicit Bitlang lowering pipeline primitives."""

from __future__ import annotations

from dataclasses import dataclass
from enum import Enum
from typing import Callable, Generic, TypeVar

Payload = TypeVar("Payload")


class ArtifactKind(str, Enum):
    SOURCE = "source"
    PREPROCESSED = "preprocessed"
    COMPILED = "compiled"
    TREE_OBJECT = "tree_object"
    VM_ASSEMBLY = "vm_assembly"


@dataclass(frozen=True, slots=True)
class Artifact(Generic[Payload]):
    kind: ArtifactKind
    payload: Payload


class PipelineError(RuntimeError):
    pass


@dataclass(frozen=True, slots=True)
class Stage:
    name: str
    input_kind: ArtifactKind
    output_kind: ArtifactKind
    transform: Callable[[object], object]

    def run(self, artifact: Artifact[object]) -> Artifact[object]:
        if artifact.kind is not self.input_kind:
            raise PipelineError(
                f"stage {self.name!r} expects {self.input_kind.value!r}, "
                f"got {artifact.kind.value!r}"
            )
        return Artifact(self.output_kind, self.transform(artifact.payload))


_ALLOWED_TRANSITIONS: set[tuple[ArtifactKind, ArtifactKind]] = {
    (ArtifactKind.SOURCE, ArtifactKind.PREPROCESSED),
    (ArtifactKind.PREPROCESSED, ArtifactKind.COMPILED),
    (ArtifactKind.COMPILED, ArtifactKind.TREE_OBJECT),
    (ArtifactKind.TREE_OBJECT, ArtifactKind.VM_ASSEMBLY),
}


class Pipeline:
    """A checked sequence of explicit Bitlang lowering stages.

    The pipeline validates adjacent artifact kinds when stages are registered.
    The VM and architecture translators intentionally live beyond VM assembly:
    VM assembly is a stable output boundary rather than an implicit extra stage.
    """

    def __init__(self) -> None:
        self._stages: list[Stage] = []

    @property
    def stages(self) -> tuple[Stage, ...]:
        return tuple(self._stages)

    def add(self, stage: Stage) -> "Pipeline":
        transition = (stage.input_kind, stage.output_kind)
        if transition not in _ALLOWED_TRANSITIONS:
            raise PipelineError(
                f"invalid Bitlang transition: {stage.input_kind.value} "
                f"-> {stage.output_kind.value}"
            )

        if self._stages:
            previous = self._stages[-1]
            if previous.output_kind is not stage.input_kind:
                raise PipelineError(
                    f"stage {stage.name!r} cannot follow {previous.name!r}: "
                    f"{previous.output_kind.value} != {stage.input_kind.value}"
                )

        self._stages.append(stage)
        return self

    def run(self, artifact: Artifact[object]) -> Artifact[object]:
        current = artifact
        for stage in self._stages:
            current = stage.run(current)
        return current
