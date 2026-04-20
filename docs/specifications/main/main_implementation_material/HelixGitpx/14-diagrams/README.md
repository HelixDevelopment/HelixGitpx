# Diagrams

This folder holds the canonical **diagram sources** for HelixGitpx. All diagrams are maintained as code (Mermaid / PlantUML / D2) so they can be version-controlled, diff-reviewed, and regenerated.

Rendering is automatic in GitHub and on the docs site. Local: `mmdc -i <src>.mmd -o <out>.svg` or `plantuml -tsvg <src>.puml`.

---

## Index

| File | Kind | Describes |
|---|---|---|
| `c4-context.mmd`         | C4 Level 1 | Whole system boundary & external actors |
| `c4-container.mmd`       | C4 Level 2 | Microservices and their shared dependencies |
| `sequence-push-fanout.mmd` | Sequence | A push replicates to every upstream |
| `sequence-webhook-ingest.mmd` | Sequence | An upstream webhook lands in Kafka |
| `sequence-conflict-resolve.mmd` | Sequence | Detection → proposal → apply |
| `sequence-login-oidc.mmd` | Sequence | Browser OIDC login flow |
| `state-conflict-case.mmd` | State | `Case` status transitions |
| `state-pr.mmd`           | State | Pull-request state machine |
| `er-core.mmd`            | ER | Core entities (org→repo→ref→upstream) |
| `deployment-regions.mmd` | Deployment | Multi-region topology |
| `data-flow-events.mmd`   | Data flow | Producers → Kafka → Projectors → Read models |

---

## c4-context.mmd

```mermaid
C4Context
    title System Context — HelixGitpx

    Person(dev, "Developer", "Works locally; pushes to one upstream")
    Person(maint, "Maintainer", "Reviews PRs, resolves conflicts")
    Person(admin, "Org Admin", "Configures upstreams, policy, billing")

    System_Ext(gh,   "GitHub",    "Upstream Git host")
    System_Ext(gl,   "GitLab",    "Upstream Git host")
    System_Ext(gt,   "Gitee",     "Upstream Git host")
    System_Ext(oth,  "Others",    "Bitbucket, Gitea, …")
    System_Ext(oidc, "Identity Provider", "OIDC / SAML IdP")
    System_Ext(pay,  "Billing Provider", "Stripe / similar")

    System(helix, "HelixGitpx", "Multi-upstream Git federation with AI-assisted conflict resolution")

    Rel(dev,  helix, "Pushes / pulls, opens PRs, uses apps")
    Rel(maint, helix, "Reviews, resolves conflicts")
    Rel(admin, helix, "Configures, manages")

    Rel(helix, gh, "Sync git + API", "HTTPS")
    Rel(helix, gl, "Sync git + API", "HTTPS")
    Rel(helix, gt, "Sync git + API", "HTTPS")
    Rel(helix, oth, "Sync git + API", "HTTPS")
    Rel(helix, oidc, "Login / SSO", "OIDC")
    Rel(helix, pay, "Invoicing", "HTTPS / webhooks")
```

---

## c4-container.mmd

```mermaid
flowchart LR
    subgraph edge[Edge]
        CF[Cloudflare / CDN / WAF]
        EG[Envoy Gateway]
    end

    subgraph svc[HelixGitpx Services]
        GW[api-gateway]
        AUTH[auth-service]
        ORG[org-service]
        REPO[repo-service]
        UP[upstream-service]
        ADP[adapter-pool]
        WH[webhook-gateway]
        GI[git-ingress]
        SO[sync-orchestrator]
        CR[conflict-resolver]
        CRDT[crdt-service]
        POL[policy-service]
        AUD[audit-service]
        LE[live-events-service]
        AI[ai-service]
        SR[search-service]
        BIL[billing-service]
        NO[notify-service]
    end

    subgraph data[Data Plane]
        PG[(PostgreSQL)]
        K[(Kafka + Karapace)]
        RE[(Redis/Dragonfly)]
        OS[(OpenSearch)]
        MS[(Meilisearch)]
        QD[(Qdrant)]
        OBJ[(Object Store)]
        VAULT[(Vault)]
    end

    subgraph ext[External Git Hosts]
        GH[GitHub]
        GL[GitLab]
        GT[Gitee]
        OTR[Others]
    end

    CF --> EG --> GW
    GW --> AUTH
    GW --> REPO
    GW --> ORG
    GW --> UP
    GW --> SR
    GW --> LE
    GW --> AI

    WH --> K
    GI --> PG
    GI --> K
    REPO --> PG
    REPO --> K
    SO --> K
    SO --> ADP
    ADP --> GH
    ADP --> GL
    ADP --> GT
    ADP --> OTR
    CR --> K
    CR --> PG
    CR --> CRDT
    AI --> QD
    LE --> K
    LE --> RE
    SR --> OS
    SR --> MS
    SR --> QD
    AUD --> PG
    AUD --> K
    POL --> PG
    ORG --> PG
    AUTH --> PG
    AUTH --> RE

    PG <-- CDC --> K
    classDef db fill:#e7f5d0,stroke:#73a834,color:#1f3d0a;
    class PG,K,RE,OS,MS,QD,OBJ,VAULT db;
```

---

## sequence-push-fanout.mmd

```mermaid
sequenceDiagram
    autonumber
    actor Dev as Developer
    participant GI as git-ingress
    participant REPO as repo-service
    participant PG as Postgres
    participant K as Kafka
    participant SO as sync-orchestrator (Temporal)
    participant ADP as adapter-pool
    participant U1 as GitHub
    participant U2 as GitLab
    participant U3 as Gitee
    participant LE as live-events
    participant UI as Web / Mobile

    Dev->>GI: git push (git-receive-pack)
    GI->>PG: persist pack + ref update
    GI->>K: produce ref.updated (outbox)
    K-->>SO: consume ref.updated
    SO->>SO: start FanOutPush workflow
    par Parallel fan-out
      SO->>ADP: PushRef(GitHub)
      ADP->>U1: git push + API
      U1-->>ADP: ok
    and
      SO->>ADP: PushRef(GitLab)
      ADP->>U2: git push + API
      U2-->>ADP: ok
    and
      SO->>ADP: PushRef(Gitee)
      ADP->>U3: git push + API
      U3-->>ADP: rate-limited (429)
      ADP-->>SO: retry scheduled
    end
    SO->>K: produce sync.completed (partial)
    K-->>LE: fan-in to subscribers
    LE-->>UI: live event (WebSocket / gRPC stream)
```

---

## sequence-webhook-ingest.mmd

```mermaid
sequenceDiagram
    autonumber
    participant U as Upstream (e.g. GitHub)
    participant WH as webhook-gateway
    participant RE as Redis (dedup)
    participant K as Kafka
    participant CR as conflict-resolver
    participant REPO as repo-service
    participant PG as Postgres
    participant LE as live-events
    participant UI as Clients

    U->>WH: POST /webhooks/github/<id> (HMAC)
    WH->>WH: verify signature
    WH->>RE: setnx dedup key (event_id)
    alt first time
      WH->>K: produce upstream.webhook.received
      K-->>CR: consume
      CR->>REPO: detect divergence
      CR->>PG: persist conflict_case if any
      CR->>K: produce conflict.detected
      K-->>LE: fan-in
      LE-->>UI: live event
    else duplicate
      WH-->>U: 204 (ack)
    end
```

---

## sequence-conflict-resolve.mmd

```mermaid
sequenceDiagram
    autonumber
    participant CR as conflict-resolver
    participant POL as policy-service
    participant AI as ai-service
    participant S as sandbox (short-lived pod)
    participant SO as sync-orchestrator
    participant PG as Postgres
    participant LE as live-events
    participant UI as Clients

    CR->>POL: evaluate(case)
    alt policy-decisive
      POL-->>CR: allow strategy=prefer_primary
      CR->>SO: apply plan (fan-out)
      SO-->>CR: applied
      CR->>PG: update case → applied, undo_until = +5m
    else policy indeterminate
      CR->>AI: propose(k=3)
      AI-->>CR: proposal[confidence=0.94, plan]
      CR->>S: run patch in sandbox (compile/lint/test)
      alt sandbox ok
        CR->>SO: apply plan
        SO-->>CR: applied
        CR->>PG: applied + undo window
      else sandbox fail
        CR->>CR: try next proposal / escalate
        CR->>PG: case → escalated
      end
    end
    CR->>LE: emit conflict.resolved (or escalated)
    LE-->>UI: notification to assignees
```

---

## sequence-login-oidc.mmd

```mermaid
sequenceDiagram
    autonumber
    actor U as User
    participant W as Web App
    participant GW as api-gateway
    participant AUTH as auth-service
    participant IDP as OIDC IdP

    U->>W: click "Login"
    W->>AUTH: StartOIDCFlow(provider)
    AUTH-->>W: authorization_url + state + code_challenge
    W-->>U: redirect to IdP
    U->>IDP: authenticate
    IDP-->>W: redirect with code
    W->>AUTH: Login(oidc{code, code_verifier})
    AUTH->>IDP: token exchange + userinfo
    IDP-->>AUTH: id_token, access_token, claims
    AUTH->>AUTH: upsert user; issue JWT access + rotating refresh
    AUTH-->>W: TokenResponse + User
    W-->>U: logged in
```

---

## state-conflict-case.mmd

```mermaid
stateDiagram-v2
    [*] --> Detected
    Detected --> Proposed: propose()
    Detected --> Escalated: policy blocks auto
    Proposed --> AutoApplying: confidence >= threshold
    Proposed --> Escalated: confidence < threshold
    AutoApplying --> Applied: sandbox ok & apply ok
    AutoApplying --> Escalated: sandbox fail
    Applied --> Resolved: undo_window expired
    Applied --> Detected: undo() within 5 min
    Escalated --> HumanResolving: assigned / opened
    HumanResolving --> Applied: human accepts proposal
    HumanResolving --> Resolved: human rejects
    Resolved --> [*]
    Detected --> Cancelled: repo archived / unbound
    Escalated --> Cancelled: repo archived / unbound
```

---

## state-pr.mmd

```mermaid
stateDiagram-v2
    [*] --> Draft
    Draft --> Open: mark ready
    Open --> ChangesRequested: review
    ChangesRequested --> Open: push / re-request
    Open --> Approved: approvals met
    Approved --> Merged: merge
    Open --> Closed: close
    Approved --> Closed: close
    Merged --> [*]
    Closed --> Open: reopen
    Closed --> [*]
```

---

## er-core.mmd

```mermaid
erDiagram
    ORG ||--o{ TEAM : has
    ORG ||--o{ REPO : owns
    ORG ||--o{ UPSTREAM : connects
    TEAM ||--o{ MEMBERSHIP : has
    USER ||--o{ MEMBERSHIP : belongs
    REPO ||--o{ REF : has
    REPO ||--o{ BRANCH_PROTECTION : has
    REPO ||--o{ BINDING : bound_to
    UPSTREAM ||--o{ BINDING : for_repo
    REPO ||--o{ PR : has
    REPO ||--o{ ISSUE : has
    REPO ||--o{ RELEASE : has
    REPO ||--o{ CONFLICT_CASE : has
    CONFLICT_CASE ||--o{ RESOLUTION : has
    RESOLUTION ||--o{ AI_FEEDBACK : has
    PR ||--o{ REVIEW : has
    PR ||--o{ COMMENT : has
    ISSUE ||--o{ COMMENT : has
```

---

## deployment-regions.mmd

```mermaid
flowchart TB
    subgraph Cloudflare
        AnyCast[Global Anycast + WAF + DDoS]
    end

    subgraph EU[eu-west-1 Frankfurt]
      EUEG[Envoy Gateway]
      EUK8S[K8s cluster eu]
      EUPG[(Postgres primary)]
      EUKF[(Kafka)]
      EUS3[(Object store - S3 / MinIO)]
    end

    subgraph US[us-east-1 Virginia]
      USEG[Envoy Gateway]
      USK8S[K8s cluster us]
      USPG[(Postgres replica)]
      USKF[(Kafka)]
      USS3[(Object store)]
    end

    AnyCast --> EUEG
    AnyCast --> USEG

    EUPG <-- streaming --> USPG
    EUKF <-- MirrorMaker 2 --> USKF
    EUS3 <-- bi-dir replication --> USS3

    EUEG --> EUK8S
    USEG --> USK8S
    EUK8S --> EUPG
    EUK8S --> EUKF
    USK8S --> USPG
    USK8S --> USKF
```

---

## data-flow-events.mmd

```mermaid
flowchart LR
    subgraph Producers
      GI[git-ingress]
      WH[webhook-gateway]
      REPO[repo-service]
      ORG[org-service]
      SO[sync-orchestrator]
      CR[conflict-resolver]
      AUTH[auth-service]
      AUD[audit-service]
      BIL[billing-service]
      AI[ai-service]
    end

    subgraph Kafka["Kafka Cluster (with Karapace Schema Registry)"]
      T1[[helixgitpx.repo.*]]
      T2[[helixgitpx.upstream.*]]
      T3[[helixgitpx.sync.*]]
      T4[[helixgitpx.conflict.*]]
      T5[[helixgitpx.audit.events]]
      T6[[helixgitpx.ai.*]]
      T7[[helixgitpx.billing.usage]]
      T8[[helixgitpx.notify.*]]
      DLQ[[*.dlq]]
    end

    subgraph Projectors
      P1[repo-projector → Postgres]
      P2[search-projector → Meilisearch / OpenSearch]
      P3[vector-projector → Qdrant]
      P4[audit-projector → Postgres + OpenSearch]
      P5[billing-meter → Postgres]
      P6[notify-dispatcher → channels]
    end

    GI --> T1
    REPO --> T1
    WH --> T2
    SO --> T3
    CR --> T4
    AUD --> T5
    AI --> T6
    BIL --> T7
    T1 --> P1
    T1 --> P2
    T1 --> P3
    T5 --> P4
    T7 --> P5
    T4 --> P6
    T1 -. failures .-> DLQ
    T2 -. failures .-> DLQ
    T3 -. failures .-> DLQ
```

---

## Contribution Notes

- Keep diagram complexity low — split into multiple files rather than one massive diagram.
- Names in diagrams must match the glossary.
- When architecture changes, update diagrams in the same PR as code.
- Review PRs: if a diagram is impacted, reviewer must confirm freshness.
