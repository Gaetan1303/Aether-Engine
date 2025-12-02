package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// CommandFactory crée des commandes selon les paramètres
// Factory Pattern pour la création de commandes
type CommandFactory struct {
	combat *domain.Combat
}

// NewCommandFactory crée une nouvelle factory
func NewCommandFactory(combat *domain.Combat) *CommandFactory {
	return &CommandFactory{
		combat: combat,
	}
}

// CreateMoveCommand crée une commande de déplacement
func (f *CommandFactory) CreateMoveCommand(actor *domain.Unite, targetX, targetY int) (Command, error) {
	position := f.combat.Grille().Position(targetX, targetY)
	return NewMoveCommand(actor, f.combat, position), nil
}

// CreateAttackCommand crée une commande d'attaque
func (f *CommandFactory) CreateAttackCommand(actor *domain.Unite, targetID domain.UnitID) (Command, error) {
	target := f.combat.TrouverUnite(targetID)
	return NewAttackCommand(actor, f.combat, target), nil
}

// CreateSkillCommand crée une commande de compétence
func (f *CommandFactory) CreateSkillCommand(actor *domain.Unite, skillID string, targetIDs []domain.UnitID) (Command, error) {
	skill := actor.ObtenirCompetence(domain.CompetenceID(skillID))

	targets := make([]*domain.Unite, 0)
	for _, targetID := range targetIDs {
		target := f.combat.TrouverUnite(targetID)
		targets = append(targets, target)
	}

	return NewSkillCommand(actor, f.combat, skill, targets), nil
}

// CreateItemCommand crée une commande d'objet
func (f *CommandFactory) CreateItemCommand(actor *domain.Unite, itemID string, targetID domain.UnitID) (Command, error) {
	itemInterface := f.combat.ObtenirObjet(itemID)
	item, ok := itemInterface.(*shared.Item)
	if !ok || item == nil {
		return nil, fmt.Errorf("objet %s non trouvé", itemID)
	}
	target := f.combat.TrouverUnite(targetID)

	return NewItemCommand(actor, f.combat, item, target), nil
}

// CreateFleeCommand crée une commande de fuite
func (f *CommandFactory) CreateFleeCommand(actor *domain.Unite) (Command, error) {
	return NewFleeCommand(actor, f.combat), nil
}

// CreateWaitCommand crée une commande d'attente
func (f *CommandFactory) CreateWaitCommand(actor *domain.Unite) (Command, error) {
	return NewWaitCommand(actor, f.combat), nil
}
