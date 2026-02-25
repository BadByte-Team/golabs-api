// Package pagination provides helpers to parse and return paginated responses.
package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage = 1
	DefaultSize = 20
	MaxSize     = 100
)

// Page holds parsed pagination parameters.
type Page struct {
	Number int // 1-indexed
	Size   int
}

// Offset returns the SQL OFFSET value.
func (p Page) Offset() int {
	return (p.Number - 1) * p.Size
}

// Parse reads ?page= and ?size= from the request query string.
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
		size = MaxSize
	}
	return Page{Number: page, Size: size}
}

// Meta is embedded in paginated responses.
type Meta struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// Response is a generic paginated response envelope.
type Response[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// New builds a paginated response.
func New[T any](data []T, p Page, total int) Response[T] {
	if data == nil {
		data = []T{}
	}
	return Response[T]{Data: data, Meta: Meta{Page: p.Number, Size: p.Size, Total: total}}
}

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
