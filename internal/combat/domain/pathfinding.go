package domain

import (
	"container/heap"
	"errors"
	"fmt"
	"math"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// PathfindingStrategy définit l'interface pour les algorithmes de pathfinding
// Strategy Pattern - permet de changer l'algorithme de pathfinding dynamiquement
// Interface Segregation Principle - interface minimale et cohérente
type PathfindingStrategy interface {
	// TrouverChemin trouve le plus court chemin entre deux positions
	// Retourne le chemin (sans la position de départ) et le coût total
	TrouverChemin(grille *shared.GrilleCombat, depart, arrivee *shared.Position, uniteOccupees map[string]bool) ([]*shared.Position, int, error)

	// GetType retourne le type de stratégie
	GetType() string
}

// Noeud représente un nœud dans l'algorithme A*
// Value Object - immuable et sans identité
type Noeud struct {
	position *shared.Position
	parent   *Noeud
	gCost    int // Coût depuis le départ
	hCost    int // Heuristique vers l'arrivée
	fCost    int // gCost + hCost (coût total estimé)
}

// NewNoeud crée un nouveau nœud
func NewNoeud(position *shared.Position, parent *Noeud, gCost, hCost int) *Noeud {
	return &Noeud{
		position: position,
		parent:   parent,
		gCost:    gCost,
		hCost:    hCost,
		fCost:    gCost + hCost,
	}
}

// Position retourne la position du nœud
func (n *Noeud) Position() *shared.Position { return n.position }

// Parent retourne le nœud parent
func (n *Noeud) Parent() *Noeud { return n.parent }

// GCost retourne le coût depuis le départ
func (n *Noeud) GCost() int { return n.gCost }

// FCost retourne le coût total estimé
func (n *Noeud) FCost() int { return n.fCost }

// PriorityQueue implémente une file de priorité pour A*
// Utilise le package container/heap de Go
type PriorityQueue []*Noeud

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Priorité au fCost le plus faible
	// En cas d'égalité, priorité au hCost le plus faible
	if pq[i].fCost == pq[j].fCost {
		return pq[i].hCost < pq[j].hCost
	}
	return pq[i].fCost < pq[j].fCost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Noeud))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// AStarManhattanStrategy implémente A* avec heuristique Manhattan (4 directions)
// Single Responsibility Principle - une seule responsabilité : pathfinding Manhattan
type AStarManhattanStrategy struct{}

// NewAStarManhattanStrategy crée une nouvelle stratégie Manhattan
func NewAStarManhattanStrategy() *AStarManhattanStrategy {
	return &AStarManhattanStrategy{}
}

// GetType retourne le type de stratégie
func (s *AStarManhattanStrategy) GetType() string {
	return "AStar-Manhattan"
}

// TrouverChemin implémente l'algorithme A* avec heuristique Manhattan
func (s *AStarManhattanStrategy) TrouverChemin(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
) ([]*shared.Position, int, error) {
	// Validation des entrées
	if !grille.EstDansLimites(depart) {
		return nil, -1, errors.New("position de départ hors limites")
	}
	if !grille.EstDansLimites(arrivee) {
		return nil, -1, errors.New("position d'arrivée hors limites")
	}
	if !grille.EstTraversable(arrivee) {
		return nil, -1, errors.New("position d'arrivée non traversable")
	}

	// Vérifier que l'arrivée n'est pas occupée par une unité
	cle := positionKey(arrivee)
	if unitesOccupees[cle] {
		return nil, -1, errors.New("position d'arrivée occupée par une unité")
	}

	// Si départ == arrivée
	if depart.Equals(arrivee) {
		return []*shared.Position{}, 0, nil
	}

	// Initialisation
	openSet := &PriorityQueue{}
	heap.Init(openSet)

	closedSet := make(map[string]bool)
	gScores := make(map[string]int)

	// Nœud de départ
	heuristique := s.heuristique(depart, arrivee)
	noeudDepart := NewNoeud(depart, nil, 0, heuristique)
	heap.Push(openSet, noeudDepart)
	gScores[positionKey(depart)] = 0

	// Boucle principale A*
	for openSet.Len() > 0 {
		// Extraire le nœud avec le plus petit fCost
		current := heap.Pop(openSet).(*Noeud)
		currentKey := positionKey(current.Position())

		// Si on a atteint l'arrivée
		if current.Position().Equals(arrivee) {
			return s.reconstruireChemin(current), current.GCost(), nil
		}

		// Marquer comme visité
		closedSet[currentKey] = true

		// Explorer les voisins (4 directions : haut, droite, bas, gauche)
		voisins := s.obtenirVoisins(grille, current.Position(), unitesOccupees)

		for _, voisinPos := range voisins {
			voisinKey := positionKey(voisinPos)

			// Ignorer si déjà visité
			if closedSet[voisinKey] {
				continue
			}

			// Calculer le coût pour atteindre ce voisin
			coutDeplacement := grille.CoutDeplacement(voisinPos)
			nouveauGCost := current.GCost() + coutDeplacement

			// Vérifier si ce chemin est meilleur
			ancienGCost, existe := gScores[voisinKey]
			if !existe || nouveauGCost < ancienGCost {
				// Mettre à jour le score
				gScores[voisinKey] = nouveauGCost

				// Créer le nouveau nœud
				heuristique := s.heuristique(voisinPos, arrivee)
				noeudVoisin := NewNoeud(voisinPos, current, nouveauGCost, heuristique)
				heap.Push(openSet, noeudVoisin)
			}
		}
	}

	// Aucun chemin trouvé
	return nil, -1, errors.New("aucun chemin trouvé")
}

// heuristique calcule la distance de Manhattan
func (s *AStarManhattanStrategy) heuristique(pos1, pos2 *shared.Position) int {
	return pos1.Distance(pos2) // Distance de Manhattan déjà implémentée
}

// obtenirVoisins retourne les positions adjacentes traversables (4 directions)
func (s *AStarManhattanStrategy) obtenirVoisins(
	grille *shared.GrilleCombat,
	pos *shared.Position,
	unitesOccupees map[string]bool,
) []*shared.Position {
	voisins := make([]*shared.Position, 0, 4)

	// 4 directions : haut, droite, bas, gauche
	directions := []struct{ dx, dy int }{
		{0, -1}, // Haut
		{1, 0},  // Droite
		{0, 1},  // Bas
		{-1, 0}, // Gauche
	}

	for _, dir := range directions {
		voisin, _ := shared.NewPosition(pos.X()+dir.dx, pos.Y()+dir.dy)
		if voisin == nil {
			continue
		}

		// Vérifier si traversable et non occupé
		if grille.EstTraversable(voisin) {
			cle := positionKey(voisin)
			if !unitesOccupees[cle] {
				voisins = append(voisins, voisin)
			}
		}
	}

	return voisins
}

// reconstruireChemin reconstruit le chemin depuis l'arrivée jusqu'au départ
func (s *AStarManhattanStrategy) reconstruireChemin(noeudArrivee *Noeud) []*shared.Position {
	chemin := make([]*shared.Position, 0)
	current := noeudArrivee

	// Remonter depuis l'arrivée jusqu'au départ
	for current.Parent() != nil {
		chemin = append(chemin, current.Position())
		current = current.Parent()
	}

	// Inverser le chemin (pour avoir départ → arrivée)
	for i, j := 0, len(chemin)-1; i < j; i, j = i+1, j-1 {
		chemin[i], chemin[j] = chemin[j], chemin[i]
	}

	return chemin
}

// AStarEuclidienStrategy implémente A* avec heuristique Euclidienne
// Permet un pathfinding plus "naturel" avec diagonales
// Single Responsibility Principle - une seule responsabilité : pathfinding Euclidien
type AStarEuclidienStrategy struct{}

// NewAStarEuclidienStrategy crée une nouvelle stratégie Euclidienne
func NewAStarEuclidienStrategy() *AStarEuclidienStrategy {
	return &AStarEuclidienStrategy{}
}

// GetType retourne le type de stratégie
func (s *AStarEuclidienStrategy) GetType() string {
	return "AStar-Euclidien"
}

// TrouverChemin implémente l'algorithme A* avec heuristique Euclidienne
func (s *AStarEuclidienStrategy) TrouverChemin(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
) ([]*shared.Position, int, error) {
	// Validation des entrées
	if !grille.EstDansLimites(depart) {
		return nil, -1, errors.New("position de départ hors limites")
	}
	if !grille.EstDansLimites(arrivee) {
		return nil, -1, errors.New("position d'arrivée hors limites")
	}
	if !grille.EstTraversable(arrivee) {
		return nil, -1, errors.New("position d'arrivée non traversable")
	}

	// Vérifier que l'arrivée n'est pas occupée
	cle := positionKey(arrivee)
	if unitesOccupees[cle] {
		return nil, -1, errors.New("position d'arrivée occupée par une unité")
	}

	// Si départ == arrivée
	if depart.Equals(arrivee) {
		return []*shared.Position{}, 0, nil
	}

	// Initialisation (identique à Manhattan)
	openSet := &PriorityQueue{}
	heap.Init(openSet)

	closedSet := make(map[string]bool)
	gScores := make(map[string]int)

	// Nœud de départ avec heuristique Euclidienne
	heuristique := s.heuristique(depart, arrivee)
	noeudDepart := NewNoeud(depart, nil, 0, heuristique)
	heap.Push(openSet, noeudDepart)
	gScores[positionKey(depart)] = 0

	// Boucle principale A*
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Noeud)
		currentKey := positionKey(current.Position())

		if current.Position().Equals(arrivee) {
			return s.reconstruireChemin(current), current.GCost(), nil
		}

		closedSet[currentKey] = true

		// Explorer les voisins (4 directions pour cohérence avec Manhattan)
		voisins := s.obtenirVoisins(grille, current.Position(), unitesOccupees)

		for _, voisinPos := range voisins {
			voisinKey := positionKey(voisinPos)

			if closedSet[voisinKey] {
				continue
			}

			coutDeplacement := grille.CoutDeplacement(voisinPos)
			nouveauGCost := current.GCost() + coutDeplacement

			ancienGCost, existe := gScores[voisinKey]
			if !existe || nouveauGCost < ancienGCost {
				gScores[voisinKey] = nouveauGCost

				heuristique := s.heuristique(voisinPos, arrivee)
				noeudVoisin := NewNoeud(voisinPos, current, nouveauGCost, heuristique)
				heap.Push(openSet, noeudVoisin)
			}
		}
	}

	return nil, -1, errors.New("aucun chemin trouvé")
}

// heuristique calcule la distance Euclidienne (arrondie à l'entier)
func (s *AStarEuclidienStrategy) heuristique(pos1, pos2 *shared.Position) int {
	distanceEuclidienne := pos1.DistanceEuclidienne(pos2)
	return int(math.Round(distanceEuclidienne))
}

// obtenirVoisins retourne les positions adjacentes (4 directions)
func (s *AStarEuclidienStrategy) obtenirVoisins(
	grille *shared.GrilleCombat,
	pos *shared.Position,
	unitesOccupees map[string]bool,
) []*shared.Position {
	voisins := make([]*shared.Position, 0, 4)

	directions := []struct{ dx, dy int }{
		{0, -1}, {1, 0}, {0, 1}, {-1, 0},
	}

	for _, dir := range directions {
		voisin, _ := shared.NewPosition(pos.X()+dir.dx, pos.Y()+dir.dy)
		if voisin == nil {
			continue
		}

		if grille.EstTraversable(voisin) {
			cle := positionKey(voisin)
			if !unitesOccupees[cle] {
				voisins = append(voisins, voisin)
			}
		}
	}

	return voisins
}

// reconstruireChemin reconstruit le chemin
func (s *AStarEuclidienStrategy) reconstruireChemin(noeudArrivee *Noeud) []*shared.Position {
	chemin := make([]*shared.Position, 0)
	current := noeudArrivee

	for current.Parent() != nil {
		chemin = append(chemin, current.Position())
		current = current.Parent()
	}

	for i, j := 0, len(chemin)-1; i < j; i, j = i+1, j-1 {
		chemin[i], chemin[j] = chemin[j], chemin[i]
	}

	return chemin
}

// AStarDiagonalStrategy implémente A* avec déplacements en diagonale (8 directions)
// Single Responsibility Principle - une seule responsabilité : pathfinding diagonal
type AStarDiagonalStrategy struct{}

// NewAStarDiagonalStrategy crée une nouvelle stratégie Diagonale
func NewAStarDiagonalStrategy() *AStarDiagonalStrategy {
	return &AStarDiagonalStrategy{}
}

// GetType retourne le type de stratégie
func (s *AStarDiagonalStrategy) GetType() string {
	return "AStar-Diagonal"
}

// TrouverChemin implémente l'algorithme A* avec déplacements diagonaux
func (s *AStarDiagonalStrategy) TrouverChemin(
	grille *shared.GrilleCombat,
	depart, arrivee *shared.Position,
	unitesOccupees map[string]bool,
) ([]*shared.Position, int, error) {
	// Validation
	if !grille.EstDansLimites(depart) {
		return nil, -1, errors.New("position de départ hors limites")
	}
	if !grille.EstDansLimites(arrivee) {
		return nil, -1, errors.New("position d'arrivée hors limites")
	}
	if !grille.EstTraversable(arrivee) {
		return nil, -1, errors.New("position d'arrivée non traversable")
	}

	cle := positionKey(arrivee)
	if unitesOccupees[cle] {
		return nil, -1, errors.New("position d'arrivée occupée par une unité")
	}

	if depart.Equals(arrivee) {
		return []*shared.Position{}, 0, nil
	}

	// Initialisation
	openSet := &PriorityQueue{}
	heap.Init(openSet)

	closedSet := make(map[string]bool)
	gScores := make(map[string]int)

	heuristique := s.heuristique(depart, arrivee)
	noeudDepart := NewNoeud(depart, nil, 0, heuristique)
	heap.Push(openSet, noeudDepart)
	gScores[positionKey(depart)] = 0

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Noeud)
		currentKey := positionKey(current.Position())

		if current.Position().Equals(arrivee) {
			return s.reconstruireChemin(current), current.GCost(), nil
		}

		closedSet[currentKey] = true

		// Explorer les voisins (8 directions avec diagonales)
		voisins := s.obtenirVoisins(grille, current.Position(), unitesOccupees)

		for _, voisinPos := range voisins {
			voisinKey := positionKey(voisinPos)

			if closedSet[voisinKey] {
				continue
			}

			// Coût diagonal = coût base * sqrt(2) ≈ 1.4 (arrondi à 1 pour simplicité)
			// On pourrait utiliser des coûts fixes : diagonal = 14, orthogonal = 10
			coutDeplacement := grille.CoutDeplacement(voisinPos)
			nouveauGCost := current.GCost() + coutDeplacement

			ancienGCost, existe := gScores[voisinKey]
			if !existe || nouveauGCost < ancienGCost {
				gScores[voisinKey] = nouveauGCost

				heuristique := s.heuristique(voisinPos, arrivee)
				noeudVoisin := NewNoeud(voisinPos, current, nouveauGCost, heuristique)
				heap.Push(openSet, noeudVoisin)
			}
		}
	}

	return nil, -1, errors.New("aucun chemin trouvé")
}

// heuristique utilise la distance de Chebyshev (adapté aux diagonales)
func (s *AStarDiagonalStrategy) heuristique(pos1, pos2 *shared.Position) int {
	dx := abs(pos1.X() - pos2.X())
	dy := abs(pos1.Y() - pos2.Y())
	// Distance de Chebyshev : max(dx, dy)
	if dx > dy {
		return dx
	}
	return dy
}

// obtenirVoisins retourne les positions adjacentes (8 directions avec diagonales)
func (s *AStarDiagonalStrategy) obtenirVoisins(
	grille *shared.GrilleCombat,
	pos *shared.Position,
	unitesOccupees map[string]bool,
) []*shared.Position {
	voisins := make([]*shared.Position, 0, 8)

	// 8 directions : N, NE, E, SE, S, SW, W, NW
	directions := []struct{ dx, dy int }{
		{0, -1},  // N
		{1, -1},  // NE
		{1, 0},   // E
		{1, 1},   // SE
		{0, 1},   // S
		{-1, 1},  // SW
		{-1, 0},  // W
		{-1, -1}, // NW
	}

	for _, dir := range directions {
		voisin, _ := shared.NewPosition(pos.X()+dir.dx, pos.Y()+dir.dy)
		if voisin == nil {
			continue
		}

		if grille.EstTraversable(voisin) {
			cle := positionKey(voisin)
			if !unitesOccupees[cle] {
				voisins = append(voisins, voisin)
			}
		}
	}

	return voisins
}

// reconstruireChemin reconstruit le chemin
func (s *AStarDiagonalStrategy) reconstruireChemin(noeudArrivee *Noeud) []*shared.Position {
	chemin := make([]*shared.Position, 0)
	current := noeudArrivee

	for current.Parent() != nil {
		chemin = append(chemin, current.Position())
		current = current.Parent()
	}

	for i, j := 0, len(chemin)-1; i < j; i, j = i+1, j-1 {
		chemin[i], chemin[j] = chemin[j], chemin[i]
	}

	return chemin
}

// Helper functions

// positionKey crée une clé unique pour une position (pour les maps)
func positionKey(pos *shared.Position) string {
	// Utilisation de fmt.Sprintf pour créer une clé unique
	// Exemple: "10,25" pour x=10, y=25
	return fmt.Sprintf("%d,%d", pos.X(), pos.Y())
}

// abs retourne la valeur absolue d'un entier
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
