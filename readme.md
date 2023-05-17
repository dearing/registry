# Registry

Simple account management with Java Web Tokens (JWT).

## About

**Warning**

This library is intentionally simplistic and not meant for production.

Java Web Tokens (JWT) are a way to encode information in a token that can be verified by a server. This library uses JWTs to encode user information, `(jwt.payload.sub)`, and verify compariing the username and password to the map of users.  This is a bcrypted hash from the accounts created before verification by the register function.  Functions to save this a Go Object Binary (GOB) file are included to ease persistence.

**Note**

No on disk encryption, or other security measures are taken.  Though, you'd potentially only expose the usernames and hashed passwords.
