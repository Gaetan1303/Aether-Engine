package handlers

import (
	joueurService "github.com/aether-engine/aether-engine/internal/joueur/application/service"
	"github.com/aether-engine/aether-engine/internal/joueur/domain"
	repo "github.com/aether-engine/aether-engine/internal/joueur/infrastructure/repository"
	"github.com/gin-gonic/gin"
)

// JoueurHandler gère les endpoints liés aux joueurs
type JoueurHandler struct {
	joueurService *joueurService.JoueurService
	jobManager    *domain.JobManager
}

// NewJoueurHandler crée un nouveau handler pour les joueurs
func NewJoueurHandler() *JoueurHandler {
	repository := repo.NewInMemoryJoueurRepository()
	joueurSvc := joueurService.NewJoueurService(repository)

	return &JoueurHandler{
		joueurService: joueurSvc,
		jobManager:    domain.NewJobManager(),
	}
}

// GetJoueurService expose le service pour les autres handlers
func (h *JoueurHandler) GetJoueurService() *joueurService.JoueurService {
	return h.joueurService
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

		// CRUD complet avec service layer
		joueurs.POST("/", h.CreateJoueur)      // RESTful : POST sur collection
		joueurs.GET("/", h.GetAllJoueurs)      // RESTful : GET collection
		joueurs.GET("/:id", h.GetJoueur)       // RESTful : GET resource
		joueurs.PUT("/:id", h.UpdateJoueur)    // RESTful : PUT resource
		joueurs.DELETE("/:id", h.DeleteJoueur) // RESTful : DELETE resource

		// Gestion des compétences - utiliser :id au lieu de :joueur_id pour éviter conflit
		joueurs.POST("/:id/competences/assigner", h.AssignerCompetence)
		joueurs.DELETE("/:id/competences/:type/:slot", h.RetirerCompetence)
		joueurs.GET("/:id/competences", h.GetCompetencesJoueur)

		// Changement de job
		joueurs.POST("/:id/jobs/changer", h.ChangerJob)

		// Endpoints legacy (deprecated mais conservés pour compatibilité)
		joueurs.POST("/create", h.CreateJoueurLegacy)
	}
}

// GetCustomisationOptions retourne toutes les options de customisation
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
		CouleursPredefinies: []CouleurOption{
			{Nom: "Peau claire", R: 255, G: 220, B: 177, Type: "peau"},
			{Nom: "Peau medium", R: 210, G: 180, B: 140, Type: "peau"},
			{Nom: "Peau foncée", R: 139, G: 69, B: 19, Type: "peau"},
			{Nom: "Cheveux noirs", R: 0, G: 0, B: 0, Type: "cheveux"},
			{Nom: "Cheveux bruns", R: 101, G: 67, B: 33, Type: "cheveux"},
			{Nom: "Cheveux blonds", R: 255, G: 255, B: 0, Type: "cheveux"},
			{Nom: "Cheveux roux", R: 165, G: 42, B: 42, Type: "cheveux"},
			{Nom: "Yeux bleus", R: 0, G: 0, B: 255, Type: "yeux"},
			{Nom: "Yeux verts", R: 0, G: 128, B: 0, Type: "yeux"},
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

	c.JSON(200, options)
}

// GetJobsDepart retourne les 5 jobs de départ
func (h *JoueurHandler) GetJobsDepart(c *gin.Context) {
	jobs := h.jobManager.GetJobsDepart()
	c.JSON(200, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// GetJobDetails retourne les détails d'un job
func (h *JoueurHandler) GetJobDetails(c *gin.Context) {
	jobID := domain.JobID(c.Param("job_id"))

	job, err := h.jobManager.GetJob(jobID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job non trouvé"})
		return
	}

	c.JSON(200, gin.H{
		"job": job,
	})
}

// AssignerCompetence assigne une compétence à un slot
func (h *JoueurHandler) AssignerCompetence(c *gin.Context) {
	c.JSON(501, NewErrorResponse("Not implemented", "Feature coming soon"))
}

// RetirerCompetence retire une compétence d'un slot
func (h *JoueurHandler) RetirerCompetence(c *gin.Context) {
	c.JSON(501, NewErrorResponse("Not implemented", "Feature coming soon"))
}

// GetCompetencesJoueur récupère les compétences d'un joueur
func (h *JoueurHandler) GetCompetencesJoueur(c *gin.Context) {
	c.JSON(501, NewErrorResponse("Not implemented", "Feature coming soon"))
}

// ChangerJob change le job d'un joueur
func (h *JoueurHandler) ChangerJob(c *gin.Context) {
	c.JSON(501, NewErrorResponse("Not implemented", "Feature coming soon"))
}

// Requête pour changer de job
type ChangerJobRequest struct {
	NouveauJob domain.JobID `json:"nouveau_job" binding:"required"`
}
