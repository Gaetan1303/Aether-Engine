package domain

// =============================================================================
// CONSTANTES DE COMBAT
// =============================================================================

// Configuration des équipes
const (
	// MinEquipesPourCombat est le nombre minimum d'équipes pour démarrer un combat
	MinEquipesPourCombat = 2

	// EquipeActiveVictoire est le nombre d'équipes actives restantes pour déclarer une victoire
	EquipeActiveVictoire = 1
)

// =============================================================================
// CONSTANTES D'ATTAQUE
// =============================================================================

// Portées de combat
const (
	// PorteeAttaqueMelee est la portée d'une attaque de mêlée basique
	PorteeAttaqueMelee = 1

	// TailleZoneEffetSingle est la taille d'une zone d'effet pour une cible unique
	TailleZoneEffetSingle = 1
)

// Propriétés d'attaque basique
const (
	// AttaqueBasiqueDegats sont les dégâts de base de l'attaque normale
	AttaqueBasiqueDegats = 10

	// AttaqueBasiqueScaling est le coefficient de scaling ATK pour l'attaque basique (50%)
	AttaqueBasiqueScaling = 0.5

	// AttaqueBasiqueCooldown est le cooldown de l'attaque basique (toujours disponible)
	AttaqueBasiqueCooldown = 1
)

// =============================================================================
// CONSTANTES DE COÛTS
// =============================================================================

// Coûts en ressources
const (
	// CoutMPNul représente une action sans coût en MP
	CoutMPNul = 0

	// CoutStaminaNul représente une action sans coût en Stamina
	CoutStaminaNul = 0
)

// =============================================================================
// CONSTANTES DE RÉGÉNÉRATION
// =============================================================================

// Taux de régénération par tour
const (
	// RegenMPDiviseur est le diviseur pour calculer la régénération de MP (10% = MP/10)
	RegenMPDiviseur = 10

	// RegenStaminaDiviseur est le diviseur pour calculer la régénération de Stamina (20% = Stamina/5)
	RegenStaminaDiviseur = 5
)

// =============================================================================
// CONSTANTES DE PROBABILITÉS
// =============================================================================

// Limites de probabilité de fuite
const (
	// FuiteProbabiliteMin est la probabilité minimale de fuite (10%)
	FuiteProbabiliteMin = 10.0

	// FuiteProbabiliteMax est la probabilité maximale de fuite (95%)
	FuiteProbabiliteMax = 95.0
)

// =============================================================================
// CONSTANTES DE VALIDATION
// =============================================================================

// Seuils de validation
const (
	// SeuilRessourceMinimum est le seuil minimum pour les ressources (HP, MP, Stamina)
	SeuilRessourceMinimum = 0
)
