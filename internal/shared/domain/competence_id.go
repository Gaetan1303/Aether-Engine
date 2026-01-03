package domain

// CompetenceID identifiant unique commun aux compétences joueur et combat
// Permet le lien entre:
// - internal/joueur/domain/competence.go (compétences possibles du joueur)
// - internal/combat/domain/competence.go (compétences utilisables par tous)
type CompetenceID string

// Compétences de base disponibles
const (
	// Compétences Guerrier
	CompetenceAttaqueBasique   CompetenceID = "attaque_basique"
	CompetenceCharge          CompetenceID = "charge"
	CompetenceCoupPuissant    CompetenceID = "coup_puissant"
	CompetenceDefenseRenforcee CompetenceID = "defense_renforcee"

	// Compétences Mage
	CompetenceBouleDefeu     CompetenceID = "boule_de_feu"
	CompetenceEclairGlace    CompetenceID = "eclair_glace"
	CompetenceSoinMineur     CompetenceID = "soin_mineur"
	CompetenceBouclierMagique CompetenceID = "bouclier_magique"

	// Compétences Archer
	CompetenceTirPrecis      CompetenceID = "tir_precis"
	CompetenceTirPerforant   CompetenceID = "tir_perforant"
	CompetencePluieFleches   CompetenceID = "pluie_fleches"
	CompetenceRecul          CompetenceID = "recul"

	// Compétences Voleur
	CompetenceAttaqueSournoise CompetenceID = "attaque_sournoise"
	CompetenceEvasion         CompetenceID = "evasion"
	CompetenceVol             CompetenceID = "vol"
	CompetenceInvisibilite    CompetenceID = "invisibilite"

	// Compétences Clerc
	CompetenceSoinMajeur     CompetenceID = "soin_majeur"
	CompetenceBenediction    CompetenceID = "benediction"
	CompetenceGuerisonStatut CompetenceID = "guerison_statut"
	CompetenceResurrection   CompetenceID = "resurrection"
)