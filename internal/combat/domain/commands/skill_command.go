package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// SkillCommand représente l'utilisation d'une compétence
type SkillCommand struct {
	*BaseCommand
	skill   *domain.Competence
	targets []*domain.Unite
}

// NewSkillCommand crée une nouvelle commande de skill
func NewSkillCommand(actor *domain.Unite, combat *domain.Combat, skill *domain.Competence, targets []*domain.Unite) *SkillCommand {
	return &SkillCommand{
		BaseCommand: NewBaseCommand(actor, combat, CommandTypeSkill),
		skill:       skill,
		targets:     targets,
	}
}

// Validate vérifie si le skill peut être utilisé
func (c *SkillCommand) Validate() error {
	// 1. Vérifier que l'acteur peut agir
	if !c.actor.PeutAgir() {
		return fmt.Errorf("l'unité %s ne peut pas agir", c.actor.Nom())
	}

	// 2. Vérifier que le skill existe
	if c.skill == nil {
		return fmt.Errorf("aucune compétence spécifiée")
	}

	// 3. Vérifier que l'acteur possède ce skill
	if c.actor.ObtenirCompetence(c.skill.ID()) == nil {
		return fmt.Errorf("l'unité %s ne possède pas la compétence %s", c.actor.Nom(), c.skill.Nom())
	}

	// 4. Vérifier le coût MP
	if c.actor.StatsActuelles().MP < c.skill.CoutMP() {
		return fmt.Errorf("MP insuffisant (coût: %d, disponible: %d)", c.skill.CoutMP(), c.actor.StatsActuelles().MP)
	}

	// 5. Vérifier le cooldown
	if !c.actor.SkillEstPret(c.skill.ID()) {
		return fmt.Errorf("compétence en cooldown")
	}

	// 6. Vérifier les cibles
	if len(c.targets) == 0 {
		return fmt.Errorf("aucune cible spécifiée")
	}

	for _, target := range c.targets {
		if target.EstEliminee() {
			return fmt.Errorf("la cible %s est déjà éliminée", target.Nom())
		}

		// Vérifier la portée
		distance := c.actor.Position().Distance(target.Position())

		if distance > c.skill.Portee() {
			return fmt.Errorf("cible %s hors de portée (distance: %d, portée: %d)", target.Nom(), distance, c.skill.Portee())
		}
	}

	// 7. Vérifier le statut Silence (interdit les skills)
	if c.actor.EstSilence() {
		return fmt.Errorf("l'unité est Silencée, impossible d'utiliser des compétences")
	}

	return nil
}

// Execute utilise la compétence
func (c *SkillCommand) Execute() (*CommandResult, error) {
	// Créer un snapshot avant modification
	c.CreateSnapshot()

	// Consommer les MP
	c.actor.ConsommerMP(c.skill.CoutMP())

	// Activer le cooldown
	c.actor.ActiverCooldown(c.skill.ID(), c.skill.Cooldown())

	// Créer le résultat
	result := &CommandResult{
		Success: true,
		Message: fmt.Sprintf("%s utilise %s", c.actor.Nom(), c.skill.Nom()),
		CostMP:  c.skill.CoutMP(),
		Effects: make([]CommandEffect, 0),
	}

	// Appliquer les effets selon le type de skill
	calculator := c.combat.GetDamageCalculator()

	for _, target := range c.targets {
		// Utiliser le type de compétence pour déterminer l'effet
		switch c.skill.Type() {
		case domain.CompetenceAttaque, domain.CompetenceMagie:
			// Compétence de dégâts
			degats := calculator.Calculate(c.actor, target, c.skill)
			target.RecevoirDegats(degats)
			result.DamageDealt += degats
			result.Effects = append(result.Effects, CommandEffect{
				Type:     EffectTypeDamage,
				TargetID: target.ID(),
				Value:    degats,
			})

		case domain.CompetenceSoin:
			// Compétence de soin - utiliser les dégâts de base comme valeur de soin
			soins := c.skill.DegatsBase()
			target.Soigner(soins)
			result.HealingDone += soins
			result.Effects = append(result.Effects, CommandEffect{
				Type:     EffectTypeHealing,
				TargetID: target.ID(),
				Value:    soins,
			})

		case domain.CompetenceUtilitaire:
			// Compétences de support (buff, debuff, statut)
			// TODO: Implémenter complètement le système de statut depuis compétence
			fmt.Printf("[Skill] Support/Statut depuis %s sur %s (à implémenter)\n", c.skill.Nom(), target.Nom())
			result.Effects = append(result.Effects, CommandEffect{
				Type:     EffectTypeStatus,
				TargetID: target.ID(),
			})

		default:
			// Autres types de compétences (Buff, Support, etc.)
			fmt.Printf("[Skill] Type de compétence %v non géré\n", c.skill.Type())
		}
	}

	return result, nil
}

// Rollback annule l'utilisation du skill
func (c *SkillCommand) Rollback() error {
	if c.snapshot == nil {
		return fmt.Errorf("aucun snapshot disponible pour rollback")
	}

	// Restaurer les MP de l'acteur
	// Note: nécessite une méthode SetMP sur Unite
	fmt.Printf("[Rollback] Impossible de restaurer les MP de %s (méthode SetMP non implémentée)\n", c.actor.Nom())

	return nil
}

// abs est défini dans attack_command.go
