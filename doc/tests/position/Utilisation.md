// Exemple : Vérifier si une unité peut attaquer une cible
func (u *Unit) CanAttack(target *Unit, attackRange int) bool {
    distance := u.Position.ManhattanDistance(target.Position)
    return distance <= attackRange
}

// Exemple : Déplacement valide
func (b *Battle) MoveUnit(unitID UnitID, newPos Position) error {
    if !newPos.IsInBounds(b.Grid.Width, b.Grid.Height) {
        return errors.New("position out of bounds")
    }
    // ... reste de la logique
}