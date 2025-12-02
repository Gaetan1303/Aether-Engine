package combatinterfaces

// CombatRepository interface pour éviter les cycles d'imports
// Les commands/observers/validators utilisent cette interface au lieu de *domain.Combat
type CombatRepository interface {
	// Méthodes de base
	ID() string
	Grille() GridInterface
	GetDamageCalculator() DamageCalculatorInterface
	GetTimestamp() int64

	// Recherche d'unités
	TrouverUnite(unitID interface{}) UniteInterface
	ObtenirPositionsOccupees(exclusionID interface{}) map[string]bool
	ObtenirEnnemis(teamID interface{}) []UniteInterface

	// Gestion fuite
	FuiteAutorisee() bool
	MarquerEquipeFuite(teamID interface{})
	AnnulerFuite(teamID interface{})

	// Victoire/Défaite
	VerifierConditionsVictoire() string
	ObtenirResultat() string
	DistribuerRecompenses()

	// Inventaire
	PossedeObjet(itemID string) bool
	ObtenirQuantiteObjet(itemID string) int
	ConsommerObjet(itemID string, quantite int)
	AjouterObjet(itemID string, quantite int)
	ObtenirObjet(itemID string) ItemInterface
}

// UniteInterface interface pour les unités
type UniteInterface interface {
	ID() interface{}
	Nom() string
	TeamID() interface{}
	Position() PositionInterface
	Stats() StatsInterface
	StatsActuelles() StatsInterface
	HPActuels() int

	// Actions
	PeutAgir() bool
	EstEliminee() bool
	EstBloqueDeplacement() bool
	EstSilence() bool
	EstStun() bool
	EstRoot() bool
	EstEmpoisonne() bool
	EstIA() bool

	// Modifications
	RecevoirDegats(degats int)
	Soigner(soin int)
	RestaurerMP(mp int)
	ConsommerMP(mp int)
	ConsommerStamina(stamina int)
	DeplacerVers(position PositionInterface)
	AjouterStatut(statut interface{})
	RetirerStatut(statusType string)
	Ressusciter(hp int)
	AppliquerBuff(stat string, value int, duree int)

	// Compétences
	ObtenirCompetence(skillID string) CompetenceInterface
	ObtenirCompetenceParDefaut() CompetenceInterface
	PossedeCompetence(skillID string) bool
	SkillEstPret(skillID string) bool
	ActiverCooldown(skillID string, duree int)

	// Statuts
	NouveauTour()
	TraiterStatuts() []interface{}

	// IA
	IAChoisirAction(combat interface{})
}

// GridInterface interface pour la grille
type GridInterface interface {
	EstDansLimites(position PositionInterface) bool
	Position(x, y int) PositionInterface
}

// PositionInterface interface pour les positions
type PositionInterface interface {
	X() int
	Y() int
	Equals(other PositionInterface) bool
}

// StatsInterface interface pour les statistiques
type StatsInterface interface {
	HP() int
	MP() int
	ATK() int
	DEF() int
	MATK() int
	MDEF() int
	SPD() int
	MOV() int
	Stamina() int
}

// CompetenceInterface interface pour les compétences
type CompetenceInterface interface {
	ID() string
	Nom() string
	Type() interface{}
	CoutMP() int
	CoutStamina() int
	Portee() int
	Cooldown() int
	Duree() int
	DegatsBase() int
	Modificateur() float64
	Effets() []interface{}
	StatBuff() string
	BuffValue() int
	CreerStatut() interface{}
}

// ItemInterface interface pour les objets
type ItemInterface interface {
	ID() string
	Nom() string
	Type() interface{}
	Portee() int
	EffectValue() int
}

// DamageCalculatorInterface interface pour le calcul de dégâts
type DamageCalculatorInterface interface {
	CalculerDegats(atk, def int, competence CompetenceInterface) int
	CalculerDegatsSkill(acteur, cible UniteInterface, competence CompetenceInterface) int
	CalculerSoins(acteur UniteInterface, competence CompetenceInterface) int
}
