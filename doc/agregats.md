
# Agrégats Principaux – Aether-Engine

> **Note de synchronisation** :
> Ce fichier centralise la définition des agrégats, Value Objects et entités du domaine. Les diagrammes et la documentation utilisent le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).

> Architecture orientée DDD (Domain Driven Design appliqué au jeu vidéo)
> Chaque agrégat protège ses invariants, garantit la cohérence métier et expose des méthodes métier pures.

---


## 1. Joueur (Player Aggregate)

**Rôle** : Représente un joueur unique, ses états persistants et ses capacités MMO.

**Contenu** :
- Identité du joueur (UUID)
- Statistiques de base : FOR, DEX, INT, VIT…
- Statistiques secondaires : Crit, Blocage, Esquive, Vitesse…
- Ressources : PV, PM, Stamina
- Classe / Spécialisation
- Progression : XP, niveau, talents
- État monde : position, zone, instance
- Effets de statut actifs (buffs / debuffs)
- Références : InventaireID, JournalQuetesID, EquipementID

**Invariants** :
- PV ≥ 0
- Aucun équipement incompatible avec la classe
- Capacité à agir déterminée par le statut
- Les effets actifs expirent correctement

---


## 2. Instance de Combat (CombatInstance Aggregate)

**Rôle** : Gestion complète du cycle de combat tour par tour (PvE, PvP, Raid, Donjon).

**Contenu** :
- ID unique d’une instance de combat
- Liste des participants (Joueurs & PNJ)
- Ordre d’initiative
- Tour en cours
- File des actions à résoudre
- Délais (anti AFK)
- Journal des événements de combat
- Mode : PvE, PvP, Raid, Donjon…

**Invariants** :
- Un seul acteur possède la priorité d’action
- Toute action doit être valide selon l’état actuel
- Aucun participant ne peut agir après sa mort
- Fin du combat : équipe A = 0 PV, équipe B = 0 PV, fuite validée

**Méthodes métier typiques** :
- debutTour(), finTour()
- resoudreCompetence()
- appliquerDegats(), appliquerSoin()
- appliquerEffetStatut()
- prochainActeur()

---


## 3. Inventaire (Inventory Aggregate)

**Rôle** : Gestion de l’inventaire, du stockage, des limites de poids et des conteneurs.

**Contenu** :
- Liste d’items
- Capacité / Poids
- Or / Monnaie
- Slots spéciaux (clé d’instance, consommables…)

**Invariants** :
- Pas d’ajout si l’inventaire dépasse la capacité
- Items stackables : respect des règles de stack
- Pas de suppression d’un item absent

---


## 4. item (Item Aggregate)

**Rôle** : Décrit un objet dans le monde (souvent Entité, mais agrégat pour les items complexes).

**Contenu** :
- ID, type, catégorie
- Bonus de statistiques
- Rareté
- Prérequis (niveau, classe…)
- Effets appliqués (en combat ou hors combat)

**Invariants** :
- Pas de statistiques négatives incohérentes
- Un objet équipé doit respecter le niveau requis

---


## 5. Equipement (EquipmentSet Aggregate)

**Rôle** : Ensemble des items équipés par un joueur.

**Contenu** :
- Tête, Torse, Jambes, Bottes, Arme…
- Résistances cumulées
- Calculs de combat dérivés

**Invariants** :
- 1 seul item par slot
- Slot compatible avec la classe du joueur

---


## 6. Journal de Quêtes / Progression (QuestLog / QuestProgress Aggregate)

**Rôle** : Gestion de l’avancement du joueur dans les quêtes.

**Contenu** :
- Liste des quêtes actives
- État (Non commencée / En cours / Terminée)
- Objectifs : Tuer, Collecter, Explorer, Déclencher événement

**Invariants** :
- Une quête ne peut être complétée que si tous les objectifs sont remplis
- Une quête ne peut être reprise après complétion (sauf mode répétable)

---


## 7. Compétence (Skill / Ability Aggregate)

**Rôle** : Représente les compétences utilisables en combat.

**Contenu** :
- Coût (mana, stamina…)
- Cooldown
- Portée
- Règles de ciblage (soi, allié, zone, ennemi…)

**Invariants** :
- Utilisation impossible hors portée
- Pas d’utilisation en cooldown
- Le coût doit être disponible

---


## 8. PNJ / Monstre (NPC / Monster Aggregate)

**Rôle** : Représente les adversaires dans une instance PvE.

**Contenu** :
- Statistiques
- Table de loot
- Etat IA
- Script de comportement

**Invariants** :
- Aucun monstre ne peut effectuer une action non définie par son IA
- Le monstre disparait lorsque PV <= 0 et loot = résolu

---


## 9. Economie / Echange (Economy / Trade Aggregate)

**Rôle** : Gestion de l’économie persistante.

**Contenu** :
- Ordres de marché
- Items en vente aux enchères
- Prix dynamiques
- Taxes / Frais

**Invariants** :
- Pas de transaction sans fonds suffisants
- Les prix doivent respecter les règles d’économie (min/max)

---


## 10. Etat du Monde (WorldState Aggregate, optionnel)

**Rôle** : Etat persistant du monde partagé.

**Contenu** :
- Evénements dynamiques
- Zones ouvertes / fermées
- Etat temporel (jour/nuit)
- Instances actives

**Invariants** :
- Cohérence entre zones
- Synchronisation entre instances

---


## TL;DR – Synthèse des agrégats clés pour Aether-Engine

- **Joueur** : stats, talents, état du joueur
- **Instance de Combat** : moteur du tour par tour
- **Inventaire** : gestion des items
- **Equipement** : équipement porté
- **Journal de Quêtes** : progression des quêtes
- **Compétence** : capacités utilisées en combat
- **PNJ/Monstre** : adversaires PvE
- **Economie** : transactions & marché
- **Etat du Monde** : état global du monde

> Ces agrégats couvrent tous les invariants et la logique métier essentielle du moteur tactique MMO Aether-Engine.