// Package pagination provee helpers para parsear parametros de paginacion desde
// el query string de la peticion HTTP y construir respuestas paginadas estandarizadas.
//
// Uso tipico en un handler:
//
//	p := pagination.Parse(r)
//	items, total, _ := repo.List(p.Offset(), p.Size)
//	apperrors.RespondJSON(w, 200, pagination.New(items, p, total))
package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage = 1   // pagina por defecto si no se especifica ?page=
	DefaultSize = 20  // tamaño de pagina por defecto si no se especifica ?size=
	MaxSize     = 100 // tamaño maximo permitido para prevenir consultas masivas
)

// Page contiene los parametros de paginacion ya parseados y validados.
type Page struct {
	Number int // numero de pagina, con base 1 (la primera pagina es 1)
	Size   int // cantidad de registros por pagina
}

// Offset calcula el valor SQL OFFSET correspondiente a esta pagina.
// Se usa directamente en queries: SELECT ... LIMIT size OFFSET offset.
func (p Page) Offset() int {
	return (p.Number - 1) * p.Size
}

// Parse lee los parametros ?page= y ?size= del query string de la peticion.
// Aplica valores por defecto y limites para evitar paginas invalidas o demasiado grandes.
//
// Entrada:  r, peticion HTTP con query string.
// Salida:   Page con valores ya validados y dentro de los limites permitidos.
func Parse(r *http.Request) Page {
	page := queryInt(r, "page", DefaultPage)
	size := queryInt(r, "size", DefaultSize)

	if page < 1 {
		page = DefaultPage
	}
	if size < 1 {
		size = DefaultSize
	}
	if size > MaxSize {
		// Limitar para protejer contra consultas que devuelvan demasiados registros.
		size = MaxSize
	}
	return Page{Number: page, Size: size}
}

// Meta contiene los metadatos de paginacion incluidos en todas las respuestas paginadas.
type Meta struct {
	Page  int `json:"page"`  // numero de pagina actual
	Size  int `json:"size"`  // tamaño de pagina devuelto
	Total int `json:"total"` // total de registros disponibles (sin paginar)
}

// Response es el envelope generico para respuestas paginadas.
// Data contiene los registros de la pagina actual; Meta tiene la informacion de paginacion.
type Response[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// New construye una Response paginada a partir de los datos y los parametros de pagina.
// Si data es nil, se normaliza a slice vacio para evitar null en el JSON.
//
// Entrada:  data, registros de la pagina; p, parametros de pagina; total, total de registros.
// Salida:   Response lista para serializar como JSON.
func New[T any](data []T, p Page, total int) Response[T] {
	if data == nil {
		data = []T{}
	}
	return Response[T]{Data: data, Meta: Meta{Page: p.Number, Size: p.Size, Total: total}}
}

// queryInt lee un parametro del query string como entero.
// Retorna def si el parametro no existe o no puede parsearse.
func queryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
