package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetObserverSubject teste la méthode SetObserverSubject()
func TestCombat_SetObserverSubject(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void avec paramètre nil, ne devrait pas crasher
	combat.SetObserverSubject(nil)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
