package users

import (
	"fmt"
	"net/http"
	//"time"
  "encoding/json"
)



type User struct {
  UserName string `json:"username"`
  PassWord string `json:"pass"`
  FirstName string `json:"firstname"`
  LastName string `json:"lastname"`
  DOB string `json:"dob"`
  City string `json:"city"`
  State string `json:"state"`
}

func decode_user(user *User, req *http.Request) error {
  decoder := json.NewDecoder(req.Body)

  err := decoder.Decode(user)
  
  return err
}

func create_user(w http.ResponseWriter, req *http.Request) {

  var new_user User;

  err := decode_user(&new_user, req)

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
  }

  // Check User in database
  // userlist = db.select(user)
  // if len(userlist) == 0 db.insert(user)
  // else return http.StatusDuplicated
  
  fmt.Printf("The user: %s was created",new_user.UserName)
}

func update_user(w http.ResponseWriter, req *http.Request) {
  var current_user User;
  user_id := req.PathValue("id")

  err := decode_user(&current_user, req)

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
  }



  // Check User is in database
  // userlist, err = db.select(user_id)
  // if len(userlist) == 0 return http.StatusNotFound
  // else update users (user_id, user_name, password, first_name, last_name, dob, city, state) VALUES()
  // return status ok and the new user in the body

  fmt.Printf("The user: %s was updated", current_user.UserName)
}

func get_user(w http.ResponseWriter, req *http.Request) {
  user_id := req.PathValue("id")

  // Check user in the database
  // userlist = db.select(user_id)
  // if len(userlist) == 0 return http.StatusNotFound
  // else w.Write(userlist[0])

  fmt.Printf("The user with id %s was gott", user_id)
}





func UserHandler(w http.ResponseWriter, req *http.Request) {
  
  switch req.Method {
    case "POST":
      create_user(w,req)
    case "PUT":
      update_user(w,req)
  }
}

func UserIdHandler(w http.ResponseWriter, req *http.Request) {
  switch req.Method {
    case "GET":
      get_user(w,req)
  }
}







