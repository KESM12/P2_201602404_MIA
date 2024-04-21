package comandUser

import (
	"P1/comandos"
	"P1/estructuras"
	"P1/inodos"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func ReadFile(tablaInodo *estructuras.TablaInodo, path string, superbloque *estructuras.SuperBloque) string {
	contenido := ""
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo estructuras.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		contenido += string(parteArchivo.B_content[:])
	}
	return contenido
}

func AppendFile(path string, superbloque *estructuras.SuperBloque, tablaInodo *estructuras.TablaInodo, contenido string) bool {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] != -1 {
			var parteArchivo estructuras.BloqueDeArchivos
			comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
			if strlen(parteArchivo.B_content) == 63 {
				continue
			} else if strlen(parteArchivo.B_content) < 63 {
				nuevoContenido := string(TrimArray(parteArchivo.B_content[:])) + contenido
				// nuevoContenidoArray := createArray([]byte(nuevoContenido))
				if StrlenBytes([]byte(nuevoContenido)) > 63 {
					copy(parteArchivo.B_content[:], nuevoContenido[:63])
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
					AppendFile(path, superbloque, tablaInodo, string(nuevoContenido[63:]))
				} else {
					copy(parteArchivo.B_content[:], []byte(nuevoContenido))
					comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
				}
				return true
			}
		} else if tablaInodo.I_block[i] == -1 {
			var nuevoBloque estructuras.BloqueDeArchivos
			nuevaPosicion := inodos.WriteInBitmapBlock(path, superbloque)
			nuevoContenido := TrimArray([]byte(contenido))
			tablaInodo.I_block[i] = nuevaPosicion
			if StrlenBytes([]byte(contenido)) > 63 {
				copy(nuevoBloque.B_content[:], nuevoContenido[:63])
				comandos.Fwrite(&nuevoBloque, path, superbloque.S_block_start+nuevaPosicion*superbloque.S_block_size)
				AppendFile(path, superbloque, tablaInodo, string(nuevoContenido[63:]))
			} else {
				copy(nuevoBloque.B_content[:], nuevoContenido)
				comandos.Fwrite(&nuevoBloque, path, superbloque.S_block_start+nuevaPosicion*superbloque.S_block_size)
			}
			return true
		}
	}
	return false
}

func modFile(tablaInodo *estructuras.TablaInodo, path string, superbloque *estructuras.SuperBloque) string {
	contenido := ""
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo estructuras.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		contenido += string(parteArchivo.B_content[:])
		parteArchivo.B_content = [64]byte{}
		// fmt.Println(parteArchivo.B_content)
		copy(parteArchivo.B_content[:], "modificar")
		comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
	}
	return contenido
}

func SetFile(tablaInodo *estructuras.TablaInodo, path string, superbloque *estructuras.SuperBloque, contenido string) bool {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var parteArchivo estructuras.BloqueDeArchivos
		comandos.Fread(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		comparador := [64]byte{}
		copy(comparador[:], []byte("modificar"))
		if bytes.Equal(parteArchivo.B_content[:], comparador[:]) {
			if StrlenBytes([]byte(contenido)) > 63 {
				copy(parteArchivo.B_content[:], contenido[:63])
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
				return SetFile(tablaInodo, path, superbloque, string(contenido[63:]))
			} else {
				copy(parteArchivo.B_content[:], []byte(contenido))
				comandos.Fwrite(&parteArchivo, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
			}
			return true
		}

	}
	return false
}

func FindAndCreateDirectories(tablaInodo *estructuras.TablaInodo, path, ruta string, superbloque *estructuras.SuperBloque, posicionActual, userId, groupId int64) {
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		return
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueCarpeta estructuras.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			// crearemos un nuevo bloque de carpeta
			LlenarBloqueCarpetaVacio(&bloqueCarpeta)
			posicionNuevaBloque := inodos.WriteInBitmapBlock(path, superbloque)
			var TablaInodoNueva estructuras.TablaInodo
			tablaInodo.I_block[i] = posicionNuevaBloque
			comandos.Fwrite(tablaInodo, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
			posicionNuevaTablaInodo := inodos.WriteInBitmapInode(path, superbloque)
			CreateNewDirectory(&TablaInodoNueva, path, superbloque, posicionNuevaTablaInodo, posicionActual, userId, groupId)
			AgregarTablaNueva(&bloqueCarpeta, rutaParts[0], posicionNuevaTablaInodo)
			comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+posicionNuevaBloque*superbloque.S_block_size)
			FindAndCreateDirectories(&TablaInodoNueva, path, rutaParts[1], superbloque, posicionNuevaTablaInodo, userId, groupId)
			return
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		num, compare := CompareDirectories(rutaParts[0], &bloqueCarpeta)
		if compare {
			var tablaAuxiliar estructuras.TablaInodo
			comandos.Fread(&tablaAuxiliar, path, superbloque.S_inode_start+num*superbloque.S_inode_size)
			FindAndCreateDirectories(&tablaAuxiliar, path, rutaParts[1], superbloque, num, userId, groupId)
			return
		}
		if FreeSpace(&bloqueCarpeta) {
			var TablaInodoNueva estructuras.TablaInodo
			posicionNueva := inodos.WriteInBitmapInode(path, superbloque)
			CreateNewDirectory(&TablaInodoNueva, path, superbloque, posicionNueva, posicionActual, userId, groupId)
			AgregarTablaNueva(&bloqueCarpeta, rutaParts[0], posicionNueva)
			comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
			FindAndCreateDirectories(&TablaInodoNueva, path, rutaParts[1], superbloque, posicionNueva, userId, groupId)
			return
		}
	}
}

func LlenarBloqueCarpetaVacio(bloqueCarpeta *estructuras.BloqueDeCarpetas) {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		bloqueCarpeta.B_content[i].B_inodo = -1
		copy(bloqueCarpeta.B_content[i].B_name[:], "")
	}
}

func FreeSpace(bloqueCarpeta *estructuras.BloqueDeCarpetas) bool {
	for _, space := range bloqueCarpeta.B_content {
		if space.B_inodo == -1 {
			return true
		}
	}
	return false
}

func AgregarTablaNueva(bloqueCarpeta *estructuras.BloqueDeCarpetas, nombreTabla string, nuevaTabla int64) {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		if bloqueCarpeta.B_content[i].B_inodo == -1 {
			copy(bloqueCarpeta.B_content[i].B_name[:], []byte(nombreTabla))
			bloqueCarpeta.B_content[i].B_inodo = int32(nuevaTabla)
			// fmt.Println("nombre de B_inodo", string(bloqueCarpeta.B_content[i].B_name[:]))
			return
		}
	}
}

func FindDirectories(AgregarTabla int64, tablaInodo *estructuras.TablaInodo, path, ruta string, superbloque *estructuras.SuperBloque, posicionActual int64) {
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		crearArchivoDentroDeTablaInodo(AgregarTabla, posicionActual, tablaInodo, superbloque, path, ruta)
		return
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueDeCarpetas estructuras.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		comandos.Fread(&bloqueDeCarpetas, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		num, compare := CompareDirectories(rutaParts[0], &bloqueDeCarpetas)
		if compare {
			var nuevaTablaInodo estructuras.TablaInodo
			comandos.Fread(&nuevaTablaInodo, path, superbloque.S_inode_start+num*superbloque.S_inode_size)
			FindDirectories(AgregarTabla, &nuevaTablaInodo, path, rutaParts[1], superbloque, num)
			return
		}
	}
}

func FindDirs(AgregarTabla int64, tablaInodo *estructuras.TablaInodo, path, ruta string, superbloque *estructuras.SuperBloque, posicionActual int64) {
	var rutaParts []string
	if !strings.Contains(ruta, "/") {
		crearCarpetaDentroDeTablaInodo(AgregarTabla, posicionActual, tablaInodo, superbloque, path, ruta)
		return
	}
	rutaParts = strings.SplitN(ruta, "/", 2)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueDeCarpetas estructuras.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		comandos.Fread(&bloqueDeCarpetas, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		num, compare := CompareDirectories(rutaParts[0], &bloqueDeCarpetas)
		if compare {
			var nuevaTablaInodo estructuras.TablaInodo
			comandos.Fread(&nuevaTablaInodo, path, superbloque.S_inode_start+num*superbloque.S_inode_size)
			FindDirs(AgregarTabla, &nuevaTablaInodo, path, rutaParts[1], superbloque, num)
			return
		}
	}
}

func crearCarpetaDentroDeTablaInodo(AgregarTabla, posicionActual int64, tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path, nombreArchivo string) {
	var tabla estructuras.TablaInodo
	comandos.Fread(&tabla, path, superbloque.S_inode_start+AgregarTabla*superbloque.S_inode_size)
	var primerbloque estructuras.BloqueDeCarpetas
	comandos.Fread(&primerbloque, path, superbloque.S_block_start+tabla.I_block[0]*superbloque.S_block_size)
	primerbloque.B_content[1].B_inodo = int32(posicionActual)
	comandos.Fwrite(&primerbloque, path, superbloque.S_block_start+tabla.I_block[0]*superbloque.S_block_size)
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueCarpeta estructuras.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			LlenarBloqueCarpetaVacio(&bloqueCarpeta)
			posicionNuevaBloque := inodos.WriteInBitmapBlock(path, superbloque)
			tablaInodo.I_block[i] = posicionNuevaBloque
			comandos.Fwrite(tablaInodo, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
			AgregarTablaNueva(&bloqueCarpeta, nombreArchivo, AgregarTabla)
			comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+posicionNuevaBloque*superbloque.S_block_size)
			return
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		if !FreeSpace(&bloqueCarpeta) {
			continue
		}
		AgregarTablaNueva(&bloqueCarpeta, nombreArchivo, AgregarTabla)
		comandos.Fwrite(bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		return
	}
}

func crearArchivoDentroDeTablaInodo(AgregarTabla, posicionActual int64, tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path, nombreArchivo string) {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		var bloqueCarpeta estructuras.BloqueDeCarpetas
		if tablaInodo.I_block[i] == -1 {
			posicionBloqueCarpeta := inodos.WriteInBitmapBlock(path, superbloque)
			LlenarBloqueCarpetaVacio(&bloqueCarpeta)
			tablaInodo.I_block[i] = posicionBloqueCarpeta
			AgregarTablaNueva(&bloqueCarpeta, nombreArchivo, AgregarTabla)
			comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+posicionBloqueCarpeta*superbloque.S_block_size)
			comandos.Fwrite(tablaInodo, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
			return
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		if !FreeSpace(&bloqueCarpeta) {
			continue
		}
		AgregarTablaNueva(&bloqueCarpeta, nombreArchivo, AgregarTabla)
		comandos.Fwrite(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		return
	}
}

func CreateNewDirectory(nuevaTabla *estructuras.TablaInodo, path string, superbloque *estructuras.SuperBloque, posicionActual, posicionPadre, userId, groupId int64) {
	nuevaTabla.I_uid = userId
	nuevaTabla.I_gid = groupId
	nuevaTabla.I_size = 0
	nuevaTabla.I_type = '0'
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
	posicionNuevoBloqueCarpetas := inodos.WriteInBitmapBlock(path, superbloque)
	nuevaTabla.I_block[0] = posicionNuevoBloqueCarpetas
	comandos.Fwrite(nuevaTabla, path, superbloque.S_inode_start+posicionActual*superbloque.S_inode_size)
	nuevoBloqueCarpetas := estructuras.BloqueDeCarpetas{}
	copy(nuevoBloqueCarpetas.B_content[0].B_name[:], ".")
	nuevoBloqueCarpetas.B_content[0].B_inodo = int32(posicionActual)
	copy(nuevoBloqueCarpetas.B_content[1].B_name[:], "..")
	nuevoBloqueCarpetas.B_content[1].B_inodo = int32(posicionPadre)
	copy(nuevoBloqueCarpetas.B_content[2].B_name[:], "")
	nuevoBloqueCarpetas.B_content[2].B_inodo = -1
	copy(nuevoBloqueCarpetas.B_content[3].B_name[:], "")
	nuevoBloqueCarpetas.B_content[3].B_inodo = -1
	comandos.Fwrite(&nuevoBloqueCarpetas, path, superbloque.S_block_start+posicionNuevoBloqueCarpetas*superbloque.S_block_size)
}

func SearchFreeSpace(tablaInodo *estructuras.TablaInodo, path string, superbloque *estructuras.SuperBloque) estructuras.BloqueDeCarpetas {
	var bloqueCarpeta estructuras.BloqueDeCarpetas
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			return bloqueCarpeta
		}
		comandos.Fread(&bloqueCarpeta, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		if !SearchFreeSpaceInBlock(&bloqueCarpeta) {
			continue
		}
	}
	return bloqueCarpeta
}

func CreateNewFile(bloqueCarpeta *estructuras.BloqueDeCarpetas, path string, superbloque *estructuras.SuperBloque) {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		if bloqueCarpeta.B_content[i].B_inodo != -1 {
			continue
		}
		posicion := inodos.WriteInBitmapBlock(path, superbloque)
		bloqueCarpeta.B_content[i].B_inodo = int32(posicion)

	}
}

func SearchFreeSpaceInBlock(bloqueCarpeta *estructuras.BloqueDeCarpetas) bool {
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		if bloqueCarpeta.B_content[i].B_inodo != -1 {
			continue
		}
		return true
	}
	return false
}

func CompareDirectories(rutaPart string, bloqueCarpeta *estructuras.BloqueDeCarpetas) (int64, bool) {
	for _, content := range bloqueCarpeta.B_content {
		if string(TrimArray(content.B_name[:])) == string(TrimArray([]byte(rutaPart[:]))) {
			return int64(content.B_inodo), true
		}
	}
	return -1, false
}

func PrintTree(tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path string) {
	fmt.Println("I_uid:", tablaInodo.I_uid)
	fmt.Println("I_gid:", tablaInodo.I_gid)
	fmt.Println("I_size:", tablaInodo.I_size)
	fmt.Println("I_atime:", string(tablaInodo.I_atime[:]))
	fmt.Println("I_ctime:", string(tablaInodo.I_ctime[:]))
	fmt.Println("I_mtime:", string(tablaInodo.I_mtime[:]))
	for i := 0; i < len(tablaInodo.I_block); i++ {
		fmt.Printf("I_block[%d]: %d\n", i, tablaInodo.I_block[i])
	}
	fmt.Println("I_type:", tablaInodo.I_type)
	fmt.Println("I_perm:", tablaInodo.I_perm)
	if tablaInodo.I_type == '0' {
		PrintBloqueDeCarpetas(tablaInodo, superbloque, path)
	} else {
		PrintBloqueDeArchivos(tablaInodo, superbloque, path)
	}
}

func PrintBloqueDeCarpetas(tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path string) {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var bloqueDeCarpetas estructuras.BloqueDeCarpetas
		comandos.Fread(&bloqueDeCarpetas, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		for _, contenido := range bloqueDeCarpetas.B_content {
			fmt.Println("B_name:", string(contenido.B_name[:]))
			fmt.Println("B_inodo:", contenido.B_inodo)
		}
	}
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var bloqueDeCarpetas estructuras.BloqueDeCarpetas
		comandos.Fread(&bloqueDeCarpetas, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		for _, contenido := range bloqueDeCarpetas.B_content {
			var nuevaTablaInodo estructuras.TablaInodo
			if contenido.B_inodo == -1 || string(TrimArray(contenido.B_name[:])) == "." || string(TrimArray(contenido.B_name[:])) == ".." {
				continue
			}
			comandos.Fread(&nuevaTablaInodo, path, superbloque.S_inode_start+int64(contenido.B_inodo)*superbloque.S_inode_size)
			PrintTree(&nuevaTablaInodo, superbloque, path)
		}
	}
}

func PrintBloqueDeArchivos(tablaInodo *estructuras.TablaInodo, superbloque *estructuras.SuperBloque, path string) {
	for i := 0; i < len(tablaInodo.I_block); i++ {
		if tablaInodo.I_block[i] == -1 {
			continue
		}
		var bloqueDeArchivos estructuras.BloqueDeArchivos
		comandos.Fread(&bloqueDeArchivos, path, superbloque.S_block_start+tablaInodo.I_block[i]*superbloque.S_block_size)
		fmt.Printf("bloque[%d]: %s\n", i, string(bloqueDeArchivos.B_content[:]))
	}
}

func strlen(arr [64]byte) int {
	count := 0
	for _, c := range arr {
		if c != 0 {
			count++
		}
	}
	return count
}

func StrlenBytes(arr []byte) int {
	count := 0
	for _, c := range arr {
		if c != 0 {
			count++
		}
	}
	return count
}

func TrimArray(arr []byte) []byte {
	var result []byte
	for _, v := range arr {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
