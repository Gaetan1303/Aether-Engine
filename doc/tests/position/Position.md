# Test unitaire – Position

## Objectif

Vérifier la validité, l'égalité, les calculs de distance et la gestion des coordonnées d'une position sur la grille tactique.

Règles métier
### Contraintes

Les coordonnées X et Y doivent être ≥ 0
Les coordonnées doivent être dans les limites de la grille (configurables)
Position est un Value Object : immutable après création

### Comportements attendus

Création : Validation à la construction
Égalité : Comparaison par valeur
Distance Manhattan : |x1-x2| + |y1-y2| (déplacement sur grille)
Distance Euclidienne : √((x1-x2)² + (y1-y2)²) (portée compétences)
Adjacence : Vérifier si deux positions sont voisines


#### Structure proposée : 
```go
type Position struct {
package domain

import "errors"

// Position représente une coordonnée sur la grille de combat
type Position struct {
    x int
    y int
}

// NewPosition crée une position valide
func NewPosition(x, y int) (Position, error) {
    if x < 0 || y < 0 {
        return Position{}, errors.New("coordinates must be non-negative")
    }
    return Position{x: x, y: y}, nil
}

// Getters (encapsulation)
func (p Position) X() int { return p.x }
func (p Position) Y() int { return p.y }

// Equals compare deux positions
func (p Position) Equals(other Position) bool {
    return p.x == other.x && p.y == other.y
}

// ManhattanDistance calcule la distance Manhattan
func (p Position) ManhattanDistance(other Position) int {
    return abs(p.x-other.x) + abs(p.y-other.y)
}

// EuclideanDistance calcule la distance euclidienne
func (p Position) EuclideanDistance(other Position) float64 {
    dx := float64(p.x - other.x)
    dy := float64(p.y - other.y)
    return math.Sqrt(dx*dx + dy*dy)
}

// IsAdjacent vérifie si deux positions sont adjacentes (distance Manhattan = 1)
func (p Position) IsAdjacent(other Position) bool {
    return p.ManhattanDistance(other) == 1
}

// IsInBounds vérifie si la position est dans les limites
func (p Position) IsInBounds(width, height int) bool {
    return p.x < width && p.y < height
}

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
```

### Test unitaire (Go + testify)

```gopackage domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "aether-engine-server/internal/combat/domain"
)

// ========== Tests de création ==========

func TestNewPosition_ValidCoordinates(t *testing.T) {
    p, err := domain.NewPosition(2, 3)
    
    assert.NoError(t, err)
    assert.Equal(t, 2, p.X())
    assert.Equal(t, 3, p.Y())
}

func TestNewPosition_NegativeX(t *testing.T) {
    _, err := domain.NewPosition(-1, 3)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "non-negative")
}

func TestNewPosition_NegativeY(t *testing.T) {
    _, err := domain.NewPosition(2, -5)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "non-negative")
}

func TestNewPosition_Origin(t *testing.T) {
    p, err := domain.NewPosition(0, 0)
    
    assert.NoError(t, err)
    assert.Equal(t, 0, p.X())
    assert.Equal(t, 0, p.Y())
}

// ========== Tests d'égalité ==========

func TestPosition_Equals_SameCoordinates(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 3)
    
    assert.True(t, p1.Equals(p2))
}

func TestPosition_Equals_DifferentX(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(5, 3)
    
    assert.False(t, p1.Equals(p2))
}

func TestPosition_Equals_DifferentY(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 7)
    
    assert.False(t, p1.Equals(p2))
}

// ========== Tests de distance Manhattan ==========

func TestPosition_ManhattanDistance_SamePosition(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 3)
    
    distance := p1.ManhattanDistance(p2)
    
    assert.Equal(t, 0, distance)
}

func TestPosition_ManhattanDistance_Horizontal(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(5, 3)
    
    distance := p1.ManhattanDistance(p2)
    
    assert.Equal(t, 3, distance)
}

func TestPosition_ManhattanDistance_Vertical(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 7)
    
    distance := p1.ManhattanDistance(p2)
    
    assert.Equal(t, 4, distance)
}

func TestPosition_ManhattanDistance_Diagonal(t *testing.T) {
    p1, _ := domain.NewPosition(1, 1)
    p2, _ := domain.NewPosition(4, 5)
    
    distance := p1.ManhattanDistance(p2)
    
    // |4-1| + |5-1| = 3 + 4 = 7
    assert.Equal(t, 7, distance)
}

// ========== Tests de distance Euclidienne ==========

func TestPosition_EuclideanDistance_SamePosition(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 3)
    
    distance := p1.EuclideanDistance(p2)
    
    assert.Equal(t, 0.0, distance)
}

func TestPosition_EuclideanDistance_Pythagorean(t *testing.T) {
    p1, _ := domain.NewPosition(0, 0)
    p2, _ := domain.NewPosition(3, 4)
    
    distance := p1.EuclideanDistance(p2)
    
    // √(3² + 4²) = √25 = 5
    assert.Equal(t, 5.0, distance)
}

func TestPosition_EuclideanDistance_Diagonal(t *testing.T) {
    p1, _ := domain.NewPosition(1, 1)
    p2, _ := domain.NewPosition(2, 2)
    
    distance := p1.EuclideanDistance(p2)
    
    // √(1² + 1²) = √2 ≈ 1.414
    assert.InDelta(t, 1.414, distance, 0.001)
}

// ========== Tests d'adjacence ==========

func TestPosition_IsAdjacent_Right(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(3, 3)
    
    assert.True(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_Left(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(1, 3)
    
    assert.True(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_Up(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 4)
    
    assert.True(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_Down(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 2)
    
    assert.True(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_Diagonal(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(3, 4)
    
    // Distance Manhattan = 2, donc pas adjacent
    assert.False(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_SamePosition(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(2, 3)
    
    assert.False(t, p1.IsAdjacent(p2))
}

func TestPosition_IsAdjacent_TooFar(t *testing.T) {
    p1, _ := domain.NewPosition(2, 3)
    p2, _ := domain.NewPosition(5, 7)
    
    assert.False(t, p1.IsAdjacent(p2))
}

// ========== Tests de limites de grille ==========

func TestPosition_IsInBounds_ValidPosition(t *testing.T) {
    p, _ := domain.NewPosition(10, 15)
    
    assert.True(t, p.IsInBounds(48, 48))
}

func TestPosition_IsInBounds_EdgeX(t *testing.T) {
    p, _ := domain.NewPosition(47, 20)
    
    assert.True(t, p.IsInBounds(48, 48))
}

func TestPosition_IsInBounds_EdgeY(t *testing.T) {
    p, _ := domain.NewPosition(20, 47)
    
    assert.True(t, p.IsInBounds(48, 48))
}

func TestPosition_IsInBounds_OutOfBoundsX(t *testing.T) {
    p, _ := domain.NewPosition(48, 20)
    
    assert.False(t, p.IsInBounds(48, 48))
}

func TestPosition_IsInBounds_OutOfBoundsY(t *testing.T) {
    p, _ := domain.NewPosition(20, 48)
    
    assert.False(t, p.IsInBounds(48, 48))
}

func TestPosition_IsInBounds_Origin(t *testing.T) {
    p, _ := domain.NewPosition(0, 0)
    
    assert.True(t, p.IsInBounds(48, 48))
}

// ========== Tests de cas limites ==========

func TestPosition_MaxCoordinates(t *testing.T) {
    p, err := domain.NewPosition(9999, 9999)
    
    assert.NoError(t, err)
    assert.Equal(t, 9999, p.X())
    assert.Equal(t, 9999, p.Y())
}

func TestPosition_Immutability(t *testing.T) {
    p1, _ := domain.NewPosition(5, 5)
    x := p1.X()
    y := p1.Y()
    
    // Tenter de modifier (ne compile pas grâce aux champs privés)
    // p1.x = 10 // ❌ Erreur de compilation
    
    // Vérifier que les valeurs n'ont pas changé
    assert.Equal(t, x, p1.X())
    assert.Equal(t, y, p1.Y())
}
```
