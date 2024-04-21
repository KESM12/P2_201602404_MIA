package analizador

import (
	"P1/comandos"
	"P1/comandos/comandUser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Analizador struct {
}

func (a *Analizador) Execute(input string) {
	ComandosyParametros := a.Split_input(input)
	var comandos string
	var parametros []string
	for i, v := range ComandosyParametros {
		if i == 0 {
			comandos = v
		} else {
			parametros = append(parametros, v)
		}
	}
	a.MatchParams(comandos, parametros)
}

func (a *Analizador) MatchParams(command string, params []string) {
	command = strings.ReplaceAll(command, " ", "")
	if command == "execute" {
		for _, v := range params {
			if strings.Contains(v, "path") {
				v = strings.Replace(v, "path=", "", 1)
				v = strings.ReplaceAll(v, "\"", "")
				a.Read(v)
			}
		}
	} else if command == "pause" {
		var option string
		fmt.Println("presione 'ENTER' para continuar: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option = scanner.Text()
		fmt.Println(option)
	} else if command == "mkdisk" {
		m := comandos.Mkdisk{}
		m.Execute(params)
	} else if command == "rmdisk" {
		r := comandos.Rmdisk{}
		r.Execute(params)
	} else if command == "fdisk" {
		f := comandos.Fdisk{}
		f.Execute(params)
	} else if command == "mount" {
		m := comandos.Mount{}
		m.Execute(params)
	} else if command == "unmount" {
		m := comandos.Unmount{}
		m.Execute(params)
	} else if command == "mkfs" {
		m := comandos.Mkfs{}
		m.Execute(params)
	} else if command == "login" {
		l := comandUser.Login{}
		l.Execute(params)
	} else if command == "logout" {
		l := comandUser.Logout{}
		l.Execute(params)
	} else if command == "mkgrp" {
		m := comandUser.Mkgrp{}
		m.Execute(params)
	} else if command == "rmgrp" {
		r := comandUser.Rmgrp{}
		r.Execute(params)
	} else if command == "mkusr" {
		m := comandUser.Mkusr{}
		m.Execute(params)
	} else if command == "rmusr" {
		r := comandUser.Rmusr{}
		r.Execute(params)
	} else if command == "mkfile" {
		m := comandUser.Mkfile{}
		m.Execute(params)
	} else if command == "mkdir" {
		m := comandUser.Mkdir{}
		m.Execute(params)
	} else if command == "rep" {
		r := comandUser.Rep{}
		r.Execute(params)
	} else if strings.Contains(command, "#") {
		contenido := command
		contenido += "\n\n"
		fmt.Println(contenido)
	}

}

func (a *Analizador) Split_input(comando string) []string {
	return strings.Split(comando, "-")
}

func (a *Analizador) Read(ruta string) {
	file, err := os.Open(ruta)
	if err != nil {
		fmt.Println("No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		comando := scanner.Text() // Lee la línea actual como un comando completo
		a.Execute(comando)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}
