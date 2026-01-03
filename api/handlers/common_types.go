package handlers

import (
	"time"
)

// ErrorResponse structure standard pour les erreurs API
type ErrorResponse struct {
	Error     string `json:"error"`
	Details   string `json:"details,omitempty"`
	Code      string `json:"code,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewErrorResponse crée une nouvelle réponse d'erreur
func NewErrorResponse(error, details string) ErrorResponse {
	return ErrorResponse{
		Error:     error,
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessResponse structure standard pour les réponses de succès
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// NewSuccessResponse crée une nouvelle réponse de succès
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CustomisationOptions options pour personnaliser l'apparence
type CustomisationOptions struct {
	Sexes                     []string           `json:"sexes"`
	Tailles                   []string           `json:"tailles"`
	StylesCheveux             []string           `json:"styles_cheveux"`
	StylesBarbe               []string           `json:"styles_barbe"`
	CouleursPredefinies       []CouleurOption    `json:"couleurs_predefinies"`
	TatouagesDisponibles      []TatouageOption   `json:"tatouages_disponibles"`
	CicatricesDisponibles     []CicatriceOption  `json:"cicatrices_disponibles"`
}

// CouleurOption option de couleur prédéfinie
type CouleurOption struct {
	Nom  string `json:"nom"`
	R    uint8  `json:"r"`
	G    uint8  `json:"g"`
	B    uint8  `json:"b"`
	Type string `json:"type"` // "peau", "cheveux", "yeux"
}

// TatouageOption option de tatouage disponible
type TatouageOption struct {
	ID        string   `json:"id"`
	Nom       string   `json:"nom"`
	Positions []string `json:"positions"`
}

// CicatriceOption option de cicatrice disponible
type CicatriceOption struct {
	ID        string   `json:"id"`
	Nom       string   `json:"nom"`
	Positions []string `json:"positions"`
}