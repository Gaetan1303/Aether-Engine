package domain

import (
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// PathfindingFactory crée des stratégies de pathfinding
// Factory Pattern - encapsule la création d'objets complexes
// Single Responsibility Principle - responsable uniquement de la création
type PathfindingFactory struct{}

// NewPathfindingFactory crée une nouvelle factory de pathfinding
func NewPathfindingFactory() *PathfindingFactory {
	return &PathfindingFactory{}
}

// CreatePathfinder crée une stratégie de pathfinding selon le type
// Open/Closed Principle - fermé à la modification, ouvert à l'extension
func (f *PathfindingFactory) CreatePathfinder(strategyType string) PathfindingStrategy {
	switch strategyType {
	case "manhattan":
		return NewAStarManhattanStrategy()
	case "euclidien":
		return NewAStarEuclidienStrategy()
	case "diagonal":
		return NewAStarDiagonalStrategy()
	default:
		// Par défaut, utiliser Manhattan (le plus utilisé dans les jeux tactiques)
		return NewAStarManhattanStrategy()
	}
}

// CreateDefaultPathfinder crée la stratégie par défaut (Manhattan)
func (f *PathfindingFactory) CreateDefaultPathfinder() PathfindingStrategy {
	return NewAStarManhattanStrategy()
}

// CreatePathfinderForTerrain choisit la stratégie selon le type de terrain
// Démontre le Strategy Pattern - sélection contextuelle de l'algorithme
func (f *PathfindingFactory) CreatePathfinderForTerrain(typeTerrain string) PathfindingStrategy {
	switch typeTerrain {
	case "grille":
		// Terrain en grille stricte : Manhattan (4 directions)
		return NewAStarManhattanStrategy()
	case "ouvert":
		// Terrain ouvert : Euclidien (plus naturel)
		return NewAStarEuclidienStrategy()
	case "flexible":
		// Terrain flexible : Diagonal (8 directions)
		return NewAStarDiagonalStrategy()
	default:
		return NewAStarManhattanStrategy()
	}
}

// CreatePathfinderForUnit crée une stratégie selon les capacités de l'unité
// Démontre la flexibilité du Strategy Pattern
func (f *PathfindingFactory) CreatePathfinderForUnit(unite *Unite) PathfindingStrategy {
	// Exemple : certaines unités peuvent se déplacer en diagonale
	// (à adapter selon les règles métier)

	// Pour l'instant, toutes les unités utilisent Manhattan
	// Ceci peut être étendu avec des capacités spéciales
	return NewAStarManhattanStrategy()
}

// PathfindingService encapsule la logique de pathfinding
// Facade Pattern - simplifie l'utilisation du pathfinding
// Single Responsibility Principle - gère uniquement le pathfinding de haut niveau
type PathfindingService struct {
	factory  *PathfindingFactory
	strategy PathfindingStrategy
}

// NewPathfindingService crée un nouveau service de pathfinding
func NewPathfindingService() *PathfindingService {
	factory := NewPathfindingFactory()
	return &PathfindingService{
		factory:  factory,
		strategy: factory.CreateDefaultPathfinder(),
	}
}

// SetStrategy change la stratégie de pathfinding
// Strategy Pattern - permet de changer d'algorithme à l'exécution
func (s *PathfindingService) SetStrategy(strategy PathfindingStrategy) {
	s.strategy = strategy
}

// SetStrategyType change la stratégie par son type
func (s *PathfindingService) SetStrategyType(strategyType string) {
	s.strategy = s.factory.CreatePathfinder(strategyType)
}

// TrouverChemin trouve un chemin avec la stratégie actuelle
func (s *PathfindingService) TrouverChemin(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
) ([]*shared.Position, int, error) {
	return s.strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)
}

// TrouverCheminAvecPortee trouve un chemin en respectant la portée de mouvement
// Ajoute une validation métier sur le pathfinding
func (s *PathfindingService) TrouverCheminAvecPortee(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
	porteeMax int,
) ([]*shared.Position, int, error) {
	// Trouver le chemin optimal
	chemin, cout, err := s.strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)
	if err != nil {
		return nil, -1, err
	}

	// Vérifier que le coût ne dépasse pas la portée
	if cout > porteeMax {
		return nil, -1, shared.NewDomainError("destination hors de portée", "HORS_PORTEE")
	}

	return chemin, cout, nil
}

// TrouverPositionsAccessibles trouve toutes les positions accessibles dans une portée
// Utile pour afficher les cases de déplacement possibles
func (s *PathfindingService) TrouverPositionsAccessibles(
	grille *shared.GrilleCombat,
	depart *shared.Position,
	unitesOccupees map[string]bool,
	porteeMax int,
) []*shared.Position {
	positions := make([]*shared.Position, 0)

	// Parcourir toute la grille
	for y := 0; y < grille.Hauteur(); y++ {
		for x := 0; x < grille.Largeur(); x++ {
			pos, err := shared.NewPosition(x, y)
			if err != nil || !grille.EstTraversable(pos) {
				continue
			}

			// Ignorer la position de départ
			if pos.Equals(depart) {
				continue
			}

			// Vérifier si la position est occupée
			cle := positionKey(pos)
			if unitesOccupees[cle] {
				continue
			}

			// Essayer de trouver un chemin
			_, cout, err := s.strategy.TrouverChemin(grille, depart, pos, unitesOccupees)
			if err == nil && cout <= porteeMax {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

// EstAccessible vérifie si une position est accessible depuis une autre
func (s *PathfindingService) EstAccessible(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
	porteeMax int,
) bool {
	_, cout, err := s.strategy.TrouverChemin(grille, depart, arrivee, unitesOccupees)
	return err == nil && cout <= porteeMax
}

// GetStrategyType retourne le type de stratégie actuelle
func (s *PathfindingService) GetStrategyType() string {
	return s.strategy.GetType()
}
