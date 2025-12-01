package domain

// ModificateurStat représente un modificateur temporaire de stat
type ModificateurStat struct {
	Stat   string // "ATK", "DEF", etc.
	Valeur int    // Peut être négatif
}

// EffetStatut représente l'effet d'un statut appliqué
type EffetStatut struct {
	Type   TypeStatut
	Valeur int // Dégâts/soin périodique
}

// ObjetID est l'identifiant d'un objet
type ObjetID string
