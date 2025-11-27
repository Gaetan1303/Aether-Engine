// Exemple : Infliger des dégâts avec formule
func (b *Battle) ResolveDamage(attacker, defender *Unit, skill Skill) {
    rawDamage := skill.BaseDamage + attacker.Stats.ATK() - defender.Stats.DEF()
    if rawDamage < 0 {
        rawDamage = 0
    }
    
    defender.Stats.TakeDamage(rawDamage)
    
    if defender.Stats.IsKO() {
        b.HandleUnitKO(defender)
    }
}

// Exemple : Lancer une compétence
func (u *Unit) CastSkill(skill Skill) error {
    if err := u.Stats.ConsumeMP(skill.MPCost); err != nil {
        return fmt.Errorf("cannot cast skill: %w", err)
    }
    // ... logique de la compétence
    return nil
}

// Exemple : Régénération en fin de tour
func (b *Battle) EndOfTurnRegeneration(unit *Unit) {
    mpRegen := unit.Stats.MAG() / 10 // 10% de MAG en MP
    unit.Stats.RestoreMP(mpRegen)
}