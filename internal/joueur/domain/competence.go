package domain

// CompetenceID identifiant d'une compétence
type CompetenceID string

// Types de compétences
type TypeCompetence string

const (
	TypeActif    TypeCompetence = "actif"    // Compétence activable en combat
	TypePassif   TypeCompetence = "passif"   // Bonus permanent
	TypeReaction TypeCompetence = "reaction" // Se déclenche automatiquement
	TypeSupport  TypeCompetence = "support"  // Compétence de support hors combat
)

// Competence définit une compétence/sort
type Competence struct {
	ID           CompetenceID   `json:"id"`
	Nom          string         `json:"nom"`
	Description  string         `json:"description"`
	Type         TypeCompetence `json:"type"`
	CoutMP       int            `json:"cout_mp"`
	CoutAP       int            `json:"cout_ap"` // Action Points pour le combat tour par tour
	Portee       int            `json:"portee"`
	Zone         string         `json:"zone"` // "single", "line", "aoe", etc.
	Cooldown     int            `json:"cooldown"`
	JobOrigine   JobID          `json:"job_origine"`
	NiveauRequis int            `json:"niveau_requis"`
	Effets       []Effet        `json:"effets"`
}

// Effet d'une compétence
type Effet struct {
	Type     string  `json:"type"` // "damage", "heal", "buff", "debuff", etc.
	Valeur   float64 `json:"valeur"`
	Duree    int     `json:"duree"`    // En tours
	Attribut string  `json:"attribut"` // Stat affectée
}

// CompetenceSlots système de slots de compétences comme FF Tactics A2
type CompetenceSlots struct {
	// Slots actifs (compétences utilisables en combat)
	SlotsActifs    []SlotCompetence `json:"slots_actifs"`
	MaxSlotsActifs int              `json:"max_slots_actifs"`

	// Slots passifs (bonus permanents)
	SlotsPassifs    []SlotCompetence `json:"slots_passifs"`
	MaxSlotsPassifs int              `json:"max_slots_passifs"`

	// Slots réaction (compétences automatiques)
	SlotsReaction    []SlotCompetence `json:"slots_reaction"`
	MaxSlotsReaction int              `json:"max_slots_reaction"`

	// Slots support (hors combat)
	SlotsSupport    []SlotCompetence `json:"slots_support"`
	MaxSlotsSupport int              `json:"max_slots_support"`

	// Compétences disponibles (débloquées via les jobs)
	CompetencesDisponibles []CompetenceID `json:"competences_disponibles"`
}

// SlotCompetence un slot de compétence avec sa compétence assignée
type SlotCompetence struct {
	Actif        bool         `json:"actif"`
	CompetenceID CompetenceID `json:"competence_id,omitempty"`
}

// NewCompetenceSlots crée un nouveau système de slots
func NewCompetenceSlots() CompetenceSlots {
	return CompetenceSlots{
		SlotsActifs:            make([]SlotCompetence, 4), // 4 slots actifs au début
		MaxSlotsActifs:         4,
		SlotsPassifs:           make([]SlotCompetence, 2), // 2 slots passifs
		MaxSlotsPassifs:        2,
		SlotsReaction:          make([]SlotCompetence, 1), // 1 slot réaction
		MaxSlotsReaction:       1,
		SlotsSupport:           make([]SlotCompetence, 2), // 2 slots support
		MaxSlotsSupport:        2,
		CompetencesDisponibles: make([]CompetenceID, 0),
	}
}

// AssignerCompetence assigne une compétence à un slot
func (cs *CompetenceSlots) AssignerCompetence(typeComp TypeCompetence, slotIndex int, competenceID CompetenceID) error {
	// Vérifier si la compétence est disponible
	if !cs.CompetenceDisponible(competenceID) {
		return NewErrCompetenceNonDisponible(competenceID)
	}

	switch typeComp {
	case TypeActif:
		if slotIndex >= len(cs.SlotsActifs) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsActifs[slotIndex] = SlotCompetence{
			Actif:        true,
			CompetenceID: competenceID,
		}

	case TypePassif:
		if slotIndex >= len(cs.SlotsPassifs) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsPassifs[slotIndex] = SlotCompetence{
			Actif:        true,
			CompetenceID: competenceID,
		}

	case TypeReaction:
		if slotIndex >= len(cs.SlotsReaction) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsReaction[slotIndex] = SlotCompetence{
			Actif:        true,
			CompetenceID: competenceID,
		}

	case TypeSupport:
		if slotIndex >= len(cs.SlotsSupport) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsSupport[slotIndex] = SlotCompetence{
			Actif:        true,
			CompetenceID: competenceID,
		}

	default:
		return NewErrTypeCompetenceInvalide(typeComp)
	}

	return nil
}

// RetirerCompetence retire une compétence d'un slot
func (cs *CompetenceSlots) RetirerCompetence(typeComp TypeCompetence, slotIndex int) error {
	switch typeComp {
	case TypeActif:
		if slotIndex >= len(cs.SlotsActifs) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsActifs[slotIndex] = SlotCompetence{Actif: false}

	case TypePassif:
		if slotIndex >= len(cs.SlotsPassifs) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsPassifs[slotIndex] = SlotCompetence{Actif: false}

	case TypeReaction:
		if slotIndex >= len(cs.SlotsReaction) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsReaction[slotIndex] = SlotCompetence{Actif: false}

	case TypeSupport:
		if slotIndex >= len(cs.SlotsSupport) || slotIndex < 0 {
			return NewErrSlotInvalide(typeComp, slotIndex)
		}
		cs.SlotsSupport[slotIndex] = SlotCompetence{Actif: false}

	default:
		return NewErrTypeCompetenceInvalide(typeComp)
	}

	return nil
}

// IsActive vérifie si une compétence est active dans les slots
func (cs *CompetenceSlots) IsActive(competenceID CompetenceID) bool {
	// Vérifier dans tous les types de slots
	slots := [][]SlotCompetence{
		cs.SlotsActifs,
		cs.SlotsPassifs,
		cs.SlotsReaction,
		cs.SlotsSupport,
	}

	for _, slotType := range slots {
		for _, slot := range slotType {
			if slot.Actif && slot.CompetenceID == competenceID {
				return true
			}
		}
	}

	return false
}

// CompetenceDisponible vérifie si une compétence est disponible
func (cs *CompetenceSlots) CompetenceDisponible(competenceID CompetenceID) bool {
	for _, compDisponible := range cs.CompetencesDisponibles {
		if compDisponible == competenceID {
			return true
		}
	}
	return false
}

// AjouterCompetenceDisponible ajoute une compétence aux disponibles
func (cs *CompetenceSlots) AjouterCompetenceDisponible(competenceID CompetenceID) {
	if !cs.CompetenceDisponible(competenceID) {
		cs.CompetencesDisponibles = append(cs.CompetencesDisponibles, competenceID)
	}
}

// GetCompetencesActives retourne toutes les compétences actives
func (cs *CompetenceSlots) GetCompetencesActives() []CompetenceID {
	competences := make([]CompetenceID, 0)

	// Ajouter les compétences actives de tous les types
	slots := [][]SlotCompetence{
		cs.SlotsActifs,
		cs.SlotsPassifs,
		cs.SlotsReaction,
		cs.SlotsSupport,
	}

	for _, slotType := range slots {
		for _, slot := range slotType {
			if slot.Actif && slot.CompetenceID != "" {
				competences = append(competences, slot.CompetenceID)
			}
		}
	}

	return competences
}

// AugmenterSlots augmente le nombre de slots disponibles (progression)
func (cs *CompetenceSlots) AugmenterSlots(typeComp TypeCompetence, nombre int) {
	switch typeComp {
	case TypeActif:
		cs.MaxSlotsActifs += nombre
		// Ajouter de nouveaux slots vides
		for i := 0; i < nombre; i++ {
			cs.SlotsActifs = append(cs.SlotsActifs, SlotCompetence{Actif: false})
		}

	case TypePassif:
		cs.MaxSlotsPassifs += nombre
		for i := 0; i < nombre; i++ {
			cs.SlotsPassifs = append(cs.SlotsPassifs, SlotCompetence{Actif: false})
		}

	case TypeReaction:
		cs.MaxSlotsReaction += nombre
		for i := 0; i < nombre; i++ {
			cs.SlotsReaction = append(cs.SlotsReaction, SlotCompetence{Actif: false})
		}

	case TypeSupport:
		cs.MaxSlotsSupport += nombre
		for i := 0; i < nombre; i++ {
			cs.SlotsSupport = append(cs.SlotsSupport, SlotCompetence{Actif: false})
		}
	}
}

// Types d'erreur pour le système de compétences
type ErrCompetenceNonDisponible struct {
	CompetenceID CompetenceID
}

func NewErrCompetenceNonDisponible(competenceID CompetenceID) error {
	return &ErrCompetenceNonDisponible{CompetenceID: competenceID}
}

func (e *ErrCompetenceNonDisponible) Error() string {
	return "compétence non disponible: " + string(e.CompetenceID)
}

type ErrSlotInvalide struct {
	Type  TypeCompetence
	Index int
}

func NewErrSlotInvalide(typeComp TypeCompetence, index int) error {
	return &ErrSlotInvalide{Type: typeComp, Index: index}
}

func (e *ErrSlotInvalide) Error() string {
	return "slot invalide pour type " + string(e.Type) + " à l'index " + string(rune(e.Index))
}

type ErrTypeCompetenceInvalide struct {
	Type TypeCompetence
}

func NewErrTypeCompetenceInvalide(typeComp TypeCompetence) error {
	return &ErrTypeCompetenceInvalide{Type: typeComp}
}

func (e *ErrTypeCompetenceInvalide) Error() string {
	return "type de compétence invalide: " + string(e.Type)
}
