package domain

// ModificateurStat représente un modificateur temporaire de stat
type ModificateurStat struct {
	Stat   string // "ATK", "DEF", etc.
	Valeur int    // Peut être négatif
}

// EffetStatut représente l'effet d'un statut appliqué
type EffetStatut struct {
	Type   TypeStatut
	Valeur int // Dégâts/soin périodique
}

// ObjetID est l'identifiant d'un objet
type ObjetID string

// Item représente un objet utilisable (Potion, Éther, etc.)
// TODO: Implémenter complètement le système d'items
type Item struct {
	ID          string
	Name        string
	Description string
	ItemType    string // "Potion", "Ether", "Antidote", "Revive", "Bomb"
	EffectVal   int    // Renommé pour éviter conflit avec méthode
	Range       int
}

// Constantes pour les types d'items
const (
	ItemTypePotion   = "Potion"
	ItemTypeEther    = "Ether"
	ItemTypeAntidote = "Antidote"
	ItemTypeRevive   = "Revive"
	ItemTypeBomb     = "Bomb"
)

// Méthodes temporaires pour Item
func (i *Item) GetID() string       { return i.ID }
func (i *Item) GetName() string     { return i.Name }
func (i *Item) GetItemType() string { return i.ItemType }
func (i *Item) EffectValue() int    { return i.EffectVal }
func (i *Item) GetRange() int       { return i.Range }

// Competence représente une compétence (temporaire, utiliser domain.Competence à la place)
// TODO: Supprimer ce type et utiliser domain.Competence partout
type Competence struct {
	ID       string
	Name     string
	MPCost   int
	Cooldown int
}

// DomainError représente une erreur métier du domaine
type DomainError struct {
	Message string
	Code    string
}

// NewDomainError crée une nouvelle erreur métier
func NewDomainError(message, code string) error {
	return &DomainError{
		Message: message,
		Code:    code,
	}
}

// Error implémente l'interface error
func (e *DomainError) Error() string {
	return e.Message
}
