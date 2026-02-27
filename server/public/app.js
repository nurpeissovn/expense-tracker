(() => {
  // ------- STATE & CONSTANTS -------
  const STORAGE_KEY = "et-transactions-v3";
  const THEME_KEY = "et-theme";
  const BUDGET_KEY = "et-budget";
  const API_BASE = `${window.location.origin}/api`;
  const todayISO = new Date().toISOString().split("T")[0];

  const els = {
    incomeTotal: byId("incomeTotal"),
    expenseTotal: byId("expenseTotal"),
    balanceTotal: byId("balanceTotal"),
    transactionsList: byId("transactionsList"),
    form: byId("transactionForm"),
    filterType: byId("filterType"),
    filterCategory: byId("filterCategory"),
    filterStart: byId("filterStart"),
    filterEnd: byId("filterEnd"),
    themeToggle: byId("themeToggle"),
    searchInput: byId("searchInput"),
    typeInput: byId("type"),
    amount: byId("amount"),
    category: byId("category"),
    categoryDot: byId("categoryDot"),
    date: byId("date"),
    note: byId("note"),
    toast: byId("toast"),
    insightsList: byId("insightsList"),
    donutCanvas: byId("donutChart"),
    lineCanvas: byId("lineChart"),
    barCanvas: byId("barChart"),
    donutLegend: byId("donutLegend"),
    budgetInput: byId("budgetInput"),
    budgetBar: byId("budgetBar"),
    budgetMessage: byId("budgetMessage"),
    budgetStatus: byId("budgetStatus"),
    quickAdd: byId("quickAdd"),
    rollupSelect: byId("rollupCategory"),
    rollupTotal: byId("rollupTotal"),
    rollupCount: byId("rollupCount"),
    amountHint: byId("amountHint"),
  };
  const tooltip = createTooltip();

  const colors = [
    "#4f8bff",
    "#2dd4bf",
    "#f59e0b",
    "#c084fc",
    "#fb7185",
    "#34d399",
    "#a78bfa",
    "#22c55e",
    "#38bdf8",
  ];

  let transactions = loadData().map(normalizeTx);
  let budget = loadBudget();
  let donutState = [];
  let linePoints = [];
  let barRects = [];

  // ------- INIT -------
  els.date.value = todayISO;
  wireTheme();
  wireForm();
  wireFilters();
  wireNav();
  wireBudget();
  setupTooltip();
  window.addEventListener("resize", () => renderCharts());
  // initial load: try server first, fallback to cache
  renderAll(true);
  refreshFromServer().finally(() => renderAll());

  // ------- HELPERS -------
  function byId(id) {
    return document.getElementById(id);
  }

  function createTooltip() {
    const el = document.createElement("div");
    el.className = "chart-tooltip";
    document.body.appendChild(el);
    return el;
  }

  function loadData() {
    try {
      return JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
    } catch {
      return [];
    }
  }

  async function refreshFromServer() {
    try {
      const res = await fetch(`${API_BASE}/transactions`);
      if (!res.ok) throw new Error("Failed to fetch transactions");
      const data = await res.json();
      transactions = data.map(normalizeTx);
      saveData();
    } catch (err) {
      console.warn("Using cached data; server unavailable", err);
    }
  }

  function saveData() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(transactions));
  }

  function loadBudget() {
    const val = localStorage.getItem(BUDGET_KEY);
    return val ? Number(val) || 0 : 0;
  }

  function saveBudget(val) {
    budget = Number(val) || 0;
    localStorage.setItem(BUDGET_KEY, budget);
    renderBudget();
  }

  function normalizeTx(tx) {
    return {
      ...tx,
      amount: Number(tx.amount) || 0,
      date: tx.date?.slice(0, 10) || todayISO,
    };
  }

  function kzt(value) {
    return new Intl.NumberFormat("kk-KZ", {
      style: "currency",
      currency: "KZT",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(Number(value) || 0);
  }

  function toast(msg) {
    els.toast.textContent = msg;
    els.toast.classList.add("show");
    setTimeout(() => els.toast.classList.remove("show"), 1800);
  }

  function showTip(text, x, y) {
    tooltip.textContent = text;
    tooltip.style.left = `${x + 12}px`;
    tooltip.style.top = `${y - 10}px`;
    tooltip.style.opacity = 1;
    tooltip.style.transform = "translateY(0)";
  }
  function hideTip() {
    tooltip.style.opacity = 0;
    tooltip.style.transform = "translateY(-4px)";
  }

  function colorForCategory(cat) {
    if (!cat) return colors[0];
    const hash = [...cat].reduce((a, c) => a + c.charCodeAt(0), 0);
    return colors[hash % colors.length];
  }

  // ------- RENDER -------
  function renderAll(withSkeleton = false) {
    if (withSkeleton) toggleSkeleton(true);
    renderFilters();
    renderTotals();
    renderTransactions();
    renderCharts();
    renderInsights();
    renderBudget();
    renderRollup();
    if (withSkeleton) setTimeout(() => toggleSkeleton(false), 300);
  }

  function renderTotals() {
    const income = sum(transactions.filter((t) => t.type === "income").map((t) => t.amount));
    const expense = sum(transactions.filter((t) => t.type === "expense").map((t) => t.amount));
    els.incomeTotal.textContent = kzt(income);
    els.expenseTotal.textContent = kzt(expense);
    els.balanceTotal.textContent = kzt(income - expense);
  }

  function renderTransactions() {
    const filtered = getFilteredTransactions();
    if (!filtered.length) {
      els.transactionsList.classList.add("empty-state");
      els.transactionsList.innerHTML = `
        <div class="empty">
          <div class="empty-icon">🗒️</div>
          <p>No transactions match your filters.</p>
          <button class="ghost-btn" id="emptyAdd">Add Transaction</button>
        </div>`;
      const btn = byId("emptyAdd");
      if (btn) btn.onclick = () => els.amount.focus();
      return;
    }
    els.transactionsList.classList.remove("empty-state");
    els.transactionsList.innerHTML = filtered
      .map((t) => {
        const sign = t.type === "expense" ? "-" : "+";
        return `
        <div class="transaction">
          <div>
            <div class="meta">
              <span class="dot" style="background:${colorForCategory(t.category)}"></span>
              <span>${t.category}</span>
              <span>${t.date}</span>
              ${t.note ? `<span>${t.note}</span>` : ""}
            </div>
            <div class="amount" style="color:${t.type === "expense" ? "var(--danger)" : "var(--success)"}">${sign}${kzt(
          Math.abs(t.amount)
        )}</div>
          </div>
          <button class="delete-btn" data-id="${t.id}">Delete</button>
        </div>`;
      })
      .join("");
    els.transactionsList.querySelectorAll(".delete-btn").forEach((btn) =>
      btn.addEventListener("click", (e) => {
        const id = e.currentTarget.getAttribute("data-id");
        deleteTransaction(id);
      })
    );
  }

  function renderFilters() {
    const categories = Array.from(new Set(transactions.map((t) => t.category))).sort();
    const current = els.filterCategory.value || "all";
    els.filterCategory.innerHTML = `<option value="all">All categories</option>`;
    categories.forEach((cat) => {
      const opt = document.createElement("option");
      opt.value = cat;
      opt.textContent = cat;
      els.filterCategory.appendChild(opt);
    });
    if (categories.includes(current)) els.filterCategory.value = current;

    if (els.rollupSelect) {
      const prev = els.rollupSelect.value || "all";
      els.rollupSelect.innerHTML = `<option value="all">All categories</option>`;
      categories.forEach((cat) => {
        const opt = document.createElement("option");
        opt.value = cat;
        opt.textContent = cat;
        els.rollupSelect.appendChild(opt);
      });
      els.rollupSelect.value = categories.includes(prev) ? prev : "all";
    }
  }

  function renderCharts() {
    resizeCanvas(els.donutCanvas);
    resizeCanvas(els.lineCanvas);
    resizeCanvas(els.barCanvas);
    const data = getFilteredTransactions().filter((t) => t.type === "expense");
    renderDonut(data);
    renderLine();
    renderBar();
  }

  function renderDonut(data) {
    const ctx = els.donutCanvas.getContext("2d");
    clearCanvas(ctx);
    if (!data.length) {
      drawEmptyState(ctx, "No expenses");
      els.donutLegend.innerHTML = "";
      donutState = [];
      return;
    }
    const rect = els.donutCanvas.getBoundingClientRect();
    const totals = groupByCategory(data);
    const sumTotal = sum(Object.values(totals));
    let start = -Math.PI / 2;
    els.donutLegend.innerHTML = "";
    donutState = [];
    const slices = Object.entries(totals).map(([cat, val], idx) => {
      const angle = (val / sumTotal) * Math.PI * 2;
      const slice = { label: cat, value: val, color: colors[idx % colors.length], start, end: start + angle };
      start += angle;
      return slice;
    });
    donutState = slices;
    slices.forEach((s) => {
      const chip = document.createElement("span");
      chip.innerHTML = `<span class="dot" style="background:${s.color}"></span>${s.label}`;
      els.donutLegend.appendChild(chip);
    });
    const cx = rect.width / 2;
    const cy = rect.height / 2;
    const radius = Math.min(cx, cy) - 20;
    const inner = radius * 0.5;
    animate(500, (p) => {
      clearCanvas(ctx);
      slices.forEach((s) => {
        const end = s.start + (s.end - s.start) * p;
        ctx.beginPath();
        ctx.moveTo(cx, cy);
        ctx.arc(cx, cy, radius, s.start, end);
        ctx.fillStyle = s.color;
        ctx.fill();
      });
      ctx.globalCompositeOperation = "destination-out";
      ctx.beginPath();
      ctx.arc(cx, cy, inner, 0, Math.PI * 2);
      ctx.fill();
      ctx.globalCompositeOperation = "source-over";
    });
  }

  function renderLine() {
    const ctx = els.lineCanvas.getContext("2d");
    clearCanvas(ctx);
    const days = lastNDays(30);
    const values = days.map((d) =>
      sum(
        transactions
          .filter((t) => t.type === "expense" && t.date === d)
          .map((t) => t.amount)
      )
    );
    if (sum(values) === 0) return drawEmptyState(ctx, "No data");
    animate(500, (p) => {
      clearCanvas(ctx);
      drawLineChart(ctx, values.map((v) => v * p));
    });
  }

  function renderBar() {
    const ctx = els.barCanvas.getContext("2d");
    clearCanvas(ctx);
    const monthDays = daysInMonth();
    const values = monthDays.map((d) =>
      sum(
        transactions
          .filter((t) => t.type === "expense" && t.date === d)
          .map((t) => t.amount)
      )
    );
    if (sum(values) === 0) return drawEmptyState(ctx, "No data");
    animate(500, (p) => {
      clearCanvas(ctx);
      drawBarChart(ctx, values.map((v) => v * p));
    });
  }

  function renderInsights() {
    const list = [];
    if (!transactions.length) {
      els.insightsList.innerHTML = "<li>Add some transactions to see insights.</li>";
      return;
    }
    const expense = transactions.filter((t) => t.type === "expense");
    const topCat = Object.entries(groupByCategory(expense)).sort((a, b) => b[1] - a[1])[0];
    if (topCat) list.push(`Top spending: ${topCat[0]} (${kzt(topCat[1])}).`);
    const avgDay =
      sum(expense.map((t) => t.amount)) /
      Math.max(
        1,
        new Set(expense.map((t) => t.date)).size
      );
    list.push(`Avg spend per day: ${kzt(avgDay)}.`);
    const thisMonth = monthRangeTotals(0);
    const lastMonth = monthRangeTotals(1);
    if (lastMonth > 0) {
      const diff = ((thisMonth - lastMonth) / lastMonth) * 100;
      list.push(`You spent ${diff >= 0 ? "↑" : "↓"}${Math.abs(diff).toFixed(1)}% vs last month.`);
    }
    const maxTx = expense.sort((a, b) => b.amount - a.amount)[0];
    if (maxTx) list.push(`Largest transaction: ${kzt(maxTx.amount)} on ${maxTx.category}.`);
    els.insightsList.innerHTML = list.map((i) => `<li>${i}</li>`).join("");
  }

  function renderBudget() {
    els.budgetInput.value = budget || "";
    const spent = sum(transactions.filter((t) => t.type === "expense").map((t) => t.amount));
    const pct = budget > 0 ? Math.min((spent / budget) * 100, 150) : 0;
    els.budgetBar.style.width = `${pct}%`;
    els.budgetBar.style.background = spent > budget && budget > 0 ? "var(--danger)" : "linear-gradient(135deg, var(--accent), var(--accent-2))";
    if (!budget) {
      els.budgetStatus.textContent = "Set";
      els.budgetMessage.textContent = "No budget set.";
    } else if (spent <= budget) {
      els.budgetStatus.textContent = "On track";
      els.budgetMessage.textContent = `${kzt(budget - spent)} remaining`;
    } else {
      els.budgetStatus.textContent = "Over";
      els.budgetMessage.textContent = `Over by ${kzt(spent - budget)}`;
    }
  }

  function renderRollup() {
    if (!els.rollupSelect) return;
    const cat = els.rollupSelect.value;
    const filtered = cat === "all" ? transactions : transactions.filter((t) => t.category === cat);
    const total = sum(filtered.map((t) => t.amount));
    els.rollupTotal.textContent = kzt(total);
    els.rollupCount.textContent = filtered.length;
  }

  async function deleteTransaction(id) {
    try {
      await fetch(`${API_BASE}/transactions/${id}`, { method: "DELETE" });
    } catch (err) {
      console.warn("Server delete failed, removing locally", err);
    } finally {
      transactions = transactions.filter((t) => t.id !== id);
      saveData();
      toast("Deleted");
      renderAll();
    }
  }

  function toggleSkeleton(state) {
    document.querySelectorAll(".canvas-wrap").forEach((el) => {
      el.classList.toggle("loading", state);
    });
  }

  // ------- EVENT WIRING -------
  function wireTheme() {
    const saved = localStorage.getItem(THEME_KEY) || "light";
    document.documentElement.setAttribute("data-theme", saved);
    els.themeToggle.checked = saved === "dark";
    els.themeToggle.onchange = () => {
      const next = els.themeToggle.checked ? "dark" : "light";
      document.documentElement.setAttribute("data-theme", next);
      localStorage.setItem(THEME_KEY, next);
    };
  }

  function wireForm() {
    els.amount.addEventListener("input", () => {
      const val = Number(els.amount.value);
      els.amountHint.textContent = val > 0 ? `Will save ${kzt(val)}` : "";
    });

    document.querySelectorAll(".type-btn").forEach((btn) =>
      btn.addEventListener("click", () => {
        document.querySelectorAll(".type-btn").forEach((b) => b.classList.remove("active"));
        btn.classList.add("active");
        els.typeInput.value = btn.dataset.type;
      })
    );

    els.category.addEventListener("input", () => {
      els.categoryDot.style.background = colorForCategory(els.category.value);
    });

    els.form.addEventListener("submit", async (e) => {
      e.preventDefault();
      const amount = Number(els.amount.value);
      if (!amount || amount <= 0) return alert("Amount must be greater than zero.");
      const tx = {
        id: crypto.randomUUID ? crypto.randomUUID() : Date.now().toString(),
        type: els.typeInput.value,
        amount,
        category: els.category.value.trim() || "Uncategorized",
        date: els.date.value || todayISO,
        note: els.note.value.trim(),
      };
      try {
        const res = await fetch(`${API_BASE}/transactions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(tx),
        });
        if (!res.ok) throw new Error(await res.text());
        const saved = normalizeTx(await res.json());
        transactions.unshift(saved);
      } catch (err) {
        console.warn("Falling back to local add", err);
        transactions.unshift(tx);
      }
      saveData();
      toast("Transaction added");
      els.form.reset();
      els.typeInput.value = "income";
      document.querySelectorAll(".type-btn").forEach((b, i) => b.classList.toggle("active", i === 0));
      els.date.value = todayISO;
      renderAll();
    });
  }

  function wireFilters() {
    [els.filterType, els.filterCategory, els.filterStart, els.filterEnd].forEach((el) => {
      el.addEventListener("change", () => renderAll());
    });
    els.searchInput.addEventListener("input", () => renderTransactions());
    if (els.quickAdd) els.quickAdd.onclick = () => els.amount.focus();
  }

  function wireNav() {
    const navBtns = document.querySelectorAll(".nav-btn");
    navBtns.forEach((btn) =>
      btn.addEventListener("click", () => {
        const targetId = btn.dataset.target;
        navBtns.forEach((b) => b.classList.remove("active"));
        btn.classList.add("active");
        if (targetId) {
          const el = document.getElementById(targetId);
          if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      })
    );
  }

  function wireBudget() {
    els.budgetInput.addEventListener("change", (e) => {
      saveBudget(e.target.value);
      toast("Budget updated");
    });
  }

  function setupTooltip() {
    els.donutCanvas.addEventListener("mousemove", handleDonutHover);
    els.donutCanvas.addEventListener("mouseleave", hideTip);
    els.lineCanvas.addEventListener("mousemove", handleLineHover);
    els.lineCanvas.addEventListener("mouseleave", hideTip);
    els.barCanvas.addEventListener("mousemove", handleBarHover);
    els.barCanvas.addEventListener("mouseleave", hideTip);
  }

  // ------- FILTERING -------
  function getFilteredTransactions() {
    const type = els.filterType.value;
    const cat = els.filterCategory.value;
    const start = els.filterStart.value;
    const end = els.filterEnd.value;
    const term = els.searchInput.value.trim().toLowerCase();
    return transactions.filter((t) => {
      const typeOk = type === "all" || t.type === type;
      const catOk = cat === "all" || t.category === cat;
      const startOk = !start || t.date >= start;
      const endOk = !end || t.date <= end;
      const termOk = !term || `${t.category} ${t.note}`.toLowerCase().includes(term);
      return typeOk && catOk && startOk && endOk && termOk;
    });
  }

  // ------- CHART UTILITIES -------
  function clearCanvas(ctx) {
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
  }

  function resizeCanvas(canvas) {
    const dpr = window.devicePixelRatio || 1;
    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;
    const ctx = canvas.getContext("2d");
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.scale(dpr, dpr);
  }

  function drawEmptyState(ctx, text) {
    ctx.fillStyle = "var(--muted)";
    ctx.font = "14px Inter";
    ctx.textAlign = "center";
    ctx.fillText(text, ctx.canvas.width / 2, ctx.canvas.height / 2);
  }

  function animate(duration, cb) {
    const start = performance.now();
    function frame(now) {
      const p = Math.min((now - start) / duration, 1);
      cb(p);
      if (p < 1) requestAnimationFrame(frame);
    }
    requestAnimationFrame(frame);
  }

  function drawLineChart(ctx, values) {
    const rect = ctx.canvas.getBoundingClientRect();
    const w = rect.width;
    const h = rect.height;
    const max = Math.max(...values, 1);
    const pad = 20;
    const step = values.length > 1 ? (w - pad * 2) / (values.length - 1) : 0;
    linePoints = [];
    ctx.strokeStyle = "var(--accent)";
    ctx.lineWidth = 3;
    ctx.beginPath();
    values.forEach((v, i) => {
      const x = pad + i * step;
      const y = h - pad - (v / max) * (h - pad * 2);
      if (i === 0) ctx.moveTo(x, y);
      else ctx.lineTo(x, y);
      linePoints.push({ x, y, value: v });
    });
    ctx.stroke();
  }

  function drawBarChart(ctx, values) {
    const rect = ctx.canvas.getBoundingClientRect();
    const w = rect.width;
    const h = rect.height;
    const max = Math.max(...values, 1);
    const pad = 20;
    const bw = (w - pad * 2) / values.length;
    barRects = [];
    values.forEach((v, i) => {
      const x = pad + i * bw;
      const height = (v / max) * (h - pad * 2);
      const y = h - pad - height;
      ctx.fillStyle = "var(--accent)";
      ctx.fillRect(x, y, bw * 0.8, height);
      barRects.push({ x, y, w: bw * 0.8, h: height, value: v });
    });
  }

  function handleDonutHover(e) {
    if (!donutState.length) return;
    const rect = els.donutCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const cx = rect.width / 2;
    const cy = rect.height / 2;
    const radius = Math.min(cx, cy) - 20;
    const dx = x - cx;
    const dy = y - cy;
    const r = Math.sqrt(dx * dx + dy * dy);
    if (r < radius * 0.5 || r > radius + 5) return hideTip();
    const angle = Math.atan2(dy, dx);
    const norm = angle < -Math.PI / 2 ? angle + Math.PI * 2 : angle;
    const slice = donutState.find((s) => norm >= s.start && norm < s.end);
    if (!slice) return hideTip();
    showTip(`${slice.label}: ${kzt(slice.value)}`, e.clientX, e.clientY);
  }

  function handleLineHover(e) {
    if (!linePoints.length) return;
    const rect = els.lineCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    let nearest = linePoints[0];
    let min = Math.abs(x - nearest.x);
    linePoints.forEach((p) => {
      const d = Math.abs(x - p.x);
      if (d < min) {
        min = d;
        nearest = p;
      }
    });
    showTip(kzt(nearest.value), e.clientX, e.clientY);
  }

  function handleBarHover(e) {
    if (!barRects.length) return;
    const rect = els.barCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const hit = barRects.find((r) => x >= r.x && x <= r.x + r.w && y >= r.y && y <= r.y + r.h);
    if (!hit) return hideTip();
    showTip(kzt(hit.value), e.clientX, e.clientY);
  }

  // ------- DATA HELPERS -------
  function sum(arr) {
    return arr.reduce((a, b) => a + (Number(b) || 0), 0);
  }

  function groupByCategory(items) {
    return items.reduce((acc, t) => {
      acc[t.category] = (acc[t.category] || 0) + (Number(t.amount) || 0);
      return acc;
    }, {});
  }

  function lastNDays(n) {
    const out = [];
    const d = new Date();
    for (let i = n - 1; i >= 0; i--) {
      const tmp = new Date(d);
      tmp.setDate(d.getDate() - i);
      out.push(tmp.toISOString().split("T")[0]);
    }
    return out;
  }

  function daysInMonth(offset = 0) {
    const d = new Date();
    d.setMonth(d.getMonth() - offset, 1);
    const month = d.getMonth();
    const year = d.getFullYear();
    const res = [];
    const last = new Date(year, month + 1, 0).getDate();
    for (let i = 1; i <= last; i++) {
      res.push(`${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`);
    }
    return res;
  }

  function monthRangeTotals(offset) {
    const days = daysInMonth(offset);
    return sum(
      transactions
        .filter((t) => t.type === "expense" && days.includes(t.date))
        .map((t) => t.amount)
    );
  }
})();
