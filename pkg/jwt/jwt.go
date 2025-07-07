package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type T struct{}

type Error struct {
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("jwt error: %s, details: %v", e.Message, e.Err)
}

type Alg string

const (
	HS512 Alg = "HS512"
)

// header.payload.signature

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub any   `json:"sub"`
	Iat int64 `json:"iat"`
	Exp int64 `json:"exp"`
}

// TODO: validate alg
func (*T) NewToken(alg string, sub any, ttl time.Duration) (string, error) {
	header := Header{
		Alg: string(alg),
		Typ: "jwt",
	}

	n := time.Now()
	p := Payload{
		Sub: sub,
		Iat: n.Unix(),
		Exp: n.Add(ttl).Unix(),
	}

	encodedHeader, err := encode(header)
	if err != nil {
		return "", err
	}

	encodedPayload, err := encode(p)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)

	return token, nil
}

func encode(data any) (encoded []byte, encErr *Error) {
	encErr = &Error{Message: "failed to encode data"}

	b, err := json.Marshal(data)
	if err != nil {
		encErr.Err = err
		return
	}

	buf := bytes.Buffer{}

	enc := base64.NewEncoder(base64.URLEncoding, &buf)
	defer func() { _ = enc.Close() }()

	if _, err := enc.Write(b); err != nil {
		encErr.Err = err
		return
	}

	encoded = buf.Bytes()
	if len(encoded) == 0 {
		encErr.Err = fmt.Errorf("encoded data is empty")
		return
	}

	// TODO: remove padding

	return encoded, nil
}
