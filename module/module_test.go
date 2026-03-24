package module

import (
	"testing"

	"github.com/murphysecurity/murphysec/model"
)

func TestGetActiveInspectors(t *testing.T) {
	t.Parallel()

	old := SkillScanEnabled
	t.Cleanup(func() {
		SkillScanEnabled = old
	})

	SkillScanEnabled = true
	if !hasInspector(GetActiveInspectors(), "Skills") {
		t.Fatal("expected Skills inspector when skill scan is enabled")
	}

	SkillScanEnabled = false
	if hasInspector(GetActiveInspectors(), "Skills") {
		t.Fatal("did not expect Skills inspector when skill scan is disabled")
	}
}

func hasInspector(inspectors []model.Inspector, name string) bool {
	for _, inspector := range inspectors {
		if inspector.String() == name {
			return true
		}
	}
	return false
}
