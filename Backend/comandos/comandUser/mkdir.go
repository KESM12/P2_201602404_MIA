package comandUser

import (
	"P1/comandos"
	"P1/estructuras"
	"P1/inodos"
	"P1/lista"
	"P1/logger"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

type ParametrosMkdir struct {
	Path string
	R    bool
}

type Mkdir struct {
	Params ParametrosMkdir
}

func (m *Mkdir) Execute(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkdir(m.Params.Path, m.Params.R) {
		fmt.Printf("\nLa carpeta %s creada correctamente\n\n", m.Params.Path)
	} else {
		fmt.Printf("")
	}
}

func (m *Mkdir) SaveParams(parametros []string) ParametrosMkdir {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "path") {
			v = strings.ReplaceAll(v, "path=", "")
			m.Params.Path = v
		} else if v == "r" {
			// v = strings.ReplaceAll(v, "r", "")
			m.Params.R = true
		}
	}
	return m.Params
}

func (m *Mkdir) Mkdir(path string, r bool) bool {
	if path == "" {
		fmt.Println("Ruta no encontrada.")
		return false
	}
	path = strings.Replace(path, "/", "", 1)
	if !logger.Log.IsLoggedIn() {
		fmt.Println("Inicie sesión primero.")
		return false
	}

	if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
		return createDirectory(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, r)
	} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL != nil {
		return createDirectory(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(estructuras.EBR{})), r)
	}
	return false
}

func createDirectory(name [10]byte, path, ruta string, whereToStart int64, r bool) bool {
	var superbloque estructuras.SuperBloque
	comandos.Fread(&superbloque, path, whereToStart)
	var tablaInodoRoot estructuras.TablaInodo
	comandos.Fread(&tablaInodoRoot, path, superbloque.S_inode_start)
	var tablaInodoUsers estructuras.TablaInodo
	comandos.Fread(&tablaInodoUsers, path, superbloque.S_inode_start+superbloque.S_inode_size)
	contenido := ReadFile(&tablaInodoUsers, path, &superbloque)
	userId := GetUserId(contenido, string(TrimArray(name[:])))
	groupdId := GetGroupId(contenido, string(TrimArray(name[:])))
	if r {
		FindAndCreateDirectories(&tablaInodoRoot, path, ruta, &superbloque, 0, userId, groupdId)
	}
	num := NewInodeDirectory(&superbloque, path, userId, groupdId)
	FindDirs(num, &tablaInodoRoot, path, ruta, &superbloque, 0)
	comandos.Fwrite(&tablaInodoRoot, path, superbloque.S_inode_start)
	comandos.Fwrite(&superbloque, path, whereToStart)
	return true
}

func NewInodeDirectory(superbloque *estructuras.SuperBloque, path string, userId, groupId int64) int64 {
	var nuevaTabla estructuras.TablaInodo
	posicionActual := inodos.WriteInBitmapInode(path, superbloque)
	// aqui llenaremos la nueva tabla de inodos
	nuevaTabla.I_uid = userId
	nuevaTabla.I_gid = groupId
	nuevaTabla.I_size = 0
	nuevaTabla.I_type = '0'
	nuevaTabla.I_perm = 664
	// llenando las fechas
	atime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_atime[i] = atime.String()[i]
	}
	ctime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_ctime[i] = ctime.String()[i]
	}
	mtime := time.Now()
	for i := 0; i < len(nuevaTabla.I_atime)-1; i++ {
		nuevaTabla.I_mtime[i] = mtime.String()[i]
	}
	// llenando a todos los bloques no utilizados
	for i := 0; i < len(nuevaTabla.I_block); i++ {
		nuevaTabla.I_block[i] = -1
	}

	// aqui escribiremos el contenido y crearemos el nuevo bloque de carpeta
	posicionNuevoBloqueCarpetas := inodos.WriteInBitmapBlock(path, superbloque)
	nuevaTabla.I_block[0] = posicionNuevoBloqueCarpetas

	nuevoBloqueCarpetas := estructuras.BloqueDeCarpetas{}

	// llenando la carpeta
	copy(nuevoBloqueCarpetas.B_content[0].B_name[:], ".")
	nuevoBloqueCarpetas.B_content[0].B_inodo = int32(posicionActual)

	copy(nuevoBloqueCarpetas.B_content[1].B_name[:], "..")
	nuevoBloqueCarpetas.B_content[1].B_inodo = -1

	copy(nuevoBloqueCarpetas.B_content[2].B_name[:], "")
	nuevoBloqueCarpetas.B_content[2].B_inodo = -1

	copy(nuevoBloqueCarpetas.B_content[3].B_name[:], "")
	nuevoBloqueCarpetas.B_content[3].B_inodo = -1
	// escribiendo la nueva tabla de inodos
	comandos.Fwrite(&nuevaTabla, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
	comandos.Fwrite(&nuevoBloqueCarpetas, path, superbloque.S_block_start+posicionNuevoBloqueCarpetas*superbloque.S_block_size)
	return posicionActual
}

func GetContent(cont string) string {
	file, err := os.Open(cont)
	if err != nil {
		return ""
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return content
}
