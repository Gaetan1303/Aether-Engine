package domain_test

import (
	"aether-engine-server/internal/combat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ========== Tests de StatusType ==========

func TestStatusType_String(t *testing.T) {
	assert.Equal(t, "Poison", domain.StatusPoison.String())
	assert.Equal(t, "Haste", domain.StatusHaste.String())
	assert.Equal(t, "Shield", domain.StatusShield.String())
}

func TestStatusType_IsDebuff(t *testing.T) {
	assert.True(t, domain.StatusPoison.IsDebuff())
	assert.True(t, domain.StatusSilence.IsDebuff())
	assert.True(t, domain.StatusSlow.IsDebuff())
	assert.False(t, domain.StatusHaste.IsDebuff())
	assert.False(t, domain.StatusRegen.IsDebuff())
}

func TestStatusType_IsBuff(t *testing.T) {
	assert.True(t, domain.StatusHaste.IsBuff())
	assert.True(t, domain.StatusRegen.IsBuff())
	assert.True(t, domain.StatusShield.IsBuff())
	assert.False(t, domain.StatusPoison.IsBuff())
	assert.False(t, domain.StatusSilence.IsBuff())
}

// ========== Tests de création Status ==========

func TestNewStatus_Valid(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")

	status, err := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPoison, status.Type())
	assert.Equal(t, 3, status.Duration())
	assert.Equal(t, 10, status.Intensity())
	assert.True(t, status.SourceID().Equals(sourceID))
}

func TestNewStatus_ZeroDuration(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")

	_, err := domain.NewStatus(domain.StatusPoison, 0, 10, sourceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

func TestNewStatus_NegativeDuration(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")

	_, err := domain.NewStatus(domain.StatusPoison, -1, 10, sourceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be negative")
}

func TestNewStatus_NegativeIntensity(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")

	_, err := domain.NewStatus(domain.StatusPoison, 3, -5, sourceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be negative")
}

func TestNewStatus_ZeroIntensity(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")

	// Valide pour les status sans intensité (Silence, Stun)
	status, err := domain.NewStatus(domain.StatusSilence, 2, 0, sourceID)

	assert.NoError(t, err)
	assert.Equal(t, 0, status.Intensity())
}

// ========== Tests de durée ==========

func TestStatus_DecrementDuration(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)

	status.DecrementDuration()
	assert.Equal(t, 2, status.Duration())

	status.DecrementDuration()
	assert.Equal(t, 1, status.Duration())

	status.DecrementDuration()
	assert.Equal(t, 0, status.Duration())
	assert.True(t, status.IsExpired())
}

func TestStatus_DecrementDuration_AtZero(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 1, 10, sourceID)

	status.DecrementDuration()
	status.DecrementDuration() // Déjà à 0

	assert.Equal(t, 0, status.Duration())
}

func TestStatus_IsExpired(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 1, 10, sourceID)

	assert.False(t, status.IsExpired())

	status.DecrementDuration()

	assert.True(t, status.IsExpired())
}

// ========== Tests de String() ==========

func TestStatus_String_WithIntensity(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 15, sourceID)

	str := status.String()

	assert.Contains(t, str, "Poison")
	assert.Contains(t, str, "duration=3")
	assert.Contains(t, str, "intensity=15")
}

func TestStatus_String_WithoutIntensity(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusSilence, 2, 0, sourceID)

	str := status.String()

	assert.Contains(t, str, "Silence")
	assert.Contains(t, str, "duration=2")
	assert.NotContains(t, str, "intensity")
}

// ========== Tests de Equals() ==========

func TestStatus_Equals_SameTypeAndSource(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status1, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	status2, _ := domain.NewStatus(domain.StatusPoison, 5, 15, sourceID)
	assert.True(t, status1.Equals(status2))
}

func TestStatus_Equals_DifferentType(t *testing.T) {
	sourceID, _ := domain.NewUnitID("unit_caster")
	status1, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	status2, _ := domain.NewStatus(domain.StatusHaste, 3, 0, sourceID)

	assert.False(t, status1.Equals(status2))
}

func TestStatus_Equals_DifferentSource(t *testing.T) {
	source1, _ := domain.NewUnitID("unit_caster1")
	source2, _ := domain.NewUnitID("unit_caster2")
	status1, _ := domain.NewStatus(domain.StatusPoison, 3, 10, source1)
	status2, _ := domain.NewStatus(domain.StatusPoison, 3, 10, source2)

	assert.False(t, status1.Equals(status2))
}

// ========== Tests de StatusCollection ==========

func TestStatusCollection_Add(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)

	err := collection.Add(status)

	assert.NoError(t, err)
	assert.True(t, collection.Has(domain.StatusPoison))
	assert.Equal(t, 1, collection.Count())
}

func TestStatusCollection_Add_ReplaceExisting(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	status1, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(status1)

	status2, _ := domain.NewStatus(domain.StatusPoison, 5, 20, sourceID)
	collection.Add(status2)

	// Doit avoir remplacé
	retrieved, _ := collection.Get(domain.StatusPoison)
	assert.Equal(t, 5, retrieved.Duration())
	assert.Equal(t, 20, retrieved.Intensity())
	assert.Equal(t, 1, collection.Count())
}

func TestStatusCollection_Add_ExpiredStatus(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 1, 10, sourceID)

	status.DecrementDuration() // Expire

	err := collection.Add(status)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestStatusCollection_Remove(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(status)

	collection.Remove(domain.StatusPoison)

	assert.False(t, collection.Has(domain.StatusPoison))
	assert.Equal(t, 0, collection.Count())
}

func TestStatusCollection_Get(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(status)

	retrieved, exists := collection.Get(domain.StatusPoison)

	assert.True(t, exists)
	assert.Equal(t, 3, retrieved.Duration())
}

func TestStatusCollection_Get_NotFound(t *testing.T) {
	collection := domain.NewStatusCollection()

	_, exists := collection.Get(domain.StatusPoison)

	assert.False(t, exists)
}

func TestStatusCollection_Has(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")
	status, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(status)

	assert.True(t, collection.Has(domain.StatusPoison))
	assert.False(t, collection.Has(domain.StatusHaste))
}

func TestStatusCollection_All(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	poison, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	haste, _ := domain.NewStatus(domain.StatusHaste, 2, 0, sourceID)

	collection.Add(poison)
	collection.Add(haste)

	all := collection.All()

	assert.Equal(t, 2, len(all))
}

func TestStatusCollection_DecrementAll(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	poison, _ := domain.NewStatus(domain.StatusPoison, 2, 10, sourceID)
	haste, _ := domain.NewStatus(domain.StatusHaste, 1, 0, sourceID)
	shield, _ := domain.NewStatus(domain.StatusShield, 3, 50, sourceID)

	collection.Add(poison)
	collection.Add(haste)
	collection.Add(shield)

	expired := collection.DecrementAll()

	// Haste doit avoir expiré (durée 1 → 0)
	assert.Equal(t, 1, len(expired))
	assert.Equal(t, domain.StatusHaste, expired[0])

	// Poison et Shield doivent être encore là
	assert.True(t, collection.Has(domain.StatusPoison))
	assert.True(t, collection.Has(domain.StatusShield))
	assert.False(t, collection.Has(domain.StatusHaste))

	// Vérifier les durées
	poisonStatus, _ := collection.Get(domain.StatusPoison)
	assert.Equal(t, 1, poisonStatus.Duration())

	shieldStatus, _ := collection.Get(domain.StatusShield)
	assert.Equal(t, 2, shieldStatus.Duration())
}

func TestStatusCollection_Clear(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	poison, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	haste, _ := domain.NewStatus(domain.StatusHaste, 2, 0, sourceID)

	collection.Add(poison)
	collection.Add(haste)

	collection.Clear()

	assert.Equal(t, 0, collection.Count())
	assert.False(t, collection.Has(domain.StatusPoison))
	assert.False(t, collection.Has(domain.StatusHaste))
}

func TestStatusCollection_Count(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	assert.Equal(t, 0, collection.Count())

	poison, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(poison)
	assert.Equal(t, 1, collection.Count())

	haste, _ := domain.NewStatus(domain.StatusHaste, 2, 0, sourceID)
	collection.Add(haste)
	assert.Equal(t, 2, collection.Count())

	collection.Remove(domain.StatusPoison)
	assert.Equal(t, 1, collection.Count())
}

// ========== Tests d'intégration ==========

func TestStatusCollection_MultipleStatusLifecycle(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_mage")

	// Tour 1 : Applique Poison (3 tours) et Haste (2 tours)
	poison, _ := domain.NewStatus(domain.StatusPoison, 3, 15, sourceID)
	haste, _ := domain.NewStatus(domain.StatusHaste, 2, 0, sourceID)
	collection.Add(poison)
	collection.Add(haste)

	assert.Equal(t, 2, collection.Count())

	// Tour 2 : Fin de tour
	expired := collection.DecrementAll()
	assert.Equal(t, 0, len(expired))
	assert.Equal(t, 2, collection.Count())

	// Tour 3 : Fin de tour
	expired = collection.DecrementAll()
	assert.Equal(t, 1, len(expired)) // Haste expire
	assert.Equal(t, domain.StatusHaste, expired[0])
	assert.Equal(t, 1, collection.Count())

	// Tour 4 : Fin de tour
	expired = collection.DecrementAll()
	assert.Equal(t, 1, len(expired)) // Poison expire
	assert.Equal(t, domain.StatusPoison, expired[0])
	assert.Equal(t, 0, collection.Count())
}

func TestStatusCollection_RefreshDuration(t *testing.T) {
	collection := domain.NewStatusCollection()
	sourceID, _ := domain.NewUnitID("unit_caster")

	// Applique Poison (3 tours)
	poison1, _ := domain.NewStatus(domain.StatusPoison, 3, 10, sourceID)
	collection.Add(poison1)

	// Fin de tour
	collection.DecrementAll()
	poisonStatus, _ := collection.Get(domain.StatusPoison)
	assert.Equal(t, 2, poisonStatus.Duration())

	// Réapplique Poison (5 tours) - refresh
	poison2, _ := domain.NewStatus(domain.StatusPoison, 5, 20, sourceID)
	collection.Add(poison2)

	// Doit avoir remplacé
	poisonStatus, _ = collection.Get(domain.StatusPoison)
	assert.Equal(t, 5, poisonStatus.Duration())
	assert.Equal(t, 20, poisonStatus.Intensity())
}
