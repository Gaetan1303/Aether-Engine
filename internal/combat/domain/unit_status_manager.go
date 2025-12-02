package domain

import (
	"errors"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// UnitStatusManager gère les statuts d'une unité
// Responsabilités: Ajout/retrait/traitement des statuts (buffs/debuffs)
// Single Responsibility Principle - Une seule raison de changer: logique des statuts
type UnitStatusManager struct {
	statuses []*shared.Statut
}

// NewUnitStatusManager crée un nouveau gestionnaire de statuts
func NewUnitStatusManager() *UnitStatusManager {
	return &UnitStatusManager{
		statuses: make([]*shared.Statut, 0),
	}
}

// Statuses retourne tous les statuts actifs
func (m *UnitStatusManager) Statuses() []*shared.Statut {
	return m.statuses
}

// AddStatus ajoute un statut à l'unité
func (m *UnitStatusManager) AddStatus(status *shared.Statut) error {
	if status == nil {
		return errors.New("statut nil")
	}

	// Vérifier si le statut existe déjà (même type)
	for i, existing := range m.statuses {
		if existing.Type() == status.Type() {
			// Remplacer par le nouveau (refresh)
			m.statuses[i] = status
			return nil
		}
	}

	// Ajouter le nouveau statut
	m.statuses = append(m.statuses, status)
	return nil
}

// RemoveStatus retire un statut par type
func (m *UnitStatusManager) RemoveStatus(statusType shared.TypeStatut) {
	for i, status := range m.statuses {
		if status.Type() == statusType {
			// Supprimer en gardant l'ordre
			m.statuses = append(m.statuses[:i], m.statuses[i+1:]...)
			return
		}
	}
}

// ProcessStatuses traite les statuts (décrémenter durée, appliquer effets)
// Retourne les effets périodiques appliqués
func (m *UnitStatusManager) ProcessStatuses(target shared.StatsModifiable) []shared.EffetStatut {
	effets := make([]shared.EffetStatut, 0)

	// Traiter chaque statut
	for i := len(m.statuses) - 1; i >= 0; i-- {
		status := m.statuses[i]

		// Appliquer l'effet périodique
		if effect := status.AppliquerEffetPeriodique(target); effect != nil {
			effets = append(effets, *effect)
		}

		// Décrémenter la durée
		status.DecrémenterDuree()

		// Retirer si expiré
		if status.EstExpire() {
			m.statuses = append(m.statuses[:i], m.statuses[i+1:]...)
		}
	}

	return effets
}

// HasStatus vérifie si un statut est actif
func (m *UnitStatusManager) HasStatus(statusType shared.TypeStatut) bool {
	for _, status := range m.statuses {
		if status.Type() == statusType {
			return true
		}
	}
	return false
}

// GetStatus retourne un statut par type
func (m *UnitStatusManager) GetStatus(statusType shared.TypeStatut) *shared.Statut {
	for _, status := range m.statuses {
		if status.Type() == statusType {
			return status
		}
	}
	return nil
}

// ClearAllStatuses retire tous les statuts
func (m *UnitStatusManager) ClearAllStatuses() {
	m.statuses = make([]*shared.Statut, 0)
}

// IsStunned vérifie si l'unité est étourdie
func (m *UnitStatusManager) IsStunned() bool {
	return m.HasStatus(shared.TypeStatutStun)
}

// IsSilenced vérifie si l'unité est silencée
func (m *UnitStatusManager) IsSilenced() bool {
	return m.HasStatus(shared.TypeStatutSilence)
}

// IsRooted vérifie si l'unité est enracinée
func (m *UnitStatusManager) IsRooted() bool {
	return m.HasStatus(shared.TypeStatutRoot)
}

// IsPoisoned vérifie si l'unité est empoisonnée
func (m *UnitStatusManager) IsPoisoned() bool {
	return m.HasStatus(shared.TypeStatutPoison)
}

// BlocksActions vérifie si un statut bloque les actions
func (m *UnitStatusManager) BlocksActions() bool {
	for _, status := range m.statuses {
		if status.BloqueActions() {
			return true
		}
	}
	return false
}

// BlocksMovement vérifie si un statut bloque le déplacement
func (m *UnitStatusManager) BlocksMovement() bool {
	for _, status := range m.statuses {
		if status.BloqueDeplacement() {
			return true
		}
	}
	return false
}

// StatusCount retourne le nombre de statuts actifs
func (m *UnitStatusManager) StatusCount() int {
	return len(m.statuses)
}
