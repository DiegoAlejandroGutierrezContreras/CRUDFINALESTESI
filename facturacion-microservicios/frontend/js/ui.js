// Utilidades de interfaz compartidas.

const currencyFormatter = new Intl.NumberFormat("es-CO", {
  style: "currency",
  currency: "COP",
  maximumFractionDigits: 0,
});

/** Formatea un número como pesos colombianos: $3.500.000 */
export function formatCurrency(value) {
  const amount = Number(value) || 0;
  return currencyFormatter.format(amount).replace(/\u00a0/g, " ");
}

/** Escapa texto para insertarlo en plantillas HTML. */
export function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

/** Iniciales de un nombre: "Carlos Pérez" -> "CP". */
export function initials(name) {
  return String(name || "")
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part.charAt(0).toUpperCase())
    .join("");
}

function getToastContainer() {
  let container = document.querySelector(".toast-container");
  if (!container) {
    container = document.createElement("div");
    container.className =
      "toast-container fixed bottom-4 right-4 z-[100] flex w-[min(360px,calc(100vw-2rem))] flex-col gap-2";
    document.body.appendChild(container);
  }
  return container;
}

/**
 * Muestra una notificación temporal.
 * @param {string} message Texto a mostrar.
 * @param {"success"|"error"|"info"} type Tipo de notificación.
 */
export function showToast(message, type = "info") {
  const container = getToastContainer();
  const toast = document.createElement("div");
  const palette = {
    success: "bg-tertiary-container text-on-tertiary border-tertiary-container",
    error: "bg-error-container text-on-error-container border-error-container",
    info: "bg-surface-container-highest text-on-surface border-outline-variant",
  }[type] || "bg-surface-container-highest text-on-surface border-outline-variant";

  const icon = {
    success: "check_circle",
    error: "error",
    info: "info",
  }[type] || "info";

  toast.className = `flex items-start gap-2 rounded-xl border px-4 py-3 shadow-lg font-body-sm text-body-sm transition-all duration-300 opacity-0 translate-y-2 ${palette}`;
  toast.innerHTML = `<span class="material-symbols-outlined text-[18px] leading-none shrink-0">${icon}</span><span class="flex-1"></span>`;
  toast.querySelector("span.flex-1").textContent = message;
  container.appendChild(toast);

  requestAnimationFrame(() => {
    toast.classList.remove("opacity-0", "translate-y-2");
  });

  setTimeout(() => {
    toast.classList.add("opacity-0", "translate-y-2");
    setTimeout(() => toast.remove(), 300);
  }, 4200);
}

/** Activa/desactiva el estado de carga de un botón. */
export function setButtonLoading(button, loading, loadingText = "Procesando...") {
  if (!button) return;

  if (loading) {
    button.dataset.originalHtml = button.innerHTML;
    button.disabled = true;
    button.classList.add("opacity-70", "pointer-events-none");
    button.innerHTML = `<span class="material-symbols-outlined animate-spin text-[18px]">progress_activity</span><span></span>`;
    button.querySelector("span:last-child").textContent = loadingText;
  } else {
    button.disabled = false;
    button.classList.remove("opacity-70", "pointer-events-none");
    if (button.dataset.originalHtml) {
      button.innerHTML = button.dataset.originalHtml;
      delete button.dataset.originalHtml;
    }
  }
}

/** Formatea una fecha ISO (YYYY-MM-DD) para mostrarla. */
export function formatDate(value) {
  if (!value) return "—";
  const parts = String(value).split("-");
  if (parts.length !== 3) return value;
  return `${parts[2]}/${parts[1]}/${parts[0]}`;
}
