package conan

import (
	"encoding/json"
	"testing"

	"github.com/murphysecurity/murphysec/model"
)

func TestConanGraphTreeSupportsDependenciesMap(t *testing.T) {
	j := []byte(`{
  "graph": {
    "nodes": {
      "0": {
        "id": "0",
        "ref": "conanfile",
        "label": "conanfile.txt",
        "dependencies": {
          "1": { "id": "1", "ref": "jsoncpp/1.9.4" },
          "2": { "id": "2", "ref": "zlib/1.2.11" }
        }
      },
      "1": { "id": "1", "ref": "jsoncpp/1.9.4", "dependencies": { "3": { "id": "3", "ref": "fmt/10.0.0" } } },
      "2": { "id": "2", "ref": "zlib/1.2.11", "dependencies": {} }
      ,"3": { "id": "3", "ref": "fmt/10.0.0", "dependencies": {} }
    }
  }
}`)

	var g _ConanGraphInfoJsonFile
	if err := json.Unmarshal(j, &g); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	tree, err := g.Tree()
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}
	if tree == nil {
		t.Fatal("tree is nil")
	}
	if len(tree.Dependencies) != 2 {
		t.Fatalf("want 2 dependencies, got %d", len(tree.Dependencies))
	}
	got := map[string]string{}
	for _, dep := range tree.Dependencies {
		got[dep.CompName] = dep.CompVersion
	}
	if got["jsoncpp"] != "1.9.4" {
		t.Fatalf("jsoncpp parse mismatch, got %q", got["jsoncpp"])
	}
	if got["zlib"] != "1.2.11" {
		t.Fatalf("zlib parse mismatch, got %q", got["zlib"])
	}
	for _, dep := range tree.Dependencies {
		if dep.DependencyRelation != model.DependencyRelationDirect {
			t.Fatalf("expected direct dependency relation, got %v for %s", dep.DependencyRelation, dep.CompName)
		}
		if dep.CompName == "jsoncpp" {
			if len(dep.Dependencies) != 1 {
				t.Fatalf("expected one transitive dependency under jsoncpp, got %d", len(dep.Dependencies))
			}
			if dep.Dependencies[0].CompName != "fmt" || dep.Dependencies[0].CompVersion != "10.0.0" {
				t.Fatalf("unexpected transitive dependency: %s/%s", dep.Dependencies[0].CompName, dep.Dependencies[0].CompVersion)
			}
			if dep.Dependencies[0].DependencyRelation != model.DependencyRelationTransitive {
				t.Fatalf("expected transitive dependency relation, got %v", dep.Dependencies[0].DependencyRelation)
			}
		}
	}
}

func TestConanGraphTreeSupportsLegacyRequires(t *testing.T) {
	j := []byte(`{
  "graph": {
    "nodes": [
      { "id": "0", "ref": "conanfile", "requires": [{ "id": "1", "ref": "openssl/3.0.0" }] },
      { "id": "1", "ref": "openssl/3.0.0", "requires": [] }
    ]
  }
}`)
	var g _ConanGraphInfoJsonFile
	if err := json.Unmarshal(j, &g); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	tree, err := g.Tree()
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}
	if tree == nil || len(tree.Dependencies) != 1 {
		t.Fatalf("want 1 dependency, got %#v", tree)
	}
	if tree.Dependencies[0].CompName != "openssl" || tree.Dependencies[0].CompVersion != "3.0.0" {
		t.Fatalf("openssl parse mismatch: %s/%s", tree.Dependencies[0].CompName, tree.Dependencies[0].CompVersion)
	}
	if tree.Dependencies[0].DependencyRelation != model.DependencyRelationDirect {
		t.Fatalf("expected direct dependency relation, got %v", tree.Dependencies[0].DependencyRelation)
	}
}
