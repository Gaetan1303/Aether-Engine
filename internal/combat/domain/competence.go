package domain

import (
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// CompetenceID est l'identifiant unique d'une compétence
type CompetenceID string

// Competence représente une compétence utilisable par une unité
type Competence struct {
	id             CompetenceID
	nom            string
	description    string
	typeCompetence TypeCompetence
	portee         int
	zone           ZoneEffet
	coutMP         int
	coutStamina    int
	cooldown       int
	cooldownActuel int
	degatsBase     int
	modificateur   float64 // Scaling (ATK, MATK, etc.)
	effets         []EffetCompetence
	cibles         TypeCible
}

// TypeCompetence énumère les types de compétences
type TypeCompetence int

const (
	CompetenceAttaque TypeCompetence = iota
	CompetenceMagie
	CompetenceSoin
	CompetenceBuff
	CompetenceDebuff
	CompetenceUtilitaire
	CompetenceInvocation
)

// ZoneEffet définit la zone d'effet d'une compétence
type ZoneEffet struct {
	forme  FormeZone
	taille int // Rayon ou dimension
}

// FormeZone énumère les formes de zones
type FormeZone int

const (
	ZoneSingle FormeZone = iota // Cible unique
	ZoneCone                    // Cône devant le lanceur
	ZoneCercle                  // Cercle autour d'un point
	ZoneLigne                   // Ligne droite
	ZoneCroix                   // Croix (+)
)

// TypeCible définit les cibles valides
type TypeCible int

const (
	CibleEnnemis TypeCible = iota
	CibleAllies
	CibleTous
	CibleSoi
)

// EffetCompetence représente un effet d'une compétence
type EffetCompetence struct {
	typeEffet TypeEffetCompetence
	valeur    int
	duree     int // Pour les statuts
	statut    *shared.TypeStatut
}

// Getters pour EffetCompetence
func (e *EffetCompetence) TypeEffet() TypeEffetCompetence { return e.typeEffet }
func (e *EffetCompetence) Valeur() int                    { return e.valeur }
func (e *EffetCompetence) Duree() int                     { return e.duree }
func (e *EffetCompetence) StatutType() *shared.TypeStatut { return e.statut }

// TypeEffetCompetence énumère les types d'effets
type TypeEffetCompetence int

const (
	EffetDegats TypeEffetCompetence = iota
	EffetSoin
	EffetStatut
	EffetDeplacement
	EffetInvocation
	EffetModificateurStat
)

// NewCompetence crée une nouvelle compétence
func NewCompetence(
	id CompetenceID,
	nom, description string,
	typeCompetence TypeCompetence,
	portee int,
	zone ZoneEffet,
	coutMP, coutStamina, cooldown int,
	degatsBase int,
	modificateur float64,
	cibles TypeCible,
) *Competence {
	return &Competence{
		id:             id,
		nom:            nom,
		description:    description,
		typeCompetence: typeCompetence,
		portee:         portee,
		zone:           zone,
		coutMP:         coutMP,
		coutStamina:    coutStamina,
		cooldown:       cooldown,
		cooldownActuel: 0,
		degatsBase:     degatsBase,
		modificateur:   modificateur,
		effets:         make([]EffetCompetence, 0),
		cibles:         cibles,
	}
}

// Getters
func (c *Competence) ID() CompetenceID          { return c.id }
func (c *Competence) Nom() string               { return c.nom }
func (c *Competence) Description() string       { return c.description }
func (c *Competence) Type() TypeCompetence      { return c.typeCompetence }
func (c *Competence) Portee() int               { return c.portee }
func (c *Competence) Zone() ZoneEffet           { return c.zone }
func (c *Competence) CoutMP() int               { return c.coutMP }
func (c *Competence) CoutStamina() int          { return c.coutStamina }
func (c *Competence) Cooldown() int             { return c.cooldown }
func (c *Competence) CooldownActuel() int       { return c.cooldownActuel }
func (c *Competence) DegatsBase() int           { return c.degatsBase }
func (c *Competence) Modificateur() float64     { return c.modificateur }
func (c *Competence) Effets() []EffetCompetence { return c.effets }
func (c *Competence) Cibles() TypeCible         { return c.cibles }

// AjouterEffet ajoute un effet à la compétence
func (c *Competence) AjouterEffet(effet EffetCompetence) {
	c.effets = append(c.effets, effet)
}

// EstEnCooldown vérifie si la compétence est en cooldown
func (c *Competence) EstEnCooldown() bool {
	return c.cooldownActuel > 0
}

// ActiverCooldown active le cooldown de la compétence
func (c *Competence) ActiverCooldown() {
	c.cooldownActuel = c.cooldown
}

// DecrémenterCooldown décrémente le cooldown
func (c *Competence) DecrémenterCooldown() {
	if c.cooldownActuel > 0 {
		c.cooldownActuel--
	}
}

// SetCooldownActuel définit le cooldown actuel (utilisé par ActiverCooldown dans Unite)
func (c *Competence) SetCooldownActuel(duree int) {
	c.cooldownActuel = duree
}

// CalculerDegats calcule les dégâts totaux en fonction des stats de l'attaquant
func (c *Competence) CalculerDegats(stats *shared.Stats) int {
	// Choisir la stat appropriée selon le type
	var statBase int
	switch c.typeCompetence {
	case CompetenceAttaque:
		statBase = stats.ATK
	case CompetenceMagie:
		statBase = stats.MATK
	default:
		statBase = 0
	}

	// Calculer: degatsBase + (statBase * modificateur)
	degats := float64(c.degatsBase) + (float64(statBase) * c.modificateur)
	return int(degats)
}

// EstCibleValide vérifie si une unité est une cible valide
func (c *Competence) EstCibleValide(lanceur, cible *Unite) bool {
	// Même équipe ?
	memeEquipe := lanceur.TeamID() == cible.TeamID()

	switch c.cibles {
	case CibleEnnemis:
		return !memeEquipe
	case CibleAllies:
		return memeEquipe
	case CibleSoi:
		return lanceur.ID() == cible.ID()
	case CibleTous:
		return true
	default:
		return false
	}
}

// ObtenirPositionsDansZone retourne les positions affectées par la zone d'effet
func (c *Competence) ObtenirPositionsDansZone(centre *shared.Position, grille *shared.GrilleCombat) []*shared.Position {
	positions := make([]*shared.Position, 0)

	switch c.zone.forme {
	case ZoneSingle:
		positions = append(positions, centre)

	case ZoneCercle:
		// Toutes les positions dans le rayon
		positions = grille.PositionsADansPortee(centre, c.zone.taille)

	case ZoneCone:
		// TODO: Implémenter forme cone
		positions = append(positions, centre)

	case ZoneLigne:
		// TODO: Implémenter forme ligne
		positions = append(positions, centre)

	case ZoneCroix:
		// TODO: Implémenter forme croix
		positions = append(positions, centre)
	}

	return positions
}

// Clone crée une copie de la compétence
func (c *Competence) Clone() *Competence {
	clone := *c
	clone.effets = make([]EffetCompetence, len(c.effets))
	copy(clone.effets, c.effets)
	return &clone
}
