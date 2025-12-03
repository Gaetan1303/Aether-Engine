package domain

import (
	"errors"
	"time"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Combat est l'agrégat racine qui gère une instance de combat tactique
// Architecture: State Machine + Command Pattern + Observer Pattern + Chain of Responsibility
type Combat struct {
	id                string
	etat              EtatCombat
	equipes           map[TeamID]*Equipe
	grille            *shared.GrilleCombat
	tourActuel        int
	version           int
	evenements        []Evenement // Uncommitted events
	createdAt         time.Time
	updatedAt         time.Time
	damageCalculator  DamageCalculator         // Strategy Pattern - algorithme de dégâts
	calculatorFactory *DamageCalculatorFactory // Factory pour créer strategies

	// Step C - Design Patterns Architecture
	// Interfaces typées pour éviter les cycles d'imports (Interface Segregation Principle)
	stateMachine    StateMachineProvider
	commandInvoker  CommandSystemProvider
	commandFactory  CommandFactoryProvider
	observerSubject ObserverProvider
	validationChain ValidationProvider
	fuiteAutorisee  bool            // Indique si la fuite est autorisée
	equipesFuites   map[TeamID]bool // Équipes ayant fui le combat
}

// NewCombat crée une nouvelle instance de combat
func NewCombat(id string, equipes []*Equipe, grille *shared.GrilleCombat) (*Combat, error) {
	if len(equipes) < MinEquipesPourCombat {
		return nil, errors.New("minimum 2 équipes requises")
	}

	combat := &Combat{
		id:                id,
		etat:              EtatAttente,
		equipes:           make(map[TeamID]*Equipe),
		grille:            grille,
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
func (c *Combat) Equipes() map[TeamID]*Equipe  { return c.equipes }
func (c *Combat) GetTimestamp() int64          { return c.updatedAt.Unix() }

// Step C - Méthodes publiques pour la State Machine
// Toute l'exécution passe par combatfacade.ExecutePlayerAction()

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
		if equipesActives == EquipeActiveVictoire {
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

// Demarrer démarre le combat via la State Machine
// Utilise CombatInitializer.InitializeAll() puis combatfacade pour démarrer
func (c *Combat) Demarrer() error {
	if c.etat != EtatAttente {
		return errors.New("le combat doit être en attente pour démarrer")
	}

	c.etat = EtatEnCours
	c.tourActuel = 1

	// La State Machine gère maintenant le démarrage
	// via la transition Idle → Initializing → Ready
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
// Utilisé pour Event Sourcing - rejoue tous les événements pour reconstruire l'état
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

// Apply applique un événement à l'agrégat (Event Sourcing)
func (c *Combat) Apply(evt Evenement) error {
	switch e := evt.(type) {
	case *CombatDemarreEvent:
		c.etat = EtatEnCours
		c.tourActuel = e.Tour
		return nil
	case *TourDemarreEvent:
		c.tourActuel = e.Tour
		return nil
	case *CombatTermineEvent:
		c.etat = EtatTermine
		return nil
	case *ActionExecuteeEvent, *DegatsInfligesEvent, *SoinApliqueEvent,
		*StatutAppliqueEvent, *UniteElimineeEvent, *CompetenceUtiliseeEvent,
		*DeplacementExecuteEvent:
		// Événements gérés par la State Machine
		return nil
	default:
		return errors.New("type d'événement inconnu")
	}
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

// Step C - Setters et getters pour les composants Step C (utilisés par combatinitializer et combatfacade)

// Step C - Setters avec interfaces typées (Dependency Inversion Principle)

func (c *Combat) SetStateMachine(sm StateMachineProvider) {
	c.stateMachine = sm
}

func (c *Combat) SetCommandInvoker(invoker CommandSystemProvider) {
	c.commandInvoker = invoker
}

func (c *Combat) SetCommandFactory(factory CommandFactoryProvider) {
	c.commandFactory = factory
}

func (c *Combat) SetObserverSubject(subject ObserverProvider) {
	c.observerSubject = subject
}

func (c *Combat) SetValidationChain(chain ValidationProvider) {
	c.validationChain = chain
}

// Step C - Getters avec interfaces typées

func (c *Combat) GetStateMachine() StateMachineProvider {
	return c.stateMachine
}

func (c *Combat) GetCommandInvoker() CommandSystemProvider {
	return c.commandInvoker
}

func (c *Combat) GetCommandFactory() CommandFactoryProvider {
	return c.commandFactory
}

func (c *Combat) GetObserverSubject() ObserverProvider {
	return c.observerSubject
}

func (c *Combat) GetValidationChain() ValidationProvider {
	return c.validationChain
}
