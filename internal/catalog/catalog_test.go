package catalog

import "testing"

func TestLoadCatalogHelloWorld(t *testing.T) {
	def, err := Load("hello-world")
	if err != nil {
		t.Fatalf("expected catalog to load, got error: %v", err)
	}

	if len(def.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(def.Steps))
	}

	if def.Steps[0] != "01-initial" {
		t.Fatalf("unexpected first step: %s", def.Steps[0])
	}
}
