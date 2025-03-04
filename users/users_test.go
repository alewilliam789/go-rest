package users

import (
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
    PassWord: "BOO",
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

  var respUser User

  decodeErr := decodeResp(&respUser,resp)

  if decodeErr != nil {
    t.Fatal(decodeErr)
  }

  jimmy.Id = respUser.Id

  if jimmy != respUser {
    t.Fatal("Response doesn't match")
  }
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
