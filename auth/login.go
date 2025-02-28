package login

import (
  "net/http"
  "encoding/json"
  "fmt"
)


type LoginAttempt struct {
  User string `json:"User"`
  Pass string `json:"Pass"`
  // token String
}

func Login(w http.ResponseWriter, req *http.Request) {
    if (req.Method != "POST") {
      w.WriteHeader(http.StatusBadRequest)
    }

    var la LoginAttempt
    decoder := json.NewDecoder(req.Body)

    err := decoder.Decode(&la)

    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      panic(err)
    }

    fmt.Printf("Hello %s \n",la.User)

    w.WriteHeader(http.StatusOK) 
}
