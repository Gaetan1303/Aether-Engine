# Agrégats Principaux – Aether-Engine

> Architecture orientée DDD (Domain Driven Design for Games)
>
> Chaque agrégat protège ses invariants, garantit la cohérence métier et expose des méthodes métier pures.

---

## 1. Player Aggregate (Joueur)

**Rôle** : Représente un joueur unique, ses états persistants et ses capacités MMO.

**Contenu** :
- Identité du joueur (UUID)
- Stats de base : STR, DEX, INT, VIT…
- Stats secondaires : Crit, Block, Dodge, Speed…
- Ressources : HP, MP, Stamina
- Classe / Spécialisation
- Progression : XP, niveau, talents
- État monde : position, zone, instance
- Status Effects actifs (buffs / debuffs)
- Références : InventoryID, QuestLogID, EquipmentSetID

**Invariants** :
- HP ≥ 0
- Aucun équipement incompatible avec la classe
- Capacité à agir déterminée par le statut
- Les effets actifs expirent correctement

---

## 2. CombatInstance Aggregate (Combat Tour par Tour)

**Rôle** : Gestion complète du cycle de combat tour par tour (PvE, PvP, Raid, Dungeon).

**Contenu** :
- ID unique d’une instance de combat
- Liste des participants (Players & NPCs)
- Initiative order
- Tour en cours
- File des actions à résoudre (action queue)
- Timeouts / Délais (anti AFK)
- Log des événements de combat
- Mode : PvE, PvP, Raid, Dungeon…

**Invariants** :
- Un seul acteur possède la priorité d’action
- Toute action doit être valide selon l’état actuel
- Aucun participant ne peut agir après sa mort
- Fin du combat : équipe A = 0 HP, équipe B = 0 HP, fuite validée

**Méthodes métier typiques** :
- StartTurn(), EndTurn()
- ResolveSkill()
- ApplyDamage(), ApplyHeal()
- ApplyStatusEffect()
- NextActor()

---

## 3. Inventory Aggregate

**Rôle** : Gestion de l’inventaire, du stockage, des limites de poids et des conteneurs.

**Contenu** :
- Liste d’items
- Capacité / Poids
- Gold / Currency
- Slots spéciaux (clé d’instance, consommables…)

**Invariants** :
- Pas d’ajout si l’inventaire dépasse la capacité
- Items stackables : respect des règles de stack
- Pas de suppression d’un item absent

---

## 4. Item Aggregate

**Rôle** : Décrit un objet dans le monde (souvent Entity, mais agrégat pour les items complexes).

**Contenu** :
- ID, type, catégorie
- Stats bonus
- Rareté
- Requirements (niveau, classe…)
- Effets appliqués (en combat ou hors combat)

**Invariants** :
- Pas de stats négatives incohérentes
- Un objet équipé doit respecter le niveau requis

---

## 5. EquipmentSet Aggregate

**Rôle** : Ensemble des objets équipés par un joueur.

**Contenu** :
- Head, Chest, Legs, Boots, Weapon…
- Résistances cumulées
- Calculs de combat dérivés

**Invariants** :
- 1 seul objet par slot
- Slot compatible avec la classe du joueur

---

## 6. QuestLog / QuestProgress Aggregate

**Rôle** : Gestion de l’avancement du joueur dans les quêtes.

**Contenu** :
- Liste des quêtes actives
- État (Not started / In progress / Completed)
- Objectifs : Kill, Collect, Explore, Trigger event

**Invariants** :
- Une quête ne peut être complétée que si tous les objectifs sont remplis
- Une quête ne peut être reprise après completion (sauf mode repeatable)

---

## 7. Skill / Ability Aggregate

**Rôle** : Représente les compétences utilisables en combat.

**Contenu** :
- Coût (mana, stamina…)
- Cooldown
- Portée
- Règles de ciblage (self, ally, area, enemy…)

**Invariants** :
- Utilisation impossible hors range
- Pas d’utilisation en cooldown
- Le coût doit être disponible

---

## 8. NPC / Monster Aggregate

**Rôle** : Représente les adversaires dans une instance PvE.

**Contenu** :
- Stats
- LootTable
- AIState
- Behavior Script

**Invariants** :
- Aucun monstre ne peut effectuer une action non définie par son AI
- Le monstre disparait lorsque HP <= 0 et drop = résolu

---

## 9. Economy / Trade Aggregate

**Rôle** : Gestion de l’économie persistante.

**Contenu** :
- Market Orders
- Auction Items
- Prix dynamiques
- Taxes / Fees

**Invariants** :
- Pas de transaction sans fonds suffisants
- Les prix doivent respecter les règles d’économie (min/max)

---

## 10. WorldState Aggregate (optionnel)

**Rôle** : État persistant du monde partagé.

**Contenu** :
- Événements dynamiques
- Zones ouvertes / fermées
- Time state (jour/nuit)
- Instances actives

**Invariants** :
- Cohérence entre zones
- Synchronisation entre instances

---

## TL;DR – Synthèse des agrégats clés pour Aether-Engine

- **Player** : stats, talents, état du joueur
- **CombatInstance** : moteur du tour par tour
- **Inventory** : gestion des items
- **EquipmentSet** : équipement porté
- **QuestLog** : progression des quêtes
- **Skill** : capacités utilisées en combat
- **NPC/Monster** : adversaires PvE
- **Economy** : transactions & marché
- **WorldState** : état global du monde

> Ces agrégats couvrent tous les invariants et la logique métier essentielle du moteur tactique MMO Aether-Engine.