package domain

import (
	"errors"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// UnitID est l'identifiant unique d'une unité
type UnitID string

// Unite représente un combattant dans le combat
type Unite struct {
	id                 UnitID
	nom                string
	teamID             TeamID
	stats              *shared.Stats
	statsActuelles     *shared.Stats // Stats actuelles (modifiées par buffs/debuffs)
	position           *shared.Position
	competences        []*Competence
	statuts            []*shared.Statut
	inventaire         []shared.ObjetID
	estEliminee        bool
	deplacementRestant int
	actionsRestantes   int
}

// NewUnite crée une nouvelle unité
func NewUnite(id UnitID, nom string, teamID TeamID, stats *shared.Stats, position *shared.Position) *Unite {
	return &Unite{
		id:                 id,
		nom:                nom,
		teamID:             teamID,
		stats:              stats,
		statsActuelles:     stats.Clone(),
		position:           position,
		competences:        make([]*Competence, 0),
		statuts:            make([]*shared.Statut, 0),
		inventaire:         make([]shared.ObjetID, 0),
		estEliminee:        false,
		deplacementRestant: stats.MOV,
		actionsRestantes:   1, // 1 action par tour par défaut
	}
}

// Getters
func (u *Unite) ID() UnitID                    { return u.id }
func (u *Unite) Nom() string                   { return u.nom }
func (u *Unite) TeamID() TeamID                { return u.teamID }
func (u *Unite) Stats() *shared.Stats          { return u.stats }
func (u *Unite) StatsActuelles() *shared.Stats { return u.statsActuelles }
func (u *Unite) Position() *shared.Position    { return u.position }
func (u *Unite) Competences() []*Competence    { return u.competences }
func (u *Unite) Statuts() []*shared.Statut     { return u.statuts }
func (u *Unite) EstEliminee() bool             { return u.estEliminee }

// PeutAgir vérifie si l'unité peut effectuer une action
func (u *Unite) PeutAgir() bool {
	if u.estEliminee {
		return false
	}

	// Vérifier les statuts bloquants (Stun, Sleep, etc.)
	for _, statut := range u.statuts {
		if statut.BloqueActions() {
			return true
		}
	}

	return u.actionsRestantes > 0
}

// PeutSeDeplacer vérifie si l'unité peut se déplacer
func (u *Unite) PeutSeDeplacer() bool {
	if u.estEliminee {
		return false
	}

	// Vérifier les statuts bloquants (Root, Stun, etc.)
	for _, statut := range u.statuts {
		if statut.BloqueDeplacement() {
			return false
		}
	}

	return u.deplacementRestant > 0
}

// RecevoirDegats applique des dégâts à l'unité
func (u *Unite) RecevoirDegats(degats int) {
	// Soustraire de HP
	u.statsActuelles.HP -= degats

	// Vérifier si éliminée
	if u.statsActuelles.HP <= 0 {
		u.statsActuelles.HP = 0
		u.estEliminee = true
	}
}

// RecevoirSoin applique un soin à l'unité
func (u *Unite) RecevoirSoin(soin int) {
	if u.estEliminee {
		return // Pas de soin si mort
	}

	u.statsActuelles.HP += soin

	// Cap aux HP max
	if u.statsActuelles.HP > u.stats.HP {
		u.statsActuelles.HP = u.stats.HP
	}
}

// AjouterStatut ajoute un statut à l'unité
func (u *Unite) AjouterStatut(statut *shared.Statut) error {
	if u.estEliminee {
		return errors.New("unité éliminée")
	}

	// Vérifier si le statut existe déjà
	for _, s := range u.statuts {
		if s.Type() == statut.Type() {
			// Refresh ou stack
			return s.Refresh(statut.Duree())
		}
	}

	// Ajouter nouveau statut
	u.statuts = append(u.statuts, statut)

	// Appliquer l'effet initial
	statut.AppliquerEffet(u)

	return nil
}

// RetirerStatut retire un statut de l'unité
func (u *Unite) RetirerStatut(typeStatut shared.TypeStatut) {
	for i, statut := range u.statuts {
		if statut.Type() == typeStatut {
			// Retirer l'effet
			statut.RetirerEffet(u)

			// Retirer du slice
			u.statuts = append(u.statuts[:i], u.statuts[i+1:]...)
			return
		}
	}
}

// TraiterStatuts traite tous les statuts actifs (début de tour)
func (u *Unite) TraiterStatuts() []shared.EffetStatut {
	effets := make([]shared.EffetStatut, 0)

	// Traiter chaque statut
	for i := len(u.statuts) - 1; i >= 0; i-- {
		statut := u.statuts[i]

		// Appliquer l'effet périodique
		effet := statut.AppliquerEffetPeriodique(u)
		if effet != nil {
			effets = append(effets, *effet)
		}

		// Décrémenter la durée
		statut.DecrémenterDuree()

		// Retirer si expiré
		if statut.EstExpire() {
			statut.RetirerEffet(u)
			u.statuts = append(u.statuts[:i], u.statuts[i+1:]...)
		}
	}

	return effets
}

// AjouterCompetence ajoute une compétence à l'unité
func (u *Unite) AjouterCompetence(comp *Competence) error {
	// Vérifier que la compétence n'existe pas déjà
	for _, c := range u.competences {
		if c.ID() == comp.ID() {
			return errors.New("compétence déjà apprise")
		}
	}

	u.competences = append(u.competences, comp)
	return nil
}

// ObtenirCompetence récupère une compétence par ID
func (u *Unite) ObtenirCompetence(id CompetenceID) *Competence {
	for _, comp := range u.competences {
		if comp.ID() == id {
			return comp
		}
	}
	return nil
}

// PeutUtiliserCompetence vérifie si l'unité peut utiliser une compétence
func (u *Unite) PeutUtiliserCompetence(compID CompetenceID) bool {
	comp := u.ObtenirCompetence(compID)
	if comp == nil {
		return false
	}

	// Vérifier les coûts
	if comp.CoutMP() > u.statsActuelles.MP {
		return false
	}

	if comp.CoutStamina() > u.statsActuelles.Stamina {
		return false
	}

	// Vérifier le cooldown
	if comp.EstEnCooldown() {
		return false
	}

	return true
}

// UtiliserCompetence utilise une compétence
func (u *Unite) UtiliserCompetence(compID CompetenceID) error {
	comp := u.ObtenirCompetence(compID)
	if comp == nil {
		return errors.New("compétence inconnue")
	}

	if !u.PeutUtiliserCompetence(compID) {
		return errors.New("impossible d'utiliser cette compétence")
	}

	// Déduire les coûts
	u.statsActuelles.MP -= comp.CoutMP()
	u.statsActuelles.Stamina -= comp.CoutStamina()

	// Activer le cooldown
	comp.ActiverCooldown()

	return nil
}

// SeDeplacer déplace l'unité vers une nouvelle position
func (u *Unite) SeDeplacer(nouvellePosition *shared.Position, coutDeplacement int) error {
	if !u.PeutSeDeplacer() {
		return errors.New("déplacement impossible")
	}

	if coutDeplacement > u.deplacementRestant {
		return errors.New("pas assez de mouvement restant")
	}

	u.position = nouvellePosition
	u.deplacementRestant -= coutDeplacement

	return nil
}

// RegenererStatut régénère les stats périodiquement
func (u *Unite) RegenererStatut() {
	if u.estEliminee {
		return
	}

	// Régénération MP (exemple: 10% par tour)
	regenMP := u.stats.MP / 10
	u.statsActuelles.MP += regenMP
	if u.statsActuelles.MP > u.stats.MP {
		u.statsActuelles.MP = u.stats.MP
	}

	// Régénération Stamina (exemple: 20% par tour)
	regenStamina := u.stats.Stamina / 5
	u.statsActuelles.Stamina += regenStamina
	if u.statsActuelles.Stamina > u.stats.Stamina {
		u.statsActuelles.Stamina = u.stats.Stamina
	}
}

// NouveauTour réinitialise les compteurs de tour
func (u *Unite) NouveauTour() {
	u.actionsRestantes = 1
	u.deplacementRestant = u.statsActuelles.MOV

	// Traiter les statuts
	u.TraiterStatuts()

	// Régénération
	u.RegenererStatut()

	// Décrémenter cooldowns des compétences
	for _, comp := range u.competences {
		comp.DecrémenterCooldown()
	}
}

// AppliquerModificateurStat applique un modificateur temporaire à une stat
func (u *Unite) AppliquerModificateurStat(modificateur *shared.ModificateurStat) {
	switch modificateur.Stat {
	case "ATK":
		u.statsActuelles.ATK += modificateur.Valeur
	case "DEF":
		u.statsActuelles.DEF += modificateur.Valeur
	case "MATK":
		u.statsActuelles.MATK += modificateur.Valeur
	case "MDEF":
		u.statsActuelles.MDEF += modificateur.Valeur
	case "SPD":
		u.statsActuelles.SPD += modificateur.Valeur
	case "MOV":
		u.statsActuelles.MOV += modificateur.Valeur
	}
}

// RetirerModificateurStat retire un modificateur temporaire
func (u *Unite) RetirerModificateurStat(modificateur *shared.ModificateurStat) {
	switch modificateur.Stat {
	case "ATK":
		u.statsActuelles.ATK -= modificateur.Valeur
	case "DEF":
		u.statsActuelles.DEF -= modificateur.Valeur
	case "MATK":
		u.statsActuelles.MATK -= modificateur.Valeur
	case "MDEF":
		u.statsActuelles.MDEF -= modificateur.Valeur
	case "SPD":
		u.statsActuelles.SPD -= modificateur.Valeur
	case "MOV":
		u.statsActuelles.MOV -= modificateur.Valeur
	}
}

// RecalculerStats recalcule les stats actuelles en appliquant tous les modificateurs
func (u *Unite) RecalculerStats() {
	// Partir des stats de base
	u.statsActuelles = u.stats.Clone()

	// Appliquer tous les modificateurs des statuts
	for _, statut := range u.statuts {
		for _, mod := range statut.Modificateurs() {
			u.AppliquerModificateurStat(&mod)
		}
	}
}

// HPActuels retourne les HP actuels de l'unité
func (u *Unite) HPActuels() int {
	return u.statsActuelles.HP
}

// ConsommerMP consomme des points de magie
func (u *Unite) ConsommerMP(montant int) {
	u.statsActuelles.MP -= montant
	if u.statsActuelles.MP < 0 {
		u.statsActuelles.MP = 0
	}
}

// ConsommerStamina consomme de l'endurance
func (u *Unite) ConsommerStamina(montant int) {
	u.statsActuelles.Stamina -= montant
	if u.statsActuelles.Stamina < 0 {
		u.statsActuelles.Stamina = 0
	}
}

// ObtenirCompetenceParDefaut retourne l'attaque basique de l'unité
func (u *Unite) ObtenirCompetenceParDefaut() *Competence {
	// Attaque basique par défaut (physique)
	// TODO: Créer une vraie compétence par défaut dans NewUnite
	attaqueBasique := NewCompetence(
		"attaque-basique",
		"Attaque Basique",
		"Attaque physique de base",
		CompetenceAttaque,
		1, // Portée 1
		ZoneEffet{forme: ZoneSingle, taille: 1},
		0,   // Pas de coût MP
		0,   // Pas de coût Stamina
		1,   // Cooldown 1
		10,  // 10 dégâts de base
		0.5, // Scaling 50% ATK
		CibleEnnemis,
	)

	// Ajouter l'effet de dégâts
	attaqueBasique.AjouterEffet(EffetCompetence{
		typeEffet: EffetDegats,
		valeur:    10,
		duree:     0,
		statut:    nil,
	})

	return attaqueBasique
}

// AppliquerStatut applique un statut à l'unité (alias de AjouterStatut)
func (u *Unite) AppliquerStatut(statut *shared.Statut) error {
	return u.AjouterStatut(statut)
}
