package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "math/rand"

	usersSql "github.com/alewilliam789/go-rest/db"
	"github.com/redis/go-redis/v9"
)


type LoginAttempt struct {
  ClientId int `json:"client_id"`
  UserName string `json:"username"`
  Password string `json:"password"`
}

func generateCode(code []byte) {
  // 25% chance of stopping
  const prob = 4
  const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~"  


  for len(code) < 128 {
    if len(code) >= 43 {
      flip := rand.Intn(prob)

      if flip == 0 {
        return
      }
    }
    
    randIndex := rand.Intn(len(characters)-1)
    randByte := byte(characters[randIndex])

    code = append(code,randByte)
  }
}

func authorize(w http.ResponseWriter, req *http.Request, cache *redis.Client, queries *usersSql.Queries) {

    ctx  := context.Background()

    var login_att LoginAttempt
    decoder := json.NewDecoder(req.Body)

    decodeErr := decoder.Decode(&login_att)

    if decodeErr != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Print("Unable to decode sent JSON")
    }

    found_user, userNotFoundErr := queries.GetUser(ctx, login_att.UserName)

    if userNotFoundErr != nil {
      w.WriteHeader(http.StatusForbidden)
      log.Print("Could not find requested user")
    }

    var access_code []byte;

    generateCode(access_code)

    if bytes.Equal(found_user.Password, []byte(login_att.Password)) {
      // Will store access_code with client id and send back to generate token    
    }

    fmt.Printf("Hello %s \n",login_att.UserName)

    w.WriteHeader(http.StatusOK) 
}



func AuthorizeHandler(w http.ResponseWriter, req *http.Request, client *redis.Client, db *sql.DB) {
  
  queries := usersSql.New(db)
  
  // Determine what HTTP METHOD for auth

}
