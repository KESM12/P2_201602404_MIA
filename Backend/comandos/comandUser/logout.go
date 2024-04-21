package comandUser

import (
	"P1/logger"
	"fmt"
)

type Logout struct {
}

func (l *Logout) Execute(parametros []string) {
	if l.Logout() {
		fmt.Println("Se cerro la sesión.")
	} else {
		fmt.Println("Cierre de sesión incorrecto.")
	}
}

func (l *Logout) Logout() bool {
	return logger.Log.Logout()
}
