package domain

// EtatCombat représente l'état global du combat
type EtatCombat int

const (
	EtatAttente EtatCombat = iota // En attente de démarrage
	EtatEnCours                   // Combat en cours
	EtatPause                     // Combat en pause
	EtatTermine                   // Combat terminé
	EtatAnnule                    // Combat annulé
)

func (e EtatCombat) String() string {
	switch e {
	case EtatAttente:
		return "Attente"
	case EtatEnCours:
		return "EnCours"
	case EtatPause:
		return "Pause"
	case EtatTermine:
		return "Terminé"
	case EtatAnnule:
		return "Annulé"
	default:
		return "Inconnu"
	}
}

// PhaseTour représente la phase actuelle du tour
type PhaseTour int

const (
	PhasePreparation     PhaseTour = iota // Préparation initiale
	PhaseDebutTour                        // Début du tour (traitement statuts, etc.)
	PhaseAttenteAction                    // Attente de l'action du joueur
	PhaseExecutionAction                  // Exécution d'une action
	PhasePostTraitement                   // Post-traitement (fin d'action)
	PhaseFinTour                          // Fin du tour
	PhaseTermine                          // Combat terminé
)

func (p PhaseTour) String() string {
	switch p {
	case PhasePreparation:
		return "Préparation"
	case PhaseDebutTour:
		return "DébutTour"
	case PhaseAttenteAction:
		return "AttenteAction"
	case PhaseExecutionAction:
		return "ExecutionAction"
	case PhasePostTraitement:
		return "PostTraitement"
	case PhaseFinTour:
		return "FinTour"
	case PhaseTermine:
		return "Terminé"
	default:
		return "Inconnu"
	}
}

// TypeStatut énumère les types de statuts (défini dans value_objects.go mais réexporté ici)
// Voir value_objects.go pour la définition complète

// TypeAction énumère les types d'actions (défini dans events.go mais réexporté ici)
// Voir events.go pour la définition complète

// TypeCompetence énumère les types de compétences (défini dans competence.go mais réexporté ici)
// Voir competence.go pour la définition complète

// TypeCellule représente le type de terrain d'une cellule (défini dans value_objects.go mais réexporté ici)
// Voir value_objects.go pour la définition complète
