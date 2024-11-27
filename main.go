package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Pelicula struct {
	ID       int    `json:"id"`
	Titulo   string `json:"titulo"`
	Director string `json:"director"`
	Anio     int    `json:"anio"`
}

var peliculas []Pelicula
var archivoJSON = "peliculas.json"

// Leer archivo JSON
func cargarPeliculas() error {
	data, err := os.ReadFile(archivoJSON)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &peliculas)
}

// Guardar en archivo JSON
func guardarPeliculas() error {
	data, err := json.MarshalIndent(peliculas, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(archivoJSON, data, 0644)
}

// Renderizar plantillas
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t := template.Must(template.ParseFiles("templates/" + tmpl))
	t.Execute(w, data)
}

func listaPeliculas(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", peliculas)
}

func editarPelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(idStr)
	for _, p := range peliculas {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Película no encontrada", http.StatusNotFound)
}

func guardarPelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	titulo := r.FormValue("titulo")
	director := r.FormValue("director")
	anioStr := r.FormValue("anio")

	id, _ := strconv.Atoi(idStr)
	anio, _ := strconv.Atoi(anioStr)

	if id == 0 {
		id = len(peliculas) + 1 // acumulador de ID incrementales
		peliculas = append(peliculas, Pelicula{ID: id, Titulo: titulo, Director: director, Anio: anio})
	} else {
		for i, p := range peliculas {
			if p.ID == id {
				peliculas[i] = Pelicula{ID: id, Titulo: titulo, Director: director, Anio: anio}
				break
			}
		}
	}

	guardarPeliculas()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func borrarPelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	for i, p := range peliculas {
		if p.ID == id {
			peliculas = append(peliculas[:i], peliculas[i+1:]...)
			break
		}
	}

	guardarPeliculas()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	if err := cargarPeliculas(); err != nil {
		log.Fatalf("Error al cargar las películas: %v", err)
	}

	http.HandleFunc("/", listaPeliculas)
	http.HandleFunc("/nuevo", editarPelicula)
	http.HandleFunc("/editar", editarPelicula)
	http.HandleFunc("/guardar", guardarPelicula)
	http.HandleFunc("/borrar", borrarPelicula)

	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
