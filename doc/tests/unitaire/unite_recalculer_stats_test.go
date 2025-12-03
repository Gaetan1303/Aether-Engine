package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_RecalculerStats teste la méthode RecalculerStats()
func TestUnite_RecalculerStats(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Oracle", "team-1", 5, 5)

	// Act
	unite.RecalculerStats()

	// Assert
	stats := unite.StatsActuelles()
	assert.NotNil(t, stats, "Les stats recalculées ne devraient pas être nil")
	assert.Equal(t, 100, stats.HP, "HP devrait être recalculé à la valeur de base")
}
