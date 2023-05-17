package registry

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTHeader represents the header section of a JWT token.
type JWTHeader struct {
	Algorithm string `json:"alg"` // Algorithm
	Type      string `json:"typ"` // Type
}

// JWTPayload represents the payload section of a JWT token.
type JWTPayload struct {
	Subject   string    `json:"sub"` // Subject
	IssuedAt  Timestamp `json:"iat"` // Issued At
	ExpiresAt Timestamp `json:"exp"` // Expiration Time
}

type Timestamp int64

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t Timestamp) HumanReadable() string {
	return time.Unix(int64(t), 0).UTC().Format(time.RFC3339)
}

// JWT represents a JSON Web Token.
type JWT struct {
	Header  JWTHeader  // Header
	Payload JWTPayload // Payload
}

// NewJWT creates a new JWT object.
func NewJWT(sub string, expiry int64) *JWT {
	return &JWT{
		Header: JWTHeader{
			Algorithm: "HS256",
			Type:      "JWT",
		},
		Payload: JWTPayload{
			Subject:   sub,
			IssuedAt:  Timestamp(time.Now().Unix()),
			ExpiresAt: Timestamp(expiry),
		},
	}
}

// Encode encodes a with your secret key, returning a full JWT token.
func (jwt *JWT) Encode(secretKey []byte) (string, error) {
	// Encode the header and payload as base64
	headerBytes, err := json.Marshal(jwt.Header)
	if err != nil {
		return "", fmt.Errorf("failed to encode header: %v", err)
	}
	payloadBytes, err := json.Marshal(jwt.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode payload: %v", err)
	}
	headerBase64 := base64.RawStdEncoding.EncodeToString(headerBytes)
	payloadBase64 := base64.RawStdEncoding.EncodeToString(payloadBytes)

	// Compute the signature using HMAC-SHA256
	signatureInput := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(signatureInput))
	signature := base64.RawStdEncoding.EncodeToString(mac.Sum(nil))

	// Join the header, payload, and signature with dots to create the complete JWT
	token := fmt.Sprintf("%s.%s.%s", headerBase64, payloadBase64, signature)

	return token, nil
}

// Decode decodes a JWT as a string and secret key, returning a JWT object.
func Decode(token string, secretKey []byte) (*JWT, error) {
	// Split the token into its header, payload, and signature parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the header and payload from base64
	headerBytes, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}
	payloadBytes, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	// Parse the header and payload as JSON
	header := JWTHeader{}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}
	payload := JWTPayload{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %v", err)
	}

	// Verify the signature using HMAC-SHA256
	signatureInput := fmt.Sprintf("%s.%s", parts[0], parts[1])
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(signatureInput))
	expectedSignature := base64.RawStdEncoding.EncodeToString(mac.Sum(nil))
	actualSignature := parts[2]
	if expectedSignature != actualSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// Verify the expiration time
	now := time.Now().Unix()
	if payload.ExpiresAt < Timestamp(now) {
		return nil, fmt.Errorf("token has expired")
	}

	// Create a new JWT object and populate its fields
	jwt := &JWT{
		Header:  header,
		Payload: payload,
	}

	return jwt, nil
}
