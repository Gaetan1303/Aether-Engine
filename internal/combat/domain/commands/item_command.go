package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// ItemCommand représente l'utilisation d'un objet (Potion, Éther, Antidote)
type ItemCommand struct {
	*BaseCommand
	item   *shared.Item
	target *domain.Unite
}

// NewItemCommand crée une nouvelle commande d'utilisation d'objet
func NewItemCommand(actor *domain.Unite, combat *domain.Combat, item *shared.Item, target *domain.Unite) *ItemCommand {
	return &ItemCommand{
		BaseCommand: NewBaseCommand(actor, combat, CommandTypeItem),
		item:        item,
		target:      target,
	}
}

// Validate vérifie si l'objet peut être utilisé
func (c *ItemCommand) Validate() error {
	// 1. Vérifier que l'acteur peut agir
	if !c.actor.PeutAgir() {
		return fmt.Errorf("l'unité %s ne peut pas agir", c.actor.Nom())
	}

	// 2. Vérifier que l'objet existe
	if c.item == nil {
		return fmt.Errorf("aucun objet spécifié")
	}

	// 3. Vérifier que l'objet est dans l'inventaire
	if !c.combat.PossedeObjet(c.item.GetID()) {
		return fmt.Errorf("objet %s non trouvé dans l'inventaire", c.item.GetName())
	}

	// 4. Vérifier la quantité disponible
	quantite := c.combat.ObtenirQuantiteObjet(c.item.GetID())
	if quantite <= 0 {
		return fmt.Errorf("quantité insuffisante pour %s", c.item.GetName())
	}

	// 5. Vérifier la cible
	if c.target == nil {
		return fmt.Errorf("aucune cible spécifiée")
	}

	if c.target.EstEliminee() {
		return fmt.Errorf("la cible %s est déjà éliminée", c.target.Nom())
	}

	// 6. Vérifier la portée (objets utilisables à distance)
	distance := abs(c.actor.Position().X()-c.target.Position().X()) +
		abs(c.actor.Position().Y()-c.target.Position().Y())

	if distance > c.item.GetRange() {
		return fmt.Errorf("cible hors de portée (distance: %d, portée: %d)", distance, c.item.GetRange())
	}

	// 7. Vérifier les restrictions d'usage
	if !c.canUseItemOnTarget() {
		return fmt.Errorf("impossible d'utiliser %s sur %s", c.item.GetName(), c.target.Nom())
	}

	return nil
}

// canUseItemOnTarget vérifie si l'objet peut être utilisé sur la cible
func (c *ItemCommand) canUseItemOnTarget() bool {
	switch c.item.GetItemType() {
	case shared.ItemTypePotion:
		// Potion: seulement sur alliés vivants avec HP < Max
		return c.actor.TeamID() == c.target.TeamID() &&
			!c.target.EstEliminee() &&
			c.target.HPActuels() < c.target.Stats().HP

	case shared.ItemTypeEther:
		// Éther: seulement sur alliés vivants avec MP < Max
		return c.actor.TeamID() == c.target.TeamID() &&
			!c.target.EstEliminee() &&
			c.target.StatsActuelles().MP < c.target.Stats().MP

	case shared.ItemTypeAntidote:
		// Antidote: seulement sur alliés empoisonnés
		return c.actor.TeamID() == c.target.TeamID() &&
			!c.target.EstEliminee() &&
			c.target.EstEmpoisonne()

	case shared.ItemTypeRevive:
		// Revive: seulement sur alliés KO
		return c.actor.TeamID() == c.target.TeamID() &&
			c.target.EstEliminee()

	case shared.ItemTypeBomb:
		// Bombe: seulement sur ennemis vivants
		return c.actor.TeamID() != c.target.TeamID() &&
			!c.target.EstEliminee()

	default:
		return true
	}
}

// Execute utilise l'objet
func (c *ItemCommand) Execute() (*CommandResult, error) {
	// Créer un snapshot avant modification
	c.CreateSnapshot()

	// Consommer l'objet
	c.combat.ConsommerObjet(c.item.GetID(), 1)

	// Créer le résultat
	result := &CommandResult{
		Success: true,
		Message: fmt.Sprintf("%s utilise %s sur %s", c.actor.Nom(), c.item.GetName(), c.target.Nom()),
		Effects: make([]CommandEffect, 0),
	}

	// Appliquer l'effet selon le type d'objet
	switch c.item.GetItemType() {
	case shared.ItemTypePotion:
		// Soigner les HP
		soins := c.item.EffectValue()
		c.target.Soigner(soins)
		result.HealingDone = soins
		result.Effects = append(result.Effects, CommandEffect{
			Type:     EffectTypeHealing,
			TargetID: c.target.ID(),
			Value:    soins,
		})

	case shared.ItemTypeEther:
		// Restaurer les MP
		mp := c.item.EffectValue()
		c.target.RestaurerMP(mp)
		result.Effects = append(result.Effects, CommandEffect{
			Type:     EffectTypeHealing,
			TargetID: c.target.ID(),
			Value:    mp,
		})

	case shared.ItemTypeAntidote:
		// Retirer le poison
		c.target.RetirerStatut(shared.TypeStatutPoison)
		result.Effects = append(result.Effects, CommandEffect{
			Type:     EffectTypeStatus,
			TargetID: c.target.ID(),
		})

	case shared.ItemTypeRevive:
		// Ressusciter l'unité
		c.target.Ressusciter(c.item.EffectValue())
		result.HealingDone = c.item.EffectValue()
		result.Effects = append(result.Effects, CommandEffect{
			Type:     EffectTypeHealing,
			TargetID: c.target.ID(),
			Value:    c.item.EffectValue(),
		})

	case shared.ItemTypeBomb:
		// Infliger des dégâts
		degats := c.item.EffectValue()
		c.target.RecevoirDegats(degats)
		result.DamageDealt = degats
		result.Effects = append(result.Effects, CommandEffect{
			Type:     EffectTypeDamage,
			TargetID: c.target.ID(),
			Value:    degats,
		})
	}

	return result, nil
}

// Rollback annule l'utilisation de l'objet
func (c *ItemCommand) Rollback() error {
	if c.snapshot == nil {
		return fmt.Errorf("aucun snapshot disponible pour rollback")
	}

	// Rendre l'objet à l'inventaire
	c.combat.AjouterObjet(c.item.GetID(), 1)

	// Restaurer l'état de la cible
	// Note: nécessite des méthodes SetHP/SetMP sur Unite
	fmt.Printf("[Rollback] Objet %s rendu à l'inventaire\n", c.item.GetName())

	return nil
}
