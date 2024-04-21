package comandos

import (
	"P1/estructuras"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ParametrosMkdisk struct {
	Size int
	Fit  byte
	Unit byte
	Path string
}

type Mkdisk struct {
	Params ParametrosMkdisk
}

func (m *Mkdisk) Execute(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mkdisk(m.Params.Size, m.Params.Fit, m.Params.Unit, m.Params.Path) {
		//fmt.Println("Comando MKDISK realizado con exito")
	} else {
		fmt.Printf("Error al crear el disco: %s\n\n", m.Params.Path)
	}
}

func (m *Mkdisk) SaveParams(parametros []string) ParametrosMkdisk {
	incial := "A"
	for letra := incial[0]; letra <= 'Z'; letra++ {
		nombreArchivo := string(letra) + ".dsk"
		// Verificar si el archivo existe
		_, err := os.Stat(nombreArchivo)
		if os.IsNotExist(err) {
			m.Params.Path = nombreArchivo
			break
		}
	}
	for _, v := range parametros {
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error al convertir el tamaño del disco.")
			}
			m.Params.Size = num
		} else if strings.Contains(v, "unit") {
			v = strings.ReplaceAll(v, "unit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.Params.Unit = v[0]
		} else if strings.Contains(v, "fit") {
			v = strings.ReplaceAll(v, "fit=", "")
			v = strings.ReplaceAll(v, " ", "")
			m.Params.Fit = v[0]
		}
	}
	return m.Params
}

func (m *Mkdisk) Mkdisk(size int, fit byte, unit byte, path string) bool {
	var fileSize = 0
	var master estructuras.MBR
	// Comprobando si existe una ruta valida para la creacion del disco
	if path == "" {
		fmt.Println("Ruta no valida.")
		return false
	}
	// comprobando el tamano del disco, debe ser mayor que cero
	if size <= 0 {
		fmt.Printf("el tamaño del disco debe ser mayor que 0\n")
		return false
	}
	// tipo de unidad a utilizar, si el parametro esta vacio se utilizaran MegaBytes como default size
	if unit == 'k' || unit == 'K' {
		fileSize = size
	} else if unit == 'm' || unit == 'M' {
		fileSize = size * 1024
	} else if unit == 0 {
		fileSize = size * 1024
	} else {
		fmt.Println("Unit incompatible.")
		return false
	}
	// definiendo el tipo de fit que el disco tendra, como default sera First Fit
	if strconv.Itoa(int(fit)) == "66" || string(fit) == "BF" {
		master.Dsk_fit = 'b'
	} else if strconv.Itoa(int(fit)) == "70" || string(fit) == "FF" {
		master.Dsk_fit = 'f'
	} else if strconv.Itoa(int(fit)) == "87" || string(fit) == "WF" {
		master.Dsk_fit = 'w'
	} else if fit == 0 {
		master.Dsk_fit = 'f'
	} else {
		fmt.Printf("Fit incompatible\n")
		return false
	}
	// llenando el buffer con '0' para indicar que esta vacio.
	bloque := make([]byte, 1024)
	for i := 0; i < len(bloque); i++ {
		bloque[i] = 0
	}

	iterator := 0
	MkDirectory(path) // creando el directorio para el disco sino existe
	binaryFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error al crear el disco.\n")
		return false
	}
	defer binaryFile.Close()
	for iterator < fileSize {
		_, err := binaryFile.Write(bloque[:])
		if err != nil {
			fmt.Printf("Se produjo un error al escribir la información dentro del disco.\n")
		}
		iterator++
	}
	master.Mbr_tamano = int64(fileSize * 1024)
	master.Mbr_dsk_signature = GetRandom()
	// formateando el tiempo
	date := time.Now()
	for i := 0; i < len(master.Mbr_fecha_creacion)-1; i++ {
		master.Mbr_fecha_creacion[i] = date.String()[i]
	}
	FillPartitions(&master)
	WriteMBR(&master, path)
	return true
}

func FillPartitions(master *estructuras.MBR) {
	for i := 0; i < len(master.Mbr_partitions); i++ {
		master.Mbr_partitions[i].Part_status = '0'
		master.Mbr_partitions[i].Part_fit = '0'
		master.Mbr_partitions[i].Part_start = 0
		master.Mbr_partitions[i].Part_size = 0
		master.Mbr_partitions[i].Part_type = '0'
		copy(master.Mbr_partitions[i].Part_name[:], "")
	}
}
