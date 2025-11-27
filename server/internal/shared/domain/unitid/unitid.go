package unitid

import (
	"fmt"
)

// UnitID représente un identifiant unique d'unité, typiquement un entier positif non nul.
type UnitID struct {
	value int
}

// New crée un UnitID valide (doit être > 0).
func New(id int) (UnitID, error) {
	if id <= 0 {
		return UnitID{}, fmt.Errorf("UnitID doit être strictement positif (reçu: %d)", id)
	}
	return UnitID{value: id}, nil
}

// NewUnchecked crée un UnitID sans validation (usage interne/tests).
func NewUnchecked(id int) UnitID {
	return UnitID{value: id}
}

// Value retourne la valeur brute de l'identifiant.
func (u UnitID) Value() int {
	return u.value
}

// Equals compare deux UnitID pour l'égalité stricte.
func (u UnitID) Equals(other UnitID) bool {
	return u.value == other.value
}
