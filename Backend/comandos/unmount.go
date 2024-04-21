package comandos

import (
	"P1/lista"
	"fmt"
	"strings"
)

type ParametrosUnmount struct {
	Path string
	ID   string
}

type Unmount struct {
	Params ParametrosUnmount
}

func (u *Unmount) Execute(parametros []string) {
	u.Params = u.SaveParams(parametros)
	if u.Unmount(u.Params.Path, u.Params.ID) {
		fmt.Printf("Particion con id '%s' desmontada con exito\n\n", u.Params.ID)
	} else {
		fmt.Printf("No se logro desmontar la particion con id '%s'\n", u.Params.ID)
	}
}

func (u *Unmount) SaveParams(parametros []string) ParametrosUnmount {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "driveletter") {
			v = strings.ReplaceAll(v, "driveletter=", "")
			v = strings.ReplaceAll(v, "\"", "")
			v = v + ".dsk"
			u.Params.Path = v
		} else if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			u.Params.ID = v
		}
	}
	return u.Params
}

func (m *Unmount) Unmount(path string, id string) bool {
	lista.ListaMount.PrintList()
	if path == "" {
		fmt.Printf("Falta el parametro name.")
		return false
	}
	if id == "" {
		fmt.Printf("Falta el parametro id.\n")
		return false
	}
	unmountptr := lista.ListaMount.UnMount(id)
	if unmountptr == nil {
		fmt.Println(fmt.Printf("No existe una partición montada con el ID '%s'\n", string(id)))
		return false
	}
	lista.ListaMount.PrintList()
	return true
}
