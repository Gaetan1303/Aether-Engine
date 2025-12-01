package domain_test

import (
	"sync"
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestSingletonPattern_UniqueInstance teste que GetIDGenerator retourne toujours la même instance
func TestSingletonPattern_UniqueInstance(t *testing.T) {
	// Act: Obtenir l'instance plusieurs fois
	gen1 := shared.GetIDGenerator()
	gen2 := shared.GetIDGenerator()
	gen3 := shared.GetIDGenerator()

	// Assert: Toutes les références pointent vers la même instance
	assert.Same(t, gen1, gen2, "gen1 et gen2 doivent être la même instance")
	assert.Same(t, gen2, gen3, "gen2 et gen3 doivent être la même instance")
	assert.Same(t, gen1, gen3, "gen1 et gen3 doivent être la même instance")
}

// TestSingletonPattern_ThreadSafety teste la thread-safety du Singleton
func TestSingletonPattern_ThreadSafety(t *testing.T) {
	// Reset pour test propre
	shared.GetIDGenerator().Reset()

	const numGoroutines = 100
	var wg sync.WaitGroup
	instances := make([]*shared.IDGenerator, numGoroutines)

	// Act: Accéder au Singleton depuis 100 goroutines simultanément
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			instances[index] = shared.GetIDGenerator()
		}(i)
	}

	wg.Wait()

	// Assert: Toutes les goroutines ont obtenu la même instance
	firstInstance := instances[0]
	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, firstInstance, instances[i], "Toutes les instances doivent être identiques")
	}
}

// TestIDGenerator_UniqueCombatIDs teste la génération d'IDs uniques pour les combats
func TestIDGenerator_UniqueCombatIDs(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer 100 IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := gen.NewCombatID()
		ids[id] = true
	}

	// Assert: Tous les IDs sont uniques
	assert.Equal(t, 100, len(ids), "Tous les IDs doivent être uniques")
}

// TestIDGenerator_UniqueUnitIDs teste la génération d'IDs uniques pour les unités
func TestIDGenerator_UniqueUnitIDs(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer 100 IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := gen.NewUnitID()
		ids[id] = true
	}

	// Assert: Tous les IDs sont uniques
	assert.Equal(t, 100, len(ids), "Tous les IDs doivent être uniques")
}

// TestIDGenerator_UniqueTeamIDs teste la génération d'IDs uniques pour les équipes
func TestIDGenerator_UniqueTeamIDs(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer 50 IDs
	ids := make(map[string]bool)
	for i := 0; i < 50; i++ {
		id := gen.NewTeamID()
		ids[id] = true
	}

	// Assert: Tous les IDs sont uniques
	assert.Equal(t, 50, len(ids), "Tous les IDs doivent être uniques")
}

// TestIDGenerator_MixedIDs teste la génération de différents types d'IDs
func TestIDGenerator_MixedIDs(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer différents types d'IDs
	combatID := gen.NewCombatID()
	unitID1 := gen.NewUnitID()
	unitID2 := gen.NewUnitID()
	teamID := gen.NewTeamID()
	skillID := gen.NewCompetenceID()
	itemID := gen.NewObjetID()
	eventID := gen.NewEventID()

	// Assert: Tous sont uniques
	ids := []string{combatID, unitID1, unitID2, teamID, skillID, itemID, eventID}
	uniqueIDs := make(map[string]bool)
	for _, id := range ids {
		uniqueIDs[id] = true
	}
	assert.Equal(t, len(ids), len(uniqueIDs), "Tous les IDs doivent être uniques")

	// Assert: Les préfixes sont corrects
	assert.Contains(t, combatID, "combat_", "Combat ID doit avoir le bon préfixe")
	assert.Contains(t, unitID1, "unit_", "Unit ID doit avoir le bon préfixe")
	assert.Contains(t, teamID, "team_", "Team ID doit avoir le bon préfixe")
	assert.Contains(t, skillID, "skill_", "Skill ID doit avoir le bon préfixe")
	assert.Contains(t, itemID, "item_", "Item ID doit avoir le bon préfixe")
	assert.Contains(t, eventID, "event_", "Event ID doit avoir le bon préfixe")
}

// TestIDGenerator_ConcurrentGeneration teste la génération concurrente d'IDs
func TestIDGenerator_ConcurrentGeneration(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	const numGoroutines = 50
	const idsPerGoroutine = 20
	var wg sync.WaitGroup

	idChan := make(chan string, numGoroutines*idsPerGoroutine)

	// Act: Générer des IDs depuis plusieurs goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				idChan <- gen.NewCombatID()
			}
		}()
	}

	wg.Wait()
	close(idChan)

	// Assert: Collecter tous les IDs et vérifier l'unicité
	ids := make(map[string]bool)
	for id := range idChan {
		ids[id] = true
	}

	expectedCount := numGoroutines * idsPerGoroutine
	assert.Equal(t, expectedCount, len(ids), "Tous les IDs générés en parallèle doivent être uniques")
}

// TestIDGenerator_GenerateUUID teste la génération d'UUIDs
func TestIDGenerator_GenerateUUID(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()

	// Act: Générer plusieurs UUIDs
	uuid1 := gen.GenerateUUID()
	uuid2 := gen.GenerateUUID()
	uuid3 := gen.GenerateUUID()

	// Assert: Les UUIDs sont uniques
	assert.NotEqual(t, uuid1, uuid2)
	assert.NotEqual(t, uuid2, uuid3)
	assert.NotEqual(t, uuid1, uuid3)

	// Assert: Format UUID valide (contient des tirets)
	assert.Contains(t, uuid1, "-", "UUID doit contenir des tirets")
	assert.Greater(t, len(uuid1), 10, "UUID doit avoir une longueur suffisante")
}

// TestIDGenerator_GetStats teste les statistiques du générateur
func TestIDGenerator_GetStats(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer quelques IDs
	gen.NewCombatID()
	gen.NewUnitID()
	gen.NewTeamID()

	stats := gen.GetStats()

	// Assert: Les stats sont cohérentes
	assert.Equal(t, uint64(3), stats.TotalGenerated, "Compteur doit être à 3")
	assert.NotEmpty(t, stats.MachineID, "MachineID ne doit pas être vide")
	assert.GreaterOrEqual(t, stats.Uptime.Milliseconds(), int64(0), "Uptime doit être >= 0")
}

// TestIDGenerator_ParseIDType teste l'extraction du type d'ID
func TestIDGenerator_ParseIDType(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	combatID := gen.NewCombatID()
	unitID := gen.NewUnitID()
	teamID := gen.NewTeamID()

	// Act & Assert
	assert.Equal(t, "combat", shared.ParseIDType(combatID))
	assert.Equal(t, "unit", shared.ParseIDType(unitID))
	assert.Equal(t, "team", shared.ParseIDType(teamID))
}

// TestIDGenerator_ValidateIDFormat teste la validation du format d'ID
func TestIDGenerator_ValidateIDFormat(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	validID := gen.NewCombatID()

	// Act & Assert: IDs valides
	assert.True(t, shared.ValidateIDFormat(validID), "ID généré doit être valide")
	assert.True(t, shared.ValidateIDFormat("combat_1234567890_abcd_0001"), "Format valide")

	// Act & Assert: IDs invalides
	assert.False(t, shared.ValidateIDFormat(""), "ID vide invalide")
	assert.False(t, shared.ValidateIDFormat("combat"), "ID sans segments invalide")
	assert.False(t, shared.ValidateIDFormat("invalid_format"), "Format incomplet invalide")
}

// TestIDGenerator_TypedHelpers teste les helpers typés
func TestIDGenerator_TypedHelpers(t *testing.T) {
	// Act: Utiliser les helpers
	combatID := shared.NewCombatIDTyped()
	unitID := shared.NewUnitIDTyped()
	teamID := shared.NewTeamIDTyped()
	skillID := shared.NewCompetenceIDTyped()
	itemID := shared.NewObjetIDTyped()

	// Assert: Tous sont générés et uniques
	assert.NotEmpty(t, combatID)
	assert.NotEmpty(t, unitID)
	assert.NotEmpty(t, teamID)
	assert.NotEmpty(t, skillID)
	assert.NotEmpty(t, itemID)

	// Assert: Préfixes corrects
	assert.Equal(t, "combat", shared.ParseIDType(combatID))
	assert.Equal(t, "unit", shared.ParseIDType(unitID))
	assert.Equal(t, "team", shared.ParseIDType(teamID))
	assert.Equal(t, "skill", shared.ParseIDType(skillID))
	assert.Equal(t, "item", shared.ParseIDType(itemID))
}

// TestIDGenerator_HighVolume teste la génération à haut volume
func TestIDGenerator_HighVolume(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()
	gen.Reset()

	// Act: Générer 10000 IDs
	ids := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		id := gen.NewCombatID()
		ids[id] = true
	}

	// Assert: Tous uniques
	assert.Equal(t, 10000, len(ids), "Tous les 10000 IDs doivent être uniques")

	// Assert: Stats cohérentes
	stats := gen.GetStats()
	assert.Equal(t, uint64(10000), stats.TotalGenerated, "Compteur doit être à 10000")
}

// TestIDGenerator_Reset teste la réinitialisation (pour tests uniquement)
func TestIDGenerator_Reset(t *testing.T) {
	// Arrange
	gen := shared.GetIDGenerator()

	// Sauvegarder le compteur avant
	statsBefore := gen.GetStats()
	countBefore := statsBefore.TotalGenerated

	// Générer quelques IDs
	gen.NewCombatID()
	gen.NewCombatID()
	gen.NewCombatID()

	statsAfterGen := gen.GetStats()
	assert.Equal(t, countBefore+3, statsAfterGen.TotalGenerated, "Compteur doit avoir augmenté de 3")

	// Act: Reset
	gen.Reset()

	// Assert: Compteur remis à zéro
	statsAfterReset := gen.GetStats()
	assert.Equal(t, uint64(0), statsAfterReset.TotalGenerated, "Compteur doit être remis à 0")
} // BenchmarkIDGenerator_NewCombatID benchmark la génération d'IDs
func BenchmarkIDGenerator_NewCombatID(b *testing.B) {
	gen := shared.GetIDGenerator()
	gen.Reset()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.NewCombatID()
	}
}

// BenchmarkIDGenerator_Concurrent benchmark la génération concurrente
func BenchmarkIDGenerator_Concurrent(b *testing.B) {
	gen := shared.GetIDGenerator()
	gen.Reset()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = gen.NewCombatID()
		}
	})
}
