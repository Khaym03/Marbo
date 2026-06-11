package domain

import "testing"

func TestLoadValidKnowledgeBase(
	t *testing.T,
) {

	kb, err := Load(
		"testdata/fixture_valid.json",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(kb.Intents) != 1 {
		t.Fatal("expected 1 intent")
	}
}
