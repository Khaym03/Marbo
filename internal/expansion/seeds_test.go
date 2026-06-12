package expansion

import (
	"testing"
)

func TestGetSeeds(t *testing.T) {
	seeds := GetSeeds()
	if len(seeds) == 0 {
		t.Fatal("expected seeds, got none")
	}

	found := false
	for _, s := range seeds {
		if s.IntentID == "INTENT_REQUISITOS_INGRESO" {
			found = true
			if len(s.Concepts) == 0 {
				t.Error("expected concepts for REQUISITOS_INGRESO")
			}
		}
	}
	if !found {
		t.Error("expected to find REQUISITOS_INGRESO in seeds")
	}
}
