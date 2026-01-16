package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DemoHandler gère les endpoints de la démo pour le frontend
type DemoHandler struct {
	activeSessions map[string]*DemoSession
	mu             sync.RWMutex
}

// DemoSession représente une session de jeu active
type DemoSession struct {
	ID            string                 `json:"id"`
	Combat        *domain.Combat         `json:"-"`
	EquipeHeros   *domain.Equipe         `json:"equipeHeros"`
	EquipeEnnemis *domain.Equipe         `json:"equipeEnnemis"`
	TourActuel    int                    `json:"tourActuel"`
	Statut        string                 `json:"statut"` // "en_cours", "victoire", "defaite"
	CreatedAt     time.Time              `json:"createdAt"`
	Stats         map[string]*CombatStat `json:"stats"`
}

// CombatStat représente les statistiques d'une unité
type CombatStat struct {
	AttaquesTotal    int `json:"attaquesTotal"`
	AttaquesReussies int `json:"attaquesReussies"`
	AttaquesRatees   int `json:"attaquesRatees"`
	DegatsInfliges   int `json:"degatsInfliges"`
	DegatsRecus      int `json:"degatsRecus"`
	CompetencesUsees int `json:"competencesUsees"`
}

// ActionRequest représente une action du joueur
type ActionRequest struct {
	UniteID      string  `json:"uniteId" binding:"required"`
	TypeAction   string  `json:"typeAction" binding:"required"` // "attaque", "competence", "deplacement"
	CibleID      string  `json:"cibleId,omitempty"`
	CompetenceID string  `json:"competenceId,omitempty"`
	Position     *PosReq `json:"position,omitempty"`
}

type PosReq struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// StartDemoRequest permet de personnaliser la démo au démarrage
type StartDemoRequest struct {
	CustomCharacter *CustomCharacterRequest `json:"customCharacter,omitempty"`
	ReplaceHeroSlot int                     `json:"replaceHeroSlot,omitempty"` // 0=Paladin, 1=Archer, 2=Mage
}

// CustomCharacterRequest représente un personnage personnalisé créé via Angular
type CustomCharacterRequest struct {
	ID          string                     `json:"id"`
	Nom         string                     `json:"nom" binding:"required"`
	Stats       *shared.Stats              `json:"stats" binding:"required"`
	Position    *PosReq                    `json:"position,omitempty"`
	Competences []*CustomCompetenceRequest `json:"competences,omitempty"`
}

// CustomCompetenceRequest représente une compétence personnalisée
type CustomCompetenceRequest struct {
	ID           string  `json:"id" binding:"required"`
	Nom          string  `json:"nom" binding:"required"`
	Description  string  `json:"description"`
	Type         string  `json:"type"` // "attaque", "magie", "soin"
	Portee       int     `json:"portee"`
	CoutMP       int     `json:"coutMP"`
	CoutStamina  int     `json:"coutStamina"`
	Cooldown     int     `json:"cooldown"`
	DegatsBase   int     `json:"degatsBase"`
	Modificateur float64 `json:"modificateur"`
}

// Response structures
type DemoStateResponse struct {
	SessionID     string                 `json:"sessionId"`
	TourActuel    int                    `json:"tourActuel"`
	Statut        string                 `json:"statut"`
	EquipeHeros   *EquipeInfo            `json:"equipeHeros"`
	EquipeEnnemis *EquipeInfo            `json:"equipeEnnemis"`
	Grille        *GrilleInfo            `json:"grille"`
	Message       string                 `json:"message,omitempty"`
	Stats         map[string]*CombatStat `json:"stats"`
}

type EquipeInfo struct {
	ID      string       `json:"id"`
	Nom     string       `json:"nom"`
	Couleur string       `json:"couleur"`
	Unites  []*UniteInfo `json:"unites"`
}

type UniteInfo struct {
	ID          string            `json:"id"`
	Nom         string            `json:"nom"`
	HP          int               `json:"hp"`
	HPMax       int               `json:"hpMax"`
	MP          int               `json:"mp"`
	MPMax       int               `json:"mpMax"`
	Position    *PosReq           `json:"position"`
	Stats       *shared.Stats     `json:"stats"`
	Competences []*CompetenceInfo `json:"competences"`
	EstVivant   bool              `json:"estVivant"`
	AJoue       bool              `json:"aJoue"`
}

type CompetenceInfo struct {
	ID          string `json:"id"`
	Nom         string `json:"nom"`
	Description string `json:"description"`
	Portee      int    `json:"portee"`
	CoutMP      int    `json:"coutMP"`
	CoutStamina int    `json:"coutStamina"`
	Cooldown    int    `json:"cooldown"`
}

type GrilleInfo struct {
	Largeur int `json:"largeur"`
	Hauteur int `json:"hauteur"`
}

// NewDemoHandler crée un nouveau handler de démo
func NewDemoHandler() *DemoHandler {
	return &DemoHandler{
		activeSessions: make(map[string]*DemoSession),
	}
}

// RegisterRoutes enregistre les routes de la démo
func (h *DemoHandler) RegisterRoutes(router *gin.RouterGroup) {
	demo := router.Group("/demo")
	{
		demo.POST("/start", h.StartDemo)
		demo.GET("/:sessionId/state", h.GetState)
		demo.POST("/:sessionId/action", h.ExecuteAction)
		demo.POST("/:sessionId/tour", h.PasserTour)
		demo.DELETE("/:sessionId", h.EndDemo)
	}
}

// StartDemo démarre une nouvelle session de démo
func (h *DemoHandler) StartDemo(c *gin.Context) {
	sessionID := uuid.New().String()

	// Lire la requête optionnelle pour un personnage personnalisé
	var req StartDemoRequest
	_ = c.ShouldBindJSON(&req) // Optionnel, donc on ignore l'erreur

	session, err := h.createDemoSession(sessionID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.mu.Lock()
	h.activeSessions[sessionID] = session
	h.mu.Unlock()

	message := "Combat démarré! À vous de jouer."
	if req.CustomCharacter != nil {
		message = fmt.Sprintf("Combat démarré avec votre personnage %s! À vous de jouer.", req.CustomCharacter.Nom)
	}

	response := h.buildStateResponse(session, message)
	c.JSON(http.StatusCreated, response)
}

// GetState récupère l'état actuel de la démo
func (h *DemoHandler) GetState(c *gin.Context) {
	sessionID := c.Param("sessionId")

	h.mu.RLock()
	session, exists := h.activeSessions[sessionID]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session non trouvée"})
		return
	}

	response := h.buildStateResponse(session, "")
	c.JSON(http.StatusOK, response)
}

// ExecuteAction exécute une action du joueur
func (h *DemoHandler) ExecuteAction(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req ActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.mu.Lock()
	session, exists := h.activeSessions[sessionID]
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session non trouvée"})
		return
	}

	message, err := h.executePlayerAction(session, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si le tour est terminé et jouer l'IA si nécessaire
	if h.allHeroesPlayed(session) {
		aiMessages := h.executeAITurn(session)
		message = message + "\n\n" + aiMessages
		session.TourActuel++
	}

	// Vérifier condition de victoire
	h.checkVictoryConditions(session)

	response := h.buildStateResponse(session, message)
	c.JSON(http.StatusOK, response)
}

// PasserTour passe le tour (toutes les unités jouent automatiquement)
func (h *DemoHandler) PasserTour(c *gin.Context) {
	sessionID := c.Param("sessionId")

	h.mu.Lock()
	session, exists := h.activeSessions[sessionID]
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session non trouvée"})
		return
	}

	messages := h.executeAITurn(session)
	session.TourActuel++
	h.checkVictoryConditions(session)

	response := h.buildStateResponse(session, messages)
	c.JSON(http.StatusOK, response)
}

// EndDemo termine une session de démo
func (h *DemoHandler) EndDemo(c *gin.Context) {
	sessionID := c.Param("sessionId")

	h.mu.Lock()
	delete(h.activeSessions, sessionID)
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Session terminée"})
}

// === Fonctions privées ===

func (h *DemoHandler) createDemoSession(sessionID string, req *StartDemoRequest) (*DemoSession, error) {
	grille, _ := shared.NewGrilleCombat(10, 10)

	// === ÉQUIPE HÉROS ===
	joueurID := "player-1"
	equipeHeros, _ := domain.NewEquipe("team-heros", "Héros de Lumière", "#00FF00", false, &joueurID)

	// Créer les héros par défaut
	heroes := h.createDefaultHeroes()

	// Si un personnage personnalisé est fourni, remplacer le héros spécifié
	if req != nil && req.CustomCharacter != nil {
		slot := req.ReplaceHeroSlot
		if slot < 0 || slot >= len(heroes) {
			slot = 0 // Par défaut, remplacer le Paladin
		}

		customHero, err := h.createCustomHero(req.CustomCharacter)
		if err == nil {
			heroes[slot] = customHero
		}
	}

	// Ajouter les héros à l'équipe
	for _, hero := range heroes {
		equipeHeros.AjouterMembre(hero)
	}

	// === ÉQUIPE ENNEMIS ===
	equipeEnnemis := h.createEnemyTeam()

	// Créer le combat
	combat, err := domain.NewCombat("demo-combat-"+sessionID, []*domain.Equipe{equipeHeros, equipeEnnemis}, grille)
	if err != nil {
		return nil, err
	}

	session := &DemoSession{
		ID:            sessionID,
		Combat:        combat,
		EquipeHeros:   equipeHeros,
		EquipeEnnemis: equipeEnnemis,
		TourActuel:    1,
		Statut:        "en_cours",
		CreatedAt:     time.Now(),
		Stats:         make(map[string]*CombatStat),
	}

	// Initialiser les stats
	for _, unite := range equipeHeros.Membres() {
		session.Stats[string(unite.ID())] = &CombatStat{}
	}
	for _, unite := range equipeEnnemis.Membres() {
		session.Stats[string(unite.ID())] = &CombatStat{}
	}

	return session, nil
}

// createDefaultHeroes crée les 3 héros par défaut
func (h *DemoHandler) createDefaultHeroes() []*domain.Unite {
	heroes := make([]*domain.Unite, 0, 3)

	// Paladin
	statsPaladin := &shared.Stats{HP: 150, MP: 50, Stamina: 100, ATK: 20, DEF: 25, MATK: 8, MDEF: 18, SPD: 8, MOV: 3, ATH: 80}
	posPaladin, _ := shared.NewPosition(1, 4)
	paladin := domain.NewUnite("hero-paladin", "Paladin", "team-heros", statsPaladin, posPaladin)
	taunt := domain.NewCompetence("taunt", "Provocation", "Force un ennemi à vous attaquer", domain.CompetenceAttaque, 3, domain.ZoneEffet{}, 15, 0, 3, 10, 1.2, domain.CibleEnnemis)
	paladin.AjouterCompetence(taunt)
	heroes = append(heroes, paladin)

	// Archer
	statsArcher := &shared.Stats{HP: 90, MP: 40, Stamina: 80, ATK: 28, DEF: 10, MATK: 5, MDEF: 8, SPD: 16, MOV: 4, ATH: 95}
	posArcher, _ := shared.NewPosition(2, 2)
	archer := domain.NewUnite("hero-archer", "Archer", "team-heros", statsArcher, posArcher)
	precisionShot := domain.NewCompetence("precision-shot", "Tir de Précision", "Tir précis avec dégâts augmentés", domain.CompetenceAttaque, 6, domain.ZoneEffet{}, 0, 15, 2, 35, 1.8, domain.CibleEnnemis)
	archer.AjouterCompetence(precisionShot)
	heroes = append(heroes, archer)

	// Mage
	statsMage := &shared.Stats{HP: 70, MP: 120, Stamina: 50, ATK: 8, DEF: 6, MATK: 35, MDEF: 22, SPD: 12, MOV: 3, ATH: 92}
	posMage, _ := shared.NewPosition(2, 6)
	mage := domain.NewUnite("hero-mage", "Mage", "team-heros", statsMage, posMage)
	fireball := domain.NewCompetence("fireball", "Boule de Feu", "Projectile enflammé dévastateur", domain.CompetenceMagie, 5, domain.ZoneEffet{}, 0, 30, 2, 40, 2.0, domain.CibleEnnemis)
	mage.AjouterCompetence(fireball)
	heroes = append(heroes, mage)

	return heroes
}

// createEnemyTeam crée l'équipe ennemie
func (h *DemoHandler) createEnemyTeam() *domain.Equipe {
	equipeEnnemis, _ := domain.NewEquipe("team-ennemis", "Forces des Ténèbres", "#FF0000", true, nil)

	// Orc Berserker
	statsOrc := &shared.Stats{HP: 120, MP: 30, Stamina: 120, ATK: 30, DEF: 18, MATK: 5, MDEF: 10, SPD: 10, MOV: 4, ATH: 85}
	posOrc, _ := shared.NewPosition(8, 4)
	orc := domain.NewUnite("enemy-orc", "Orc Berserker", "team-ennemis", statsOrc, posOrc)
	rage := domain.NewCompetence("rage", "Rage", "Attaque puissante qui ignore une partie de la défense", domain.CompetenceAttaque, 2, domain.ZoneEffet{}, 25, 0, 3, 15, 2.2, domain.CibleEnnemis)
	orc.AjouterCompetence(rage)
	equipeEnnemis.AjouterMembre(orc)

	// Gobelin Archer
	statsGobelin := &shared.Stats{HP: 60, MP: 20, Stamina: 60, ATK: 22, DEF: 8, MATK: 3, MDEF: 6, SPD: 18, MOV: 5, ATH: 88}
	posGobelin, _ := shared.NewPosition(7, 2)
	gobelin := domain.NewUnite("enemy-gobelin", "Gobelin Archer", "team-ennemis", statsGobelin, posGobelin)
	poisonArrow := domain.NewCompetence("poison-arrow", "Flèche Empoisonnée", "Flèche qui inflige des dégâts sur la durée", domain.CompetenceAttaque, 5, domain.ZoneEffet{}, 0, 10, 2, 20, 1.3, domain.CibleEnnemis)
	gobelin.AjouterCompetence(poisonArrow)
	equipeEnnemis.AjouterMembre(gobelin)

	// Nécromancien
	statsNecro := &shared.Stats{HP: 80, MP: 100, Stamina: 40, ATK: 10, DEF: 8, MATK: 32, MDEF: 20, SPD: 11, MOV: 3, ATH: 90}
	posNecro, _ := shared.NewPosition(7, 6)
	necro := domain.NewUnite("enemy-necro", "Nécromancien", "team-ennemis", statsNecro, posNecro)
	darkBolt := domain.NewCompetence("dark-bolt", "Éclair Noir", "Magie noire destructrice", domain.CompetenceMagie, 5, domain.ZoneEffet{}, 0, 25, 2, 35, 1.9, domain.CibleEnnemis)
	necro.AjouterCompetence(darkBolt)
	equipeEnnemis.AjouterMembre(necro)

	return equipeEnnemis
}

// createCustomHero crée un héros personnalisé à partir des données Angular
func (h *DemoHandler) createCustomHero(customChar *CustomCharacterRequest) (*domain.Unite, error) {
	// Position par défaut ou fournie
	var pos *shared.Position
	if customChar.Position != nil {
		pos, _ = shared.NewPosition(customChar.Position.X, customChar.Position.Y)
	} else {
		pos, _ = shared.NewPosition(1, 4) // Position par défaut
	}

	// ID par défaut si non fourni
	uniteID := customChar.ID
	if uniteID == "" {
		uniteID = "custom-hero-" + uuid.New().String()[:8]
	}

	// Créer l'unité
	unite := domain.NewUnite(domain.UnitID(uniteID), customChar.Nom, "team-heros", customChar.Stats, pos)

	// Ajouter les compétences personnalisées
	for _, compReq := range customChar.Competences {
		typeComp := h.parseCompetenceType(compReq.Type)
		competence := domain.NewCompetence(
			shared.CompetenceID(compReq.ID),
			compReq.Nom,
			compReq.Description,
			typeComp,
			compReq.Portee,
			domain.ZoneEffet{},
			compReq.CoutMP,
			compReq.CoutStamina,
			compReq.Cooldown,
			compReq.DegatsBase,
			compReq.Modificateur,
			domain.CibleEnnemis,
		)
		unite.AjouterCompetence(competence)
	}

	return unite, nil
}

// parseCompetenceType convertit le type de compétence string en TypeCompetence
func (h *DemoHandler) parseCompetenceType(typeStr string) domain.TypeCompetence {
	switch typeStr {
	case "magie":
		return domain.CompetenceMagie
	case "soin":
		return domain.CompetenceSoin
	case "buff":
		return domain.CompetenceBuff
	case "debuff":
		return domain.CompetenceDebuff
	default:
		return domain.CompetenceAttaque
	}
}

func (h *DemoHandler) buildStateResponse(session *DemoSession, message string) *DemoStateResponse {
	return &DemoStateResponse{
		SessionID:     session.ID,
		TourActuel:    session.TourActuel,
		Statut:        session.Statut,
		EquipeHeros:   h.buildEquipeInfo(session.EquipeHeros),
		EquipeEnnemis: h.buildEquipeInfo(session.EquipeEnnemis),
		Grille: &GrilleInfo{
			Largeur: 10,
			Hauteur: 10,
		},
		Message: message,
		Stats:   session.Stats,
	}
}

func (h *DemoHandler) buildEquipeInfo(equipe *domain.Equipe) *EquipeInfo {
	unites := make([]*UniteInfo, 0)
	for _, unite := range equipe.Membres() {
		competences := make([]*CompetenceInfo, 0)
		for _, comp := range unite.Competences() {
			competences = append(competences, &CompetenceInfo{
				ID:          string(comp.ID()),
				Nom:         comp.Nom(),
				Description: comp.Description(),
				Portee:      comp.Portee(),
				CoutMP:      comp.CoutMP(),
				CoutStamina: comp.CoutStamina(),
				Cooldown:    comp.Cooldown(),
			})
		}

		pos := unite.Position()
		unites = append(unites, &UniteInfo{
			ID:          string(unite.ID()),
			Nom:         unite.Nom(),
			HP:          unite.Stats().HP,
			HPMax:       unite.Stats().HP,
			MP:          unite.Stats().MP,
			MPMax:       unite.Stats().MP,
			Position:    &PosReq{X: pos.X(), Y: pos.Y()},
			Stats:       unite.Stats(),
			Competences: competences,
			EstVivant:   !unite.EstEliminee(),
			AJoue:       false, // TODO: tracking
		})
	}

	return &EquipeInfo{
		ID:      string(equipe.ID()),
		Nom:     equipe.Nom(),
		Couleur: equipe.Couleur(),
		Unites:  unites,
	}
}

func (h *DemoHandler) executePlayerAction(session *DemoSession, req *ActionRequest) (string, error) {
	unite := session.EquipeHeros.ObtenirMembre(domain.UnitID(req.UniteID))
	if unite == nil {
		return "", fmt.Errorf("unité non trouvée")
	}

	if unite.EstEliminee() {
		return "", fmt.Errorf("unité morte")
	}

	switch req.TypeAction {
	case "attaque":
		return h.executeAttack(session, unite, req.CibleID)
	case "competence":
		return h.executeSkill(session, unite, req.CompetenceID, req.CibleID)
	case "deplacement":
		if req.Position == nil {
			return "", fmt.Errorf("position requise pour déplacement")
		}
		newPos, _ := shared.NewPosition(req.Position.X, req.Position.Y)
		// Note: La méthode Deplacer du Combat n'est peut-être pas disponible
		// Pour simplifier, on déplace directement l'unité
		unite.DeplacerVers(newPos)
		return fmt.Sprintf("%s se déplace vers (%d, %d)", unite.Nom(), req.Position.X, req.Position.Y), nil
	default:
		return "", fmt.Errorf("type d'action invalide")
	}
}

func (h *DemoHandler) executeAttack(session *DemoSession, attaquant *domain.Unite, cibleID string) (string, error) {
	cible := session.EquipeEnnemis.ObtenirMembre(domain.UnitID(cibleID))
	if cible == nil {
		return "", fmt.Errorf("cible non trouvée")
	}

	stats := session.Stats[string(attaquant.ID())]
	stats.AttaquesTotal++

	// Calcul du test d'ATH
	ath := attaquant.Stats().ATH
	roll := rand.Intn(100) + 1

	message := fmt.Sprintf("🎯 %s attaque %s (ATH: %d, Roll: %d)\n", attaquant.Nom(), cible.Nom(), ath, roll)

	if roll <= ath {
		// Touché!
		degats := attaquant.Stats().ATK - cible.Stats().DEF/2
		if degats < 1 {
			degats = 1
		}

		cible.RecevoirDegats(degats)
		stats.AttaquesReussies++
		stats.DegatsInfliges += degats
		session.Stats[cibleID].DegatsRecus += degats

		message += fmt.Sprintf("✅ Touché! %d dégâts infligés\n", degats)
		if cible.EstEliminee() {
			message += fmt.Sprintf("💀 %s est vaincu!\n", cible.Nom())
		} else {
			message += fmt.Sprintf("   HP restants: %d\n", cible.HPActuels())
		}
	} else {
		// Raté!
		stats.AttaquesRatees++
		message += "❌ Attaque ratée!\n"
	}

	return message, nil
}

func (h *DemoHandler) executeSkill(session *DemoSession, attaquant *domain.Unite, competenceID, cibleID string) (string, error) {
	competence := attaquant.ObtenirCompetence(shared.CompetenceID(competenceID))
	if competence == nil {
		return "", fmt.Errorf("compétence non trouvée")
	}

	cible := session.EquipeEnnemis.ObtenirMembre(domain.UnitID(cibleID))
	if cible == nil {
		return "", fmt.Errorf("cible non trouvée")
	}

	// Vérifier les coûts
	if attaquant.Stats().MP < competence.CoutMP() {
		return "", fmt.Errorf("pas assez de MP")
	}
	if attaquant.Stats().Stamina < competence.CoutStamina() {
		return "", fmt.Errorf("pas assez de stamina")
	}

	attaquant.ConsommerMP(competence.CoutMP())
	attaquant.ConsommerStamina(competence.CoutStamina())

	stats := session.Stats[string(attaquant.ID())]
	stats.CompetencesUsees++

	// Calcul des dégâts
	baseDamage := competence.DegatsBase()
	if competence.Type() == domain.CompetenceMagie {
		baseDamage = int(float64(attaquant.Stats().MATK) * competence.Modificateur())
	} else {
		baseDamage = int(float64(attaquant.Stats().ATK) * competence.Modificateur())
	}

	degats := baseDamage - cible.Stats().DEF/3
	if degats < 1 {
		degats = 1
	}

	cible.RecevoirDegats(degats)
	stats.DegatsInfliges += degats
	session.Stats[cibleID].DegatsRecus += degats

	message := fmt.Sprintf("✨ %s utilise %s sur %s!\n", attaquant.Nom(), competence.Nom(), cible.Nom())
	message += fmt.Sprintf("   💥 %d dégâts infligés\n", degats)

	if cible.EstEliminee() {
		message += fmt.Sprintf("   💀 %s est vaincu!\n", cible.Nom())
	} else {
		message += fmt.Sprintf("   HP restants: %d\n", cible.HPActuels())
	}

	return message, nil
}

func (h *DemoHandler) allHeroesPlayed(session *DemoSession) bool {
	// Pour simplifier, on considère qu'un héros a joué après une action
	// Dans une vraie implémentation, il faudrait tracker l'état de chaque unité
	return false // Le frontend contrôle quand passer au tour suivant
}

func (h *DemoHandler) executeAITurn(session *DemoSession) string {
	messages := "🤖 Tour de l'IA:\n\n"

	for _, unite := range session.EquipeEnnemis.Membres() {
		if unite.EstEliminee() {
			continue
		}

		// IA simple: attaquer la cible la plus faible
		var cibleChoisie *domain.Unite
		minHP := 9999
		for _, hero := range session.EquipeHeros.Membres() {
			if !hero.EstEliminee() && hero.HPActuels() < minHP {
				minHP = hero.HPActuels()
				cibleChoisie = hero
			}
		}

		if cibleChoisie != nil {
			// 50% de chance d'utiliser une compétence si disponible
			if len(unite.Competences()) > 0 && rand.Intn(100) < 50 {
				comp := unite.Competences()[0]
				if unite.Stats().MP >= comp.CoutMP() {
					msg, _ := h.executeSkill(session, unite, string(comp.ID()), string(cibleChoisie.ID()))
					messages += msg + "\n"
					continue
				}
			}

			// Sinon attaque normale
			msg, _ := h.executeAttack(session, unite, string(cibleChoisie.ID()))
			messages += msg + "\n"
		}
	}

	return messages
}

func (h *DemoHandler) checkVictoryConditions(session *DemoSession) {
	herosVivants := 0
	ennemisVivants := 0

	for _, unite := range session.EquipeHeros.Membres() {
		if !unite.EstEliminee() {
			herosVivants++
		}
	}

	for _, unite := range session.EquipeEnnemis.Membres() {
		if !unite.EstEliminee() {
			ennemisVivants++
		}
	}

	if herosVivants == 0 {
		session.Statut = "defaite"
	} else if ennemisVivants == 0 {
		session.Statut = "victoire"
	}
}
