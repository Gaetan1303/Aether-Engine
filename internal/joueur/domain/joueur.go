package domain

import (
	"time"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// JoueurID identifiant unique d'un joueur
type JoueurID string

// Sexe enumération pour le sexe du personnage
type Sexe string

const (
	SexeMasculin Sexe = "masculin"
	SexeFeminin  Sexe = "feminin"
	SexeAutre    Sexe = "autre"
)

// TaillePersonnage enumération pour la taille
type TaillePersonnage string

const (
	TaillePetite  TaillePersonnage = "petite"
	TailleMoyenne TaillePersonnage = "moyenne"
	TailleGrande  TaillePersonnage = "grande"
)

// CouleurPersonnage structure pour les couleurs (RGB)
type CouleurPersonnage struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

// ApparencePhysique structure pour customiser l'apparence
type ApparencePhysique struct {
	Sexe           Sexe              `json:"sexe"`
	Taille         TaillePersonnage  `json:"taille"`
	CouleurPeau    CouleurPersonnage `json:"couleur_peau"`
	CouleurCheveux CouleurPersonnage `json:"couleur_cheveux"`
	CouleurYeux    CouleurPersonnage `json:"couleur_yeux"`
	StyleCheveux   string            `json:"style_cheveux"` // "courts", "longs", "chauve", etc.
	StyleBarbe     string            `json:"style_barbe"`   // "aucune", "courte", "longue", etc.
	Tatouages      []Tatouage        `json:"tatouages"`
	Maquillage     *Maquillage       `json:"maquillage,omitempty"`
	Cicatrices     []Cicatrice       `json:"cicatrices"`
}

// Tatouage structure pour les tatouages
type Tatouage struct {
	ID          string            `json:"id"`
	Nom         string            `json:"nom"`
	Position    string            `json:"position"` // "bras_gauche", "dos", "visage", etc.
	Couleur     CouleurPersonnage `json:"couleur"`
	Taille      string            `json:"taille"` // "petit", "moyen", "grand"
	Description string            `json:"description"`
}

// Maquillage structure pour le maquillage
type Maquillage struct {
	CouleurLevres    *CouleurPersonnage `json:"couleur_levres,omitempty"`
	CouleurJoues     *CouleurPersonnage `json:"couleur_joues,omitempty"`
	CouleurPaupieres *CouleurPersonnage `json:"couleur_paupieres,omitempty"`
	Intensite        string             `json:"intensite"` // "leger", "moyen", "fort"
}

// Cicatrice structure pour les cicatrices
type Cicatrice struct {
	ID          string `json:"id"`
	Position    string `json:"position"` // "visage", "bras", etc.
	Type        string `json:"type"`     // "fine", "large", "profonde"
	Taille      string `json:"taille"`   // "petite", "moyenne", "grande"
	Description string `json:"description"`
}

// Joueur entité principale du joueur avec customisation
type Joueur struct {
	ID        JoueurID          `json:"id"`
	Nom       string            `json:"nom"`
	Apparence ApparencePhysique `json:"apparence"`

	// Progression
	Niveau        int   `json:"niveau"`
	Experience    int64 `json:"experience"`
	ExperienceMax int64 `json:"experience_max"`

	// Job System
	JobActuel      JobID           `json:"job_actuel"`
	JobsDebloquees []JobID         `json:"jobs_debloquees"`
	ExperienceJobs map[JobID]int64 `json:"experience_jobs"`
	NiveauJobs     map[JobID]int   `json:"niveau_jobs"`

	// Système de slots compétences (comme FF Tactics A2)
	SlotsCompetences CompetenceSlots `json:"slots_competences"`

	// Stats de base (indépendantes du job)
	StatsBase shared.Stats `json:"stats_base"`

	// Métadonnées
	DateCreation time.Time     `json:"date_creation"`
	DernierLogin time.Time     `json:"dernier_login"`
	TempsJeu     time.Duration `json:"temps_jeu"`

	// Localisation actuelle (pour l'open world)
	ZoneActuelle string `json:"zone_actuelle"` // "ville_depart", "foret_nord", etc.
	SousZone     string `json:"sous_zone"`     // "quartier_marchand", "auberge", etc.

	// Inventaire et équipement (à développer plus tard)
	Inventaire []string          `json:"inventaire"` // IDs des objets
	Equipement map[string]string `json:"equipement"` // slot -> itemID
}

// NewJoueur crée un nouveau joueur avec apparence personnalisée
func NewJoueur(id JoueurID, nom string, apparence ApparencePhysique, jobInitial JobID) *Joueur {
	now := time.Now()

	// Stats de base selon l'apparence (taille influence certaines stats)
	statsBase, _ := shared.NewStats(10, 10, 10, 10, 10, 10, 10, 10, 1, 85) // Stats neutres

	// Bonus selon la taille
	switch apparence.Taille {
	case TaillePetite:
		statsBase, _ = shared.NewStats(8, 10, 10, 12, 10, 12, 10, 12, 1, 85) // -2 ATK, +2 SPD, +2 MATK
	case TailleGrande:
		statsBase, _ = shared.NewStats(12, 10, 10, 12, 12, 8, 12, 8, 1, 85) // +2 HP, +2 ATK, +2 DEF, -2 MATK, -2 SPD
	}

	return &Joueur{
		ID:               id,
		Nom:              nom,
		Apparence:        apparence,
		Niveau:           1,
		Experience:       0,
		ExperienceMax:    100,
		JobActuel:        jobInitial,
		JobsDebloquees:   []JobID{jobInitial},
		ExperienceJobs:   map[JobID]int64{jobInitial: 0},
		NiveauJobs:       map[JobID]int{jobInitial: 1},
		SlotsCompetences: NewCompetenceSlots(),
		StatsBase:        *statsBase,
		DateCreation:     now,
		DernierLogin:     now,
		TempsJeu:         0,
		ZoneActuelle:     "ville_depart",
		SousZone:         "place_centrale",
		Inventaire:       make([]string, 0),
		Equipement:       make(map[string]string),
	}
}

// GetStatsEffectives calcule les stats effectives (base + job + équipement)
func (j *Joueur) GetStatsEffectives() *shared.Stats {
	// Pour l'instant, retourne juste les stats de base
	// Plus tard : ajouter bonus du job + équipement
	return &j.StatsBase
}

// CanUseCompetence vérifie si le joueur peut utiliser une compétence
func (j *Joueur) CanUseCompetence(competenceID CompetenceID) bool {
	return j.SlotsCompetences.IsActive(competenceID)
}

// ChangerJob change le job actuel du joueur
func (j *Joueur) ChangerJob(nouveauJob JobID) error {
	// Vérifier si le job est débloqué
	for _, job := range j.JobsDebloquees {
		if job == nouveauJob {
			j.JobActuel = nouveauJob
			return nil
		}
	}

	return NewErrJobNonDebloque(nouveauJob)
}

// DebloquerJob débloque un nouveau job
func (j *Joueur) DebloquerJob(jobID JobID) {
	// Vérifier si déjà débloqué
	for _, job := range j.JobsDebloquees {
		if job == jobID {
			return
		}
	}

	j.JobsDebloquees = append(j.JobsDebloquees, jobID)
	j.ExperienceJobs[jobID] = 0
	j.NiveauJobs[jobID] = 1
}

// GagnerExperience fait gagner de l'expérience au joueur et à son job actuel
func (j *Joueur) GagnerExperience(xp int64) {
	j.Experience += xp
	j.ExperienceJobs[j.JobActuel] += xp

	// Vérifier montée de niveau
	j.verifierMonteeNiveau()
}

func (j *Joueur) verifierMonteeNiveau() {
	// Montée de niveau général
	for j.Experience >= j.ExperienceMax {
		j.Niveau++
		j.Experience -= j.ExperienceMax
		j.ExperienceMax = int64(j.Niveau * 100) // Formule simple

		// Améliorer stats de base
		j.StatsBase.ATK++
		j.StatsBase.HP += 5
	}

	// Montée de niveau du job
	jobXP := j.ExperienceJobs[j.JobActuel]
	jobNiveau := j.NiveauJobs[j.JobActuel]
	jobXPMax := int64(jobNiveau * 50)

	if jobXP >= jobXPMax {
		j.NiveauJobs[j.JobActuel]++
		j.ExperienceJobs[j.JobActuel] -= jobXPMax

		// Débloquer nouvelles compétences du job (à implémenter)
		// j.SlotsCompetences.DebloquerCompetences(j.JobActuel, j.NiveauJobs[j.JobActuel])
	}
}
