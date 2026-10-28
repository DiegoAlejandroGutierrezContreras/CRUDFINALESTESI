// Detalle de factura, modal JSON e impresión.

import { API } from "./config.js";
import { apiGet } from "./api.js";
import { initLayout } from "./layout.js";
import { escapeHTML, formatCurrency, formatDate, showToast } from "./ui.js";

let facturaActual = null;

function renderFactura(factura) {
  facturaActual = factura;

  document.title = `Factura #${factura.factura_id} — Factura`;
  document.getElementById("invoice-number").textContent = `Factura #${factura.factura_id}`;
  document.getElementById("invoice-date").textContent = formatDate(factura.fecha);
  document.getElementById("invoice-client-name").textContent = factura.cliente?.nombre || "—";
  document.getElementById("invoice-client-email").textContent = factura.cliente?.email || "—";
  document.getElementById("invoice-client-id").textContent = `ID: ${factura.cliente?.id ?? "—"}`;
  document.getElementById("invoice-user-name").textContent = factura.usuario?.nombre || "—";
  document.getElementById("invoice-total").textContent = formatCurrency(factura.total);
  document.getElementById("json-title").textContent = `Factura #${factura.factura_id} JSON`;
  document.getElementById("json-code").textContent = JSON.stringify(factura, null, 2);

  const tbody = document.getElementById("invoice-items");
  tbody.innerHTML = "";

  (factura.productos || []).forEach((item) => {
    const row = document.createElement("tr");
    row.className = "hover:bg-surface-container-low/40 transition-colors";
    row.innerHTML = `
      <td class="py-space-md px-space-md font-label-md text-label-md text-on-surface item-name"></td>
      <td class="py-space-md px-space-md text-center tabular-nums">${item.cantidad}</td>
      <td class="py-space-md px-space-md text-right font-tabular-amount-md text-tabular-amount-md text-on-surface-variant">${formatCurrency(item.precio_unitario)}</td>
      <td class="py-space-md px-space-md text-right font-tabular-amount-md text-tabular-amount-md text-on-surface font-semibold">${formatCurrency(item.subtotal)}</td>
    `;
    row.querySelector(".item-name").textContent = item.nombre;
    tbody.appendChild(row);
  });
}

function showError(message) {
  document.getElementById("invoice-sheet")?.classList.add("hidden");
  const errorBox = document.getElementById("invoice-error");
  if (errorBox) {
    errorBox.classList.remove("hidden");
    document.getElementById("invoice-error-text").textContent = message;
  }
}

function initJsonModal() {
  const modal = document.getElementById("json-modal");
  const openButton = document.getElementById("btn-view-json");
  const closeButton = document.getElementById("close-json-modal");
  const closeFooterButton = document.getElementById("btn-close-footer");

  openButton?.addEventListener("click", () => {
    if (facturaActual) modal?.showModal();
  });
  closeButton?.addEventListener("click", () => modal?.close());
  closeFooterButton?.addEventListener("click", () => modal?.close());

  modal?.addEventListener("click", (event) => {
    const rect = modal.getBoundingClientRect();
    const inside =
      event.clientY >= rect.top &&
      event.clientY <= rect.top + rect.height &&
      event.clientX >= rect.left &&
      event.clientX <= rect.left + rect.width;
    if (!inside) modal.close();
  });
}

async function loadFactura() {
  const id = Number(new URLSearchParams(window.location.search).get("id"));
  if (!id || !Number.isFinite(id) || id <= 0) {
    showError("No se indicó una factura válida.");
    return;
  }

  try {
    const factura = await apiGet(API.facturacion, `/facturas/${id}`);
    renderFactura(factura);
  } catch (error) {
    showError(error.message);
    showToast(error.message, "error");
  }
}

initLayout("facturas");
initJsonModal();
document.getElementById("btn-print")?.addEventListener("click", () => window.print());
loadFactura();
