```mermaid
sequenceDiagram
    participant Client
    participant API Gateway
    participant Auth Middleware
    participant Reverse Proxy
    participant Comments Service
    participant Database

    Client->>API Gateway: POST / (comment data)
    API Gateway->>Auth Middleware: Validate authentication
    Auth Middleware-->>API Gateway: Auth result (success/failure)
    API Gateway-->>Client: 401 Unauthorized (if failure)
    API Gateway->>Reverse Proxy: Forward request (if success)
    Reverse Proxy->>Comments Service: Forward request 
    Comments Service->>Database: Insert comment into DB
    Database-->>Comments Service: Acknowledgment
    Comments Service-->>Reverse Proxy: Response (success/failure)
    Reverse Proxy-->>API Gateway: Response (success/failure)
    API Gateway-->>Client: Response (success/failure)


```