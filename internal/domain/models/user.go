package models

type User struct {
	ID       int64
	IsAdmin  bool
	Email    string
	PassHash []byte
}
