package pagination

import (
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"
)

// Cursor opaco (base64)
type Cursor string

// Pagination parámetros de entrada
type Pagination struct {
	Limit      int
	NextCursor *Cursor
	PrevCursor *Cursor
	OrderBy    string // columna por la que se pagina (ej: "id" o "created_at")
	SortAsc    bool   // true = ASC, false = DESC
}

// PaginationOpts opciones para ejecutar la query
type PaginationOpts[T any] struct {
	DB         *gorm.DB
	Where      interface{}
	Result     *[]T // slice destino
	TableModel interface{}
}

// Page salida paginada tipada
type Page[T any] struct {
	Items       []T     `json:"items"`
	NextCursor  *Cursor `json:"next_cursor,omitempty"`
	PrevCursor  *Cursor `json:"prev_cursor,omitempty"`
	HasNext     bool    `json:"has_next"`
	HasPrevious bool    `json:"has_previous"`
}

// EncodeCursor convierte string → cursor
func EncodeCursor(value string) Cursor {
	return Cursor(base64.StdEncoding.EncodeToString([]byte(value)))
}

// DecodeCursor convierte cursor → string
func DecodeCursor(c Cursor) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(c))
	if err != nil {
		return "", fmt.Errorf("invalid cursor: %w", err)
	}
	return string(decoded), nil
}

// helper: invertir slice in-place
func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Paginate función genérica
func Paginate[T any](p Pagination, opts PaginationOpts[T], extractCursor func(item T) string) (Page[T], error) {
	db := opts.DB.Model(opts.TableModel)

	// filtros extra
	if opts.Where != nil {
		db = db.Where(opts.Where)
	}

	// orden base
	order := p.OrderBy
	if !p.SortAsc {
		order += " DESC"
	}

	// nextCursor
	if p.NextCursor != nil {
		cursorVal, err := DecodeCursor(*p.NextCursor)
		if err != nil {
			return Page[T]{}, err
		}
		if p.SortAsc {
			db = db.Where(fmt.Sprintf("%s > ?", p.OrderBy), cursorVal)
		} else {
			db = db.Where(fmt.Sprintf("%s < ?", p.OrderBy), cursorVal)
		}
	}

	// prevCursor (trae hacia atrás y cambia el orden para “retroceder”)
	if p.PrevCursor != nil {
		cursorVal, err := DecodeCursor(*p.PrevCursor)
		if err != nil {
			return Page[T]{}, err
		}
		if p.SortAsc {
			db = db.Where(fmt.Sprintf("%s < ?", p.OrderBy), cursorVal)
			order = p.OrderBy + " DESC"
		} else {
			db = db.Where(fmt.Sprintf("%s > ?", p.OrderBy), cursorVal)
			order = p.OrderBy + " ASC"
		}
	}

	// aplicar orden y límite (+1 para detectar si hay más)
	db = db.Order(order).Limit(p.Limit + 1)

	// ejecutar
	if err := db.Find(opts.Result).Error; err != nil {
		return Page[T]{}, err
	}

	itemsFull := *opts.Result
	page := Page[T]{}

	// hay más si vino limit+1
	hasMore := len(itemsFull) > p.Limit
	if hasMore {
		itemsFull = itemsFull[:p.Limit]
		page.HasNext = (p.PrevCursor == nil) // si ibas hacia atrás, "más" significa páginas previas, no next
	}

	// si fue prevCursor invertimos para devolver orden estable (el “natural” según SortAsc)
	if p.PrevCursor != nil {
		reverseSlice(itemsFull)
		// Nota: podrías setear HasPrevious=true si querés indicar que hay más hacia atrás;
		// aquí lo marcamos abajo con p.PrevCursor != nil.
	}

	page.Items = itemsFull

	// cursores a partir de los items retornados (ya recortados e invertidos si correspondía)
	if len(page.Items) > 0 {
		first := page.Items[0]
		last := page.Items[len(page.Items)-1]

		firstCursor := EncodeCursor(extractCursor(first))
		lastCursor := EncodeCursor(extractCursor(last))

		page.PrevCursor = &firstCursor // para ir hacia atrás desde el primero mostrado
		page.NextCursor = &lastCursor  // para ir hacia adelante desde el último mostrado
	}

	// heurística simple para flags (ajústalo si querés comportamiento distinto)
	page.HasPrevious = p.PrevCursor != nil

	return page, nil
}
