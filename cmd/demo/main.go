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

// Couleurs ANSI pour un rendu sympa
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
)

type GameDemo struct {
	combat        *domain.Combat
	equipeHeros   *domain.Equipe
	equipeEnnemis *domain.Equipe
	uniteActuelle *domain.Unite
	reader        *bufio.Reader
}

func main() {
	// Initialiser le générateur aléatoire
	rand.Seed(time.Now().UnixNano())

	fmt.Println(ColorBold + ColorCyan + "╔══════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║   AETHER ENGINE - COMBAT DEMO CLI   ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "╚══════════════════════════════════════╝" + ColorReset)
	fmt.Println()

	game := &GameDemo{
		reader: bufio.NewReader(os.Stdin),
	}

	if err := game.initialiserCombat(); err != nil {
		fmt.Printf(ColorRed+"Erreur initialisation: %v\n"+ColorReset, err)
		return
	}

	game.afficherIntro()
	game.boucleDeJeu()
}

func (g *GameDemo) initialiserCombat() error {
	// Créer la grille 8x8
	grille, _ := shared.NewGrilleCombat(8, 8)

	// === ÉQUIPE HÉROS ===
	joueurID := "player-1"
	equipeHeros, _ := domain.NewEquipe("team-heros", "Héros", "#00FF00", false, &joueurID)

	// Guerrier
	statsGuerrier := &shared.Stats{
		HP:      120,
		MP:      30,
		Stamina: 100,
		ATK:     25,
		DEF:     15,
		MATK:    5,
		MDEF:    8,
		SPD:     12,
		MOV:     4,
		ATH:     85, // 85% de chance de toucher
	}
	posGuerrier, _ := shared.NewPosition(1, 3)
	guerrier := domain.NewUnite("hero-guerrier", "Guerrier", "team-heros", statsGuerrier, posGuerrier)

	// Mage
	statsMage := &shared.Stats{
		HP:      80,
		MP:      100,
		Stamina: 60,
		ATK:     10,
		DEF:     8,
		MATK:    30,
		MDEF:    20,
		SPD:     10,
		MOV:     3,
		ATH:     90, // 90% de chance de toucher (magie plus précise)
	}
	posMage, _ := shared.NewPosition(1, 4)
	mage := domain.NewUnite("hero-mage", "Mage", "team-heros", statsMage, posMage)

	// Ajouter compétence Boule de Feu au Mage
	fireball := domain.NewCompetence(
		"fireball",
		"Boule de Feu",
		"Projette une boule de feu dévastatrice",
		domain.CompetenceMagie,
		5,                   // portée
		domain.ZoneEffet{},  // zone
		20,                  // coût MP
		0,                   // coût Stamina
		2,                   // cooldown
		35,                  // dégâts de base
		1.5,                 // modificateur
		domain.CibleEnnemis, // cibles
	)
	mage.AjouterCompetence(fireball)

	equipeHeros.AjouterMembre(guerrier)
	equipeHeros.AjouterMembre(mage)

	// === ÉQUIPE GOBELINS ===
	equipeGobelins, _ := domain.NewEquipe("team-gobelins", "Gobelins", "#FF0000", true, nil)

	// Gobelin Guerrier 1
	statsGobelin := &shared.Stats{
		HP:      70,
		MP:      10,
		Stamina: 80,
		ATK:     18,
		DEF:     10,
		MATK:    3,
		MDEF:    5,
		SPD:     15,
		MOV:     5,
		ATH:     75, // 75% de chance de toucher
	}
	posGob1, _ := shared.NewPosition(6, 3)
	gobelin1 := domain.NewUnite("gobelin-1", "Gobelin Guerrier", "team-gobelins", statsGobelin, posGob1)

	// Gobelin Archer
	statsArcher := &shared.Stats{
		HP:      60,
		MP:      20,
		Stamina: 70,
		ATK:     22,
		DEF:     8,
		MATK:    5,
		MDEF:    6,
		SPD:     18,
		MOV:     4,
		ATH:     80, // 80% de chance de toucher (archer précis)
	}
	posGob2, _ := shared.NewPosition(6, 4)
	gobelin2 := domain.NewUnite("gobelin-2", "Gobelin Archer", "team-gobelins", statsArcher, posGob2)

	equipeGobelins.AjouterMembre(gobelin1)
	equipeGobelins.AjouterMembre(gobelin2)

	// Créer le combat
	combat, err := domain.NewCombat("demo-combat", []*domain.Equipe{equipeHeros, equipeGobelins}, grille)
	if err != nil {
		return err
	}

	// Démarrer le combat
	if err := combat.Demarrer(); err != nil {
		return err
	}

	g.combat = combat
	g.equipeHeros = equipeHeros
	g.equipeEnnemis = equipeGobelins

	return nil
}

func (g *GameDemo) afficherIntro() {
	fmt.Println(ColorYellow + "📜 Scénario:" + ColorReset)
	fmt.Println("Deux héros courageux font face à une bande de gobelins malfaisants!")
	fmt.Println()

	fmt.Println(ColorGreen + "⚔️  ÉQUIPE HÉROS:" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		stats := u.Stats()
		fmt.Printf("  • %s - HP:%d ATK:%d DEF:%d MATK:%d ATH:%d%% SPD:%d\n",
			u.Nom(), stats.HP, stats.ATK, stats.DEF, stats.MATK, stats.ATH, stats.SPD)
	}
	fmt.Println()

	fmt.Println(ColorRed + "👹 ÉQUIPE GOBELINS:" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		stats := u.Stats()
		fmt.Printf("  • %s - HP:%d ATK:%d DEF:%d ATH:%d%% SPD:%d\n",
			u.Nom(), stats.HP, stats.ATK, stats.DEF, stats.ATH, stats.SPD)
	}
	fmt.Println()

	fmt.Println(ColorCyan + " Commandes disponibles:" + ColorReset)
	fmt.Println("  attack <cible>  - Attaquer une cible (ex: attack gobelin-1)")
	fmt.Println("  skill <nom> <cible> - Utiliser une compétence (ex: skill fireball gobelin-1)")
	fmt.Println("  move <x> <y>    - Se déplacer (ex: move 3 4)")
	fmt.Println("  pass            - Passer son tour")
	fmt.Println("  help            - Afficher l'aide")
	fmt.Println("  quit            - Quitter")
	fmt.Println()
}

func (g *GameDemo) boucleDeJeu() {
	tourJoueur := 0

	for {
		// Vérifier conditions de victoire
		resultat := g.combat.VerifierConditionsVictoire()
		if resultat != "CONTINUE" {
			g.afficherFinCombat(resultat)
			break
		}

		// Tour du joueur
		tourJoueur++
		g.afficherEtatCombat(tourJoueur)

		// Tour des héros
		for _, unite := range g.equipeHeros.MembresVivants() {
			if !g.jouerTourHero(unite) {
				return // Quit
			}
		}

		// Tour des gobelins (IA simple)
		for _, unite := range g.equipeEnnemis.MembresVivants() {
			g.jouerTourIA(unite)
		}

		// Nouveau tour pour toutes les unités
		for _, unite := range g.equipeHeros.Membres() {
			unite.NouveauTour()
		}
		for _, unite := range g.equipeEnnemis.Membres() {
			unite.NouveauTour()
		}
	}
}

func (g *GameDemo) jouerTourHero(unite *domain.Unite) bool {
	fmt.Println()
	fmt.Println(ColorBold + ColorGreen + "═══════════════════════════════════════" + ColorReset)
	fmt.Printf(ColorGreen+"🗡️  Tour de %s\n"+ColorReset, unite.Nom())
	fmt.Println(ColorGreen + "═══════════════════════════════════════" + ColorReset)

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

		if input == "pass" {
			fmt.Println(ColorYellow + "⏭️  Tour passé" + ColorReset)
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

	// Vérifier la portée (attaque de base = corps-à-corps = distance 1)
	posAttaquant := attaquant.Position()
	posCible := cible.Position()
	distance := abs(posAttaquant.X()-posCible.X()) + abs(posAttaquant.Y()-posCible.Y())

	if distance > 1 {
		fmt.Printf(ColorRed+"❌ Cible trop éloignée! Distance:%d (portée CAC:1)\n"+ColorReset, distance)
		fmt.Printf(ColorYellow + "💡 Utilisez 'move' pour vous rapprocher ou 'skill' pour une compétence à distance\n" + ColorReset)
		return false
	}

	// Vérifier chance de toucher
	ath := attaquant.Stats().ATH
	chanceToucher := rand.Intn(100) + 1 // 1-100

	if chanceToucher > ath {
		fmt.Printf(ColorYellow+"⚔️  %s attaque %s mais RATE! (ATH:%d%% vs jet:%d)\n"+ColorReset,
			attaquant.Nom(), cible.Nom(), ath, chanceToucher)
		return true
	}

	// Attaque de base
	competence := attaquant.ObtenirCompetenceParDefaut()
	degats := g.combat.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)

	cible.RecevoirDegats(degats)

	fmt.Printf(ColorYellow+"⚔️  %s attaque %s et inflige %d dégâts! (ATH:%d%%)\n"+ColorReset,
		attaquant.Nom(), cible.Nom(), degats, ath)

	if cible.EstEliminee() {
		fmt.Printf(ColorRed+"💀 %s a été vaincu!\n"+ColorReset, cible.Nom())
	} else {
		fmt.Printf(ColorYellow+"   %s: %d/%d HP\n"+ColorReset,
			cible.Nom(), cible.HPActuels(), cible.Stats().HP)
	}

	return true
}

func (g *GameDemo) executerCompetence(attaquant *domain.Unite, skillName string, cibleID string) bool {
	// Trouver la compétence
	var competence *domain.Competence
	for _, c := range attaquant.Competences() {
		if strings.ToLower(string(c.ID())) == strings.ToLower(skillName) {
			competence = c
			break
		}
	}

	if competence == nil {
		fmt.Printf(ColorRed+"❌ Compétence '%s' introuvable\n"+ColorReset, skillName)
		fmt.Println(ColorCyan + "Compétences disponibles:" + ColorReset)
		for _, c := range attaquant.Competences() {
			fmt.Printf("  • %s (MP:%d, Cooldown:%d)\n", c.ID(), c.CoutMP(), c.Cooldown())
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

	// Vérifier chance de toucher (magie plus fiable +10%)
	ath := attaquant.Stats().ATH + 10
	if ath > 100 {
		ath = 100
	}
	chanceToucher := rand.Intn(100) + 1 // 1-100

	if chanceToucher > ath {
		// Consommer quand même les ressources
		if err := attaquant.UtiliserCompetence(competence.ID()); err != nil {
			fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
			return false
		}
		fmt.Printf(ColorPurple+"✨ %s lance %s sur %s mais RATE! (ATH:%d%% vs jet:%d)\n"+ColorReset,
			attaquant.Nom(), competence.Nom(), cible.Nom(), ath, chanceToucher)
		return true
	}

	// Utiliser la compétence
	if err := attaquant.UtiliserCompetence(competence.ID()); err != nil {
		fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
		return false
	}

	degats := g.combat.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)
	cible.RecevoirDegats(degats)

	fmt.Printf(ColorPurple+"✨ %s lance %s sur %s et inflige %d dégâts! (ATH:%d%%)\n"+ColorReset,
		attaquant.Nom(), competence.Nom(), cible.Nom(), degats, ath)

	if cible.EstEliminee() {
		fmt.Printf(ColorRed+"💀 %s a été vaincu!\n"+ColorReset, cible.Nom())
	} else {
		fmt.Printf(ColorYellow+"   %s: %d/%d HP\n"+ColorReset,
			cible.Nom(), cible.HPActuels(), cible.Stats().HP)
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

	// Vérifier limites grille
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		fmt.Println(ColorRed + "❌ Position hors de la grille (0-7)" + ColorReset)
		return false
	}

	// Calculer distance Manhattan
	posActuelle := unite.Position()
	distance := abs(nouvellePos.X()-posActuelle.X()) + abs(nouvellePos.Y()-posActuelle.Y())

	if distance > unite.Stats().MOV {
		fmt.Printf(ColorRed+"❌ Trop loin! Distance:%d, Mouvement:%d\n"+ColorReset,
			distance, unite.Stats().MOV)
		return false
	}

	if err := unite.SeDeplacer(nouvellePos, distance); err != nil {
		fmt.Printf(ColorRed+"❌ Erreur: %v\n"+ColorReset, err)
		return false
	}

	fmt.Printf(ColorCyan+"🏃 %s se déplace en (%d, %d)\n"+ColorReset, unite.Nom(), x, y)
	return true
}

func (g *GameDemo) jouerTourIA(unite *domain.Unite) {
	fmt.Println()
	fmt.Printf(ColorRed+"👹 Tour de %s\n"+ColorReset, unite.Nom())

	// IA simple : attaquer le héros le plus proche
	var cibleProche *domain.Unite
	distanceMin := 999

	for _, hero := range g.equipeHeros.MembresVivants() {
		posIA := unite.Position()
		posHero := hero.Position()
		distance := abs(posIA.X()-posHero.X()) + abs(posIA.Y()-posHero.Y())

		if distance < distanceMin {
			distanceMin = distance
			cibleProche = hero
		}
	}

	if cibleProche == nil {
		return
	}

	// Si à portée, attaquer
	if distanceMin <= 1 {
		// Vérifier chance de toucher
		ath := unite.Stats().ATH
		chanceToucher := rand.Intn(100) + 1 // 1-100

		if chanceToucher > ath {
			fmt.Printf(ColorRed+"⚔️  %s attaque %s mais RATE! (ATH:%d%% vs jet:%d)\n"+ColorReset,
				unite.Nom(), cibleProche.Nom(), ath, chanceToucher)
		} else {
			competence := unite.ObtenirCompetenceParDefaut()
			degats := g.combat.GetDamageCalculator().CalculerDegats(unite, cibleProche, competence)
			cibleProche.RecevoirDegats(degats)

			fmt.Printf(ColorRed+"⚔️  %s attaque %s et inflige %d dégâts! (ATH:%d%%)\n"+ColorReset,
				unite.Nom(), cibleProche.Nom(), degats, ath)

			if cibleProche.EstEliminee() {
				fmt.Printf(ColorRed+"💀 %s a été vaincu!\n"+ColorReset, cibleProche.Nom())
			}
		}
	} else {
		// Se rapprocher
		posIA := unite.Position()
		posCible := cibleProche.Position()

		newX := posIA.X()
		newY := posIA.Y()

		if posIA.X() < posCible.X() {
			newX++
		} else if posIA.X() > posCible.X() {
			newX--
		}

		if posIA.Y() < posCible.Y() {
			newY++
		} else if posIA.Y() > posCible.Y() {
			newY--
		}

		nouvellePos, _ := shared.NewPosition(newX, newY)
		unite.DeplacerVers(nouvellePos)

		fmt.Printf(ColorRed+"🏃 %s se rapproche en (%d, %d)\n"+ColorReset, unite.Nom(), newX, newY)
	}
}

func (g *GameDemo) afficherEtatCombat(tour int) {
	fmt.Println()
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════" + ColorReset)
	fmt.Printf(ColorBold+ColorWhite+"        TOUR %d\n"+ColorReset, tour)
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════" + ColorReset)
	fmt.Println()

	// Héros
	fmt.Println(ColorGreen + "⚔️  HÉROS:" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		g.afficherBarreHP(u, ColorGreen)
	}
	fmt.Println()

	// Gobelins
	fmt.Println(ColorRed + "👹 GOBELINS:" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		g.afficherBarreHP(u, ColorRed)
	}
}

func (g *GameDemo) afficherBarreHP(unite *domain.Unite, couleur string) {
	if unite.EstEliminee() {
		fmt.Printf("  %s💀 [VAINCU]%s\n", ColorRed, ColorReset)
		return
	}

	hpActuel := unite.HPActuels()
	hpMax := unite.Stats().HP
	pourcentage := float64(hpActuel) / float64(hpMax)

	barreLength := 20
	rempli := int(pourcentage * float64(barreLength))

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
	fmt.Printf("  %s%-15s%s %s %d/%d HP (Pos: %d,%d)\n",
		couleur, unite.Nom(), ColorReset, barre, hpActuel, hpMax, pos.X(), pos.Y())
}

func (g *GameDemo) afficherInfoUnite(unite *domain.Unite) {
	stats := unite.StatsActuelles()
	pos := unite.Position()

	fmt.Printf("📊 Stats: HP:%d/%d MP:%d/%d ATK:%d DEF:%d MATK:%d ATH:%d%%\n",
		unite.HPActuels(), stats.HP, stats.MP, unite.Stats().MP,
		stats.ATK, stats.DEF, stats.MATK, stats.ATH)
	fmt.Printf("📍 Position: (%d, %d) | Mouvement: %d\n", pos.X(), pos.Y(), stats.MOV)

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
				status = ColorYellow + "MP" + ColorReset
			}
			fmt.Printf("  • %s (MP:%d) %s\n", c.Nom(), c.CoutMP(), status)
		}
	}
}

func (g *GameDemo) afficherAide() {
	fmt.Println()
	fmt.Println(ColorCyan + "🎮 COMMANDES:" + ColorReset)
	fmt.Println("  attack <cible>       - Attaque de base")
	fmt.Println("  skill <nom> <cible>  - Utiliser compétence")
	fmt.Println("  move <x> <y>         - Se déplacer")
	fmt.Println("  pass                 - Passer tour")
	fmt.Println("  help                 - Afficher aide")
	fmt.Println("  quit                 - Quitter")
	fmt.Println()
	fmt.Println(ColorYellow + "💡 ASTUCES:" + ColorReset)
	fmt.Println("  • Les IDs de cibles: gobelin-1, gobelin-2")
	fmt.Println("  • Grille: 0-7 en X et Y")
	fmt.Println("  • Le Mage a la compétence 'fireball'")
	fmt.Println()
}

func (g *GameDemo) afficherFinCombat(resultat string) {
	fmt.Println()
	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════" + ColorReset)

	if resultat == "VICTORY" {
		fmt.Println(ColorBold + ColorGreen + "        🎉 VICTOIRE! 🎉" + ColorReset)
		fmt.Println(ColorGreen + "Les héros ont triomphé des gobelins!" + ColorReset)
	} else if resultat == "DEFEAT" {
		fmt.Println(ColorBold + ColorRed + "        💀 DÉFAITE 💀" + ColorReset)
		fmt.Println(ColorRed + "Les gobelins ont vaincu les héros..." + ColorReset)
	}

	fmt.Println(ColorBold + ColorWhite + "═══════════════════════════════════════" + ColorReset)
	fmt.Println()

	// Résumé
	fmt.Println(ColorCyan + "📊 RÉSUMÉ DU COMBAT:" + ColorReset)
	fmt.Println()
	fmt.Println(ColorGreen + "HÉROS:" + ColorReset)
	for _, u := range g.equipeHeros.Membres() {
		if u.EstEliminee() {
			fmt.Printf("  💀 %s - VAINCU\n", u.Nom())
		} else {
			fmt.Printf("  ⚔️  %s - %d/%d HP\n", u.Nom(), u.HPActuels(), u.Stats().HP)
		}
	}
	fmt.Println()
	fmt.Println(ColorRed + "GOBELINS:" + ColorReset)
	for _, u := range g.equipeEnnemis.Membres() {
		if u.EstEliminee() {
			fmt.Printf("  💀 %s - VAINCU\n", u.Nom())
		} else {
			fmt.Printf("  👹 %s - %d/%d HP\n", u.Nom(), u.HPActuels(), u.Stats().HP)
		}
	}
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
