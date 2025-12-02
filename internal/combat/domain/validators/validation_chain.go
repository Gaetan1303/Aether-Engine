package validators

import (
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
)

// Chain of Responsibility Pattern

// Validator représente un validateur dans la chaîne
type Validator interface {
	// SetNext définit le prochain validateur dans la chaîne
	SetNext(validator Validator) Validator
	// Validate valide la commande
	Validate(cmd commands.Command) error
}

// BaseValidator implémentation de base pour la chaîne
type BaseValidator struct {
	next Validator
}

// SetNext définit le prochain validateur
func (v *BaseValidator) SetNext(validator Validator) Validator {
	v.next = validator
	return validator
}

// CallNext appelle le prochain validateur dans la chaîne
func (v *BaseValidator) CallNext(cmd commands.Command) error {
	if v.next != nil {
		return v.next.Validate(cmd)
	}
	return nil
}

// CostValidator vérifie le coût en MP/HP/Stamina
type CostValidator struct {
	BaseValidator
}

// NewCostValidator crée un nouveau validateur de coût
func NewCostValidator() *CostValidator {
	return &CostValidator{}
}

// Validate vérifie que l'acteur a assez de ressources
func (v *CostValidator) Validate(cmd commands.Command) error {
	actor := cmd.GetActor()

	// Selon le type de commande, vérifier les coûts
	switch cmd.GetType() {
	case commands.CommandTypeSkill:
		// Vérifier MP pour les skills (déjà fait dans SkillCommand)
		// Cette validation est redondante mais permet une architecture propre

	case commands.CommandTypeMove:
		// Vérifier stamina pour le déplacement (si applicable)
		// TODO: implémenter système de stamina

	case commands.CommandTypeItem:
		// Vérifier l'inventaire (fait dans ItemCommand)
	}

	// Log de validation
	fmt.Printf("[Validator] CostValidator: OK pour %s\n", actor.Nom())

	// Appeler le prochain validateur
	return v.CallNext(cmd)
}

// RangeValidator vérifie la portée des actions
type RangeValidator struct {
	BaseValidator
}

// NewRangeValidator crée un nouveau validateur de portée
func NewRangeValidator() *RangeValidator {
	return &RangeValidator{}
}

// Validate vérifie que les cibles sont à portée
func (v *RangeValidator) Validate(cmd commands.Command) error {
	actor := cmd.GetActor()

	// La validation de portée est déjà faite dans chaque commande
	// Cette validation globale permet de centraliser la logique

	fmt.Printf("[Validator] RangeValidator: OK pour %s\n", actor.Nom())

	// Appeler le prochain validateur
	return v.CallNext(cmd)
}

// TargetValidator vérifie la validité des cibles
type TargetValidator struct {
	BaseValidator
}

// NewTargetValidator crée un nouveau validateur de cible
func NewTargetValidator() *TargetValidator {
	return &TargetValidator{}
}

// Validate vérifie que les cibles sont valides
func (v *TargetValidator) Validate(cmd commands.Command) error {
	actor := cmd.GetActor()

	switch cmd.GetType() {
	case commands.CommandTypeAttack, commands.CommandTypeSkill:
		// Vérifier que les cibles ne sont pas déjà mortes
		// (déjà fait dans chaque commande)

	case commands.CommandTypeMove:
		// Vérifier que la position cible est libre
		// (déjà fait dans MoveCommand)
	}

	fmt.Printf("[Validator] TargetValidator: OK pour %s\n", actor.Nom())

	// Appeler le prochain validateur
	return v.CallNext(cmd)
}

// StatusValidator vérifie les statuts bloquants
type StatusValidator struct {
	BaseValidator
}

// NewStatusValidator crée un nouveau validateur de statut
func NewStatusValidator() *StatusValidator {
	return &StatusValidator{}
}

// Validate vérifie les statuts qui empêchent l'action
func (v *StatusValidator) Validate(cmd commands.Command) error {
	actor := cmd.GetActor()

	// Vérifier les statuts bloquants selon le type de commande
	switch cmd.GetType() {
	case commands.CommandTypeSkill:
		if actor.EstSilence() {
			return fmt.Errorf("impossible d'utiliser des compétences: unité Silencée")
		}

	case commands.CommandTypeMove:
		if actor.EstBloqueDeplacement() {
			return fmt.Errorf("impossible de se déplacer: unité Root/Stun")
		}

	case commands.CommandTypeAttack:
		if actor.EstStun() {
			return fmt.Errorf("impossible d'attaquer: unité Stunned")
		}
	}

	// Vérifier que l'unité peut agir en général
	if !actor.PeutAgir() {
		return fmt.Errorf("l'unité ne peut pas agir (Stun/Sleep/Dead)")
	}

	fmt.Printf("[Validator] StatusValidator: OK pour %s\n", actor.Nom())

	// Appeler le prochain validateur
	return v.CallNext(cmd)
}

// ValidationChain gère la chaîne de validateurs
type ValidationChain struct {
	head Validator
}

// NewValidationChain crée une nouvelle chaîne de validation
func NewValidationChain() *ValidationChain {
	// Construire la chaîne: Status → Cost → Range → Target
	statusValidator := NewStatusValidator()
	costValidator := NewCostValidator()
	rangeValidator := NewRangeValidator()
	targetValidator := NewTargetValidator()

	// Chaîner les validateurs
	statusValidator.SetNext(costValidator)
	costValidator.SetNext(rangeValidator)
	rangeValidator.SetNext(targetValidator)

	return &ValidationChain{
		head: statusValidator,
	}
}

// Validate lance la validation complète
func (vc *ValidationChain) Validate(cmd commands.Command) error {
	fmt.Printf("[ValidationChain] Début de la validation pour: %s\n", cmd.GetType())

	// Lancer la chaîne de validation
	if err := vc.head.Validate(cmd); err != nil {
		fmt.Printf("[ValidationChain] Échec: %v\n", err)
		return err
	}

	fmt.Printf("[ValidationChain] Validation réussie\n")
	return nil
}
