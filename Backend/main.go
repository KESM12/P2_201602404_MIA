package main

import (
	"P1/analizador"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

//var logued = false

type DatosEntrada struct {
	Comandos []string `json:"comandos"`
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", inicial).Methods("GET")
	router.HandleFunc("/", analizadorweb).Methods("POST")

	handler := allowCORS(router)
	fmt.Println("Server on port :4000")
	log.Fatal(http.ListenAndServe(":4000", handler))
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handler.ServeHTTP(w, r)
	})
}

func inicial(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>¡Hola Desde el servidor!</h1>")
}

func analizadorweb(w http.ResponseWriter, r *http.Request) {
	var analizador analizador.Analizador
	var datos DatosEntrada
	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = guardarDatos("./prueba.script", datos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ejecutar el archivo de script
	analizador.Execute("./prueba.script")
	fmt.Fprintf(w, "Script ejecutado exitosamente")
}

func guardarDatos(archivo string, datos DatosEntrada) error {
	// Abrir o crear el archivo
	file, err := os.Create(archivo)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir los comandos en el archivo
	for _, comando := range datos.Comandos {
		_, err := file.WriteString(strings.TrimSpace(comando) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// func main() {
// 	var analizador analizador.Analizador
// 	ejecutar := true
// 	fmt.Println("*********** Kevin Estuardo Secaida Molina ***********")
// 	for ejecutar {
// 		var opcion string
// 		fmt.Printf("\n")
// 		fmt.Printf("Comando: ")
// 		scanner := bufio.NewScanner(os.Stdin)
// 		scanner.Scan()
// 		opcion = scanner.Text()
// 		if opcion == "exit" {
// 			ejecutar = false
// 		} else {
// 			analizador.Execute(opcion)
// 		}
// 	}
// }
