# Structure de base Gin – Go

Ce dossier contient un squelette minimal pour démarrer un serveur Gin en Go.

## Lancer le serveur

```bash
cd server
go run main.go
```

## Points de départ
- `main.go` : point d’entrée, route `/ping` pour test
- `go.mod` : module Go initialisé avec Gin

## Pour aller plus loin
- Ajoute tes routes dans `main.go` ou structure le projet en dossiers (`handlers`, `routes`, `models`, etc.)
- Consulte la doc Gin : https://gin-gonic.com/docs/

---

**Exemple de requête :**

```bash
curl http://localhost:8080/ping
# Réponse : {"message":"pong"}
```
