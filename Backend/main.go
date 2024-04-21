package main

import (
	"P1/analizador"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var analizador analizador.Analizador
	ejecutar := true
	fmt.Println("*********** Kevin Estuardo Secaida Molina ***********")
	for ejecutar {
		var opcion string
		fmt.Printf("\n")
		fmt.Printf("Comando: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		opcion = scanner.Text()
		if opcion == "exit" {
			ejecutar = false
		} else {
			analizador.Execute(opcion)
		}
	}
}
