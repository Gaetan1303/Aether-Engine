package handlers

import (
	"net/http"

	combatApp "github.com/aether-engine/aether-engine/internal/combat/application"
	joueurService "github.com/aether-engine/aether-engine/internal/joueur/application/service"
	"github.com/gin-gonic/gin"
)

// GameHandler gère les endpoints de gameplay (lien joueur <-> combat)
type GameHandler struct {
	joueurService *joueurService.JoueurService
	combatEngine  combatApp.CombatEngine
}

// NewGameHandler crée un nouveau handler de gameplay
func NewGameHandler(joueurSvc *joueurService.JoueurService, combatEng combatApp.CombatEngine) *GameHandler {
	return &GameHandler{
		joueurService: joueurSvc,
		combatEngine:  combatEng,
	}
}

// RegisterRoutes enregistre les routes de gameplay
func (h *GameHandler) RegisterRoutes(router *gin.RouterGroup) {
	game := router.Group("/game")
	{
		// Lancer un combat avec personnages créés
		game.POST("/combat/start", h.StartCombat)
		
		// Joindre un combat existant
		game.POST("/combat/:combat_id/join", h.JoinCombat)
		
		// Actions en combat (utilise le personnage créé)
		game.POST("/combat/:combat_id/action", h.ExecuteAction)
		
		// Statut d'un combat
		game.GET("/combat/:combat_id/status", h.GetCombatStatus)
		
		// Liste des combats disponibles
		game.GET("/combats", h.GetAvailableCombats)
	}
}

// StartCombatRequest structure pour démarrer un combat
type StartCombatRequest struct {
	JoueurID     string `json:"joueur_id" binding:"required"`
	MapID        string `json:"map_id" binding:"required"`
	CombatType   string `json:"combat_type"` // "pve", "pvp", "training"
	MaxParticipants int `json:"max_participants"`
}

// StartCombat démarre un nouveau combat avec un personnage créé
// @Summary Démarrer un combat
// @Description Lance un combat en utilisant le personnage du joueur
// @Tags Game
// @Accept json
// @Produce json
// @Param request body StartCombatRequest true "Configuration du combat"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /aether/v1/game/combat/start [post]
func (h *GameHandler) StartCombat(c *gin.Context) {
	var request StartCombatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// 1. Récupérer le joueur
	_, err := h.joueurService.GetJoueur(c.Request.Context(), request.JoueurID)
	if err != nil {
		c.JSON(http.StatusNotFound, NewErrorResponse(
			"Joueur not found",
			err.Error(),
		))
		return
	}

	// TODO: 2. Convertir Joueur en Unite via ToUnite()
	// TODO: 3. Créer le combat avec CombatEngine
	// TODO: 4. Retourner l'ID du combat créé
	
	c.JSON(http.StatusNotImplemented, NewErrorResponse(
		"Not implemented yet",
		"Combat integration coming soon",
	))
}

// JoinCombat permet à un joueur de rejoindre un combat existant
func (h *GameHandler) JoinCombat(c *gin.Context) {
	_ = c.Param("combat_id") // Éviter l'erreur "declared and not used"
	
	c.JSON(http.StatusNotImplemented, NewErrorResponse(
		"Not implemented yet", 
		"Join combat feature coming soon",
	))
}

// ExecuteAction exécute une action en combat
func (h *GameHandler) ExecuteAction(c *gin.Context) {
	_ = c.Param("combat_id") // Éviter l'erreur "declared and not used"
	
	c.JSON(http.StatusNotImplemented, NewErrorResponse(
		"Not implemented yet",
		"Combat actions coming soon",
	))
}

// GetCombatStatus récupère le statut d'un combat
func (h *GameHandler) GetCombatStatus(c *gin.Context) {
	_ = c.Param("combat_id") // Éviter l'erreur "declared and not used"
	
	c.JSON(http.StatusNotImplemented, NewErrorResponse(
		"Not implemented yet",
		"Combat status coming soon",
	))
}

// GetAvailableCombats liste les combats disponibles
func (h *GameHandler) GetAvailableCombats(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, NewErrorResponse(
		"Not implemented yet",
		"Combat list coming soon",
	))
}