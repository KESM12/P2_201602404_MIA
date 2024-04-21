package comandos

import (
	"fmt"
	"os"
	"strings"
)

type ParametrosRmdisk struct {
	Path string
}

type Rmdisk struct {
	Params ParametrosRmdisk
}

func (r *Rmdisk) Execute(parametros []string) {
	r.Params = r.SaveParams(parametros)
	fmt.Printf("¿Quieres eliminar el disco %s? [Y/N]: ", r.Params.Path)
	var respuesta string
	fmt.Scanln(&respuesta)
	// Convertir la respuesta a mayúsculas para manejar entradas en minúsculas
	respuesta = strings.ToUpper(respuesta)

	if respuesta == "Y" || respuesta == "S" {
		if r.Rmdisk(r.Params.Path) {
			fmt.Printf("\n Disco eliminado correctamente: %s\n\n", r.Params.Path)
		}
	} else if respuesta == "N" {
		fmt.Println("Archivo no eliminado.")
	} else {
		fmt.Println("Respuesta no válida. Por favor, ingresa 'Y', 'S' o 'N'.")
	}
}

func (r *Rmdisk) SaveParams(parametros []string) ParametrosRmdisk {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "driveletter") {
			v = strings.ReplaceAll(v, "driveletter=", "")
			v = strings.ReplaceAll(v, "\"", "")
			v = v + ".dsk"
			r.Params.Path = v
		}
	}
	return r.Params
}

func (r *Rmdisk) Rmdisk(path string) bool {
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		fmt.Printf("Ruta no encontrada.\n")
		return false
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Disco no eliminado: %s", r.Params.Path)
		return false
	}
	return true
}
