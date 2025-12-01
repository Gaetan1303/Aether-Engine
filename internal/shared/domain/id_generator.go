package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// IDGenerator génère des identifiants uniques pour les entités du domaine
// Implémente le Singleton Pattern (thread-safe) avec sync.Once
type IDGenerator struct {
	mu        sync.Mutex
	counter   uint64
	machineID string
	startTime time.Time
}

var (
	// Instance unique du générateur (Singleton)
	idGeneratorInstance *IDGenerator

	// sync.Once garantit une initialisation unique et thread-safe
	idGeneratorOnce sync.Once
)

// GetIDGenerator retourne l'instance singleton du générateur d'IDs
// Thread-safe grâce à sync.Once - appelé plusieurs fois, n'initialise qu'une fois
//
// Exemple d'utilisation:
//
//	gen := domain.GetIDGenerator()
//	combatID := gen.NewCombatID()
//	uniteID := gen.NewUnitID()
func GetIDGenerator() *IDGenerator {
	idGeneratorOnce.Do(func() {
		idGeneratorInstance = &IDGenerator{
			counter:   0,
			machineID: generateMachineID(),
			startTime: time.Now(),
		}
	})
	return idGeneratorInstance
}

// NewCombatID génère un nouvel ID unique pour un combat
// Format: combat_<timestamp>_<machine>_<counter>
// Exemple: combat_1701432123_a3b4_0001
func (g *IDGenerator) NewCombatID() string {
	return g.generateID("combat")
}

// NewUnitID génère un nouvel ID unique pour une unité
// Format: unit_<timestamp>_<machine>_<counter>
// Exemple: unit_1701432123_a3b4_0002
func (g *IDGenerator) NewUnitID() string {
	return g.generateID("unit")
}

// NewTeamID génère un nouvel ID unique pour une équipe
// Format: team_<timestamp>_<machine>_<counter>
// Exemple: team_1701432123_a3b4_0003
func (g *IDGenerator) NewTeamID() string {
	return g.generateID("team")
}

// NewCompetenceID génère un nouvel ID unique pour une compétence
// Format: skill_<timestamp>_<machine>_<counter>
// Exemple: skill_1701432123_a3b4_0004
func (g *IDGenerator) NewCompetenceID() string {
	return g.generateID("skill")
}

// NewObjetID génère un nouvel ID unique pour un objet
// Format: item_<timestamp>_<machine>_<counter>
// Exemple: item_1701432123_a3b4_0005
func (g *IDGenerator) NewObjetID() string {
	return g.generateID("item")
}

// NewEventID génère un nouvel ID unique pour un événement
// Format: event_<timestamp>_<machine>_<counter>
// Exemple: event_1701432123_a3b4_0006
func (g *IDGenerator) NewEventID() string {
	return g.generateID("event")
}

// generateID est la méthode interne qui génère les IDs
// Thread-safe grâce au mutex
func (g *IDGenerator) generateID(prefix string) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Incrémenter le compteur
	g.counter++

	// Timestamp en secondes depuis epoch
	timestamp := time.Now().Unix()

	// Format: prefix_timestamp_machine_counter
	// Exemple: combat_1701432123_a3b4_0001
	return fmt.Sprintf("%s_%d_%s_%04d", prefix, timestamp, g.machineID, g.counter)
}

// GenerateUUID génère un UUID v4 simple (pour compatibilité)
// Alternative à generateID pour un format plus standard
func (g *IDGenerator) GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback sur timestamp + counter si crypto/rand échoue
		g.mu.Lock()
		defer g.mu.Unlock()
		g.counter++
		return fmt.Sprintf("%d-%s-%04d", time.Now().UnixNano(), g.machineID, g.counter)
	}

	// Version 4 UUID
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// GetStats retourne des statistiques sur le générateur
// Utile pour monitoring et debugging
type GeneratorStats struct {
	TotalGenerated uint64
	MachineID      string
	Uptime         time.Duration
	StartTime      time.Time
}

// GetStats retourne les statistiques du générateur
func (g *IDGenerator) GetStats() GeneratorStats {
	g.mu.Lock()
	defer g.mu.Unlock()

	return GeneratorStats{
		TotalGenerated: g.counter,
		MachineID:      g.machineID,
		Uptime:         time.Since(g.startTime),
		StartTime:      g.startTime,
	}
}

// Reset réinitialise le compteur (UNIQUEMENT POUR LES TESTS!)
// Ne devrait JAMAIS être appelé en production
func (g *IDGenerator) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter = 0
	g.startTime = time.Now()
}

// generateMachineID génère un ID unique pour cette instance/machine
// Utilisé pour garantir l'unicité des IDs entre plusieurs instances
func generateMachineID() string {
	// Générer 4 bytes aléatoires
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		// Fallback sur timestamp si crypto/rand échoue
		return fmt.Sprintf("%x", time.Now().UnixNano()%0xFFFF)
	}
	return hex.EncodeToString(b)[:4]
}

// Helper functions pour créer des IDs typés (type-safe wrappers)

// NewCombatIDTyped retourne un ID typé pour Combat
func NewCombatIDTyped() string {
	return GetIDGenerator().NewCombatID()
}

// NewUnitIDTyped retourne un ID typé pour Unite
func NewUnitIDTyped() string {
	return GetIDGenerator().NewUnitID()
}

// NewTeamIDTyped retourne un ID typé pour Equipe
func NewTeamIDTyped() string {
	return GetIDGenerator().NewTeamID()
}

// NewCompetenceIDTyped retourne un ID typé pour Competence
func NewCompetenceIDTyped() string {
	return GetIDGenerator().NewCompetenceID()
}

// NewObjetIDTyped retourne un ID typé pour Objet
func NewObjetIDTyped() string {
	return GetIDGenerator().NewObjetID()
}

// ParseIDType extrait le type d'un ID généré
// Exemple: "combat_1701432123_a3b4_0001" -> "combat"
func ParseIDType(id string) string {
	for i, c := range id {
		if c == '_' {
			return id[:i]
		}
	}
	return ""
}

// ValidateIDFormat valide le format d'un ID généré
// Retourne true si l'ID respecte le format attendu
func ValidateIDFormat(id string) bool {
	parts := len(id) > 0
	if !parts {
		return false
	}

	// Format minimal: prefix_timestamp_machine_counter
	// Exemple: combat_1701432123_a3b4_0001 (au moins 4 segments)
	segmentCount := 0
	for _, c := range id {
		if c == '_' {
			segmentCount++
		}
	}

	return segmentCount >= 3
}
