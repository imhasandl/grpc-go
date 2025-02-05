package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}

type App struct {
	ID     int64
	Name   string
	Secret string 
}
