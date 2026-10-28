# Sistema de Facturación Digital (Microservicios)

Proyecto de Programación Avanzada 2026_2. Sistema de facturación basado en arquitectura de microservicios desarrollado en Go utilizando exclusivamente la librería estándar `net/http`, base de datos persistente en archivos JSON y una interfaz de usuario web multipágina desarrollada en JavaScript.

El código y los microservicios se encuentran en la carpeta `facturacion-microservicios/`.

> **Alcance académico:** el sistema no usa JWT, tokens, sesiones, cookies, OAuth ni roles. El login valida las credenciales y el frontend guarda únicamente los datos del usuario en `localStorage` para personalizar la interfaz. Las contraseñas se almacenan en texto plano por decisión del equipo; **no es un diseño apto para producción**.

---

## Requisitos

- **Go 1.22** o superior instalado en el sistema.
- Navegador web moderno (Chrome, Edge, Firefox, Safari) con conexión a Internet para cargar Tailwind CDN y las fuentes.

---

## Puertos y Servicios

| Servicio | Puerto | Descripción | Archivo de Persistencia |
| :--- | :---: | :--- | :--- |
| **Autenticación** | `8081` | Registro, login y consulta de usuarios (sin tokens ni roles) | `data/usuarios.json` |
| **Productos** | `8082` | CRUD de catálogo de productos y descuento atómico de stock | `data/productos.json` |
| **Clientes** | `8083` | CRUD de clientes con nombre, cédula, correo y dirección | `data/clientes.json` |
| **Facturación** | `8084` | Orquestador que consume los 3 anteriores y emite facturas | `data/facturas.json` |
| **Frontend Web** | `3000` | Servidor web estático que sirve la interfaz de usuario | Carpeta `frontend/` |

---

## Credenciales de Prueba

Usuarios registrados en `data/usuarios.json` (contraseña en texto plano, sin campo `rol`):

- `ricardo@example.com` / `admin123`
- `maria@example.com` / `cliente123`

Ambos usuarios tienen exactamente las mismas capacidades dentro del sistema.

---

## Ejecución

Abrir una terminal por servicio desde la carpeta `facturacion-microservicios/`:

```bash
# Terminal 1 - Autenticación
go run ./cmd/auth

# Terminal 2 - Productos
go run ./cmd/productos

# Terminal 3 - Clientes
go run ./cmd/clientes

# Terminal 4 - Facturación
go run ./cmd/facturacion

# Terminal 5 - Frontend Web
go run ./cmd/web
```

Abrir en el navegador: **http://localhost:3000**

---

## Pantallas del Frontend

| Página | Descripción |
| :--- | :--- |
| `index.html` | Inicio de sesión |
| `registro.html` | Creación de cuenta |
| `productos.html` | Catálogo de productos (solo lectura) |
| `clientes.html` | Directorio de clientes y modal de creación |
| `nueva-factura.html` | Selección de cliente, carrito y validación de stock |
| `facturas.html` | Historial de facturas |
| `factura.html` | Detalle de factura, modal JSON e impresión / PDF |

---

## Pruebas

Desde `facturacion-microservicios/`:

```bash
go build ./...
go vet ./...
go test ./...
```

---

## Documentación

- `facturacion-microservicios/README.md`: detalle de endpoints, flujos y estructura.
- `facturacion-microservicios/docs/DOCUMENTO_TECNICO.md`: documento técnico de arquitectura y flujos.
