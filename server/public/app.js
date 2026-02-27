(() => {
  const storageKey = "et-transactions-v2";
  const themeKey = "et-theme";
  const API_BASE = `${window.location.origin}/api`;
  const offlineMode = window.location.protocol === "file:";

  const incomeTotalEl = document.getElementById("incomeTotal");
  const expenseTotalEl = document.getElementById("expenseTotal");
  const balanceTotalEl = document.getElementById("balanceTotal");
  const transactionsListEl = document.getElementById("transactionsList");
  const form = document.getElementById("transactionForm");
  const filterTypeEl = document.getElementById("filterType");
  const filterCategoryEl = document.getElementById("filterCategory");
  const themeToggleBtn = document.getElementById("themeToggle");
  const dateInput = document.getElementById("date");

  let transactions = [];

  const today = new Date().toISOString().split("T")[0];
  dateInput.value = today;

  const savedTheme = localStorage.getItem(themeKey) || "light";
  setTheme(savedTheme);

  // seed from cache while network loads
  loadCache();
  renderFilters();
  renderTransactions();
  renderTotals();
  if (!offlineMode) {
    refreshFromServer();
  }

  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const type = document.getElementById("type").value;
    const amount = parseFloat(document.getElementById("amount").value);
    const category = document.getElementById("category").value.trim();
    const date = document.getElementById("date").value;
    const note = document.getElementById("note").value.trim();

    if (!amount || amount <= 0 || !category || !date) {
      alert("Please enter a valid amount, category, and date.");
      return;
    }

    const newTx = {
      type,
      amount: Math.round(amount * 100) / 100,
      category,
      date,
      note,
    };

    if (offlineMode) {
      addLocal(newTx);
      return;
    }

    createTransaction(newTx)
      .then((saved) => addLocal(saved))
      .catch((err) => {
        console.error(err);
        alert("Could not save transaction. Check your connection and try again.");
      });
  });

  filterTypeEl.addEventListener("change", () => {
    renderTransactions();
  });

  filterCategoryEl.addEventListener("change", () => {
    renderTransactions();
  });

  themeToggleBtn.addEventListener("click", () => {
    const next = document.documentElement.getAttribute("data-theme") === "light" ? "dark" : "light";
    setTheme(next);
  });

  function setTheme(mode) {
    document.documentElement.setAttribute("data-theme", mode);
    localStorage.setItem(themeKey, mode);
    themeToggleBtn.textContent = mode === "light" ? "🌙 Dark" : "☀️ Light";
  }

  function loadCache() {
    try {
      const stored = localStorage.getItem(storageKey);
      transactions = stored ? JSON.parse(stored) : [];
    } catch (err) {
      console.error("Failed to load cache", err);
      transactions = [];
    }
  }

  function saveCache() {
    localStorage.setItem(storageKey, JSON.stringify(transactions));
  }

  async function refreshFromServer() {
    try {
      const res = await fetch(`${API_BASE}/transactions`);
      if (!res.ok) throw new Error("Network response was not ok");
      const data = await res.json();
      transactions = data;
      saveCache();
      renderFilters();
      renderTransactions();
      renderTotals();
    } catch (err) {
      console.warn("Using cached data; failed to fetch from server", err);
    }
  }

  function addLocal(tx) {
    const withId = tx.id
      ? tx
      : {
          ...tx,
          id: crypto.randomUUID ? crypto.randomUUID() : Date.now().toString(),
        };
    transactions.unshift(withId);
    saveCache();
    renderFilters();
    renderTransactions();
    renderTotals();
    form.reset();
    document.getElementById("type").value = tx.type;
    dateInput.value = today;
  }

  async function createTransaction(payload) {
    const res = await fetch(`${API_BASE}/transactions`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const msg = await res.text();
      throw new Error(msg || "Failed to create");
    }
    return res.json();
  }

  async function deleteRemote(id) {
    const res = await fetch(`${API_BASE}/transactions/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error("Delete failed");
  }

  function formatCurrency(val) {
    return `$${Number(val).toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })}`;
  }

  function renderTotals() {
    const income = transactions
      .filter((t) => t.type === "income")
      .reduce((sum, t) => sum + t.amount, 0);
    const expense = transactions
      .filter((t) => t.type === "expense")
      .reduce((sum, t) => sum + t.amount, 0);
    const balance = income - expense;

    incomeTotalEl.textContent = formatCurrency(income);
    expenseTotalEl.textContent = formatCurrency(expense);
    balanceTotalEl.textContent = formatCurrency(balance);
  }

  function renderFilters() {
    const categories = Array.from(new Set(transactions.map((t) => t.category))).sort();
    const current = filterCategoryEl.value || "all";
    filterCategoryEl.innerHTML = `<option value="all">All categories</option>`;
    categories.forEach((cat) => {
      const opt = document.createElement("option");
      opt.value = cat;
      opt.textContent = cat;
      filterCategoryEl.appendChild(opt);
    });
    if (categories.includes(current)) {
      filterCategoryEl.value = current;
    }
  }

  function renderTransactions() {
    if (!transactions.length) {
      transactionsListEl.classList.add("empty-state");
      transactionsListEl.innerHTML = "<p>No transactions yet. Add one to get started.</p>";
      return;
    }

    const typeFilter = filterTypeEl.value;
    const categoryFilter = filterCategoryEl.value;

    const filtered = transactions.filter((t) => {
      const typeMatch = typeFilter === "all" ? true : t.type === typeFilter;
      const catMatch = categoryFilter === "all" ? true : t.category === categoryFilter;
      return typeMatch && catMatch;
    });

    if (!filtered.length) {
      transactionsListEl.classList.add("empty-state");
      transactionsListEl.innerHTML = "<p>No transactions match the filters.</p>";
      return;
    }

    transactionsListEl.classList.remove("empty-state");
    transactionsListEl.innerHTML = filtered
      .map((t) => {
        const sign = t.type === "expense" ? "-" : "+";
        const amountColor = t.type === "expense" ? "var(--danger)" : "var(--success)";
        return `
          <div class="transaction" data-id="${t.id}">
            <div>
              <div class="meta">
                <span class="badge">${t.category}</span>
                <span>${t.date}</span>
                ${t.note ? `<span>${t.note}</span>` : ""}
              </div>
              <div class="amount" style="color:${amountColor}">${sign}${formatCurrency(t.amount).slice(1)}</div>
            </div>
            <button class="delete-btn" aria-label="Delete" data-id="${t.id}">Delete</button>
          </div>
        `;
      })
      .join("");

    transactionsListEl.querySelectorAll(".delete-btn").forEach((btn) => {
      btn.addEventListener("click", (e) => {
        const id = e.currentTarget.getAttribute("data-id");
        deleteTransaction(id);
      });
    });
  }

  function deleteTransaction(id) {
    if (offlineMode) {
      transactions = transactions.filter((t) => t.id !== id);
      saveCache();
      renderFilters();
      renderTransactions();
      renderTotals();
      return;
    }

    deleteRemote(id)
      .then(() => {
        transactions = transactions.filter((t) => t.id !== id);
        saveCache();
        renderFilters();
        renderTransactions();
        renderTotals();
      })
      .catch((err) => {
        console.error(err);
        alert("Could not delete. Check your connection and try again.");
      });
  }
})();
