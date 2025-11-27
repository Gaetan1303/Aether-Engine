// internal/shared/domain/position/position_test.go

package position_test

import (
	"testing"

	"aether-engine-server/internal/shared/domain/position" 
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
    // Définir la taille du volume de la grille : 10x10 au sol, hauteur max 5.
    const GridW, GridH, GridMaxZ = 10, 10, 5 
	
    cases := []struct {
        name string
        x, y, z int
        expected bool
    }{
        // Cas normal
        {"Min corner", 0, 0, 0, true},
        {"Max corner", 9, 9, 5, true},
        {"Middle", 5, 5, 3, true},

        // Cas critique : Z trop haut par rapport au MAX du VOLUME
        {"Z too high", 5, 5, 6, false}, // Z > MaxZ

        // Cas X hors limites
        {"X too high", 10, 5, 2, false},
        {"X negative", -1, 5, 2, false},

        // Cas Y hors limites
        {"Y too high", 5, 10, 2, false},
        {"Y negative", 5, -1, 2, false},

        // Cas Z à la limite du volume
        {"Z at Max limit", 5, 5, 5, true}, // Z = MaxZ
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            p := position.NewUnchecked(tc.x, tc.y, tc.z)
            result := p.InBounds(GridW, GridH, GridMaxZ)
            
            if result != tc.expected {
                t.Errorf("InBounds(%d, %d, %d): Expected %t, got %t", 
                    tc.x, tc.y, tc.z, tc.expected, result)
            }
        })
    }
}

func TestPosition3D_IsAdjacent(t *testing.T) {
	a := position.NewUnchecked(2, 2, 2)
	b := position.NewUnchecked(3, 2, 2) // Adjacent X
	c := position.NewUnchecked(2, 3, 2) // Adjacent Y
	d := position.NewUnchecked(2, 2, 3) // Adjacent Z
	e := position.NewUnchecked(3, 3, 2) // Diagonal XY
	f := position.NewUnchecked(4, 4, 4) // Non adjacent

	if !a.IsAdjacent(b) || !a.IsAdjacent(c) || !a.IsAdjacent(d) || !a.IsAdjacent(e) {
		t.Fatal("Positions should be adjacent")
	}
	if a.IsAdjacent(f) {
		t.Fatal("Positions should not be adjacent")
	}
}
