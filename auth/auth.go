package login

import (
	// "bytes"
	"bytes"
	"context"
	"crypto"
	"database/sql"
	"encoding/json"
  "errors"
	// "fmt"
	"crypto/rsa"
	// "crypto/sha256"
	"log"
	"net/http"

	usersSql "github.com/alewilliam789/go-rest/db"
	"github.com/redis/go-redis/v9"
)

const COOKIE_NAME = "user-cookie"

type LoginAttempt struct {
  UserName string `json:"username"`
  Password string `json:"password"`
}

func getUserMessage(cookieVal string, client *redis.Client) ([]byte, error) {
  
  ctx := context.Background()

  userMessage, redisErr := client.Get(ctx,cookieVal).Result()

  if redisErr != nil {
    return []byte(""), redisErr
  }
  
  return []byte(userMessage), nil
}


func checkUserCookie(msgBytes []byte, keys *rsa.PrivateKey) error {

  verifyErr := rsa.VerifyPKCS1v15(&keys.PublicKey,crypto.SHA256,msgBytes,nil)
  
  if verifyErr != nil {
    return verifyErr
  }

  return nil
}

func checkUserCredentials(loginAttempt *LoginAttempt, queries *usersSql.Queries) error {
  
  ctx := context.Background()

  foundUser, userNotFoundErr := queries.GetUser(ctx, loginAttempt.UserName)

  if userNotFoundErr != nil {
    return userNotFoundErr
  }

  if !bytes.Equal(foundUser.Password, []byte(loginAttempt.Password)) {
    return errors.New("Mismatching credentials") 
  }

  return nil
}

func generateCookie() (*http.Cookie, error) {

  return nil, nil
}


func login(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries, client *redis.Client, key *rsa.PrivateKey) {
    
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
      if verificationErr := checkUserCredentials(&loginAttempt,queries); verificationErr != nil {
        w.WriteHeader(http.StatusUnauthorized)
        log.Print("Incorrect User Credentials")
        return
      }
      newCookie, cookieErr := generateCookie()
      
      if cookieErr != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Print("Issue generating cookie")
        return
      }

      http.SetCookie(w, newCookie)
      w.WriteHeader(http.StatusOK)
      return
    }
    
    _, redisErr := getUserMessage(cookie.Value, client)

    if redisErr != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Print(redisErr)
      return
    }

    // Verify that the correct key was used and send on merry way
  
    w.WriteHeader(http.StatusOK) 
}


func AuthorizeHandler(w http.ResponseWriter, req *http.Request, db *sql.DB, client *redis.Client, keys *rsa.PrivateKey) {
  
  queries := usersSql.New(db)
  
  if req.Method == "POST" {
    login(w,req,queries,client,keys)
  }
}
