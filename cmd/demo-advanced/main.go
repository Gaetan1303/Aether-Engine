package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Couleurs ANSI
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorOrange = "\033[38;5;208m"
	ColorPink   = "\033[38;5;205m"
)

type CombatStats struct {
	AttaquesTotal    int
	AttaquesReussies int
	AttaquesRatees   int
	DegatsInfliges   int
	DegatsRecus      int
	CompetencesUsees int
	ToursJoues       int
}

type GameDemo struct {
	combat        *domain.Combat
	equipeHeros   *domain.Equipe
	equipeEnnemis *domain.Equipe
	reader        *bufio.Reader
	stats         map[string]*CombatStats
	tourActuel    int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	afficherBanniere()

	game := &GameDemo{
		reader: bufio.NewReader(os.Stdin),
		stats:  make(map[string]*CombatStats),
	}

	if err := game.initialiserCombat(); err != nil {
		fmt.Printf(ColorRed+"Erreur initialisation: %v\n"+ColorReset, err)
		return
	}

	game.afficherIntro()
	game.boucleDeJeu()
}

func afficherBanniere() {
	fmt.Println(ColorBold + ColorCyan + "╔════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║                                                ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║     🏰  AETHER ENGINE - DEMO AVANCÉE  ⚔️      ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║                                                ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║          Combat Tactique 3v3 Épique            ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║                                                ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "╚════════════════════════════════════════════════╝" + ColorReset)
	fmt.Println()
	fmt.Println(ColorYellow + "✨ Nouvelles fonctionnalités :" + ColorReset)
	fmt.Println("   • Système ATH (chances de toucher)")
	fmt.Println("   • 6 unités uniques (3 vs 3)")
	fmt.Println("   • 5+ compétences variées")
	fmt.Println("   • Statistiques de combat en temps réel")
	fmt.Println("   • IA améliorée avec priorisation")
	fmt.Println()
}

func (g *GameDemo) initialiserCombat() error {
	grille, _ := shared.NewGrilleCombat(10, 10)

	// === ÉQUIPE HÉROS ===
	joueurID := "player-1"
	equipeHeros, _ := domain.NewEquipe("team-heros", "Héros de Lumière", "#00FF00", false, &joueurID)

	// 1. Paladin Tank (haute DEF, faible ATH pour équilibrer)
	statsPaladin := &shared.Stats{
		HP:      150,
		MP:      50,
		Stamina: 100,
		ATK:     20,
		DEF:     25,
		MATK:    8,
		MDEF:    18,
		SPD:     8,
		MOV:     3,
		ATH:     80, // Tank moins précis
	}
	posPaladin, _ := shared.NewPosition(1, 4)
	paladin := domain.NewUnite("hero-paladin", "Paladin", "team-heros", statsPaladin, posPaladin)

	// Compétence 1: Provocation (attire ennemis)
	taunt := domain.NewCompetence(
		"taunt",
		"Provocation",
		"Force un ennemi à vous attaquer",
		domain.CompetenceAttaque,
		3,
		domain.ZoneEffet{},
		15,
		0,
		3,
		10,
		1.2,
		domain.CibleEnnemis,
	)
	paladin.AjouterCompetence(taunt)

	// Compétence 2: Soin Divin
	heal := domain.NewCompetence(
		"heal",
		"Soin Divin",
		"Restaure les HP d'un allié",
		domain.CompetenceSoin,
		4,
		domain.ZoneEffet{},
		20,
		0,
		3,
		50,
		1.0,
		domain.CibleAllies,
	)
	paladin.AjouterCompetence(heal)

	// 2. Archer Sniper (haute ATH, haute portée)
	statsArcher := &shared.Stats{
		HP:      90,
		MP:      40,
		Stamina: 80,
		ATK:     28,
		DEF:     10,
		MATK:    5,
		MDEF:    8,
		SPD:     16,
		MOV:     4,
		ATH:     95, // Archer très précis
	}
	posArcher, _ := shared.NewPosition(2, 2)
	archer := domain.NewUnite("hero-archer", "Archer", "team-heros", statsArcher, posArcher)

	// Compétence: Tir de Précision
	precisionShot := domain.NewCompetence(
		"precision-shot",
		"Tir de Précision",
		"Tir précis avec dégâts augmentés",
		domain.CompetenceAttaque,
		6,
		domain.ZoneEffet{},
		0,
		15,
		2,
		35,
		1.8,
		domain.CibleEnnemis,
	)
	archer.AjouterCompetence(precisionShot)

	// 3. Mage Élémentaliste (haute MATK, compétences variées)
	statsMage := &shared.Stats{
		HP:      70,
		MP:      120,
		Stamina: 50,
		ATK:     8,
		DEF:     6,
		MATK:    35,
		MDEF:    22,
		SPD:     12,
		MOV:     3,
		ATH:     92, // Magie précise
	}
	posMage, _ := shared.NewPosition(2, 6)
	mage := domain.NewUnite("hero-mage", "Mage", "team-heros", statsMage, posMage)

	// Compétence 1: Boule de Feu
	fireball := domain.NewCompetence(
		"fireball",
		"Boule de Feu",
		"Projectile enflammé dévastateur",
		domain.CompetenceMagie,
		5,
		domain.ZoneEffet{},
		25,
		0,
		2,
		40,
		1.6,
		domain.CibleEnnemis,
	)
	mage.AjouterCompetence(fireball)

	// Compétence 2: Éclair
	lightning := domain.NewCompetence(
		"lightning",
		"Éclair",
		"Frappe éclair rapide",
		domain.CompetenceMagie,
		7,
		domain.ZoneEffet{},
		20,
		0,
		1,
		28,
		1.4,
		domain.CibleEnnemis,
	)
	mage.AjouterCompetence(lightning)

	// Compétence 3: Sommeil
	sleep := domain.NewCompetence(
		"sleep",
		"Sommeil",
		"Endort un ennemi (immobilisé 2 tours)",
		domain.CompetenceDebuff,
		5,
		domain.ZoneEffet{},
		15,
		0,
		4,
		0,
		1.0,
		domain.CibleEnnemis,
	)
	mage.AjouterCompetence(sleep)

	// Compétence 4: Boost Magique
	boost := domain.NewCompetence(
		"boost",
		"Boost Magique",
		"Augmente MATK d'un allié (+15 MATK, 3 tours)",
		domain.CompetenceBuff,
		4,
		domain.ZoneEffet{},
		18,
		0,
		3,
		0,
		1.0,
		domain.CibleAllies,
	)
	mage.AjouterCompetence(boost)

	equipeHeros.AjouterMembre(paladin)
	equipeHeros.AjouterMembre(archer)
	equipeHeros.AjouterMembre(mage)

	// === ÉQUIPE ENNEMIS ===
	equipeEnnemis, _ := domain.NewEquipe("team-ennemis", "Horde Gobeline", "#FF0000", true, nil)

	// 1. Chef Gobelin (équilibré, dangereux)
	statsChef := &shared.Stats{
		HP:      100,
		MP:      30,
		Stamina: 90,
		ATK:     24,
		DEF:     15,
		MATK:    10,
		MDEF:    12,
		SPD:     14,
		MOV:     4,
		ATH:     82, // Chef plus précis
	}
	posChef, _ := shared.NewPosition(8, 4)
	chef := domain.NewUnite("ennemi-chef", "Chef Gobelin", "team-ennemis", statsChef, posChef)

	// Compétence: Cri de Guerre (boost moral)
	warCry := domain.NewCompetence(
		"war-cry",
		"Cri de Guerre",
		"Boost l'attaque des alliés proches",
		domain.CompetenceBuff,
		0,
		domain.ZoneEffet{},
		10,
		20,
		4,
		0,
		1.0,
		domain.CibleAllies,
	)
	chef.AjouterCompetence(warCry)

	// 2. Gobelin Berserker (haute ATK, faible DEF)
	statsBerserker := &shared.Stats{
		HP:      85,
		MP:      10,
		Stamina: 100,
		ATK:     30,
		DEF:     8,
		MATK:    3,
		MDEF:    5,
		SPD:     15,
		MOV:     5,
		ATH:     72, // Berserker imprécis
	}
	posBerserker, _ := shared.NewPosition(8, 2)
	berserker := domain.NewUnite("ennemi-berserker", "Gobelin Berserker", "team-ennemis", statsBerserker, posBerserker)

	// 3. Shaman Gobelin (magie noire)
	statsShaman := &shared.Stats{
		HP:      65,
		MP:      80,
		Stamina: 60,
		ATK:     10,
		DEF:     8,
		MATK:    28,
		MDEF:    18,
		SPD:     11,
		MOV:     3,
		ATH:     88, // Shaman précis
	}
	posShaman, _ := shared.NewPosition(8, 6)
	shaman := domain.NewUnite("ennemi-shaman", "Shaman Gobelin", "team-ennemis", statsShaman, posShaman)

	// Compétence 1: Ombre Malefique
	darkBolt := domain.NewCompetence(
		"dark-bolt",
		"Éclair Sombre",
		"Magie noire corrompue",
		domain.CompetenceMagie,
		5,
		domain.ZoneEffet{},
		18,
		0,
		2,
		32,
		1.5,
		domain.CibleEnnemis,
	)
	shaman.AjouterCompetence(darkBolt)

	// Compétence 2: Paralysie
	paralyze := domain.NewCompetence(
		"paralyze",
		"Paralysie",
		"Paralyse un ennemi (immobilisé 2 tours)",
		domain.CompetenceDebuff,
		4,
		domain.ZoneEffet{},
		15,
		0,
		3,
		0,
		1.0,
		domain.CibleEnnemis,
	)
	shaman.AjouterCompetence(paralyze)

	// Compétence 3: Poison
	poison := domain.NewCompetence(
		"poison",
		"Poison",
		"Empoisonne un ennemi (10 dégâts/tour, 3 tours)",
		domain.CompetenceDebuff,
		5,
		domain.ZoneEffet{},
		12,
		0,
		4,
		0,
		1.0,
		domain.CibleEnnemis,
	)
	shaman.AjouterCompetence(poison)

	equipeEnnemis.AjouterMembre(chef)
	equipeEnnemis.AjouterMembre(berserker)
	equipeEnnemis.AjouterMembre(shaman)

	// Initialiser les stats
	for _, u := range equipeHeros.Membres() {
		g.stats[string(u.ID())] = &CombatStats{}
	}
	for _, u := range equipeEnnemis.Membres() {
		g.stats[string(u.ID())] = &CombatStats{}
	}

	combat, err := domain.NewCombat("demo-combat-advanced", []*domain.Equipe{equipeHeros, equipeEnnemis}, grille)
	if err != nil {
		return err
	}

	if err := combat.Demarrer(); err != nil {
		return err
	}

	g.combat = combat
	g.equipeHeros = equipeHeros
	g.equipeEnnemis = equipeEnnemis

	return nil
}

func (g *GameDemo) afficherIntro() {
	fmt.Println(ColorBold + ColorYellow + "📖 SCÉNARIO" + ColorReset)
	fmt.Println(ColorYellow + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Println("Trois héros courageux font face à une horde gobeline menée")
	fmt.Println("par un chef rusé. Le destin du royaume se joue ici !")
	fmt.Println()

	fmt.Println(ColorGreen + "⚔️  HÉROS DE LUMIÈRE (Vous)" + ColorReset)
	fmt.Println(ColorGreen + "─────────────────────────────────" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		stats := u.Stats()
		role := g.getRoleDescription(u)
		fmt.Printf(ColorGreen+"  ⭐ %s"+ColorReset+" - %s\n", u.Nom(), role)
		fmt.Printf("     HP:%d ATK:%d DEF:%d MATK:%d ATH:%d%% SPD:%d MOV:%d\n",
			stats.HP, stats.ATK, stats.DEF, stats.MATK, stats.ATH, stats.SPD, stats.MOV)
		if len(u.Competences()) > 1 {
			fmt.Printf(ColorCyan + "     Compétences: " + ColorReset)
			for i, c := range u.Competences() {
				if c.ID() == "attaque-basique" {
					continue
				}
				if i > 1 {
					fmt.Print(", ")
				}
				fmt.Printf("%s", c.Nom())
			}
			fmt.Println()
		}
	}
	fmt.Println()

	fmt.Println(ColorRed + "👹 HORDE GOBELINE (IA)" + ColorReset)
	fmt.Println(ColorRed + "─────────────────────────────────" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		stats := u.Stats()
		role := g.getRoleDescription(u)
		fmt.Printf(ColorRed+"  💀 %s"+ColorReset+" - %s\n", u.Nom(), role)
		fmt.Printf("     HP:%d ATK:%d DEF:%d MATK:%d ATH:%d%% SPD:%d\n",
			stats.HP, stats.ATK, stats.DEF, stats.MATK, stats.ATH, stats.SPD)
	}
	fmt.Println()

	fmt.Println(ColorCyan + "🎮 COMMANDES DISPONIBLES" + ColorReset)
	fmt.Println(ColorCyan + "─────────────────────────────────" + ColorReset)
	fmt.Println("  attack <cible>         - Attaque de base")
	fmt.Println("  skill <nom> <cible>    - Utiliser une compétence")
	fmt.Println("  move <x> <y>           - Se déplacer")
	fmt.Println("  stats                  - Voir statistiques détaillées")
	fmt.Println("  pass                   - Passer son tour")
	fmt.Println("  help                   - Aide")
	fmt.Println("  quit                   - Quitter")
	fmt.Println()

	fmt.Println(ColorPurple + "💡 ASTUCE" + ColorReset)
	fmt.Println("Les IDs des unités: hero-paladin, hero-archer, hero-mage")
	fmt.Println("                    ennemi-chef, ennemi-berserker, ennemi-shaman")
	fmt.Println()

	g.attendreAppui()
}

func (g *GameDemo) getRoleDescription(u *domain.Unite) string {
	stats := u.Stats()
	nom := u.Nom()

	switch {
	case strings.Contains(nom, "Paladin"):
		return "🛡️  Tank - Encaisse les dégâts"
	case strings.Contains(nom, "Archer"):
		return "🏹 DPS Physique - Précis et mortel"
	case strings.Contains(nom, "Mage"):
		return "✨ DPS Magique - Sorts dévastateurs"
	case strings.Contains(nom, "Chef"):
		return "👑 Leader - Équilibré et dangereux"
	case strings.Contains(nom, "Berserker"):
		return "⚡ Assaut - Rapide et brutal"
	case strings.Contains(nom, "Shaman"):
		return "🔮 Sorcier - Magie noire"
	}

	if stats.ATK > stats.MATK {
		return "⚔️  Combattant"
	}
	return "🔮 Mage"
}

func (g *GameDemo) attendreAppui() {
	fmt.Print(ColorYellow + "Appuyez sur ENTRÉE pour commencer..." + ColorReset)
	g.reader.ReadString('\n')
	fmt.Println()
}

func (g *GameDemo) boucleDeJeu() {
	for {
		resultat := g.combat.VerifierConditionsVictoire()
		if resultat != "CONTINUE" {
			g.afficherFinCombat(resultat)
			break
		}

		g.tourActuel++
		g.afficherEtatCombat()

		// Tour des héros
		for _, unite := range g.equipeHeros.MembresVivants() {
			g.stats[string(unite.ID())].ToursJoues++
			if !g.jouerTourHero(unite) {
				return // Quit
			}

			// Vérifier victoire après chaque action
			if g.combat.VerifierConditionsVictoire() != "CONTINUE" {
				break
			}
		}

		// Vérifier victoire
		if g.combat.VerifierConditionsVictoire() != "CONTINUE" {
			continue
		}

		// Tour des ennemis
		fmt.Println()
		fmt.Println(ColorRed + "👹 Phase Ennemie..." + ColorReset)
		time.Sleep(500 * time.Millisecond)

		for _, unite := range g.equipeEnnemis.MembresVivants() {
			g.stats[string(unite.ID())].ToursJoues++
			g.jouerTourIA(unite)
			time.Sleep(800 * time.Millisecond)

			if g.combat.VerifierConditionsVictoire() != "CONTINUE" {
				break
			}
		}

		// Nouveau tour
		for _, unite := range g.equipeHeros.Membres() {
			unite.NouveauTour()
		}
		for _, unite := range g.equipeEnnemis.Membres() {
			unite.NouveauTour()
		}

		fmt.Println()
		fmt.Println(ColorYellow + "═══ Fin du tour " + fmt.Sprint(g.tourActuel) + " ═══" + ColorReset)
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}

func (g *GameDemo) jouerTourHero(unite *domain.Unite) bool {
	fmt.Println()
	fmt.Println(ColorBold + ColorGreen + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Printf(ColorGreen+"🗡️  Tour de %s\n"+ColorReset, unite.Nom())
	fmt.Println(ColorGreen + "═══════════════════════════════════════════════" + ColorReset)

	g.afficherInfoUnite(unite)
	fmt.Println()

	for {
		fmt.Print(ColorCyan + "> " + ColorReset)
		input, _ := g.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			return false
		}

		if input == "help" {
			g.afficherAide()
			continue
		}

		if input == "map" {
			g.afficherGrille()
			continue
		}

		if input == "stats" {
			g.afficherStatistiques()
			continue
		}

		if input == "pass" {
			fmt.Println(ColorYellow + "⏭️  Tour passé" + ColorReset)
			// Vider le buffer pour éviter les bugs de commande
			g.reader = bufio.NewReader(os.Stdin)
			return true
		}

		if g.executerCommande(unite, input) {
			return true
		}

		fmt.Println(ColorRed + "❌ Commande invalide. Tapez 'help' pour l'aide." + ColorReset)
	}
}

func (g *GameDemo) executerCommande(unite *domain.Unite, input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	cmd := parts[0]

	switch cmd {
	case "attack":
		if len(parts) < 2 {
			fmt.Println(ColorRed + "Usage: attack <cible-id>" + ColorReset)
			return false
		}
		return g.executerAttaque(unite, parts[1])

	case "skill":
		if len(parts) < 3 {
			fmt.Println(ColorRed + "Usage: skill <competence> <cible-id>" + ColorReset)
			return false
		}
		return g.executerCompetence(unite, parts[1], parts[2])

	case "move":
		if len(parts) < 3 {
			fmt.Println(ColorRed + "Usage: move <x> <y>" + ColorReset)
			return false
		}
		return g.executerDeplacement(unite, parts[1], parts[2])

	default:
		return false
	}
}

func (g *GameDemo) executerAttaque(attaquant *domain.Unite, cibleID string) bool {
	cible := g.trouverUnite(cibleID)
	if cible == nil {
		fmt.Printf(ColorRed+"❌ Cible '%s' introuvable\n"+ColorReset, cibleID)
		return false
	}

	if cible.EstEliminee() {
		fmt.Println(ColorRed + "❌ La cible est déjà éliminée" + ColorReset)
		return false
	}

	// Déterminer la portée d'attaque selon le type d'unité
	porteeAttaque := 1 // Mêlée par défaut
	typeArme := "⚔️"

	// Archer a une portée de 4 pour ses attaques de base
	if attaquant.ID() == "hero-archer" {
		porteeAttaque = 4
		typeArme = "🏹"
	}

	// Vérifier la portée
	posAttaquant := attaquant.Position()
	posCible := cible.Position()
	distance := abs(posAttaquant.X()-posCible.X()) + abs(posAttaquant.Y()-posCible.Y())

	if distance > porteeAttaque {
		if porteeAttaque == 1 {
			fmt.Printf(ColorRed+"❌ Cible trop éloignée! Distance:%d (portée CAC:1)\n"+ColorReset, distance)
			fmt.Printf(ColorYellow + "💡 Utilisez 'move' pour vous rapprocher ou une compétence à distance\n" + ColorReset)
		} else {
			fmt.Printf(ColorRed+"❌ Cible trop éloignée! Distance:%d (portée:%d)\n"+ColorReset, distance, porteeAttaque)
		}
		return false
	}

	g.stats[string(attaquant.ID())].AttaquesTotal++

	// Vérifier chance de toucher
	ath := attaquant.Stats().ATH
	chanceToucher := rand.Intn(100) + 1

	if chanceToucher > ath {
		g.stats[string(attaquant.ID())].AttaquesRatees++
		fmt.Printf(ColorYellow+"%s  %s attaque %s mais "+ColorBold+"RATE"+ColorReset+ColorYellow+"! (ATH:%d%% vs jet:%d)\n"+ColorReset,
			typeArme, attaquant.Nom(), cible.Nom(), ath, chanceToucher)
		return true
	}

	g.stats[string(attaquant.ID())].AttaquesReussies++

	competence := attaquant.ObtenirCompetenceParDefaut()
	degats := g.combat.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)

	g.stats[string(attaquant.ID())].DegatsInfliges += degats
	g.stats[string(cible.ID())].DegatsRecus += degats

	cible.RecevoirDegats(degats)

	fmt.Printf(ColorYellow+"%s  %s attaque %s et inflige "+ColorBold+"%d dégâts"+ColorReset+ColorYellow+"! (ATH:%d%%, portée:%d)\n"+ColorReset,
		typeArme, attaquant.Nom(), cible.Nom(), degats, ath, porteeAttaque)

	if cible.EstEliminee() {
		fmt.Printf(ColorRed+"💀 "+ColorBold+"%s a été vaincu!"+ColorReset+"\n", cible.Nom())
	} else {
		g.afficherBarreHP(cible, ColorYellow)
	}

	return true
}

func (g *GameDemo) executerCompetence(attaquant *domain.Unite, skillName string, cibleID string) bool {
	var competence *domain.Competence
	for _, c := range attaquant.Competences() {
		if strings.ToLower(string(c.ID())) == strings.ToLower(skillName) ||
			strings.ToLower(c.Nom()) == strings.ToLower(skillName) {
			competence = c
			break
		}
	}

	if competence == nil {
		fmt.Printf(ColorRed+"❌ Compétence '%s' introuvable\n"+ColorReset, skillName)
		fmt.Println(ColorCyan + "Compétences disponibles:" + ColorReset)
		for _, c := range attaquant.Competences() {
			if c.ID() == "attaque-basique" {
				continue
			}
			fmt.Printf("  • %s (ID: %s, MP:%d, Cooldown:%d)\n", c.Nom(), c.ID(), c.CoutMP(), c.Cooldown())
		}
		return false
	}

	if !attaquant.PeutUtiliserCompetence(competence.ID()) {
		if competence.EstEnCooldown() {
			fmt.Printf(ColorRed+"❌ %s est en cooldown (%d tours restants)\n"+ColorReset,
				competence.Nom(), competence.CooldownActuel())
		} else if attaquant.StatsActuelles().MP < competence.CoutMP() {
			fmt.Printf(ColorRed+"❌ MP insuffisants (requis: %d, disponible: %d)\n"+ColorReset,
				competence.CoutMP(), attaquant.StatsActuelles().MP)
		}
		return false
	}

	cible := g.trouverUnite(cibleID)
	if cible == nil {
		fmt.Printf(ColorRed+"❌ Cible '%s' introuvable\n"+ColorReset, cibleID)
		return false
	}

	if cible.EstEliminee() {
		fmt.Println(ColorRed + "❌ La cible est déjà éliminée" + ColorReset)
		return false
	}

	g.stats[string(attaquant.ID())].AttaquesTotal++
	g.stats[string(attaquant.ID())].CompetencesUsees++

	// Vérifier chance de toucher (magie +10%)
	ath := attaquant.Stats().ATH + 10
	if ath > 100 {
		ath = 100
	}
	chanceToucher := rand.Intn(100) + 1

	if chanceToucher > ath {
		// Consommer quand même les ressources
		if err := attaquant.UtiliserCompetence(competence.ID()); err != nil {
			fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
			return false
		}
		g.stats[string(attaquant.ID())].AttaquesRatees++
		fmt.Printf(ColorPurple+"✨ %s lance %s sur %s mais "+ColorBold+"RATE"+ColorReset+ColorPurple+"! (ATH:%d%% vs jet:%d)\n"+ColorReset,
			attaquant.Nom(), competence.Nom(), cible.Nom(), ath, chanceToucher)
		return true
	}

	g.stats[string(attaquant.ID())].AttaquesReussies++

	if err := attaquant.UtiliserCompetence(competence.ID()); err != nil {
		fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
		return false
	}

	// Gérer les différents types de compétences
	typeComp := competence.Type()

	if typeComp == domain.CompetenceSoin {
		// Compétence de soin
		soin := int(float64(competence.DegatsBase()) * competence.Modificateur())
		cible.RecevoirSoin(soin)
		fmt.Printf(ColorGreen+"💚 %s lance "+ColorBold+"%s"+ColorReset+ColorGreen+" sur %s et restaure "+ColorBold+"%d HP"+ColorReset+ColorGreen+"!\n"+ColorReset,
			attaquant.Nom(), competence.Nom(), cible.Nom(), soin)
		g.afficherBarreHP(cible, ColorGreen)

	} else if typeComp == domain.CompetenceDebuff {
		// Compétences de statut (Sleep, Paralyse, Poison)
		skillID := competence.ID()
		if skillID == "sleep" {
			// Sommeil
			statut := shared.NewStatut(shared.StatutSommeil, 2, 0)
			cible.AjouterStatut(statut)
			fmt.Printf(ColorPurple+"😴 %s lance "+ColorBold+"%s"+ColorReset+ColorPurple+" sur %s! "+ColorBold+"Endormi 2 tours!"+ColorReset+"\n",
				attaquant.Nom(), competence.Nom(), cible.Nom())
		} else if skillID == "paralyze" {
			// Paralysie
			statut := shared.NewStatut(shared.StatutStun, 2, 0)
			cible.AjouterStatut(statut)
			fmt.Printf(ColorPurple+"⚡ %s lance "+ColorBold+"%s"+ColorReset+ColorPurple+" sur %s! "+ColorBold+"Paralysé 2 tours!"+ColorReset+"\n",
				attaquant.Nom(), competence.Nom(), cible.Nom())
		} else if skillID == "poison" {
			// Poison (dégâts sur la durée)
			statut := shared.NewStatut(shared.StatutPoison, 3, 10)
			cible.AjouterStatut(statut)
			fmt.Printf(ColorPurple+"☠️  %s lance "+ColorBold+"%s"+ColorReset+ColorPurple+" sur %s! "+ColorBold+"Empoisonné (10 dégâts/tour, 3 tours)!"+ColorReset+"\n",
				attaquant.Nom(), competence.Nom(), cible.Nom())
		}

	} else if typeComp == domain.CompetenceBuff {
		// Boost MATK
		if competence.ID() == "boost" {
			modificateur := &shared.ModificateurStat{
				Stat:   "MATK",
				Valeur: 15,
			}
			cible.AppliquerModificateurStat(modificateur)
			fmt.Printf(ColorCyan+"✨ %s lance "+ColorBold+"%s"+ColorReset+ColorCyan+" sur %s! "+ColorBold+"MATK +15 pendant 3 tours!"+ColorReset+"\n",
				attaquant.Nom(), competence.Nom(), cible.Nom())
		} else {
			// Autres buffs
			fmt.Printf(ColorCyan+"✨ %s lance "+ColorBold+"%s"+ColorReset+ColorCyan+" sur %s!\n"+ColorReset,
				attaquant.Nom(), competence.Nom(), cible.Nom())
		}
	} else {
		// Compétences de dégâts (Magie, Attaque)

		// Fireball a un effet AOE en croix
		if competence.ID() == "fireball" {
			posCible := cible.Position()
			ciblesAOE := []*domain.Unite{cible}

			// Définir les positions en croix (haut, bas, gauche, droite)
			directions := []struct{ dx, dy int }{
				{-1, 0}, // Haut
				{1, 0},  // Bas
				{0, -1}, // Gauche
				{0, 1},  // Droite
			}

			// Chercher les unités dans les cases adjacentes en croix
			for _, dir := range directions {
				posAdjacente, _ := shared.NewPosition(posCible.X()+dir.dx, posCible.Y()+dir.dy)

				// Vérifier toutes les unités (héros et ennemis)
				toutesLesUnites := []*domain.Unite{}
				toutesLesUnites = append(toutesLesUnites, g.equipeHeros.Membres()...)
				toutesLesUnites = append(toutesLesUnites, g.equipeEnnemis.Membres()...)

				for _, unite := range toutesLesUnites {
					if !unite.EstEliminee() && unite.ID() != cible.ID() {
						posUnite := unite.Position()
						if posUnite.X() == posAdjacente.X() && posUnite.Y() == posAdjacente.Y() {
							ciblesAOE = append(ciblesAOE, unite)
						}
					}
				}
			}

			// Appliquer les dégâts à toutes les cibles touchées
			fmt.Printf(ColorPurple+"🔥 %s lance "+ColorBold+"%s"+ColorReset+ColorPurple+" en zone! (ATH:%d%%)\n"+ColorReset,
				attaquant.Nom(), competence.Nom(), ath)

			for _, cibleAOE := range ciblesAOE {
				degats := g.combat.GetDamageCalculator().CalculerDegats(attaquant, cibleAOE, competence)

				g.stats[string(attaquant.ID())].DegatsInfliges += degats
				g.stats[string(cibleAOE.ID())].DegatsRecus += degats

				cibleAOE.RecevoirDegats(degats)

				fmt.Printf(ColorPurple+"  💥 %s subit "+ColorBold+"%d dégâts"+ColorReset+ColorPurple+"!\n"+ColorReset,
					cibleAOE.Nom(), degats)

				if cibleAOE.EstEliminee() {
					fmt.Printf(ColorRed+"  💀 "+ColorBold+"%s a été vaincu!"+ColorReset+"\n", cibleAOE.Nom())
				}
			}

			// Afficher les HP restants des survivants
			for _, cibleAOE := range ciblesAOE {
				if !cibleAOE.EstEliminee() {
					g.afficherBarreHP(cibleAOE, ColorPurple)
				}
			}
		} else {
			// Compétences de dégâts normales (cible unique)
			degats := g.combat.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)

			g.stats[string(attaquant.ID())].DegatsInfliges += degats
			g.stats[string(cible.ID())].DegatsRecus += degats

			cible.RecevoirDegats(degats)

			fmt.Printf(ColorPurple+"✨ %s lance "+ColorBold+"%s"+ColorReset+ColorPurple+" sur %s et inflige "+ColorBold+"%d dégâts"+ColorReset+ColorPurple+"! (ATH:%d%%)\n"+ColorReset,
				attaquant.Nom(), competence.Nom(), cible.Nom(), degats, ath)

			if cible.EstEliminee() {
				fmt.Printf(ColorRed+"💀 "+ColorBold+"%s a été vaincu!"+ColorReset+"\n", cible.Nom())
			} else {
				g.afficherBarreHP(cible, ColorPurple)
			}
		}
	}

	return true
}

func (g *GameDemo) executerDeplacement(unite *domain.Unite, xStr, yStr string) bool {
	x, err1 := strconv.Atoi(xStr)
	y, err2 := strconv.Atoi(yStr)

	if err1 != nil || err2 != nil {
		fmt.Println(ColorRed + "❌ Coordonnées invalides" + ColorReset)
		return false
	}

	nouvellePos, _ := shared.NewPosition(x, y)

	if x < 0 || x >= 10 || y < 0 || y >= 10 {
		fmt.Println(ColorRed + "❌ Position hors de la grille (0-9)" + ColorReset)
		return false
	}

	posActuelle := unite.Position()
	distance := abs(nouvellePos.X()-posActuelle.X()) + abs(nouvellePos.Y()-posActuelle.Y())

	if distance > unite.Stats().MOV {
		fmt.Printf(ColorRed+"❌ Trop loin! Distance:%d, Mouvement:%d\n"+ColorReset,
			distance, unite.Stats().MOV)
		return false
	}

	// Vérifier qu'aucune unité n'occupe déjà cette case
	for _, u := range g.equipeHeros.Membres() {
		if !u.EstEliminee() && u.ID() != unite.ID() {
			pos := u.Position()
			if pos.X() == x && pos.Y() == y {
				fmt.Printf(ColorRed+"❌ Case occupée par %s!\n"+ColorReset, u.Nom())
				return false
			}
		}
	}
	for _, u := range g.equipeEnnemis.Membres() {
		if !u.EstEliminee() {
			pos := u.Position()
			if pos.X() == x && pos.Y() == y {
				fmt.Printf(ColorRed+"❌ Case occupée par %s!\n"+ColorReset, u.Nom())
				return false
			}
		}
	}

	if err := unite.SeDeplacer(nouvellePos, distance); err != nil {
		fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
		return false
	}

	fmt.Printf(ColorCyan+"🏃 %s se déplace en (%d, %d)\n"+ColorReset, unite.Nom(), x, y)
	return true
}

func (g *GameDemo) jouerTourIA(unite *domain.Unite) {
	fmt.Printf(ColorRed+"👹 Tour de %s\n"+ColorReset, unite.Nom())

	// IA améliorée: prioriser les cibles faibles
	var ciblePrioritaire *domain.Unite
	prioriteMin := 999999

	for _, hero := range g.equipeHeros.MembresVivants() {
		posIA := unite.Position()
		posHero := hero.Position()
		distance := abs(posIA.X()-posHero.X()) + abs(posIA.Y()-posHero.Y())

		// Priorité = HP restants + distance (cible faible et proche = mieux)
		priorite := hero.HPActuels() + (distance * 5)

		if priorite < prioriteMin {
			prioriteMin = priorite
			ciblePrioritaire = hero
		}
	}

	if ciblePrioritaire == nil {
		return
	}

	posIA := unite.Position()
	posCible := ciblePrioritaire.Position()
	distance := abs(posIA.X()-posCible.X()) + abs(posIA.Y()-posCible.Y())

	// Essayer d'utiliser une compétence si disponible
	for _, comp := range unite.Competences() {
		if comp.ID() == "attaque-basique" {
			continue
		}

		if unite.PeutUtiliserCompetence(comp.ID()) && distance <= comp.Portee() {
			g.stats[string(unite.ID())].AttaquesTotal++
			g.stats[string(unite.ID())].CompetencesUsees++

			ath := unite.Stats().ATH + 10
			if ath > 100 {
				ath = 100
			}
			chanceToucher := rand.Intn(100) + 1

			if chanceToucher > ath {
				unite.UtiliserCompetence(comp.ID())
				g.stats[string(unite.ID())].AttaquesRatees++
				fmt.Printf(ColorRed+"✨ %s lance %s mais RATE! (ATH:%d%%)\n"+ColorReset,
					unite.Nom(), comp.Nom(), ath)
				return
			}

			g.stats[string(unite.ID())].AttaquesReussies++
			unite.UtiliserCompetence(comp.ID())
			degats := g.combat.GetDamageCalculator().CalculerDegats(unite, ciblePrioritaire, comp)

			g.stats[string(unite.ID())].DegatsInfliges += degats
			g.stats[string(ciblePrioritaire.ID())].DegatsRecus += degats

			ciblePrioritaire.RecevoirDegats(degats)

			fmt.Printf(ColorRed+"✨ %s lance %s sur %s et inflige %d dégâts!\n"+ColorReset,
				unite.Nom(), comp.Nom(), ciblePrioritaire.Nom(), degats)

			if ciblePrioritaire.EstEliminee() {
				fmt.Printf(ColorRed+"💀 %s a été vaincu!\n"+ColorReset, ciblePrioritaire.Nom())
			}
			return
		}
	}

	// Si à portée, attaquer
	if distance <= 1 {
		g.stats[string(unite.ID())].AttaquesTotal++

		ath := unite.Stats().ATH
		chanceToucher := rand.Intn(100) + 1

		if chanceToucher > ath {
			g.stats[string(unite.ID())].AttaquesRatees++
			fmt.Printf(ColorRed+"⚔️  %s attaque %s mais RATE! (ATH:%d%% vs jet:%d)\n"+ColorReset,
				unite.Nom(), ciblePrioritaire.Nom(), ath, chanceToucher)
		} else {
			g.stats[string(unite.ID())].AttaquesReussies++
			competence := unite.ObtenirCompetenceParDefaut()
			degats := g.combat.GetDamageCalculator().CalculerDegats(unite, ciblePrioritaire, competence)

			g.stats[string(unite.ID())].DegatsInfliges += degats
			g.stats[string(ciblePrioritaire.ID())].DegatsRecus += degats

			ciblePrioritaire.RecevoirDegats(degats)

			fmt.Printf(ColorRed+"⚔️  %s attaque %s et inflige %d dégâts! (ATH:%d%%)\n"+ColorReset,
				unite.Nom(), ciblePrioritaire.Nom(), degats, ath)

			if ciblePrioritaire.EstEliminee() {
				fmt.Printf(ColorRed+"💀 %s a été vaincu!\n"+ColorReset, ciblePrioritaire.Nom())
			}
		}
	} else {
		// Se rapprocher intelligemment
		newX := posIA.X()
		newY := posIA.Y()

		diffX := posCible.X() - posIA.X()
		diffY := posCible.Y() - posIA.Y()

		// Se déplacer dans la direction la plus efficace
		if abs(diffX) > abs(diffY) {
			if diffX > 0 {
				newX++
			} else if diffX < 0 {
				newX--
			}
		} else {
			if diffY > 0 {
				newY++
			} else if diffY < 0 {
				newY--
			}
		}

		nouvellePos, _ := shared.NewPosition(newX, newY)
		unite.DeplacerVers(nouvellePos)

		fmt.Printf(ColorRed+"🏃 %s se rapproche de %s en (%d, %d)\n"+ColorReset,
			unite.Nom(), ciblePrioritaire.Nom(), newX, newY)
	}
}

func (g *GameDemo) afficherEtatCombat() {
	fmt.Println()
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Printf(ColorBold+ColorWhite+"        TOUR %d - SITUATION\n"+ColorReset, g.tourActuel)
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Println()

	// Héros
	fmt.Println(ColorGreen + "⚔️  HÉROS" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		g.afficherBarreHP(u, ColorGreen)
	}
	fmt.Println()

	// Ennemis
	fmt.Println(ColorRed + "👹 ENNEMIS" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		g.afficherBarreHP(u, ColorRed)
	}
	fmt.Println()
}

func (g *GameDemo) afficherBarreHP(unite *domain.Unite, couleur string) {
	if unite.EstEliminee() {
		fmt.Printf("  %s💀 %s [VAINCU]%s\n", ColorRed, unite.Nom(), ColorReset)
		return
	}

	hpActuel := unite.HPActuels()
	hpMax := unite.Stats().HP
	pourcentage := float64(hpActuel) / float64(hpMax)

	barreLength := 25
	rempli := int(pourcentage * float64(barreLength))

	// Couleur de la barre selon HP
	couleurBarre := ColorGreen
	if pourcentage < 0.3 {
		couleurBarre = ColorRed
	} else if pourcentage < 0.6 {
		couleurBarre = ColorYellow
	}

	barre := "["
	for i := 0; i < barreLength; i++ {
		if i < rempli {
			barre += "█"
		} else {
			barre += "░"
		}
	}
	barre += "]"

	pos := unite.Position()
	mpActuel := unite.StatsActuelles().MP
	mpMax := unite.Stats().MP

	fmt.Printf("  %s%-18s%s %s%s%s %d/%d HP",
		couleur, unite.Nom(), ColorReset,
		couleurBarre, barre, ColorReset,
		hpActuel, hpMax)

	if mpMax > 0 {
		fmt.Printf(" | MP:%d/%d", mpActuel, mpMax)
	}
	fmt.Printf(" | (%d,%d)\n", pos.X(), pos.Y())
}

func (g *GameDemo) afficherInfoUnite(unite *domain.Unite) {
	stats := unite.StatsActuelles()
	statsBase := unite.Stats()
	pos := unite.Position()

	fmt.Printf("📊 Stats: HP:%d/%d MP:%d/%d ATK:%d DEF:%d MATK:%d ATH:%d%% SPD:%d\n",
		unite.HPActuels(), statsBase.HP, stats.MP, statsBase.MP,
		stats.ATK, stats.DEF, stats.MATK, statsBase.ATH, stats.SPD)
	fmt.Printf("📍 Position: (%d, %d) | Mouvement: %d cases\n", pos.X(), pos.Y(), stats.MOV)

	if len(unite.Competences()) > 1 {
		fmt.Println("✨ Compétences:")
		for _, c := range unite.Competences() {
			if c.ID() == "attaque-basique" {
				continue
			}
			status := ColorGreen + "✓" + ColorReset
			if c.EstEnCooldown() {
				status = ColorRed + fmt.Sprintf("CD:%d", c.CooldownActuel()) + ColorReset
			} else if stats.MP < c.CoutMP() {
				status = ColorYellow + "MP↓" + ColorReset
			}
			fmt.Printf("  • %s (ID:%s, MP:%d, Portée:%d) %s\n",
				c.Nom(), c.ID(), c.CoutMP(), c.Portee(), status)
		}
	}
}

func (g *GameDemo) afficherGrille() {
	fmt.Println()
	fmt.Println(ColorBold + ColorCyan + "🗺️  GRILLE TACTIQUE 10x10" + ColorReset)
	fmt.Println(ColorCyan + "═══════════════════════════════════════════════" + ColorReset)

	// Créer une carte des positions
	grille := make([][]string, 10)
	for i := range grille {
		grille[i] = make([]string, 10)
		for j := range grille[i] {
			grille[i][j] = "·"
		}
	}

	// Placer les héros
	for _, u := range g.equipeHeros.Membres() {
		if !u.EstEliminee() {
			pos := u.Position()
			x, y := pos.X(), pos.Y()
			if x >= 0 && x < 10 && y >= 0 && y < 10 {
				switch u.ID() {
				case "hero-paladin":
					grille[x][y] = ColorGreen + "P" + ColorReset
				case "hero-archer":
					grille[x][y] = ColorGreen + "A" + ColorReset
				case "hero-mage":
					grille[x][y] = ColorGreen + "M" + ColorReset
				}
			}
		}
	}

	// Placer les ennemis
	for _, u := range g.equipeEnnemis.Membres() {
		if !u.EstEliminee() {
			pos := u.Position()
			x, y := pos.X(), pos.Y()
			if x >= 0 && x < 10 && y >= 0 && y < 10 {
				switch u.ID() {
				case "ennemi-chef":
					grille[x][y] = ColorRed + "C" + ColorReset
				case "ennemi-berserker":
					grille[x][y] = ColorRed + "B" + ColorReset
				case "ennemi-shaman":
					grille[x][y] = ColorRed + "S" + ColorReset
				}
			}
		}
	}

	// Afficher la grille
	fmt.Print("   ")
	for j := 0; j < 10; j++ {
		fmt.Printf("%d ", j)
	}
	fmt.Println()

	for i := 0; i < 10; i++ {
		fmt.Printf(" %d ", i)
		for j := 0; j < 10; j++ {
			fmt.Printf("%s ", grille[i][j])
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(ColorGreen + "  HÉROS: " + ColorReset + "P=Paladin, A=Archer, M=Mage")
	fmt.Println(ColorRed + "  ENNEMIS: " + ColorReset + "C=Chef, B=Berserker, S=Shaman")
	fmt.Println()
}

func (g *GameDemo) afficherStatistiques() {
	fmt.Println()
	fmt.Println(ColorBold + ColorCyan + "📊 STATISTIQUES DE COMBAT" + ColorReset)
	fmt.Println(ColorCyan + "═══════════════════════════════════════════════" + ColorReset)

	fmt.Println(ColorGreen + "\n⚔️  HÉROS" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		g.afficherStatsUnite(u)
	}

	fmt.Println(ColorRed + "\n👹 ENNEMIS" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		g.afficherStatsUnite(u)
	}
	fmt.Println()
}

func (g *GameDemo) afficherStatsUnite(u *domain.Unite) {
	stats := g.stats[string(u.ID())]
	precision := 0.0
	if stats.AttaquesTotal > 0 {
		precision = float64(stats.AttaquesReussies) / float64(stats.AttaquesTotal) * 100
	}

	status := ColorGreen + "Vivant" + ColorReset
	if u.EstEliminee() {
		status = ColorRed + "Vaincu" + ColorReset
	}

	fmt.Printf("\n  %s (%s)\n", u.Nom(), status)
	fmt.Printf("    Tours joués: %d\n", stats.ToursJoues)
	fmt.Printf("    Attaques: %d (✓%d ✗%d) - Précision: %.1f%%\n",
		stats.AttaquesTotal, stats.AttaquesReussies, stats.AttaquesRatees, precision)
	fmt.Printf("    Compétences utilisées: %d\n", stats.CompetencesUsees)
	fmt.Printf("    Dégâts infligés: %d\n", stats.DegatsInfliges)
	fmt.Printf("    Dégâts subis: %d\n", stats.DegatsRecus)
}

func (g *GameDemo) afficherAide() {
	fmt.Println()
	fmt.Println(ColorCyan + "🎮 AIDE - COMMANDES DISPONIBLES" + ColorReset)
	fmt.Println(ColorCyan + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Println("  attack <cible-id>       - Attaque de base")
	fmt.Println("  skill <nom> <cible-id>  - Utiliser une compétence")
	fmt.Println("  move <x> <y>            - Se déplacer sur la grille")
	fmt.Println("  map                     - Afficher la grille tactique ⭐ NOUVEAU")
	fmt.Println("  stats                   - Voir les statistiques détaillées")
	fmt.Println("  pass                    - Passer son tour")
	fmt.Println("  help                    - Afficher cette aide")
	fmt.Println("  quit                    - Quitter le jeu")
	fmt.Println()
	fmt.Println(ColorYellow + "💡 EXEMPLES" + ColorReset)
	fmt.Println("  attack ennemi-chef")
	fmt.Println("  skill fireball ennemi-shaman")
	fmt.Println("  skill precision-shot ennemi-berserker")
	fmt.Println("  move 5 5")
	fmt.Println()
	fmt.Println(ColorPurple + "🎯 IDS DES UNITÉS" + ColorReset)
	fmt.Println(ColorGreen + "  Héros:" + ColorReset)
	fmt.Println("    hero-paladin, hero-archer, hero-mage")
	fmt.Println(ColorRed + "  Ennemis:" + ColorReset)
	fmt.Println("    ennemi-chef, ennemi-berserker, ennemi-shaman")
	fmt.Println()
}

func (g *GameDemo) afficherFinCombat(resultat string) {
	fmt.Println()
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Println()

	if resultat == "VICTORY" {
		fmt.Println(ColorBold + ColorGreen + "        🎉 VICTOIRE HÉROÏQUE! 🎉" + ColorReset)
		fmt.Println()
		fmt.Println(ColorGreen + "Les héros ont triomphé de la horde gobeline!" + ColorReset)
		fmt.Println(ColorGreen + "Le royaume peut dormir tranquille cette nuit." + ColorReset)
	} else if resultat == "DEFEAT" {
		fmt.Println(ColorBold + ColorRed + "        💀 DÉFAITE AMÈRE 💀" + ColorReset)
		fmt.Println()
		fmt.Println(ColorRed + "Les gobelins ont vaincu les héros..." + ColorReset)
		fmt.Println(ColorRed + "L'obscurité s'abat sur le royaume." + ColorReset)
	}

	fmt.Println()
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════════════" + ColorReset)
	fmt.Println()

	// Statistiques finales
	fmt.Println(ColorCyan + "📊 STATISTIQUES FINALES" + ColorReset)
	fmt.Println(ColorCyan + "─────────────────────────────────────" + ColorReset)
	fmt.Printf("Tours total: %d\n", g.tourActuel)
	fmt.Println()

	// MVP (Most Valuable Player)
	var mvp *domain.Unite
	var mvpDegats int

	for _, u := range append(g.equipeHeros.Membres(), g.equipeEnnemis.Membres()...) {
		stats := g.stats[string(u.ID())]
		if stats.DegatsInfliges > mvpDegats {
			mvpDegats = stats.DegatsInfliges
			mvp = u
		}
	}

	if mvp != nil {
		fmt.Println(ColorYellow + "🏆 MVP (Most Valuable Player)" + ColorReset)
		fmt.Printf("   %s - %d dégâts infligés\n", mvp.Nom(), mvpDegats)
		fmt.Println()
	}

	fmt.Println(ColorGreen + "HÉROS:" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		stats := g.stats[string(u.ID())]
		precision := 0.0
		if stats.AttaquesTotal > 0 {
			precision = float64(stats.AttaquesReussies) / float64(stats.AttaquesTotal) * 100
		}

		if u.EstEliminee() {
			fmt.Printf("  💀 %s - VAINCU\n", u.Nom())
		} else {
			fmt.Printf("  ⚔️  %s - %d/%d HP\n", u.Nom(), u.HPActuels(), u.Stats().HP)
		}
		fmt.Printf("      Dégâts: %d | Précision: %.1f%% | Compétences: %d\n",
			stats.DegatsInfliges, precision, stats.CompetencesUsees)
	}

	fmt.Println()
	fmt.Println(ColorRed + "ENNEMIS:" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		stats := g.stats[string(u.ID())]
		precision := 0.0
		if stats.AttaquesTotal > 0 {
			precision = float64(stats.AttaquesReussies) / float64(stats.AttaquesTotal) * 100
		}

		if u.EstEliminee() {
			fmt.Printf("  💀 %s - VAINCU\n", u.Nom())
		} else {
			fmt.Printf("  👹 %s - %d/%d HP\n", u.Nom(), u.HPActuels(), u.Stats().HP)
		}
		fmt.Printf("      Dégâts: %d | Précision: %.1f%%\n",
			stats.DegatsInfliges, precision)
	}
	fmt.Println()

	fmt.Println(ColorPurple + "Merci d'avoir joué ! 🎮" + ColorReset)
	fmt.Println()
}

func (g *GameDemo) trouverUnite(id string) *domain.Unite {
	for _, u := range g.equipeHeros.Membres() {
		if string(u.ID()) == id {
			return u
		}
	}
	for _, u := range g.equipeEnnemis.Membres() {
		if string(u.ID()) == id {
			return u
		}
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
