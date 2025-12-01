Extensions possibles pour UnitID
===============================
Le système actuel de génération d'UnitID utilise un singleton avec un mutex pour garantir l'unicité des identifiants dans un environnement concurrent. Cependant, plusieurs extensions et améliorations sont possibles pour répondre à des besoins futurs ou spécifiques.

### Extensions possibles
1. Support de différents types d'IDs
```go
type IDGenerator interface {
    Generate() UnitID
}

type UUIDGenerator struct{}
type SequentialGenerator struct{}
type CustomGenerator struct{}
```

2. Validation plus stricte
```go
func (id UnitID) IsValidUUID() bool {
    parts := strings.Split(id.value, "_")
    if len(parts) != 2 {
        return false
    }
    _, err := uuid.Parse(parts[1])
    return err == nil
}
```


3. Métriques de génération
```go
type UnitIDGenerator struct {
    mutex     sync.Mutex
    generated int64 // Compteur d'IDs générés
}

func (g *UnitIDGenerator) TotalGenerated() int64 {
    g.mutex.Lock()
    defer g.mutex.Unlock()
    return g.generated
}
```