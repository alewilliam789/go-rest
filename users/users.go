package users

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	usersSql "github.com/alewilliam789/go-rest/db"
)



type User struct {
  id int32 `json:"id,omitempty"`
  UserName string `json:"username"`
  PassWord string `json:"pass"`
  FirstName string `json:"firstname"`
  LastName string `json:"lastname"`
  DOB string `json:"dob"`
  City string `json:"city"`
  State string `json:"state"`
}

func decodeUser(user *User, req *http.Request) error {
  decoder := json.NewDecoder(req.Body)

  err := decoder.Decode(user)
  
  return err
}

func encodeUser(user *User, w http.ResponseWriter) error {
  encoder := json.NewEncoder(w)

  err := encoder.Encode(user)

  return err
}

func hashPass(user *User) error {
  hasher := sha256.New()

  _, err := hasher.Write([]byte(user.PassWord))

  if err != nil {
    return err
  }

  user.PassWord = string(hasher.Sum(nil))

  return nil
}

func createUser(w http.ResponseWriter, req *http.Request, ctx context.Context, queries *usersSql.Queries) {

  var newUser User;

  jsonErr := decodeUser(&newUser, req)

  if jsonErr != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Panic(jsonErr)
  }

  hashErr := hashPass(&newUser)

  if hashErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Panic(hashErr)
  }

  newUserParams := usersSql.CreateUserParams {
    Username: newUser.UserName,
    Password: newUser.PassWord,
    Firstname: sql.NullString{String: newUser.FirstName,Valid:true},
    Lastname: sql.NullString{String: newUser.LastName, Valid:true},
    Dob: sql.NullString{String: newUser.DOB, Valid:true},
    City: sql.NullString{String: newUser.City, Valid:true},
    State: sql.NullString{String: newUser.State, Valid:true},
  }

  _, sqlErr := queries.CreateUser(ctx, newUserParams)

  
  if sqlErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(sqlErr)
  }

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type","application/json")
  encodeErr := encodeUser(&newUser,w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(encodeErr)
  }

  fmt.Printf("The user: %s was created",newUser.UserName)
}

func updateUser(w http.ResponseWriter, req *http.Request, ctx context.Context, queries *usersSql.Queries) {
  var currentUser User;

  jsonErr := decodeUser(&currentUser, req)

  if jsonErr != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Fatal(jsonErr)
  }

  hashErr := hashPass(&currentUser)

  if hashErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(hashErr)
  }

  updateUserParams := usersSql.UpdateUserParams{
    ID: currentUser.id,
    Password: currentUser.PassWord,
    Firstname: sql.NullString{String: currentUser.FirstName, Valid:true},
    Lastname: sql.NullString{String: currentUser.LastName, Valid:true},
    Dob: sql.NullString{String: currentUser.DOB, Valid: true},
    City: sql.NullString{String: currentUser.City, Valid: true},
    State: sql.NullString{String: currentUser.State, Valid: true},
  }

  sqlErr := queries.UpdateUser(ctx,updateUserParams)

  if sqlErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(sqlErr)
  }

  w.WriteHeader(http.StatusAccepted)
  w.Header().Set("Content-Type","application/json")
  encodeErr := encodeUser(&currentUser,w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(encodeErr)
  }

  fmt.Printf("The user: %s was updated", currentUser.UserName)
}

func getUser(w http.ResponseWriter, req *http.Request, ctx context.Context, queries *usersSql.Queries) {
  userId := req.PathValue("id")

  // Check user in the database
  // userlist = db.select(user_id)
  // if len(userlist) == 0 return http.StatusNotFound
  // else w.Write(userlist[0])

  fmt.Printf("The user with id %s was gott", userId)
}





func UserHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {
  
  ctx := context.Background()

  queries := usersSql.New(db)

  switch req.Method {
    case "POST":
      createUser(w,req, ctx, queries)
    case "PUT":
      updateUser(w,req, ctx, queries)
  }
}

func UserIdHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {

  ctx := context.Background()

  queries := usersSql.New(db)

  switch req.Method {
    case "GET":
      getUser(w,req, ctx, queries)
  }
}







