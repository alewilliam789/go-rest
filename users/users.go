package users

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "strconv"

	usersSql "github.com/alewilliam789/go-rest/db"
)



type User struct {
  Id int32 `json:"id,omitempty"`
  UserName string `json:"username"`
  PassWord string `json:"password"`
  FirstName string `json:"firstname"`
  LastName string `json:"lastname"`
  DOB string `json:"dob"`
  City string `json:"city"`
  State string `json:"state"`
}

func (user *User) FromDB(newUser *usersSql.User) {
  user.Id = newUser.UserID
  user.UserName = newUser.UserName
  user.PassWord = newUser.Password
  user.FirstName = newUser.FirstName.String
  user.LastName = newUser.LastName.String
  user.DOB = newUser.Dob.String
  user.City = newUser.City.String
  user.State = newUser.State.String
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

func createUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {

  var newUser User;

  ctx := context.Background()
  defer ctx.Done()

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
    UserName: newUser.UserName,
    Password: newUser.PassWord,
    FirstName: sql.NullString{String: newUser.FirstName,Valid:true},
    LastName: sql.NullString{String: newUser.LastName, Valid:true},
    Dob: sql.NullString{String: newUser.DOB, Valid:true},
    City: sql.NullString{String: newUser.City, Valid:true},
    State: sql.NullString{String: newUser.State, Valid:true},
  }

  result, sqlErr := queries.CreateUser(ctx, newUserParams)

  
  if sqlErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(sqlErr)
  }

  new_id, idErr := result.LastInsertId()

  if idErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(idErr)
  }

  newUser.Id = int32(new_id)

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type","application/json")
  encodeErr := encodeUser(&newUser,w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(encodeErr)
  }

  fmt.Printf("The user: %s was created",newUser.UserName)
}

func updateUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  var currentUser User;
  ctx := context.Background()
  defer ctx.Done()

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
    UserID: currentUser.Id,
    Password: currentUser.PassWord,
    FirstName: sql.NullString{String: currentUser.FirstName, Valid:true},
    LastName: sql.NullString{String: currentUser.LastName, Valid:true},
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

func getUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  userId, convErr := strconv.ParseInt(req.PathValue("id"),10,32)
  
  ctx := context.Background()
  defer ctx.Done()

  if convErr != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Fatal(convErr)
  }
  
  foundUser, dbErr := queries.GetUser(ctx,int32(userId))

  if dbErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(dbErr)
  }

  var currentUser User

  currentUser.FromDB(&foundUser)

  w.WriteHeader(http.StatusFound)
  w.Header().Set("Content-Type", "application/json")
  encodeErr := encodeUser(&currentUser, w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(encodeErr)
  }

  fmt.Printf("The user with id %d was gott", userId)
}

func deleteUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  userId, parseErr := strconv.ParseInt(req.PathValue("id"),10,32)
  
  ctx := context.Background()

  if parseErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Fatal(parseErr)
  }

  sqlErr := queries.DeleteUser(ctx, int32(userId))
  
  if sqlErr != nil {
    w.WriteHeader(http.StatusAccepted)
    log.Fatal(sqlErr)
  }
  
  fmt.Printf("The user with id %d was deleted",userId)
}

func UserHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {
  
  queries := usersSql.New(db)

  switch req.Method {
    case "POST":
      createUser(w, req, queries)
    case "PUT":
      updateUser(w, req, queries)
  }
}

func UserIdHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {

  queries := usersSql.New(db)

  switch req.Method {
    case "GET":
      getUser(w,req, queries)
    case "DELETE":
      deleteUser(w,req,queries)
  }
}







