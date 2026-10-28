package jsondb

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

type registro struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

func nuevoStore(t *testing.T) (*Store[registro], string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "registros.json")
	store, err := New(path, 1,
		func(r *registro) int { return r.ID },
		func(r *registro, id int) { r.ID = id },
	)
	if err != nil {
		t.Fatalf("no se pudo crear el store: %v", err)
	}

	return store, path
}

func TestCreateAsignaIDsSecuenciales(t *testing.T) {
	store, _ := nuevoStore(t)

	primero, err := store.Create(registro{Nombre: "uno"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	segundo, err := store.Create(registro{Nombre: "dos"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if primero.ID != 1 || segundo.ID != 2 {
		t.Fatalf("IDs esperados 1 y 2, obtenidos %d y %d", primero.ID, segundo.ID)
	}
}

func TestTransactConfirmaCambios(t *testing.T) {
	store, _ := nuevoStore(t)
	creado, err := store.Create(registro{Nombre: "original"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	err = store.Transact(func(items []registro) error {
		for i := range items {
			if items[i].ID == creado.ID {
				items[i].Nombre = "actualizado"
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transact: %v", err)
	}

	got, err := store.GetByID(creado.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.Nombre != "actualizado" {
		t.Fatalf("cambio no persistido: %+v", got)
	}
}

func TestTransactRevierteSiFnFalla(t *testing.T) {
	store, _ := nuevoStore(t)
	creado, err := store.Create(registro{Nombre: "original"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	errEsperado := errors.New("fallo simulado")
	err = store.Transact(func(items []registro) error {
		items[0].Nombre = "no debería guardarse"
		return errEsperado
	})
	if !errors.Is(err, errEsperado) {
		t.Fatalf("se esperaba el error simulado, se obtuvo: %v", err)
	}

	got, err := store.GetByID(creado.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.Nombre != "original" {
		t.Fatalf("el cambio debió revertirse, se obtuvo: %+v", got)
	}
}

func TestCreateConcurrenteGeneraIDsUnicos(t *testing.T) {
	store, _ := nuevoStore(t)

	const total = 25
	ids := make(chan int, total)
	var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			item, err := store.Create(registro{Nombre: fmt.Sprintf("item-%d", n)})
			if err != nil {
				t.Errorf("create concurrente: %v", err)
				return
			}
			ids <- item.ID
		}(i)
	}

	wg.Wait()
	close(ids)

	vistos := make(map[int]bool, total)
	for id := range ids {
		if vistos[id] {
			t.Fatalf("ID duplicado: %d", id)
		}
		vistos[id] = true
	}
	if len(vistos) != total {
		t.Fatalf("se esperaban %d IDs, se obtuvieron %d", total, len(vistos))
	}
}

func TestNextIDContinuaDesdeElArchivoExistente(t *testing.T) {
	store, path := nuevoStore(t)
	if _, err := store.Create(registro{Nombre: "uno"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := store.Create(registro{Nombre: "dos"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	reabierto, err := New(path, 1,
		func(r *registro) int { return r.ID },
		func(r *registro, id int) { r.ID = id },
	)
	if err != nil {
		t.Fatalf("reabrir: %v", err)
	}

	nuevo, err := reabierto.Create(registro{Nombre: "tres"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if nuevo.ID != 3 {
		t.Fatalf("se esperaba ID 3, se obtuvo %d", nuevo.ID)
	}
}
