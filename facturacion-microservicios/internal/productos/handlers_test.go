package productos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"facturacion-microservicios/internal/jsondb"
)

func nuevoHandlerPrueba(t *testing.T) (*Handler, *jsondb.Store[Producto]) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "productos.json")
	store, err := jsondb.New(path, 101,
		func(p *Producto) int { return p.ID },
		func(p *Producto, id int) { p.ID = id },
	)
	if err != nil {
		t.Fatalf("no se pudo crear el store: %v", err)
	}

	seed := []Producto{
		{Nombre: "Laptop Lenovo ThinkPad", Precio: 3500000, Stock: 14},
		{Nombre: "Mouse inalámbrico Logitech", Precio: 12000, Stock: 47},
	}
	for _, p := range seed {
		if _, err := store.Create(p); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	return NewHandler(store), store
}

func muxDe(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func solicitar(t *testing.T, mux *http.ServeMux, method, path string, body any) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()

	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var decoded map[string]any
	if rec.Body.Len() > 0 && bytes.HasPrefix(bytes.TrimSpace(rec.Body.Bytes()), []byte("{")) {
		if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
			t.Fatalf("respuesta no es JSON: %v (%s)", err, rec.Body.String())
		}
	}

	return rec, decoded
}

func stockDe(t *testing.T, store *jsondb.Store[Producto], id int) int {
	t.Helper()

	p, err := store.GetByID(id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p == nil {
		t.Fatalf("producto %d no existe", id)
	}
	return p.Stock
}

func TestListarYConsultarProductos(t *testing.T) {
	handler, _ := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, _ := solicitar(t, mux, http.MethodGet, "/productos", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}

	var lista []Producto
	if err := json.Unmarshal(rec.Body.Bytes(), &lista); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("se esperaban 2 productos, hay %d", len(lista))
	}

	rec, body := solicitar(t, mux, http.MethodGet, "/productos/101", nil)
	if rec.Code != http.StatusOK || body["nombre"] != "Laptop Lenovo ThinkPad" {
		t.Fatalf("detalle inesperado (%d): %s", rec.Code, rec.Body.String())
	}

	rec, _ = solicitar(t, mux, http.MethodGet, "/productos/999", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestDescontarStockDescuentaLote(t *testing.T) {
	handler, store := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, _ := solicitar(t, mux, http.MethodPost, "/productos/descontar-stock", DescontarStockRequest{
		Items: []DescontarStockItem{
			{ProductoID: 101, Cantidad: 1},
			{ProductoID: 102, Cantidad: 2},
		},
	})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	if got := stockDe(t, store, 101); got != 13 {
		t.Fatalf("stock 101 esperado 13, obtenido %d", got)
	}
	if got := stockDe(t, store, 102); got != 45 {
		t.Fatalf("stock 102 esperado 45, obtenido %d", got)
	}
}

func TestDescontarStockEsAtomicoSiFaltaStock(t *testing.T) {
	handler, store := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, _ := solicitar(t, mux, http.MethodPost, "/productos/descontar-stock", DescontarStockRequest{
		Items: []DescontarStockItem{
			{ProductoID: 101, Cantidad: 5},
			{ProductoID: 102, Cantidad: 999},
		},
	})

	if rec.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	if got := stockDe(t, store, 101); got != 14 {
		t.Fatalf("el stock 101 no debió cambiar, está en %d", got)
	}
	if got := stockDe(t, store, 102); got != 47 {
		t.Fatalf("el stock 102 no debió cambiar, está en %d", got)
	}
}

func TestDescontarStockProductoInexistenteNoModificaNada(t *testing.T) {
	handler, store := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, _ := solicitar(t, mux, http.MethodPost, "/productos/descontar-stock", DescontarStockRequest{
		Items: []DescontarStockItem{
			{ProductoID: 101, Cantidad: 1},
			{ProductoID: 999, Cantidad: 1},
		},
	})

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	if got := stockDe(t, store, 101); got != 14 {
		t.Fatalf("el stock 101 no debió cambiar, está en %d", got)
	}
}

func TestDescontarStockCombinaItemsRepetidos(t *testing.T) {
	handler, store := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, _ := solicitar(t, mux, http.MethodPost, "/productos/descontar-stock", DescontarStockRequest{
		Items: []DescontarStockItem{
			{ProductoID: 102, Cantidad: 2},
			{ProductoID: 102, Cantidad: 3},
		},
	})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	if got := stockDe(t, store, 102); got != 42 {
		t.Fatalf("stock 102 esperado 42, obtenido %d", got)
	}
}

func TestDescontarStockValidaCantidades(t *testing.T) {
	handler, store := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	casos := []DescontarStockRequest{
		{},
		{Items: []DescontarStockItem{{ProductoID: 102, Cantidad: 0}}},
		{Items: []DescontarStockItem{{ProductoID: 102, Cantidad: -1}}},
		{Items: []DescontarStockItem{{ProductoID: 0, Cantidad: 1}}},
	}

	for _, caso := range casos {
		rec, _ := solicitar(t, mux, http.MethodPost, "/productos/descontar-stock", caso)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("se esperaba 400 para %+v, se obtuvo %d: %s", caso, rec.Code, rec.Body.String())
		}
	}
	if got := stockDe(t, store, 102); got != 47 {
		t.Fatalf("el stock no debió cambiar, está en %d", got)
	}
}

func TestCrearYActualizarProducto(t *testing.T) {
	handler, _ := nuevoHandlerPrueba(t)
	mux := muxDe(handler)

	rec, body := solicitar(t, mux, http.MethodPost, "/productos", CreateRequest{
		Nombre: "  Teclado mecánico  ",
		Precio: 150000,
		Stock:  5,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
	if body["id"] != float64(103) || body["nombre"] != "Teclado mecánico" {
		t.Fatalf("producto inesperado: %v", body)
	}

	rec, _ = solicitar(t, mux, http.MethodPut, "/productos/103", UpdateRequest{
		Nombre: "Teclado mecánico RGB",
		Precio: 175000,
		Stock:  3,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}

	rec, _ = solicitar(t, mux, http.MethodDelete, "/productos/103", nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", rec.Code)
	}
}
