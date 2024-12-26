package shared

import (
	"testing"
)

func TestCircularDetector_Put(t *testing.T) {
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	if len(cd.path) != 1 {
		t.Fatalf("expected path length 1, got %d", len(cd.path))
	}
	if cd.path[0] != [2]string{"a", "1.0"} {
		t.Fatalf("expected path [a 1.0], got %v", cd.path[0])
	}
}

func TestCircularDetector_Put_CircularDependency(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic due to circular dependency")
		}
	}()
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	cd.Put("a", "1.0")
}

func TestCircularDetector_Leave(t *testing.T) {
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	cd.Leave("a", "1.0")
	if len(cd.path) != 0 {
		t.Fatalf("expected path length 0, got %d", len(cd.path))
	}
}

func TestCircularDetector_Leave_LengthUnderflow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic due to length underflow")
		}
	}()
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Leave("a", "1.0")
}

func TestCircularDetector_Leave_NameVersionInconsistent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic due to name@version inconsistent")
		}
	}()
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	cd.Leave("b", "1.0")
}

func TestCircularDetector_Contains(t *testing.T) {
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	if !cd.Contains("a", "1.0") {
		t.Fatalf("expected to contain a@1.0")
	}
	if cd.Contains("b", "1.0") {
		t.Fatalf("expected not to contain b@1.0")
	}
}

func TestCircularDetector_Error(t *testing.T) {
	cd := &CircularDetector{m: make(map[[2]string]struct{})}
	cd.Put("a", "1.0")
	err := cd.Error("a", "1.0")
	if err == nil {
		t.Fatalf("expected error due to circular dependency")
	}
	if err.Error() != "circle detected[length:1]: a@1.0 -> a@1.0" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}
