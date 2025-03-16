package main

import (
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

func setupRedisConn() *redis.Conn {
  client := redis.NewClient(&redis.Options{
    Addr: os.Getenv("RSADDR"),
    Password: os.Getenv("RSPASS"),
    DB: 0,
    Protocol: 2,
  })

  return client.Conn()
}


func main() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal(err)
  }


  userDb := setupDbConn()
  defer userDb.Close()

  authCache := setupRedisConn()
  defer authCache.Close()

  http.HandleFunc("/login", func(w http.ResponseWriter ,r *http.Request) {
    auth.Login(w, r)
  })
  http.HandleFunc("/user",func(w http.ResponseWriter, r *http.Request) {
    users.UserHandler(w,r,userDb)  
  })
  http.HandleFunc("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
    users.UserNameHandler(w,r,userDb)
  })

  fmt.Printf("Starting server on 8080 \n")
  httpErr := http.ListenAndServe(":8080",nil)
  
  if httpErr != nil {
    log.Fatal(httpErr)
  }
}
