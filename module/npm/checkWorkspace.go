package npm

import (
	"context"
	"os"
	"path/filepath"
)

func getWorkspace(packagePath string, mp map[string]bool) error {
	by, err := os.ReadFile(packagePath)
	if err != nil {
		return err
	}
	pkg, err := parsePkgFile(by)
	if err != nil {
		return err
	}

	for _, k := range pkg.Workspaces {
		p := filepath.Join(filepath.Dir(packagePath), k)
		globPath, err := filepath.Glob(p)
		if err != nil {
			return err
		}
		for _, j := range globPath {
			mp[j] = true
		}
	}
	return nil
}
func getGlobPath(ctx context.Context, dir string, mp map[string]bool) error {
	basePath := ctx.Value("basePath").(string)
	if basePath != "" && filepath.Clean(filepath.Dir(basePath)) == filepath.Clean(dir) {
		return nil
	}
	packagePath := filepath.Join(dir, "package.json")
	_, err := os.Stat(packagePath)
	if os.IsNotExist(err) {
		return getGlobPath(ctx, filepath.Dir(dir), mp)
	}
	if err := getWorkspace(packagePath, mp); err != nil {
		return err
	}
	return getGlobPath(ctx, filepath.Dir(dir), mp)
}
func skipWorkspaceDirectory(ctx context.Context, dir string) bool {
	mp := make(map[string]bool)
	if err := getGlobPath(ctx, dir, mp); err != nil {
		return false
	}
	for k, _ := range mp {
		if filepath.Clean(k) == filepath.Clean(dir) {
			return false
		}
	}
	return true
}
