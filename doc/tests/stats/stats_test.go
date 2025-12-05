// internal/shared/domain/stats/stats_test.go

package stats_test

import (
	"testing"

	"aether-engine-server/internal/shared/domain/stats" // Assurez-vous que le chemin est correct
)

// NOTE: Suppose que NewStats(maxHP, maxMP, ATK, DEF, SPD, MAG, RES) est la signature.

func TestNewStats_InvalidBaseStats(t *testing.T) {
	// Testons des invariants de base comme HP > 0, ATK >= 1
	_, err := stats.NewStats(0, 10, 10, 10, 10, 10, 10, 10, 3, 80)
	if err == nil {
		t.Fatal("Expected error for MaxHP <= 0")
	}
	_, err = stats.NewStats(100, 10, 10, 0, 10, 0, 10, 10, 3, 80)
	if err == nil {
		t.Fatal("Expected error for ATK <= 0 (pour une unité de base)")
	}
}

func TestStats_TakeDamage_Invariant(t *testing.T) {
	// Création: HP=100
	s, _ := stats.NewStats(100, 10, 10, 10, 10, 10, 10, 10, 3, 80)

	// 1. Dégâts normaux
	s.TakeDamage(50)
	if s.CurrentHP() != 50 {
		t.Fatalf("Expected HP 50, got %d", s.CurrentHP())
	}

	// 2. Dégâts excessifs (plancher à 0)
	s.TakeDamage(100)
	if s.CurrentHP() != 0 || !s.IsKO() {
		t.Fatalf("Expected HP 0 and IsKO=true, got HP %d, IsKO %t", s.CurrentHP(), s.IsKO())
	}

	// 3. Dégâts sur KO (doit rester à 0)
	s.TakeDamage(10)
	if s.CurrentHP() != 0 {
		t.Fatalf("HP should remain 0, got %d", s.CurrentHP())
	}
}

func TestStats_Heal_Invariant(t *testing.T) {
	// Création: HP=100
	s, _ := stats.NewStats(100, 10, 10, 10, 10, 10, 10, 10, 3, 80)
	s.TakeDamage(80) // HP=20

	// 1. Soins normaux
	s.Heal(50) // HP=70
	if s.CurrentHP() != 70 {
		t.Fatalf("Expected HP 70, got %d", s.CurrentHP())
	}

	// 2. Soins excessifs (plafond à MaxHP)
	s.Heal(100) // HP=100
	if s.CurrentHP() != 100 {
		t.Fatalf("Expected HP 100, got %d", s.CurrentHP())
	}
}

func TestStats_ConsumeMP(t *testing.T) {
	// Création: MP=50
	s, _ := stats.NewStats(100, 50, 10, 10, 10, 10, 10, 10, 3, 80)

	// 1. Consommation valide
	if err := s.ConsumeMP(20); err != nil {
		t.Fatal("Expected no error on valid MP consumption")
	}
	if s.CurrentMP() != 30 {
		t.Fatalf("Expected MP 30, got %d", s.CurrentMP())
	}

	// 2. Consommation insuffisante
	if err := s.ConsumeMP(40); err == nil {
		t.Fatal("Expected error for insufficient MP")
	}
	// Le MP ne doit pas avoir changé
	if s.CurrentMP() != 30 {
		t.Fatalf("MP should remain 30, got %d", s.CurrentMP())
	}
}
