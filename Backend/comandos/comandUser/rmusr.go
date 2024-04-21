package comandUser

import (
	"P1/comandos"
	"P1/estructuras"
	"P1/lista"
	"P1/logger"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

type ParametrosRmusr struct {
	User string
}

type Rmusr struct {
	params ParametrosRmusr
}

func (r *Rmusr) Execute(parametros []string) {
	r.params = r.SaveParams(parametros)
	if r.Rmusr(r.params.User) {
		fmt.Printf("Usuario \"%s\" eliminado con éxito\n\n", r.params.User)
	} else {
		fmt.Printf("")
	}
}

func (r *Rmusr) SaveParams(parametros []string) ParametrosRmusr {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			r.params.User = v
		}
	}
	return r.params
}

func (r *Rmusr) Rmusr(user string) bool {
	if user == "" {
		fmt.Println("Usuario no encontrado.")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmusrPartition(user, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return r.RmusrPartition(user, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(estructuras.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (r *Rmusr) RmusrPartition(user string, whereToStart int64, path string) bool {
	// superbloque de la particion
	var superbloque estructuras.SuperBloque
	comandos.Fread(&superbloque, path, whereToStart)

	// tabla de inodos de archivo Users.txt
	var tablaInodo estructuras.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(estructuras.TablaInodo{})))
	// modificar la fecha en la que se esta modificando el inodo
	mtime := time.Now()
	for i := 0; i < len(tablaInodo.I_mtime); i++ {
		tablaInodo.I_mtime[i] = mtime.String()[i]
	}
	if !r.ExisteUsuario(ReadFile(&tablaInodo, path, &superbloque), user) {
		fmt.Printf("Usuario \"%s\" no encontrado.\n", user)
		return false
	}
	contenido := modFile(&tablaInodo, path, &superbloque)
	nuevoContenido := r.DesactivarUsuario(contenido, user)
	// fmt.Println(nuevoContenido)
	if SetFile(&tablaInodo, path, &superbloque, nuevoContenido) {
		fmt.Println(ReadFile(&tablaInodo, path, &superbloque))
		return true
	}
	return false
}

func (r *Rmusr) DesactivarUsuario(contenido string, userName string) string {
	nuevoContenido := ""
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] == "U" {
			if parametros[3] == userName {
				parametros[0] = "0"
			}
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "," + parametros[3] + "," + parametros[4] + "\n"
		} else if parametros[1] == "G" {
			nuevoContenido += parametros[0] + "," + parametros[1] + "," + parametros[2] + "\n"
		}
	}
	return nuevoContenido
}

func (r *Rmusr) ExisteUsuario(contenido string, userName string) bool {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		if strings.Compare(parametros[3], userName) == 0 {
			return true
		}
	}
	return false
}
