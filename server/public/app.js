(() => {
  const apiBase =
    (document.querySelector('meta[name="api-base"]')?.content || "").trim() ||
    (typeof window !== "undefined" ? window.FINSET_API_BASE : "") ||
    location.origin;
  const API = `${apiBase.replace(/\/$/, "")}/api/transactions`;
  const STORAGE_KEY = "transactions";
  const demo = [
    { id: 1, type: "income", amount: 8500, category: "Salary", dateISO: "2024-01-05" },
    { id: 2, type: "expense", amount: 1200, category: "Food & Groceries", dateISO: "2024-01-10" },
    { id: 3, type: "expense", amount: 800, category: "Cafe & Restaurants", dateISO: "2024-02-02" },
    { id: 4, type: "income", amount: 8500, category: "Salary", dateISO: "2024-02-05" },
    { id: 5, type: "expense", amount: 2200, category: "Rent", dateISO: "2024-03-01" },
    { id: 6, type: "expense", amount: 600, category: "Entertainment", dateISO: "2024-03-12" },
    { id: 7, type: "income", amount: 8500, category: "Salary", dateISO: "2024-03-05" },
    { id: 8, type: "expense", amount: 300, category: "Transport", dateISO: "2024-04-08" },
    { id: 9, type: "expense", amount: 900, category: "Food & Groceries", dateISO: "2024-04-18" },
    { id: 10, type: "income", amount: 8500, category: "Salary", dateISO: "2024-04-05" },
    { id: 11, type: "expense", amount: 1400, category: "Shopping", dateISO: "2024-05-14" },
    { id: 12, type: "expense", amount: 700, category: "Cafe & Restaurants", dateISO: "2024-06-07" },
    { id: 13, type: "income", amount: 8500, category: "Salary", dateISO: "2024-06-05" },
    { id: 14, type: "expense", amount: 1100, category: "Health & Beauty", dateISO: "2024-07-03" },
  ];

  const els = {
    kpiBalance: document.getElementById("kpi-balance"),
    kpiIncome: document.getElementById("kpi-income"),
    kpiExpense: document.getElementById("kpi-expense"),
    kpiSavings: document.getElementById("kpi-savings"),
    kpiBalanceChange: document.getElementById("kpi-balance-change"),
    kpiIncomeChange: document.getElementById("kpi-income-change"),
    kpiExpenseChange: document.getElementById("kpi-expense-change"),
    kpiSavingsChange: document.getElementById("kpi-savings-change"),
    bars: document.getElementById("cBars"),
    xl: document.getElementById("cXL"),
    txBody: document.getElementById("txBody"),
    donutTotal: document.querySelector(".total-val"),
  };

  let transactions = [];

  init();

  async function init() {
    await loadTransactions();
    updateUI();
    wireUI();
    window.addEventListener("resize", renderChart);
  }

  async function loadTransactions() {
    try {
      const res = await fetch(API);
      if (!res.ok) throw new Error("bad api");
      const rows = (await res.json()) || [];
      transactions = rows.map((r, idx) => ({
        ...r,
        amount: Number(r.amount),
        type: r.type || r.kind,
        category: r.category || r.note || "Uncategorized",
        dateISO: r.dateISO || r.date || r.dateiso,
        id: r.id || idx,
      }));
      localStorage.setItem(STORAGE_KEY, JSON.stringify(transactions));
    } catch {
      try {
        transactions = JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
      } catch {
        transactions = [];
      }
      if (!transactions.length) {
        transactions = demo;
        localStorage.setItem(STORAGE_KEY, JSON.stringify(transactions));
      }
    }
  }

  function updateUI() {
    renderKPIs();
    renderChart();
    renderTransactions();
    renderBudget();
  }

  function renderKPIs() {
    const thisM = monthTotals(0);
    const lastM = monthTotals(1);
    const balance = sumType("income") - sumType("expense");
    const savings = balance;

    setMoney(els.kpiBalance, balance);
    setMoney(els.kpiIncome, thisM.income);
    setMoney(els.kpiExpense, thisM.expense);
    setMoney(els.kpiSavings, savings);

    setChange(els.kpiBalanceChange, balance, lastM.income - lastM.expense, true);
    setChange(els.kpiIncomeChange, thisM.income, lastM.income, true);
    setChange(els.kpiExpenseChange, thisM.expense, lastM.expense, false);
    setChange(els.kpiSavingsChange, savings, lastM.income - lastM.expense, true);
  }

  function renderTransactions() {
    if (!els.txBody) return;
    const rows = [...transactions]
      .sort((a, b) => new Date(b.dateISO) - new Date(a.dateISO))
      .slice(0, 20)
      .map((t) => {
        const cls = t.type === "income" ? "pos" : "neg";
        const sign = t.type === "income" ? "+" : "−";
        return `<tr>
          <td class="td-d">${fmtDate(t.dateISO)}</td>
          <td class="${cls}">${sign} $${fmtMoney(t.amount)}</td>
          <td><span class="txi" style="background:var(--th-bg);">💳</span><span class="txn">${t.category}</span></td>
          <td><span class="mchip"><svg viewBox="0 0 24 24"><rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/></svg>Card ••••</span></td>
          <td><span class="cchip">${t.category}</span></td>
        </tr>`;
      })
      .join("");
    els.txBody.innerHTML = rows || `<tr><td colspan="5" class="td-d">No transactions</td></tr>`;
  }

  function renderChart() {
    if (!els.bars || !els.xl) return;
    els.bars.innerHTML = "";
    els.xl.innerHTML = "";
    const months = lastMonths(7);
    const data = months.map((m) => monthData(m));
    const max = Math.max(1, ...data.flatMap((d) => [d.income, d.expense]));

    data.forEach((d, i) => {
      const col = document.createElement("div");
      col.className = "bc";
      const pair = document.createElement("div");
      pair.className = "bp";

      const inc = document.createElement("div");
      inc.className = "bar inc";
      inc.style.height = Math.max(4, (d.income / max) * 100) + "%";
      inc.innerHTML = `<div class="btt">$${fmtMoney(d.income)}</div>`;

      const exp = document.createElement("div");
      exp.className = "bar exp";
      exp.style.height = Math.max(4, (d.expense / max) * 100) + "%";
      exp.innerHTML = `<div class="btt">$${fmtMoney(d.expense)}</div>`;

      pair.append(inc, exp);
      col.append(pair);
      els.bars.append(col);

      const lbl = document.createElement("div");
      lbl.className = "xl";
      lbl.textContent = months[i].label;
      els.xl.append(lbl);
    });
  }

  function renderBudget() {
    if (!els.donutTotal) return;
    const expense = monthTotals(0).expense;
    els.donutTotal.textContent = `$${fmtMoney(expense)}`;
  }

  // helpers
  function sumType(type) {
    return transactions.filter((t) => t.type === type).reduce((a, b) => a + Number(b.amount || 0), 0);
  }

  function monthTotals(offset) {
    const d = new Date();
    d.setMonth(d.getMonth() - offset);
    return monthData({ month: d.getMonth(), year: d.getFullYear() });
  }

  function monthData({ month, year }) {
    return transactions.reduce(
      (acc, t) => {
        const d = new Date(t.dateISO);
        if (d.getMonth() === month && d.getFullYear() === year) {
          if (t.type === "income") acc.income += Number(t.amount) || 0;
          else acc.expense += Number(t.amount) || 0;
        }
        return acc;
      },
      { income: 0, expense: 0 }
    );
  }

  function lastMonths(n) {
    const arr = [];
    const now = new Date();
    for (let i = n - 1; i >= 0; i--) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
      arr.push({ month: d.getMonth(), year: d.getFullYear(), label: d.toLocaleString("en", { month: "short" }) });
    }
    return arr;
  }

  function setMoney(el, value) {
    if (!el) return;
    const abs = Math.abs(value);
    const [d, c] = abs.toFixed(2).split(".");
    const pref = value < 0 ? "− $" : "$";
    el.innerHTML = `${pref}${Number(d).toLocaleString()}<span class="c">.${c}</span>`;
  }

  function setChange(el, cur, prev, positiveGood) {
    if (!el) return;
    if (!prev) {
      el.textContent = "—";
      el.className = "chg";
      return;
    }
    const diff = ((cur - prev) / prev) * 100;
    const up = diff >= 0;
    el.textContent = `${up ? "↑" : "↓"} ${Math.abs(diff).toFixed(1)}%`;
    el.className = "chg " + (positiveGood ? (up ? "up" : "dn") : (up ? "dn" : "up"));
  }

  function fmtMoney(n) {
    return Number(n || 0).toLocaleString();
  }

  function fmtDate(iso) {
    const d = new Date(iso);
    return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  }

  function wireUI() {
    document.querySelectorAll(".nav-item").forEach((a) =>
      a.addEventListener("click", () => {
        document.querySelectorAll(".nav-item").forEach((x) => x.classList.remove("active"));
        a.classList.add("active");
      })
    );
    document.querySelectorAll(".pill").forEach((p) =>
      p.addEventListener("click", () => {
        document.querySelectorAll(".pill").forEach((x) => x.classList.remove("active"));
        p.classList.add("active");
      })
    );

    const tog = document.getElementById("themeToggle");
    const lbl = document.getElementById("themeLabel");
    const html = document.documentElement;
    const body = document.body;
    const applyTheme = (theme) => {
      html.setAttribute("data-theme", theme);
      html.dataset.theme = theme;
      body.setAttribute("data-theme", theme);
      body.classList.toggle("dark", theme === "dark");
      body.classList.toggle("light", theme !== "dark");
      if (lbl) lbl.textContent = theme === "dark" ? "Dark mode" : "Light mode";
      try {
        localStorage.setItem("finset-theme", theme);
      } catch (e) {}
    };

    const saved = localStorage.getItem("finset-theme") || html.getAttribute("data-theme") || "light";
    applyTheme(saved);

    if (tog) {
      tog.addEventListener("click", () => {
        const next = (html.getAttribute("data-theme") || "light") === "light" ? "dark" : "light";
        applyTheme(next);
      });
    }

    // animate goals
    document.querySelectorAll(".prog-f").forEach((b) => {
      const w = b.dataset.w || b.style.width || "30";
      b.style.width = "0%";
      setTimeout(() => {
        b.style.width = w + "%";
      }, 400);
    });
  }
})();
</script>
