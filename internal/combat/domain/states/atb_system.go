package states

import (
	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// ATBSystem gère le système de jauge ATB (Active Time Battle)
// Système de gestion du temps actif pour déterminer l'ordre des tours
type ATBSystem struct {
	gauges map[domain.UnitID]*ATBGauge
}

// ATBGauge représente la jauge ATB d'une unité
type ATBGauge struct {
	UnitID domain.UnitID
	Value  int // 0-100
	Speed  int // Vitesse de remplissage (basée sur SPD stat)
	Active bool
}

// NewATBSystem crée un nouveau système ATB
func NewATBSystem() *ATBSystem {
	return &ATBSystem{
		gauges: make(map[domain.UnitID]*ATBGauge),
	}
}

// InitializeGauge initialise la jauge d'une unité
func (atb *ATBSystem) InitializeGauge(unitID domain.UnitID, speed int) {
	atb.gauges[unitID] = &ATBGauge{
		UnitID: unitID,
		Value:  0,
		Speed:  speed,
		Active: true,
	}
}

// Tick fait progresser toutes les jauges actives
func (atb *ATBSystem) Tick() {
	for _, gauge := range atb.gauges {
		if gauge.Active && gauge.Value < 100 {
			gauge.Value += gauge.Speed
			if gauge.Value > 100 {
				gauge.Value = 100
			}
		}
	}
}

// GetReadyUnits retourne les unités avec ATB >= 100
func (atb *ATBSystem) GetReadyUnits() []domain.UnitID {
	ready := make([]domain.UnitID, 0)
	for _, gauge := range atb.gauges {
		if gauge.Active && gauge.Value >= 100 {
			ready = append(ready, gauge.UnitID)
		}
	}
	return ready
}

// ResetGauge remet la jauge d'une unité à 0
func (atb *ATBSystem) ResetGauge(unitID domain.UnitID) {
	if gauge, exists := atb.gauges[unitID]; exists {
		gauge.Value = 0
	}
}

// DeactivateGauge désactive la jauge d'une unité (morte, KO)
func (atb *ATBSystem) DeactivateGauge(unitID domain.UnitID) {
	if gauge, exists := atb.gauges[unitID]; exists {
		gauge.Active = false
	}
}

// GetGaugeValue retourne la valeur actuelle de la jauge
func (atb *ATBSystem) GetGaugeValue(unitID domain.UnitID) int {
	if gauge, exists := atb.gauges[unitID]; exists {
		return gauge.Value
	}
	return 0
}
