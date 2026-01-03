package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aether-engine/aether-engine/internal/joueur/domain"
	"github.com/gin-gonic/gin"
)

// JoueurHandler gère les endpoints liés aux joueurs
type JoueurHandler struct {
	jobManager *domain.JobManager
	// Note: Pas de persistance BDD dans ce projet
	// La sauvegarde se fera via API externe plus tard
}

// NewJoueurHandler crée un nouveau handler pour les joueurs
func NewJoueurHandler() *JoueurHandler {
	return &JoueurHandler{
		jobManager: domain.NewJobManager(),
	}
}

// RegisterRoutes enregistre les routes du handler
func (h *JoueurHandler) RegisterRoutes(router *gin.RouterGroup) {
	joueurs := router.Group("/joueurs")
	{
		// Phase de customisation
		joueurs.GET("/customisation/options", h.GetCustomisationOptions)

		// Jobs disponibles
		joueurs.GET("/jobs", h.GetJobsDepart)
		joueurs.GET("/jobs/:job_id", h.GetJobDetails)

		// Création de joueur
		joueurs.POST("/create", h.CreateJoueur)

		// Gestion des compétences (théorycraft)
		joueurs.POST("/:joueur_id/competences/assigner", h.AssignerCompetence)
		joueurs.DELETE("/:joueur_id/competences/:type/:slot", h.RetirerCompetence)
		joueurs.GET("/:joueur_id/competences", h.GetCompetencesJoueur)

		// Changement de job
		joueurs.POST("/:joueur_id/jobs/changer", h.ChangerJob)

		// TODO: Les autres endpoints CRUD seront implémentés avec l'API BDD externe
		// joueurs.GET("/", h.GetAllJoueurs)
		// joueurs.GET("/:id", h.GetJoueur)
		// joueurs.PUT("/:id", h.UpdateJoueur)
		// joueurs.DELETE("/:id", h.DeleteJoueur)
	}
}

// GetCustomisationOptions retourne toutes les options de customisation
// @Summary Options de customisation physique
// @Description Retourne toutes les options disponibles pour personnaliser l'apparence
// @Tags Joueurs
// @Produce json
// @Success 200 {object} CustomisationOptions
// @Router /aether/v1/joueurs/customisation/options [get]
func (h *JoueurHandler) GetCustomisationOptions(c *gin.Context) {
	options := CustomisationOptions{
		Sexes: []string{
			string(domain.SexeMasculin),
			string(domain.SexeFeminin),
			string(domain.SexeAutre),
		},
		Tailles: []string{
			string(domain.TaillePetite),
			string(domain.TailleMoyenne),
			string(domain.TailleGrande),
		},
		StylesCheveux: []string{
			"chauve", "courts", "mi-longs", "longs", "tresses", "chignon",
		},
		StylesBarbe: []string{
			"aucune", "courte", "moyenne", "longue", "bouc", "moustache",
		},
		CouleursPredefinies: []domain.CouleurPersonnage{
			{R: 255, G: 220, B: 177}, // Peau claire
			{R: 210, G: 180, B: 140}, // Peau medium
			{R: 139, G: 69, B: 19},   // Peau foncée
			{R: 0, G: 0, B: 0},       // Noir (cheveux)
			{R: 101, G: 67, B: 33},   // Brun
			{R: 255, G: 255, B: 0},   // Blond
			{R: 165, G: 42, B: 42},   // Roux
			{R: 0, G: 0, B: 255},     // Bleu (yeux)
			{R: 0, G: 128, B: 0},     // Vert (yeux)
		},
		TatouagesDisponibles: []TatouageOption{
			{ID: "dragon", Nom: "Dragon", Positions: []string{"dos", "bras_gauche", "bras_droit"}},
			{ID: "tribal", Nom: "Tribal", Positions: []string{"bras_gauche", "bras_droit", "jambe_gauche"}},
			{ID: "rune", Nom: "Rune mystique", Positions: []string{"front", "main_gauche", "main_droite"}},
		},
		CicatricesDisponibles: []CicatriceOption{
			{ID: "balafre", Nom: "Balafre", Positions: []string{"visage", "joue_gauche", "joue_droite"}},
			{ID: "combat", Nom: "Marque de combat", Positions: []string{"bras", "torse", "dos"}},
		},
	}

	c.JSON(http.StatusOK, options)
}

// GetJobsDepart retourne les 5 jobs de départ
// @Summary Jobs de départ disponibles
// @Description Retourne les 5 jobs de base sélectionnables à la création
// @Tags Joueurs
// @Produce json
// @Success 200 {array} domain.Job
// @Router /aether/v1/joueurs/jobs [get]
func (h *JoueurHandler) GetJobsDepart(c *gin.Context) {
	jobs := h.jobManager.GetJobsDepart()
	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// GetJobDetails retourne les détails d'un job
// @Summary Détails d'un job
// @Description Retourne les informations complètes d'un job spécifique
// @Tags Joueurs
// @Param job_id path string true "ID du job"
// @Produce json
// @Success 200 {object} domain.Job
// @Failure 404 {object} map[string]string
// @Router /aether/v1/joueurs/jobs/{job_id} [get]
func (h *JoueurHandler) GetJobDetails(c *gin.Context) {
	jobID := domain.JobID(c.Param("job_id"))

	job, err := h.jobManager.GetJob(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job non trouvé"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CreateJoueur crée un nouveau joueur avec customisation
// @Summary Crée un joueur personnalisé
// @Description Crée un nouveau joueur avec apparence customisée et job de départ
// @Tags Joueurs
// @Accept json
// @Produce json
// @Param joueur body CreateJoueurRequest true "Données du joueur"
// @Success 201 {object} domain.Joueur
// @Failure 400 {object} map[string]string
// @Router /aether/v1/joueurs/create [post]
func (h *JoueurHandler) CreateJoueur(c *gin.Context) {
	var req CreateJoueurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Validation
	if req.Nom == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le nom est requis"})
		return
	}

	// Vérifier que le job existe et est un job de départ
	job, err := h.jobManager.GetJob(req.JobInitial)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job invalide"})
		return
	}

	jobsDepart := h.jobManager.GetJobsDepart()
	jobValide := false
	for _, jobDepart := range jobsDepart {
		if jobDepart.ID == req.JobInitial {
			jobValide = true
			break
		}
	}

	if !jobValide {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job non disponible à la création"})
		return
	}

	// Générer un ID unique
	timestamp := time.Now().UnixNano()
	joueurID := domain.JoueurID("joueur-" + strconv.FormatInt(timestamp, 10))

	// Créer le joueur
	joueur := domain.NewJoueur(joueurID, req.Nom, req.Apparence, req.JobInitial)

	// Ajouter les compétences de base du job
	for _, competenceID := range job.CompetencesBase {
		joueur.SlotsCompetences.AjouterCompetenceDisponible(competenceID)
	}

	// Note pour Angular: ce joueur n'est pas sauvegardé en BDD
	c.Header("X-Note", "Joueur créé en mémoire uniquement, pas sauvegardé en BDD")

	c.JSON(http.StatusCreated, joueur)
}

// AssignerCompetence assigne une compétence à un slot
func (h *JoueurHandler) AssignerCompetence(c *gin.Context) {
	var req AssignerCompetenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// TODO: Récupérer le joueur depuis la BDD externe
	// Pour l'instant, retourne une erreur car pas de persistance
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":      "Fonctionnalité nécessite une BDD externe pour récupérer le joueur",
		"joueur_id":  c.Param("joueur_id"),
		"competence": req,
	})
}

// RetirerCompetence retire une compétence d'un slot
func (h *JoueurHandler) RetirerCompetence(c *gin.Context) {
	joueurID := c.Param("joueur_id")
	typeComp := domain.TypeCompetence(c.Param("type"))
	slot, err := strconv.Atoi(c.Param("slot"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slot invalide"})
		return
	}

	// TODO: Récupérer le joueur depuis la BDD externe
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":     "Fonctionnalité nécessite une BDD externe pour récupérer le joueur",
		"joueur_id": joueurID,
		"type":      typeComp,
		"slot":      slot,
	})
}

// GetCompetencesJoueur retourne les compétences d'un joueur
func (h *JoueurHandler) GetCompetencesJoueur(c *gin.Context) {
	joueurID := c.Param("joueur_id")

	// TODO: Récupérer le joueur depuis la BDD externe
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":     "Fonctionnalité nécessite une BDD externe pour récupérer le joueur",
		"joueur_id": joueurID,
	})
}

// ChangerJob change le job d'un joueur
func (h *JoueurHandler) ChangerJob(c *gin.Context) {
	var req ChangerJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	joueurID := c.Param("joueur_id")

	// TODO: Récupérer le joueur depuis la BDD externe
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":       "Fonctionnalité nécessite une BDD externe pour récupérer le joueur",
		"joueur_id":   joueurID,
		"nouveau_job": req.NouveauJob,
	})
}

// DTOs pour les requêtes
type CreateJoueurRequest struct {
	Nom        string                   `json:"nom" binding:"required"`
	Apparence  domain.ApparencePhysique `json:"apparence" binding:"required"`
	JobInitial domain.JobID             `json:"job_initial" binding:"required"`
}

type AssignerCompetenceRequest struct {
	Type         domain.TypeCompetence `json:"type" binding:"required"`
	SlotIndex    int                   `json:"slot_index" binding:"required"`
	CompetenceID domain.CompetenceID   `json:"competence_id" binding:"required"`
}

type ChangerJobRequest struct {
	NouveauJob domain.JobID `json:"nouveau_job" binding:"required"`
}

// DTOs pour les réponses
type CustomisationOptions struct {
	Sexes                 []string                   `json:"sexes"`
	Tailles               []string                   `json:"tailles"`
	StylesCheveux         []string                   `json:"styles_cheveux"`
	StylesBarbe           []string                   `json:"styles_barbe"`
	CouleursPredefinies   []domain.CouleurPersonnage `json:"couleurs_predefinies"`
	TatouagesDisponibles  []TatouageOption           `json:"tatouages_disponibles"`
	CicatricesDisponibles []CicatriceOption          `json:"cicatrices_disponibles"`
}

type TatouageOption struct {
	ID        string   `json:"id"`
	Nom       string   `json:"nom"`
	Positions []string `json:"positions"`
}

type CicatriceOption struct {
	ID        string   `json:"id"`
	Nom       string   `json:"nom"`
	Positions []string `json:"positions"`
}
