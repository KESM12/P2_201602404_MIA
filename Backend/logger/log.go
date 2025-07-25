package logger

import (
	"P1/estructuras"
	"fmt"
)

type User struct {
	Grupo, User, Pass [10]byte
	Id                string
}

func (u *User) GetName() [10]byte {
	return u.User
}

func (u *User) GetId() string {
	return u.Id
}

type Logger struct {
	LoggedIn bool
	Usr      *User
}

func (l *Logger) Login(usr *User) bool {
	if !l.IsLoggedIn() {
		l.Usr = usr
		return true
	}
	fmt.Println("ya existe un usuario registrado")
	return false
}

func (l *Logger) IsLoggedIn() bool {
	return l.Usr != nil
}

func (l *Logger) Logout() bool {
	if l.IsLoggedIn() {
		l.Usr = nil
		return true
	}
	fmt.Println("no hay usuario loggeado")
	return false
}

func (l *Logger) UserIsRoot() bool {
	return estructuras.Equal(l.Usr.GetName(), "root")
}

func (l *Logger) GetUserName() [10]byte {
	return l.Usr.GetName()
}

func (l *Logger) GetUserId() string {
	return l.Usr.GetId()
}

var Log = Logger{
	LoggedIn: false,
	Usr:      nil,
}
