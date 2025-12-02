package domain

import (
	"errors"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// UnitID est l'identifiant unique d'une unité
type UnitID string

// Unite représente un combattant dans le combat
// Composition Pattern - Délègue aux composants spécialisés
// Single Responsibility Principle - Coordonne, ne gère pas directement
type Unite struct {
	// Identité
	id     UnitID
	nom    string
	teamID TeamID

	// Position
	position *shared.Position

	// Composants (Composition Pattern)
	combat    *UnitCombatBehavior
	statuses  *UnitStatusManager
	inventory *UnitInventory

	// État du tour
	deplacementRestant int
	actionsRestantes   int
}

// NewUnite crée une nouvelle unité
// Utilise la Composition pour déléguer aux composants spécialisés
func NewUnite(id UnitID, nom string, teamID TeamID, stats *shared.Stats, position *shared.Position) *Unite {
	return &Unite{
		id:       id,
		nom:      nom,
		teamID:   teamID,
		position: position,

		// Composition Pattern - Créer les composants
		combat:    NewUnitCombatBehavior(stats),
		statuses:  NewUnitStatusManager(),
		inventory: NewUnitInventory(),

		// État initial du tour
		deplacementRestant: stats.MOV,
		actionsRestantes:   1,
	}
}

// Getters basiques
func (u *Unite) ID() UnitID                 { return u.id }
func (u *Unite) Nom() string                { return u.nom }
func (u *Unite) TeamID() TeamID             { return u.teamID }
func (u *Unite) Position() *shared.Position { return u.position }

// Getters délégués aux composants (Composition Pattern)
func (u *Unite) Stats() *shared.Stats          { return u.combat.BaseStats() }
func (u *Unite) StatsActuelles() *shared.Stats { return u.combat.CurrentStats() }
func (u *Unite) HPActuels() int                { return u.combat.CurrentHP() }
func (u *Unite) Competences() []*Competence    { return u.inventory.Skills() }
func (u *Unite) Statuts() []*shared.Statut     { return u.statuses.Statuses() }
func (u *Unite) EstEliminee() bool             { return u.combat.IsEliminated() }

// PeutAgir vérifie si l'unité peut effectuer une action
func (u *Unite) PeutAgir() bool {
	if u.combat.IsEliminated() {
		return false
	}

	// Déléguer aux composants
	if u.statuses.BlocksActions() {
		return false
	}

	return u.actionsRestantes > 0
}

// PeutSeDeplacer vérifie si l'unité peut se déplacer
func (u *Unite) PeutSeDeplacer() bool {
	if u.combat.IsEliminated() {
		return false
	}

	// Déléguer aux composants
	if u.statuses.BlocksMovement() {
		return false
	}

	return u.deplacementRestant > 0
}

// EstBloqueDeplacement vérifie si l'unité est bloquée pour le déplacement
func (u *Unite) EstBloqueDeplacement() bool {
	if u.combat.IsEliminated() {
		return true
	}

	// Déléguer au gestionnaire de statuts
	return u.statuses.BlocksMovement()
}

// DeplacerVers déplace l'unité vers une position spécifique
func (u *Unite) DeplacerVers(nouvellePosition *shared.Position) {
	u.position = nouvellePosition
}

// RecevoirDegats applique des dégâts à l'unité (délègue au composant combat)
func (u *Unite) RecevoirDegats(degats int) {
	u.combat.TakeDamage(degats)
}

// RecevoirSoin applique un soin à l'unité (délègue au composant combat)
func (u *Unite) RecevoirSoin(soin int) {
	u.combat.Heal(soin)
}

// Soigner est un alias de RecevoirSoin (pour compatibilité)
func (u *Unite) Soigner(soin int) {
	u.combat.Heal(soin)
}

// AjouterStatut ajoute un statut à l'unité (délègue au composant)
func (u *Unite) AjouterStatut(statut *shared.Statut) error {
	if u.combat.IsEliminated() {
		return errors.New("unité éliminée")
	}

	// Déléguer au gestionnaire de statuts
	return u.statuses.AddStatus(statut)
}

// RetirerStatut retire un statut de l'unité (délègue au composant)
func (u *Unite) RetirerStatut(typeStatut shared.TypeStatut) {
	u.statuses.RemoveStatus(typeStatut)
}

// TraiterStatuts traite tous les statuts actifs (délègue au composant)
func (u *Unite) TraiterStatuts() []shared.EffetStatut {
	// Déléguer au gestionnaire de statuts
	return u.statuses.ProcessStatuses(u)
}

// AjouterCompetence ajoute une compétence à l'unité (délègue au composant)
func (u *Unite) AjouterCompetence(comp *Competence) error {
	return u.inventory.AddSkill(comp)
}

// ObtenirCompetence récupère une compétence par ID (délègue au composant)
func (u *Unite) ObtenirCompetence(id CompetenceID) *Competence {
	return u.inventory.GetSkill(id)
}

// PeutUtiliserCompetence vérifie si l'unité peut utiliser une compétence
func (u *Unite) PeutUtiliserCompetence(compID CompetenceID) bool {
	comp := u.inventory.GetSkill(compID)
	if comp == nil {
		return false
	}

	// Vérifier les coûts via le composant combat
	if comp.CoutMP() > u.combat.CurrentMP() {
		return false
	}

	if comp.CoutStamina() > u.combat.CurrentStats().Stamina {
		return false
	}

	// Vérifier le cooldown via le composant inventory
	return u.inventory.IsSkillReady(compID)
}

// UtiliserCompetence utilise une compétence
func (u *Unite) UtiliserCompetence(compID CompetenceID) error {
	if !u.PeutUtiliserCompetence(compID) {
		return errors.New("impossible d'utiliser cette compétence")
	}

	comp := u.inventory.GetSkill(compID)
	if comp == nil {
		return errors.New("compétence inconnue")
	}

	// Déduire les coûts via le composant combat
	if err := u.combat.ConsumeMP(comp.CoutMP()); err != nil {
		return err
	}
	if err := u.combat.ConsumeStamina(comp.CoutStamina()); err != nil {
		return err
	}

	// Activer le cooldown via le composant inventory
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
	if u.combat.IsEliminated() {
		return
	}

	// Régénération MP (exemple: 10% par tour)
	baseStats := u.combat.BaseStats()
	regenMP := baseStats.MP / 10
	u.combat.RestoreMP(regenMP)

	// Régénération Stamina (exemple: 20% par tour)
	// Note: Stamina sera géré via combat component
	regenStamina := baseStats.Stamina / 5
	currentStats := u.combat.CurrentStats()
	currentStats.Stamina += regenStamina
	if currentStats.Stamina > baseStats.Stamina {
		currentStats.Stamina = baseStats.Stamina
	}
}

// NouveauTour réinitialise les compteurs de tour
func (u *Unite) NouveauTour() {
	u.actionsRestantes = 1
	u.deplacementRestant = u.combat.CurrentStats().MOV

	// Traiter les statuts via composant
	u.TraiterStatuts()

	// Décrémenter les cooldowns
	u.inventory.DecrementAllCooldowns()

	// Régénération
	u.RegenererStatut()

	// Décrémenter cooldowns des compétences
	u.inventory.DecrementAllCooldowns()
}

// AppliquerModificateurStat applique un modificateur temporaire à une stat
func (u *Unite) AppliquerModificateurStat(modificateur *shared.ModificateurStat) {
	stats := u.combat.CurrentStats()
	switch modificateur.Stat {
	case "ATK":
		stats.ATK += modificateur.Valeur
	case "DEF":
		stats.DEF += modificateur.Valeur
	case "MATK":
		stats.MATK += modificateur.Valeur
	case "MDEF":
		stats.MDEF += modificateur.Valeur
	case "SPD":
		stats.SPD += modificateur.Valeur
	case "MOV":
		stats.MOV += modificateur.Valeur
	}
}

// RetirerModificateurStat retire un modificateur temporaire
func (u *Unite) RetirerModificateurStat(modificateur *shared.ModificateurStat) {
	stats := u.combat.CurrentStats()
	switch modificateur.Stat {
	case "ATK":
		stats.ATK -= modificateur.Valeur
	case "DEF":
		stats.DEF -= modificateur.Valeur
	case "MATK":
		stats.MATK -= modificateur.Valeur
	case "MDEF":
		stats.MDEF -= modificateur.Valeur
	case "SPD":
		stats.SPD -= modificateur.Valeur
	case "MOV":
		stats.MOV -= modificateur.Valeur
	}
}

// RecalculerStats recalcule les stats actuelles en appliquant tous les modificateurs
func (u *Unite) RecalculerStats() {
	// Réinitialiser aux stats de base
	baseStats := u.combat.BaseStats()
	currentStats := u.combat.CurrentStats()

	// Copier les stats de base
	currentStats.HP = baseStats.HP
	currentStats.MP = baseStats.MP
	currentStats.ATK = baseStats.ATK
	currentStats.DEF = baseStats.DEF
	currentStats.SPD = baseStats.SPD
	currentStats.MATK = baseStats.MATK
	currentStats.MDEF = baseStats.MDEF
	currentStats.MOV = baseStats.MOV
	currentStats.Stamina = baseStats.Stamina

	// Appliquer tous les modificateurs des statuts
	for _, statut := range u.statuses.Statuses() {
		for _, mod := range statut.Modificateurs() {
			u.AppliquerModificateurStat(&mod)
		}
	}
}

// HPActuels est déjà défini dans les getters (ligne 63)

// ConsommerMP consomme des points de magie (délègue au composant)
func (u *Unite) ConsommerMP(montant int) error {
	return u.combat.ConsumeMP(montant)
}

// ConsommerStamina consomme de l'endurance (délègue au composant)
func (u *Unite) ConsommerStamina(montant int) error {
	return u.combat.ConsumeStamina(montant)
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

// Step C - Méthodes additionnelles pour les patterns

// SetHP définit les HP actuels (utilisé pour rollback des commandes)
func (u *Unite) SetHP(hp int) {
	currentStats := u.combat.CurrentStats()
	baseStats := u.combat.BaseStats()

	currentStats.HP = hp
	if currentStats.HP < 0 {
		currentStats.HP = 0
	}
	if currentStats.HP > baseStats.HP {
		currentStats.HP = baseStats.HP
	}

	// Mettre à jour le statut d'élimination via le combat component
	if currentStats.HP == 0 {
		u.combat.TakeDamage(0) // Déclenche la logique KO
	}
}

// SetMP définit les MP actuels (utilisé pour rollback des commandes)
func (u *Unite) SetMP(mp int) {
	currentStats := u.combat.CurrentStats()
	baseStats := u.combat.BaseStats()

	currentStats.MP = mp
	if currentStats.MP < 0 {
		currentStats.MP = 0
	}
	if currentStats.MP > baseStats.MP {
		currentStats.MP = baseStats.MP
	}
}

// RestaurerMP restaure des points de magie
func (u *Unite) RestaurerMP(montant int) {
	currentStats := u.combat.CurrentStats()
	baseStats := u.combat.BaseStats()

	currentStats.MP += montant
	if currentStats.MP > baseStats.MP {
		currentStats.MP = baseStats.MP
	}
}

// Ressusciter ressuscite l'unité avec un montant de HP
func (u *Unite) Ressusciter(hp int) {
	if !u.combat.IsEliminated() {
		return // Déjà vivante
	}

	// Utiliser la méthode Revive du combat component
	u.combat.Revive(hp)

	// Retirer le statut "Dead" si présent
	u.RetirerStatut(shared.TypeStatutMort)
}

// EstIA vérifie si l'unité est contrôlée par l'IA
func (u *Unite) EstIA() bool {
	// TODO: Ajouter un champ isAI dans la structure Unite
	// Pour l'instant, on peut considérer que toutes les unités ennemies sont IA
	// ou utiliser une map externe dans Combat
	return false // Placeholder
}

// IAChoisirAction fait choisir une action à l'IA
func (u *Unite) IAChoisirAction(combat interface{}) interface{} {
	// TODO: Implémenter logique IA (système d'IA sera dans Step D ou E)
	// Pour l'instant, retourner WaitCommand par défaut
	return nil // Placeholder
}

// EstSilence vérifie si l'unité est sous l'effet Silence (bloque compétences)
func (u *Unite) EstSilence() bool {
	return u.statuses.IsSilenced()
}

// EstStun vérifie si l'unité est sous l'effet Stun (bloque toutes actions)
func (u *Unite) EstStun() bool {
	return u.statuses.IsStunned()
}

// EstRoot vérifie si l'unité est sous l'effet Root (bloque déplacement)
func (u *Unite) EstRoot() bool {
	return u.statuses.IsRooted()
}

// EstEmpoisonne vérifie si l'unité est empoisonnée
func (u *Unite) EstEmpoisonne() bool {
	return u.statuses.IsPoisoned()
}

// SkillEstPret vérifie si une compétence est prête (pas en cooldown et ressources suffisantes)
func (u *Unite) SkillEstPret(skillID CompetenceID) bool {
	return u.PeutUtiliserCompetence(skillID)
}

// ActiverCooldown active le cooldown d'une compétence
func (u *Unite) ActiverCooldown(skillID CompetenceID, duree int) {
	comp := u.ObtenirCompetence(skillID)
	if comp != nil {
		comp.SetCooldownActuel(duree)
	}
}

// AppliquerBuff applique un buff temporaire sur une stat
func (u *Unite) AppliquerBuff(stat string, valeur int, duree int) {
	// Créer un modificateur de stat
	mod := shared.ModificateurStat{
		Stat:   stat,
		Valeur: valeur,
	}

	// Appliquer immédiatement
	u.AppliquerModificateurStat(&mod)

	// TODO: Gérer la durée avec un système de buffs temporaires
	// Pour l'instant, on applique directement sur les stats actuelles
}
