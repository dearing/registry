# Registry

simple account management with Java Web Tokens (JWT)

Java Web Tokens (JWT) are a way to encode information in a token that can be verified by a server. This library uses JWTs to encode user information, `(jwt.payload.sub)`, and verify the username and password to the map of users.  This is a bcrypted hash from the accounts created before verification by the register function.  Functions to save and load, taking a `io.Writer, io.Reader` respectively, the accounts map as a Go Object Binary (GOB) file are included to ease persistence.

```mermaid
sequenceDiagram
    participant Client
    participant Server

    Client ->> Server: 1. Register (Send registration data)
    Server ->> Server: 2. Process registration
    alt Registration Successful
        Server -->> Client: 3. Return registration success message
    else Registration Failed
        Server -->> Client: 3. Return registration failure message
    end

    Client ->> Server: 4. Login (Send login credentials)
    Server ->> Server: 5. Authenticate user
    alt Authentication Successful
        Server -->> Client: 6. Return JWT (JSON Web Token)
    else Authentication Failed
        Server -->> Client: 6. Return authentication failure message
    end

    Client ->> Server: 7. Request session information (Send JWT in request)
    Server ->> Server: 8. Verify JWT
    alt JWT Verification Successful
        Server -->> Client: 9. Return session information
    else JWT Verification Failed
        Server -->> Client: 9. Return JWT verification failure message
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
