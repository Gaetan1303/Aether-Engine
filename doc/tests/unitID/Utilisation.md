// Création d'une nouvelle unité
func NewUnit(name string, stats Stats, position Position) *Unit {
    gen := domain.GetUnitIDGenerator()
    
    return &Unit{
        id:       gen.Generate(),
        name:     name,
        stats:    stats,
        position: position,
    }
}

// Création avec suffixe pour debugging
func NewHeroUnit(name string) *Unit {
    gen := domain.GetUnitIDGenerator()
    
    return &Unit{
        id:   gen.GenerateWithSuffix("hero"),
        name: name,
        // ...
    }
}

// Recherche d'unité dans une bataille
func (b *Battle) FindUnit(id UnitID) (*Unit, error) {
    for _, unit := range b.AllUnits() {
        if unit.ID().Equals(id) {
            return unit, nil
        }
    }
    return nil, fmt.Errorf("unit %s not found", id)
}

// Logging
func (u *Unit) Attack(target *Unit) {
    log.Printf("Unit %s attacks unit %s", u.id, target.id)
}