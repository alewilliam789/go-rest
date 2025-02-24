package users

import (
	"net/http"
	"time"
)

type User struct {
  UserName string `json:"user_name"`
  PassWord string `json:"pass"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  City string `json:"city"`
  State string `json:"state"`
}





