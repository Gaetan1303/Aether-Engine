package application

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
)

// MockEventStore pour les tests
type MockEventStore struct {
	events map[string][]domain.Evenement
}

func NewMockEventStore() *MockEventStore {
	return &MockEventStore{
		events: make(map[string][]domain.Evenement),
	}
}

func (m *MockEventStore) AppendEvents(aggregateID string, events []domain.Evenement, expectedVersion int) error {
	if m.events[aggregateID] == nil {
		m.events[aggregateID] = make([]domain.Evenement, 0)
	}
	m.events[aggregateID] = append(m.events[aggregateID], events...)
	return nil
}

func (m *MockEventStore) LoadEvents(aggregateID string) ([]domain.Evenement, error) {
	return m.events[aggregateID], nil
}

func (m *MockEventStore) LoadEventsFromVersion(aggregateID string, fromVersion int) ([]domain.Evenement, error) {
	return m.events[aggregateID], nil
}

func (m *MockEventStore) SaveSnapshot(aggregateID string, version int, data []byte) error {
	return nil
}

func (m *MockEventStore) LoadSnapshot(aggregateID string) (version int, data []byte, err error) {
	return 0, nil, nil
}

// MockEventPublisher pour les tests
type MockEventPublisher struct {
	publishedEvents []domain.Evenement
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		publishedEvents: make([]domain.Evenement, 0),
	}
}

func (m *MockEventPublisher) Publish(event domain.Evenement) error {
	m.publishedEvents = append(m.publishedEvents, event)
	return nil
}

// TestDemarrerCombatWithStepC vérifie que tous les patterns Step C sont initialisés
func TestDemarrerCombatWithStepC(t *testing.T) {
	eventStore := NewMockEventStore()
	publisher := NewMockEventPublisher()
	engine := NewCombatEngine(eventStore, publisher)

	// Créer un combat minimal
	joueur1 := "player1"
	joueur2 := "player2"

	equipe1 := EquipeDTO{
		ID:       "team1",
		Nom:      "Test",
		JoueurID: &joueur1,
		Membres: []UniteDTO{
			{
				ID:       "unit1",
				Nom:      "Test Unit",
				TeamID:   "team1",
				Stats:    StatsDTO{HP: 100, MP: 50, SPD: 10, MOV: 3, ATK: 10, DEF: 5},
				Position: PositionDTO{X: 0, Y: 0},
			},
		},
	}

	equipe2 := EquipeDTO{
		ID:       "team2",
		Nom:      "Test2",
		JoueurID: &joueur2,
		Membres: []UniteDTO{
			{
				ID:       "unit2",
				Nom:      "Test Unit 2",
				TeamID:   "team2",
				Stats:    StatsDTO{HP: 100, MP: 50, SPD: 10, MOV: 3, ATK: 10, DEF: 5},
				Position: PositionDTO{X: 5, Y: 5},
			},
		},
	}

	cmd := CommandeDemarrerCombat{
		CombatID: "combat-step-c-test",
		Equipes:  []EquipeDTO{equipe1, equipe2},
		Grille:   GrilleDTO{Largeur: 10, Hauteur: 10},
	}

	combatDTO, err := engine.DemarrerCombat(cmd)
	if err != nil {
		t.Fatalf("Erreur démarrage: %v", err)
	}

	if combatDTO.ID != "combat-step-c-test" {
		t.Errorf("ID incorrect: %s", combatDTO.ID)
	}

	if combatDTO.Etat != "actif" {
		t.Errorf("État incorrect: %s (attendu: actif)", combatDTO.Etat)
	}

	// Vérifier qu'il y a bien 2 équipes
	if len(combatDTO.Equipes) != 2 {
		t.Errorf("Nombre d'équipes incorrect: %d (attendu: 2)", len(combatDTO.Equipes))
	}

	// Vérifier que les événements ont été sauvegardés
	events, err := eventStore.LoadEvents("combat-step-c-test")
	if err != nil {
		t.Fatalf("Erreur chargement événements: %v", err)
	}

	if len(events) == 0 {
		t.Error("❌ Aucun événement sauvegardé")
	} else {
		t.Logf("✅ %d événements sauvegardés", len(events))
	}

	// Vérifier que les événements ont été publiés
	if len(publisher.publishedEvents) == 0 {
		t.Error("❌ Aucun événement publié")
	} else {
		t.Logf("✅ %d événements publiés", len(publisher.publishedEvents))
	}

	t.Logf("✅ Combat démarré avec Step C: %s (État: %s, Tour: %d)",
		combatDTO.ID, combatDTO.Etat, combatDTO.TourActuel)
}
