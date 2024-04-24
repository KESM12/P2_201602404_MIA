package main

import (
	"P1/analizador"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
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

	// Agregar el middleware CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(http.HandlerFunc(router.ServeHTTP))

	fmt.Println("Server on port :4000")
	http.ListenAndServe(":4000", corsHandler)
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
	// Crear un mensaje de respuesta
	mensaje := struct {
		Mensaje string `json:"mensaje"`
	}{
		Mensaje: "Script ejecutado exitosamente",
	}
	// Convertir el mensaje a JSON
	respuesta, err := json.Marshal(mensaje)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado de contenido JSON
	w.Header().Set("Content-Type", "application/json")

	// Escribir la respuesta JSON en la respuesta HTTP
	fmt.Fprintf(w, string(respuesta))
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
