flowchart TD

subgraph Ticker["P2.T2 - BattleTicker (ATB Engine)"]
    Tick[Tick() : Avancer le temps] --> ReadyUnits["1. Ready Units (Triées par Initiative)"]
end

subgraph Input["Input Externe (Client/API)"]
    CommandInput[Command Input (OpID, UnitID, Action)]
end

subgraph Resolver["P2.T4 - Turn Resolver (L'Orchestrateur Battle Aggregate)"]
    Queue[Active Actor Queue (File d'Attente d'Action)]
    
    ReadyUnits --> Queue
    
    subgraph ResolutionPipeline["P3 - Pipeline de Résolution"]
        CommandInput --> Validate[2. Valider (Tour, Portée, Coût)]
        Validate -- Valid --> Execute[3. Exécuter la Command]
        Execute --> Aggregate[Battle Aggregate (L'État)]
        Execute --> ResetATB[4. Réinitialiser ATB (Appel Ticker.ResetATB)]
        Execute --> Event[5. Émettre Événement (P4 - Event Sourcing)]
    end
end

Aggregate --> Aggregate
ResetATB --> Aggregate

style Aggregate fill:#f9f,stroke:#333
style Queue fill:#ccf,stroke:#333
style ReadyUnits fill:#ccf,stroke:#333
