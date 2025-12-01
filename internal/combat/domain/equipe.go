package domain

import (
	"errors"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// TeamID est l'identifiant unique d'une équipe
type TeamID string

// Equipe représente un groupe d'unités combattant ensemble
type Equipe struct {
	id       TeamID
	nom      string
	couleur  string // Hex color (#FF0000)
	membres  []*Unite
	isIA     bool
	joueurID *string // nil si IA
}

// NewEquipe crée une nouvelle équipe
func NewEquipe(id TeamID, nom string, couleur string, isIA bool, joueurID *string) (*Equipe, error) {
	if nom == "" {
		return nil, errors.New("le nom ne peut pas être vide")
	}

	// Vérifier que joueurID est fourni si pas IA
	if !isIA && joueurID == nil {
		return nil, errors.New("joueurID requis pour équipe non-IA")
	}

	return &Equipe{
		id:       id,
		nom:      nom,
		couleur:  couleur,
		membres:  make([]*Unite, 0),
		isIA:     isIA,
		joueurID: joueurID,
	}, nil
}

// Getters
func (e *Equipe) ID() TeamID        { return e.id }
func (e *Equipe) Nom() string       { return e.nom }
func (e *Equipe) Couleur() string   { return e.couleur }
func (e *Equipe) Membres() []*Unite { return e.membres }
func (e *Equipe) IsIA() bool        { return e.isIA }
func (e *Equipe) JoueurID() *string { return e.joueurID }

// AjouterMembre ajoute une unité à l'équipe
func (e *Equipe) AjouterMembre(unite *Unite) error {
	// Vérifier que l'unité n'est pas déjà dans l'équipe
	for _, m := range e.membres {
		if m.ID() == unite.ID() {
			return errors.New("unité déjà dans l'équipe")
		}
	}

	// Vérifier que l'unité appartient à cette équipe
	if unite.TeamID() != e.id {
		return errors.New("l'unité n'appartient pas à cette équipe")
	}

	e.membres = append(e.membres, unite)
	return nil
}

// RetirerMembre retire une unité de l'équipe
func (e *Equipe) RetirerMembre(uniteID UnitID) error {
	for i, membre := range e.membres {
		if membre.ID() == uniteID {
			e.membres = append(e.membres[:i], e.membres[i+1:]...)
			return nil
		}
	}
	return errors.New("unité non trouvée")
}

// ObtenirMembre récupère un membre par ID
func (e *Equipe) ObtenirMembre(uniteID UnitID) *Unite {
	for _, membre := range e.membres {
		if membre.ID() == uniteID {
			return membre
		}
	}
	return nil
}

// NombreMembres retourne le nombre de membres dans l'équipe
func (e *Equipe) NombreMembres() int {
	return len(e.membres)
}

// ADesMembresVivants vérifie si l'équipe a au moins un membre vivant
func (e *Equipe) ADesMembresVivants() bool {
	for _, membre := range e.membres {
		if !membre.EstEliminee() {
			return true
		}
	}
	return false
}

// MembresVivants retourne tous les membres vivants
func (e *Equipe) MembresVivants() []*Unite {
	vivants := make([]*Unite, 0)
	for _, membre := range e.membres {
		if !membre.EstEliminee() {
			vivants = append(vivants, membre)
		}
	}
	return vivants
}

// MembresElimines retourne tous les membres éliminés
func (e *Equipe) MembresElimines() []*Unite {
	elimines := make([]*Unite, 0)
	for _, membre := range e.membres {
		if membre.EstEliminee() {
			elimines = append(elimines, membre)
		}
	}
	return elimines
}

// StatsMoyennes calcule les stats moyennes de l'équipe
func (e *Equipe) StatsMoyennes() *shared.Stats {
	if len(e.membres) == 0 {
		return &shared.Stats{}
	}

	somme := &shared.Stats{}
	for _, membre := range e.membres {
		stats := membre.Stats()
		somme.HP += stats.HP
		somme.MP += stats.MP
		somme.Stamina += stats.Stamina
		somme.ATK += stats.ATK
		somme.DEF += stats.DEF
		somme.MATK += stats.MATK
		somme.MDEF += stats.MDEF
		somme.SPD += stats.SPD
		somme.MOV += stats.MOV
	}

	count := len(e.membres)
	return &shared.Stats{
		HP:      somme.HP / count,
		MP:      somme.MP / count,
		Stamina: somme.Stamina / count,
		ATK:     somme.ATK / count,
		DEF:     somme.DEF / count,
		MATK:    somme.MATK / count,
		MDEF:    somme.MDEF / count,
		SPD:     somme.SPD / count,
		MOV:     somme.MOV / count,
	}
}

// PuissanceTotale calcule la puissance totale de l'équipe
func (e *Equipe) PuissanceTotale() int {
	total := 0
	for _, membre := range e.membres {
		stats := membre.Stats()
		// Formule simple: somme de toutes les stats
		total += stats.HP + stats.MP + stats.Stamina + stats.ATK +
			stats.DEF + stats.MATK + stats.MDEF + stats.SPD + stats.MOV
	}
	return total
}

// EstComplete vérifie si l'équipe a au moins un membre
func (e *Equipe) EstComplete() bool {
	return len(e.membres) > 0
}

// EstEnnemie vérifie si l'équipe est ennemie d'une autre
func (e *Equipe) EstEnnemie(autre *Equipe) bool {
	return e.id != autre.id
}

// ContientUnite vérifie si l'équipe contient une unité donnée
func (e *Equipe) ContientUnite(uniteID UnitID) bool {
	return e.ObtenirMembre(uniteID) != nil
}
