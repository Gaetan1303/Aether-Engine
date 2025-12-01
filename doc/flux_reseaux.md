> **Note de synchronisation** :
> Ce diagramme de flux réseau utilise le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/agregats.md`.

```mermaid
flowchart TD
    %% Clients
    Client[Client Angular / Frontend] -->|WebSocket / RPC| Passerelle[API Gateway]

    %% Gateway
    Passerelle -->|Auth & Routage| ServiceJoueur[Service Joueur / Compte]
    Passerelle -->|Routage| ServiceCombat[Service Combat / Gameplay]
    Passerelle -->|Routage| ServiceMonde[Service Monde / World]
    Passerelle -->|Routage| ServiceQuete[Service Quêtes / Missions]
    Passerelle -->|Routage| ServiceEconomie[Service Économie / Craft]
    Passerelle -->|Routage| ServiceChat[Service Communication / Chat]

    %% Event Bus
    ServiceCombat -->|Evénement combat terminé| BusEvenement[Broker / Event Bus]
    ServiceQuete -->|Evénement progression| BusEvenement
    ServiceEconomie -->|Evénement craft / loot| BusEvenement
    BusEvenement --> ServiceJoueur
    BusEvenement --> ServiceMonde
    BusEvenement --> ServiceQuete
    BusEvenement --> ServiceEconomie

    %% Persistence
    ServiceJoueur -->|CRUD| DBJoueur[PostgreSQL / DB joueurs]
    ServiceMonde -->|Persiste état| DBMonde[PostgreSQL / DB monde]
    ServiceCombat -->|Log combats| EventStore[Event Store]
    ServiceQuete -->|Log progression| EventStore
    ServiceEconomie -->|Log économie| EventStore
    ServiceChat -->|Stocke messages| DBChat[NoSQL / Redis]

    %% Cache
    ServiceJoueur -->|Cache stats & inventaire| CacheRedis[Redis]
    ServiceMonde -->|Cache positions & état| CacheRedis
    ServiceCombat -->|Cache état combat| CacheRedis

    %% Monitoring
    BusEvenement -->|Métriques & logs| Monitoring[Prometheus / Grafana]

    %% Notes
    subgraph "Infrastructure"
        Passerelle
        BusEvenement
        DBJoueur
        DBMonde
        EventStore
        DBChat
        CacheRedis
        Monitoring
    end

```