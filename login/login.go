package login

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
  "crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"log"
	"net/http"
	"strings"

	usersSql "github.com/alewilliam789/go-rest/db"
	webToken "github.com/alewilliam789/go-rest/jwt"
)

const COOKIE_NAME = "user-cookie"

type LoginAttempt struct {
  UserName string `json:"username"`
  Password string `json:"password"`
  Hashed []byte `json:"-"`
}

func hashPass(loginAttempt *LoginAttempt) error {
  hasher := sha256.New()

  _, err := hasher.Write([]byte(loginAttempt.Password))

  if err != nil {
    return err
  }

  loginAttempt.Hashed = hasher.Sum(nil)

  return nil
}

func checkUserJWT(jwt string, keys *rsa.PrivateKey) error {

  jwtSegments := strings.Split(jwt,".")

  if len(jwtSegments) != 3 {
     return errors.New("Malformed JWT provided")
  }

  signingString := fmt.Sprintf("%s.%s",jwtSegments[0:0],jwtSegments[1:1])
  
  hasher := sha512.New()

  hashed := hasher.Sum([]byte(signingString))

  verifyErr := rsa.VerifyPKCS1v15(&keys.PublicKey,crypto.SHA512,hashed,[]byte(jwtSegments[2]))
  
  if verifyErr != nil {
    return verifyErr
  }

  return nil
}

func checkUserCredentials(loginAttempt *LoginAttempt, queries *usersSql.Queries) error {
  
  ctx := context.Background()
  defer ctx.Done()

  foundUser, userNotFoundErr := queries.GetUser(ctx, loginAttempt.UserName)

  if userNotFoundErr != nil {
    return userNotFoundErr
  }

  hashPass(loginAttempt)

  if !bytes.Equal(foundUser.Password, loginAttempt.Hashed) {
    return errors.New("Mismatching credentials") 
  }
  
  return nil
}

func generateCookie(userLogin *LoginAttempt,key *rsa.PrivateKey) (*http.Cookie, error) {
  
  var jwt webToken.JWT;

  jwt.New(userLogin.UserName)

  sigErr := jwt.GenerateSig(key)

  if sigErr != nil {
    return nil, sigErr
  }

  newCookie := http.Cookie {
    HttpOnly: true,
    Name: COOKIE_NAME,
    Value : jwt.Token,
  }

  duration, parseErr := time.ParseDuration("48h")

  if parseErr != nil {
    return nil, parseErr
  }

  newCookie.Expires.Add(duration)

  return &newCookie, nil
}


func login(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries, key *rsa.PrivateKey) {
    
    // Decode request body into login attempt
    decoder := json.NewDecoder(req.Body)

    var loginAttempt LoginAttempt

    decodeErr := decoder.Decode(&loginAttempt)

    if decodeErr != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Print("Unable to decode sent JSON")
    }

    // Try to retrieve cookie and check/generate
    cookie, cookieErr := req.Cookie(COOKIE_NAME)
  
    if cookieErr != nil {
      verificationErr := checkUserCredentials(&loginAttempt,queries)
      
      if verificationErr != nil {
        w.WriteHeader(http.StatusUnauthorized)
        log.Print("Incorrect User Credentials")
        return
      }

      newCookie, cookieErr := generateCookie(&loginAttempt,key)
      
      if cookieErr != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Print("Issue generating cookie")
        return
      }

      http.SetCookie(w, newCookie)
      w.WriteHeader(http.StatusOK)
      return
    }
    

    // Verify that the correct key was used and send on merry way
    signErr := checkUserJWT(cookie.Value, key)

    if signErr != nil {
      verificationErr := checkUserCredentials(&loginAttempt, queries)

      if verificationErr != nil {
        w.WriteHeader(http.StatusUnauthorized)
        log.Print("Incorrect User Credentials")
        return
      }

      newCookie, cookieErr := generateCookie(&loginAttempt, key)

      if cookieErr != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Print("Issue generating cookie")
        return
      }

      http.SetCookie(w, newCookie)
    }

    w.WriteHeader(http.StatusOK) 
}


func AuthorizeHandler(w http.ResponseWriter, req *http.Request, db *sql.DB, keys *rsa.PrivateKey) {
  
  queries := usersSql.New(db)
  
  if req.Method == "POST" {
    login(w,req,queries,keys)
  }
}
