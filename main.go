package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/alewilliam789/go-rest/auth"
	users "github.com/alewilliam789/go-rest/users"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func setupConn() *sql.DB {
  cfg := mysql.Config{
    User: os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net: "tcp",
    Addr: os.Getenv("ADDR"),
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


func main() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal(err)
  }


  userDb := setupConn()
  defer userDb.Close()

  http.HandleFunc("/login", func(w http.ResponseWriter ,r *http.Request) {
    auth.Login(w, r)
  })
  http.HandleFunc("/user",func(w http.ResponseWriter, r *http.Request) {
    users.UserHandler(w,r,userDb)  
  })
  http.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
    users.UserIdHandler(w,r,userDb)
  })

  fmt.Printf("Starting server on 8080 \n")
  httpErr := http.ListenAndServe(":8080",nil)
  
  if httpErr != nil {
    log.Fatal(httpErr)
  }
}
