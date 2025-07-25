package comandos

import (
	"P1/estructuras"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type ParametrosFdisk struct {
	Size int
	Unit byte
	Path string
	Type byte
	Fit  byte
	Name [16]byte
}

type Fdisk struct {
	Params ParametrosFdisk
}

func (f *Fdisk) Execute(parametros []string) {
	f.Params = f.SaveParams(parametros)
	if f.Fdisk(f.Params.Name, f.Params.Path, f.Params.Size, f.Params.Unit, f.Params.Fit, f.Params.Type) {
		fmt.Printf("\n FDISK realizado con exito para la ruta: %s y particion: %s\n\n", f.Params.Path, string(f.Params.Name[:]))
	} else {
		fmt.Printf("\n No se logro realizar el comando FDISK para la ruta: %s\n\n", f.Params.Path)
	}
}

func (f *Fdisk) SaveParams(parametros []string) ParametrosFdisk {
	for _, v := range parametros {
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "driveletter") {
			v = strings.ReplaceAll(v, "driveletter=", "")
			v = strings.ReplaceAll(v, "\"", "")
			v = v + ".dsk"
			f.Params.Path = v
		} else if strings.Contains(v, "size") {
			v = strings.ReplaceAll(v, "size=", "")
			v = strings.ReplaceAll(v, " ", "")
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error al convertir a entero")
			}
			f.Params.Size = num
		} else if strings.Contains(v, "unit") {
			v = strings.ReplaceAll(v, "unit=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.Params.Unit = ' '
			} else {
				f.Params.Unit = v[0]
			}
		} else if strings.Contains(v, "fit") {
			v = strings.ReplaceAll(v, "fit=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.Params.Fit = ' '
			} else {
				f.Params.Fit = v[0]
			}

		} else if strings.Contains(v, "type") {
			v = strings.ReplaceAll(v, "type=", "")
			v = strings.ReplaceAll(v, " ", "")
			if v == "" {
				f.Params.Type = ' '
			} else {
				f.Params.Type = v[0]
			}
		} else if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			copy(f.Params.Name[:], v)
		}
	}
	return f.Params
}

func (f *Fdisk) Fdisk(name [16]byte, path string, size int, unit byte, fit byte, t byte) bool {
	if path == "" {
		fmt.Printf("Ruta no encontrad.\n")
		return false
	}
	master := GetMBR(path)
	newPartition := estructuras.Partition{}
	fileSize := 0
	if unit == 'b' || unit == 'B' {
		fileSize = size
	} else if unit == 'k' || unit == 'K' {
		fileSize = size * 1024
	} else if unit == 'm' || unit == 'M' {
		fileSize = size * 1024 * 1024
	} else if unit == 0 {
		fileSize = size * 1024
	} else {
		fmt.Printf("Unit incorrecta.\n")
		return false
	}
	// se debe comprobar que no exista ninguna particion con el mismo nombre
	if ExisteParticion(&master, name) {
		fmt.Printf("Partición ya existente: \"%s\"\n", string(estructuras.TrimArray(name[:])))
		return false
	}
	// comprobando el tamano de la particion, este debe ser mayor que cero
	if size <= 0 {
		fmt.Printf("El parametro size no puede ser 0 o menor a 0.\n")
		return false
	}
	// definiendo el tipo de fit que la particion tendra, como default se utilizara Worst Fit
	if strconv.Itoa(int(fit)) == "87" || fit == 'W' {
		newPartition.Part_fit = 'w'
	} else if strconv.Itoa(int(fit)) == "66" || fit == 'B' {
		newPartition.Part_fit = 'b'
	} else if strconv.Itoa(int(fit)) == "70" || fit == 'F' {
		newPartition.Part_fit = 'f'
	} else if fit == 0 {
		newPartition.Part_fit = 'w'
	} else {
		fmt.Printf("Fit invalido.\n")
		return false
	}
	// verificando que el tamano de la particion a crear sea menor
	// o igual que el tamano que queda en el disco.
	totalSize := int(unsafe.Sizeof(estructuras.MBR{}))
	for _, v := range master.Mbr_partitions {
		if v.Part_status == '1' {
			totalSize += int(v.Part_size)
		}
	}
	// fmt.Println("espacio disponible, espacio a utilizar:", int(master.Mbr_tamano)-totalSize, fileSize)
	if t != 'l' && t != 'L' {
		if fileSize > int(master.Mbr_tamano)-int(totalSize) {
			fmt.Printf("Tamaño de la partición muy grande.\n")
			return false
		}
	}

	// indicando el tipo de particion
	if t == 0 {
		t = 'p'
	} else if t != 'p' && t != 'e' && t != 'l' && t != 'P' && t != 'E' && t != 'L' {
		fmt.Printf("Tipo de partición invalido: \"%c\"\n", t)
		return false
	}
	newPartition.Part_size = int64(fileSize)
	newPartition.Part_type = t
	newPartition.Part_status = '1'
	copy(newPartition.Part_name[:], name[:])

	// revisando que no exista mas de una particion Extendida y que Exista en caso de que se vaya a crear una particion logica
	existeParticionExtendida := false //esta variable se utiliza para encontrar si existe una particion extendida
	var whereToStart int              // con este valor le vamos a pasar a la particion logica donde comienza la particion extendida
	var partitionSize int             // con este valor le indicamos a la particion logica cuanto espacio ocupa la particion extendida
	var extendedFit byte              // con este valor le indicamos a la particion logica el tipo de ajuste que tiene la particion extendida
	var extendedName [16]byte         // con este valor le indicamos el nombre de la particion extendida a la particion logica
	// aqui le agregamos a las variables anteriores su correspondiente valor
	for _, v := range master.Mbr_partitions {
		if v.Part_type == 'e' || v.Part_type == 'E' {
			copy(extendedName[:], v.Part_name[:])
			existeParticionExtendida = true
			extendedFit = v.Part_fit
			whereToStart = int(v.Part_start)
			partitionSize = int(v.Part_size)
		}
	}

	// comprobamos que exista una particion libre
	existeParticionLibre := false
	if t != 'l' && t != 'L' {
		for _, v := range master.Mbr_partitions {
			if v.Part_status == '0' {
				existeParticionLibre = true
			}
		}
	} else if t == 'l' || t == 'L' {
		existeParticionLibre = true
	}
	// sino se encuentra un espacio libre para particion
	if !existeParticionLibre {
		fmt.Printf("Disco lleno, particiones completas.\n")
		return false
	}
	// comprobamos que tipo de particion es, luego la creamos
	if t == 'p' || t == 'P' {
		f.CreatePrimaryPartition(&master, newPartition)
	} else if t == 'e' || t == 'E' {
		if existeParticionExtendida {
			fmt.Printf("Solo se permite una partición extendida.\n")
			return false
		}
		f.CreateExtendedPartition(&master, newPartition, path)
	} else if t == 'l' || t == 'L' {
		if !existeParticionExtendida {
			fmt.Printf("Partición extendida no encontrada.")
			return false
		}
		particionLogica := estructuras.EBR{}
		particionLogica.Part_fit = newPartition.Part_fit
		particionLogica.Part_next = -1
		particionLogica.Part_size = newPartition.Part_size
		particionLogica.Part_status = newPartition.Part_status
		copy(particionLogica.Part_name[:], newPartition.Part_name[:])
		return f.CreateLogicPartition(&particionLogica, path, whereToStart, partitionSize, extendedFit, extendedName)

	}
	WriteMBR(&master, path)
	PrintPartitions(&master)
	return true
}

func (f *Fdisk) CreatePrimaryPartition(master *estructuras.MBR, newPartition estructuras.Partition) {
	// Asignacion de que particion es la que se utilizara
	if master.Dsk_fit == 'b' {
		BestFit(master, &newPartition)
	} else if master.Dsk_fit == 'w' {
		WorstFit(master, &newPartition)
	} else if master.Dsk_fit == 'f' {
		FirstFit(master, &newPartition)
	}
}

func (f Fdisk) CreateExtendedPartition(master *estructuras.MBR, newPartition estructuras.Partition, path string) {
	// Asignacion de que particion es la que se utilizara
	if master.Dsk_fit == 'b' {
		BestFit(master, &newPartition)
	} else if master.Dsk_fit == 'w' {
		WorstFit(master, &newPartition)
	} else if master.Dsk_fit == 'f' {
		FirstFit(master, &newPartition)
	}
	temp := estructuras.EBR{}
	temp.Part_status = '0'
	temp.Part_fit = '0'
	temp.Part_start = newPartition.Part_start
	temp.Part_size = 0
	temp.Part_next = -1
	copy(temp.Part_name[:], "vacio")
	WriteEBR(&temp, path, newPartition.Part_start)
}

func BestFit(master *estructuras.MBR, newPartition *estructuras.Partition) {
	bestFit := 0

	encontroParticion := false
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size {
			if i != bestFit {
				if v.Part_size < master.Mbr_partitions[bestFit].Part_size {
					encontroParticion = true
					bestFit = i
				}
			}
		}
	}
	if !encontroParticion {
		for i, v := range master.Mbr_partitions {
			if v.Part_start == 0 {
				bestFit = i
				break
			}
		}
	}
	master.Mbr_partitions[bestFit] = *newPartition
	if bestFit == 0 {
		master.Mbr_partitions[bestFit].Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
		newPartition.Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
	} else {
		master.Mbr_partitions[bestFit].Part_start = master.Mbr_partitions[bestFit-1].Part_start + master.Mbr_partitions[bestFit-1].Part_size
		newPartition.Part_start = master.Mbr_partitions[bestFit-1].Part_start + master.Mbr_partitions[bestFit-1].Part_size
	}
}

func WorstFit(master *estructuras.MBR, newPartition *estructuras.Partition) {
	worstFit := 0
	encontroParticion := false
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size {
			if i != worstFit {
				if v.Part_size > master.Mbr_partitions[worstFit].Part_size {
					worstFit = i
					encontroParticion = true
				}
			}
		}
	}
	if !encontroParticion {
		for i, v := range master.Mbr_partitions {
			fmt.Println(v.Part_start)
			if v.Part_start == 0 {
				worstFit = i
				break
			}
		}
	}
	master.Mbr_partitions[worstFit] = *newPartition
	if worstFit == 0 {
		master.Mbr_partitions[worstFit].Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
		newPartition.Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
	} else {
		master.Mbr_partitions[worstFit].Part_start = master.Mbr_partitions[worstFit-1].Part_start + master.Mbr_partitions[worstFit-1].Part_size
		newPartition.Part_start = master.Mbr_partitions[worstFit-1].Part_start + master.Mbr_partitions[worstFit-1].Part_size
	}
}

func FirstFit(master *estructuras.MBR, newPartition *estructuras.Partition) {
	firstFit := 0
	for i, v := range master.Mbr_partitions {
		if v.Part_status == '5' && v.Part_size >= newPartition.Part_size || v.Part_start == 0 {
			firstFit = i
			break
		}
	}
	master.Mbr_partitions[firstFit] = *newPartition
	if firstFit == 0 {
		master.Mbr_partitions[firstFit].Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
		newPartition.Part_start = int64(unsafe.Sizeof(estructuras.MBR{}))
	} else {
		master.Mbr_partitions[firstFit].Part_start = master.Mbr_partitions[firstFit-1].Part_start + master.Mbr_partitions[firstFit-1].Part_size
		newPartition.Part_start = master.Mbr_partitions[firstFit-1].Part_start + master.Mbr_partitions[firstFit-1].Part_size
	}
}

func (f *Fdisk) CreateLogicPartition(logicPartition *estructuras.EBR, path string, whereToStart int, partitionSize int, extendedFit byte, extendedName [16]byte) bool {
	if extendedFit == 'f' {
		return FirstFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	} else if extendedFit == 'b' {
		return BestFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	} else if extendedFit == 'w' {
		return WorstFitLogicPart(logicPartition, path, whereToStart, partitionSize, extendedName)
	}
	return false
}

func FirstFitLogicPart(logicPartition *estructuras.EBR, path string, whereToStart int, partitionSize int, extendedName [16]byte) bool {
	temp := estructuras.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("La partición logica no puede ser mas grande que la partición extendida.")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
		} else if temp.Part_status == '5' {
			if temp.Part_size >= logicPartition.Part_size {
				logicPartition.Part_start = temp.Part_start
				logicPartition.Part_next = temp.Part_next
				WriteEBR(logicPartition, path, temp.Part_start)
				flag = false
			}
		} else if temp.Part_next == -1 {
			totalSize += int(temp.Part_size)
			if partitionSize < totalSize {
				fmt.Println("Espacio insuficiente")
				return false
			}
			temp.Part_next = temp.Part_start + temp.Part_size
			logicPartition.Part_start = temp.Part_next
			WriteEBR(&temp, path, temp.Part_start)
			WriteEBR(logicPartition, path, temp.Part_next)
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}

	PrintLogicPartitions(path, int64(whereToStart), int64(partitionSize), extendedName)
	return true
}

func BestFitLogicPart(logicPartition *estructuras.EBR, path string, whereToStart int, partitionSize int, extendedName [16]byte) bool {
	var particionesLogicas []estructuras.EBR
	temp := estructuras.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	Wrote := false
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("La partición logica no puede ser mas grande que la partición extendida.")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
			Wrote = true
		} else if temp.Part_status == '5' {
			particionesLogicas = append(particionesLogicas, temp)
		} else if temp.Part_next == -1 {
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}
	bestFit := 0
	tempSize := 0
	if len(particionesLogicas) != 0 {
		for i, v := range particionesLogicas {
			if tempSize != 0 {
				bestFit = i
			} else if tempSize > int(v.Part_size) && v.Part_size >= logicPartition.Part_size {
				tempSize = int(v.Part_size)
				bestFit = i
			}
		}
		logicPartition.Part_start = particionesLogicas[bestFit].Part_start
		logicPartition.Part_next = particionesLogicas[bestFit].Part_next
		WriteEBR(logicPartition, path, logicPartition.Part_start)
		Wrote = true
	}
	if !Wrote {
		totalSize = int(logicPartition.Part_size)
		temp = ReadEBR(path, int64(whereToStart))
		flag2 := true
		for flag2 {
			if temp.Part_next == -1 {
				totalSize += int(temp.Part_size)
				if partitionSize < totalSize {
					fmt.Println("Espacio insuficiente")
					return false
				}
				temp.Part_next = temp.Part_start + temp.Part_size
				logicPartition.Part_start = temp.Part_next
				WriteEBR(&temp, path, temp.Part_start)
				WriteEBR(logicPartition, path, temp.Part_next)
				flag2 = false
			} else {
				totalSize += int(temp.Part_size)
				temp = ReadEBR(path, temp.Part_next)
			}
		}
	}

	PrintLogicPartitions(path, int64(whereToStart), int64(partitionSize), extendedName)
	return true
}

func WorstFitLogicPart(logicPartition *estructuras.EBR, path string, whereToStart int, partitionSize int, extendedName [16]byte) bool {
	var particionesLogicas []estructuras.EBR
	temp := estructuras.EBR{}
	totalSize := 0
	totalSize += int(logicPartition.Part_size)
	temp = ReadEBR(path, int64(whereToStart))
	Wrote := false
	flag := true
	for flag {
		if temp.Part_size == 0 {
			if partitionSize < int(logicPartition.Part_size) {
				fmt.Println("La partición logica no puede ser mas grande que la partición extendida.")
				return false
			}
			logicPartition.Part_start = int64(whereToStart)
			WriteEBR(logicPartition, path, int64(whereToStart))
			flag = false
			Wrote = true
		} else if temp.Part_status == '5' {
			particionesLogicas = append(particionesLogicas, temp)
		} else if temp.Part_next == -1 {
			flag = false
		} else {
			totalSize += int(temp.Part_size)
			temp = ReadEBR(path, temp.Part_next)
		}
	}
	worstFit := 0
	tempSize := 0
	if len(particionesLogicas) != 0 {
		for i, v := range particionesLogicas {
			if tempSize != 0 {
				worstFit = i
			} else if tempSize < int(v.Part_size) && v.Part_size >= logicPartition.Part_size {
				tempSize = int(v.Part_size)
				worstFit = i
			}
		}
		logicPartition.Part_start = particionesLogicas[worstFit].Part_start
		logicPartition.Part_next = particionesLogicas[worstFit].Part_next
		WriteEBR(logicPartition, path, logicPartition.Part_start)
		Wrote = true
	}
	if !Wrote {
		totalSize = int(logicPartition.Part_size)
		temp = ReadEBR(path, int64(whereToStart))
		flag2 := true
		for flag2 {
			if temp.Part_next == -1 {
				totalSize += int(temp.Part_size)
				if partitionSize < totalSize {
					fmt.Println("Espacio insuficiente")
					return false
				}
				temp.Part_next = temp.Part_start + temp.Part_size
				logicPartition.Part_start = temp.Part_next
				WriteEBR(&temp, path, temp.Part_start)
				WriteEBR(logicPartition, path, temp.Part_next)
				flag2 = false
			} else {
				totalSize += int(temp.Part_size)
				temp = ReadEBR(path, temp.Part_next)
			}
		}
	}

	PrintLogicPartitions(path, int64(whereToStart), int64(partitionSize), extendedName)
	return true
}

func ExisteParticion(master *estructuras.MBR, name [16]byte) bool {
	for _, v := range master.Mbr_partitions {
		if bytes.Equal(v.Part_name[:], name[:]) {
			return true
		}
	}
	return false
}

func PrintPartitions(master *estructuras.MBR) {
	str := ""
	for i := 0; i < 70; i++ {
		str += "-"
	}
	contenido := ""
	contenido += fmt.Sprintf("%s\n", str)
	contenido += fmt.Sprintf("%-20s%-10s%-10s%-10s%-10s%-10s\n", "Name", "Type", "Fit", "Start", "Size", "Status")
	for _, part := range master.Mbr_partitions {
		contenido += fmt.Sprintf("%s\n", str)
		if string(estructuras.TrimArray(part.Part_name[:])) == "" {
			contenido += fmt.Sprintf("%-20s", "-")
		} else {
			contenido += fmt.Sprintf("%-20s", string(estructuras.TrimArray(part.Part_name[:])))
		}
		if part.Part_type == '0' {
			contenido += fmt.Sprintf("%-10c", '-')
		} else {
			contenido += fmt.Sprintf("%-10c", part.Part_type)
		}
		if part.Part_fit == '0' {
			contenido += fmt.Sprintf("%-10c", '-')
		} else {
			contenido += fmt.Sprintf("%-10c", part.Part_fit)
		}
		contenido += fmt.Sprintf("%-10d", part.Part_start)
		contenido += fmt.Sprintf("%-10d", part.Part_size)
		contenido += fmt.Sprintf("%-10c\n", part.Part_status)
	}
	contenido += fmt.Sprintf("%s\n", str)
	fmt.Println(contenido)
}

func PrintLogicPartitions(path string, whereToStart, PartitionSize int64, extendedName [16]byte) {
	str := ""
	for i := 0; i < 70; i++ {
		str += "-"
	}
	contenido := ""
	contenido += fmt.Sprintf("Partition name: %s\n", string(estructuras.TrimArray(extendedName[:])))
	contenido += fmt.Sprintf("Partition size: %d\n", PartitionSize)
	contenido += fmt.Sprintf("%s\n", str)
	contenido += fmt.Sprintf("%-20s%-12s%-10s%-10s%-10s%-10s\n", "Name", "Next Part", "Fit", "Start", "Size", "Status")
	var temp estructuras.EBR
	Fread(&temp, path, whereToStart)
	flag := true
	for flag {
		contenido += fmt.Sprintf("%s\n", str)
		if string(estructuras.TrimArray(temp.Part_name[:])) == "" {
			contenido += fmt.Sprintf("%-20s", "Disponible")
		} else {
			contenido += fmt.Sprintf("%-20s", string(estructuras.TrimArray(temp.Part_name[:])))
		}
		contenido += fmt.Sprintf("%-12d", temp.Part_next)
		if temp.Part_fit == '0' {
			contenido += fmt.Sprintf("%-10c", '-')
		} else {
			contenido += fmt.Sprintf("%-10c", temp.Part_fit)
		}
		contenido += fmt.Sprintf("%-10d", temp.Part_start)
		contenido += fmt.Sprintf("%-10d", temp.Part_size)
		contenido += fmt.Sprintf("%-10c\n", temp.Part_status)
		if temp.Part_next == -1 {
			flag = false
		} else {
			Fread(&temp, path, temp.Part_next)
		}
	}
	contenido += fmt.Sprintf("%s\n", str)
	fmt.Println(contenido)
}
