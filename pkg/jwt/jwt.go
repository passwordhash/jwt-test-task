package jwt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Err struct {
	reason string
	err    error
}

func (e *Err) Error() string {
	return fmt.Sprintf("jwt error: %s, details: %v", e.reason, e.err)
}

type Alg string

const (
	HS512 Alg = "HS512"
)

const (
	JWTType = "jwt"
)

type Header struct {
	Alg Alg    `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub any   `json:"sub"`
	Iat int64 `json:"iat"`
	Exp int64 `json:"exp"`
}

func NewToken(alg string, sub any, ttl time.Duration, secret string) (string, error) {
	var err error

	header := Header{
		Alg: Alg(alg), // TODO: validate alg
		Typ: JWTType,
	}

	headerBase64, err := encodeBase64(header)
	if err != nil {
		return "", &Err{reason: "encoding header in base 64 failed", err: err}
	}

	n := time.Now()
	p := Payload{
		Sub: sub,
		Iat: n.Unix(),
		Exp: n.Add(ttl).Unix(),
	}

	payloadBase64, err := encodeBase64(p)
	if err != nil {
		return "", &Err{reason: "encoding payload in base 64 failed", err: err}
	}

	signature, err := sign(header.Alg, headerBase64, payloadBase64, []byte(secret)) // TODO: returns string or byte slice ?
	if err != nil {
		return "", &Err{reason: "signing failed", err: err}
	}

	signatureBase64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(signature)))
	base64.RawURLEncoding.Encode(signatureBase64, signature)

	token := fmt.Sprintf("%s.%s.%s", string(headerBase64), string(payloadBase64), string(signatureBase64))

	return token, nil
}

func sign(alg Alg, encodedHeader, encodedPayload, secret []byte) ([]byte, error) {
	switch alg {
	case HS512:
		mac := hmac.New(sha512.New, secret)

		mac.Write(encodedHeader)
		mac.Write([]byte{'.'})
		mac.Write(encodedPayload)

		return mac.Sum(nil), nil
	default:
		return nil, &Err{reason: "unsupported algorithm", err: nil}
	}
}

func encodeBase64(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, &Err{reason: "failed to encode in base 64", err: err}
	}

	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(buf, b)

	return buf, nil
}

func ParseToken(token, secret string) (*Payload, error) {
	ErrParseToken := errors.New("failed to parse token")

	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		return nil, &Err{reason: "invalid token format", err: ErrParseToken}
	}

	header := []byte(parts[0])
	payload := []byte(parts[1])
	signature := []byte(parts[2])

	decodedSignature, err := decodeBase64(signature)
	if err != nil {
		return nil, &Err{reason: err.Error(), err: ErrParseToken}
	}

	// TODO: alg from header
	computedSignature, err := sign(HS512, header, payload, []byte(secret))
	if err != nil {
		return nil, &Err{reason: err.Error(), err: ErrParseToken}
	}

	if !hmac.Equal(decodedSignature, computedSignature) {
		// TODO: change
		return nil, &Err{reason: "invalid signature", err: ErrParseToken}
	}

	decodedPayload, err := decodeBase64(payload)
	if err != nil {
		return nil, &Err{reason: err.Error(), err: ErrParseToken}
	}

	p := new(Payload)
	if err := json.Unmarshal(decodedPayload, p); err != nil {
		return nil, &Err{reason: err.Error(), err: ErrParseToken}
	}

	return p, nil
}

func decodeBase64(data []byte) ([]byte, error) {
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	_, err := base64.RawURLEncoding.Decode(buf, data)

	return buf, err
}
