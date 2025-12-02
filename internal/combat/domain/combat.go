package domain

import (
	"errors"
	"time"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Combat est l'agrégat racine qui gère une instance de combat tactique
type Combat struct {
	id                string
	etat              EtatCombat
	equipes           map[TeamID]*Equipe
	grille            *shared.GrilleCombat
	ordreDeJeu        []UnitID
	uniteActive       UnitID
	phase             PhaseTour
	tourActuel        int
	actionEnCours     *ActionCombat
	version           int
	evenements        []Evenement // Uncommitted events
	createdAt         time.Time
	updatedAt         time.Time
	damageCalculator  DamageCalculator         // Strategy Pattern - algorithme de dégâts
	calculatorFactory *DamageCalculatorFactory // Factory pour créer strategies

	// Step C - State Machine + Command Pattern + Observer Pattern
	// Utilisation d'interfaces pour éviter les cycles d'imports
	stateMachine    interface{}     // CombatStateMachine interface
	commandInvoker  interface{}     // CommandInvoker interface
	commandFactory  interface{}     // CommandFactory interface
	observerSubject interface{}     // ObserverSubject interface
	validationChain interface{}     // ValidationChain interface
	fuiteAutorisee  bool            // Indique si la fuite est autorisée
	equipesFuites   map[TeamID]bool // Équipes ayant fui le combat
}

// NewCombat crée une nouvelle instance de combat
func NewCombat(id string, equipes []*Equipe, grille *shared.GrilleCombat) (*Combat, error) {
	if len(equipes) < 2 {
		return nil, errors.New("minimum 2 équipes requises")
	}

	combat := &Combat{
		id:                id,
		etat:              EtatAttente,
		equipes:           make(map[TeamID]*Equipe),
		grille:            grille,
		ordreDeJeu:        make([]UnitID, 0),
		phase:             PhasePreparation,
		tourActuel:        0,
		version:           0,
		evenements:        make([]Evenement, 0),
		createdAt:         time.Now(),
		updatedAt:         time.Now(),
		calculatorFactory: NewDamageCalculatorFactory(),  // Factory Pattern
		damageCalculator:  NewPhysicalDamageCalculator(), // Strategy par défaut

		// Step C - Les patterns seront initialisés via CombatInitializer
		fuiteAutorisee: true, // Par défaut, fuite autorisée
		equipesFuites:  make(map[TeamID]bool),
	}

	// Ajouter les équipes
	for _, equipe := range equipes {
		combat.equipes[equipe.ID()] = equipe
	}

	return combat, nil
}

// Getters
func (c *Combat) ID() string                   { return c.id }
func (c *Combat) Etat() EtatCombat             { return c.etat }
func (c *Combat) Grille() *shared.GrilleCombat { return c.grille }
func (c *Combat) TourActuel() int              { return c.tourActuel }
func (c *Combat) Version() int                 { return c.version }
func (c *Combat) UniteActive() UnitID          { return c.uniteActive }
func (c *Combat) Equipes() map[TeamID]*Equipe  { return c.equipes }
func (c *Combat) GetTimestamp() int64          { return c.updatedAt.Unix() }

// Step C - Méthodes simples sans dépendances cycliques
// (Les méthodes complexes sont dans combat_step_c.go)

func (c *Combat) FuiteAutorisee() bool {
	return c.fuiteAutorisee
}

func (c *Combat) SetFuiteAutorisee(autorisee bool) {
	c.fuiteAutorisee = autorisee
}

// Strategy Pattern - Setters pour changer l'algorithme de dégâts
// Permet de changer dynamiquement la stratégie de calcul (Open/Closed Principle)

// SetDamageCalculator change la stratégie de calcul de dégâts
func (c *Combat) SetDamageCalculator(calculator DamageCalculator) {
	c.damageCalculator = calculator
}

// SetPhysicalDamageMode active le mode dégâts physiques
func (c *Combat) SetPhysicalDamageMode() {
	c.damageCalculator = NewPhysicalDamageCalculator()
}

// SetMagicalDamageMode active le mode dégâts magiques
func (c *Combat) SetMagicalDamageMode() {
	c.damageCalculator = NewMagicalDamageCalculator()
}

// SetHybridDamageMode active le mode dégâts hybrides
func (c *Combat) SetHybridDamageMode(physicalRatio, magicalRatio float64) {
	c.damageCalculator = c.calculatorFactory.CreateHybridCalculator(physicalRatio, magicalRatio)
}

// GetDamageCalculator retourne la stratégie actuelle
func (c *Combat) GetDamageCalculator() DamageCalculator {
	return c.damageCalculator
}

// Step C - Méthodes helper publiques (sans imports cycliques)
// Les méthodes complexes avec state machine sont dans combat_step_c.go

// TrouverUnite trouve une unité par son ID (méthode publique)
func (c *Combat) TrouverUnite(id UnitID) *Unite {
	return c.trouverUnite(id)
}

// ObtenirPositionsOccupees retourne les positions occupées (sauf une unité)
func (c *Combat) ObtenirPositionsOccupees(exclusionID UnitID) map[string]bool {
	return c.obtenirPositionsOccupees(exclusionID)
}

// MarquerEquipeFuite marque une équipe comme ayant fui
func (c *Combat) MarquerEquipeFuite(teamID TeamID) {
	c.equipesFuites[teamID] = true
}

// AnnulerFuite annule la fuite d'une équipe (pour rollback)
func (c *Combat) AnnulerFuite(teamID TeamID) {
	delete(c.equipesFuites, teamID)
}

// ObtenirEnnemis retourne tous les ennemis d'une équipe
func (c *Combat) ObtenirEnnemis(teamID TeamID) []*Unite {
	ennemis := make([]*Unite, 0)
	for _, equipe := range c.equipes {
		if equipe.ID() != teamID {
			for _, membre := range equipe.Membres() {
				if !membre.EstEliminee() {
					ennemis = append(ennemis, membre)
				}
			}
		}
	}
	return ennemis
}

// VerifierConditionsVictoire vérifie les conditions de victoire/défaite
func (c *Combat) VerifierConditionsVictoire() string {
	// Vérifier les équipes ayant fui
	equipesActives := 0
	for teamID, equipe := range c.equipes {
		// Si l'équipe a fui, elle n'est pas active
		if c.equipesFuites[teamID] {
			continue
		}
		// Si l'équipe a des membres vivants, elle est active
		if equipe.ADesMembresVivants() {
			equipesActives++
		}
	}

	if equipesActives <= 1 {
		// Vérifier si c'est une fuite
		for teamID := range c.equipesFuites {
			if c.equipesFuites[teamID] {
				return "FLED"
			}
		}
		// Sinon victoire/défaite
		if equipesActives == 1 {
			return "VICTORY"
		}
		return "DEFEAT"
	}

	return "CONTINUE"
}

// ObtenirResultat retourne le résultat du combat
func (c *Combat) ObtenirResultat() string {
	return c.VerifierConditionsVictoire()
}

// DistribuerRecompenses distribue les récompenses (XP, loots)
func (c *Combat) DistribuerRecompenses() {
	// TODO: Implémenter la distribution de récompenses
	// Pour l'instant, juste un placeholder
}

// Méthodes pour le système d'inventaire (Item commands)
func (c *Combat) PossedeObjet(itemID string) bool {
	// TODO: Implémenter vérification inventaire
	return true // Placeholder
}

func (c *Combat) ObtenirQuantiteObjet(itemID string) int {
	// TODO: Implémenter récupération quantité
	return 1 // Placeholder
}

func (c *Combat) ConsommerObjet(itemID string, quantite int) {
	// TODO: Implémenter consommation objet
}

func (c *Combat) AjouterObjet(itemID string, quantite int) {
	// TODO: Implémenter ajout objet
}

func (c *Combat) ObtenirObjet(itemID string) interface{} {
	// TODO: Implémenter récupération objet (retournera *shared.Item quand implémenté)
	return nil // Placeholder
}

// Demarrer démarre le combat et calcule l'ordre d'initiative
// DEPRECATED: Utilisez InitializeCombatWithStateMachine() via CombatInitializer pour la nouvelle state machine (Step C)
func (c *Combat) Demarrer() error {
	if c.etat != EtatAttente {
		return errors.New("le combat doit être en attente pour démarrer")
	}

	// Calculer l'ordre d'initiative basé sur la vitesse
	c.calculerOrdreInitiative()

	c.etat = EtatEnCours
	c.phase = PhaseDebutTour
	c.tourActuel = 1

	if len(c.ordreDeJeu) > 0 {
		c.uniteActive = c.ordreDeJeu[0]
	}

	// Raise event
	c.RaiseEvent(NewCombatDemarreEvent(c.id, c.tourActuel, c.ordreDeJeu))

	return nil
}

// ExecuterAction exécute une action de combat
func (c *Combat) ExecuterAction(action *ActionCombat) error {
	if c.etat != EtatEnCours {
		return errors.New("le combat doit être en cours")
	}

	if c.phase != PhaseAttenteAction {
		return errors.New("pas en phase d'attente d'action")
	}

	// Valider que l'acteur est bien l'unité active
	if action.ActeurID != c.uniteActive {
		return errors.New("ce n'est pas le tour de cette unité")
	}

	// Récupérer l'acteur
	acteur := c.trouverUnite(action.ActeurID)
	if acteur == nil {
		return errors.New("acteur introuvable")
	}

	// Valider l'action
	if err := c.validerAction(acteur, action); err != nil {
		return err
	}

	// Résoudre l'action
	c.phase = PhaseExecutionAction
	c.actionEnCours = action

	resultat, err := c.resoudreAction(acteur, action)
	if err != nil {
		return err
	}

	// Appliquer les effets
	c.appliquerResultatAction(resultat)

	// Raise event
	c.RaiseEvent(NewActionExecuteeEvent(c.id, c.tourActuel, action, resultat))

	// Passer à la phase suivante
	c.phase = PhasePostTraitement

	return nil
}

// TourSuivant passe au tour suivant
func (c *Combat) TourSuivant() error {
	if c.etat != EtatEnCours {
		return errors.New("le combat doit être en cours")
	}

	// Vérifier condition de fin
	if c.estTermine() {
		return c.Terminer()
	}

	// Passer à l'unité suivante
	c.passerUniteeSuivante()

	// Si on revient au début de l'ordre, nouveau tour
	if c.ordreDeJeu[0] == c.uniteActive {
		c.tourActuel++
		c.RaiseEvent(NewTourDemarreEvent(c.id, c.tourActuel))
	}

	c.phase = PhaseAttenteAction
	c.actionEnCours = nil

	return nil
}

// Terminer termine le combat
func (c *Combat) Terminer() error {
	if c.etat == EtatTermine {
		return errors.New("le combat est déjà terminé")
	}

	c.etat = EtatTermine
	c.phase = PhaseTermine

	// Déterminer le vainqueur
	vainqueur := c.determinerVainqueur()

	c.RaiseEvent(NewCombatTermineEvent(c.id, c.tourActuel, vainqueur))

	return nil
}

// Event Sourcing methods

// RaiseEvent ajoute un événement à la liste des événements non committés
func (c *Combat) RaiseEvent(evt Evenement) {
	evt.SetAggregateID(c.id)
	evt.SetAggregateVersion(c.version + len(c.evenements) + 1)
	evt.SetTimestamp(time.Now())
	c.evenements = append(c.evenements, evt)
}

// GetUncommittedEvents retourne les événements non committés
func (c *Combat) GetUncommittedEvents() []Evenement {
	return c.evenements
}

// ClearUncommittedEvents vide la liste des événements non committés
func (c *Combat) ClearUncommittedEvents() {
	c.evenements = make([]Evenement, 0)
}

// ReconstruireDepuisEvenements reconstruit l'état du combat depuis les événements
func ReconstruireDepuisEvenements(events []Evenement) (*Combat, error) {
	if len(events) == 0 {
		return nil, errors.New("aucun événement fourni")
	}

	// Premier événement doit être CombatDemarre
	firstEvent, ok := events[0].(*CombatDemarreEvent)
	if !ok {
		return nil, errors.New("le premier événement doit être CombatDemarre")
	}

	combat := &Combat{
		id:         firstEvent.AggregateID(),
		equipes:    make(map[TeamID]*Equipe),
		evenements: make([]Evenement, 0),
	}

	// Appliquer tous les événements
	for _, evt := range events {
		if err := combat.Apply(evt); err != nil {
			return nil, err
		}
		combat.version = evt.AggregateVersion()
	}

	return combat, nil
}

// Apply applique un événement à l'agrégat
func (c *Combat) Apply(evt Evenement) error {
	switch e := evt.(type) {
	case *CombatDemarreEvent:
		return c.applyCombatDemarre(e)
	case *ActionExecuteeEvent:
		return c.applyActionExecutee(e)
	case *TourDemarreEvent:
		return c.applyTourDemarre(e)
	case *CombatTermineEvent:
		return c.applyCombatTermine(e)
	default:
		return errors.New("type d'événement inconnu")
	}
}

// Private methods

func (c *Combat) applyCombatDemarre(evt *CombatDemarreEvent) error {
	c.etat = EtatEnCours
	c.tourActuel = evt.Tour
	c.ordreDeJeu = evt.OrdreInitiative
	if len(c.ordreDeJeu) > 0 {
		c.uniteActive = c.ordreDeJeu[0]
	}
	c.phase = PhaseAttenteAction
	return nil
}

func (c *Combat) applyActionExecutee(evt *ActionExecuteeEvent) error {
	// Mettre à jour l'état selon le résultat de l'action
	// TODO: Implémenter la logique de mise à jour
	return nil
}

func (c *Combat) applyTourDemarre(evt *TourDemarreEvent) error {
	c.tourActuel = evt.Tour
	return nil
}

func (c *Combat) applyCombatTermine(evt *CombatTermineEvent) error {
	c.etat = EtatTermine
	c.phase = PhaseTermine
	return nil
}

func (c *Combat) calculerOrdreInitiative() {
	// Collecter toutes les unités
	unites := make([]*Unite, 0)
	for _, equipe := range c.equipes {
		unites = append(unites, equipe.Membres()...)
	}

	// Trier par vitesse (SPD) décroissante
	// TODO: Implémenter le tri avec randomisation pour égalités

	// Construire l'ordre
	c.ordreDeJeu = make([]UnitID, len(unites))
	for i, unite := range unites {
		c.ordreDeJeu[i] = unite.ID()
	}
}

func (c *Combat) validerAction(acteur *Unite, action *ActionCombat) error {
	// Vérifier que l'unité peut agir
	if !acteur.PeutAgir() {
		return errors.New("l'unité ne peut pas agir")
	}

	// Vérifier les coûts (MP, Stamina)
	// TODO: Implémenter validation des coûts

	// Vérifier la portée
	// TODO: Implémenter validation de portée

	// Vérifier les cibles valides
	// TODO: Implémenter validation des cibles

	return nil
}

func (c *Combat) resoudreAction(acteur *Unite, action *ActionCombat) (*ResultatAction, error) {
	// Créer le résultat
	resultat := &ResultatAction{
		ActeurID:   action.ActeurID,
		TypeAction: action.Type,
		Effets:     make([]EffetAction, 0),
	}

	// Résoudre selon le type d'action
	switch action.Type {
	case TypeActionAttaque:
		c.resoudreAttaque(acteur, action, resultat)
	case TypeActionCompetence:
		c.resoudreCompetence(acteur, action, resultat)
	case TypeActionDeplacement:
		c.resoudreDeplacement(acteur, action, resultat)
	case TypeActionObjet:
		c.resoudreObjet(acteur, action, resultat)
	default:
		return nil, errors.New("type d'action inconnu")
	}

	return resultat, nil
}

func (c *Combat) resoudreAttaque(acteur *Unite, action *ActionCombat, resultat *ResultatAction) {
	if action.CibleID == nil {
		return
	}

	// Trouver la cible
	cible := c.trouverUnite(*action.CibleID)
	if cible == nil {
		return
	}

	// Utiliser l'attaque basique (compétence par défaut)
	// Pour une attaque basique, on utilise le calculator par défaut (physique)
	attaqueBasique := acteur.ObtenirCompetenceParDefaut()

	// Calculer les dégâts avec le Strategy Pattern
	degats := c.damageCalculator.Calculate(acteur, cible, attaqueBasique)

	// Appliquer les dégâts à la cible
	cible.RecevoirDegats(degats)

	// Ajouter l'effet au résultat
	resultat.Effets = append(resultat.Effets, EffetAction{
		Type:    TypeEffetActionDegats,
		CibleID: *action.CibleID,
		Valeur:  degats,
	})

	// Raise event
	c.RaiseEvent(NewDegatsInfligesEvent(c.id, c.tourActuel, action.ActeurID, *action.CibleID, degats))

	// Vérifier si la cible est éliminée
	if cible.EstEliminee() {
		c.RaiseEvent(NewUniteElimineeEvent(c.id, c.tourActuel, *action.CibleID))
	}
}

func (c *Combat) resoudreCompetence(acteur *Unite, action *ActionCombat, resultat *ResultatAction) {
	if action.CompetenceID == nil {
		return
	}

	// Trouver la compétence
	competence := acteur.ObtenirCompetence(*action.CompetenceID)
	if competence == nil {
		return
	}

	// Vérifier le coût en MP/Stamina
	statsActeur := acteur.Stats()
	if statsActeur.MP < competence.CoutMP() || statsActeur.Stamina < competence.CoutStamina() {
		return // Pas assez de ressources
	}

	// Consommer les ressources
	acteur.ConsommerMP(competence.CoutMP())
	acteur.ConsommerStamina(competence.CoutStamina())

	// Créer le calculator approprié selon le type de compétence
	calculator := c.calculatorFactory.CreateCalculator(competence)

	// Résoudre selon le type d'effet
	for _, effet := range competence.Effets() {
		switch effet.TypeEffet() {
		case EffetDegats:
			c.resoudreEffetDegats(acteur, action, competence, calculator, resultat)
		case EffetSoin:
			c.resoudreEffetSoin(acteur, action, competence, resultat)
		case EffetStatut:
			c.resoudreEffetStatut(acteur, action, competence, resultat)
		case EffetInvocation:
			c.resoudreEffetInvocation(acteur, action, competence, resultat)
		}
	}

	// Raise event
	cibles := make([]UnitID, 0)
	if action.CibleID != nil {
		cibles = append(cibles, *action.CibleID)
	}
	c.RaiseEvent(NewCompetenceUtiliseeEvent(c.id, c.tourActuel, action.ActeurID, *action.CompetenceID, cibles))
}

func (c *Combat) resoudreEffetDegats(acteur *Unite, action *ActionCombat, competence *Competence, calculator DamageCalculator, resultat *ResultatAction) {
	if action.CibleID == nil {
		return
	}

	// Trouver la cible
	cible := c.trouverUnite(*action.CibleID)
	if cible == nil {
		return
	}

	// Calculer les dégâts avec le Strategy Pattern
	degats := calculator.Calculate(acteur, cible, competence)

	// Appliquer les dégâts
	cible.RecevoirDegats(degats)

	// Ajouter l'effet au résultat
	resultat.Effets = append(resultat.Effets, EffetAction{
		Type:    TypeEffetActionDegats,
		CibleID: *action.CibleID,
		Valeur:  degats,
	})

	// Raise event
	c.RaiseEvent(NewDegatsInfligesEvent(c.id, c.tourActuel, action.ActeurID, *action.CibleID, degats))

	// Vérifier si la cible est éliminée
	if cible.EstEliminee() {
		c.RaiseEvent(NewUniteElimineeEvent(c.id, c.tourActuel, *action.CibleID))
	}
}

func (c *Combat) resoudreEffetSoin(acteur *Unite, action *ActionCombat, competence *Competence, resultat *ResultatAction) {
	if action.CibleID == nil {
		return
	}

	// Trouver la cible
	cible := c.trouverUnite(*action.CibleID)
	if cible == nil {
		return
	}

	// Calculer le soin (base + scaling)
	statsActeur := acteur.Stats()
	soinBase := competence.DegatsBase() // Réutilise le champ pour le montant de soin
	scaling := competence.Modificateur() * float64(statsActeur.MATK)
	soin := soinBase + int(scaling)

	// Appliquer le soin
	cible.RecevoirSoin(soin)

	// Ajouter l'effet au résultat
	resultat.Effets = append(resultat.Effets, EffetAction{
		Type:    TypeEffetActionSoin,
		CibleID: *action.CibleID,
		Valeur:  soin,
	})

	// Raise event
	c.RaiseEvent(NewSoinApliqueEvent(c.id, c.tourActuel, action.ActeurID, *action.CibleID, soin))
}

func (c *Combat) resoudreEffetStatut(acteur *Unite, action *ActionCombat, competence *Competence, resultat *ResultatAction) {
	if action.CibleID == nil {
		return
	}

	// Trouver la cible
	cible := c.trouverUnite(*action.CibleID)
	if cible == nil {
		return
	}

	// Appliquer chaque statut de la compétence
	for _, effet := range competence.Effets() {
		if effet.TypeEffet() == EffetStatut && effet.StatutType() != nil {
			// Créer un statut à partir du type (duree, puissance)
			statut := shared.NewStatut(*effet.StatutType(), effet.Duree(), effet.Valeur())

			cible.AppliquerStatut(statut)

			// Ajouter l'effet au résultat
			resultat.Effets = append(resultat.Effets, EffetAction{
				Type:    TypeEffetActionStatut,
				CibleID: *action.CibleID,
				Valeur:  0, // Pas de valeur numérique pour un statut
			})

			// Raise event
			c.RaiseEvent(NewStatutAppliqueEvent(c.id, c.tourActuel, action.ActeurID, *action.CibleID, statut))
		}
	}
}

func (c *Combat) resoudreEffetInvocation(acteur *Unite, action *ActionCombat, competence *Competence, resultat *ResultatAction) {
	// TODO: Implémenter invocation (Step B avancé)
}

func (c *Combat) resoudreDeplacement(acteur *Unite, action *ActionCombat, resultat *ResultatAction) {
	// Validation: l'action doit avoir une position cible
	if action.PositionCible == nil {
		resultat.Succes = false
		resultat.MessageErreur = "position cible requise pour le déplacement"
		return
	}

	positionActuelle := acteur.Position()
	positionCible := action.PositionCible

	// Vérifier que la destination est différente
	if positionActuelle.Equals(positionCible) {
		resultat.Succes = false
		resultat.MessageErreur = "l'unité est déjà à cette position"
		return
	}

	// Vérifier que l'unité peut se déplacer (pas Root ou Stun)
	if acteur.EstBloqueDeplacement() {
		resultat.Succes = false
		resultat.MessageErreur = "l'unité ne peut pas se déplacer (statut bloquant)"
		return
	}

	// Créer une map des positions occupées par les autres unités
	unitesOccupees := c.obtenirPositionsOccupees(acteur.ID())

	// Utiliser le service de pathfinding (Strategy Pattern)
	pathfindingService := NewPathfindingService()

	// Choisir la stratégie selon le contexte (Factory Pattern)
	// Par défaut : Manhattan (déplacement tactique en grille)
	pathfindingService.SetStrategyType("manhattan")

	// Trouver le chemin avec respect de la portée de mouvement
	porteeMax := acteur.Stats().MOV
	chemin, cout, err := pathfindingService.TrouverCheminAvecPortee(
		c.grille,
		positionActuelle,
		positionCible,
		unitesOccupees,
		porteeMax,
	)

	if err != nil {
		resultat.Succes = false
		resultat.MessageErreur = err.Error()
		return
	}

	// Déplacer l'unité
	acteur.DeplacerVers(positionCible)

	// Ajouter les informations du déplacement au résultat
	resultat.Succes = true
	resultat.CoutDeplacement = cout
	resultat.CheminParcouru = chemin

	// Raise event
	c.RaiseEvent(NewDeplacementExecuteEvent(
		c.id,
		c.tourActuel,
		acteur.ID(),
		positionActuelle,
		positionCible,
		chemin,
		cout,
	))
}

func (c *Combat) resoudreObjet(acteur *Unite, action *ActionCombat, resultat *ResultatAction) {
	// TODO: Implémenter résolution objet
}

func (c *Combat) appliquerResultatAction(resultat *ResultatAction) {
	// Appliquer tous les effets
	for _, effet := range resultat.Effets {
		switch effet.Type {
		case TypeEffetActionDegats:
			c.appliquerDegats(effet.CibleID, effet.Valeur)
		case TypeEffetActionSoin:
			c.appliquerSoin(effet.CibleID, effet.Valeur)
		case TypeEffetActionStatut:
			c.appliquerStatut(effet.CibleID, effet.Statut)
		}
	}
}

func (c *Combat) appliquerDegats(cibleID UnitID, degats int) {
	cible := c.trouverUnite(cibleID)
	if cible == nil {
		return
	}

	cible.RecevoirDegats(degats)

	// Raise event
	c.RaiseEvent(NewDegatsInfligesEvent(c.id, c.tourActuel, c.uniteActive, cibleID, degats))

	// Vérifier si l'unité est éliminée
	if cible.EstEliminee() {
		c.RaiseEvent(NewUniteElimineeEvent(c.id, c.tourActuel, cibleID))
	}
}

func (c *Combat) appliquerSoin(cibleID UnitID, soin int) {
	cible := c.trouverUnite(cibleID)
	if cible == nil {
		return
	}

	cible.RecevoirSoin(soin)

	// Raise event
	c.RaiseEvent(NewSoinApliqueEvent(c.id, c.tourActuel, c.uniteActive, cibleID, soin))
}

func (c *Combat) appliquerStatut(cibleID UnitID, statut *shared.Statut) {
	cible := c.trouverUnite(cibleID)
	if cible == nil {
		return
	}

	cible.AjouterStatut(statut)

	// Raise event
	c.RaiseEvent(NewStatutAppliqueEvent(c.id, c.tourActuel, c.uniteActive, cibleID, statut))
}

func (c *Combat) trouverUnite(id UnitID) *Unite {
	for _, equipe := range c.equipes {
		for _, unite := range equipe.Membres() {
			if unite.ID() == id {
				return unite
			}
		}
	}
	return nil
}

func (c *Combat) passerUniteeSuivante() {
	// Trouver l'index de l'unité active
	currentIndex := -1
	for i, id := range c.ordreDeJeu {
		if id == c.uniteActive {
			currentIndex = i
			break
		}
	}

	// Passer à la suivante (avec wrap around)
	nextIndex := (currentIndex + 1) % len(c.ordreDeJeu)
	c.uniteActive = c.ordreDeJeu[nextIndex]
}

func (c *Combat) obtenirPositionsOccupees(exclusionID UnitID) map[string]bool {
	// Crée une map des positions occupées par toutes les unités sauf celle exclue
	positions := make(map[string]bool)

	for _, equipe := range c.equipes {
		for _, unite := range equipe.Membres() {
			if unite.ID() != exclusionID && !unite.EstEliminee() {
				pos := unite.Position()
				cle := positionKey(pos)
				positions[cle] = true
			}
		}
	}

	return positions
}

func (c *Combat) estTermine() bool {
	// Compter les équipes actives (avec au moins 1 membre vivant)
	equipesActives := 0
	for _, equipe := range c.equipes {
		if equipe.ADesMembresVivants() {
			equipesActives++
		}
	}

	// Combat terminé si 0 ou 1 équipe active
	return equipesActives <= 1
}

func (c *Combat) determinerVainqueur() *TeamID {
	for _, equipe := range c.equipes {
		if equipe.ADesMembresVivants() {
			id := equipe.ID()
			return &id
		}
	}
	return nil
}

// Step C - Setters et getters pour les composants (utilisés par combatinitializer et combatfacade)

func (c *Combat) SetStateMachine(sm interface{}) {
	c.stateMachine = sm
}

func (c *Combat) SetCommandInvoker(invoker interface{}) {
	c.commandInvoker = invoker
}

func (c *Combat) SetCommandFactory(factory interface{}) {
	c.commandFactory = factory
}

func (c *Combat) SetObserverSubject(subject interface{}) {
	c.observerSubject = subject
}

func (c *Combat) SetValidationChain(chain interface{}) {
	c.validationChain = chain
}

func (c *Combat) GetStateMachineRaw() interface{} {
	return c.stateMachine
}

func (c *Combat) GetCommandInvokerRaw() interface{} {
	return c.commandInvoker
}

func (c *Combat) GetCommandFactoryRaw() interface{} {
	return c.commandFactory
}

func (c *Combat) GetObserverSubjectRaw() interface{} {
	return c.observerSubject
}

func (c *Combat) GetValidationChainRaw() interface{} {
	return c.validationChain
}
