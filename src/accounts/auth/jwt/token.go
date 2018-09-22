// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package jwt implements a Json Web Tokens
package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"strings"
	"time"
)

// JWT is used to sign and validate a token.
type JWT struct {
	signingHash hash.Hash
	algorithm   string
}

// NewHeader returns a new Header object.
func (token *JWT) NewHeader() *Header {
	return &Header{
		Typ: "JWT",
		Alg: token.algorithm,
	}
}

func (token *JWT) sum(data []byte) []byte {
	return token.signingHash.Sum(data)
}

func (token *JWT) reset() {
	token.signingHash.Reset()
}

func (token *JWT) write(data []byte) (int, error) {
	return token.signingHash.Write(data)
}

// Sign signs the token with the given hash, and key
func (token *JWT) Sign(unsignedToken string) ([]byte, error) {
	_, err := token.write([]byte(unsignedToken))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to write to %s", token.algorithm))
	}
	encodedToken := token.sum(nil)
	token.reset()
	return encodedToken, nil
}

// Encode returns an encoded JWT token from a header, payload, and secret
func (token *JWT) Encode(payload *Claims) (string, error) {
	header := token.NewHeader()
	jsonTokenHeader, err := json.Marshal(header)
	if err != nil {
		return "", ErrTokenUnableToMarshallHeader
	}

	b64TokenHeader := base64.RawURLEncoding.EncodeToString(jsonTokenHeader)
	jsonTokenPayload, err := json.Marshal(payload.claims)
	if err != nil {
		return "", ErrTokenUnableToMarshallPayload
	}

	b64TokenPayload := base64.RawURLEncoding.EncodeToString(jsonTokenPayload)
	unsignedSignature := b64TokenHeader + "." + b64TokenPayload
	signature, err := token.Sign(unsignedSignature)
	if err != nil {
		return "", ErrTokenUnableToSign
	}
	b64Signature := base64.RawURLEncoding.EncodeToString([]byte(signature))
	jwtToken := b64TokenHeader + "." + b64TokenPayload + "." + b64Signature

	return jwtToken, nil
}

// Decode returns a map representing the token's claims. DOESN'T validate the claims though.
func (token *JWT) Decode(encoded string) (*Claims, error) {
	encryptedComponents := strings.Split(encoded, ".")
	if len(encryptedComponents) != 3 {
		return nil, ErrTokenIsMalformed
	}
	b64Payload := encryptedComponents[1]

	var claims map[string]interface{}
	payload, err := base64.RawURLEncoding.DecodeString(b64Payload)
	if err != nil {
		return nil, ErrTokenUnableToDecodeB64Payload
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrTokenUnableToUnmarshallPayload
	}
	return &Claims{
		claims: claims,
	}, nil
}

// Validate verifies a token's validity. It returns nil if it is valid, and an error if invalid.
func (token *JWT) Validate(encoded string) error {
	_, err := token.DecodeAndValidate(encoded)
	return err
}

// DecodeAndValidate returns a map representing the token's claims, and it's valid.
func (token *JWT) DecodeAndValidate(encoded string) (claims *Claims, err error) {
	claims, err = token.Decode(encoded)
	if err != nil {
		return
	}
	if err = token.validateSignature(encoded); err != nil {
		err = errors.New(fmt.Sprintf("failed to validate signature: %s", err.Error()))
		return
	}
	if err = token.validateExp(claims); err != nil {
		err = errors.New(fmt.Sprintf("failed to validate exp: %s", err.Error()))
		return
	}
	if err = token.validateNbf(claims); err != nil {
		err = errors.New(fmt.Sprintf("failed to validate nbf: %s", err.Error()))
	}
	return
}

// validateSignature verifies a token's signature.
func (token *JWT) validateSignature(encoded string) error {
	encryptedComponents := strings.Split(encoded, ".")

	b64Header := encryptedComponents[0]
	b64Payload := encryptedComponents[1]
	b64Signature := encryptedComponents[2]

	unsignedAttempt := b64Header + "." + b64Payload
	signedAttempt, err := token.Sign(unsignedAttempt)
	if err != nil {
		return ErrTokenUnableToSign
	}

	b64SignedAttempt := base64.RawURLEncoding.EncodeToString([]byte(signedAttempt))

	if !hmac.Equal([]byte(b64Signature), []byte(b64SignedAttempt)) {
		return ErrTokenInvalidSignature
	}

	return nil
}

// validateExp verifies a token's exp claim.
func (token *JWT) validateExp(claims *Claims) error {
	if claims.Contains("exp") {
		exp, err := claims.GetTime("exp")
		if err != nil {
			return err
		}
		if exp.Before(time.Now()) {
			return ErrTokenHasExpired
		}
	}
	return nil
}

// validateNbf verifies a token's nbf claim.
func (token *JWT) validateNbf(claims *Claims) error {
	if claims.Contains("nbf") {
		nbf, err := claims.GetTime("nbf")
		if err != nil {
			return err
		}
		if nbf.After(time.Now()) {
			return ErrTokenNotValid
		}
	}
	return nil
}
