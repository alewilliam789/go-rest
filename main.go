package main

import (
	"fmt"
	"net/http"

	auth "github.com/alewilliam789/go-rest/auth"
)




func main() {

  http.HandleFunc("/login", auth.Login)

  fmt.Printf("Starting server on 8080 \n")
  http.ListenAndServe(":8080",nil)
}
