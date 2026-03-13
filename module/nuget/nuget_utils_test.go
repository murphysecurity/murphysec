package nuget

import (
	"os"
	"path/filepath"
	"testing"
)

func touchFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

func TestFindCLNListPreferSln(t *testing.T) {
	root := t.TempDir()
	sln := filepath.Join(root, "a.sln")
	csproj := filepath.Join(root, "src", "a.csproj")
	touchFile(t, sln)
	touchFile(t, csproj)

	got, err := findCLNList(root)
	if err != nil {
		t.Fatalf("findCLNList failed: %v", err)
	}
	if len(got) != 1 || got[0] != sln {
		t.Fatalf("unexpected targets: %#v", got)
	}
}

func TestFindCLNListFallbackProject(t *testing.T) {
	root := t.TempDir()
	csproj := filepath.Join(root, "src", "a.csproj")
	touchFile(t, csproj)

	got, err := findCLNList(root)
	if err != nil {
		t.Fatalf("findCLNList failed: %v", err)
	}
	if len(got) != 1 || got[0] != csproj {
		t.Fatalf("unexpected targets: %#v", got)
	}
}

func TestValidateBuildTargetsRejectInvalid(t *testing.T) {
	root := t.TempDir()
	if err := validateBuildTargets([]string{root}); err == nil {
		t.Fatal("expected error for directory target, got nil")
	}
}

func TestFindMissingProjectReferences(t *testing.T) {
	root := t.TempDir()
	projectPath := filepath.Join(root, "Web", "Web.csproj")
	if err := os.MkdirAll(filepath.Dir(projectPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	content := `<Project Sdk="Microsoft.NET.Sdk.Web">
  <ItemGroup>
    <ProjectReference Include="../ApplicationCore/ApplicationCore.csproj" />
    <ProjectReference Include="../Infrastructure/Infrastructure.csproj" />
  </ItemGroup>
</Project>`
	if err := os.WriteFile(projectPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write csproj failed: %v", err)
	}
	touchFile(t, filepath.Join(root, "ApplicationCore", "ApplicationCore.csproj"))
	missing, err := findMissingProjectReferences(projectPath)
	if err != nil {
		t.Fatalf("findMissingProjectReferences failed: %v", err)
	}
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing reference, got %v", missing)
	}
	want := filepath.Join(root, "Infrastructure", "Infrastructure.csproj")
	if missing[0] != want {
		t.Fatalf("unexpected missing reference: got %s want %s", missing[0], want)
	}
}

func TestFindMissingProjectReferencesNone(t *testing.T) {
	root := t.TempDir()
	projectPath := filepath.Join(root, "Web", "Web.csproj")
	if err := os.MkdirAll(filepath.Dir(projectPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	content := `<Project Sdk="Microsoft.NET.Sdk.Web">
  <ItemGroup>
    <ProjectReference Include="../ApplicationCore/ApplicationCore.csproj" />
  </ItemGroup>
</Project>`
	if err := os.WriteFile(projectPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write csproj failed: %v", err)
	}
	touchFile(t, filepath.Join(root, "ApplicationCore", "ApplicationCore.csproj"))
	missing, err := findMissingProjectReferences(projectPath)
	if err != nil {
		t.Fatalf("findMissingProjectReferences failed: %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("expected no missing references, got %v", missing)
	}
}
