package registry

import (
	"testing"
	"time"
)

var secretKey = []byte("chumpy")

func TestJWT(t *testing.T) {

	// Create a new JWT object and populate its fields
	jwt := NewJWT("user123", time.Now().Add(time.Hour).Unix())

	token, err := jwt.Encode(secretKey)
	if err != nil {
		t.Error("Failed to encode JWT:", err)
	}

	decodedJwt, err := Decode(token, secretKey)
	if err != nil {
		t.Error("Failed to decode JWT:", err)
	}

	// Check that the decoded JWT matches the original one

	if jwt.Header != decodedJwt.Header {
		t.Errorf("Both JWT token headers should be equal  %v != %v", jwt, decodedJwt)
	}

	if jwt.Payload != decodedJwt.Payload {
		t.Errorf("Both JWT token payloads should be equal  %v != %v", jwt, decodedJwt)
	}

	// Check that the decoded JWT has our chosen constants (constraints)

	if decodedJwt.Header.Algorithm != "HS256" {
		t.Errorf("JWT token algorithm should only be HS256  %v != %v", jwt.Header, "HS256")
	}

	if decodedJwt.Header.Type != "JWT" {
		t.Errorf("JWT token type should only be JWT  %v != %v", jwt.Header, "JWT")
	}

	// Check that the decoded JWT is handling as expected

	if decodedJwt.Payload.ExpiresAt <= decodedJwt.Payload.IssuedAt {
		t.Errorf("JWT token expiry should be later than issued date  %v >= %v", jwt.Payload.ExpiresAt, jwt.Payload.IssuedAt)
	}

	// check that the expiry is one hour after the issued date
	if time.Unix(int64(decodedJwt.Payload.ExpiresAt), 0).Sub(time.Unix(int64(decodedJwt.Payload.IssuedAt), 0)) != time.Hour {
		t.Errorf("JWT token expiry should be 1 hour after issued date  %v != %v", time.Unix(int64(decodedJwt.Payload.ExpiresAt), 0).Sub(time.Unix(int64(decodedJwt.Payload.IssuedAt), 0)), time.Hour)
	}

}
