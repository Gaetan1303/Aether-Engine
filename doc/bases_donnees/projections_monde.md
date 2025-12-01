# Projections Monde - Modèles de Lecture

## Vue d'ensemble

Les projections monde regroupent les modèles de lecture pour les éléments du monde du jeu : items, compétences, quêtes, économie et état du monde. Ces projections sont construites à partir des événements des agrégats correspondants.

## Architecture

- **Base de données**: PostgreSQL
- **Pattern**: CQRS - Modèles de lecture optimisés pour les requêtes
- **Agrégats Sources**: Item, Competence, Quete, Economy, WorldState
- **Handlers**: ItemProjection, SkillProjection, QuestProjection, EconomyProjection, WorldProjection
- **Normalisation**: Dénormalisation pour performance des lectures

## Schéma des Projections

```mermaid
erDiagram
    ITEMS {
        UUID item_id PK "Identifiant unique de l'item"
        VARCHAR nom "Nom de l'item"
        VARCHAR description "Description"
        VARCHAR type_item "ARME, ARMURE, CONSOMMABLE, MATERIAU, QUETE"
        VARCHAR rarete "COMMUN, RARE, EPIQUE, LEGENDAIRE"
        INTEGER niveau_requis "Niveau minimum pour utiliser"
        JSONB stats "Statistiques de l'item"
        JSONB effets "Effets spéciaux"
        INTEGER valeur_or "Valeur marchande en or"
        INTEGER poids "Poids en unités"
        BOOLEAN empilable "Peut être empilé"
        INTEGER stack_max "Taille max de la pile"
        BOOLEAN echangeable "Peut être échangé"
        BOOLEAN destructible "Peut être détruit"
        VARCHAR icone_url "URL de l'icône"
        JSONB recette_craft "Recette de fabrication (si applicable)"
        BIGINT last_event_sequence "Dernière séquence traitée"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    COMPETENCES {
        UUID competence_id PK "Identifiant unique de la compétence"
        VARCHAR nom "Nom de la compétence"
        VARCHAR description "Description"
        VARCHAR type_competence "ACTIVE, PASSIVE, ULTIMATE"
        VARCHAR element "FEU, EAU, TERRE, AIR, LUMIERE, OMBRE, NEUTRE"
        VARCHAR cible "SOI, ALLIE, ENNEMI, ZONE"
        INTEGER niveau_requis "Niveau requis"
        INTEGER cout_mana "Coût en mana"
        INTEGER cout_stamina "Coût en endurance"
        INTEGER cooldown_tours "Temps de recharge en tours"
        INTEGER portee "Portée en cases"
        JSONB effets "Effets de la compétence"
        JSONB modificateurs_niveau "Scaling par niveau"
        VARCHAR animation "Identifiant de l'animation"
        VARCHAR icone_url "URL de l'icône"
        JSONB pre_requis "Prérequis (autres compétences, classe, etc.)"
        BIGINT last_event_sequence "Dernière séquence traitée"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    QUETES {
        UUID quete_id PK "Identifiant unique de la quête"
        VARCHAR nom "Nom de la quête"
        VARCHAR description "Description"
        VARCHAR type_quete "PRINCIPALE, SECONDAIRE, QUOTIDIENNE, EVENEMENT"
        INTEGER niveau_requis "Niveau minimum"
        INTEGER niveau_recommande "Niveau recommandé"
        JSONB pre_requis "Prérequis (quêtes, niveau, items, etc.)"
        JSONB objectifs "Liste des objectifs"
        JSONB recompenses "Récompenses (XP, or, items, etc.)"
        VARCHAR zone_depart "Zone de départ"
        VARCHAR pnj_donneur "PNJ qui donne la quête"
        INTEGER duree_limite_minutes "Durée limite (si applicable)"
        BOOLEAN repete "Peut être répétée"
        VARCHAR frequence_repetition "QUOTIDIEN, HEBDOMADAIRE"
        JSONB dialogues "Dialogues associés"
        BIGINT last_event_sequence "Dernière séquence traitée"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    ETAT_MONDE {
        UUID monde_id PK "Identifiant du monde (singleton)"
        INTEGER cycle_jour_nuit "Position dans le cycle (0-23)"
        VARCHAR saison "PRINTEMPS, ETE, AUTOMNE, HIVER"
        INTEGER jour_saison "Jour dans la saison (1-30)"
        JSONB evenements_actifs "Événements mondiaux en cours"
        JSONB boss_vaincus "Boss mondiaux vaincus"
        JSONB zones_decouvertes "Zones découvertes par la communauté"
        JSONB ressources_mondiales "État des ressources globales"
        JSONB facteurs_economiques "Facteurs influençant l'économie"
        JSONB classements "Leaderboards mondiaux"
        TIMESTAMP derniere_mise_a_jour_cycle "Dernière update du cycle"
        BIGINT last_event_sequence "Dernière séquence traitée"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    ORDRES_ECONOMIE {
        UUID ordre_id PK "Identifiant unique de l'ordre"
        VARCHAR type_ordre "VENTE, ACHAT"
        UUID joueur_id "ID du joueur"
        VARCHAR pseudo_joueur "Pseudo du joueur"
        UUID item_id FK "Référence vers l'item"
        INTEGER quantite "Quantité totale"
        INTEGER quantite_restante "Quantité non encore échangée"
        INTEGER prix_unitaire "Prix par unité en or"
        VARCHAR statut "ACTIF, PARTIELLEMENT_EXECUTE, COMPLETE, ANNULE, EXPIRE"
        TIMESTAMP cree_a "Date de création"
        TIMESTAMP expire_a "Date d'expiration"
        TIMESTAMP complete_a "Date de complétion"
        BIGINT last_event_sequence "Dernière séquence traitée"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    TRANSACTIONS_ECONOMIE {
        UUID transaction_id PK "Identifiant unique de la transaction"
        UUID ordre_vente_id FK "Ordre de vente"
        UUID ordre_achat_id FK "Ordre d'achat"
        UUID vendeur_id "ID du vendeur"
        VARCHAR pseudo_vendeur "Pseudo du vendeur"
        UUID acheteur_id "ID de l'acheteur"
        VARCHAR pseudo_acheteur "Pseudo de l'acheteur"
        UUID item_id FK "Item échangé"
        INTEGER quantite "Quantité échangée"
        INTEGER prix_unitaire "Prix unitaire de la transaction"
        INTEGER montant_total "Montant total de la transaction"
        INTEGER taxe "Taxe prélevée"
        TIMESTAMP executee_a "Date d'exécution"
        BIGINT event_sequence "Séquence de l'événement source"
    }

    PRIX_MARCHE {
        UUID prix_id PK "Identifiant unique"
        UUID item_id FK "Référence vers l'item"
        INTEGER prix_moyen_24h "Prix moyen sur 24h"
        INTEGER prix_min_24h "Prix minimum sur 24h"
        INTEGER prix_max_24h "Prix maximum sur 24h"
        INTEGER volume_24h "Volume échangé sur 24h"
        INTEGER nombre_transactions_24h "Nombre de transactions"
        JSONB historique_prix "Historique des prix (points horaires)"
        TIMESTAMP derniere_transaction "Date de la dernière transaction"
        TIMESTAMP updated_at "Dernière mise à jour"
    }

    ITEMS ||--o{ ORDRES_ECONOMIE : "fait l'objet de"
    ITEMS ||--o{ TRANSACTIONS_ECONOMIE : "est échangé dans"
    ITEMS ||--o{ PRIX_MARCHE : "a un prix"
    ORDRES_ECONOMIE ||--o{ TRANSACTIONS_ECONOMIE : "génère"
```

## Tables Détaillées

### ITEMS

Catalogue complet des items du jeu.

```sql
CREATE TABLE items (
    item_id UUID PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    description TEXT,
    type_item VARCHAR(20) NOT NULL CHECK (type_item IN ('ARME', 'ARMURE', 'CONSOMMABLE', 'MATERIAU', 'QUETE', 'AUTRE')),
    rarete VARCHAR(15) NOT NULL CHECK (rarete IN ('COMMUN', 'RARE', 'EPIQUE', 'LEGENDAIRE', 'MYTHIQUE')),
    niveau_requis INTEGER DEFAULT 1,
    stats JSONB NOT NULL DEFAULT '{}'::jsonb,
    effets JSONB NOT NULL DEFAULT '[]'::jsonb,
    valeur_or INTEGER NOT NULL DEFAULT 0,
    poids INTEGER NOT NULL DEFAULT 1,
    empilable BOOLEAN DEFAULT false,
    stack_max INTEGER DEFAULT 1,
    echangeable BOOLEAN DEFAULT true,
    destructible BOOLEAN DEFAULT true,
    icone_url VARCHAR(255),
    recette_craft JSONB,
    last_event_sequence BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_items_type ON items(type_item);
CREATE INDEX idx_items_rarete ON items(rarete);
CREATE INDEX idx_items_niveau ON items(niveau_requis);
CREATE INDEX idx_items_nom ON items(nom);
CREATE INDEX idx_items_echangeable ON items(echangeable) WHERE echangeable = true;
CREATE INDEX idx_items_stats ON items USING GIN(stats);
```

**Événements Sources**:
- ItemCree → Création de l'item
- ItemModifie → Mise à jour des propriétés
- ProprietesModifiees → Mise à jour stats/effets

### COMPETENCES

Catalogue des compétences disponibles.

```sql
CREATE TABLE competences (
    competence_id UUID PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    description TEXT,
    type_competence VARCHAR(10) NOT NULL CHECK (type_competence IN ('ACTIVE', 'PASSIVE', 'ULTIMATE')),
    element VARCHAR(10) CHECK (element IN ('FEU', 'EAU', 'TERRE', 'AIR', 'LUMIERE', 'OMBRE', 'NEUTRE')),
    cible VARCHAR(10) NOT NULL CHECK (cible IN ('SOI', 'ALLIE', 'ENNEMI', 'ZONE', 'TOUS')),
    niveau_requis INTEGER DEFAULT 1,
    cout_mana INTEGER DEFAULT 0,
    cout_stamina INTEGER DEFAULT 0,
    cooldown_tours INTEGER DEFAULT 0,
    portee INTEGER DEFAULT 1,
    effets JSONB NOT NULL DEFAULT '[]'::jsonb,
    modificateurs_niveau JSONB DEFAULT '{}'::jsonb,
    animation VARCHAR(50),
    icone_url VARCHAR(255),
    pre_requis JSONB DEFAULT '{}'::jsonb,
    last_event_sequence BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_competences_type ON competences(type_competence);
CREATE INDEX idx_competences_element ON competences(element);
CREATE INDEX idx_competences_niveau ON competences(niveau_requis);
CREATE INDEX idx_competences_nom ON competences(nom);
```

**Événements Sources**:
- CompetenceApprise (au niveau agrégat Joueur, mais création dans catalogue)
- CompetenceAmelioree → Mise à jour modificateurs_niveau

### QUETES

Catalogue des quêtes disponibles.

```sql
CREATE TABLE quetes (
    quete_id UUID PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    description TEXT,
    type_quete VARCHAR(15) NOT NULL CHECK (type_quete IN ('PRINCIPALE', 'SECONDAIRE', 'QUOTIDIENNE', 'EVENEMENT')),
    niveau_requis INTEGER DEFAULT 1,
    niveau_recommande INTEGER,
    pre_requis JSONB DEFAULT '{}'::jsonb,
    objectifs JSONB NOT NULL,
    recompenses JSONB NOT NULL DEFAULT '{}'::jsonb,
    zone_depart VARCHAR(100),
    pnj_donneur VARCHAR(100),
    duree_limite_minutes INTEGER,
    repete BOOLEAN DEFAULT false,
    frequence_repetition VARCHAR(15) CHECK (frequence_repetition IN ('QUOTIDIEN', 'HEBDOMADAIRE')),
    dialogues JSONB DEFAULT '[]'::jsonb,
    last_event_sequence BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quetes_type ON quetes(type_quete);
CREATE INDEX idx_quetes_niveau ON quetes(niveau_requis);
CREATE INDEX idx_quetes_zone ON quetes(zone_depart);
CREATE INDEX idx_quetes_nom ON quetes(nom);
CREATE INDEX idx_quetes_repete ON quetes(repete) WHERE repete = true;
```

**Événements Sources**:
- QueteCreee → Création de la quête
- ObjectifProgresse → Mise à jour objectifs (au niveau Joueur)
- QueteTerminee → Stats globales de complétion

### ETAT_MONDE

État global du monde (singleton).

```sql
CREATE TABLE etat_monde (
    monde_id UUID PRIMARY KEY DEFAULT '00000000-0000-0000-0000-000000000000',
    cycle_jour_nuit INTEGER NOT NULL DEFAULT 0 CHECK (cycle_jour_nuit BETWEEN 0 AND 23),
    saison VARCHAR(10) NOT NULL DEFAULT 'PRINTEMPS' CHECK (saison IN ('PRINTEMPS', 'ETE', 'AUTOMNE', 'HIVER')),
    jour_saison INTEGER NOT NULL DEFAULT 1 CHECK (jour_saison BETWEEN 1 AND 30),
    evenements_actifs JSONB NOT NULL DEFAULT '[]'::jsonb,
    boss_vaincus JSONB NOT NULL DEFAULT '[]'::jsonb,
    zones_decouvertes JSONB NOT NULL DEFAULT '[]'::jsonb,
    ressources_mondiales JSONB NOT NULL DEFAULT '{}'::jsonb,
    facteurs_economiques JSONB NOT NULL DEFAULT '{}'::jsonb,
    classements JSONB NOT NULL DEFAULT '{}'::jsonb,
    derniere_mise_a_jour_cycle TIMESTAMP NOT NULL DEFAULT NOW(),
    last_event_sequence BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- S'assurer qu'il n'y a qu'une seule ligne
CREATE UNIQUE INDEX idx_etat_monde_singleton ON etat_monde((monde_id = '00000000-0000-0000-0000-000000000000'));
```

**Événements Sources**:
- EvenementMondialDeclenche → Ajout dans evenements_actifs
- ZoneDecouverte → Ajout dans zones_decouvertes
- BossVaincu → Ajout dans boss_vaincus
- RessourceRegeneree → Mise à jour ressources_mondiales

### ORDRES_ECONOMIE

Ordres d'achat/vente sur le marché.

```sql
CREATE TABLE ordres_economie (
    ordre_id UUID PRIMARY KEY,
    type_ordre VARCHAR(10) NOT NULL CHECK (type_ordre IN ('VENTE', 'ACHAT')),
    joueur_id UUID NOT NULL,
    pseudo_joueur VARCHAR(50) NOT NULL,
    item_id UUID NOT NULL REFERENCES items(item_id),
    quantite INTEGER NOT NULL CHECK (quantite > 0),
    quantite_restante INTEGER NOT NULL CHECK (quantite_restante >= 0),
    prix_unitaire INTEGER NOT NULL CHECK (prix_unitaire > 0),
    statut VARCHAR(25) NOT NULL DEFAULT 'ACTIF' CHECK (statut IN ('ACTIF', 'PARTIELLEMENT_EXECUTE', 'COMPLETE', 'ANNULE', 'EXPIRE')),
    cree_a TIMESTAMP NOT NULL,
    expire_a TIMESTAMP,
    complete_a TIMESTAMP,
    last_event_sequence BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ordres_economie_joueur ON ordres_economie(joueur_id);
CREATE INDEX idx_ordres_economie_item ON ordres_economie(item_id, statut);
CREATE INDEX idx_ordres_economie_type ON ordres_economie(type_ordre, statut);
CREATE INDEX idx_ordres_economie_statut ON ordres_economie(statut);
CREATE INDEX idx_ordres_economie_prix ON ordres_economie(item_id, type_ordre, prix_unitaire) WHERE statut = 'ACTIF';
```

**Événements Sources**:
- OrdreVenteCree → Création ordre de vente
- OrdreAchatCree → Création ordre d'achat
- TransactionExecutee → Mise à jour quantite_restante, statut
- OrdreAnnule → Mise à jour statut

### TRANSACTIONS_ECONOMIE

Historique des transactions du marché.

```sql
CREATE TABLE transactions_economie (
    transaction_id UUID PRIMARY KEY,
    ordre_vente_id UUID NOT NULL REFERENCES ordres_economie(ordre_id),
    ordre_achat_id UUID NOT NULL REFERENCES ordres_economie(ordre_id),
    vendeur_id UUID NOT NULL,
    pseudo_vendeur VARCHAR(50) NOT NULL,
    acheteur_id UUID NOT NULL,
    pseudo_acheteur VARCHAR(50) NOT NULL,
    item_id UUID NOT NULL REFERENCES items(item_id),
    quantite INTEGER NOT NULL,
    prix_unitaire INTEGER NOT NULL,
    montant_total INTEGER NOT NULL,
    taxe INTEGER NOT NULL DEFAULT 0,
    executee_a TIMESTAMP NOT NULL,
    event_sequence BIGINT NOT NULL
);

CREATE INDEX idx_transactions_vendeur ON transactions_economie(vendeur_id);
CREATE INDEX idx_transactions_acheteur ON transactions_economie(acheteur_id);
CREATE INDEX idx_transactions_item ON transactions_economie(item_id, executee_a DESC);
CREATE INDEX idx_transactions_date ON transactions_economie(executee_a);
```

**Événements Sources**:
- TransactionExecutee → Création de la transaction

### PRIX_MARCHE

Statistiques et historique des prix.

```sql
CREATE TABLE prix_marche (
    prix_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(item_id),
    prix_moyen_24h INTEGER,
    prix_min_24h INTEGER,
    prix_max_24h INTEGER,
    volume_24h INTEGER DEFAULT 0,
    nombre_transactions_24h INTEGER DEFAULT 0,
    historique_prix JSONB DEFAULT '[]'::jsonb,
    derniere_transaction TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_prix_marche_item ON prix_marche(item_id);
CREATE INDEX idx_prix_marche_volume ON prix_marche(volume_24h DESC);
```

**Événements Sources**:
- TransactionExecutee → Recalcul des statistiques
- PrixMisAJour → Mise à jour des données

## Requêtes Typiques

### 1. Recherche d'Items par Filtres

```sql
SELECT *
FROM items
WHERE type_item = $1
  AND rarete = ANY($2::text[])
  AND niveau_requis <= $3
  AND echangeable = true
ORDER BY nom;
```

### 2. Compétences Disponibles pour un Niveau

```sql
SELECT *
FROM competences
WHERE niveau_requis <= $1
  AND type_competence = 'ACTIVE'
ORDER BY element, nom;
```

### 3. Quêtes Disponibles pour un Joueur

```sql
SELECT q.*
FROM quetes q
WHERE q.niveau_requis <= $1
  AND q.niveau_recommande <= $1 + 5
  AND (q.repete = true OR q.quete_id NOT IN (
    -- Quêtes déjà complétées par le joueur
    SELECT quete_id FROM quetes_joueur WHERE joueur_id = $2 AND statut = 'COMPLETE'
  ))
ORDER BY q.type_quete, q.niveau_requis;
```

### 4. Meilleurs Prix d'Achat/Vente

```sql
-- Meilleurs prix de vente (pour acheter)
SELECT o.*, i.nom as item_nom
FROM ordres_economie o
JOIN items i ON o.item_id = i.item_id
WHERE o.type_ordre = 'VENTE'
  AND o.statut = 'ACTIF'
  AND o.item_id = $1
ORDER BY o.prix_unitaire ASC
LIMIT 10;

-- Meilleurs prix d'achat (pour vendre)
SELECT o.*, i.nom as item_nom
FROM ordres_economie o
JOIN items i ON o.item_id = i.item_id
WHERE o.type_ordre = 'ACHAT'
  AND o.statut = 'ACTIF'
  AND o.item_id = $1
ORDER BY o.prix_unitaire DESC
LIMIT 10;
```

### 5. Historique de Prix d'un Item

```sql
SELECT 
    pm.prix_moyen_24h,
    pm.prix_min_24h,
    pm.prix_max_24h,
    pm.volume_24h,
    pm.historique_prix
FROM prix_marche pm
WHERE pm.item_id = $1;
```

### 6. État Actuel du Monde

```sql
SELECT 
    cycle_jour_nuit,
    saison,
    jour_saison,
    evenements_actifs,
    classements
FROM etat_monde
WHERE monde_id = '00000000-0000-0000-0000-000000000000';
```

### 7. Transactions Récentes d'un Joueur

```sql
SELECT 
    t.*,
    i.nom as item_nom,
    i.icone_url
FROM transactions_economie t
JOIN items i ON t.item_id = i.item_id
WHERE t.vendeur_id = $1 OR t.acheteur_id = $1
ORDER BY t.executee_a DESC
LIMIT 50;
```

## Maintenance

### Recalcul des Prix du Marché

```sql
-- Job exécuté toutes les heures
UPDATE prix_marche pm
SET 
    prix_moyen_24h = subq.prix_moyen,
    prix_min_24h = subq.prix_min,
    prix_max_24h = subq.prix_max,
    volume_24h = subq.volume,
    nombre_transactions_24h = subq.nb_trans,
    derniere_transaction = subq.derniere,
    updated_at = NOW()
FROM (
    SELECT 
        item_id,
        CAST(AVG(prix_unitaire) AS INTEGER) as prix_moyen,
        MIN(prix_unitaire) as prix_min,
        MAX(prix_unitaire) as prix_max,
        SUM(quantite) as volume,
        COUNT(*) as nb_trans,
        MAX(executee_a) as derniere
    FROM transactions_economie
    WHERE executee_a >= NOW() - INTERVAL '24 hours'
    GROUP BY item_id
) subq
WHERE pm.item_id = subq.item_id;
```

### Expiration des Ordres

```sql
-- Job exécuté périodiquement
UPDATE ordres_economie
SET 
    statut = 'EXPIRE',
    updated_at = NOW()
WHERE statut = 'ACTIF'
  AND expire_a IS NOT NULL
  AND expire_a < NOW();
```

### Archivage des Transactions

```sql
-- Archiver les transactions > 90 jours
INSERT INTO transactions_economie_archive
SELECT * FROM transactions_economie
WHERE executee_a < NOW() - INTERVAL '90 days';

DELETE FROM transactions_economie
WHERE executee_a < NOW() - INTERVAL '90 days';
```

## Références

- **event_store.md**: Source de vérité
- **event_handlers.md**: ItemProjection, SkillProjection, QuestProjection, EconomyProjection, WorldProjection
- **matrice_evenements.md**: Structures des événements Item/Quete/Competence/Economy/WorldState
