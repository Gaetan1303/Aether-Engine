package states

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// InitializingState gère l'initialisation du combat
// - Charge Teams + Units
// - Initialise BattleGrid
// - Configure Ticker (ATB)
// - Validation des règles
type InitializingState struct {
	BaseState
}

// NewInitializingState crée un nouvel état Initializing
func NewInitializingState() *InitializingState {
	return &InitializingState{
		BaseState: BaseState{
			name: "Initializing",
			allowedTransitions: map[string]bool{
				"Ready":  true,
				"Failed": true,
			},
		},
	}
}

// Enter initialise le combat
func (s *InitializingState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// 1. Valider les équipes
	if err := s.validateTeams(ctx); err != nil {
		return fmt.Errorf("validation des équipes échouée: %w", err)
	}

	// 2. Initialiser les jauges ATB
	if err := s.initializeATB(ctx); err != nil {
		return fmt.Errorf("initialisation ATB échouée: %w", err)
	}

	// 3. Vérifier la grille de combat
	if err := s.validateGrid(ctx); err != nil {
		return fmt.Errorf("validation de la grille échouée: %w", err)
	}

	fmt.Printf("[State] Initialisation complète\n")
	return nil
}

// Exit est appelé lors de la sortie
func (s *InitializingState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Initializing
func (s *InitializingState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventSetupComplete:
		// Initialisation réussie, passer à Ready
		return NewReadyState(), nil

	case EventValidationError:
		// Erreur de validation, passer à Failed
		return NewFailedState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// validateTeams vérifie que les équipes sont valides
func (s *InitializingState) validateTeams(ctx *CombatContext) error {
	equipes := ctx.Combat.Equipes()

	if len(equipes) < domain.MinEquipesPourCombat {
		return fmt.Errorf("au moins 2 équipes requises pour un combat")
	}

	for _, equipe := range equipes {
		membres := equipe.Membres()
		if len(membres) == 0 {
			return fmt.Errorf("l'équipe %s n'a aucun membre", equipe.ID())
		}

		// Vérifier que chaque membre a une position valide
		for _, membre := range membres {
			if membre.Position() == nil {
				return fmt.Errorf("l'unité %s n'a pas de position définie", membre.ID())
			}
		}
	}

	return nil
}

// initializeATB initialise les jauges ATB de toutes les unités
func (s *InitializingState) initializeATB(ctx *CombatContext) error {
	equipes := ctx.Combat.Equipes()

	for _, equipe := range equipes {
		for _, unite := range equipe.Membres() {
			// Vitesse de remplissage basée sur SPD
			speed := s.calculateATBSpeed(unite)
			ctx.ATBSystem.InitializeGauge(unite.ID(), speed)
		}
	}

	return nil
}

// calculateATBSpeed calcule la vitesse de remplissage ATB basée sur SPD
func (s *InitializingState) calculateATBSpeed(unite *domain.Unite) int {
	// Formule: Speed = SPD / 10 (min 1, max 10)
	spd := unite.Stats().SPD
	speed := spd / 10

	if speed < 1 {
		speed = 1
	}
	if speed > 10 {
		speed = 10
	}

	return speed
}

// validateGrid vérifie que la grille de combat est valide
func (s *InitializingState) validateGrid(ctx *CombatContext) error {
	grille := ctx.Combat.Grille()

	if grille == nil {
		return fmt.Errorf("grille de combat non initialisée")
	}

	if grille.Largeur() <= 0 || grille.Hauteur() <= 0 {
		return fmt.Errorf("dimensions de grille invalides: %dx%d",
			grille.Largeur(), grille.Hauteur())
	}

	return nil
}
