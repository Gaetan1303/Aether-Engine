package domain

import (
	"fmt"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// ClassePersonnage représente les classes de personnage disponibles
type ClassePersonnage string

const (
	ClasseGuerrier ClassePersonnage = "GUERRIER"
	ClassePaladin  ClassePersonnage = "PALADIN"
	ClasseArcher   ClassePersonnage = "ARCHER"
	ClasseMage     ClassePersonnage = "MAGE"
	ClasseVoleur   ClassePersonnage = "VOLEUR"
	ClasseClerc    ClassePersonnage = "CLERC"
)

// PersonnageTemplate définit les stats de base par classe
type PersonnageTemplate struct {
	Classe      ClassePersonnage
	NomAffiche  string
	Description string
	StatsBase   *shared.Stats
	Competences []string // IDs des compétences de base
}

// CharacterCreator gère la création de personnages
type CharacterCreator struct {
	templates map[ClassePersonnage]*PersonnageTemplate
}

// NewCharacterCreator crée un nouveau créateur de personnages
func NewCharacterCreator() *CharacterCreator {
	cc := &CharacterCreator{
		templates: make(map[ClassePersonnage]*PersonnageTemplate),
	}
	cc.initialiserTemplates()
	return cc
}

// initialiserTemplates configure les templates de classes
func (cc *CharacterCreator) initialiserTemplates() {
	// Guerrier - Tank équilibré
	statsGuerrier, _ := shared.NewStats(120, 30, 100, 25, 15, 5, 8, 12, 4, 85)
	cc.templates[ClasseGuerrier] = &PersonnageTemplate{
		Classe:      ClasseGuerrier,
		NomAffiche:  "Guerrier",
		Description: "Tank robuste avec haute défense",
		StatsBase:   statsGuerrier,
		Competences: []string{"charge", "defense-renforcee"},
	}

	// Paladin - Tank magique
	statsPaladin, _ := shared.NewStats(150, 50, 100, 20, 25, 8, 18, 8, 3, 80)
	cc.templates[ClassePaladin] = &PersonnageTemplate{
		Classe:      ClassePaladin,
		NomAffiche:  "Paladin",
		Description: "Tank saint avec capacités de soin",
		StatsBase:   statsPaladin,
		Competences: []string{"provocation", "soin-divin", "benediction"},
	}

	// Archer - DPS physique
	statsArcher, _ := shared.NewStats(90, 40, 80, 28, 10, 5, 8, 16, 4, 95)
	cc.templates[ClasseArcher] = &PersonnageTemplate{
		Classe:      ClasseArcher,
		NomAffiche:  "Archer",
		Description: "DPS précis avec attaques à distance",
		StatsBase:   statsArcher,
		Competences: []string{"tir-precision", "pluie-fleches"},
	}

	// Mage - DPS magique
	statsMage, _ := shared.NewStats(70, 120, 60, 8, 6, 35, 25, 10, 3, 92)
	cc.templates[ClasseMage] = &PersonnageTemplate{
		Classe:      ClasseMage,
		NomAffiche:  "Mage",
		Description: "DPS magique avec sorts destructeurs",
		StatsBase:   statsMage,
		Competences: []string{"fireball", "lightning", "boost-magique"},
	}

	// Voleur - DPS rapide
	statsVoleur, _ := shared.NewStats(80, 25, 70, 22, 8, 3, 5, 18, 5, 90)
	cc.templates[ClasseVoleur] = &PersonnageTemplate{
		Classe:      ClasseVoleur,
		NomAffiche:  "Voleur",
		Description: "DPS agile avec attaques critiques",
		StatsBase:   statsVoleur,
		Competences: []string{"attaque-sournoise", "esquive"},
	}

	// Clerc - Support/Heal
	statsClerc, _ := shared.NewStats(100, 100, 80, 12, 20, 25, 22, 9, 3, 88)
	cc.templates[ClasseClerc] = &PersonnageTemplate{
		Classe:      ClasseClerc,
		NomAffiche:  "Clerc",
		Description: "Support avec soins et buffs",
		StatsBase:   statsClerc,
		Competences: []string{"soin", "benediction", "resurrection"},
	}
}

// GetClassesDisponibles retourne la liste des classes disponibles
func (cc *CharacterCreator) GetClassesDisponibles() []ClassePersonnage {
	classes := make([]ClassePersonnage, 0, len(cc.templates))
	for classe := range cc.templates {
		classes = append(classes, classe)
	}
	return classes
}

// GetTemplate retourne le template d'une classe
func (cc *CharacterCreator) GetTemplate(classe ClassePersonnage) (*PersonnageTemplate, error) {
	template, exists := cc.templates[classe]
	if !exists {
		return nil, fmt.Errorf("classe %s non trouvée", classe)
	}
	return template, nil
}

// CreerPersonnage crée un nouveau personnage basé sur une classe
func (cc *CharacterCreator) CreerPersonnage(
	id UnitID,
	nom string,
	classe ClassePersonnage,
	teamID TeamID,
	position *shared.Position,
) (*Unite, error) {
	template, err := cc.GetTemplate(classe)
	if err != nil {
		return nil, err
	}

	// Créer une copie des stats pour éviter les modifications partagées
	stats, _ := shared.NewStats(
		template.StatsBase.HP,
		template.StatsBase.MP,
		template.StatsBase.Stamina,
		template.StatsBase.ATK,
		template.StatsBase.DEF,
		template.StatsBase.MATK,
		template.StatsBase.MDEF,
		template.StatsBase.SPD,
		template.StatsBase.MOV,
		template.StatsBase.ATH,
	)

	// Créer l'unité
	unite := NewUnite(id, nom, teamID, stats, position)

	// Ajouter les compétences de base
	for _, compID := range template.Competences {
		competence := cc.creerCompetenceParID(compID)
		if competence != nil {
			unite.AjouterCompetence(competence)
		}
	}

	return unite, nil
}

// creerCompetenceParID crée une compétence basée sur son ID
func (cc *CharacterCreator) creerCompetenceParID(compID string) *Competence {
	switch compID {
	case "charge":
		return NewCompetence(
			CompetenceID("charge"), "Charge", "Attaque de mêlée puissante",
			CompetenceAttaque, 1, ZoneEffet{}, 0, 15, 2,
			20, 1.2, CibleEnnemis,
		)
	case "defense-renforcee":
		return NewCompetence(
			CompetenceID("defense-renforcee"), "Défense Renforcée", "Augmente la défense temporairement",
			CompetenceBuff, 0, ZoneEffet{}, 0, 20, 3,
			0, 1.0, CibleSoi,
		)
	case "provocation":
		return NewCompetence(
			CompetenceID("provocation"), "Provocation", "Force les ennemis à attaquer le Paladin",
			CompetenceBuff, 2, ZoneEffet{}, 10, 0, 3,
			0, 1.0, CibleEnnemis,
		)
	case "soin-divin":
		return NewCompetence(
			CompetenceID("soin-divin"), "Soin Divin", "Soigne un allié",
			CompetenceSoin, 3, ZoneEffet{}, 15, 0, 1,
			30, 1.5, CibleAllies,
		)
	case "tir-precision":
		return NewCompetence(
			CompetenceID("tir-precision"), "Tir de Précision", "Tir ultra-précis",
			CompetenceAttaque, 6, ZoneEffet{}, 0, 25, 2,
			25, 1.3, CibleEnnemis,
		)
	case "fireball":
		return NewCompetence(
			CompetenceID("fireball"), "Boule de Feu", "Projectile enflammé",
			CompetenceMagie, 5, ZoneEffet{}, 25, 0, 2,
			40, 1.6, CibleEnnemis,
		)
	case "lightning":
		return NewCompetence(
			CompetenceID("lightning"), "Éclair", "Frappe éclair rapide",
			CompetenceMagie, 7, ZoneEffet{}, 20, 0, 1,
			28, 1.4, CibleEnnemis,
		)
	default:
		return nil
	}
}

// PersonnageRequest structure pour les requêtes de création
type PersonnageRequest struct {
	Nom      string           `json:"nom"`
	Classe   ClassePersonnage `json:"classe"`
	TeamID   string           `json:"team_id,omitempty"`
	JoueurID string           `json:"joueur_id,omitempty"`
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
}

// PersonnageResponse structure pour les réponses
type PersonnageResponse struct {
	ID          string           `json:"id"`
	Nom         string           `json:"nom"`
	Classe      ClassePersonnage `json:"classe"`
	TeamID      string           `json:"team_id"`
	Stats       *shared.Stats    `json:"stats"`
	Position    *shared.Position `json:"position"`
	Competences []string         `json:"competences"`
}

// ToResponse convertit une unité en réponse API
func (cc *CharacterCreator) ToResponse(unite *Unite, classe ClassePersonnage) *PersonnageResponse {
	compIDs := make([]string, 0, len(unite.Competences()))
	for _, comp := range unite.Competences() {
		compIDs = append(compIDs, string(comp.ID()))
	}

	return &PersonnageResponse{
		ID:          string(unite.ID()),
		Nom:         unite.Nom(),
		Classe:      classe,
		TeamID:      string(unite.TeamID()),
		Stats:       unite.Stats(),
		Position:    unite.Position(),
		Competences: compIDs,
	}
}
