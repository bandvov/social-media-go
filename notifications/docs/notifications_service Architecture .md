graph TD
    A[User] -->|Triggers Notification| B[Http Handler]
    B -->|Sends Request| C[Redis Request Channel]
    C -->|Notification Request| D[Notification Request Worker]
    D -->|Processes Request| E[Notification Service]
    E -->|Store Notification| F[PostgreSQL Database]
    E -->|Publish Processed Notification| G[Redis Processed Channel]
    G -->|Processed Notification| H[Notification Processed Worker]
    H -->|Send to SSE| I[SSE Handler]
    I -->|Real-time Notification| A

    style A fill:#f9f,stroke:#333,stroke-width:2px
    style B fill:#ccf,stroke:#333,stroke-width:2px
    style C fill:#ccf,stroke:#333,stroke-width:2px
    style D fill:#cfc,stroke:#333,stroke-width:2px
    style E fill:#ff9,stroke:#333,stroke-width:2px
    style F fill:#ff9,stroke:#333,stroke-width:2px
    style G fill:#ff9,stroke:#333,stroke-width:2px
    style H fill:#cfc,stroke:#333,stroke-width:2px
    style I fill:#ff9,stroke:#333,stroke-width:2px
