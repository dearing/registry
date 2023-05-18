# Registry

Simple account management with Java Web Tokens (JWT).

## About

**Warning** : *This library is intentionally simplistic and not meant for production.*

Java Web Tokens (JWT) are a way to encode information in a token that can be verified by a server. This library uses JWTs to encode user information, `(jwt.payload.sub)`, and verify compariing the username and password to the map of users.  This is a bcrypted hash from the accounts created before verification by the register function.  Functions to save and load `io.Writer, io.Reader` this as a Go Object Binary (GOB) file are included to ease persistence.

## Notes
- No on-disk encryption, or other security measures are taken.
- JWT algo is arguable, here it is hard coded to HS256
- review tests for implementation details
