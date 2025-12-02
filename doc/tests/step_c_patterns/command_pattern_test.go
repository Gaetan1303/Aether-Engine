package step_c_patterns_test

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Test de MoveCommand avec pathfinding
func TestMoveCommand_WithPathfinding(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 50)
	unitPos, _ := shared.NewPosition(0, 0)
	unit.DeplacerVers(unitPos)
	addUnitToCombat(combat, unit)

	factory := commands.NewCommandFactory(combat)

	// Act - Créer une commande de déplacement vers (5, 5)
	cmd, err := factory.CreateMoveCommand(unit, 5, 5)

	// Assert
	if err != nil {
		t.Fatalf("Erreur lors de la création de MoveCommand: %v", err)
	}

	// Validate - Vérifier que le chemin est valide
	err = cmd.Validate()
	if err != nil {
		t.Errorf("MoveCommand devrait être valide: %v", err)
	}

	// Execute
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("Erreur lors de l'exécution de MoveCommand: %v", err)
	}

	// Vérifier que l'unité s'est déplacée
	if result != nil && result.Success {
		newPos := unit.Position()
		if newPos.X() != 5 || newPos.Y() != 5 {
			t.Errorf("Position attendue (5,5), obtenue (%d,%d)", newPos.X(), newPos.Y())
		}
	}
}

// Test de MoveCommand avec obstruction
func TestMoveCommand_ObstructedPath(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 50)
	unitPos, _ := shared.NewPosition(0, 0)
	unit.DeplacerVers(unitPos)

	// Ajouter une unité bloquante
	blocker := createTestUnit("U2", 50)
	blockerPos, _ := shared.NewPosition(2, 2)
	blocker.DeplacerVers(blockerPos)

	addUnitToCombat(combat, unit)
	addUnitToCombat(combat, blocker)

	factory := commands.NewCommandFactory(combat)

	// Act - Essayer de se déplacer à travers le bloqueur
	cmd, err := factory.CreateMoveCommand(unit, 2, 2)

	// Assert
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	err = cmd.Validate()
	if err == nil {
		t.Errorf("MoveCommand devrait échouer (position occupée)")
	}
}

// Test de AttackCommand avec portée valide
func TestAttackCommand_ValidRange(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	attacker := createTestUnit("A1", 50)
	attackerPos, _ := shared.NewPosition(0, 0)
	attacker.DeplacerVers(attackerPos)

	target := createTestUnit("E1", 50)
	targetPos, _ := shared.NewPosition(1, 0)
	target.DeplacerVers(targetPos) // Adjacent
	target.SetHP(100)

	addUnitToCombat(combat, attacker)
	addUnitToCombat(combat, target)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateAttackCommand(attacker, target.ID())
	if err != nil {
		t.Fatalf("Erreur lors de la création de AttackCommand: %v", err)
	}

	// Validate
	err = cmd.Validate()
	if err != nil {
		t.Errorf("AttackCommand devrait être valide (cible adjacente): %v", err)
	}

	// Execute
	initialHP := target.StatsActuelles().HP
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("Erreur lors de l'exécution: %v", err)
	}

	// Assert
	if result != nil && result.Success {
		if target.StatsActuelles().HP >= initialHP {
			t.Errorf("La cible devrait avoir perdu des HP")
		}
		if result.DamageDealt <= 0 {
			t.Errorf("DamageDealt devrait être > 0, obtenu: %d", result.DamageDealt)
		}
	}
}

// Test de AttackCommand hors portée
func TestAttackCommand_OutOfRange(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	attacker := createTestUnit("A1", 50)
	attackerPos, _ := shared.NewPosition(0, 0)
	attacker.DeplacerVers(attackerPos)

	target := createTestUnit("E1", 50)
	targetPos, _ := shared.NewPosition(9, 9)
	target.DeplacerVers(targetPos) // Trop loin

	addUnitToCombat(combat, attacker)
	addUnitToCombat(combat, target)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateAttackCommand(attacker, target.ID())
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	// Validate
	err = cmd.Validate()

	// Assert
	if err == nil {
		t.Errorf("AttackCommand devrait échouer (hors portée)")
	}
}

// Test de SkillCommand avec MP suffisants
func TestSkillCommand_SufficientMP(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	caster := createTestUnit("U1", 50)
	caster.SetMP(100)
	casterPos, _ := shared.NewPosition(0, 0)
	caster.DeplacerVers(casterPos)

	target := createTestUnit("U2", 50)
	targetPos, _ := shared.NewPosition(1, 1)
	target.DeplacerVers(targetPos)
	target.SetHP(100)

	addUnitToCombat(combat, caster)
	addUnitToCombat(combat, target)

	// Créer une compétence de test (ex: attaque magique coûtant 30 MP)
	skill := createTestSkill("fireball", 30, domain.CompetenceMagie)
	caster.AjouterCompetence(skill)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateSkillCommand(caster, "fireball", []domain.UnitID{target.ID()})
	if err != nil {
		t.Fatalf("Erreur lors de la création de SkillCommand: %v", err)
	}

	// Validate
	err = cmd.Validate()
	if err != nil {
		t.Errorf("SkillCommand devrait être valide (MP suffisants): %v", err)
	}

	// Execute
	initialMP := caster.StatsActuelles().MP
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("Erreur lors de l'exécution: %v", err)
	}

	// Assert
	if result != nil && result.Success {
		if caster.StatsActuelles().MP != initialMP-30 {
			t.Errorf("MP attendus: %d, obtenus: %d", initialMP-30, caster.StatsActuelles().MP)
		}
		if result.CostMP != 30 {
			t.Errorf("CostMP attendu: 30, obtenu: %d", result.CostMP)
		}
	}
}

// Test de SkillCommand avec MP insuffisants
func TestSkillCommand_InsufficientMP(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	caster := createTestUnit("U1", 50)
	caster.SetMP(10) // Pas assez pour la compétence

	target := createTestUnit("U2", 50)
	addUnitToCombat(combat, caster)
	addUnitToCombat(combat, target)

	skill := createTestSkill("fireball", 30, domain.CompetenceMagie)
	caster.AjouterCompetence(skill)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateSkillCommand(caster, "fireball", []domain.UnitID{target.ID()})
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	// Validate
	err = cmd.Validate()

	// Assert
	if err == nil {
		t.Errorf("SkillCommand devrait échouer (MP insuffisants)")
	}
}

// Test de SkillCommand avec cooldown actif
func TestSkillCommand_CooldownActive(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	caster := createTestUnit("U1", 50)
	caster.SetMP(100)

	target := createTestUnit("U2", 50)
	addUnitToCombat(combat, caster)
	addUnitToCombat(combat, target)

	skill := createTestSkill("ultimate", 50, domain.CompetenceMagie)
	caster.AjouterCompetence(skill)

	// Activer le cooldown
	caster.ActiverCooldown(skill.ID(), 3)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateSkillCommand(caster, string(skill.ID()), []domain.UnitID{target.ID()})
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	// Validate
	err = cmd.Validate()

	// Assert
	if err == nil {
		t.Errorf("SkillCommand devrait échouer (cooldown actif)")
	}
}

// Test de SkillCommand multi-cibles (AoE)
func TestSkillCommand_MultiTarget(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	caster := createTestUnit("U1", 50)
	caster.SetMP(100)
	casterPos, _ := shared.NewPosition(5, 5)
	caster.DeplacerVers(casterPos)

	target1 := createTestUnit("E1", 50)
	target1Pos, _ := shared.NewPosition(6, 5)
	target1.DeplacerVers(target1Pos)
	target1.SetHP(100)

	target2 := createTestUnit("E2", 50)
	target2Pos, _ := shared.NewPosition(6, 6)
	target2.DeplacerVers(target2Pos)
	target2.SetHP(100)

	addUnitToCombat(combat, caster)
	addUnitToCombat(combat, target1)
	addUnitToCombat(combat, target2)

	skill := createTestSkill("aoe_spell", 40, domain.CompetenceMagie)
	caster.AjouterCompetence(skill)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateSkillCommand(caster, string(skill.ID()),
		[]domain.UnitID{target1.ID(), target2.ID()})
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	// Execute
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("Erreur lors de l'exécution multi-cibles: %v", err)
	}

	// Assert
	if result != nil && result.Success {
		if len(result.Effects) != 2 {
			t.Errorf("Nombre d'effets attendu: 2, obtenu: %d", len(result.Effects))
		}

		// Vérifier que les deux cibles ont pris des dégâts
		if target1.StatsActuelles().HP >= 100 || target2.StatsActuelles().HP >= 100 {
			t.Errorf("Les deux cibles devraient avoir perdu des HP")
		}
	}
}

// Test de ItemCommand - Potion
func TestItemCommand_Potion(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	user := createTestUnit("U1", 50)
	user.SetHP(50) // Blessé
	addUnitToCombat(combat, user)

	// Ajouter une potion à l'inventaire
	combat.AjouterObjet("potion", 1)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateItemCommand(user, "potion", user.ID())
	if err != nil {
		t.Fatalf("Erreur lors de la création de ItemCommand: %v", err)
	}

	// Validate
	err = cmd.Validate()
	if err != nil {
		t.Errorf("ItemCommand devrait être valide: %v", err)
	}

	// Execute
	initialHP := user.StatsActuelles().HP
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("Erreur lors de l'exécution: %v", err)
	}

	// Assert
	if result != nil && result.Success {
		if user.StatsActuelles().HP <= initialHP {
			t.Errorf("HP devraient augmenter après utilisation de potion")
		}
		if result.HealingDone <= 0 {
			t.Errorf("HealingDone devrait être > 0")
		}

		// Vérifier que l'objet a été consommé
		if combat.ObtenirQuantiteObjet("potion") != 0 {
			t.Errorf("La potion devrait être consommée")
		}
	}
}

// Test de ItemCommand - Objet non disponible
func TestItemCommand_ItemNotInInventory(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	user := createTestUnit("U1", 50)
	addUnitToCombat(combat, user)

	factory := commands.NewCommandFactory(combat)

	// Act - Essayer d'utiliser un objet qui n'existe pas
	cmd, err := factory.CreateItemCommand(user, "nonexistent", user.ID())
	if err != nil {
		t.Fatalf("Erreur lors de la création: %v", err)
	}

	// Validate
	err = cmd.Validate()

	// Assert
	if err == nil {
		t.Errorf("ItemCommand devrait échouer (objet non disponible)")
	}
}

// Test de FleeCommand avec probabilité de réussite
func TestFleeCommand_Success(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 100) // Speed élevé = plus de chances
	addUnitToCombat(combat, unit)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateFleeCommand(unit)
	if err != nil {
		t.Fatalf("Erreur lors de la création de FleeCommand: %v", err)
	}

	// Execute plusieurs fois pour tester la probabilité
	successCount := 0
	for i := 0; i < 100; i++ {
		result, err := cmd.Execute()
		if err == nil && result != nil && result.Success {
			successCount++
		}

		// Recréer la commande pour chaque tentative
		cmd, _ = factory.CreateFleeCommand(unit)
	}

	// Assert - Avec speed 100, devrait réussir au moins 50% du temps
	if successCount < 30 {
		t.Errorf("FleeCommand devrait réussir plus souvent avec speed élevé, succès: %d/100", successCount)
	}
}

// Test de WaitCommand
func TestWaitCommand_Execute(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 50)
	addUnitToCombat(combat, unit)

	factory := commands.NewCommandFactory(combat)

	// Act
	cmd, err := factory.CreateWaitCommand(unit)
	if err != nil {
		t.Fatalf("Erreur lors de la création de WaitCommand: %v", err)
	}

	// Execute
	result, err := cmd.Execute()
	if err != nil {
		t.Errorf("WaitCommand ne devrait jamais échouer: %v", err)
	}

	// Assert
	if result == nil || !result.Success {
		t.Errorf("WaitCommand devrait toujours réussir")
	}
	if result.Message == "" {
		t.Errorf("WaitCommand devrait retourner un message")
	}
}

// Test de CommandInvoker - History
func TestCommandInvoker_History(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 50)
	addUnitToCombat(combat, unit)

	factory := commands.NewCommandFactory(combat)
	invoker := commands.NewCommandInvoker(100) // Max 100 commandes dans l'historique

	// Act - Exécuter plusieurs commandes
	cmd1, _ := factory.CreateWaitCommand(unit)
	cmd2, _ := factory.CreateWaitCommand(unit)
	cmd3, _ := factory.CreateWaitCommand(unit)

	invoker.Execute(cmd1)
	invoker.Execute(cmd2)
	invoker.Execute(cmd3)

	// Assert
	history := invoker.GetHistory()
	if len(history) != 3 {
		t.Errorf("L'historique devrait contenir 3 commandes, obtenu: %d", len(history))
	}
}

// Test de Command Rollback
func TestCommand_Rollback(t *testing.T) {
	// Arrange
	combat := createTestCombat()
	unit := createTestUnit("U1", 50)
	initialPos, _ := shared.NewPosition(0, 0)
	unit.DeplacerVers(initialPos)
	addUnitToCombat(combat, unit)

	factory := commands.NewCommandFactory(combat)

	// Act - Exécuter puis Rollback
	cmd, _ := factory.CreateMoveCommand(unit, 5, 5)
	result, err := cmd.Execute()
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution: %v", err)
	}

	if result != nil && result.Success {
		// Rollback
		err = cmd.Rollback()
		if err != nil {
			t.Errorf("Erreur lors du rollback: %v", err)
		}

		// Assert - L'unité devrait être revenue à sa position initiale
		pos := unit.Position()
		if pos.X() != initialPos.X() || pos.Y() != initialPos.Y() {
			t.Errorf("Position après rollback: (%d,%d), attendue: (%d,%d)",
				pos.X(), pos.Y(), initialPos.X(), initialPos.Y())
		}
	}
}
