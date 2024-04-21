package comandos

import (
	"P1/estructuras"
	"P1/lista"
	"bytes"
	"fmt"
	"strings"
)

type ParametrosMount struct {
	Path string
	Name [16]byte
}

type Mount struct {
	Params ParametrosMount
}

func (m *Mount) Execute(parametros []string) {
	m.Params = m.SaveParams(parametros)
	if m.Mount(m.Params.Path, m.Params.Name) {
		fmt.Printf("Particion %s montada con éxito\n\n", m.Params.Path)
	} else {
		fmt.Printf("Partición no montada. %s\n", m.Params.Path)
	}
}

func (m *Mount) SaveParams(parametros []string) ParametrosMount {
	// fmt.Println(parametros)
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "driveletter") {
			v = strings.ReplaceAll(v, "driveletter=", "")
			v = strings.ReplaceAll(v, "\"", "")
			v = v + ".dsk"
			m.Params.Path = v
		} else if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			copy(m.Params.Name[:], v[:])
		}
	}
	return m.Params
}

func (m *Mount) Mount(path string, name [16]byte) bool {
	// comprobando que el parametro "path" sea diferente de ""
	if path == "" {
		fmt.Printf("Ruta incorrecta.")
		return false
	}
	// comprobando que el parametro "name" sea diferente de ""
	if bytes.Equal(name[:], []byte("")) {
		fmt.Printf("Name no puede estar vacio.\n")
		return false
	}
	master := GetMBR(path)
	partitionMounted := false
	particionEncontrada := false
	for _, particion := range master.Mbr_partitions {
		// si entro aqui es porque si leyo el MBR del disco
		if bytes.Equal(particion.Part_name[:], name[:]) {
			// comprobaremos que la particion no se haya montado previamente
			particionEncontrada = true
			if particion.Part_status == '2' {
				fmt.Printf("Partición ya montada con anterioridad.\n")
				return false
			}
			if particion.Part_type == 'e' || particion.Part_type == 'E' {
				fmt.Printf("Las particiones extendidas no pueden ser montadas.\n")
				return false
			}
			particion.Part_type = '2'
			var part *estructuras.Partition = new(estructuras.Partition)
			part = &particion
			lista.ListaMount.Mount(path, 04, part, nil)
			partitionMounted = true
			lista.ListaMount.PrintList()
			break

		}
	}
	if !particionEncontrada {
		// buscaremos si existe una particion logica con ese nombre
		for _, particion := range master.Mbr_partitions {
			if particion.Part_type == 'e' || particion.Part_type == 'E' {
				partitionMounted = true
				m.MountParticionLogica(path, int(particion.Part_start), name)
				// tener un metodo de Mount List que agregue un texto a la consola
				lista.ListaMount.PrintList()
			}
		}
	}
	if !partitionMounted {
		fmt.Printf("No se encontro la partición: %s\n", string(estructuras.TrimArray(name[:])))
		return false
	}
	WriteMBR(&master, path)
	return true
}

func (m *Mount) MountParticionLogica(path string, whereToStart int, name [16]byte) {
	logicPartitionMounted := false
	temp := ReadEBR(path, int64(whereToStart))
	flag := true
	for flag {
		if bytes.Equal(temp.Part_name[:], name[:]) {
			temp.Part_status = '2'
			var partL *estructuras.EBR = new(estructuras.EBR)
			partL = &temp
			lista.ListaMount.Mount(path, 04, nil, partL)
			logicPartitionMounted = true
			flag = false
			break
		} else if temp.Part_next != -1 {
			temp = ReadEBR(path, temp.Part_next)
		} else {
			flag = false
		}
	}
	if !logicPartitionMounted {
		fmt.Printf("No se encontro la partición: %s\n", string(estructuras.TrimArray(name[:])))
		return
	}
}
