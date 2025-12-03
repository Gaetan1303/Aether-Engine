package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewCompetence teste la création d'une nouvelle compétence
func TestNewCompetence(t *testing.T) {
	// Arrange
	id := domain.CompetenceID("fireball")
	nom := "Boule de Feu"
	description := "Lance une boule de feu enflammée"
	typeComp := domain.CompetenceMagie
	portee := 5
	zone := domain.ZoneEffet{}
	coutMP := 15
	coutStamina := 0
	cooldown := 3
	degatsBase := 30
	modificateur := 0.8
	cibles := domain.CibleEnnemis

	// Act
	comp := domain.NewCompetence(id, nom, description, typeComp, portee, zone, coutMP, coutStamina, cooldown, degatsBase, modificateur, cibles)

	// Assert
	assert.NotNil(t, comp, "La compétence ne devrait pas être nil")
	assert.Equal(t, id, comp.ID(), "L'ID devrait correspondre")
	assert.Equal(t, nom, comp.Nom(), "Le nom devrait correspondre")
	assert.Equal(t, description, comp.Description(), "La description devrait correspondre")
	assert.Equal(t, typeComp, comp.Type(), "Le type devrait correspondre")
	assert.Equal(t, portee, comp.Portee(), "La portée devrait correspondre")
	assert.Equal(t, coutMP, comp.CoutMP(), "Le coût MP devrait correspondre")
	assert.Equal(t, coutStamina, comp.CoutStamina(), "Le coût Stamina devrait correspondre")
	assert.Equal(t, cooldown, comp.Cooldown(), "Le cooldown devrait correspondre")
	assert.Equal(t, degatsBase, comp.DegatsBase(), "Les dégâts de base devraient correspondre")
	assert.Equal(t, modificateur, comp.Modificateur(), "Le modificateur devrait correspondre")
	assert.Equal(t, cibles, comp.Cibles(), "Le type de cibles devrait correspondre")
	assert.Equal(t, 0, comp.CooldownActuel(), "Le cooldown actuel devrait être 0 au départ")
}
