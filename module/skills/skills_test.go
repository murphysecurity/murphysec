package skills

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/murphysecurity/murphysec/model"
)

func TestScanDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agent", "hidden", "SKILL.md"), "hidden\r\nskill\r\n")
	writeFile(t, filepath.Join(root, "alpha", "SKILL.md"), "alpha\r\nintro\r\n")
	writeFile(t, filepath.Join(root, "alpha", "docs", "guide.md"), "guide\nbody\n")
	writeFile(t, filepath.Join(root, "beta", "SKILL.md"), "beta\n")
	writeFile(t, filepath.Join(root, "beta", "notes.txt"), "ignored")
	writeFile(t, filepath.Join(root, ".git", "SKILL.md"), "ignored")

	module, ok, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir failed: %v", err)
	}
	if !ok {
		t.Fatal("expected skills module")
	}

	if module.ModuleName != skillsModuleName {
		t.Fatalf("module name mismatch: got %q", module.ModuleName)
	}
	if module.ModuleVersion != skillsModuleUUID {
		t.Fatalf("module version mismatch: got %q", module.ModuleVersion)
	}
	if module.PackageManager != skillsPackageManger {
		t.Fatalf("package manager mismatch: got %q", module.PackageManager)
	}
	if len(module.Dependencies) != 3 {
		t.Fatalf("dependency count mismatch: got %d", len(module.Dependencies))
	}

	hidden := module.Dependencies[0]
	if hidden.CompName != "hidden" {
		t.Fatalf("first skill mismatch: got %q", hidden.CompName)
	}
	if hidden.SkillFiles == nil || len(*hidden.SkillFiles) != 1 {
		t.Fatalf("hidden skill files mismatch: %+v", hidden.SkillFiles)
	}
	assertSkillFileHashes(t, (*hidden.SkillFiles)[0], "SKILL.md", "hidden\r\nskill\r\n", "hidden\nskill\n")

	alpha := module.Dependencies[1]
	if alpha.CompName != "alpha" {
		t.Fatalf("second skill mismatch: got %q", alpha.CompName)
	}
	if alpha.SkillFiles == nil || len(*alpha.SkillFiles) != 2 {
		t.Fatalf("alpha skill files mismatch: %+v", alpha.SkillFiles)
	}
	assertSkillFileHashes(t, (*alpha.SkillFiles)[0], "SKILL.md", "alpha\r\nintro\r\n", "alpha\nintro\n")
	assertSkillFileHashes(t, (*alpha.SkillFiles)[1], "docs/guide.md", "guide\nbody\n", "guide\nbody\n")

	beta := module.Dependencies[2]
	if beta.CompName != "beta" {
		t.Fatalf("third skill mismatch: got %q", beta.CompName)
	}
	if beta.SkillFiles == nil || len(*beta.SkillFiles) != 1 {
		t.Fatalf("beta skill files mismatch: %+v", beta.SkillFiles)
	}
	assertSkillFileHashes(t, (*beta.SkillFiles)[0], "SKILL.md", "beta\n", "beta\n")
}

func TestScanDirNoSkills(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "not a skill")

	_, ok, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir failed: %v", err)
	}
	if ok {
		t.Fatal("expected no skills module")
	}
}

func assertSkillFileHashes(t *testing.T, file model.SkillFile, wantPath, raw, normalized string) {
	t.Helper()

	if file.RelativePath != wantPath {
		t.Fatalf("relative path mismatch: got %q want %q", file.RelativePath, wantPath)
	}
	if len(file.SHA256Hashes) != 2 {
		t.Fatalf("hash count mismatch: got %d", len(file.SHA256Hashes))
	}
	if got := file.SHA256Hashes[0]; got != model.ComputeSHA256Hash([]byte(raw)) {
		t.Fatalf("raw hash mismatch for %s: got %s want %s", wantPath, got.String(), model.ComputeSHA256Hash([]byte(raw)).String())
	}
	if got := file.SHA256Hashes[1]; got != model.ComputeSHA256Hash([]byte(normalized)) {
		t.Fatalf("normalized hash mismatch for %s: got %s want %s", wantPath, got.String(), model.ComputeSHA256Hash([]byte(normalized)).String())
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
