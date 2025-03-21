package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	user "github.com/alewilliam789/go-rest/users"
)


func createUser() error {

  jimmy :=user.User {
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
    return marshalErr
  }

  resp, reqErr := http.Post("http://localhost:8080/v1/user","application/json", bytes.NewBuffer(jimmyBytes))

  if reqErr != nil {
    return reqErr
  }

  if resp.StatusCode != 201 {
    return errors.New("User was not created")
  }

  return nil
}

func deleteUser() error {

  client := &http.Client{}

  req, reqErr := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/user/jimmy789",nil)
  
  if reqErr != nil {
    return reqErr
  }
  
  resp, respErr := client.Do(req)
  
  if respErr != nil {
    return respErr
  }
  
  if resp.StatusCode != http.StatusOK {
    return errors.New("User not deleted")
  }

  return nil
}

var testCookie *http.Cookie;


func TestRightCredentialsNoCookieLogin(t *testing.T){
  
  createErr := createUser()

  defer deleteUser()

  if createErr != nil {
    t.Fatal(createErr)
  }
  
  newLogin := LoginAttempt {
    UserName: "jimmy789",
    Password: "BOO",
  }

  loginBytes, marshalErr := json.Marshal(newLogin)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }
  
  resp, reqErr := http.Post("http://localhost:8080/v1/login","application/json", bytes.NewBuffer(loginBytes))

  if reqErr != nil {
    t.Fatal(reqErr)
  }

  fmt.Print(resp.StatusCode)

  if resp.StatusCode != 200 {
    t.Fatal("Login with correct credentials was not correct") 
  }

  cookies := resp.Cookies()

  correctName := false

  for _, cookie := range(cookies) { 
    if cookie.Name == "user-cookie" {
      testCookie = cookie
      correctName = true
      fmt.Print(cookie.Value)
    }
  }

  if !correctName {
    t.Fatal("Cookie improperly set")
  }
}


func TestRightCredentialsCookieLogin(t *testing.T){

  
  newLogin := LoginAttempt {
    UserName: "jimmy789",
    Password: "BOO",
  }

  loginBytes, marshalErr := json.Marshal(newLogin)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }

  client := &http.Client{}
  
  req, reqErr := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/login",bytes.NewBuffer(loginBytes))
  
  if reqErr != nil {
    t.Fatal(reqErr)
  }

  if testCookie != nil {
    t.Fatal("Test cookie is not set")
  }

  req.AddCookie(testCookie)

  resp, respErr := client.Do(req)

  if respErr != nil {
    t.Fatal(respErr)
  }

  if resp.StatusCode != 200 {
    t.Fatal("Login with correct credentials was not correct") 
  }

  cookies := resp.Cookies()

  for _, cookie := range(cookies) { 
    if cookie.Name == "user-cookie" && (testCookie.Value != cookie.Value){
      t.Fatal("Non-matching cookie")
    }
  }
}

func TestWrongCredentialsNoCookieLogin(t *testing.T){
  
  newLogin := LoginAttempt {
    UserName: "jimmy789",
    Password: "GHOST",
  }

  loginBytes, marshalErr := json.Marshal(newLogin)

  if marshalErr != nil {
    t.Fatal(marshalErr)
  }
  
  resp, reqErr := http.Post("http://localhost:8080/v1/login","application/json", bytes.NewBuffer(loginBytes))

  if reqErr != nil {
    t.Fatal(reqErr)
  }

  if resp.StatusCode != 401 {
    t.Fatal("Login with incorrect credentials worked!") 
  }
}




