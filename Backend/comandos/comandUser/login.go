package comandUser

import (
	"P1/comandos"
	"P1/estructuras"
	"P1/lista"
	"P1/logger"
	"bytes"
	"fmt"
	"strings"
	"unsafe"
)

type ParametrosLogin struct {
	User [10]byte
	Pwd  [10]byte
	Id   string
}

type Login struct {
	Params ParametrosLogin
}

func (l *Login) Execute(parametros []string) {
	l.Params = l.SaveParams(parametros)
	if l.Login(l.Params.User, l.Params.Pwd, l.Params.Id) {
		usuario := estructuras.TrimArray(l.Params.User[:])
		fmt.Printf("Usuario \"%s\" loggeado con éxito.\n\n", usuario)
	} else {
		usuario := estructuras.TrimArray(l.Params.User[:])
		fmt.Printf("el usuario \"%s\" no pudo iniciar sesión.\n", usuario)
	}
}

func (l *Login) SaveParams(parametros []string) ParametrosLogin {
	for _, v := range parametros {
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "user") {
			v = strings.ReplaceAll(v, "user=", "")
			copy(l.Params.User[:], v[:])
		} else if strings.Contains(v, "pass") {
			v = strings.ReplaceAll(v, "pass=", "")
			copy(l.Params.Pwd[:], v[:])
		} else if strings.Contains(v, "id") {
			v = strings.ReplaceAll(v, "id=", "")
			l.Params.Id = v
		}
	}
	return l.Params
}

func (l *Login) Login(User [10]byte, Pwd [10]byte, Id string) bool {
	if bytes.Equal(User[:], []byte("")) {
		fmt.Println("Usuario no reconocido.")
		return false
	}
	if bytes.Equal(Pwd[:], []byte("")) {
		fmt.Println("Ingrese la contraseña por favor.")
		return false
	}
	if Id == "" {
		fmt.Println("No se encuentra el id de la particion montada")
		return false
	}

	node := lista.ListaMount.GetNodeById(Id)
	if node == nil {
		fmt.Printf("Id: %s no coincide con ninguna partición montada.\n", Id)
		return false
	}
	if node.Value != nil {
		return l.LoginInPrimaryPartition(node.Ruta, User, Pwd, Id, node.Value)
	} else if node.ValueL != nil {
		return l.LoginInLogicPartition(node.Ruta, User, Pwd, Id, node.ValueL)
	} else {
		// no deberia de entrar aqui nunca
		fmt.Println("No hay particion montada")
	}
	return false
}

func (l *Login) LoginInPrimaryPartition(path string, User [10]byte, Pwd [10]byte, Id string, partition *estructuras.Partition) bool {
	// leyendo el superbloque
	var superbloque estructuras.SuperBloque
	comandos.Fread(&superbloque, path, partition.Part_start)

	// tabla de inodos del archivo
	var tablaInodo estructuras.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(estructuras.TablaInodo{})))

	// vamos a recorrer la tabla de inodos del archivo Users.txt
	var contenido string
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo estructuras.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(estructuras.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
	}
	// leeremos el archivo por linea que se encuentre dentro del archivo
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		grupo := parametros[2]
		username := parametros[3]
		password := parametros[4]
		if !(string(estructuras.TrimArray(User[:])) == string(estructuras.TrimArray([]byte(username)))) || !(string(estructuras.TrimArray(Pwd[:])) == string(estructuras.TrimArray([]byte(password[:])))) {
			continue
		}
		user := &logger.User{
			User: User,
			Pass: Pwd,
			Id:   Id,
		}
		copy(user.Grupo[:], grupo)
		return logger.Log.Login(user)
	}
	fmt.Print("Usuario no encontrado.\n")
	return false
}
func (l *Login) LoginInLogicPartition(path string, User [10]byte, Pwd [10]byte, Id string, partition *estructuras.EBR) bool {
	// leyendo el superbloque
	var superbloque estructuras.SuperBloque
	comandos.Fread(&superbloque, path, partition.Part_start+int64(unsafe.Sizeof(estructuras.EBR{})))

	// tabla de inodos del archivo
	var tablaInodo estructuras.TablaInodo
	comandos.Fread(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(estructuras.TablaInodo{})))

	// vamos a recorrer la tabla de inodos del archivo Users.Txt
	var contenido string
	for i := 0; i < len(tablaInodo.I_block); i++ {
		// fmt.Println(tablaInodo.I_block[i])
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo estructuras.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*int64(unsafe.Sizeof(estructuras.BloqueDeArchivos{})))
		contenido += string(parteArchivo.B_content[:])
	}
	// leeremos el archivo por linea que se encuentre dentro del archivo
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		grupo := parametros[2]
		username := parametros[3]
		password := parametros[4]
		if !estructuras.Equal(User, username) || !estructuras.Equal(Pwd, password) {
			continue
		}
		user := &logger.User{
			User: User,
			Pass: Pwd,
		}
		copy(user.Grupo[:], grupo)
		return logger.Log.Login(user)
	}
	fmt.Println("Usuario no encontrado.")
	return false
}
