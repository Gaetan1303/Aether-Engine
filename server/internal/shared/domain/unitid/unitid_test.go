package unitid_test

import (
	"testing"

	uid "aether-engine-server/internal/shared/domain/unitid"
)

func TestNewUnitID_Invalid(t *testing.T) {
	cases := []struct {
		name string
		id   int
	}{
		{"Zero", 0},
		{"Négatif", -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := uid.New(tc.id); err == nil {
				t.Fatalf("Attendu une erreur pour %s, reçu nil", tc.name)
			}
		})
	}
}

func TestNewUnitID_Valid(t *testing.T) {
	id, err := uid.New(42)
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}
	if id.Value() != 42 {
		t.Fatalf("Value attendue 42, reçu %d", id.Value())
	}
}

func TestUnitID_Equals(t *testing.T) {
	a := uid.NewUnchecked(7)
	b := uid.NewUnchecked(7)
	c := uid.NewUnchecked(8)

	if !a.Equals(b) {
		t.Fatal("UnitID identiques devraient être égaux")
	}
	if a.Equals(c) {
		t.Fatal("UnitID différents ne devraient pas être égaux")
	}
}
