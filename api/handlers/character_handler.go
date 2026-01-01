package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/gin-gonic/gin"
)

// CharacterHandler gère les endpoints liés à la création de personnages
type CharacterHandler struct {
	characterCreator *domain.CharacterCreator
	// Note: Pas de persistance BDD dans ce projet
	// La sauvegarde se fera via API externe plus tard
}

// NewCharacterHandler crée un nouveau handler pour les personnages
func NewCharacterHandler() *CharacterHandler {
	return &CharacterHandler{
		characterCreator: domain.NewCharacterCreator(),
	}
}

// RegisterRoutes enregistre les routes du handler
func (h *CharacterHandler) RegisterRoutes(router *gin.RouterGroup) {
	characters := router.Group("/characters")
	{
		// Endpoints de consultation (pas de BDD nécessaire)
		characters.GET("/classes", h.GetClasses)
		characters.GET("/classes/:classe", h.GetClasseTemplate)
		characters.GET("/templates", h.GetAllTemplates)
		
		// Endpoint de création (retourne le personnage créé sans sauvegarde)
		characters.POST("/create", h.CreateCharacter)
		
		// TODO: Les autres endpoints CRUD seront implémentés avec l'API BDD externe
		// characters.GET("/", h.GetAllCharacters)
		// characters.GET("/:id", h.GetCharacter)
		// characters.PUT("/:id", h.UpdateCharacter)
		// characters.DELETE("/:id", h.DeleteCharacter)
		// characters.GET("/player/:joueur_id", h.GetCharactersByPlayer)
	}
}

// GetClasses retourne la liste des classes disponibles
// @Summary Liste des classes disponibles
// @Description Retourne toutes les classes de personnage disponibles
// @Tags Characters
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /aether/v1/characters/classes [get]
func (h *CharacterHandler) GetClasses(c *gin.Context) {
	classes := h.characterCreator.GetClassesDisponibles()

	classesStr := make([]string, len(classes))
	for i, classe := range classes {
		classesStr[i] = string(classe)
	}

	c.JSON(http.StatusOK, gin.H{
		"classes": classesStr,
		"count":   len(classesStr),
	})
}

// GetClasseTemplate retourne les informations d'une classe spécifique
// @Summary Template d'une classe
// @Description Retourne les stats de base et compétences d'une classe
// @Tags Characters
// @Param classe path string true "Nom de la classe"
// @Produce json
// @Success 200 {object} domain.PersonnageTemplate
// @Failure 404 {object} map[string]string
// @Router /aether/v1/characters/classes/{classe} [get]
func (h *CharacterHandler) GetClasseTemplate(c *gin.Context) {
	classeStr := c.Param("classe")
	classe := domain.ClassePersonnage(classeStr)

	template, err := h.characterCreator.GetTemplate(classe)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Classe non trouvée"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetAllTemplates retourne tous les templates de classes
// @Summary Tous les templates de classes
// @Description Retourne les informations de toutes les classes disponibles
// @Tags Characters
// @Produce json
// @Success 200 {object} map[string]domain.PersonnageTemplate
// @Router /aether/v1/characters/templates [get]
func (h *CharacterHandler) GetAllTemplates(c *gin.Context) {
	classes := h.characterCreator.GetClassesDisponibles()
	templates := make(map[string]*domain.PersonnageTemplate)

	for _, classe := range classes {
		template, err := h.characterCreator.GetTemplate(classe)
		if err == nil {
			templates[string(classe)] = template
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"count":     len(templates),
	})
}

// CreateCharacter crée un nouveau personnage
// @Summary Crée un personnage
// @Description Crée un nouveau personnage basé sur une classe (sans sauvegarde BDD)
// @Tags Characters
// @Accept json
// @Produce json
// @Param character body domain.PersonnageRequest true "Données du personnage"
// @Success 201 {object} domain.PersonnageResponse
// @Failure 400 {object} map[string]string
// @Router /aether/v1/characters/create [post]
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req domain.PersonnageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Validation
	if req.Nom == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le nom est requis"})
		return
	}

	if req.Classe == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La classe est requise"})
		return
	}

	// Générer un ID unique basé sur timestamp
	timestamp := time.Now().UnixNano()
	characterID := domain.UnitID("char-" + strconv.FormatInt(timestamp, 10))
	
	// TeamID par défaut si non fourni
	teamID := domain.TeamID(req.TeamID)
	if teamID == "" {
		teamID = "team-player"
	}

	// Position
	position, err := shared.NewPosition(req.Position.X, req.Position.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Position invalide"})
		return
	}

	// Créer le personnage
	personnage, err := h.characterCreator.CreerPersonnage(
		characterID,
		req.Nom,
		req.Classe,
		teamID,
		position,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir en réponse (sans sauvegarde BDD)
	response := h.characterCreator.ToResponse(personnage, req.Classe)
	
	// Note pour Angular: ce personnage n'est pas sauvegardé en BDD
	// Il faudra l'envoyer à l'API BDD externe pour la persistance
	c.Header("X-Note", "Personnage créé en mémoire uniquement, pas sauvegardé en BDD")

	c.JSON(http.StatusCreated, response)
}