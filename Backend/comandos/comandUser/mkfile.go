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
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type ParametrosMkfile struct {
	Path string
	R    bool
	Size int
	Cont string
}

type Mkfile struct {
	Params ParametrosMkfile
}

func (m *Mkfile) Execute(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkfile(m.Params.Path, m.Params.R, m.Params.Size, m.Params.Cont) {
		fmt.Printf("\nel archivo con ruta %s se creó correctamente\n\n", m.Params.Path)
	} else {
		fmt.Printf("\nel archivo con ruta %s no se pudo crear\n\n", m.Params.Path)

	}
}

func (m *Mkfile) SaveParams(parametros []string) ParametrosMkfile {
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
		} else if strings.Contains(v, "cont") {
			v = strings.ReplaceAll(v, "cont=", "")
			m.Params.Cont = v
		} else if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("hubo un error al convertir a int", err.Error())
			}
			m.Params.Size = num
		}
	}
	return m.Params
}

func (m *Mkfile) Mkfile(path string, r bool, size int, cont string) bool {
	if path == "" {
		fmt.Println("no se encontro una ruta.")
		return false
	}
	path = strings.Replace(path, "/", "", 1)
	if !logger.Log.IsLoggedIn() {
		fmt.Println("no se encuentra un usuario loggeado para crear un archivo.")
		return false
	}
	if size < 0 {
		fmt.Println(("el size no puede ser negativo"))
		return false
	}

	if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
		return createFile(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, r, size, cont)
	} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL != nil {
		return createFile(logger.Log.GetUserName(), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta, path, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(estructuras.EBR{})), r, size, cont)
	}
	return false
}

func createFile(name [10]byte, path, ruta string, whereToStart int64, r bool, size int, cont string) bool {
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
	content := ""
	if cont != "" {
		content = getContent(cont)
	}
	if size != 0 && size > 0 {
		if StrlenBytes([]byte(content)) != 0 && StrlenBytes([]byte(content)) < size {
			contador := 0
			for i := StrlenBytes([]byte(content)); i < size; i++ {
				if contador != 9 {
					contador++
				} else {
					contador = 0
				}
				content += strconv.Itoa(contador)

			}
		} else {
			contador := 0
			for i := 0; i < size; i++ {
				if contador != 9 {
					contador++
				} else {
					contador = 0
				}
				content += strconv.Itoa(contador)
			}
		}
	}
	fmt.Println("************ARCHIVO************")
	fmt.Println("\"" + content + "\"\n")
	num := NewInodeFile(&superbloque, path, userId, groupdId, content)
	FindDirectories(num, &tablaInodoRoot, path, ruta, &superbloque, 0)
	comandos.Fwrite(&tablaInodoRoot, path, superbloque.S_inode_start)
	comandos.Fwrite(&superbloque, path, whereToStart)
	return true
}

func GetGroupId(contenido, name string) int64 {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	groupName := ""
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		if parametros[3] == name {
			groupName = parametros[2]
		}
	}
	if groupName == "" {
		return -1
	}
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "G" {
			continue
		}
		if parametros[2] == groupName {
			num, _ := strconv.Atoi(parametros[0])
			return int64(num)
		}
	}
	return -1
}

func GetUserId(contenido, name string) int64 {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		linea = strings.ReplaceAll(linea, "\x00", "")
		parametros := strings.Split(linea, ",")
		if parametros[1] != "U" {
			continue
		}
		if parametros[3] == name {
			num, _ := strconv.Atoi(parametros[0])
			return int64(num)
		}
	}
	return -1
}

func NewInodeFile(superbloque *estructuras.SuperBloque, path string, userId, groupId int64, contenido string) int64 {
	var nuevaTabla estructuras.TablaInodo
	nuevaPosicion := inodos.WriteInBitmapInode(path, superbloque)
	nuevaTabla.I_uid = userId
	nuevaTabla.I_gid = groupId
	nuevaTabla.I_size = int64(StrlenBytes([]byte(contenido)))
	nuevaTabla.I_type = '1'
	nuevaTabla.I_perm = 664
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
	for i := 0; i < len(nuevaTabla.I_block); i++ {
		nuevaTabla.I_block[i] = -1
	}
	llenarTablaDeInodoDeArchivos(&nuevaTabla, superbloque, path, contenido)
	comandos.Fwrite(&nuevaTabla, path, superbloque.S_inode_start+nuevaPosicion*superbloque.S_inode_size)
	return nuevaPosicion
}

func getContent(cont string) string {
	// aqui hay que leer el archivo y ejecutarlo
	file, err := os.Open(cont)
	if err != nil {
		fmt.Printf("Error al intentar abrir el archivo: %s\n", cont)

		return ""
	}

	defer file.Close()

	// Crear un scanner para luego leer linea por linea el archivo de entrada
	scanner := bufio.NewScanner(file)
	content := ""
	// Leyendo linea por linea
	for scanner.Scan() {
		// obteniendo la linea actual
		content += scanner.Text() + "\n"
	}

	// comprobar que no hubo error al leer el archivo
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error al leer el archivo: %s\n", cont)

		return ""
	}
	return content
}

func llenarTablaDeInodoDeArchivos(tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path, contenido string) {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			var bloqueArchivo estructuras.BloqueDeArchivos
			if StrlenBytes([]byte(contenido)) > 63 {
				posicionBloqueDeArchivo := inodos.WriteInBitmapBlock(path, superbloque)
				tablaInodo.I_block[i] = posicionBloqueDeArchivo
				copy(bloqueArchivo.B_content[:], []byte(contenido[:63]))
				comandos.Fwrite(&bloqueArchivo, path, superbloque.S_block_start+posicionBloqueDeArchivo*superbloque.S_block_size)
				llenarTablaDeInodoDeArchivos(tablaInodo, superbloque, path, contenido[63:])
				return
			} else {
				posicionBloqueDeArchivo := inodos.WriteInBitmapBlock(path, superbloque)
				tablaInodo.I_block[i] = posicionBloqueDeArchivo
				copy(bloqueArchivo.B_content[:], []byte(contenido[:]))
				comandos.Fwrite(&bloqueArchivo, path, superbloque.S_block_start+posicionBloqueDeArchivo*superbloque.S_block_size)
				return
			}
		}
	}
}
