package comandUser

import (
	"P1/comandos"
	"P1/estructuras"
	"P1/lista"
	"P1/logger"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type ParametrosMkgrp struct {
	Name string
}

type Mkgrp struct {
	params ParametrosMkgrp
}

func (m *Mkgrp) Execute(parametros []string) {
	m.params = m.SaveParams(parametros)
	if m.Mkgrp(m.params.Name) {
		fmt.Printf("\nGrupo \"%s\" creado con éxito\n\n", m.params.Name)
	} else {
		fmt.Printf("")
	}
}

func (m *Mkgrp) SaveParams(parametros []string) ParametrosMkgrp {
	for _, v := range parametros {
		// fmt.Println(v)
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, " ")
		v = strings.ReplaceAll(v, "\"", "")
		if strings.Contains(v, "name") {
			v = strings.ReplaceAll(v, "name=", "")
			m.params.Name = v
		}
	}
	return m.params
}

func (m *Mkgrp) Mkgrp(name string) bool {
	if name == "" {
		fmt.Println("Parametro name no encontrado.")
		return true
	}
	if logger.Log.IsLoggedIn() && logger.Log.UserIsRoot() {
		if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value.Part_start, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		} else if lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Value != nil {
			return m.MkgrpPartition(name, lista.ListaMount.GetNodeById(logger.Log.GetUserId()).ValueL.Part_start+int64(unsafe.Sizeof(estructuras.EBR{})), lista.ListaMount.GetNodeById(logger.Log.GetUserId()).Ruta)
		}
	}
	return false
}

func (m *Mkgrp) MkgrpPartition(name string, whereToStart int64, path string) bool {
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
	if m.ExisteGrupo(ReadFile(&tablaInodo, path, &superbloque), name) {
		fmt.Printf("Grupo ya existente: \"%s\"\n", name)
		return false
	}
	numero := m.ContarGrupos(ReadFile(&tablaInodo, path, &superbloque))
	grupo := m.AgregarGrupo(numero, name)
	if AppendFile(path, &superbloque, &tablaInodo, grupo) {
		comandos.Fwrite(&tablaInodo, path, superbloque.S_inode_start+int64(unsafe.Sizeof(estructuras.TablaInodo{})))
		fmt.Println(ReadFile(&tablaInodo, path, &superbloque))
		comandos.Fwrite(&superbloque, path, whereToStart)
		return true
	}
	return false
}

func (m *Mkgrp) AgregarGrupo(groupNumber int, groupName string) string {
	return strconv.Itoa(groupNumber) + ",G," + groupName + "\n"
}

func (m *Mkgrp) ContarGrupos(contenido string) int {
	contador := 1
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "G" {
			continue
		}
		contador++
	}
	return contador
}

func (m *Mkgrp) ExisteGrupo(contenido string, groupName string) bool {
	lineas := strings.Split(contenido, "\n")
	lineas = lineas[:len(lineas)-1]
	for _, linea := range lineas {
		parametros := strings.Split(linea, ",")
		if parametros[1] != "G" {
			continue
		}
		if parametros[2] == groupName {
			return true
		}
	}
	return false
}
