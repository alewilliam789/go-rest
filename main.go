package main

import (
	"context"
	"crypto/rsa"
  "crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	auth "github.com/alewilliam789/go-rest/auth"
	users "github.com/alewilliam789/go-rest/users"
)

func setupDbConn() *sql.DB {
  cfg := mysql.Config{
    User: os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net: "tcp",
    Addr: os.Getenv("DBADDR"),
    DBName: os.Getenv("DB"),
  }

  db, err := sql.Open("mysql", cfg.FormatDSN())
  if err != nil {
    log.Fatal(err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connected")

  return db
}

func setupRedisConn() *redis.Client {
  client := redis.NewClient(&redis.Options{
    Addr: os.Getenv("RSADDR"),
    Password: os.Getenv("RSPASS"),
    DB: 0,
    Protocol: 2,
  })

  ctx := context.Background()

  status := client.Ping(ctx)

  if statusErr := status.Err(); statusErr != nil {
    log.Fatal(statusErr)
  }

  return client
}

func generateKeys() *rsa.PrivateKey {

  private_key, keyErr := rsa.GenerateKey(rand.Reader,2048)

  if keyErr != nil {
    log.Fatal(keyErr)
  }

  return private_key
}


func main() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal(err)
  }

  userDb := setupDbConn()
  defer userDb.Close()

  authClient := setupRedisConn()
  defer authClient.Close()

  // Filler for now while I write more for keys
  keys := generateKeys()

  http.HandleFunc("/v1/login", func(w http.ResponseWriter ,r *http.Request) {
    auth.AuthorizeHandler(w, r, userDb, authClient, keys)
  })
  http.HandleFunc("/v1/user",func(w http.ResponseWriter, r *http.Request) {
    users.UserHandler(w,r,userDb)  
  })
  http.HandleFunc("/v1/user/{username}", func(w http.ResponseWriter, r *http.Request) {
    users.UserNameHandler(w,r,userDb)
  })

  fmt.Printf("Starting server on 8080 \n")
  httpErr := http.ListenAndServe(":8080",nil)
  
  if httpErr != nil {
    log.Fatal(httpErr)
  }
}
