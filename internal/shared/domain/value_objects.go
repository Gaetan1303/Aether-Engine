package domain

import (
	"errors"
	"math"
)

// Position représente une position sur la grille (Value Object)
type Position struct {
	x int
	y int
}

// NewPosition crée une nouvelle position
func NewPosition(x, y int) (*Position, error) {
	if x < 0 || y < 0 {
		return nil, errors.New("les coordonnées doivent être positives")
	}
	return &Position{x: x, y: y}, nil
}

// X retourne la coordonnée X
func (p *Position) X() int { return p.x }

// Y retourne la coordonnée Y
func (p *Position) Y() int { return p.y }

// Equals vérifie l'égalité avec une autre position
func (p *Position) Equals(autre *Position) bool {
	return p.x == autre.x && p.y == autre.y
}

// Distance calcule la distance de Manhattan vers une autre position
func (p *Position) Distance(autre *Position) int {
	return abs(p.x-autre.x) + abs(p.y-autre.y)
}

// DistanceEuclidienne calcule la distance euclidienne vers une autre position
func (p *Position) DistanceEuclidienne(autre *Position) float64 {
	dx := float64(p.x - autre.x)
	dy := float64(p.y - autre.y)
	return math.Sqrt(dx*dx + dy*dy)
}

// EstAdjacente vérifie si une autre position est adjacente (4 directions)
func (p *Position) EstAdjacente(autre *Position) bool {
	return p.Distance(autre) == 1
}

// EstAdjacenteDiagonale vérifie si une autre position est adjacente (8 directions)
func (p *Position) EstAdjacenteDiagonale(autre *Position) bool {
	dx := abs(p.x - autre.x)
	dy := abs(p.y - autre.y)
	return (dx <= 1 && dy <= 1) && !(dx == 0 && dy == 0)
}

// Stats représente les statistiques d'une unité (Value Object)
type Stats struct {
	HP      int // Points de vie
	MP      int // Points de magie
	Stamina int // Endurance
	ATK     int // Attaque physique
	DEF     int // Défense physique
	MATK    int // Attaque magique
	MDEF    int // Défense magique
	SPD     int // Vitesse (initiative)
	MOV     int // Mouvement (cases par tour)
	ATH     int // Attack Hit (chance de toucher en %)
}

// NewStats crée un nouveau set de stats
func NewStats(hp, mp, stamina, atk, def, matk, mdef, spd, mov, ath int) (*Stats, error) {
	if hp <= 0 || mp < 0 || stamina < 0 {
		return nil, errors.New("HP doit être > 0, MP et Stamina >= 0")
	}
	if atk < 0 || def < 0 || matk < 0 || mdef < 0 {
		return nil, errors.New("les stats offensives/défensives doivent être >= 0")
	}
	if spd <= 0 || mov <= 0 {
		return nil, errors.New("SPD et MOV doivent être > 0")
	}
	if ath < 0 || ath > 100 {
		return nil, errors.New("ATH doit être entre 0 et 100")
	}

	return &Stats{
		HP:      hp,
		MP:      mp,
		Stamina: stamina,
		ATK:     atk,
		DEF:     def,
		MATK:    matk,
		MDEF:    mdef,
		SPD:     spd,
		MOV:     mov,
		ATH:     ath,
	}, nil
}

// Clone crée une copie des stats
func (s *Stats) Clone() *Stats {
	return &Stats{
		HP:      s.HP,
		MP:      s.MP,
		Stamina: s.Stamina,
		ATH:     s.ATH,
		ATK:     s.ATK,
		DEF:     s.DEF,
		MATK:    s.MATK,
		MDEF:    s.MDEF,
		SPD:     s.SPD,
		MOV:     s.MOV,
	}
}

// Equals vérifie l'égalité avec d'autres stats
func (s *Stats) Equals(autre *Stats) bool {
	return s.HP == autre.HP &&
		s.MP == autre.MP &&
		s.Stamina == autre.Stamina &&
		s.ATK == autre.ATK &&
		s.DEF == autre.DEF &&
		s.MATK == autre.MATK &&
		s.MDEF == autre.MDEF &&
		s.SPD == autre.SPD &&
		s.MOV == autre.MOV &&
		s.ATH == autre.ATH
}

// GrilleCombat représente la grille de combat (Value Object)
type GrilleCombat struct {
	largeur  int
	hauteur  int
	cellules [][]TypeCellule
}

// TypeCellule représente le type de terrain d'une cellule
type TypeCellule int

const (
	CelluleNormale TypeCellule = iota
	CelluleObstacle
	CelluleDifficile // Coût de mouvement x2
	CelluleDanger    // Inflige dégâts périodiques
	CelluleSoin      // Soigne périodiquement
)

// NewGrilleCombat crée une nouvelle grille de combat
func NewGrilleCombat(largeur, hauteur int) (*GrilleCombat, error) {
	if largeur <= 0 || hauteur <= 0 {
		return nil, errors.New("dimensions doivent être > 0")
	}

	// Initialiser avec des cellules normales
	cellules := make([][]TypeCellule, hauteur)
	for i := range cellules {
		cellules[i] = make([]TypeCellule, largeur)
		for j := range cellules[i] {
			cellules[i][j] = CelluleNormale
		}
	}

	return &GrilleCombat{
		largeur:  largeur,
		hauteur:  hauteur,
		cellules: cellules,
	}, nil
}

// Largeur retourne la largeur de la grille
func (g *GrilleCombat) Largeur() int { return g.largeur }

// Hauteur retourne la hauteur de la grille
func (g *GrilleCombat) Hauteur() int { return g.hauteur }

// EstDansLimites vérifie si une position est dans la grille
func (g *GrilleCombat) EstDansLimites(pos *Position) bool {
	return pos.X() >= 0 && pos.X() < g.largeur &&
		pos.Y() >= 0 && pos.Y() < g.hauteur
}

// ObtenirTypeCellule retourne le type de cellule à une position
func (g *GrilleCombat) ObtenirTypeCellule(pos *Position) (TypeCellule, error) {
	if !g.EstDansLimites(pos) {
		return 0, errors.New("position hors limites")
	}
	return g.cellules[pos.Y()][pos.X()], nil
}

// DefinirTypeCellule définit le type de cellule à une position
func (g *GrilleCombat) DefinirTypeCellule(pos *Position, typeCellule TypeCellule) error {
	if !g.EstDansLimites(pos) {
		return errors.New("position hors limites")
	}
	g.cellules[pos.Y()][pos.X()] = typeCellule
	return nil
}

// EstTraversable vérifie si une cellule peut être traversée
func (g *GrilleCombat) EstTraversable(pos *Position) bool {
	if !g.EstDansLimites(pos) {
		return false
	}
	typeCellule := g.cellules[pos.Y()][pos.X()]
	return typeCellule != CelluleObstacle
}

// Position crée une nouvelle position à partir de coordonnées (pour CommandFactory)
func (g *GrilleCombat) Position(x, y int) *Position {
	pos, _ := NewPosition(x, y)
	return pos
}

// CoutDeplacement retourne le coût de déplacement pour une cellule
func (g *GrilleCombat) CoutDeplacement(pos *Position) int {
	if !g.EstDansLimites(pos) {
		return -1 // Invalide
	}

	typeCellule := g.cellules[pos.Y()][pos.X()]
	switch typeCellule {
	case CelluleNormale, CelluleDanger, CelluleSoin:
		return 1
	case CelluleDifficile:
		return 2
	case CelluleObstacle:
		return -1 // Non traversable
	default:
		return 1
	}
}

// PositionsAdjacentes retourne les positions adjacentes traversables (4 directions)
func (g *GrilleCombat) PositionsAdjacentes(pos *Position) []*Position {
	adjacentes := make([]*Position, 0, 4)

	directions := []struct{ dx, dy int }{
		{0, -1}, // Haut
		{1, 0},  // Droite
		{0, 1},  // Bas
		{-1, 0}, // Gauche
	}

	for _, dir := range directions {
		nouvellePos, _ := NewPosition(pos.X()+dir.dx, pos.Y()+dir.dy)
		if g.EstTraversable(nouvellePos) {
			adjacentes = append(adjacentes, nouvellePos)
		}
	}

	return adjacentes
}

// PositionsADansPortee retourne toutes les positions dans une portée donnée
func (g *GrilleCombat) PositionsADansPortee(centre *Position, portee int) []*Position {
	positions := make([]*Position, 0)

	for y := 0; y < g.hauteur; y++ {
		for x := 0; x < g.largeur; x++ {
			pos, _ := NewPosition(x, y)
			if centre.Distance(pos) <= portee {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

// ChercheChemin trouve le plus court chemin entre deux positions (A*)
func (g *GrilleCombat) ChercheChemin(depart, arrivee *Position) ([]*Position, int) {
	// TODO: Implémenter A*
	// Pour l'instant, retourner un chemin vide
	return []*Position{}, -1
}

// Statut représente un effet de statut sur une unité (Value Object)
type Statut struct {
	typeStatut        TypeStatut
	duree             int // Tours restants
	puissance         int // Intensité de l'effet
	modificateurs     []ModificateurStat
	bloqueActions     bool
	bloqueDeplacement bool
}

// TypeStatut énumère les types de statuts
type TypeStatut int

const (
	StatutPoison TypeStatut = iota
	StatutBrulure
	StatutGel
	StatutParalysie
	StatutSommeil
	StatutStun
	StatutRoot
	StatutBuff
	StatutDebuff
	StatutRegeneration
	StatutBouclier
	StatutMort    // Ajouté pour Step C
	StatutSilence // Ajouté pour Step C
)

// Alias pour compatibilité avec code Step C
const (
	TypeStatutMort    = StatutMort
	TypeStatutSilence = StatutSilence
	TypeStatutStun    = StatutStun
	TypeStatutRoot    = StatutRoot
	TypeStatutPoison  = StatutPoison
)

// NewStatut crée un nouveau statut
func NewStatut(typeStatut TypeStatut, duree, puissance int) *Statut {
	s := &Statut{
		typeStatut:        typeStatut,
		duree:             duree,
		puissance:         puissance,
		modificateurs:     make([]ModificateurStat, 0),
		bloqueActions:     false,
		bloqueDeplacement: false,
	}

	// Définir les propriétés selon le type
	switch typeStatut {
	case StatutStun:
		s.bloqueActions = true
		s.bloqueDeplacement = true
	case StatutRoot:
		s.bloqueDeplacement = true
	case StatutSommeil:
		s.bloqueActions = true
		s.bloqueDeplacement = true
	case StatutParalysie:
		s.modificateurs = append(s.modificateurs, ModificateurStat{
			Stat:   "SPD",
			Valeur: -puissance,
		})
	}

	return s
}

// Getters
func (s *Statut) Type() TypeStatut                  { return s.typeStatut }
func (s *Statut) Duree() int                        { return s.duree }
func (s *Statut) Puissance() int                    { return s.puissance }
func (s *Statut) Modificateurs() []ModificateurStat { return s.modificateurs }
func (s *Statut) BloqueActions() bool               { return s.bloqueActions }
func (s *Statut) BloqueDeplacement() bool           { return s.bloqueDeplacement }

// DecrémenterDuree décrémente la durée du statut
func (s *Statut) DecrémenterDuree() {
	if s.duree > 0 {
		s.duree--
	}
}

// EstExpire vérifie si le statut a expiré
func (s *Statut) EstExpire() bool {
	return s.duree <= 0
}

// Refresh renouvelle le statut
func (s *Statut) Refresh(nouvelleDuree int) error {
	if nouvelleDuree > s.duree {
		s.duree = nouvelleDuree
	}
	return nil
}

// AppliquerEffet applique l'effet initial du statut
// Dependency Inversion Principle (SOLID) - Dépend d'une interface, pas d'un type concret
func (s *Statut) AppliquerEffet(cible StatsModifiable) {
	// Appliquer les modificateurs de stats
	for _, mod := range s.modificateurs {
		cible.AppliquerModificateurStat(&mod)
	}
}

// RetirerEffet retire l'effet du statut
// Dependency Inversion Principle (SOLID) - Dépend d'une interface, pas d'un type concret
func (s *Statut) RetirerEffet(cible StatsModifiable) {
	// Retirer les modificateurs de stats
	for _, mod := range s.modificateurs {
		cible.RetirerModificateurStat(&mod)
	}
}

// AppliquerEffetPeriodique applique l'effet périodique (début de tour)
// Dependency Inversion Principle (SOLID) - Dépend d'une interface, pas d'un type concret
func (s *Statut) AppliquerEffetPeriodique(cible StatsModifiable) *EffetStatut {
	switch s.typeStatut {
	case StatutPoison, StatutBrulure:
		degats := s.puissance
		cible.RecevoirDegats(degats)
		return &EffetStatut{Type: s.typeStatut, Valeur: degats}
	case StatutRegeneration:
		soin := s.puissance
		cible.RecevoirSoin(soin)
		return &EffetStatut{Type: s.typeStatut, Valeur: soin}
	}
	return nil
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
