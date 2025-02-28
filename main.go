package main

import (
	"fmt"
	"net/http"

	auth "github.com/alewilliam789/go-rest/auth"
  users "github.com/alewilliam789/go-rest/users"
)




func main() {

  http.HandleFunc("/login", auth.Login)
  http.HandleFunc("/user",users.UserHandler)
  http.HandleFunc("/user/{id}", users.UserIdHandler)

  fmt.Printf("Starting server on 8080 \n")
  http.ListenAndServe(":8080",nil)
}
