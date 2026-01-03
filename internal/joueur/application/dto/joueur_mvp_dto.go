package dto

// JoueurCreateMVPDTO version simplifiée pour MVP
// Retire les éléments complexes (tatouages, cicatrices, maquillage, slots multiples)
type JoueurCreateMVPDTO struct {
	Nom       string              `json:"nom" binding:"required,min=3,max=20"`
	Apparence ApparenceSimpleDTO  `json:"apparence" binding:"required"`
	JobInitial string             `json:"job_initial" binding:"required,oneof=guerrier mage archer voleur clerc"`
}

// ApparenceSimpleDTO version MVP - uniquement l'essentiel
type ApparenceSimpleDTO struct {
	Sexe           string     `json:"sexe" binding:"required,oneof=masculin feminin autre"`
	Taille         string     `json:"taille" binding:"required,oneof=petite moyenne grande"`
	CouleurPeau    CouleurDTO `json:"couleur_peau" binding:"required"`
	CouleurCheveux CouleurDTO `json:"couleur_cheveux" binding:"required"`
	CouleurYeux    CouleurDTO `json:"couleur_yeux" binding:"required"`
	StyleCheveux   string     `json:"style_cheveux" binding:"required,oneof=chauve courts longs"`
	StyleBarbe     string     `json:"style_barbe" binding:"omitempty,oneof=aucune courte longue"`
}

// CompetencesSimpleMVPDTO version MVP - seulement slots actifs
type CompetencesSimpleMVPDTO struct {
	SlotsActifs            []SlotCompetenceDTO `json:"slots_actifs"`
	MaxSlotsActifs         int                 `json:"max_slots_actifs"`
	CompetencesDisponibles []string            `json:"competences_disponibles"`
}

// CustomisationOptionsMVP options simplifiées pour MVP
type CustomisationOptionsMVP struct {
	Sexes                  []string                        `json:"sexes"`
	Tailles                []string                        `json:"tailles"`
	StylesCheveux          []string                        `json:"styles_cheveux"`
	StylesBarbe            []string                        `json:"styles_barbe"`
	CouleursPredefinies    []CouleurDTO                    `json:"couleurs_predefinies"`
	JobsDepart             []JobBasiqueDTO                 `json:"jobs_depart"`
}

// JobBasiqueDTO informations de base sur un job pour MVP
type JobBasiqueDTO struct {
	ID          string `json:"id"`
	Nom         string `json:"nom"`
	Description string `json:"description"`
	StatsBonus  string `json:"stats_bonus"`
}