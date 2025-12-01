```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> SelectingSkill
    SelectingSkill --> SelectingTarget
    SelectingTarget --> Confirming
    Confirming --> Submitted
    Submitted --> [*]
    Confirming --> Idle: Cancel
```