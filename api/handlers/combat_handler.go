package handlers

import (
	"net/http"

	"github.com/aether-engine/aether-engine/internal/combat/application"
	"github.com/gin-gonic/gin"
)

// CombatHandler gère les endpoints liés aux combats
type CombatHandler struct {
	engine application.CombatEngine
}

// NewCombatHandler crée une nouvelle instance de CombatHandler
func NewCombatHandler(engine application.CombatEngine) *CombatHandler {
	return &CombatHandler{
		engine: engine,
	}
}

// DemarrerCombat démarre un nouveau combat
// POST /api/v1/combats
func (h *CombatHandler) DemarrerCombat(c *gin.Context) {
	var cmd application.CommandeDemarrerCombat

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	combat, err := h.engine.DemarrerCombat(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, combat)
}

// ExecuterAction exécute une action dans un combat
// POST /api/v1/combats/:id/actions
func (h *CombatHandler) ExecuterAction(c *gin.Context) {
	combatID := c.Param("id")

	var cmd application.CommandeExecuterAction
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.CombatID = combatID

	resultat, err := h.engine.ExecuterAction(cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultat)
}

// PasserTour passe au tour suivant
// POST /api/v1/combats/:id/tour-suivant
func (h *CombatHandler) PasserTour(c *gin.Context) {
	combatID := c.Param("id")

	cmd := application.CommandePasserTour{
		CombatID: combatID,
	}

	if err := h.engine.PasserTour(cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tour suivant"})
}

// TerminerCombat termine un combat
// POST /api/v1/combats/:id/terminer
func (h *CombatHandler) TerminerCombat(c *gin.Context) {
	combatID := c.Param("id")

	cmd := application.CommandeTerminerCombat{
		CombatID: combatID,
	}

	if err := h.engine.TerminerCombat(cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Combat terminé"})
}

// ObtenirCombat récupère l'état d'un combat
// GET /api/v1/combats/:id
func (h *CombatHandler) ObtenirCombat(c *gin.Context) {
	combatID := c.Param("id")

	query := application.QueryObtenirCombat{
		CombatID: combatID,
	}

	combat, err := h.engine.ObtenirCombat(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, combat)
}

// RegisterRoutes enregistre toutes les routes du CombatHandler
func (h *CombatHandler) RegisterRoutes(router *gin.RouterGroup) {
	combats := router.Group("/combats")
	{
		combats.POST("", h.DemarrerCombat)
		combats.GET("/:id", h.ObtenirCombat)
		combats.POST("/:id/actions", h.ExecuterAction)
		combats.POST("/:id/tour-suivant", h.PasserTour)
		combats.POST("/:id/terminer", h.TerminerCombat)
	}
}
