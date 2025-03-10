package users

import (
	// "bytes"
	"bytes"
	"encoding/json"
	"fmt"
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

  resp, reqErr := http.Post("http://localhost:8080/user","application/json", bytes.NewBuffer(jimmyBytes))
  
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

  jimmy.Id = respUser.Id

  if !jimmy.DeepEqual(&respUser) {
    t.Fatal("Response doesn't match")
  }
}

// func TestUpdateUser(t *testing.T) {
//   jimmy := User {
//     UserName: "jimmy789",
//     PassWord: []byte("HELLO"),
//     FirstName: "Jimmy",
//     LastName: "Halpert",
//     DOB: "02/23/1997",
//     City: "Los Angeles",
//     State: "CA",
//     Id: 2,
//   }
//   
//   jimmyBytes, marshalErr := json.Marshal(jimmy)
//
//   if marshalErr != nil {
//     t.Fatal(marshalErr)
//   }
//
//   client := &http.Client{}
//   
//   req, reqErr := http.NewRequest(http.MethodPut, "http://localhost:8080/user",bytes.NewBuffer(jimmyBytes))
//   
//   if reqErr != nil {
//     t.Fatal(reqErr)
//   }
//
//   resp, respErr := client.Do(req)
//
//   fmt.Print(resp.Status)
//
//   if respErr != nil {
//     t.Fatal(respErr)
//   }
//
//   var returned_jimmy User;
//
//   decodeErr := decodeResp(&returned_jimmy, resp)
//
//   if decodeErr != nil {
//     t.Fatal(decodeErr)
//   }
//
//   if !jimmy.DeepEqual(&returned_jimmy) {
//     t.Fatal("The user was not returned correctly")
//   }
//
//   fmt.Print(returned_jimmy)
// }

func TestGetUser(t *testing.T) {
  jimmy := User {
    UserName: "jimmy789",
    PassWord: []byte("BOO"),
    FirstName: "Jimmy",
    LastName: "Halpert",
    DOB: "02/23/1197",
    City: "Scranton",
    State: "PA",
  }

  var gotUser User

  resp, respErr := http.Get("http://localhost:8080/user/")

  if respErr != nil {
    t.Fatal(respErr) 
  }

  decodeErr := decodeResp(&gotUser, resp)

  if decodeErr != nil {
    t.Fatal(decodeErr)
  }

  if !jimmy.DeepEqual(&gotUser) {
    t.Fatal("The users don't match")
  }

  fmt.Print(gotUser)
}


func TestDeleteUser(t *testing.T) {
  client := &http.Client{}

  req, reqErr := http.NewRequest(http.MethodDelete, "http://localhost:8080/user/1",nil)
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }
  
  resp, respErr := client.Do(req)
  
  if respErr != nil {
    t.Fatal(respErr)
  }
  
  fmt.Println("Response Status:", resp.Status)

  if resp.StatusCode != http.StatusOK {
    t.Fatal("User not deleted")
  }
}
