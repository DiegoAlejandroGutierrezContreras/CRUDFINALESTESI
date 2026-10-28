// Catálogo de productos (solo lectura).

import { API } from "./config.js";
import { apiGet } from "./api.js";
import { initLayout } from "./layout.js";
import { escapeHTML, formatCurrency, showToast } from "./ui.js";

function renderRows(products) {
  const tbody = document.getElementById("catalog-rows");
  const empty = document.getElementById("catalog-empty");
  const table = document.getElementById("catalog-table");
  if (!tbody) return;

  tbody.innerHTML = "";

  if (!products.length) {
    table?.classList.add("hidden");
    empty?.classList.remove("hidden");
    return;
  }

  table?.classList.remove("hidden");
  empty?.classList.add("hidden");

  products.forEach((product, index) => {
    const row = document.createElement("tr");
    row.className = `h-14 ${
      index % 2 === 0 ? "bg-surface-container-lowest" : ""
    } hover:bg-surface-container-low/60 transition-colors`;
    row.innerHTML = `
      <td class="pl-space-lg pr-space-md py-space-sm font-tabular-amount-md text-tabular-amount-md text-secondary">${product.id}</td>
      <td class="px-space-md py-space-sm font-body-md text-body-md text-on-surface"></td>
      <td class="px-space-md py-space-sm font-tabular-amount-md text-tabular-amount-md text-on-surface text-right">${formatCurrency(product.precio)}</td>
      <td class="pl-space-md pr-space-lg py-space-sm font-tabular-amount-md text-tabular-amount-md text-on-surface text-right">
        <span class="${product.stock > 0 ? "text-on-surface" : "text-error"}">${product.stock}</span>
      </td>
    `;
    row.querySelector("td:nth-child(2)").textContent = product.nombre;
    tbody.appendChild(row);
  });
}

async function loadProducts() {
  const tbody = document.getElementById("catalog-rows");
  if (tbody) {
    tbody.innerHTML = `<tr><td class="py-12 text-center text-on-surface-variant" colspan="4">Cargando productos...</td></tr>`;
  }

  try {
    const products = await apiGet(API.productos, "/productos");
    renderRows(Array.isArray(products) ? products : []);
  } catch (error) {
    showToast(error.message, "error");
    if (tbody) {
      tbody.innerHTML = `<tr><td class="py-12 text-center text-error" colspan="4"></td></tr>`;
      tbody.querySelector("td").textContent = error.message;
    }
  }
}

initLayout("productos");
loadProducts();
