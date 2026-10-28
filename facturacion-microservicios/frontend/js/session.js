// Sesión del usuario sin tokens: solo se guarda el objeto usuario en localStorage.

const STORAGE_KEY = "facturacion_usuario";

export function getUsuario() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function setUsuario(usuario) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(usuario));
}

export function clearUsuario() {
  localStorage.removeItem(STORAGE_KEY);
}

/**
 * Guardia de navegación. No es seguridad real: solo evita mostrar las
 * pantallas internas cuando no se ha iniciado sesión.
 */
export function requireSession() {
  const usuario = getUsuario();
  if (!usuario) {
    window.location.replace("index.html");
    return null;
  }
  return usuario;
}

export function redirectIfLoggedIn() {
  if (getUsuario()) {
    window.location.replace("productos.html");
  }
}
