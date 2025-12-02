package domain

import (
	"errors"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// UnitInventory gère l'inventaire d'une unité
// Responsabilités: Compétences et objets
// Single Responsibility Principle - Une seule raison de changer: gestion inventaire
type UnitInventory struct {
	skills []*Competence
	items  []shared.ObjetID
}

// NewUnitInventory crée un nouveau gestionnaire d'inventaire
func NewUnitInventory() *UnitInventory {
	return &UnitInventory{
		skills: make([]*Competence, 0),
		items:  make([]shared.ObjetID, 0),
	}
}

// Skills retourne toutes les compétences
func (inv *UnitInventory) Skills() []*Competence {
	return inv.skills
}

// Items retourne tous les objets
func (inv *UnitInventory) Items() []shared.ObjetID {
	return inv.items
}

// AddSkill ajoute une compétence à l'inventaire
func (inv *UnitInventory) AddSkill(skill *Competence) error {
	if skill == nil {
		return errors.New("compétence nil")
	}

	// Vérifier si la compétence existe déjà
	for _, existing := range inv.skills {
		if existing.ID() == skill.ID() {
			return errors.New("compétence déjà apprise")
		}
	}

	inv.skills = append(inv.skills, skill)
	return nil
}

// GetSkill retourne une compétence par ID
func (inv *UnitInventory) GetSkill(skillID CompetenceID) *Competence {
	for _, skill := range inv.skills {
		if skill.ID() == skillID {
			return skill
		}
	}
	return nil
}

// HasSkill vérifie si une compétence est disponible
func (inv *UnitInventory) HasSkill(skillID CompetenceID) bool {
	return inv.GetSkill(skillID) != nil
}

// IsSkillReady vérifie si une compétence est prête
func (inv *UnitInventory) IsSkillReady(skillID CompetenceID) bool {
	skill := inv.GetSkill(skillID)
	if skill == nil {
		return false
	}
	return !skill.EstEnCooldown()
}

// ActivateSkillCooldown active le cooldown d'une compétence
func (inv *UnitInventory) ActivateSkillCooldown(skillID CompetenceID, duration int) error {
	skill := inv.GetSkill(skillID)
	if skill == nil {
		return errors.New("compétence introuvable")
	}

	skill.ActiverCooldown()
	return nil
}

// DecrementAllCooldowns décrémente les cooldowns de toutes les compétences
func (inv *UnitInventory) DecrementAllCooldowns() {
	for _, skill := range inv.skills {
		skill.DecrémenterCooldown()
	}
}

// AddItem ajoute un objet à l'inventaire
func (inv *UnitInventory) AddItem(itemID shared.ObjetID, quantity int) {
	for i := 0; i < quantity; i++ {
		inv.items = append(inv.items, itemID)
	}
}

// RemoveItem retire un objet de l'inventaire
func (inv *UnitInventory) RemoveItem(itemID shared.ObjetID, quantity int) error {
	removed := 0

	for i := len(inv.items) - 1; i >= 0 && removed < quantity; i-- {
		if inv.items[i] == itemID {
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			removed++
		}
	}

	if removed < quantity {
		return errors.New("quantité insuffisante")
	}

	return nil
}

// HasItem vérifie si un objet est dans l'inventaire
func (inv *UnitInventory) HasItem(itemID shared.ObjetID) bool {
	for _, item := range inv.items {
		if item == itemID {
			return true
		}
	}
	return false
}

// CountItem compte le nombre d'exemplaires d'un objet
func (inv *UnitInventory) CountItem(itemID shared.ObjetID) int {
	count := 0
	for _, item := range inv.items {
		if item == itemID {
			count++
		}
	}
	return count
}

// SkillCount retourne le nombre de compétences
func (inv *UnitInventory) SkillCount() int {
	return len(inv.skills)
}

// ItemCount retourne le nombre total d'objets
func (inv *UnitInventory) ItemCount() int {
	return len(inv.items)
}

// ClearItems vide l'inventaire d'objets
func (inv *UnitInventory) ClearItems() {
	inv.items = make([]shared.ObjetID, 0)
}
