
> **Note de synchronisation** :
> Cette matrice d'événements utilise le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/agregats.md`.

```mermaid
flowchart LR


    subgraph Joueur[Agrégat Joueur]
        PE1(JoueurMonteNiveau)
        PE2(StatsJoueurModifiées)
        PE3(JoueurDéplacé)
        PE4(EtatJoueurModifié)
        PE5(JoueurMort)
        PE6(JoueurEquipeItem)
        PE7(JoueurDesequipeItem)
    end

    subgraph Combat[Agrégat Instance de Combat]
        CE1(CombatCommencé)
        CE2(TourCommencé)
        CE3(TourTerminé)
        CE4(ActionDéclarée)
        CE5(ActionRésolue)
        CE6(DégâtsAppliqués)
        CE7(SoinAppliqué)
        CE8(StatutAppliqué)
        CE9(StatutExpiré)
        CE10(CombatTerminé)
    end

    subgraph Inventaire[Agrégat Inventaire]
        IE1(ItemAjoutéInventaire)
        IE2(ItemRetiréInventaire)
        IE3(CapacitéInventaireAtteinte)
    end

    subgraph item[Agrégat item]
        IT1(ItemCréé)
        IT2(ItemConsommé)
        IT3(EffetItemDéclenché)
    end

    subgraph Equipement[Agrégat Equipement]
        EQ1(EquipementMisAJour)
        EQ2(BonusEquipementModifié)
    end

    subgraph Quete[Agrégat Journal de Quêtes]
        QE1(QueteCommencée)
        QE2(ObjectifQueteProgressé)
        QE3(QueteTerminée)
    end

    subgraph Compétence[Agrégat Compétence]
        SK1(CooldownCompétenceCommencé)
        SK2(CooldownCompétenceTerminé)
        SK3(CompétenceUtilisée)
    end

    subgraph PNJ[Agrégat PNJ/Monstre]
        NE1(ActionPnjEffectuée)
        NE2(PnjMort)
        NE3(PnjApparu)
    end

    subgraph Economy[Economy Aggregate]
        EC1(MarketOrderPlaced)
        EC2(MarketOrderFulfilled)
        EC3(MarketPriceUpdated)
    end

    subgraph World[WorldState Aggregate]
        WS1(WorldEventTriggered)
        WS2(InstanceCreated)
        WS3(InstanceClosed)
        WS4(DayNightCycleChanged)
    end
```