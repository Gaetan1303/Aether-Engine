
# Plan d’Architecture Global – Aether‑Engine MMO

## 1. Principes clés

- **Serveur autoritatif** : toute logique de combat, craft, quêtes, économie est validée côté serveur.
- **Microservices modulaires** : chaque domaine (combat, joueur, monde, économie, craft, quêtes, matchmaking) est un service indépendant, scalable.
- **Persistance robuste** : DB relationnelle pour les données critiques (PostgreSQL), cache pour performance (Redis), event store pour actions critiques et audit.
- **Scalabilité horizontale** : support des shards / instances / zones pour répartir les joueurs.
- **API client / protocole réseau** : WebSocket ou RPC binaire pour faible latence.

---

## 2. Couche Client

- **Client Frontend Angular** : affichage, interface utilisateur, animations, HUD, interactions.
- **Gestion locale** : cache local (ex. stockage temporaire d’inventaire pour UX fluide).
- **Communication réseau** : WebSocket / RPC → serveur (protocole optimisé pour jeux tour‑par‑tour MMO).

---

## 3. Couche Serveur

### a) Gateway / API Gateway
	- Reçoit toutes les requêtes du client
	- Authentification & validation de token
	- Routage vers les services internes (combat, joueur, monde…)
	- Rate limiting & logging

### b) Service Joueur / Account Service
	- Gestion des comptes, profils, stats, progression, inventaire
	- Persistance : PostgreSQL + Redis (cache des stats temps réel)
	- Validation côté serveur : anti-triche, vérification de cohérence des données

### c) Service Monde / World Service
	- Gestion des zones / instances / shards
	- Synchronisation des positions, état des objets persistants
	- Event sourcing pour modifications d’état du monde
	- Communication inter-shard si nécessaire

### d) Service Combat / Gameplay
	- Calcul des combats tour par tour
	- Gestion des effets de statut, initiative, dégâts, compétences, classes
	- Validation strictement côté serveur
	- Peut être scalable par zone ou par instance de combat

### e) Service Quêtes / Missions
	- Tracking de progression des quêtes
	- Déclencheurs et conditions de complétion
	- Interactions avec économie, combat, loot

### f) Service Économie / Craft
	- Gestion du craft, loot, commerce, drop tables
	- Suivi des ressources disponibles, équilibrage de l’économie
	- Event store pour audit et rollback

### g) Service Communication / Chat
	- Chat global, guildes, notifications
	- Peut être un microservice indépendant scalable horizontalement
	- Historique stocké dans DB NoSQL ou Redis

### h) Broker / Event Bus
	- Kafka / RabbitMQ pour communication asynchrone entre microservices
	- Événements critiques : combat terminé, craft terminé, objets créés, progression quêtes
	- Permet de découpler services et scaler indépendamment

---

## 4. Couche Persistance

- **PostgreSQL** : comptes joueurs, inventaire, progression, objets persistants
- **Redis** : cache temps réel pour stats joueurs, positions, effets actifs
- **Event Store** : logs d’actions critiques, audit et possibilité de rollback
- **Blob Storage / CDN** : pour assets persistants (images, sons, modèles)

---

## 5. Infrastructure & Scalabilité

- **Docker + Kubernetes** : déploiement multi-service, scaling horizontal
- **Shards / Zones** : chaque zone est une instance de world service
- **Load balancer** : répartit les connexions clients sur les gateways
- **Monitoring / Logs / Metrics** : Prometheus + Grafana + ELK
- **CI/CD** : pipeline automatisé avec tests unitaires et d’intégration

---

## 6. Flux de données typique

1. Client envoie action → Gateway → Service Joueur / Combat
2. Combat calculé → résultat → Service Monde pour mise à jour de l’état
3. Event généré → Event Bus → Service Économie ou Quêtes si nécessaire
4. Mise à jour persistée → Redis cache → retour client
5. Logs et métriques envoyés au monitoring

---

## 7. Points critiques à surveiller pour MMO

- **Synchronisation monde** : verrouillage ou transaction distribuée pour éviter les conflits
- **Latence réseau** : protocole optimisé pour tour‑par‑tour
- **Rollback / recovery** : en cas de crash serveur ou d’erreurs critique
- **Scalabilité de combat & zone** : partitionner combats et instances pour éviter goulets d’étranglement
- **Sécurité et validation côté serveur** : chaque action doit être validée