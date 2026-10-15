package bitlangci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tomiya7688/Bitlang/internal/bitlang"
)

//	{
//	  責務: [
//	    familyPipelineManifest: Bitlang familyのstage、repository、移行規則を表す
//	  ]
//	  フィールド: [
//	    SchemaVersion: JSON manifest形式のversion
//	    ContractVersion: family間で共有するpipeline契約のversion
//	    Family: manifestが所有するlanguage family名
//	    GateStatus: cross-repository conformance gateの実装状態
//	    Repositories: 各pipeline責務を持つGitHub repository
//	    PublicStages: machine-readableな公開stage identity
//	    Transitions: 隣接stageとfixtureの明示的な実装状態
//	    InternalRepresentations: 公開boundaryではない内部表現
//	    LegacyAliases: 旧名の変換規則と削除version
//	  ]
//	}
type familyPipelineManifest struct {
	SchemaVersion   int    `json:"schema_version"`
	ContractVersion int    `json:"contract_version"`
	Family          string `json:"family"`
	GateStatus      string `json:"gate_status"`
	Repositories    []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		GitHub string `json:"github"`
		Role   string `json:"role"`
	} `json:"repositories"`
	PublicStages []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Repository string `json:"repository"`
	} `json:"public_stages"`
	Transitions []struct {
		From          string `json:"from"`
		To            string `json:"to"`
		FixturePath   string `json:"fixture_path"`
		FixtureStatus string `json:"fixture_status"`
	} `json:"transitions"`
	InternalRepresentations []struct {
		ID              string `json:"id"`
		OwnerRepository string `json:"owner_repository"`
		Optional        bool   `json:"optional"`
		PublicStage     bool   `json:"public_stage"`
	} `json:"internal_representations"`
	LegacyAliases []struct {
		Literal                 string `json:"literal"`
		Canonical               string `json:"canonical"`
		GoIdentifier            string `json:"go_identifier"`
		RemoveIdentifierAtMajor string `json:"remove_identifier_at_major"`
		AcceptLiteral           bool   `json:"accept_literal"`
	} `json:"legacy_aliases"`
}

//	{
//	  責務: [
//	    loadFamilyPipelineManifest: repository manifestを厳密に読み取り現在のBitlang stage contractと照合する
//	  ]
//	  引数: [
//	    root: 検査対象repositoryのroot directory
//	  ]
//	  戻り値: [
//	    error: file、schema、実装contract、または未実装fixtureの失敗
//	  ]
//	}
func loadFamilyPipelineManifest(root string) error {
	path := filepath.Join(root, "spec", "family-pipeline.json")
	// #nosec G304 -- path is derived from the repository root and a fixed manifest path.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read family pipeline manifest: %w", err)
	}
	if err := validateFamilyPipelineManifestJSONFields(data); err != nil {
		return fmt.Errorf("decode family pipeline manifest: %w", err)
	}
	var manifest familyPipelineManifest
	if err := bitlang.DecodeStrictJSON(data, &manifest); err != nil {
		return fmt.Errorf("decode family pipeline manifest: %w", err)
	}
	if err := manifest.validate(root); err != nil {
		return fmt.Errorf("validate family pipeline manifest: %w", err)
	}
	return nil
}

//	{
//	  責務: [
//	    validate: repository manifestと実装中のcanonical stage contractの一致を確認する
//	  ]
//	  引数: [
//	    root: 検査対象repositoryのroot directory
//	  ]
//	  戻り値: [
//	    error: 不一致、またはfixture未実装でgateを成功扱いできない場合の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validate(root string) error {
	if err := manifest.validateIdentity(); err != nil {
		return err
	}
	if err := manifest.validateRepositories(); err != nil {
		return err
	}
	if err := manifest.validateStages(); err != nil {
		return err
	}
	unfinishedFixtures, err := manifest.validateTransitions(root)
	if err != nil {
		return err
	}
	if err := manifest.validateInternalRepresentations(); err != nil {
		return err
	}
	if err := manifest.validateAliases(); err != nil {
		return err
	}
	return validateFamilyGateReadiness(manifest.GateStatus, unfinishedFixtures)
}

//	{
//	  責務: [
//	    validateIdentity: manifestのschemaとfamily identityを検証する
//	  ]
//	  戻り値: [
//	    error: 未対応versionまたはfamily名の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validateIdentity() error {
	if manifest.SchemaVersion != 1 || manifest.ContractVersion != 1 || manifest.Family != "Bitlang" {
		return fmt.Errorf("unsupported family manifest identity: schema=%d contract=%d family=%q", manifest.SchemaVersion, manifest.ContractVersion, manifest.Family)
	}
	if manifest.GateStatus != "implemented" && manifest.GateStatus != "not_implemented" {
		return fmt.Errorf("invalid gate_status %q", manifest.GateStatus)
	}
	return nil
}

//	{
//	  責務: [
//	    validateRepositories: manifest内のrepository一覧がBitlang familyと一致することを検証する
//	  ]
//	  戻り値: [
//	    error: repository不足、余分、またはidentity不一致の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validateRepositories() error {
	expected := [][4]string{
		{"source", "Bitlang", "tomiya7688/Bitlang", "source"},
		{"explicit", "Bitlang Explicit", "tomiya7688/Bitlang-Explicit", "explicit"},
		{"low", "Bitlang Low", "tomiya7688/Bitlang_low", "low"},
		{"vm", "Bitlang VM", "tomiya7688/Bitlang-VM", "vm_backend"},
	}
	if len(manifest.Repositories) != len(expected) {
		return fmt.Errorf("family manifest must list exactly %d repositories", len(expected))
	}
	for index, want := range expected {
		got := manifest.Repositories[index]
		if got.ID != want[0] || got.Name != want[1] || got.GitHub != want[2] || got.Role != want[3] {
			return fmt.Errorf("repository entry %d does not match the canonical family manifest", index)
		}
	}
	return nil
}

//	{
//	  責務: [
//	    validateStages: manifestの公開stage identityがGo実装と一致することを検証する
//	  ]
//	  戻り値: [
//	    error: 公開stageの不足、余分、またはidentity不一致の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validateStages() error {
	expected := [][3]string{
		{string(bitlang.ArtifactSource), "Bitlang source", "source"},
		{string(bitlang.ArtifactExplicit), "Bitlang Explicit", "explicit"},
		{string(bitlang.ArtifactLow), "Bitlang Low", "low"},
		{string(bitlang.ArtifactVMAssembly), "Bitlang VM Assembly", "vm"},
	}
	if len(manifest.PublicStages) != len(expected) {
		return fmt.Errorf("family manifest must list exactly %d public stages", len(expected))
	}
	for index, want := range expected {
		got := manifest.PublicStages[index]
		if got.ID != want[0] || got.Name != want[1] || got.Repository != want[2] {
			return fmt.Errorf("public stage entry %d does not match the canonical stage identities", index)
		}
	}
	return nil
}

//	{
//	  責務: [
//	    validateTransitions: canonical隣接遷移とPipeline実装の一致およびfixture状態を検証する
//	  ]
//	  戻り値: [
//	    unfinished: not_implemented状態のfixture遷移
//	    error: 遷移graphまたはPipeline実装の不一致
//	  ]
//	}
func (manifest familyPipelineManifest) validateTransitions(root string) ([]string, error) {
	expected := [][3]string{
		{string(bitlang.ArtifactSource), string(bitlang.ArtifactExplicit), "fixtures/family/source-explicit.json"},
		{string(bitlang.ArtifactExplicit), string(bitlang.ArtifactLow), "fixtures/family/explicit-low.json"},
		{string(bitlang.ArtifactLow), string(bitlang.ArtifactVMAssembly), "fixtures/family/low-vm-assembly.json"},
	}
	if len(manifest.Transitions) != len(expected) {
		return nil, fmt.Errorf("family manifest must list exactly %d adjacent public transitions", len(expected))
	}
	var pipeline bitlang.Pipeline
	var unfinished []string
	for index, want := range expected {
		got := manifest.Transitions[index]
		if got.From != want[0] || got.To != want[1] || got.FixturePath != want[2] {
			return nil, fmt.Errorf("transition entry %d is not the canonical adjacent transition", index)
		}
		if got.FixtureStatus != "implemented" && got.FixtureStatus != "not_implemented" {
			return nil, fmt.Errorf("transition %s -> %s has invalid fixture_status %q", got.From, got.To, got.FixtureStatus)
		}
		if got.FixtureStatus == "not_implemented" {
			unfinished = append(unfinished, got.From+" -> "+got.To)
		} else if err := validateFamilyFixturePath(root, got.FixturePath); err != nil {
			return nil, fmt.Errorf("transition %s -> %s fixture: %w", got.From, got.To, err)
		}
		stage := bitlang.Stage{
			Name:       got.From + " -> " + got.To,
			InputKind:  bitlang.ArtifactKind(got.From),
			OutputKind: bitlang.ArtifactKind(got.To),
			Transform:  func(value any) (any, error) { return value, nil },
		}
		if err := pipeline.Add(stage); err != nil {
			return nil, fmt.Errorf("implementation rejects manifest transition %s: %w", stage.Name, err)
		}
	}
	return unfinished, nil
}

//	{
//	  責務: [
//	    validateFamilyFixturePath: 実装済みと宣言されたfamily fixtureがrepository内に存在することを確認する
//	  ]
//	  引数: [
//	    root: 検査対象repositoryのroot directory
//	    relativePath: manifestで固定されたfixture path
//	  ]
//	  戻り値: [
//	    error: 不正path、欠落、または通常fileでない場合の失敗
//	  ]
//	}
func validateFamilyFixturePath(root string, relativePath string) error {
	if !filepath.IsLocal(relativePath) {
		return fmt.Errorf("fixture path must remain local to the repository")
	}
	rootPath, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve repository root: %w", err)
	}
	path := filepath.Join(rootPath, filepath.FromSlash(relativePath))
	fixturePath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("resolve fixture path %q: %w", relativePath, err)
	}
	localPath, err := filepath.Rel(rootPath, fixturePath)
	if err != nil {
		return fmt.Errorf("compare fixture path to repository root: %w", err)
	}
	if localPath == ".." || strings.HasPrefix(localPath, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("fixture path %q resolves outside the repository", relativePath)
	}
	// #nosec G304 -- path is resolved inside the repository before reading file metadata.
	info, err := os.Stat(fixturePath)
	if err != nil {
		return fmt.Errorf("read fixture path %q: %w", relativePath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("fixture path %q is not a regular file", relativePath)
	}
	return nil
}

//	{
//	  責務: [
//	    validateInternalRepresentations: 内部表現が公開stageとして誤認されないことを検証する
//	  ]
//	  戻り値: [
//	    error: 内部表現のscope不一致の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validateInternalRepresentations() error {
	if len(manifest.InternalRepresentations) != 1 {
		return fmt.Errorf("family manifest must describe exactly one current internal representation")
	}
	internal := manifest.InternalRepresentations[0]
	if internal.ID != string(bitlang.ArtifactTreeObject) || internal.OwnerRepository != "vm" || !internal.Optional || internal.PublicStage {
		return fmt.Errorf("TreeObject must remain an optional VM-internal representation outside public stages")
	}
	return nil
}

//	{
//	  責務: [
//	    validateAliases: 旧artifact識別子の変換規則とGo alias実装を照合する
//	  ]
//	  戻り値: [
//	    error: 旧名の変換先、削除version、または実装不一致の失敗
//	  ]
//	}
func (manifest familyPipelineManifest) validateAliases() error {
	expected := [][4]string{
		{"preprocessed", string(bitlang.ArtifactExplicit), "ArtifactPreprocessed", "v1"},
		{"compiled", string(bitlang.ArtifactLow), "ArtifactCompiled", "v1"},
	}
	if len(manifest.LegacyAliases) != len(expected) {
		return fmt.Errorf("family manifest must describe exactly %d legacy aliases", len(expected))
	}
	for index, want := range expected {
		got := manifest.LegacyAliases[index]
		if got.Literal != want[0] || got.Canonical != want[1] || got.GoIdentifier != want[2] || got.RemoveIdentifierAtMajor != want[3] || got.AcceptLiteral {
			return fmt.Errorf("legacy alias entry %d does not match the canonical migration policy", index)
		}
	}
	if bitlang.ArtifactPreprocessed != bitlang.ArtifactExplicit || bitlang.ArtifactCompiled != bitlang.ArtifactLow {
		return fmt.Errorf("Go source aliases do not resolve to canonical artifact identities")
	}
	return nil
}

//	{
//	  責務: [
//	    validateFamilyGateReadiness: 未実装fixtureやgate状態をCI成功扱いしない
//	  ]
//	  引数: [
//	    gateStatus: manifestが宣言するgate実装状態
//	    unfinishedFixtures: 未実装として列挙された遷移fixture
//	  ]
//	  戻り値: [
//	    error: family gate、fixture runnerが未完了または状態不整合の場合の失敗
//	  ]
//	}
func validateFamilyGateReadiness(gateStatus string, unfinishedFixtures []string) error {
	if len(unfinishedFixtures) > 0 {
		return fmt.Errorf("cross-repository acceptance fixtures are not implemented: %s", strings.Join(unfinishedFixtures, ", "))
	}
	if gateStatus != "implemented" {
		return fmt.Errorf("cross-repository acceptance gate is not implemented")
	}
	return fmt.Errorf("cross-repository acceptance fixture runner is not implemented")
}
