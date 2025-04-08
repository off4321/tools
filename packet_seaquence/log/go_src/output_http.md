sequenceDiagram
    participant Unknown
    participant Remote
    participant 65_208_228_223
    participant 145_254_160_237
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: HTTP
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: DNS
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: DNS
    Unknown->>Remote: HTTP
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 80→3371
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP
    Unknown->>Remote: HTTP
    Unknown->>Remote: TCP 3371→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3371
    Unknown->>Remote: TCP 28→80
    65_208_228_223->>145_254_160_237: HTTP/XML
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 3372→80
    Unknown->>Remote: TCP 80→3372
