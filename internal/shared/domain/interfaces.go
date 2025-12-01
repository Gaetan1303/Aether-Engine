package domain

// StatsModifiable représente une entité dont les stats peuvent être modifiées
// Interface Segregation Principle (SOLID) - Interface spécifique pour la modification des stats
type StatsModifiable interface {
	AppliquerModificateurStat(modificateur *ModificateurStat)
	RetirerModificateurStat(modificateur *ModificateurStat)
	RecevoirDegats(degats int)
	RecevoirSoin(soin int)
}

// Combattant représente une entité capable de combattre
// Interface Segregation Principle (SOLID) - Interface pour les comportements de combat
type Combattant interface {
	PeutAgir() bool
	PeutSeDeplacer() bool
	EstEliminee() bool
}
