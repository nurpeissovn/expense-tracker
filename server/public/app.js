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
  const pulseTotalEl = document.getElementById("pulseTotal");
  const pulseChangeEl = document.getElementById("pulseChange");
  const pulseAvgEl = document.getElementById("pulseAvg");
  const pulseTopCatEl = document.getElementById("pulseTopCat");
  const pulseSparkEl = document.getElementById("pulseSpark");
  const pulseRangeLabelEl = document.getElementById("pulseRangeLabel");

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
  renderPulse();
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
      renderPulse();
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
    renderPulse();
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
      let msg = await res.text();
      try {
        const parsed = JSON.parse(msg);
        msg = parsed.error || parsed.message || msg;
      } catch (_e) {
        // keep text
      }
      throw new Error(msg || "Failed to create");
    }
    return res.json();
  }

  async function deleteRemote(id) {
    const res = await fetch(`${API_BASE}/transactions/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error("Delete failed");
  }

  const currencyFormatter = (() => {
    try {
      return new Intl.NumberFormat("kk-KZ", {
        style: "currency",
        currency: "KZT",
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    } catch (err) {
      console.warn("Falling back currency formatter", err);
      return new Intl.NumberFormat(undefined, {
        style: "currency",
        currency: "KZT",
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    }
  })();

  function formatCurrency(val) {
    return currencyFormatter.format(Number(val) || 0);
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

  function renderPulse() {
    if (!pulseTotalEl) return;
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const toDate = (str) => {
      const d = new Date(str + "T00:00:00");
      d.setHours(0, 0, 0, 0);
      return d;
    };

    const expenses = transactions.filter((t) => t.type === "expense").map((t) => ({ ...t, _d: toDate(t.date) }));

    const rangeTotals = (daysBack, length) => {
      const end = new Date(today);
      end.setDate(end.getDate() - daysBack);
      const start = new Date(end);
      start.setDate(start.getDate() - (length - 1));
      const total = expenses
        .filter((t) => t._d >= start && t._d <= end)
        .reduce((s, t) => s + t.amount, 0);
      return total;
    };

    const last30Total = rangeTotals(0, 30);
    const prev30Total = rangeTotals(30, 30);
    const changeAbs = last30Total - prev30Total;
    const changePct = prev30Total === 0 ? (last30Total > 0 ? 100 : 0) : (changeAbs / prev30Total) * 100;

    const avgDaily = last30Total / 30 || 0;

    const last30Start = new Date();
    last30Start.setHours(0, 0, 0, 0);
    last30Start.setDate(last30Start.getDate() - 29);

    const topCategoryMap = expenses
      .filter((t) => t._d >= last30Start)
      .reduce((acc, t) => {
        acc[t.category] = (acc[t.category] || 0) + t.amount;
        return acc;
      }, {});
    const topCat = Object.entries(topCategoryMap).sort((a, b) => b[1] - a[1])[0];

    pulseTotalEl.textContent = formatCurrency(last30Total);
    const changeText = `${changeAbs >= 0 ? "+" : "−"}${Math.abs(changePct).toFixed(1)}%`;
    pulseChangeEl.textContent = `${changeText} vs prev 30d`;
    pulseChangeEl.style.color = changeAbs >= 0 ? "var(--success)" : "var(--danger)";
    pulseAvgEl.textContent = formatCurrency(avgDaily);
    pulseTopCatEl.textContent = topCat ? `Top category: ${topCat[0]} (${formatCurrency(topCat[1])})` : "Top category: —";

    // sparkline for last 14 days
    const days = 14;
    const dailyTotals = Array.from({ length: days }, (_, i) => {
      const d = new Date();
      d.setHours(0, 0, 0, 0);
      d.setDate(d.getDate() - (days - 1 - i));
      const total = expenses
        .filter((t) => t._d.getTime() === d.getTime())
        .reduce((s, t) => s + t.amount, 0);
      return total;
    });

    pulseRangeLabelEl.textContent = `Last ${days} days`;
    drawSparkline(dailyTotals);
  }

  function drawSparkline(values) {
    if (!pulseSparkEl) return;
    const width = 240;
    const height = 80;
    pulseSparkEl.innerHTML = "";
    const max = Math.max(...values, 1);
    const step = values.length > 1 ? width / (values.length - 1) : width;

    const points = values.map((v, i) => {
      const x = i * step;
      const y = height - (v / max) * (height - 6) - 3;
      return [x, y];
    });

    if (!points.length) return;

    const pathD = points.map((p, i) => (i === 0 ? `M ${p[0]} ${p[1]}` : `L ${p[0]} ${p[1]}`)).join(" ");
    const areaPath = `${pathD} L ${points[points.length - 1][0]} ${height} L 0 ${height} Z`;

    const area = document.createElementNS("http://www.w3.org/2000/svg", "path");
    area.setAttribute("d", areaPath);
    area.setAttribute("class", "sparkline-area");

    const line = document.createElementNS("http://www.w3.org/2000/svg", "path");
    line.setAttribute("d", pathD);
    line.setAttribute("class", "sparkline");

    pulseSparkEl.appendChild(area);
    pulseSparkEl.appendChild(line);
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
        const formattedAmount = formatCurrency(Math.abs(t.amount));
        const amountColor = t.type === "expense" ? "var(--danger)" : "var(--success)";
        return `
          <div class="transaction" data-id="${t.id}">
            <div>
              <div class="meta">
                <span class="badge">${t.category}</span>
                <span>${t.date}</span>
                ${t.note ? `<span>${t.note}</span>` : ""}
              </div>
              <div class="amount" style="color:${amountColor}">${sign}${formattedAmount}</div>
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
      renderPulse();
      return;
    }

    deleteRemote(id)
      .then(() => {
        transactions = transactions.filter((t) => t.id !== id);
        saveCache();
        renderFilters();
        renderTransactions();
        renderTotals();
        renderPulse();
      })
      .catch((err) => {
        console.error(err);
        alert("Could not delete. Check your connection and try again.");
      });
  }
})();
