#  Plan de Simplification - Pivot Pragmatique

**Date :** 5 décembre 2025  
**Objectif :** Transformer Aether Engine d'un projet en démo jouable

---

##  Ce Qui Est Déjà Fait

### 1. Démo CLI Fonctionnelle (`cmd/demo/main.go`)

**Résultat immédiat :** Combat 2v2 jouable en 600 lignes
-  Interface texte interactive
-  Combat tour par tour fonctionnel
-  IA basique mais effective
-  Compétences avec MP et cooldowns
-  Système de dégâts complet
-  Conditions de victoire/défaite

**Temps de développement :** 30 minutes  
**Utilise le code existant :** 100% (zéro refactoring du domaine)

**Commande pour tester :**
```bash
go build -o bin/demo cmd/demo/main.go
./bin/demo
```

---

##  Simplifications à Appliquer (Prochaines Étapes)

### Phase 1 : Supprimer l'Infrastructure Complexe (Gain 40% complexité)

#### A Supprimer Complètement

1. **Event Sourcing distribué**
   - ❌ `internal/combat/infrastructure/event_store_postgres.go` (248 lignes)
   - ❌ `pkg/eventbus/kafka_publisher.go` (123 lignes)
   - ❌ Toute la stack PostgreSQL/Kafka/Redis
   - ✅ **Garder** : Events in-memory pour undo/replay local

2. **Patterns inutilisés** (7 sur 11)
   - ❌ Observer Pattern (pas de subscribers)
   - ❌ Chain of Responsibility (validation simple suffit)
   - ❌ Factory Pattern pour DamageCalculator (injection directe)
   - ❌ Command Pattern complexe (fonctions simples suffisent)
   - ✅ **Garder** : Strategy (DamageCalculator), State Machine (états combat), Composition (Unite)

3. **Abstractions excessives**
   - ❌ `StateMachineProvider` interface
   - ❌ `CommandSystemProvider` interface
   - ❌ `ObserverProvider` interface
   - ❌ `ValidationProvider` interface
   - ✅ **Garder** : Interfaces utiles (DamageCalculator, EventStore simplifié)

**Fichiers à supprimer :**
```bash
rm -rf internal/combat/infrastructure/
rm -rf internal/combat/combatfacade/
rm -rf internal/combat/combatinitializer/
rm -rf pkg/eventbus/
rm -rf internal/combat/domain/states/
```

---

### Phase 2 : Simplifier le Domaine

#### Refactorings Ciblés

1. **Combat.go** (382 lignes → ~200 lignes)
   ```go
   // AVANT (11 dépendances)
   type Combat struct {
       stateMachine    StateMachineProvider
       commandInvoker  CommandSystemProvider
       commandFactory  CommandFactoryProvider
       observerSubject ObserverProvider
       validationChain ValidationProvider
       // ... 6 autres champs
   }
   
   // APRÈS (4 dépendances)
   type Combat struct {
       id              string
       etat            EtatCombat
       equipes         map[TeamID]*Equipe
       grille          *GrilleCombat
       tourActuel      int
       damageCalculator DamageCalculator
   }
   ```

2. **Supprimer les TODOs non critiques**
   -  `DistribuerRecompenses()` (système de loot)
   -  `PossedeObjet()`, `ConsommerObjet()` (inventaire)
   -  **Garder** : Mécaniques de combat core

3. **Simplifier Unite.go** (485 lignes → ~300 lignes)
   - Garder la composition (bien fait)
   - Supprimer les hooks pour patterns externes

---

### Phase 3 : Implémenter les Manques Critiques

#### Statuts Fonctionnels (Actuellement NON implémentés)

**Fichier :** `internal/combat/domain/unit_status_manager.go`

```go
// À IMPLÉMENTER
func (m *UnitStatusManager) EstStun() bool {
    // Actuellement retourne toujours false
    // TODO: Vérifier si statut Stun actif
}

func (m *UnitStatusManager) EstRoot() bool {
    // Actuellement retourne toujours false
    // TODO: Vérifier si statut Root actif
}

func (m *UnitStatusManager) AppliquerEffetsPoison() {
    // Actuellement NE FAIT RIEN
    // TODO: Infliger dégâts par tour
}
```

**Priorité :** HAUTE (rend le combat plus intéressant)  
**Temps estimé :** 2-3 heures


---

##  Nouvelle Roadmap Pragmatique

### Semaine 1-2 : Démo Solide
- [x] Démo CLI jouable (FAIT)
- [ ] Implémenter statuts (Poison, Stun, Root)
- [ ] Ajouter 3-4 compétences variées
- [ ] Améliorer l'IA (pathfinding basique)

### Semaine 3-4 : Polish & Documentation
- [ ] Tests d'acceptation pour la démo
- [ ] README attractif avec GIFs
- [ ] Architecture simplifiée documentée
- [ ] Benchmark de performance

### Semaine 5-6 : UI Graphique (Optionnel)
- [ ] Port vers Raylib ou Ebiten
- [ ] Sprites et animations simples
- [ ] Menu principal
- [ ] Sauvegarde/Chargement

### Au-delà : Features Gameplay
- [ ] Système de progression (XP, levels)
- [ ] 5-6 types d'unités
- [ ] 3-4 cartes différentes
- [ ] Mode campagne (3-5 combats)

---

##  Leçons Apprises

### Ce Qui Était Bon
 Architecture DDD propre (séparation domaine/infra)  
 Tests unitaires solides (190 tests)  
 Patterns bien implémentés techniquement  
 Code lisible et maintenable  

### Ce Qui Était Problématique
 Over-engineering sans besoin business  
 Abstractions prématurées (YAGNI violation)  
 Focus sur l'architecture, pas sur le jeu  
 Aucune validation par un joueur réel  

### Nouveau Mindset
 **Démo d'abord, architecture ensuite**  
 **Jouabilité > Pureté architecturale**  
 **Patterns quand douleur, pas par anticipation**  
 **Mesure de succès = "Quelqu'un peut jouer"**  

---

##  Pitch

**Avant simplification :**
> "Serveur de combat tactique avec Event Sourcing, 11 design patterns, architecture hexagonale..."
> 
> **Réaction recruteur :**  "Over-engineered, pas de démo"

**Après simplification :**
> "Combat tactique tour par tour jouable en CLI. Démo en 30 secondes :
> ```bash
> ./bin/demo
> > attack gobelin-1
>   Guerrier attaque Gobelin et inflige 15 dégâts!
> ```
> Architecture DDD propre, 190 tests, code production-ready."
>
> **Réaction recruteur :**  "Pragmatique, ça marche, je peux essayer"

---

##  Réponse

**Question :** "Le projet est-il condamné ?"

**Réponse :** Non, il est **sauvable** car :
1.  Le domaine fonctionne (preuve : démo en 600 lignes)
2.  Les tests sont solides (190 tests passent)
3.  L'architecture est propre (séparation concerns)
4.  Pivot rapide possible (suppression = 1 jour)

**Pronostic révisé :**
- **Avant :** Projet enlisé dans refactorings infinis
- **Après :** Démo jouable → Portfolio employable → Features itératives

**Engagement :**
-  **Semaine 1-2 :** Démo polish + statuts
-  **Semaine 3-4 :** UI graphique basique
-  **Semaine 5-6 :** 3-5 combats jouables (mini-campagne)

---

**Signature :** El Miminette For Ever !!
**Date :** 5 décembre 2025
