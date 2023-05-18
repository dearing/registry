# Registry

*This library is intentionally simplistic and not meant for production.*

simple account management with Java Web Tokens (JWT)

Java Web Tokens (JWT) are a way to encode information in a token that can be verified by a server. This library uses JWTs to encode user information, `(jwt.payload.sub)`, and verify the username and password to the map of users.  This is a bcrypted hash from the accounts created before verification by the register function.  Functions to save and load `io.Writer, io.Reader` this as a Go Object Binary (GOB) file are included to ease persistence.

## Typical flow diagram (mermaid)
```mermaid
sequenceDiagram
    participant Client
    participant Server

    Client ->> Server: 1. Sends credentials (e.g., username, password)
    Server ->> Server: 2. Validates credentials
    alt Credentials Valid
        Server ->> Server: 3. Generates JWT
        Server -->> Client: 4. Sets JWT as a secure cookie
    else Credentials Invalid
        Server -->> Client: 4. Returns error message
    end
    Client ->> Server: 5. Sends requests with JWT cookie
    Server ->> Server: 6. Verifies JWT signature and integrity
    alt JWT Valid
        Server -->> Client: 7. Returns requested data
    else JWT Invalid
        Server -->> Client: 7. Returns authentication error
    end

```

## Notes
- JWT algo is arguable, here it is hard coded to HS256
- review tests for implementation details
- `go run cmd/*.go` to see an http client demo

```text
go run cmd/*.go
server online
Response Status: 200 OK
Response Body: User dearing Registered

Response Status: 200 OK
Response Body: Welcome back dearing, your login will expire at 2023-05-19T07:08:15Z

Response Status: 200 OK
Response Body: {"Header":{"alg":"HS256","typ":"JWT"},"Payload":{"sub":"dearing","iat":1684393695,"exp":1684480095}}
```
