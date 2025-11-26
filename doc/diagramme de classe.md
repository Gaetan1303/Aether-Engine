# Diagramme de classe – Tactical RPG

```mermaid
classDiagram
class Unite {
  +id : IdentifiantUnite
  +statistiques : Statistiques
  +statuts : Statut[]
  +position : Position
  +appliquerDegats()
  +appliquerStatut()
}
class Equipe {
  +id : IdentifiantEquipe
  +membres : Unite[]
}
class GrilleDeCombat {
  +largeur
  +hauteur
  +cases
}
class Combat {
  +id : IdentifiantCombat
  +equipes : Equipe[]
  +grille : GrilleDeCombat
  +ordreDeTour : OrdreDeTour
}
Unite -- Equipe
Equipe -- Combat
Combat -- GrilleDeCombat
```

---

**Légende** :
- `+` : attribut ou méthode publique
- Les relations montrent la composition/agrégation entre les entités principales du moteur de combat.

Ce diagramme modélise la structure objet du cœur du serveur tactical RPG.