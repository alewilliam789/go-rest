package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
  "time"
)

type JWTHeader struct {
  Alg string `json:"alg"`
  Type string `json:"typ"`
  Base64 string `json:"-"`
}

type JWTPayload struct {
  Iss string `json:"iss"`
  Sub string `json:"sub"`
  Iat int64 `json:"iat"`
  Base64 string `json:"-"`
}

type JWT struct {
  Header JWTHeader
  Payload JWTPayload
  Token string
  Signature string
}

func (jwt *JWT) New(userName string) {
  
  header := JWTHeader {
    Alg : "SHA256",
    Type : "JWT",
  }

  payload := JWTPayload {
    Iss: "user-api",
    Sub: userName,
    Iat: time.Now().Unix(),
  }

  jwt.Header = header
  jwt.Payload = payload
}

func (jwt *JWT) GenerateSig(key *rsa.PrivateKey) error {
  
  headerBytes, marshalErr := json.Marshal(jwt.Header)

  if marshalErr != nil {
    return errors.New("Issue marshaling JWT Header")
  }

  jwt.Header.Base64 = base64.StdEncoding.EncodeToString(headerBytes)
  
  payloadBytes, marshalErr := json.Marshal(jwt.Payload)

  if marshalErr != nil {
    return errors.New("Issue marshaling JWT Payload")
  }
  
  jwt.Payload.Base64 = base64.StdEncoding.EncodeToString(payloadBytes)
  
  signingString := fmt.Sprintf("%s.%s",jwt.Header.Base64,jwt.Payload.Base64)
  
  hasher := sha256.New()

  hashed := hasher.Sum([]byte(signingString))

  sigBytes, signingErr := rsa.SignPKCS1v15(nil,key,crypto.SHA256,hashed)

  if signingErr != nil {
    return errors.New("Issue signing the JWT with the private key")
  }

  jwt.Signature = base64.StdEncoding.EncodeToString(sigBytes)

  jwt.Token = fmt.Sprintf("%s.%s.%s",jwt.Header.Base64,jwt.Payload.Base64,jwt.Signature)
    
  return nil
}

