package handlers

import (
	"net/http"

	"github.com/aether-engine/aether-engine/internal/joueur/application/dto"
	"github.com/gin-gonic/gin"
)

// CreateJoueur crée un nouveau joueur (RESTful)
// @Summary Créer un nouveau joueur
// @Description Crée un joueur avec apparence et job initial
// @Tags Joueurs
// @Accept json
// @Produce json
// @Param joueur body dto.JoueurCreateDTO true "Données du joueur"
// @Success 201 {object} dto.JoueurResponseDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Nom déjà utilisé"
// @Router /aether/v1/joueurs [post]
func (h *JoueurHandler) CreateJoueur(c *gin.Context) {
	var createDTO dto.JoueurCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Details: err.Error(),
		})
		return
	}

	joueur, err := h.joueurService.CreateJoueur(c.Request.Context(), &createDTO)
	if err != nil {
		if err.Error() == "nom already exists" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Name already exists",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create joueur",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, joueur)
}

// GetAllJoueurs récupère tous les joueurs
// @Summary Lister tous les joueurs
// @Description Récupère la liste de tous les joueurs créés
// @Tags Joueurs
// @Produce json
// @Success 200 {array} dto.JoueurResponseDTO
// @Router /aether/v1/joueurs [get]
func (h *JoueurHandler) GetAllJoueurs(c *gin.Context) {
	joueurs, err := h.joueurService.GetAllJoueurs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get joueurs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, joueurs)
}

// GetJoueur récupère un joueur par ID
// @Summary Récupérer un joueur
// @Description Récupère les détails d'un joueur par son ID
// @Tags Joueurs
// @Produce json
// @Param id path string true "ID du joueur"
// @Success 200 {object} dto.JoueurResponseDTO
// @Failure 404 {object} ErrorResponse
// @Router /aether/v1/joueurs/{id} [get]
func (h *JoueurHandler) GetJoueur(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID is required",
			Details: "L'ID du joueur est requis",
		})
		return
	}

	joueur, err := h.joueurService.GetJoueur(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Joueur not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, joueur)
}

// UpdateJoueur met à jour un joueur
// @Summary Mettre à jour un joueur
// @Description Met à jour les données d'un joueur existant
// @Tags Joueurs
// @Accept json
// @Produce json
// @Param id path string true "ID du joueur"
// @Success 200 {object} dto.JoueurResponseDTO
// @Failure 404 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse "Not implemented yet"
// @Router /aether/v1/joueurs/{id} [put]
func (h *JoueurHandler) UpdateJoueur(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID is required",
			Details: "L'ID du joueur est requis",
		})
		return
	}

	// TODO: Implémenter la mise à jour
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "Update not implemented",
		Details: "Waiting for external BDD API integration",
	})
}

// DeleteJoueur supprime un joueur
// @Summary Supprimer un joueur
// @Description Supprime définitivement un joueur
// @Tags Joueurs
// @Produce json
// @Param id path string true "ID du joueur"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /aether/v1/joueurs/{id} [delete]
func (h *JoueurHandler) DeleteJoueur(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID is required",
			Details: "L'ID du joueur est requis",
		})
		return
	}

	err := h.joueurService.DeleteJoueur(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Joueur not found",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetJoueurByNom récupère un joueur par nom
func (h *JoueurHandler) GetJoueurByNom(c *gin.Context) {
	nom := c.Param("nom")
	if nom == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Nom is required",
			Details: "Le nom du joueur est requis",
		})
		return
	}

	// TODO: Implémenter avec service layer
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "Not implemented yet",
		Details: "Feature coming soon",
	})
}

// GetJoueursByZone récupère les joueurs dans une zone
func (h *JoueurHandler) GetJoueursByZone(c *gin.Context) {
	zone := c.Param("zone")
	if zone == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Zone is required",
			Details: "Le nom de la zone est requis",
		})
		return
	}

	// TODO: Implémenter avec service layer
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "Not implemented yet",
		Details: "Feature coming soon",
	})
}

// CreateJoueurLegacy endpoint legacy pour compatibilité
// @Summary Créer un joueur (legacy)
// @Description Endpoint legacy, utiliser POST /joueurs à la place
// @Tags Joueurs
// @Deprecated true
// @Accept json
// @Produce json
// @Success 301 {object} ErrorResponse
// @Router /aether/v1/joueurs/create [post]
func (h *JoueurHandler) CreateJoueurLegacy(c *gin.Context) {
	c.JSON(http.StatusMovedPermanently, ErrorResponse{
		Error:   "Deprecated endpoint",
		Details: "Use POST /aether/v1/joueurs instead",
	})
}