package nuget

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNugetRestoreBinlogPath(t *testing.T) {
	root := t.TempDir()
	solutionPath := filepath.Join(root, "src", "App.sln")

	got := nugetRestoreBinlogPath(solutionPath)
	want := filepath.Join(root, "src", nugetRestoreBinlogDir, "App.restore.binlog")
	if got != want {
		t.Fatalf("unexpected binlog path: got %s want %s", got, want)
	}
}

func TestNugetRestoreTmpBinlogPath(t *testing.T) {
	root := t.TempDir()
	solutionPath := filepath.Join(root, "src", "App.sln")

	got := nugetRestoreTmpBinlogPath(solutionPath)
	if !strings.HasPrefix(got, filepath.Join(os.TempDir(), nugetTmpBinlogDir)+string(filepath.Separator)) {
		t.Fatalf("unexpected tmp binlog root: %s", got)
	}
	if !strings.HasSuffix(got, filepath.Join("App.restore.binlog")) {
		t.Fatalf("unexpected tmp binlog filename: %s", got)
	}
}

func TestDotnetRestoreArgsEnableDetailedBinlog(t *testing.T) {
	root := t.TempDir()
	solutionPath := filepath.Join(root, "src", "App.csproj")

	args := dotnetRestoreArgs(solutionPath)
	if len(args) < 5 {
		t.Fatalf("unexpected restore args length: %v", args)
	}
	if args[0] != "restore" || args[1] != solutionPath {
		t.Fatalf("unexpected restore command prefix: %v", args)
	}
	if args[2] != "-v" || args[3] != "detailed" {
		t.Fatalf("expected detailed verbosity, got %v", args)
	}
	if !strings.Contains(args[4], "/bl:") {
		t.Fatalf("expected binary log argument, got %s", args[4])
	}
	if !strings.Contains(args[4], "ProjectImports=Embed") {
		t.Fatalf("expected embedded project imports in binlog, got %s", args[4])
	}
	if !strings.Contains(args[4], "App.restore.binlog") {
		t.Fatalf("expected restore binlog filename, got %s", args[4])
	}
	if runtime.GOOS == "linux" && !strings.Contains(strings.Join(args, " "), "EnableWindowsTargeting=true") {
		t.Fatalf("expected linux restore args to include EnableWindowsTargeting=true, got %v", args)
	}
}

func TestArchiveNugetRestoreBinlog(t *testing.T) {
	root := t.TempDir()
	solutionPath := filepath.Join(root, "src", "App.csproj")
	srcPath := nugetRestoreBinlogPath(solutionPath)
	if err := os.MkdirAll(filepath.Dir(srcPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	wantContent := []byte("binlog")
	if err := os.WriteFile(srcPath, wantContent, 0o644); err != nil {
		t.Fatalf("write source binlog failed: %v", err)
	}

	got, err := archiveNugetRestoreBinlog(nil, solutionPath)
	if err != nil {
		t.Fatalf("archive binlog failed: %v", err)
	}
	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatalf("read archived binlog failed: %v", err)
	}
	if string(data) != string(wantContent) {
		t.Fatalf("unexpected archived content: got %q want %q", string(data), string(wantContent))
	}
}
