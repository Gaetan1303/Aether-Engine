

# Diagramme de classe global – Tactical RPG (francisé)

> **Note de synchronisation** : Ce diagramme utilise le nommage français pour tous les concepts, sauf les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).
> Les définitions d'agrégats et Value Objects sont centralisées dans `/doc/agregats.md`.


```mermaid
classDiagram
class IdentifiantUnite {
  +valeur : string
  +equals(autre : IdentifiantUnite) bool
}
class Position {
  +x : int
  +y : int
  +equals(autre : Position) bool
}
class Statistiques {
  +pv : int
  +pm : int
  +atk : int
  +def : int
  +vit : int
  +mag : int
  +res : int
}
class Unite {
  +id : IdentifiantUnite
  +nom : string
  +position : Position
  +stats : Statistiques
  +statuts : Statut[]
  +appliquerDegats(valeur : int)
  +appliquerStatut(statut : Statut)
  +deplacer(nouvellePosition : Position)
}
class Statut {
  +type : string
  +duree : int
}
class Equipe {
  +id : string
  +membres : Unite[]
  +ajouterUnite(unite : Unite)
  +retirerUnite(unite : Unite)
}
class Competence {
  +nom : string
  +puissance : int
  +cible : string
  +appliquer(cible : Unite)
}
class GrilleDeCombat {
  +largeur : int
  +hauteur : int
  +cases : Case[][]
  +estAccessible(position : Position) bool
}
class Case {
  +type : string
  +altitude : int
}
class OrdreDeTour {
  +ordre : Unite[]
  +prochain() : Unite
}
class Combat {
  +id : string
  +equipes : Equipe[]
  +grille : GrilleDeCombat
  +ordreDeTour : OrdreDeTour
  +etat : string
  +demarrer()
  +executerTour()
  +terminer()
}
IdentifiantUnite <|-- Unite
Position <|-- Unite
Statistiques <|-- Unite
Statut <|-- Unite
Unite -- Equipe
Equipe -- Combat
Combat -- GrilleDeCombat
GrilleDeCombat -- Case
Combat -- OrdreDeTour
Unite -- Competence
```

---

