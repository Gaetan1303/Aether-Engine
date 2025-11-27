// internal/shared/domain/position/position_test.go

package position_test

import (
	"testing"

	// "aether-engine-server/internal/shared/domain/position" // Assurez-vous que le chemin est correct
)

func TestNewPosition3D_InvalidCoordinates(t *testing.T) {
	cases := []struct {
		name string
		x, y, z int
	}{
		{"Negative X", -1, 0, 0},
		{"Negative Y", 0, -2, 0},
		{"Negative Z", 0, 0, -5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := position.New(tc.x, tc.y, tc.z); err == nil {
				t.Fatalf("Expected error for %s, but got nil", tc.name)
			}
		})
	}
}

func TestPosition3D_Equals(t *testing.T) {
	a := position.NewUnchecked(1, 2, 0)
	b := position.NewUnchecked(1, 2, 0)
	c := position.NewUnchecked(2, 2, 0)
	d := position.NewUnchecked(1, 2, 1)

	if !a.Equals(b) {
		t.Fatal("Positions should be equal")
	}
	if a.Equals(c) || a.Equals(d) {
		t.Fatal("Positions should not be equal")
	}
}

func TestPosition3D_DistanceManhattan(t *testing.T) {
	a := position.NewUnchecked(0, 0, 0)
	// Déplacement : 3(x) + 4(y) + 1(z) = 8
	b := position.NewUnchecked(3, 4, 1) 
	
	d := a.DistanceManhattan(b)
	if d != 8 {
		t.Fatalf("Manhattan distance: Expected 8 got %d", d)
	}

	// Test inverse
	dInv := b.DistanceManhattan(a)
	if dInv != 8 {
		t.Fatalf("Manhattan distance (Inverse): Expected 8 got %d", dInv)
	}
}

func TestPosition3D_DistanceChebyshev(t *testing.T) {
	a := position.NewUnchecked(0, 0, 0)
	// max(3, 4) + 2(z) = 6
	b := position.NewUnchecked(3, 4, 2) 

	d := a.DistanceChebyshev(b)
	if d != 6 {
		t.Fatalf("Chebyshev distance: Expected 6 got %d", d)
	}
	
	// Test avec changement de Z plus grand que XY
	c := position.NewUnchecked(1, 1, 5) // max(1,1) + 5 = 6
	d2 := a.DistanceChebyshev(c)
	if d2 != 6 {
		t.Fatalf("Chebyshev distance 2: Expected 6 got %d", d2)
	}
}

func TestPosition3D_InBounds(t *testing.T) {
	p := position.NewUnchecked(2, 3, 1)
	
	// Cas 1 : Dans les limites (Width=10, Height=10, MaxZ=3)
	if !p.InBounds(10, 10, 3) {
		t.Fatal("Should be in bounds (2, 3, 1)")
	}

	// Cas 2 : X est hors limites (Width=2 signifie X valide: 0, 1)
	if p.InBounds(2, 4, 1) {
		t.Fatal("X >= Width should be out of bounds")
	}
	
	// Cas 3 : Z est hors limites (MaxZ=0 signifie Z valide: 0)
	if p.InBounds(10, 10, 0) {
		t.Fatal("Z > MaxZ should be out of bounds")
	}

	// Cas 4 : Coordonnée négative (même si InBounds vérifie > 0)
	pNeg := position.NewUnchecked(-1, 0, 0)
	if pNeg.InBounds(10, 10, 1) {
		t.Fatal("Negative coordinate should fail InBounds")
	}
}
