package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}


//database connection pool
type UserModel struct{
  DB *sql.DB
}


func (u *UserModel) Insert(name, email, password string) error {
  return nil
}


func (m *UserModel) Authenticate(email, password string)(int, error) {
  return 0, nil
}


func (m *UserModel) Exists(id int) (bool, error) {
  return false, nil
}
