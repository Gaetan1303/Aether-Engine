package domain

import (
	"fmt"
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Test Manhattan Strategy - Chemin simple sans obstacles
func TestAStarManhattanStrategy_CheminSimple(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(3, 3)
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err, "Le pathfinding ne devrait pas retourner d'erreur")
	assert.NotNil(t, chemin, "Le chemin ne devrait pas être nil")
	assert.Equal(t, 6, cout, "Le coût devrait être 6 (3 droite + 3 haut)")
	assert.Equal(t, 6, len(chemin), "Le chemin devrait contenir 6 positions")

	// Vérifier que la dernière position est bien la destination
	derniere := chemin[len(chemin)-1]
	assert.Equal(t, arrivee.X(), derniere.X(), "X final devrait correspondre")
	assert.Equal(t, arrivee.Y(), derniere.Y(), "Y final devrait correspondre")
}

// Test 2: Test Euclidien Strategy - Même cas que Manhattan
func TestAStarEuclidienStrategy_CheminSimple(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	strategy := &AStarEuclidienStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(3, 3)
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	assert.Equal(t, 6, cout, "Le coût devrait être 6 avec distance euclidienne")
	assert.Equal(t, 6, len(chemin))
}

// Test 3: Test Diagonal Strategy - Doit trouver un chemin plus court avec diagonales
func TestAStarDiagonalStrategy_CheminDiagonal(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	strategy := &AStarDiagonalStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(3, 3)
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	// Avec diagonales: 3 déplacements diagonaux = coût de 3
	assert.Equal(t, 3, len(chemin), "Le chemin diagonal devrait contenir 3 positions")
	assert.Equal(t, 3, cout, "Le coût diagonal devrait être 3")
}

// Test 4: Test avec obstacles - Doit contourner un obstacle
func TestAStarManhattan_AvecObstacle(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(0, 4)

	// Créer un mur vertical en x=0, y=2
	pos, _ := shared.NewPosition(0, 2)
	grille.DefinirTypeCellule(pos, shared.CelluleObstacle)

	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	// Doit contourner l'obstacle, donc coût > distance directe
	assert.Greater(t, cout, 4, "Le coût devrait être supérieur à 4 car il faut contourner")

	// Vérifier qu'aucune position du chemin n'est sur l'obstacle
	for _, pos := range chemin {
		if pos.X() == 0 && pos.Y() == 2 {
			t.Errorf("Le chemin passe par l'obstacle en (0,2)")
		}
	}
}

// Test 5: Test avec unités occupant des positions - Doit les éviter
func TestAStarManhattan_AvecUnitesOccupees(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(2, 0)

	// Unité bloquante en (1,0)
	unitesOccupees := map[string]bool{
		"1,0": true,
	}

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	// Doit contourner l'unité, donc coût > 2
	assert.Greater(t, cout, 2, "Le coût devrait être supérieur à 2")

	// Vérifier qu'aucune position n'est occupée
	for _, pos := range chemin {
		key := fmt.Sprintf("%d,%d", pos.X(), pos.Y())
		assert.False(t, unitesOccupees[key], "Le chemin ne devrait pas passer par une position occupée")
	}
}

// Test 6: Test avec terrain difficile - Coût multiplié
func TestAStarManhattan_TerrainDifficile(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(2, 0)

	// Terrain difficile en (1,0) - coût x2
	pos, _ := shared.NewPosition(1, 0)
	grille.DefinirTypeCellule(pos, shared.CelluleDifficile)

	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	// Chemin direct: 2 cases avec coût doublé sur case 1 = 1 + 2 = 3
	// OU contournement: dépend de l'algo
	assert.GreaterOrEqual(t, cout, 3, "Le coût devrait être au moins 3")
}

// Test 7: Test aucun chemin possible - Destination entièrement bloquée
func TestAStarManhattan_AucunChemin(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(2, 2)

	// Entourer complètement la destination d'obstacles
	pos1, _ := shared.NewPosition(1, 2)
	grille.DefinirTypeCellule(pos1, shared.CelluleObstacle)
	pos2, _ := shared.NewPosition(3, 2)
	grille.DefinirTypeCellule(pos2, shared.CelluleObstacle)
	pos3, _ := shared.NewPosition(2, 1)
	grille.DefinirTypeCellule(pos3, shared.CelluleObstacle)
	pos4, _ := shared.NewPosition(2, 3)
	grille.DefinirTypeCellule(pos4, shared.CelluleObstacle)

	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.Error(t, err, "Devrait retourner une erreur si aucun chemin n'existe")
	assert.Nil(t, chemin)
	assert.Equal(t, 0, cout)
}

// Test 8: Test position de départ = position d'arrivée
func TestAStarManhattan_MemePosition(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	position, _ := shared.NewPosition(2, 2)
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, position, position, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, chemin, "Le chemin devrait être vide")
	assert.Equal(t, 0, cout, "Le coût devrait être 0")
}

// Test 9: Test hors limites - Position d'arrivée invalide
func TestAStarManhattan_HorsLimites(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(5, 5)
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(2, 2)
	arrivee, _ := shared.NewPosition(10, 10) // Hors grille
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.Error(t, err, "Devrait retourner une erreur pour position hors limites")
	assert.Nil(t, chemin)
	assert.Equal(t, 0, cout)
}

// Test 10: Test PathfindingFactory - Création de stratégies
func TestPathfindingFactory_CreationStrategies(t *testing.T) {
	// Arrange
	factory := &PathfindingFactory{}

	// Act & Assert - Manhattan
	manhattan := factory.CreatePathfinder("manhattan")
	assert.NotNil(t, manhattan)
	assert.Equal(t, "manhattan", manhattan.GetType())

	// Act & Assert - Euclidien
	euclidien := factory.CreatePathfinder("euclidien")
	assert.NotNil(t, euclidien)
	assert.Equal(t, "euclidien", euclidien.GetType())

	// Act & Assert - Diagonal
	diagonal := factory.CreatePathfinder("diagonal")
	assert.NotNil(t, diagonal)
	assert.Equal(t, "diagonal", diagonal.GetType())

	// Act & Assert - Default
	defaultStrategy := factory.CreateDefaultPathfinder()
	assert.NotNil(t, defaultStrategy)
	assert.Equal(t, "manhattan", defaultStrategy.GetType())

	// Act & Assert - Type invalide (devrait retourner Manhattan par défaut)
	invalide := factory.CreatePathfinder("invalid")
	assert.NotNil(t, invalide)
	assert.Equal(t, "manhattan", invalide.GetType())
}

// Test 11: Test PathfindingService - TrouverCheminAvecPortee
func TestPathfindingService_AvecPortee(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	service := NewPathfindingService()
	service.SetStrategyType("manhattan")

	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(5, 5) // Distance Manhattan = 10
	unitesOccupees := make(map[string]bool)
	porteeMax := 10

	// Act
	chemin, cout, err := service.TrouverCheminAvecPortee(grille, depart, arrivee, unitesOccupees, porteeMax)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	assert.LessOrEqual(t, cout, porteeMax, "Le coût ne devrait pas dépasser la portée max")
}

// Test 12: Test PathfindingService - TrouverCheminAvecPortee insuffisante
func TestPathfindingService_PorteeInsuffisante(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	service := NewPathfindingService()
	service.SetStrategyType("manhattan")

	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(5, 5) // Distance = 10
	unitesOccupees := make(map[string]bool)
	porteeMax := 5 // Portée insuffisante

	// Act
	chemin, cout, err := service.TrouverCheminAvecPortee(grille, depart, arrivee, unitesOccupees, porteeMax)

	// Assert
	require.Error(t, err, "Devrait retourner une erreur si portée insuffisante")
	assert.Contains(t, err.Error(), "HORS_PORTEE", "L'erreur devrait indiquer que c'est hors portée")
	assert.Nil(t, chemin)
	assert.Equal(t, 0, cout)
}

// Test 13: Test PathfindingService - TrouverPositionsAccessibles
func TestPathfindingService_PositionsAccessibles(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	service := NewPathfindingService()
	service.SetStrategyType("manhattan")

	depart, _ := shared.NewPosition(5, 5)
	unitesOccupees := make(map[string]bool)
	porteeMax := 2

	// Act
	positions := service.TrouverPositionsAccessibles(grille, depart, unitesOccupees, porteeMax)

	// Assert
	assert.NotNil(t, positions)
	assert.Greater(t, len(positions), 0, "Devrait trouver au moins quelques positions accessibles")

	// Toutes les positions devraient être à distance <= porteeMax
	for _, pos := range positions {
		distance := depart.Distance(pos)
		assert.LessOrEqual(t, distance, porteeMax, "Toutes les positions devraient être dans la portée")
	}
}

// Test 14: Test PathfindingService - EstAccessible
func TestPathfindingService_EstAccessible(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(10, 10)
	service := NewPathfindingService()
	service.SetStrategyType("manhattan")

	depart, _ := shared.NewPosition(0, 0)
	proche, _ := shared.NewPosition(2, 2) // Distance = 4
	loin, _ := shared.NewPosition(8, 8)   // Distance = 16
	unitesOccupees := make(map[string]bool)
	porteeMax := 5

	// Act & Assert - Position proche accessible
	accessible := service.EstAccessible(grille, depart, proche, unitesOccupees, porteeMax)
	assert.True(t, accessible, "La position proche devrait être accessible")

	// Act & Assert - Position loin inaccessible
	inaccessible := service.EstAccessible(grille, depart, loin, unitesOccupees, porteeMax)
	assert.False(t, inaccessible, "La position loin devrait être inaccessible")
}

// Test 15: Test de performance - Grande grille
func TestAStarManhattan_PerformanceGrandeGrille(t *testing.T) {
	// Arrange
	grille, _ := shared.NewGrilleCombat(50, 50) // Grande grille 50x50
	strategy := &AStarManhattanStrategy{}
	depart, _ := shared.NewPosition(0, 0)
	arrivee, _ := shared.NewPosition(49, 49)
	unitesOccupees := make(map[string]bool)

	// Act
	chemin, cout, err := strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chemin)
	assert.Equal(t, 98, cout, "Le coût devrait être 98 (49+49)")

	// Le test doit se terminer rapidement (pas de timeout)
	t.Log("Performance test passed - grande grille traitée efficacement")
}
