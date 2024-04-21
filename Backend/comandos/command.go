package comandos

import (
	"P1/estructuras"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"reflect"
	"time"
)

func FileExist(path string) bool {
	fmt.Printf("-%s-\n", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func WriteMBR(master *estructuras.MBR, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("no se pudo abrir el archivo para escribir el MBR %s\n", err.Error())
		return
	}
	// Posicionandonos en el principio del archivo
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("no se pudo posicionar en el principio del archivo %s]n", err.Error())
		return
	}
	// Escribiendo el MBR
	// var masterBuffer bytes.Buffer
	err = binary.Write(file, binary.LittleEndian, master)
	if err != nil {
		fmt.Printf("no se pudo escribir el MBR %s\n", err.Error())
		file.Close()
		return
	}
	// fmt.Printf("se escribio correctamente! :D")
	defer file.Close()
}

func GetMBR(path string) estructuras.MBR {
	var mbr estructuras.MBR
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("no se pudo abrir el archivo para obtener el MBR, %s\n", err.Error())
		return mbr
	}

	defer file.Close()

	// leyendo el mbr del archivo
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Printf("no se pudo obtener la informacion del archivo para obtener el MBR %s\n", err.Error())
		return mbr
	}
	return mbr
}

func WriteEBR(ebr *estructuras.EBR, path string, position int64) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("no se pudo abrir el archivo para escribir el MBR %s\n", err.Error())
		return
	}
	// Posicionandonos en el principio del archivo
	_, err = file.Seek(position, 0)
	if err != nil {
		fmt.Printf("no se pudo posicionar en el principio del archivo %s\n", err.Error())
		return
	}
	// Escribiendo el MBR
	// var masterBuffer bytes.Buffer
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		fmt.Printf("no se pudo escribir el MBR %s\n", err.Error())
		file.Close()
		return
	}
	// fmt.Printf("se escribio correctamente! :D")
	defer file.Close()
}

func ReadEBR(path string, position int64) estructuras.EBR {
	var ebr estructuras.EBR
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("no se pudo abrir el archivo para obtener el MBR %s\n", err.Error())
		return ebr
	}

	defer file.Close()

	// leyendo el mbr del archivo
	file.Seek(position, 0)
	err = binary.Read(file, binary.LittleEndian, &ebr)
	if err != nil {
		fmt.Printf("No se puedo obtener el MBR %s\n", err.Error())
		return ebr
	}
	return ebr
}

func MkDirectory(fullPath string) {
	directory := path.Dir(fullPath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0777)
		if err != nil {
			fmt.Printf("No se puedo crear: %s\n", err.Error())
		}
	}
}

func GetRandom() int64 {
	rand.Seed(time.Now().UnixNano())
	n := 150
	randomNumber := rand.Intn(n)
	return int64(randomNumber)
}

// funcion general para escribir SuperBloque, TablaInodo, BloqueDeArchivos, BloqueDeCarpetas
func Fwrite(estructura interface{}, path string, position int64) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("No se pudo escribir la estructura. %s\n", err.Error())
		return
	}
	// Posicionandonos en donde necesitamos dentro del archivo
	_, err = file.Seek(position, 0)
	if err != nil {
		fmt.Printf("Posición incorrecta: %d, %s\n", position, err.Error())
		return
	}
	// Escribiendo la estructura
	err = binary.Write(file, binary.LittleEndian, estructura)
	if err != nil {
		fmt.Printf("No se pudo escribir la estructura. %s\n", err.Error())
		file.Close()
		return
	}
	// fmt.Printf("se escribio correctamente! :D")
	defer file.Close()
}

func Fread(estructura interface{}, path string, position int64) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("No se pudo abrir el archivo. %s\n", err.Error())
		return
	}
	// Posicionandonos en donde necesitamos dentro del archivo
	_, err = file.Seek(position, 0)
	if err != nil {
		fmt.Printf("Posición incorrecta: %d, %s\n", position, err.Error())
		return
	}
	// Leyendo La estructura
	err = binary.Read(file, binary.LittleEndian, estructura)
	if err != nil {
		fmt.Printf("Estructura incorrecta, %s, %s, %s\n", reflect.TypeOf(estructura).String(), ":", err.Error())
		return
	}
	defer file.Close()
}

func Fopen(path, contenido string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Archivo: %s, %s\n", path, err.Error())
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(contenido))
	if err != nil {
		fmt.Printf("Archivo: %s, %s\n", path, err.Error())
	}

	fmt.Printf("Archivo creado con exito: %s\n", path)
}
