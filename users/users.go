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
  Id int32 `json:"id"`
  UserName string `json:"username"`
  PassWord []byte `json:"password, omitempty"`
  FirstName string `json:"firstname"`
  LastName string `json:"lastname"`
  DOB string `json:"dob"`
  City string `json:"city"`
  State string `json:"state"`
}

func (user *User) FromDB(newUser *usersSql.User) {
  user.Id = newUser.UserID
  user.UserName = newUser.UserName
  user.FirstName = newUser.FirstName.String
  user.LastName = newUser.LastName.String
  user.DOB = newUser.Dob.String
  user.City = newUser.City.String
  user.State = newUser.State.String
}

func (user *User) DeepEqual(otherUser *User) bool {
  if user.UserName == otherUser.UserName && user.FirstName == otherUser.FirstName && user.DOB == otherUser.DOB && user.City == otherUser.City && user.State == otherUser.State {
    return true
  }

  return false
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

  user.PassWord = hasher.Sum(nil)

  return nil
}

func checkUserExists(user *User,ctx context.Context, queries *usersSql.Queries) (bool, error){
  foundUser, dbErr := queries.GetUser(ctx,user.UserName)

  userFound := true

  var checkUser User

  checkUser.FromDB(&foundUser)

  if dbErr == sql.ErrNoRows {
    userFound = false
    return userFound, nil
  }
  
  return userFound, dbErr
}

func createUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {

  var newUser User;

  ctx := context.Background()
  defer ctx.Done()

  jsonErr := decodeUser(&newUser, req)


  if jsonErr != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Print(jsonErr)
    return
  }

  hashErr := hashPass(&newUser)

  if hashErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(hashErr)
    return 
  }

  doesExist, checkErr := checkUserExists(&newUser, ctx, queries)
  
  if checkErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(checkErr)
    return
  }

  if doesExist {
    w.WriteHeader(http.StatusConflict)
    log.Printf("The user with username %s already exists\n", newUser.UserName)
    return
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
    log.Print(sqlErr)
    return
  }

  new_id, idErr := result.LastInsertId()

  if idErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(idErr)
    return
  }

  newUser.Id = int32(new_id)
  newUser.PassWord = []byte("")

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type","application/json")
  encodeErr := encodeUser(&newUser,w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(encodeErr)
    return
  }

  log.Printf("The user: %s was created\n",newUser.UserName)
}

func updateUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  var currentUser User;
  ctx := context.Background()
  defer ctx.Done()

  jsonErr := decodeUser(&currentUser, req)

  if jsonErr != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Print(jsonErr)
    return
  }

  hashErr := hashPass(&currentUser)

  if hashErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(hashErr)
    return
  }

  doesExist, checkErr := checkUserExists(&currentUser, ctx, queries)

  if checkErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(checkErr)
    return
  }

  if !doesExist {
    w.WriteHeader(http.StatusNotFound)
    log.Printf("The user with username %s does not exist\n", currentUser.UserName)
    return
  }

  updateUserParams := usersSql.UpdateUserParams{
    UserName: currentUser.UserName,
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
    log.Print(sqlErr)
    return
  }

  currentUser.PassWord = []byte("")

  w.WriteHeader(http.StatusAccepted)
  w.Header().Set("Content-Type","application/json")
  encodeErr := encodeUser(&currentUser,w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(encodeErr)
    return
  }

  fmt.Printf("The user: %s was updated\n", currentUser.UserName)
}


func getUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  username := req.PathValue("username")
  
  ctx := context.Background()
  defer ctx.Done()

  
  foundUser, dbErr := queries.GetUser(ctx,username)

  if dbErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(dbErr)
    return
  }

  var currentUser User

  currentUser.FromDB(&foundUser)

  currentUser.PassWord = []byte("")

  w.WriteHeader(http.StatusFound)
  w.Header().Set("Content-Type", "application/json")
  encodeErr := encodeUser(&currentUser, w)

  if encodeErr != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Print(encodeErr)
    return
  }

  fmt.Printf("The user with username %s was got\n", username)
}  

func deleteUser(w http.ResponseWriter, req *http.Request, queries *usersSql.Queries) {
  username := req.PathValue("username")
  
  ctx := context.Background()

  sqlErr := queries.DeleteUser(ctx, username)
  
  if sqlErr != nil {
    w.WriteHeader(http.StatusAccepted)
    log.Print(sqlErr)
    return
  }
  
  fmt.Printf("The user with username %s was deleted\n",username)
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

func UserNameHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {

  queries := usersSql.New(db)

  switch req.Method {
    case "GET":
      getUser(w,req, queries)
    case "DELETE":
      deleteUser(w,req,queries)
  }
}







