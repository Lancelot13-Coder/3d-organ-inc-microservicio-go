// Microservicio en Go (4to lenguaje) para 3D-Organ-Inc.
//
// Hace lo mismo que los otros tres (Python, Node.js y Java): CRUD
// completo (crear, editar, eliminar) sobre la tabla 'modelos' en
// Supabase. Usa solo la librería estándar de Go para el servidor HTTP
// (net/http) y el driver "lib/pq" para hablar con Postgres.
//
//	GET    /api/modelos       -> listar (opcional ?categoria=...)
//	GET    /api/modelos/{id}  -> obtener uno
//	POST   /api/modelos       -> crear
//	PUT    /api/modelos/{id}  -> editar
//	DELETE /api/modelos/{id}  -> eliminar
//	GET    /                  -> confirma que el servicio está vivo
//
// Variables de entorno (se configuran en Render):
//
//	DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
//	PORT (Render la define sola)
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Modelo3D struct {
	ID            int       `json:"id"`
	Nombre        string    `json:"nombre"`
	Descripcion   string    `json:"descripcion"`
	Categoria     string    `json:"categoria"`
	FechaRegistro time.Time `json:"fecha_registro"`
	Activo        bool      `json:"activo"`
}

func main() {
	host := os.Getenv("DB_HOST")
	puertoDb := getenvDefault("DB_PORT", "5432")
	nombreDb := getenvDefault("DB_NAME", "postgres")
	usuario := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=require",
		host, puertoDb, nombreDb, usuario, password,
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("no se pudo preparar la conexión a la base de datos: %v", err)
	}

	http.HandleFunc("/", manejarRaiz)
	http.HandleFunc("/api/modelos", manejarModelos)
	http.HandleFunc("/api/modelos/", manejarModelos)
	registrarDocs() // Swagger en /docs

	puerto := getenvDefault("PORT", "8080")
	log.Printf("Microservicio Go escuchando en el puerto %s", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}

func getenvDefault(clave, porDefecto string) string {
	valor := os.Getenv(clave)
	if valor == "" {
		return porDefecto
	}
	return valor
}

func agregarCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func enviarJSON(w http.ResponseWriter, codigo int, cuerpo interface{}) {
	agregarCors(w)
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(cuerpo)
}

func enviarError(w http.ResponseWriter, codigo int, detalle string) {
	enviarJSON(w, codigo, map[string]string{"detail": detalle})
}

func manejarRaiz(w http.ResponseWriter, r *http.Request) {
	enviarJSON(w, 200, map[string]string{
		"servicio": "microservicio-modelos-3d-go",
		"estado":   "ok",
	})
}

func manejarModelos(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		agregarCors(w)
		w.WriteHeader(200)
		return
	}

	if r.Method == http.MethodGet && r.URL.Path == "/api/modelos" {
		listar(w, r)
		return
	}

	if r.Method == http.MethodPost && r.URL.Path == "/api/modelos" {
		crear(w, r)
		return
	}

	prefijo := "/api/modelos/"
	if strings.HasPrefix(r.URL.Path, prefijo) {
		idTexto := strings.TrimPrefix(r.URL.Path, prefijo)
		id, err := strconv.Atoi(idTexto)
		if err != nil {
			enviarError(w, 400, "El id debe ser un número")
			return
		}
		switch r.Method {
		case http.MethodGet:
			obtener(w, id)
			return
		case http.MethodPut:
			editar(w, r, id)
			return
		case http.MethodDelete:
			eliminar(w, id)
			return
		}
	}

	enviarError(w, 405, "Método no permitido o ruta inválida")
}

func listar(w http.ResponseWriter, r *http.Request) {
	consulta := `SELECT id, nombre, COALESCE(descripcion, ''), categoria, fecha_registro, activo FROM modelos`
	var args []interface{}
	if categoria := r.URL.Query().Get("categoria"); categoria != "" {
		consulta += " WHERE categoria = $1"
		args = append(args, categoria)
	}
	consulta += " ORDER BY id"

	filas, err := db.Query(consulta, args...)
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}
	defer filas.Close()

	modelos := []Modelo3D{}
	for filas.Next() {
		var m Modelo3D
		if err := filas.Scan(&m.ID, &m.Nombre, &m.Descripcion, &m.Categoria, &m.FechaRegistro, &m.Activo); err != nil {
			enviarError(w, 500, "Error de base de datos: "+err.Error())
			return
		}
		modelos = append(modelos, m)
	}
	enviarJSON(w, 200, modelos)
}

func obtener(w http.ResponseWriter, id int) {
	var m Modelo3D
	consulta := `SELECT id, nombre, COALESCE(descripcion, ''), categoria, fecha_registro, activo
	             FROM modelos WHERE id = $1`
	err := db.QueryRow(consulta, id).
		Scan(&m.ID, &m.Nombre, &m.Descripcion, &m.Categoria, &m.FechaRegistro, &m.Activo)
	if err == sql.ErrNoRows {
		enviarError(w, 404, "Modelo 3D no encontrado")
		return
	}
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}
	enviarJSON(w, 200, m)
}

type entradaCrear struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Categoria   string `json:"categoria"`
	Activo      *bool  `json:"activo"`
}

func crear(w http.ResponseWriter, r *http.Request) {
	var entrada entradaCrear
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		enviarError(w, 400, "JSON inválido en el cuerpo de la petición")
		return
	}

	if entrada.Nombre == "" || entrada.Categoria == "" {
		enviarError(w, 422, "Los campos 'nombre' y 'categoria' son obligatorios")
		return
	}

	activo := true
	if entrada.Activo != nil {
		activo = *entrada.Activo
	}

	var m Modelo3D
	sql := `INSERT INTO modelos (nombre, descripcion, categoria, activo)
	        VALUES ($1, $2, $3, $4)
	        RETURNING id, nombre, descripcion, categoria, fecha_registro, activo`

	err := db.QueryRow(sql, entrada.Nombre, entrada.Descripcion, entrada.Categoria, activo).
		Scan(&m.ID, &m.Nombre, &m.Descripcion, &m.Categoria, &m.FechaRegistro, &m.Activo)
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}

	enviarJSON(w, 201, m)
}

func editar(w http.ResponseWriter, r *http.Request, id int) {
	var cambios map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&cambios); err != nil {
		enviarError(w, 400, "JSON inválido en el cuerpo de la petición")
		return
	}

	campos := []string{"nombre", "descripcion", "categoria", "activo"}
	var asignaciones []string
	var valores []interface{}
	contador := 1
	for _, campo := range campos {
		if valor, existe := cambios[campo]; existe {
			asignaciones = append(asignaciones, fmt.Sprintf("%s = $%d", campo, contador))
			valores = append(valores, valor)
			contador++
		}
	}

	if len(asignaciones) == 0 {
		enviarError(w, 400, "No se envió ningún campo para actualizar")
		return
	}

	valores = append(valores, id)
	sqlUpdate := fmt.Sprintf(
		"UPDATE modelos SET %s WHERE id = $%d",
		strings.Join(asignaciones, ", "), contador,
	)

	resultado, err := db.Exec(sqlUpdate, valores...)
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}
	filas, _ := resultado.RowsAffected()
	if filas == 0 {
		enviarError(w, 404, "Modelo 3D no encontrado")
		return
	}

	var m Modelo3D
	sqlSelect := `SELECT id, nombre, COALESCE(descripcion, ''), categoria, fecha_registro, activo
	              FROM modelos WHERE id = $1`
	err = db.QueryRow(sqlSelect, id).
		Scan(&m.ID, &m.Nombre, &m.Descripcion, &m.Categoria, &m.FechaRegistro, &m.Activo)
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}

	enviarJSON(w, 200, m)
}

func eliminar(w http.ResponseWriter, id int) {
	resultado, err := db.Exec("DELETE FROM modelos WHERE id = $1", id)
	if err != nil {
		enviarError(w, 500, "Error de base de datos: "+err.Error())
		return
	}
	filas, _ := resultado.RowsAffected()
	if filas == 0 {
		enviarError(w, 404, "Modelo 3D no encontrado")
		return
	}
	agregarCors(w)
	w.WriteHeader(204)
}