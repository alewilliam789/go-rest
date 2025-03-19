package users

import (
  "bytes"
	"encoding/json"
	"net/http"
	"testing"
)


func decodeResp(user *User, resp *http.Response) error {
  
  decoder := json.NewDecoder(resp.Body)
  
  err := decoder.Decode(user)

  return err
}

func TestCreateUser(t *testing.T) {

  jimmy := User {
    UserName: "jimmy789",
    PassWord: []byte("BOO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1197",
    City: "Scranton",
    State: "PA",
  }

  jimmyBytes, marshalErr := json.Marshal(jimmy)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }

  resp, reqErr := http.Post("http://localhost:8080/v1/user","application/json", bytes.NewBuffer(jimmyBytes))
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }

  hashErr := hashPass(&jimmy)

  if hashErr != nil {
    t.Fatal(hashErr)
  }

  var respUser User

  decodeErr := decodeResp(&respUser,resp)

  if decodeErr != nil {
    t.Fatal(decodeErr)
  }

  if !jimmy.DeepEqual(&respUser) {
    t.Fatal("Response doesn't match")
  }
}

func TestCreateExistingUser(t *testing.T) {

  jimmy := User {
    UserName: "jimmy789",
    PassWord: []byte("BOO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1197",
    City: "Scranton",
    State: "PA",
  }

  jimmyBytes, marshalErr := json.Marshal(jimmy)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }

  resp, reqErr := http.Post("http://localhost:8080/v1/user","application/json", bytes.NewBuffer(jimmyBytes))
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }

  if resp.StatusCode != 409 {
    t.Fatalf("Expected Status Code 409, received %d",resp.StatusCode)
  }
}

func TestUpdateUser(t *testing.T) {
  jimmy := User {
    UserName: "jimmy789",
    PassWord: []byte("HELLO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1997",
    City: "Los Angeles",
    State: "CA",
  }
  
  jimmyBytes, marshalErr := json.Marshal(jimmy)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }

  client := &http.Client{}
  
  req, reqErr := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/user",bytes.NewBuffer(jimmyBytes))
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }

  resp, respErr := client.Do(req)

  if respErr != nil {
    t.Fatal(respErr)
  }

  hashErr := hashPass(&jimmy)

  if hashErr != nil {
    t.Fatal(hashErr)
  }

  var returned_jimmy User;

  decodeErr := decodeResp(&returned_jimmy, resp)

  if decodeErr != nil {
    t.Fatal(decodeErr)
  }

  if !jimmy.DeepEqual(&returned_jimmy) {
    t.Fatal("The user was not returned correctly")
  }
}

func TestUpdateMissingUser(t *testing.T) {
  jimmy := User {
    UserName: "immy789",
    PassWord: []byte("HELLO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1997",
    City: "Los Angeles",
    State: "CA",
  }
  
  jimmyBytes, marshalErr := json.Marshal(jimmy)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }

  client := &http.Client{}
  
  req, reqErr := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/user",bytes.NewBuffer(jimmyBytes))
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }

  resp, respErr := client.Do(req)

  if respErr != nil {
    t.Fatal(respErr)
  }

  if resp.StatusCode != 404 {
    t.Fatalf("Expected Status Code 404, received %d", resp.StatusCode)
  }
}

func TestGetUser(t *testing.T) {
  jimmy := User {
    UserName: "jimmy789",
    PassWord: []byte("HELLO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1997",
    City: "Los Angeles",
    State: "CA",
  }

  var gotUser User

  resp, respErr := http.Get("http://localhost:8080/user/jimmy789")

  if respErr != nil {
    t.Fatal(respErr) 
  }

  decodeErr := decodeResp(&gotUser, resp)

  if decodeErr != nil {
    t.Fatal(decodeErr)
  }

  hashErr := hashPass(&jimmy)

  if hashErr != nil {
    t.Fatal(hashErr)
  }

  if !jimmy.DeepEqual(&gotUser) {
    t.Fatal("The users don't match")
  }

}

func TestDeleteUser(t *testing.T) {
  client := &http.Client{}

  req, reqErr := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/user/jimmy789",nil)
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }
  
  resp, respErr := client.Do(req)
  
  if respErr != nil {
    t.Fatal(respErr)
  }
  
  if resp.StatusCode != http.StatusOK {
    t.Fatal("User not deleted")
  }
}
