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

