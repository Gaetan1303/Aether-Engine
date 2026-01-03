package domain

// JobID identifiant d'un job/classe
type JobID string

// Jobs de départ disponibles
const (
	JobGuerrier JobID = "guerrier"
	JobMage     JobID = "mage"
	JobArcher   JobID = "archer"
	JobVoleur   JobID = "voleur"
	JobClerc    JobID = "clerc"
)

// Job définit une classe/job avec ses caractéristiques
type Job struct {
	ID                  JobID                  `json:"id"`
	Nom                 string                 `json:"nom"`
	Description         string                 `json:"description"`
	StatsPrincipales    []string               `json:"stats_principales"`
	CompetencesBase     []CompetenceID         `json:"competences_base"`
	CompetencesNiveau   map[int][]CompetenceID `json:"competences_niveau"`
	JobsPrerequisNiveau map[JobID]int          `json:"jobs_prerequis_niveau,omitempty"`
}

// JobManager gère les jobs et leurs relations
type JobManager struct {
	jobs map[JobID]*Job
}

// NewJobManager crée un nouveau gestionnaire de jobs
func NewJobManager() *JobManager {
	jm := &JobManager{
		jobs: make(map[JobID]*Job),
	}

	jm.initialiserJobsDeBase()
	return jm
}

func (jm *JobManager) initialiserJobsDeBase() {
	// Guerrier
	jm.jobs[JobGuerrier] = &Job{
		ID:               JobGuerrier,
		Nom:              "Guerrier",
		Description:      "Combattant au corps à corps, tank et DPS physique",
		StatsPrincipales: []string{"Force", "Vie", "Defense"},
		CompetencesBase:  []CompetenceID{"attaque_simple", "garde"},
		CompetencesNiveau: map[int][]CompetenceID{
			2:  {"charge"},
			3:  {"coup_puissant"},
			5:  {"berserker"},
			7:  {"contre_attaque"},
			10: {"maitre_epee"},
		},
	}

	// Mage
	jm.jobs[JobMage] = &Job{
		ID:               JobMage,
		Nom:              "Mage",
		Description:      "Utilisateur de magie, DPS à distance et support",
		StatsPrincipales: []string{"Intelligence", "Mana", "Magie"},
		CompetencesBase:  []CompetenceID{"boule_feu", "meditation"},
		CompetencesNiveau: map[int][]CompetenceID{
			2:  {"eclair"},
			3:  {"soin"},
			5:  {"explosion"},
			7:  {"teleportation"},
			10: {"meteor"},
		},
	}

	// Archer
	jm.jobs[JobArcher] = &Job{
		ID:               JobArcher,
		Nom:              "Archer",
		Description:      "Combattant à distance, DPS et utilitaire",
		StatsPrincipales: []string{"Dextérité", "Agilité", "Précision"},
		CompetencesBase:  []CompetenceID{"tir_simple", "viser"},
		CompetencesNiveau: map[int][]CompetenceID{
			2:  {"tir_multiple"},
			3:  {"fleche_empoisonnee"},
			5:  {"pluie_fleches"},
			7:  {"tir_perçant"},
			10: {"fleche_legendaire"},
		},
	}

	// Voleur
	jm.jobs[JobVoleur] = &Job{
		ID:               JobVoleur,
		Nom:              "Voleur",
		Description:      "Assassin furtif, DPS critique et mobilité",
		StatsPrincipales: []string{"Agilité", "Dextérité", "Chance"},
		CompetencesBase:  []CompetenceID{"attaque_sournoise", "furtivite"},
		CompetencesNiveau: map[int][]CompetenceID{
			2:  {"vol"},
			3:  {"coup_critique"},
			5:  {"invisibilite"},
			7:  {"attaque_fatale"},
			10: {"maitre_assassin"},
		},
	}

	// Clerc
	jm.jobs[JobClerc] = &Job{
		ID:               JobClerc,
		Nom:              "Clerc",
		Description:      "Soigneur et support, magie divine",
		StatsPrincipales: []string{"Sagesse", "Mana", "Magie"},
		CompetencesBase:  []CompetenceID{"soin_leger", "benediction"},
		CompetencesNiveau: map[int][]CompetenceID{
			2:  {"soin_groupe"},
			3:  {"purification"},
			5:  {"resurrection"},
			7:  {"protection_divine"},
			10: {"miracle"},
		},
	}
}

// GetJob retourne un job par son ID
func (jm *JobManager) GetJob(jobID JobID) (*Job, error) {
	job, exists := jm.jobs[jobID]
	if !exists {
		return nil, NewErrJobInexistant(jobID)
	}
	return job, nil
}

// GetJobsDepart retourne les 5 jobs de départ disponibles
func (jm *JobManager) GetJobsDepart() []*Job {
	jobsDepart := []JobID{JobGuerrier, JobMage, JobArcher, JobVoleur, JobClerc}
	jobs := make([]*Job, 0, len(jobsDepart))

	for _, jobID := range jobsDepart {
		if job, exists := jm.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// GetAllJobs retourne tous les jobs disponibles
func (jm *JobManager) GetAllJobs() []*Job {
	jobs := make([]*Job, 0, len(jm.jobs))
	for _, job := range jm.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// PeutChangerJob vérifie si un joueur peut changer vers un job
func (jm *JobManager) PeutChangerJob(joueur *Joueur, nouveauJobID JobID) error {
	job, err := jm.GetJob(nouveauJobID)
	if err != nil {
		return err
	}

	// Vérifier si le job est dans les jobs débloqués du joueur
	for _, jobDebloque := range joueur.JobsDebloquees {
		if jobDebloque == nouveauJobID {
			return nil
		}
	}

	// Vérifier les prérequis si c'est un job avancé
	if job.JobsPrerequisNiveau != nil {
		for jobPrerequisID, niveauRequis := range job.JobsPrerequisNiveau {
			niveauJoueur, exists := joueur.NiveauJobs[jobPrerequisID]
			if !exists || niveauJoueur < niveauRequis {
				return NewErrPrerequisJobNonRempli(nouveauJobID, jobPrerequisID, niveauRequis)
			}
		}

		// Si tous les prérequis sont remplis, débloquer le job
		return nil
	}

	return NewErrJobNonDebloque(nouveauJobID)
}

// Types d'erreur pour le système de jobs
type ErrJobInexistant struct {
	JobID JobID
}

func NewErrJobInexistant(jobID JobID) error {
	return &ErrJobInexistant{JobID: jobID}
}

func (e *ErrJobInexistant) Error() string {
	return "job inexistant: " + string(e.JobID)
}

type ErrJobNonDebloque struct {
	JobID JobID
}

func NewErrJobNonDebloque(jobID JobID) error {
	return &ErrJobNonDebloque{JobID: jobID}
}

func (e *ErrJobNonDebloque) Error() string {
	return "job non débloqué: " + string(e.JobID)
}

type ErrPrerequisJobNonRempli struct {
	JobID          JobID
	JobPrerequisID JobID
	NiveauRequis   int
}

func NewErrPrerequisJobNonRempli(jobID, prerequisID JobID, niveauRequis int) error {
	return &ErrPrerequisJobNonRempli{
		JobID:          jobID,
		JobPrerequisID: prerequisID,
		NiveauRequis:   niveauRequis,
	}
}

func (e *ErrPrerequisJobNonRempli) Error() string {
	return string(e.JobID) + " nécessite " + string(e.JobPrerequisID) + " niveau " + string(rune(e.NiveauRequis))
}
