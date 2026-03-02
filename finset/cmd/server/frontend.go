package main

const indexHTML = `<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>FinSet – Finance Dashboard</title>
<link href="https://fonts.googleapis.com/css2?family=Sora:wght@300;400;500;600;700&family=DM+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
/* ═══════════════════════════════════════════
   THEME TOKENS
═══════════════════════════════════════════ */
[data-theme="light"] {
  --bg:           #ede9ff;
  --sidebar-bg:   #ffffff;
  --card-bg:      #ffffff;
  --accent:       #7c5cfc;
  --accent-soft:  #ede9ff;
  --accent-hover: #6a4aed;
  --text-1:       #1a1340;
  --text-2:       #9691b0;
  --text-3:       #c4bedd;
  --border:       #ede9ff;
  --row-hover:    #f9f8ff;
  --th-bg:        #f5f3ff;
  --chip-bg:      #f0effe;
  --chip-border:  #e3deff;
  --green-bg:     #dcfce7;
  --green-text:   #15803d;
  --red-bg:       #fee2e2;
  --red-text:     #dc2626;
  --shadow:       0 2px 16px rgba(124,92,252,.08);
  --shadow-card:  0 1px 8px rgba(124,92,252,.06);
  --bar-expense:  #c4b5fd;
  --bar-income:   #7c5cfc;
  --grid-line:    #ede9ff;
  --scroll-thumb: #ddd6fe;
  --modal-overlay:#0f0c1e99;
  --input-bg:     #f9f8ff;
}
[data-theme="dark"] {
  --bg:           #0f0c1e;
  --sidebar-bg:   #17122b;
  --card-bg:      #1e1836;
  --accent:       #7c5cfc;
  --accent-soft:  #2d2250;
  --accent-hover: #9775fa;
  --text-1:       #f0ecff;
  --text-2:       #7b6fa0;
  --text-3:       #453b6b;
  --border:       #2a2244;
  --row-hover:    #221d3a;
  --th-bg:        #1a1530;
  --chip-bg:      #221d3a;
  --chip-border:  #332d52;
  --green-bg:     #052e16;
  --green-text:   #4ade80;
  --red-bg:       #2c0e0e;
  --red-text:     #f87171;
  --shadow:       0 2px 20px rgba(0,0,0,.45);
  --shadow-card:  0 1px 12px rgba(0,0,0,.3);
  --bar-expense:  #4a3880;
  --bar-income:   #7c5cfc;
  --grid-line:    #2a2244;
  --scroll-thumb: #2d2250;
  --modal-overlay:#00000088;
  --input-bg:     #17122b;
}

/* ═══════════════════════════════════════════
   RESET + BASE
═══════════════════════════════════════════ */
*, *::before, *::after { margin:0; padding:0; box-sizing:border-box; }
::-webkit-scrollbar { width:4px; height:4px; }
::-webkit-scrollbar-track { background:transparent; }
::-webkit-scrollbar-thumb { background:var(--scroll-thumb); border-radius:99px; }
html { font-size:14px; }
body {
  font-family:'Sora',sans-serif;
  background:var(--bg);
  color:var(--text-1);
  min-height:100vh;
  display:flex;
  overflow:hidden;
  transition:background .3s, color .3s;
}

/* ═══════════════════════════════════════════
   SIDEBAR
═══════════════════════════════════════════ */
.sidebar {
  width:210px; min-width:210px;
  background:var(--sidebar-bg);
  border-right:1px solid var(--border);
  display:flex; flex-direction:column;
  padding:22px 0 16px;
  height:100vh;
  box-shadow:var(--shadow);
  position:relative; z-index:10;
  transition:width .25s cubic-bezier(.4,0,.2,1), min-width .25s cubic-bezier(.4,0,.2,1), background .3s, border-color .3s;
  overflow:hidden;
}
/* COLLAPSED STATE */
.sidebar.collapsed { width:60px; min-width:60px; }
.sidebar.collapsed .logo span,
.sidebar.collapsed .nav-label,
.sidebar.collapsed .nav-badge-wrap,
.sidebar.collapsed .tog-label,
.sidebar.collapsed .tog { display:none; }
.sidebar.collapsed .logo { padding:0 12px 26px; justify-content:center; }
.sidebar.collapsed .nav-item { justify-content:center; padding:9px 0; }
.sidebar.collapsed .nav-item:hover .nav-tooltip { display:block; }
.sidebar.collapsed .nav-bottom .nav-item { justify-content:center; padding:9px 0; }
.sidebar.collapsed .toggle-row { justify-content:center; }
.sidebar.collapsed .divider { margin:8px 8px; }

.logo {
  display:flex; align-items:center; gap:10px;
  padding:0 18px 26px;
  transition:padding .25s, justify-content .25s;
}
.logo-icon {
  width:34px; height:34px; border-radius:10px;
  background:var(--accent);
  display:flex; align-items:center; justify-content:center; flex-shrink:0;
}
.logo span { font-weight:700; font-size:17px; white-space:nowrap; transition:opacity .2s; }

.nav { flex:1; padding:0 10px; display:flex; flex-direction:column; gap:2px; overflow:hidden; }
.nav-item {
  display:flex; align-items:center; gap:10px;
  padding:9px 12px; border-radius:10px;
  color:var(--text-2); font-size:13.5px; font-weight:500;
  cursor:pointer; transition:all .2s; text-decoration:none; user-select:none;
  position:relative; white-space:nowrap;
}
.nav-item svg { width:16px; height:16px; flex-shrink:0; stroke:currentColor; fill:none; stroke-width:1.8; stroke-linecap:round; stroke-linejoin:round; }
.nav-item:hover { background:var(--accent-soft); color:var(--accent); }
.nav-item.active { background:var(--accent); color:#fff; }
.nav-label { flex:1; }
/* Collapse tooltip */
.nav-tooltip {
  display:none; position:absolute; left:calc(100% + 10px); top:50%; transform:translateY(-50%);
  background:var(--text-1); color:var(--bg); font-size:11.5px; font-weight:600;
  padding:4px 10px; border-radius:7px; white-space:nowrap; z-index:100;
  pointer-events:none;
}
.nav-tooltip::before {
  content:''; position:absolute; right:100%; top:50%; transform:translateY(-50%);
  border:5px solid transparent; border-right-color:var(--text-1);
}
/* Nav badge */
.nav-badge {
  background:#f43f5e; color:#fff;
  font-size:9px; font-weight:700;
  padding:1px 5px; border-radius:99px; min-width:18px; text-align:center;
}
.nav-badge-wrap { display:flex; align-items:center; }

.divider { height:1px; background:var(--border); margin:8px 10px; transition:background .3s; }
.nav-bottom { padding:0 10px; display:flex; flex-direction:column; gap:2px; }
.toggle-row {
  display:flex; align-items:center; gap:10px;
  padding:10px 12px 0;
}
.tog {
  width:42px; height:22px; border-radius:99px;
  background:var(--border); border:none; cursor:pointer; position:relative;
  transition:background .3s; flex-shrink:0; outline:none;
}
[data-theme="dark"] .tog { background:var(--accent); }
.tog::after {
  content:''; position:absolute; top:3px; left:3px;
  width:16px; height:16px; border-radius:50%;
  background:#fff; box-shadow:0 1px 4px rgba(0,0,0,.25);
  transition:transform .3s;
}
[data-theme="dark"] .tog::after { transform:translateX(20px); }
.tog-label { font-size:11.5px; color:var(--text-2); white-space:nowrap; }

/* Collapse button */
.collapse-btn {
  position:absolute; top:24px; right:-12px;
  width:24px; height:24px; border-radius:50%;
  background:var(--card-bg); border:1px solid var(--border);
  display:flex; align-items:center; justify-content:center;
  cursor:pointer; z-index:20;
  transition:all .2s, background .3s, border-color .3s;
  box-shadow:var(--shadow-card);
}
.collapse-btn:hover { background:var(--accent-soft); border-color:var(--accent); }
.collapse-btn svg { width:12px; height:12px; stroke:var(--text-2); fill:none; stroke-width:2; stroke-linecap:round; transition:stroke .2s, transform .25s; }
.collapse-btn:hover svg { stroke:var(--accent); }
.sidebar.collapsed .collapse-btn { right:-12px; }
.sidebar.collapsed .collapse-btn svg { transform:rotate(180deg); }

/* ═══════════════════════════════════════════
   MAIN LAYOUT
═══════════════════════════════════════════ */
.main { flex:1; display:flex; flex-direction:column; overflow:auto; height:100vh; }

/* TOPBAR */
.topbar {
  padding:20px 26px 0;
  display:flex; align-items:center; justify-content:space-between; gap:16px;
}
.topbar-left h1 { font-size:22px; font-weight:700; line-height:1.2; }
.topbar-left p  { font-size:12.5px; color:var(--text-2); margin-top:3px; }
.topbar-right   { display:flex; align-items:center; gap:8px; }

.icon-btn {
  width:36px; height:36px; border-radius:10px;
  background:var(--card-bg); border:1px solid var(--border);
  display:flex; align-items:center; justify-content:center;
  cursor:pointer; transition:all .2s; position:relative;
}
.icon-btn svg { width:16px; height:16px; stroke:var(--text-2); fill:none; stroke-width:1.8; stroke-linecap:round; stroke-linejoin:round; transition:stroke .3s; }
.icon-btn:hover { background:var(--accent-soft); border-color:var(--accent); }
.icon-btn:hover svg { stroke:var(--accent); }

.n-badge {
  position:absolute; top:-4px; right:-4px;
  width:15px; height:15px; border-radius:50%;
  background:#f43f5e; font-size:8px; font-weight:700; color:#fff;
  display:flex; align-items:center; justify-content:center;
  border:2px solid var(--card-bg); transition:border-color .3s;
}
.user-chip {
  display:flex; align-items:center; gap:8px;
  background:var(--card-bg); border:1px solid var(--border);
  border-radius:12px; padding:5px 12px 5px 5px;
  cursor:pointer; transition:all .2s;
}
.user-chip:hover { border-color:var(--accent); }
.u-av {
  width:28px; height:28px; border-radius:50%;
  background:linear-gradient(135deg,#c4b5fd,#7c5cfc);
  display:flex; align-items:center; justify-content:center;
  color:#fff; font-size:11px; font-weight:700; flex-shrink:0;
}
.u-name  { font-size:12.5px; font-weight:600; line-height:1.2; }
.u-email { font-size:10px; color:var(--text-2); }

/* FILTER BAR (dashboard) */
.filter-bar {
  padding:14px 26px 0;
  display:flex; align-items:center; justify-content:space-between;
}
.filter-l { display:flex; gap:6px; flex-wrap:wrap; }
.filter-r { display:flex; gap:6px; }
.pill {
  display:flex; align-items:center; gap:5px;
  padding:6px 14px; border-radius:99px;
  border:1px solid var(--border); background:var(--card-bg);
  font-size:12.5px; font-weight:500; color:var(--text-2);
  cursor:pointer; transition:all .2s; font-family:'Sora',sans-serif; outline:none;
}
.pill svg { width:13px; height:13px; stroke:currentColor; fill:none; stroke-width:1.8; }
.pill:hover, .pill.active { border-color:var(--accent); color:var(--accent); background:var(--accent-soft); }

.btn-outline {
  display:flex; align-items:center; gap:5px;
  padding:7px 14px; border-radius:10px;
  border:1px solid var(--border); background:var(--card-bg);
  font-size:12.5px; font-weight:500; color:var(--text-2);
  cursor:pointer; transition:all .2s; font-family:'Sora',sans-serif; outline:none;
}
.btn-outline svg { width:14px; height:14px; stroke:currentColor; fill:none; stroke-width:1.8; }
.btn-outline:hover { border-color:var(--accent); color:var(--accent); }

.btn-primary {
  display:flex; align-items:center; gap:5px;
  padding:7px 16px; border-radius:10px;
  border:none; background:var(--accent); color:#fff;
  font-size:12.5px; font-weight:600;
  cursor:pointer; transition:all .2s; font-family:'Sora',sans-serif; outline:none;
}
.btn-primary:hover { background:var(--accent-hover); transform:translateY(-1px); box-shadow:0 4px 14px rgba(124,92,252,.35); }
.btn-danger {
  display:flex; align-items:center; gap:5px;
  padding:7px 16px; border-radius:10px;
  border:none; background:#f43f5e; color:#fff;
  font-size:12.5px; font-weight:600;
  cursor:pointer; transition:all .2s; font-family:'Sora',sans-serif; outline:none;
}
.btn-danger:hover { background:#e11d48; transform:translateY(-1px); }

/* ═══════════════════════════════════════════
   PAGES (router)
═══════════════════════════════════════════ */
.page { display:none; flex:1; overflow:auto; flex-direction:column; }
.page.active { display:flex; }

/* ═══════════════════════════════════════════
   CONTENT GRID (dashboard)
═══════════════════════════════════════════ */
.content {
  padding:14px 26px 26px;
  display:grid; grid-template-columns:1fr 270px; gap:14px; flex:1;
}
.col-l { display:flex; flex-direction:column; gap:14px; }
.col-r { display:flex; flex-direction:column; gap:14px; }

/* STAT CARDS */
.stats-row { display:grid; grid-template-columns:repeat(4,1fr); gap:12px; }
.sc {
  background:var(--card-bg); border-radius:16px;
  padding:16px 18px; border:1px solid var(--border);
  box-shadow:var(--shadow-card);
  transition:transform .2s, box-shadow .2s, background .3s, border-color .3s;
}
.sc:hover { transform:translateY(-2px); box-shadow:0 8px 24px rgba(124,92,252,.12); }
.sc-head { display:flex; justify-content:space-between; align-items:center; margin-bottom:10px; }
.sc-lbl { font-size:12.5px; font-weight:500; color:var(--text-2); }
.sc-arrow { color:var(--text-3); }
.sc-arrow svg { width:14px; height:14px; stroke:currentColor; fill:none; stroke-width:1.8; transition:stroke .2s; }
.sc:hover .sc-arrow svg { stroke:var(--accent); }
.sc-val {
  font-size:20px; font-weight:700;
  font-family:'DM Mono',monospace; letter-spacing:-.5px;
  margin-bottom:8px; line-height:1;
}
.sc-val .c { font-size:13px; color:var(--text-2); font-weight:400; }
.chg {
  display:inline-flex; align-items:center; gap:3px;
  font-size:11px; font-weight:600;
  padding:3px 7px; border-radius:99px;
}
.chg svg { width:10px; height:10px; stroke:currentColor; fill:none; stroke-width:2.5; }
.chg.up { background:var(--green-bg); color:var(--green-text); }
.chg.dn { background:var(--red-bg);   color:var(--red-text); }
.chg-vs { font-size:11px; color:var(--text-2); margin-left:4px; }

/* GENERIC CARD */
.card {
  background:var(--card-bg); border-radius:16px;
  padding:18px 20px; border:1px solid var(--border);
  box-shadow:var(--shadow-card);
  transition:background .3s, border-color .3s;
}
.card-hd { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.card-hd h3 { font-size:14.5px; font-weight:700; }
.card-hd .ctl { display:flex; align-items:center; gap:8px; }
.leg { display:flex; gap:12px; }
.leg-i { display:flex; align-items:center; gap:5px; font-size:11.5px; color:var(--text-2); }
.leg-dot { width:8px; height:8px; border-radius:50%; flex-shrink:0; }
.sel {
  border:1px solid var(--border); border-radius:8px;
  padding:4px 8px; font-size:11.5px; color:var(--text-2);
  background:var(--chip-bg); cursor:pointer;
  font-family:'Sora',sans-serif; outline:none; transition:all .2s;
}
.sel:hover { border-color:var(--accent); }

/* CHART */
.chart-wrap { position:relative; padding-left:42px; height:150px; }
.chart-y {
  position:absolute; top:0; left:0; bottom:22px;
  display:flex; flex-direction:column; justify-content:space-between;
  pointer-events:none;
}
.yt { font-size:10px; color:var(--text-2); font-family:'DM Mono',monospace; text-align:right; width:38px; line-height:1; }
.chart-grid {
  position:absolute; top:0; left:42px; right:0; bottom:22px;
  display:flex; flex-direction:column; justify-content:space-between;
  pointer-events:none;
}
.gl { border-top:1px dashed var(--grid-line); transition:border-color .3s; }
.chart-bars {
  position:absolute; top:0; left:42px; right:0; bottom:22px;
  display:flex; align-items:flex-end; gap:8px;
}
.bc { flex:1; display:flex; align-items:flex-end; height:100%; }
.bp { display:flex; align-items:flex-end; gap:4px; width:100%; height:100%; }
.bar {
  flex:1; border-radius:6px 6px 0 0;
  position:relative; cursor:pointer; min-height:4px;
  transition:opacity .2s, filter .2s;
}
.bar:hover { opacity:.82; filter:brightness(1.1); }
.bar.inc { background:var(--bar-income); }
.bar.exp { background:var(--bar-expense); }
.btt {
  position:absolute; bottom:calc(100% + 6px); left:50%; transform:translateX(-50%);
  background:var(--text-1); color:var(--bg);
  font-size:10px; font-weight:600; padding:3px 7px; border-radius:6px;
  white-space:nowrap; display:none; pointer-events:none; z-index:20;
  transition:background .3s, color .3s;
}
.btt::after {
  content:''; position:absolute; top:100%; left:50%; transform:translateX(-50%);
  border:4px solid transparent; border-top-color:var(--text-1); transition:border-color .3s;
}
.bar:hover .btt { display:block; }
.chart-xl {
  position:absolute; bottom:0; left:42px; right:0;
  display:flex; gap:8px;
}
.xl { flex:1; text-align:center; font-size:10px; color:var(--text-2); font-family:'DM Mono',monospace; }

/* TRANSACTIONS TABLE */
.tx-tbl { width:100%; border-collapse:collapse; }
.tx-tbl thead th {
  padding:7px 10px; font-size:10.5px; font-weight:600;
  color:var(--text-2); text-transform:uppercase; letter-spacing:.6px;
  text-align:left; background:var(--th-bg); transition:background .3s;
}
.tx-tbl thead th:first-child { border-radius:8px 0 0 8px; }
.tx-tbl thead th:last-child  { border-radius:0 8px 8px 0; }
.tx-tbl tbody td {
  padding:10px 10px; font-size:12.5px;
  border-bottom:1px solid var(--border); vertical-align:middle;
  transition:background .3s, border-color .3s;
}
.tx-tbl tbody tr:last-child td { border-bottom:none; }
.tx-tbl tbody tr:hover td { background:var(--row-hover); }
.txi {
  width:30px; height:30px; border-radius:8px;
  display:inline-flex; align-items:center; justify-content:center;
  margin-right:8px; vertical-align:middle; flex-shrink:0;
}
.txn { font-weight:500; vertical-align:middle; }
.td-d  { color:var(--text-2); font-size:11.5px; white-space:nowrap; }
.neg   { color:#f43f5e; font-weight:600; font-family:'DM Mono',monospace; white-space:nowrap; }
.pos   { color:#22c55e; font-weight:600; font-family:'DM Mono',monospace; white-space:nowrap; }
.mchip {
  display:inline-flex; align-items:center; gap:4px;
  background:var(--chip-bg); border:1px solid var(--chip-border);
  padding:3px 8px; border-radius:6px;
  font-size:11px; color:var(--text-2); transition:all .3s;
}
.mchip svg { width:11px; height:11px; stroke:currentColor; fill:none; stroke-width:1.8; }
.cchip {
  background:var(--accent-soft); color:var(--accent);
  padding:3px 9px; border-radius:6px; font-size:11px; font-weight:500;
  white-space:nowrap;
}
.cchip.g { background:var(--green-bg); color:var(--green-text); }
.cchip.r { background:var(--red-bg);   color:var(--red-text); }

/* BUDGET / DONUT */
.bgt-wrap { display:flex; gap:14px; align-items:center; }
.donut-c { position:relative; width:106px; height:106px; flex-shrink:0; }
.donut-c svg { width:100%; height:100%; transform:rotate(-90deg); }
.donut-mid {
  position:absolute; top:50%; left:50%; transform:translate(-50%,-50%);
  text-align:center; pointer-events:none;
}
.donut-mid .dm-l { font-size:9px; color:var(--text-2); line-height:1.4; }
.donut-mid .dm-v { font-size:13px; font-weight:700; font-family:'DM Mono',monospace; }
.bgt-leg { display:flex; flex-direction:column; gap:6px; flex:1; }
.bli { display:flex; align-items:center; gap:7px; }
.bldot { width:7px; height:7px; border-radius:2px; flex-shrink:0; }
.blname { font-size:11.5px; color:var(--text-2); }

/* SAVING GOALS */
.sg { margin-bottom:13px; }
.sg:last-of-type { margin-bottom:0; }
.sg-top { display:flex; justify-content:space-between; align-items:baseline; margin-bottom:6px; }
.sg-n { font-size:13px; font-weight:600; }
.sg-a { font-size:13px; font-weight:700; font-family:'DM Mono',monospace; color:var(--accent); }
.prog-t { height:6px; background:var(--border); border-radius:99px; overflow:hidden; transition:background .3s; }
.prog-f { height:100%; border-radius:99px; background:var(--accent); transition:width .8s cubic-bezier(.4,0,.2,1); }
.sg-p { font-size:10.5px; color:var(--text-2); margin-top:4px; }

/* ═══════════════════════════════════════════
   PAGE WRAPPERS
═══════════════════════════════════════════ */
.page-header {
  padding:20px 26px 0;
  display:flex; align-items:center; justify-content:space-between; gap:16px;
}
.page-header h1 { font-size:22px; font-weight:700; }
.page-header p  { font-size:12.5px; color:var(--text-2); margin-top:3px; }
.page-body { padding:14px 26px 26px; display:flex; flex-direction:column; gap:14px; flex:1; }

/* TRANSACTIONS PAGE */
.tx-filters {
  display:flex; gap:8px; flex-wrap:wrap; align-items:center;
}
.tx-search {
  flex:1; min-width:180px;
  display:flex; align-items:center; gap:8px;
  background:var(--card-bg); border:1px solid var(--border);
  border-radius:10px; padding:7px 12px;
  transition:border-color .2s;
}
.tx-search:focus-within { border-color:var(--accent); }
.tx-search svg { width:14px; height:14px; stroke:var(--text-2); fill:none; stroke-width:1.8; flex-shrink:0; }
.tx-search input {
  border:none; background:transparent; outline:none;
  font-family:'Sora',sans-serif; font-size:12.5px; color:var(--text-1);
  width:100%;
}
.tx-search input::placeholder { color:var(--text-2); }
.filter-sel {
  border:1px solid var(--border); border-radius:10px;
  padding:7px 10px; font-size:12.5px; color:var(--text-2);
  background:var(--card-bg); cursor:pointer;
  font-family:'Sora',sans-serif; outline:none; transition:all .2s;
}
.filter-sel:hover, .filter-sel:focus { border-color:var(--accent); color:var(--text-1); }
input[type="date"].filter-sel { color:var(--text-1); }

.tx-empty {
  text-align:center; padding:40px 20px;
  color:var(--text-2); font-size:13px;
}
.tx-empty svg { width:40px; height:40px; stroke:var(--text-3); fill:none; stroke-width:1.4; margin-bottom:10px; }

.del-btn {
  background:none; border:none; cursor:pointer; padding:4px;
  color:var(--text-3); transition:color .2s; display:flex; align-items:center;
}
.del-btn:hover { color:#f43f5e; }
.del-btn svg { width:14px; height:14px; stroke:currentColor; fill:none; stroke-width:2; }

/* ANALYTICS PAGE */
.analytics-grid { display:grid; grid-template-columns:1fr 1fr; gap:14px; }
.analytics-wide { grid-column:1/-1; }
.stat-pill {
  background:var(--card-bg); border:1px solid var(--border);
  border-radius:16px; padding:16px 20px;
  box-shadow:var(--shadow-card);
}
.stat-pill h4 { font-size:12px; color:var(--text-2); font-weight:500; margin-bottom:8px; }
.stat-pill .big { font-size:24px; font-weight:700; font-family:'DM Mono',monospace; }
.cat-bars { display:flex; flex-direction:column; gap:10px; margin-top:4px; }
.cat-bar-row { display:flex; align-items:center; gap:10px; }
.cat-bar-label { font-size:12px; color:var(--text-2); width:110px; flex-shrink:0; }
.cat-bar-track { flex:1; height:8px; background:var(--border); border-radius:99px; overflow:hidden; }
.cat-bar-fill { height:100%; border-radius:99px; background:var(--accent); }
.cat-bar-val { font-size:11.5px; font-family:'DM Mono',monospace; color:var(--text-1); width:70px; text-align:right; flex-shrink:0; }

/* WALLET PAGE */
.wallet-cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(220px, 1fr)); gap:14px; margin-bottom:4px; }
.wallet-card {
  border-radius:18px; padding:20px 22px;
  color:#fff; position:relative; overflow:hidden;
  min-height:130px; display:flex; flex-direction:column; justify-content:space-between;
  box-shadow:0 8px 24px rgba(0,0,0,.15);
  cursor:pointer; transition:transform .2s;
}
.wallet-card:hover { transform:translateY(-3px); }
.wc-type { font-size:11px; font-weight:600; opacity:.8; text-transform:uppercase; letter-spacing:.5px; }
.wc-bal { font-size:22px; font-weight:700; font-family:'DM Mono',monospace; margin-top:8px; }
.wc-num { font-size:12px; opacity:.7; letter-spacing:1px; margin-top:4px; }
.wc-shine {
  position:absolute; top:-30px; right:-30px;
  width:100px; height:100px; border-radius:50%;
  background:rgba(255,255,255,.08);
}
.wc-shine2 {
  position:absolute; bottom:-40px; left:-20px;
  width:120px; height:120px; border-radius:50%;
  background:rgba(255,255,255,.05);
}

/* GOALS PAGE */
.goals-grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(260px,1fr)); gap:14px; }
.goal-card {
  background:var(--card-bg); border:1px solid var(--border);
  border-radius:16px; padding:18px 20px;
  box-shadow:var(--shadow-card);
}
.goal-icon { font-size:28px; margin-bottom:10px; display:block; }
.goal-name { font-size:14px; font-weight:700; margin-bottom:2px; }
.goal-target { font-size:11.5px; color:var(--text-2); margin-bottom:12px; }
.goal-amount { font-family:'DM Mono',monospace; font-size:13px; font-weight:700; color:var(--accent); }
.goal-pct { font-size:11px; color:var(--text-2); margin-top:4px; }

/* BUDGET PAGE */
.budget-page-grid { display:grid; grid-template-columns:1fr 260px; gap:14px; }
.budget-cat-row {
  display:flex; align-items:center; gap:12px;
  padding:12px 0; border-bottom:1px solid var(--border);
}
.budget-cat-row:last-child { border-bottom:none; }
.budget-cat-dot { width:10px; height:10px; border-radius:3px; flex-shrink:0; }
.budget-cat-name { font-size:13px; font-weight:500; flex:1; }
.budget-cat-bar-wrap { width:120px; height:6px; background:var(--border); border-radius:99px; overflow:hidden; }
.budget-cat-bar { height:100%; border-radius:99px; }
.budget-cat-amt { font-size:12px; font-family:'DM Mono',monospace; color:var(--text-1); width:70px; text-align:right; }

/* SETTINGS PAGE */
.settings-section { margin-bottom:6px; }
.settings-section h4 { font-size:12px; font-weight:600; text-transform:uppercase; letter-spacing:.7px; color:var(--text-2); margin-bottom:10px; }
.settings-row {
  display:flex; align-items:center; justify-content:space-between;
  padding:12px 0; border-bottom:1px solid var(--border);
}
.settings-row:last-child { border-bottom:none; }
.settings-row-left h5 { font-size:13px; font-weight:600; }
.settings-row-left p  { font-size:11.5px; color:var(--text-2); margin-top:2px; }

/* ═══════════════════════════════════════════
   MODAL SYSTEM
═══════════════════════════════════════════ */
.modal-overlay {
  display:none; position:fixed; inset:0;
  background:var(--modal-overlay);
  z-index:1000; align-items:center; justify-content:center;
  backdrop-filter:blur(4px);
}
.modal-overlay.open { display:flex; }
.modal {
  background:var(--card-bg); border-radius:20px;
  padding:28px; width:90%; max-width:480px;
  border:1px solid var(--border);
  box-shadow:0 24px 60px rgba(0,0,0,.2);
  animation:modalIn .25s cubic-bezier(.34,1.56,.64,1) both;
  max-height:90vh; overflow-y:auto;
}
@keyframes modalIn {
  from { opacity:0; transform:scale(.94) translateY(12px); }
  to   { opacity:1; transform:scale(1) translateY(0); }
}
.modal-hd {
  display:flex; align-items:center; justify-content:space-between; margin-bottom:22px;
}
.modal-hd h2 { font-size:17px; font-weight:700; }
.modal-close {
  width:30px; height:30px; border-radius:8px;
  background:var(--chip-bg); border:none; cursor:pointer;
  display:flex; align-items:center; justify-content:center;
  transition:all .2s;
}
.modal-close:hover { background:var(--red-bg); }
.modal-close svg { width:14px; height:14px; stroke:var(--text-2); fill:none; stroke-width:2.2; stroke-linecap:round; }
.modal-close:hover svg { stroke:var(--red-text); }

/* Form */
.form-row { margin-bottom:16px; }
.form-row label { display:block; font-size:12px; font-weight:600; color:var(--text-2); margin-bottom:6px; text-transform:uppercase; letter-spacing:.5px; }
.form-input, .form-select, .form-textarea {
  width:100%; padding:10px 14px;
  background:var(--input-bg); border:1px solid var(--border);
  border-radius:10px; font-family:'Sora',sans-serif; font-size:13px; color:var(--text-1);
  outline:none; transition:border-color .2s;
}
.form-input:focus, .form-select:focus, .form-textarea:focus { border-color:var(--accent); }
.form-textarea { resize:vertical; min-height:70px; }
.form-error { font-size:11px; color:var(--red-text); margin-top:4px; display:none; }
.form-error.show { display:block; }
.form-input.err, .form-select.err { border-color:var(--red-text); }
/* type toggle */
.type-toggle { display:flex; gap:8px; }
.type-btn {
  flex:1; padding:9px; border-radius:10px;
  border:2px solid var(--border); background:var(--input-bg);
  font-family:'Sora',sans-serif; font-size:13px; font-weight:600; cursor:pointer;
  transition:all .2s; color:var(--text-2);
}
.type-btn.sel-income  { border-color:var(--green-text); background:var(--green-bg); color:var(--green-text); }
.type-btn.sel-expense { border-color:var(--red-text); background:var(--red-bg); color:var(--red-text); }
.form-actions { display:flex; gap:8px; margin-top:22px; justify-content:flex-end; }

/* TOAST */
.toast {
  position:fixed; bottom:24px; right:24px; z-index:2000;
  background:var(--text-1); color:var(--bg);
  padding:12px 20px; border-radius:12px;
  font-size:13px; font-weight:600;
  box-shadow:0 8px 24px rgba(0,0,0,.2);
  animation:toastIn .3s ease both;
  pointer-events:none;
}
@keyframes toastIn {
  from { opacity:0; transform:translateY(20px); }
  to   { opacity:1; transform:translateY(0); }
}
.toast.out { animation:toastOut .3s ease both; }
@keyframes toastOut {
  to { opacity:0; transform:translateY(20px); }
}

/* ═══════════════════════════════════════════
   ANIMATIONS
═══════════════════════════════════════════ */
@keyframes fadeUp {
  from { opacity:0; transform:translateY(12px); }
  to   { opacity:1; transform:translateY(0); }
}
.anim { animation:fadeUp .4s ease both; }

/* ── Analytics redesign ─────────────────────────── */
.an-stat-card { padding: 18px 20px; }
.an-card-label { font-size: 13px; font-weight: 600; color: var(--text-1); }
.an-currency-badge { font-size: 11px; padding: 2px 8px; border-radius: 20px; border: 1px solid var(--border); color: var(--text-2); }
.an-big-val { font-family: 'DM Mono', monospace; font-size: 30px; font-weight: 700; color: var(--text-1); margin: 8px 0 6px; line-height: 1; }
.an-big-val .an-cents { font-size: 18px; color: var(--text-2); }
.an-card-meta { display: flex; gap: 16px; align-items: center; flex-wrap: wrap; margin-bottom: 4px; }
.an-change { font-size: 12px; font-weight: 600; }
.an-change.pos { color: var(--green-text); }
.an-change.neg { color: var(--red-text); }
.an-tx-count { font-size: 12px; color: var(--text-2); }
.an-card-sub { font-size: 12px; color: var(--text-2); margin-top: 4px; }
.an-cat-count { font-size: 12px; color: var(--text-2); margin-top: 2px; }
@media (max-width: 900px) {
  .an-charts-row { grid-template-columns: 1fr !important; }
}
</style>
</head>
<body>

<!-- ═══════════════════════════════════════════
     SIDEBAR
═══════════════════════════════════════════ -->
<aside class="sidebar" id="sidebar">

  <!-- Collapse button -->
  <button class="collapse-btn" id="collapseBtn" title="Collapse sidebar">
    <svg viewBox="0 0 24 24"><polyline points="15 18 9 12 15 6"/></svg>
  </button>

  <!-- Logo -->
  <div class="logo">
    <div class="logo-icon">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="2" y="3" width="20" height="14" rx="3"/><path d="M8 21h8M12 17v4"/>
      </svg>
    </div>
    <span>FinSet</span>
  </div>

  <!-- Main nav -->
  <nav class="nav" id="mainNav">
    <a class="nav-item active" data-page="dashboard">
      <svg viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/></svg>
      <span class="nav-label">Dashboard</span>
      <span class="nav-tooltip">Dashboard</span>
    </a>
    <a class="nav-item" data-page="transactions">
      <svg viewBox="0 0 24 24"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
      <span class="nav-label">Transactions</span>
      <span class="nav-badge-wrap"><span class="nav-badge" id="txBadge">0</span></span>
      <span class="nav-tooltip">Transactions</span>
    </a>
    <a class="nav-item" data-page="wallet">
      <svg viewBox="0 0 24 24"><rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/></svg>
      <span class="nav-label">Wallet</span>
      <span class="nav-tooltip">Wallet</span>
    </a>
    <a class="nav-item" data-page="goals">
      <svg viewBox="0 0 24 24"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
      <span class="nav-label">Goals</span>
      <span class="nav-badge-wrap"><span class="nav-badge" id="goalsBadge">0</span></span>
      <span class="nav-tooltip">Goals</span>
    </a>
    <a class="nav-item" data-page="budget">
      <svg viewBox="0 0 24 24"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/></svg>
      <span class="nav-label">Budget</span>
      <span class="nav-tooltip">Budget</span>
    </a>
    <a class="nav-item" data-page="analytics">
      <svg viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
      <span class="nav-label">Analytics</span>
      <span class="nav-tooltip">Analytics</span>
    </a>
    <a class="nav-item" data-page="settings">
      <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.07 4.93a10 10 0 010 14.14M4.93 4.93a10 10 0 000 14.14"/></svg>
      <span class="nav-label">Settings</span>
      <span class="nav-tooltip">Settings</span>
    </a>
  </nav>

  <div class="divider"></div>

  <nav class="nav-bottom">
    <a class="nav-item">
      <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3M12 17h.01"/></svg>
      <span class="nav-label">Help</span>
      <span class="nav-tooltip">Help</span>
    </a>
    <a class="nav-item">
      <svg viewBox="0 0 24 24"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
      <span class="nav-label">Log out</span>
      <span class="nav-tooltip">Log out</span>
    </a>
  </nav>

  <!-- Theme toggle -->
  <div class="toggle-row">
    <button class="tog" id="themeToggle" aria-label="Toggle theme"></button>
    <span class="tog-label" id="themeLabel">Light mode</span>
  </div>
</aside>

<!-- ═══════════════════════════════════════════
     MAIN
═══════════════════════════════════════════ -->
<main class="main" id="mainContent">

  <!-- ─────────────── DASHBOARD PAGE ─────────────── -->
  <section class="page active" data-page="dashboard">
    <!-- Topbar -->
    <div class="topbar">
      <div class="topbar-left">
        <h1>Welcome back, Adaline! 👋</h1>
        <p>It is the best time to manage your finances.</p>
      </div>
      <div class="topbar-right">
        <div class="icon-btn">
          <svg viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </div>
        <div class="icon-btn" style="position:relative;">
          <svg viewBox="0 0 24 24"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>
          <span class="n-badge">2</span>
        </div>
        <div class="user-chip">
          <div class="u-av">AL</div>
          <div>
            <div class="u-name">Adaline Lively</div>
            <div class="u-email"><a href="/cdn-cgi/l/email-protection" class="__cf_email__" data-cfemail="a0c1c4c1ccc9cec5e0c5cdc1c9cc8ec3cfcd">[email&#160;protected]</a></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter bar -->
    <div class="filter-bar">
      <div class="filter-l">
        <button class="pill active" data-period="month">
          <svg viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
          This month
        </button>
        <button class="pill" data-period="week">Last 7 days</button>
        <button class="pill" data-period="all">All time</button>
      </div>
      <div class="filter-r">
        <button class="btn-outline" id="manageWidgetsBtn">
          <svg viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
          Manage widgets
        </button>
        <button class="btn-primary" id="addTxBtnDash">
          <svg viewBox="0 0 24 24" width="13" height="13" stroke="#fff" fill="none" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          Add transaction
        </button>
      </div>
    </div>

    <!-- Content grid -->
    <div class="content">
      <div class="col-l">
        <!-- Stat cards (data-driven) -->
        <div class="stats-row" id="dashStats"></div>

        <!-- Money flow chart -->
        <div class="card">
          <div class="card-hd">
            <h3>Money flow</h3>
            <div class="ctl">
              <div class="leg">
                <div class="leg-i"><div class="leg-dot" style="background:var(--bar-income)"></div>Income</div>
                <div class="leg-i"><div class="leg-dot" style="background:var(--bar-expense)"></div>Expense</div>
              </div>
            </div>
          </div>
          <div class="chart-wrap">
            <div class="chart-y" id="chartY">
              <span class="yt" id="ytMax"></span>
              <span class="yt" id="ytMid"></span>
              <span class="yt" id="ytLow"></span>
              <span class="yt">$0</span>
            </div>
            <div class="chart-grid">
              <div class="gl"></div><div class="gl"></div><div class="gl"></div><div class="gl"></div>
            </div>
            <div class="chart-bars" id="cBars"></div>
            <div class="chart-xl"  id="cXL"></div>
          </div>
        </div>

        <!-- Recent transactions -->
        <div class="card">
          <div class="card-hd">
            <h3>Recent transactions</h3>
            <div class="ctl">
              <button class="btn-outline" id="seeAllBtn" style="padding:5px 12px;font-size:11.5px;">See all →</button>
            </div>
          </div>
          <table class="tx-tbl">
            <thead>
              <tr>
                <th>Date</th><th>Amount</th><th>Description</th><th>Method</th><th>Category</th>
              </tr>
            </thead>
            <tbody id="dashTxBody"></tbody>
          </table>
        </div>
      </div>

      <div class="col-r">
        <!-- Budget -->
        <div class="card">
          <div class="card-hd">
            <h3>Budget</h3>
            <span style="cursor:pointer;color:var(--text-3);"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M7 17L17 7"/><path d="M7 7h10v10"/></svg></span>
          </div>
          <div class="bgt-wrap">
            <div class="donut-c" id="donutWrap">
              <svg viewBox="0 0 36 36" id="donutSvg">
                <circle cx="18" cy="18" r="14" fill="none" stroke="var(--border)" stroke-width="4.5"/>
              </svg>
              <div class="donut-mid">
                <div class="dm-l">Total<br>this month</div>
                <div class="dm-v" id="donutTotal">$0</div>
              </div>
            </div>
            <div class="bgt-leg" id="donutLegend"></div>
          </div>
        </div>

        <!-- Saving goals -->
        <div class="card">
          <div class="card-hd">
            <h3>Saving goals</h3>
            <button class="btn-outline" onclick="setActivePage('goals')" style="padding:4px 10px;font-size:11px;">View all</button>
          </div>
          <div class="sg">
            <div class="sg-top"><span class="sg-n">💻 MacBook Pro</span><span class="sg-a">$1,650</span></div>
            <div class="prog-t"><div class="prog-f" data-w="25" style="width:0%"></div></div>
            <div class="sg-p">25% of $6,600</div>
          </div>
          <div class="sg">
            <div class="sg-top"><span class="sg-n">🚗 New car</span><span class="sg-a">$60,000</span></div>
            <div class="prog-t"><div class="prog-f" data-w="42" style="width:0%"></div></div>
            <div class="sg-p">42% of $142,857</div>
          </div>
          <div class="sg">
            <div class="sg-top"><span class="sg-n">🏠 New house</span><span class="sg-a">$150,000</span></div>
            <div class="prog-t"><div class="prog-f" data-w="3" style="width:0%"></div></div>
            <div class="sg-p">3% of $5,000,000</div>
          </div>
          <button class="btn-primary" style="width:100%;justify-content:center;margin-top:12px;" onclick="setActivePage('goals')">
            <svg viewBox="0 0 24 24" width="13" height="13" stroke="#fff" fill="none" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            Add saving goal
          </button>
        </div>
      </div>
    </div>
  </section>

  <!-- ─────────────── TRANSACTIONS PAGE ─────────────── -->
  <section class="page" data-page="transactions">
    <div class="page-header">
      <div>
        <h1>Transactions</h1>
        <p>All your income and expenses in one place.</p>
      </div>
      <button class="btn-primary" id="addTxBtnTx">
        <svg viewBox="0 0 24 24" width="13" height="13" stroke="#fff" fill="none" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Add transaction
      </button>
    </div>
    <div class="page-body">
      <!-- Filters -->
      <div class="tx-filters">
        <div class="tx-search">
          <svg viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <input type="text" placeholder="Search by note or category…" id="txSearch">
        </div>
        <select class="filter-sel" id="txTypeFilter">
          <option value="">All types</option>
          <option value="income">Income</option>
          <option value="expense">Expense</option>
        </select>
        <input type="date" class="filter-sel" id="txDateFrom" title="From date">
        <input type="date" class="filter-sel" id="txDateTo"   title="To date">
      </div>
      <!-- Table -->
      <div class="card" style="padding:0;overflow:hidden;">
        <table class="tx-tbl">
          <thead>
            <tr>
              <th style="padding-left:20px;">Date ↓</th>
              <th>Amount</th>
              <th>Description</th>
              <th>Method</th>
              <th>Category</th>
              <th>Type</th>
              <th></th>
            </tr>
          </thead>
          <tbody id="txPageBody"></tbody>
        </table>
        <div id="txEmpty" class="tx-empty" style="display:none;">
          <svg viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
          <div>No transactions found.</div>
        </div>
      </div>
    </div>
  </section>

  <!-- ─────────────── WALLET PAGE ─────────────── -->
  <section class="page" data-page="wallet">
    <div class="page-header">
      <div><h1>Wallet</h1><p>Your connected accounts and cards.</p></div>
      <button class="btn-primary">
        <svg viewBox="0 0 24 24" width="13" height="13" stroke="#fff" fill="none" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Add account
      </button>
    </div>
    <div class="page-body">
      <div class="wallet-cards">
        <div class="wallet-card" style="background:linear-gradient(135deg,#7c5cfc,#a78bfa);">
          <div class="wc-shine"></div><div class="wc-shine2"></div>
          <div class="wc-type">Main account</div>
          <div class="wc-bal" id="walletMainBal">$0.00</div>
          <div class="wc-num">•••• •••• •••• 3254</div>
        </div>
        <div class="wallet-card" style="background:linear-gradient(135deg,#0ea5e9,#38bdf8);">
          <div class="wc-shine"></div><div class="wc-shine2"></div>
          <div class="wc-type">Total Income</div>
          <div class="wc-bal" id="walletTotalInc">$0.00</div>
          <div class="wc-num">All time earnings</div>
        </div>
        <div class="wallet-card" style="background:linear-gradient(135deg,#f43f5e,#fb7185);">
          <div class="wc-shine"></div><div class="wc-shine2"></div>
          <div class="wc-type">Total Expenses</div>
          <div class="wc-bal" id="walletTotalExp">$0.00</div>
          <div class="wc-num">All time spending</div>
        </div>
      </div>
      <div class="card">
        <div class="card-hd"><h3>Recent wallet activity</h3></div>
        <table class="tx-tbl"><thead><tr><th>Date</th><th>Amount</th><th>Description</th><th>Method</th><th>Category</th></tr></thead>
        <tbody id="walletTxBody"></tbody></table>
      </div>
    </div>
  </section>

  <!-- ─────────────── GOALS PAGE ─────────────── -->
  <div class="modal-overlay" id="goalModal">
    <div class="modal">
      <div class="modal-hd">
        <h2 id="goalModalTitle">Add goal</h2>
        <button class="modal-close" id="goalModalClose">
          <svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="form-row">
        <label>Icon (emoji)</label>
        <input class="form-input" type="text" id="gIcon" placeholder="💻" maxlength="2" value="🎯">
      </div>
      <div class="form-row">
        <label>Goal name *</label>
        <input class="form-input" type="text" id="gName" placeholder="e.g. MacBook Pro">
        <div class="form-error" id="gNameErr">Please enter a goal name.</div>
      </div>
      <div class="form-row">
        <label>Target amount *</label>
        <input class="form-input" type="number" id="gTarget" placeholder="0.00" min="1" step="0.01">
        <div class="form-error" id="gTargetErr">Please enter a valid target amount.</div>
      </div>
      <div class="form-row">
        <label>Amount saved so far</label>
        <input class="form-input" type="number" id="gSaved" placeholder="0.00" min="0" step="0.01" value="0">
      </div>
      <input type="hidden" id="gEditId">
      <div class="form-actions">
        <button class="btn-outline" id="goalModalCancel">Cancel</button>
        <button class="btn-primary" id="goalModalSave">Save goal</button>
      </div>
    </div>
  </div>

  <section class="page" data-page="goals">
    <div class="page-header">
      <div><h1>Saving goals</h1><p>Track your financial milestones.</p></div>
      <button class="btn-primary" id="addGoalBtn">
        <svg viewBox="0 0 24 24" width="13" height="13" stroke="#fff" fill="none" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Add goal
      </button>
    </div>
    <div class="page-body">
      <div class="goals-grid" id="goalsGrid"></div>
    </div>
  </section>

  <!-- ─────────────── BUDGET PAGE ─────────────── -->
  <section class="page" data-page="budget">
    <div class="page-header">
      <div><h1>Budget</h1><p>Monthly spending by category.</p></div>
    </div>
    <div class="page-body">
      <div class="budget-page-grid">
        <div class="card">
          <div class="card-hd"><h3>Spending by category</h3></div>
          <div id="budgetCatList"></div>
        </div>
        <div class="col-r" style="gap:14px;display:flex;flex-direction:column;">
          <div class="card">
            <div class="card-hd"><h3>Budget overview</h3></div>
            <div class="bgt-wrap">
              <div class="donut-c">
                <svg viewBox="0 0 36 36" id="budgetDonutSvg">
                  <circle cx="18" cy="18" r="14" fill="none" stroke="var(--border)" stroke-width="4.5"/>
                </svg>
                <div class="donut-mid">
                  <div class="dm-l">Total<br>expenses</div>
                  <div class="dm-v" id="budgetDonutTotal">$0</div>
                </div>
              </div>
              <div class="bgt-leg" id="budgetDonutLegend"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- ─────────────── ANALYTICS PAGE ─────────────── -->
  <section class="page" data-page="analytics">
    <div class="page-header" style="flex-wrap:wrap;gap:10px;">
      <div>
        <h1>Analytics</h1>
        <p>Detailed overview of your financial situation</p>
      </div>
      <div style="display:flex;gap:8px;align-items:center;">
        <button class="pill active" id="anPeriodMonth" data-anperiod="month" style="padding:6px 14px;">This month</button>
        <button class="pill" id="anPeriodAll" data-anperiod="all" style="padding:6px 14px;">All time</button>
        <button class="btn-outline" id="anExportCsv" style="padding:6px 12px;font-size:12px;">
          <svg viewBox="0 0 24 24" width="13" height="13" stroke="currentColor" fill="none" stroke-width="1.8"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Export CSV
        </button>
      </div>
    </div>
    <div class="page-body">

      <!-- Top stat cards -->
      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:14px;margin-bottom:16px;">
        <div class="card an-stat-card" id="anCardBalance">
          <div style="display:flex;justify-content:space-between;align-items:flex-start;">
            <div class="an-card-label">Total balance</div>
            <span class="an-currency-badge">USD</span>
          </div>
          <div class="an-big-val" id="anTotalBal">$0<span class="an-cents">.00</span></div>
          <div class="an-card-meta">
            <span class="an-change pos" id="anBalChange"></span>
            <span class="an-tx-count" id="anBalTx"></span>
          </div>
          <div class="an-card-sub" id="anBalSub"></div>
          <div class="an-cat-count" id="anBalCats"></div>
        </div>
        <div class="card an-stat-card" id="anCardIncome">
          <div style="display:flex;justify-content:space-between;align-items:flex-start;">
            <div class="an-card-label">Income</div>
            <span class="an-currency-badge">USD</span>
          </div>
          <div class="an-big-val" id="anTotalInc">$0<span class="an-cents">.00</span></div>
          <div class="an-card-meta">
            <span class="an-change pos" id="anIncChange"></span>
            <span class="an-tx-count" id="anIncTx"></span>
          </div>
          <div class="an-card-sub" id="anIncSub"></div>
          <div class="an-cat-count" id="anIncCats"></div>
        </div>
        <div class="card an-stat-card" id="anCardExpense">
          <div style="display:flex;justify-content:space-between;align-items:flex-start;">
            <div class="an-card-label">Expense</div>
            <span class="an-currency-badge">USD</span>
          </div>
          <div class="an-big-val" id="anTotalExp">$0<span class="an-cents">.00</span></div>
          <div class="an-card-meta">
            <span class="an-change neg" id="anExpChange"></span>
            <span class="an-tx-count" id="anExpTx"></span>
          </div>
          <div class="an-card-sub" id="anExpSub"></div>
          <div class="an-cat-count" id="anExpCats"></div>
        </div>
      </div>

      <!-- Charts row -->
      <div style="display:grid;grid-template-columns:1fr 300px;gap:14px;margin-bottom:16px;" class="an-charts-row">

        <!-- Line chart: Balance overview -->
        <div class="card" style="padding:20px;">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;flex-wrap:wrap;gap:8px;">
            <h3 style="font-size:14px;font-weight:700;">Total balance overview</h3>
            <div style="display:flex;gap:16px;align-items:center;font-size:11px;color:var(--text-2);">
              <span><span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:var(--accent);margin-right:4px;"></span>This month</span>
              <span><span style="display:inline-block;width:10px;height:10px;border-radius:50%;border:2px dashed var(--text-3);margin-right:4px;"></span>Last month</span>
            </div>
          </div>
          <div style="position:relative;height:200px;overflow:hidden;">
            <canvas id="anLineChart" style="width:100%;height:100%;"></canvas>
            <div id="anLineEmpty" style="display:none;position:absolute;inset:0;display:flex;align-items:center;justify-content:center;color:var(--text-2);font-size:13px;">Add transactions to see chart</div>
          </div>
          <div id="anLineXL" style="display:flex;justify-content:space-between;margin-top:8px;font-size:11px;color:var(--text-2);padding:0 4px;"></div>
        </div>

        <!-- Statistics donut panel -->
        <div class="card" style="padding:20px;">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px;">
            <h3 style="font-size:14px;font-weight:700;">Statistics</h3>
          </div>
          <p style="font-size:11px;color:var(--text-2);margin-bottom:16px;" id="anStatSubtitle">Expense breakdown by category</p>
          <div style="position:relative;width:160px;height:160px;margin:0 auto 16px;">
            <svg viewBox="0 0 36 36" id="anDonutSvg" style="width:100%;height:100%;transform:rotate(-90deg);">
              <circle cx="18" cy="18" r="13" fill="none" stroke="var(--border)" stroke-width="5"/>
            </svg>
            <div style="position:absolute;inset:0;display:flex;flex-direction:column;align-items:center;justify-content:center;text-align:center;">
              <div style="font-size:10px;color:var(--text-2);line-height:1.3;">This month<br>expense</div>
              <div style="font-family:'DM Mono',monospace;font-size:15px;font-weight:700;color:var(--text-1);margin-top:2px;" id="anDonutTotal">$0</div>
            </div>
          </div>
          <div id="anDonutLegend" style="display:flex;flex-wrap:wrap;gap:6px;"></div>
        </div>
      </div>

      <!-- Bar chart: Monthly income vs expense -->
      <div class="card" style="padding:20px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;flex-wrap:wrap;gap:8px;">
          <h3 style="font-size:14px;font-weight:700;">Income vs Expense by month</h3>
          <div style="display:flex;gap:16px;font-size:11px;color:var(--text-2);">
            <span><span style="display:inline-block;width:10px;height:10px;border-radius:3px;background:var(--accent);margin-right:4px;"></span>Income</span>
            <span><span style="display:inline-block;width:10px;height:10px;border-radius:3px;background:var(--bar-expense);margin-right:4px;"></span>Expense</span>
          </div>
        </div>
        <div class="chart-wrap" style="height:180px;">
          <div class="chart-y">
            <span class="yt" id="anYtMax"></span>
            <span class="yt" id="anYtMid"></span>
            <span class="yt" id="anYtLow"></span>
            <span class="yt">$0</span>
          </div>
          <div class="chart-grid"><div class="gl"></div><div class="gl"></div><div class="gl"></div><div class="gl"></div></div>
          <div class="chart-bars" id="anBars"></div>
          <div class="chart-xl" id="anXL"></div>
        </div>
      </div>

    </div>
  </section>

  <!-- ─────────────── SETTINGS PAGE ─────────────── -->
  <section class="page" data-page="settings">
    <div class="page-header">
      <div><h1>Settings</h1><p>Manage your app preferences and data.</p></div>
    </div>
    <div class="page-body">
      <div class="card" style="max-width:640px;">
        <div class="settings-section">
          <h4>Appearance</h4>
          <div class="settings-row">
            <div class="settings-row-left">
              <h5>Dark mode</h5>
              <p>Switch between light and dark theme</p>
            </div>
            <div style="display:flex;align-items:center;gap:10px;">
              <button class="tog" id="themeToggle2" aria-label="Toggle theme"></button>
              <span class="tog-label" id="themeLabel2">Light mode</span>
            </div>
          </div>
        </div>
        <div class="settings-section" style="margin-top:10px;">
          <h4>Data management</h4>
          <div class="settings-row">
            <div class="settings-row-left">
              <h5>Export transactions</h5>
              <p>Download all your transactions as a JSON file</p>
            </div>
            <button class="btn-outline" id="exportBtn">
              <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="1.8"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Export JSON
            </button>
          </div>
          <div class="settings-row">
            <div class="settings-row-left">
              <h5>Import transactions</h5>
              <p>Upload a previously exported JSON file</p>
            </div>
            <label class="btn-outline" style="cursor:pointer;">
              <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="1.8"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              Import JSON
              <input type="file" accept=".json" id="importFile" style="display:none;">
            </label>
          </div>
          <div class="settings-row">
            <div class="settings-row-left">
              <h5>Reset to demo data</h5>
              <p>Clears all data and reloads the demo dataset</p>
            </div>
            <button class="btn-danger" id="resetDemoBtn">
              <svg viewBox="0 0 24 24" width="14" height="14" stroke="#fff" fill="none" stroke-width="2"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
              Reset demo
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>

</main><!-- end #mainContent -->

<!-- ═══════════════════════════════════════════
     ADD TRANSACTION MODAL
═══════════════════════════════════════════ -->
<div class="modal-overlay" id="txModal">
  <div class="modal">
    <div class="modal-hd">
      <h2>Add transaction</h2>
      <button class="modal-close" id="modalClose">
        <svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>
    <div class="form-row">
      <label>Type</label>
      <div class="type-toggle">
        <button class="type-btn sel-income" id="typeBtnIncome">↑ Income</button>
        <button class="type-btn"            id="typeBtnExpense">↓ Expense</button>
      </div>
    </div>
    <div class="form-row">
      <label>Amount *</label>
      <input class="form-input" type="number" id="fAmount" placeholder="0.00" min="0.01" step="0.01">
      <div class="form-error" id="fAmountErr">Please enter a valid positive amount.</div>
    </div>
    <div class="form-row">
      <label>Category *</label>
      <select class="form-select" id="fCategory">
        <option value="">Select category…</option>
        <option>Food</option><option>Transport</option><option>Bills</option>
        <option>Shopping</option><option>Health</option><option>Education</option>
        <option>Entertainment</option><option>Salary</option><option>Freelance</option>
        <option>Investment</option><option>Other</option>
      </select>
      <div class="form-error" id="fCategoryErr">Please select a category.</div>
    </div>
    <div class="form-row" id="fCustomCatRow" style="display:none;">
      <label>Custom category</label>
      <input class="form-input" type="text" id="fCustomCat" placeholder="Enter category name…">
    </div>
    <div class="form-row">
      <label>Payment method</label>
      <select class="form-select" id="fMethod">
        <option>Cash</option><option>Card</option><option>Transfer</option>
      </select>
    </div>
    <div class="form-row">
      <label>Date</label>
      <input class="form-input" type="date" id="fDate">
    </div>
    <div class="form-row">
      <label>Note (optional)</label>
      <textarea class="form-textarea" id="fNote" placeholder="Add a note…"></textarea>
    </div>
    <div class="form-actions">
      <button class="btn-outline" id="modalCancelBtn">Cancel</button>
      <button class="btn-primary" id="modalSubmitBtn">Save transaction</button>
    </div>
  </div>
</div>

<!-- TOAST -->
<div id="toastContainer"></div>

<!-- ═══════════════════════════════════════════
     JAVASCRIPT — API-DRIVEN (PostgreSQL via Go backend)
═══════════════════════════════════════════ -->
<script>
'use strict';

/* ══════════════════════════════════
   API CLIENT
   All data comes from /api/* — no localStorage for transactions.
   Theme + UI state (sidebar, page) still use localStorage.
══════════════════════════════════ */
const API = {
  async get(path) {
    const r = await fetch(path);
    if (!r.ok) throw new Error(` + "`" + `GET ${path} → ${r.status}` + "`" + `);
    return r.json();
  },
  async post(path, body) {
    const r = await fetch(path, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify(body),
    });
    const data = await r.json();
    if (!r.ok) throw new Error(data.error || ` + "`" + `POST ${path} → ${r.status}` + "`" + `);
    return data;
  },
  async del(path) {
    const r = await fetch(path, {method:'DELETE'});
    if (r.status === 404) throw new Error('not found');
    if (!r.ok) throw new Error(` + "`" + `DELETE ${path} → ${r.status}` + "`" + `);
  },
};

/* ══════════════════════════════════
   CONSTANTS
══════════════════════════════════ */
const LS_THEME   = 'finset-theme';
const LS_PAGE    = 'finset-page';
const LS_SIDEBAR = 'finset-sidebar';
const PAGES      = ['dashboard','transactions','wallet','goals','budget','analytics','settings'];
const CAT_COLORS = {
  'Food':'#7c5cfc','Transport':'#9567f8','Bills':'#a78bfa','Shopping':'#b08cf6',
  'Health':'#c9aaf5','Education':'#dfc8f4','Entertainment':'#38bdf8',
  'Salary':'#22c55e','Freelance':'#4ade80','Investment':'#f59e0b','Other':'#6b7280',
};

/* ══════════════════════════════════
   IN-MEMORY CACHE
   Transactions are fetched once per page-visit and cached.
   Cache is invalidated after write operations.
══════════════════════════════════ */
let _txCache = null;   // null = dirty, must re-fetch
let _activePeriod = 'month';
let _selectedType  = 'income';

async function getTx() {
  if (_txCache === null) {
    _txCache = await API.get('/api/transactions');
  }
  return _txCache;
}
function invalidateCache() { _txCache = null; }

/* ══════════════════════════════════
   ROUTER
══════════════════════════════════ */
function setActivePage(page) {
  if (!PAGES.includes(page)) page = 'dashboard';
  document.querySelectorAll('.page').forEach(p =>
    p.classList.toggle('active', p.dataset.page === page));
  document.querySelectorAll('#mainNav [data-page]').forEach(el =>
    el.classList.toggle('active', el.dataset.page === page));
  history.pushState({page}, '', '#' + page);
  localStorage.setItem(LS_PAGE, page);
  renderPage(page);
  setTimeout(animateProgressBars, 120);
}

async function renderPage(page) {
  switch(page) {
    case 'dashboard':    await renderDashboard();          break;
    case 'transactions': await renderTransactionsPage();   break;
    case 'analytics':    await renderAnalyticsPage();      break;
    case 'budget':       await renderBudgetPage();         break;
    case 'wallet':       await renderWalletPage();         break;
    case 'goals':        renderGoalsPage();                break;
  }
}

window.addEventListener('popstate', e => {
  const page = (e.state && e.state.page) || location.hash.slice(1) || 'dashboard';
  setActivePage(page);
});

/* ══════════════════════════════════
   DASHBOARD
══════════════════════════════════ */
async function renderDashboard() {
  showPageSpinner('dashStats');

  // Fetch each independently so one failure doesn't break the whole dashboard
  const [txsResult, flowResult, catResult] = await Promise.allSettled([
    getTx(),
    API.get('/api/monthly-flow?months=7'),
    API.get('/api/category-breakdown'),
  ]);

  const txs      = txsResult.status  === 'fulfilled' ? txsResult.value  : [];
  const flow     = flowResult.status === 'fulfilled' ? flowResult.value : [];
  const catBreak = catResult.status  === 'fulfilled' ? catResult.value  : [];

  if (txsResult.status === 'rejected') {
    showError('dashStats', 'transactions: ' + txsResult.reason.message);
  } else {
    renderDashStats(txs);
    renderDashRecentTx(txs);
    updateTxBadge(txs.length);
  }

  if (flowResult.status === 'rejected') {
    console.error('monthly-flow error:', flowResult.reason.message);
    const barsEl = document.getElementById('cBars');
    if (barsEl) barsEl.innerHTML = '<div style="color:var(--red-text);font-size:12px;padding:8px;">Chart error: ' + flowResult.reason.message + '</div>';
  } else {
    renderMoneyFlowChart(flow);
  }

  if (catResult.status === 'fulfilled') {
    renderDonut('donutSvg','donutTotal','donutLegend', catBreak);
  }
}

function filterTxByPeriod(txs, period) {
  const now = new Date();
  if (period === 'week') {
    const cut = new Date(now); cut.setDate(cut.getDate()-7);
    return txs.filter(t => new Date(t.date) >= cut);
  }
  if (period === 'month') {
    const cut = new Date(now); cut.setDate(cut.getDate()-30);
    return txs.filter(t => new Date(t.date) >= cut);
  }
  return txs;
}

function renderDashStats(allTxs) {
  const txs    = filterTxByPeriod(allTxs, _activePeriod);
  const income  = txs.filter(t=>t.type==='income').reduce((s,t)=>s+t.amount,0);
  const expense = txs.filter(t=>t.type==='expense').reduce((s,t)=>s+t.amount,0);
  const balance = income - expense;
  const totalInc = allTxs.filter(t=>t.type==='income').reduce((s,t)=>s+t.amount,0);
  const totalExp = allTxs.filter(t=>t.type==='expense').reduce((s,t)=>s+t.amount,0);
  const savings  = Math.max(0, totalInc - totalExp);

  const el = document.getElementById('dashStats');
  if (!el) return;
  el.innerHTML = ` + "`" + `
    ${statCard('Total balance', balance, balance>=0?'up':'dn', 'period net')}
    ${statCard('Income', income, 'up', 'period total')}
    ${statCard('Expense', expense, 'dn', 'period total')}
    ${statCard('Net savings', savings, 'up', 'accumulated')}` + "`" + `;
}

function statCard(label, val, dir, sub) {
  const arrow = dir==='up'
    ? '<polyline points="17 11 12 6 7 11"/><line x1="12" y1="6" x2="12" y2="18"/>'
    : '<polyline points="7 13 12 18 17 13"/><line x1="12" y1="18" x2="12" y2="6"/>';
  return ` + "`" + `<div class="sc">
    <div class="sc-head"><span class="sc-lbl">${esc(label)}</span>
      <span class="sc-arrow"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M7 17L17 7"/><path d="M7 7h10v10"/></svg></span>
    </div>
    <div class="sc-val">${fmtCurrency(val)}</div>
    <span class="chg ${dir}"><svg viewBox="0 0 24 24" stroke="currentColor" fill="none" stroke-width="2.5">${arrow}</svg></span>
    <span class="chg-vs">${esc(sub)}</span>
  </div>` + "`" + `;
}

function renderDashRecentTx(txs) {
  const tbody = document.getElementById('dashTxBody');
  if (!tbody) return;
  const sorted = [...txs].sort((a,b)=>new Date(b.date)-new Date(a.date)).slice(0,5);
  tbody.innerHTML = sorted.length
    ? sorted.map(t => txRowHtml(t, false)).join('')
    : ` + "`" + `<tr><td colspan="5" style="text-align:center;padding:24px;color:var(--text-2);">No transactions yet. Add one!</td></tr>` + "`" + `;
}

/* ══════════════════════════════════
   MONEY FLOW CHART
══════════════════════════════════ */
function renderMoneyFlowChart(rows) {
  const barsEl = document.getElementById('cBars');
  const xlEl   = document.getElementById('cXL');
  if (!barsEl) return;
  barsEl.innerHTML = ''; xlEl.innerHTML = '';

  const maxVal = Math.max(...rows.map(r=>Math.max(r.income,r.expense)), 1);
  const setYt = (id,v) => { const el=document.getElementById(id); if(el) el.textContent='$'+fmtK(v); };
  setYt('ytMax',maxVal); setYt('ytMid',Math.round(maxVal*.66)); setYt('ytLow',Math.round(maxVal*.33));

  rows.forEach(row => {
    const col = document.createElement('div'); col.className='bc';
    const pair= document.createElement('div'); pair.className='bp';
    const ih  = Math.max(2,(row.income /maxVal)*100);
    const eh  = Math.max(2,(row.expense/maxVal)*100);
    pair.innerHTML = ` + "`" + `
      <div class="bar inc" style="height:${ih}%"><div class="btt">${fmtCurrency(row.income)}</div></div>
      <div class="bar exp" style="height:${eh}%"><div class="btt">${fmtCurrency(row.expense)}</div></div>` + "`" + `;
    col.append(pair); barsEl.append(col);
    const lbl=document.createElement('div'); lbl.className='xl'; lbl.textContent=row.month;
    xlEl.append(lbl);
  });
}

/* ══════════════════════════════════
   DONUT CHART
══════════════════════════════════ */
function renderDonut(svgId, totalId, legendId, catRows) {
  const svgEl  = document.getElementById(svgId);
  const totEl  = document.getElementById(totalId);
  const legEl  = document.getElementById(legendId);
  if (!svgEl) return;

  const total = catRows.reduce((s,r)=>s+r.total, 0);
  if (totEl) totEl.textContent = fmtCurrency(total);

  const C = 2 * Math.PI * 14; // circumference at r=14
  svgEl.innerHTML = ` + "`" + `<circle cx="18" cy="18" r="14" fill="none" stroke="var(--border)" stroke-width="4.5"/>` + "`" + `;
  if (legEl) legEl.innerHTML = '';

  if (total === 0) {
    if (legEl) legEl.innerHTML = '<span style="font-size:12px;color:var(--text-2);">No expense data yet.</span>';
    return;
  }

  const colorsArr = Object.values(CAT_COLORS);
  let offset = 0;
  catRows.forEach((row, i) => {
    const pct  = row.total / total;
    const dash = pct * C;
    const color= CAT_COLORS[row.category] || colorsArr[i % colorsArr.length];
    const circle = document.createElementNS('http://www.w3.org/2000/svg','circle');
    circle.setAttribute('cx','18'); circle.setAttribute('cy','18'); circle.setAttribute('r','14');
    circle.setAttribute('fill','none'); circle.setAttribute('stroke',color); circle.setAttribute('stroke-width','4.5');
    circle.setAttribute('stroke-dasharray',` + "`" + `${dash.toFixed(2)} ${(C-dash).toFixed(2)}` + "`" + `);
    circle.setAttribute('stroke-dashoffset', (-offset).toFixed(2));
    svgEl.appendChild(circle);
    offset += dash;
    if (legEl) {
      const row2=document.createElement('div'); row2.className='bli';
      row2.innerHTML=` + "`" + `<div class="bldot" style="background:${color}"></div><span class="blname">${esc(row.category)}</span>` + "`" + `;
      legEl.appendChild(row2);
    }
  });
}

/* ══════════════════════════════════
   TRANSACTIONS PAGE
══════════════════════════════════ */
async function renderTransactionsPage() {
  invalidateCache(); // always fresh on page visit
  await applyTxFilters();
}

async function applyTxFilters() {
  const search   = (document.getElementById('txSearch')?.value||'').toLowerCase();
  const type     = document.getElementById('txTypeFilter')?.value||'';
  const dateFrom = document.getElementById('txDateFrom')?.value||'';
  const dateTo   = document.getElementById('txDateTo')?.value||'';

  let txs = await getTx();
  txs = [...txs].sort((a,b)=>new Date(b.date)-new Date(a.date));

  if (search)   txs = txs.filter(t=>(t.note||'').toLowerCase().includes(search)||(t.category||'').toLowerCase().includes(search));
  if (type)     txs = txs.filter(t=>t.type===type);
  if (dateFrom) txs = txs.filter(t=>t.date>=dateFrom);
  if (dateTo)   txs = txs.filter(t=>t.date<=dateTo);

  const tbody = document.getElementById('txPageBody');
  const empty = document.getElementById('txEmpty');
  if (!tbody) return;

  if (txs.length === 0) {
    tbody.innerHTML = '';
    if (empty) empty.style.display='block';
  } else {
    if (empty) empty.style.display='none';
    tbody.innerHTML = txs.map(t => txRowHtml(t, true)).join('');
    tbody.querySelectorAll('.del-btn').forEach(btn =>
      btn.addEventListener('click', () => deleteTransaction(btn.dataset.id)));
  }
  updateTxBadge(txs.length);
}

async function deleteTransaction(id) {
  if (!confirm('Delete this transaction?')) return;
  try {
    await API.del(` + "`" + `/api/transactions/${id}` + "`" + `);
    invalidateCache();
    await applyTxFilters();
    await renderDashboard();
    showToast('Transaction deleted.');
  } catch(e) {
    showToast('Delete failed: ' + e.message);
  }
}

/* ══════════════════════════════════
   ANALYTICS PAGE
══════════════════════════════════ */
let _anPeriod = 'month';

async function renderAnalyticsPage() {
  try {
    const [txs, catBreak, flow] = await Promise.all([
      getTx(),
      API.get('/api/category-breakdown'),
      API.get('/api/monthly-flow?months=7'),
    ]);

    // Filter for period
    const now = new Date();
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1);
    const lastMonthStart = new Date(now.getFullYear(), now.getMonth()-1, 1);
    const lastMonthEnd   = new Date(now.getFullYear(), now.getMonth(), 0);

    const periodTxs = _anPeriod === 'month'
      ? txs.filter(t => new Date(t.date) >= monthStart)
      : txs;
    const lastTxs = txs.filter(t => {
      const d = new Date(t.date);
      return d >= lastMonthStart && d <= lastMonthEnd;
    });

    const inc    = periodTxs.filter(t=>t.type==='income').reduce((s,t)=>s+t.amount,0);
    const exp    = periodTxs.filter(t=>t.type==='expense').reduce((s,t)=>s+t.amount,0);
    const bal    = inc - exp;
    const lastInc= lastTxs.filter(t=>t.type==='income').reduce((s,t)=>s+t.amount,0);
    const lastExp= lastTxs.filter(t=>t.type==='expense').reduce((s,t)=>s+t.amount,0);
    const lastBal= lastInc - lastExp;

    const incTxCount = periodTxs.filter(t=>t.type==='income').length;
    const expTxCount = periodTxs.filter(t=>t.type==='expense').length;
    const incCats = [...new Set(periodTxs.filter(t=>t.type==='income').map(t=>t.category))].length;
    const expCats = [...new Set(periodTxs.filter(t=>t.type==='expense').map(t=>t.category))].length;

    // Helper: format big value split at decimal
    function setVal(id, val) {
      const el = document.getElementById(id);
      if (!el) return;
      const parts = Math.abs(val).toFixed(2).split('.');
      el.innerHTML = (val < 0 ? '-' : '') + '$' + Number(parts[0]).toLocaleString() + '<span class="an-cents">.' + parts[1] + '</span>';
    }

    function setPct(id, cur, prev, cls) {
      const el = document.getElementById(id);
      if (!el) return;
      if (prev === 0) { el.textContent = ''; return; }
      const pct = ((cur - prev) / Math.abs(prev) * 100);
      const arrow = pct >= 0 ? '↑' : '↓';
      el.textContent = arrow + ' ' + Math.abs(pct).toFixed(1) + '%';
      el.className = 'an-change ' + (cls === 'exp' ? (pct > 0 ? 'neg' : 'pos') : (pct >= 0 ? 'pos' : 'neg'));
    }

    function setT(id, v) { const el=document.getElementById(id); if(el) el.textContent=v; }

    setVal('anTotalBal', bal);
    setVal('anTotalInc', inc);
    setVal('anTotalExp', exp);

    setPct('anBalChange', bal, lastBal, 'bal');
    setPct('anIncChange', inc, lastInc, 'inc');
    setPct('anExpChange', exp, lastExp, 'exp');

    setT('anBalTx',  periodTxs.length + ' transactions');
    setT('anIncTx',  incTxCount + ' transactions');
    setT('anExpTx',  expTxCount + ' transactions');

    setT('anBalCats', (incCats + expCats) + ' categories');
    setT('anIncCats', incCats + ' categories');
    setT('anExpCats', expCats + ' categories');

    const extra = bal - lastBal;
    setT('anBalSub', extra >= 0
      ? 'You have extra ' + fmtCurrency(Math.abs(extra)) + ' vs last month'
      : 'You have ' + fmtCurrency(Math.abs(extra)) + ' less vs last month');
    const extraInc = inc - lastInc;
    setT('anIncSub', extraInc >= 0
      ? 'You earn extra ' + fmtCurrency(Math.abs(extraInc)) + ' vs last month'
      : 'You earn ' + fmtCurrency(Math.abs(extraInc)) + ' less vs last month');
    const extraExp = exp - lastExp;
    setT('anExpSub', extraExp > 0
      ? 'You spent extra ' + fmtCurrency(Math.abs(extraExp)) + ' vs last month'
      : 'You spent ' + fmtCurrency(Math.abs(extraExp)) + ' less vs last month');

    // Line chart (canvas-based smooth lines)
    renderAnLineChart(flow);

    // Donut statistics
    const monthCat = catBreak; // already filtered server-side to expense
    renderAnDonut(monthCat, exp);

    // Bar chart
    renderAnBarChart(flow);

  } catch(e) { console.error('analytics:', e); }
}

function renderAnLineChart(flow) {
  const canvas = document.getElementById('anLineChart');
  const xlEl   = document.getElementById('anLineXL');
  if (!canvas) return;

  const W = canvas.parentElement.offsetWidth || 600;
  const H = 200;
  canvas.width  = W;
  canvas.height = H;
  const ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, W, H);

  if (!flow.length) {
    if (xlEl) xlEl.innerHTML = '';
    return;
  }

  const vals = flow.map(r => r.income - r.expense);
  const maxV = Math.max(...vals.map(Math.abs), 1);
  const pad  = { t: 20, b: 10, l: 10, r: 10 };
  const cW   = W - pad.l - pad.r;
  const cH   = H - pad.t - pad.b;

  const toX = i => pad.l + (i / Math.max(flow.length-1,1)) * cW;
  const toY = v => pad.t + cH - ((v + maxV) / (2 * maxV)) * cH;

  // Grid lines
  ctx.strokeStyle = getComputedStyle(document.documentElement).getPropertyValue('--grid-line').trim() || '#ede9ff';
  ctx.lineWidth = 1;
  [0.25, 0.5, 0.75, 1].forEach(t => {
    const y = pad.t + t * cH;
    ctx.beginPath(); ctx.moveTo(pad.l, y); ctx.lineTo(W - pad.r, y); ctx.stroke();
  });

  // Fill gradient
  const accent = '#7c5cfc';
  const grad = ctx.createLinearGradient(0, pad.t, 0, H);
  grad.addColorStop(0, 'rgba(124,92,252,0.25)');
  grad.addColorStop(1, 'rgba(124,92,252,0)');

  ctx.beginPath();
  vals.forEach((v, i) => {
    const x = toX(i), y = toY(v);
    i === 0 ? ctx.moveTo(x, y) : ctx.bezierCurveTo(toX(i-0.5), toY(vals[i-1]), toX(i-0.5), y, x, y);
  });
  ctx.lineTo(toX(vals.length-1), H);
  ctx.lineTo(toX(0), H);
  ctx.closePath();
  ctx.fillStyle = grad;
  ctx.fill();

  // Line
  ctx.beginPath();
  ctx.strokeStyle = accent;
  ctx.lineWidth = 2.5;
  ctx.lineJoin = 'round';
  vals.forEach((v, i) => {
    const x = toX(i), y = toY(v);
    i === 0 ? ctx.moveTo(x, y) : ctx.bezierCurveTo(toX(i-0.5), toY(vals[i-1]), toX(i-0.5), y, x, y);
  });
  ctx.stroke();

  // Dots + tooltip on highest value
  const maxIdx = vals.indexOf(Math.max(...vals));
  vals.forEach((v, i) => {
    ctx.beginPath();
    ctx.arc(toX(i), toY(v), i === maxIdx ? 5 : 3.5, 0, Math.PI*2);
    ctx.fillStyle = accent;
    ctx.fill();
    ctx.strokeStyle = '#fff';
    ctx.lineWidth = 2;
    ctx.stroke();

    if (i === maxIdx) {
      const x = toX(i), y = toY(v);
      const label = fmtCurrency(v);
      const bw = ctx.measureText(label).width + 16;
      const bh = 24;
      const bx = Math.min(x - bw/2, W - bw - 4);
      const by = y - bh - 10;
      ctx.fillStyle = '#fff';
      ctx.shadowColor = 'rgba(0,0,0,0.12)';
      ctx.shadowBlur = 8;
      ctx.beginPath();
      ctx.roundRect(bx, by, bw, bh, 6);
      ctx.fill();
      ctx.shadowBlur = 0;
      ctx.fillStyle = '#1a1340';
      ctx.font = '600 12px Sora, sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(label, bx + bw/2, by + 16);
    }
  });

  // X labels
  if (xlEl) {
    xlEl.innerHTML = flow.map(r => ` + "`" + `<span>${esc(r.month)}</span>` + "`" + `).join('');
  }
}

function renderAnDonut(catRows, total) {
  const svgEl  = document.getElementById('anDonutSvg');
  const totEl  = document.getElementById('anDonutTotal');
  const legEl  = document.getElementById('anDonutLegend');
  if (!svgEl) return;

  const parts = total.toFixed(2).split('.');
  if (totEl) totEl.innerHTML = '$' + Number(parts[0]).toLocaleString() + '<span style="font-size:11px;color:var(--text-2)">.' + parts[1] + '</span>';

  const C = 2 * Math.PI * 13;
  // Keep base circle
  svgEl.innerHTML = '<circle cx="18" cy="18" r="13" fill="none" stroke="var(--border)" stroke-width="5"/>';
  if (legEl) legEl.innerHTML = '';

  if (!total || !catRows.length) return;

  const colorsArr = Object.values(CAT_COLORS);
  let offset = 0;
  catRows.slice(0,6).forEach((row, i) => {
    const pct  = row.total / total;
    const dash = pct * C;
    const gap  = C * 0.02; // small gap between segments
    const color = CAT_COLORS[row.category] || colorsArr[i % colorsArr.length];
    const circle = document.createElementNS('http://www.w3.org/2000/svg','circle');
    circle.setAttribute('cx','18'); circle.setAttribute('cy','18'); circle.setAttribute('r','13');
    circle.setAttribute('fill','none'); circle.setAttribute('stroke',color); circle.setAttribute('stroke-width','5');
    circle.setAttribute('stroke-dasharray',` + "`" + `${Math.max(0,dash-gap).toFixed(2)} ${(C-dash+gap).toFixed(2)}` + "`" + `);
    circle.setAttribute('stroke-dashoffset', (-offset).toFixed(2));
    svgEl.appendChild(circle);
    offset += dash;

    if (legEl) {
      const pctLbl = (pct*100).toFixed(0)+'%';
      const item = document.createElement('div');
      item.style.cssText = 'display:flex;align-items:center;gap:5px;font-size:11px;color:var(--text-2);';
      item.innerHTML = ` + "`" + `<span style="width:8px;height:8px;border-radius:50%;background:${color};flex-shrink:0;"></span>${esc(row.category)} <span style="color:var(--text-1);font-weight:600;">${pctLbl}</span>` + "`" + `;
      legEl.appendChild(item);
    }
  });
}

function renderAnBarChart(flow) {
  const barsEl = document.getElementById('anBars');
  const xlEl   = document.getElementById('anXL');
  if (!barsEl) return;
  barsEl.innerHTML = ''; xlEl.innerHTML = '';

  const mMax = Math.max(...flow.map(r=>Math.max(r.income,r.expense)), 1);
  const setYt = (id,v) => { const el=document.getElementById(id); if(el) el.textContent='$'+fmtK(v); };
  setYt('anYtMax',mMax); setYt('anYtMid',Math.round(mMax*.66)); setYt('anYtLow',Math.round(mMax*.33));

  flow.forEach(row => {
    const col  = document.createElement('div'); col.className='bc';
    const pair = document.createElement('div'); pair.className='bp';
    const ih   = Math.max(2,(row.income /mMax)*100);
    const eh   = Math.max(2,(row.expense/mMax)*100);
    pair.innerHTML = ` + "`" + `
      <div class="bar inc" style="height:${ih}%;border-radius:6px 6px 0 0;"><div class="btt">${fmtCurrency(row.income)}</div></div>
      <div class="bar exp" style="height:${eh}%;border-radius:6px 6px 0 0;"><div class="btt">${fmtCurrency(row.expense)}</div></div>` + "`" + `;
    col.append(pair); barsEl.append(col);
    const lbl = document.createElement('div'); lbl.className='xl'; lbl.textContent=row.month;
    xlEl.append(lbl);
  });
}

/* ══════════════════════════════════
   BUDGET PAGE
══════════════════════════════════ */
async function renderBudgetPage() {
  try {
    const catBreak = await API.get('/api/category-breakdown');
    const total    = catBreak.reduce((s,r)=>s+r.total,0);
    const listEl   = document.getElementById('budgetCatList');
    const colorsArr= Object.values(CAT_COLORS);
    if (listEl) {
      listEl.innerHTML = catBreak.length
        ? catBreak.map((r,i)=>{
            const c=CAT_COLORS[r.category]||colorsArr[i%colorsArr.length];
            const pct=total?(r.total/total*100).toFixed(1):0;
            return ` + "`" + `<div class="budget-cat-row">
              <div class="budget-cat-dot" style="background:${c}"></div>
              <div class="budget-cat-name">${esc(r.category)}</div>
              <div class="budget-cat-bar-wrap"><div class="budget-cat-bar" style="width:${pct}%;background:${c};"></div></div>
              <div class="budget-cat-amt">${fmtCurrency(r.total)}</div>
            </div>` + "`" + `;
          }).join('')
        : '<div style="color:var(--text-2);font-size:13px;padding:12px 0;">No expense data yet.</div>';
    }
    renderDonut('budgetDonutSvg','budgetDonutTotal','budgetDonutLegend', catBreak);
  } catch(e) { console.error('budget:', e); }
}

/* ══════════════════════════════════
   WALLET PAGE
══════════════════════════════════ */
async function renderWalletPage() {
  try {
    const txs = await getTx();
    const inc = txs.filter(t=>t.type==='income').reduce((s,t)=>s+t.amount,0);
    const exp = txs.filter(t=>t.type==='expense').reduce((s,t)=>s+t.amount,0);
    const net = inc - exp;

    const el = document.getElementById('walletMainBal');
    if (el) el.textContent = fmtCurrency(Math.max(0, net));

    // Update wallet card totals dynamically
    const incEl = document.getElementById('walletTotalInc');
    const expEl = document.getElementById('walletTotalExp');
    if (incEl) incEl.textContent = fmtCurrency(inc);
    if (expEl) expEl.textContent = fmtCurrency(exp);

    const tbody = document.getElementById('walletTxBody');
    if (tbody) {
      const sorted = [...txs].sort((a,b)=>new Date(b.date)-new Date(a.date)).slice(0,8);
      tbody.innerHTML = sorted.length
        ? sorted.map(t=>txRowHtml(t,false)).join('')
        : '<tr><td colspan="5" style="text-align:center;padding:24px;color:var(--text-2);">No transactions yet.</td></tr>';
    }
  } catch(e) { console.error('wallet:', e); }
}

/* ══════════════════════════════════
   GOALS — localStorage based
══════════════════════════════════ */
const LS_GOALS = 'finset-goals';

function loadGoals() {
  try { return JSON.parse(localStorage.getItem(LS_GOALS)||'[]'); } catch{ return []; }
}
function saveGoals(goals) {
  localStorage.setItem(LS_GOALS, JSON.stringify(goals));
}

function renderGoalsPage() {
  const goals = loadGoals();
  const grid  = document.getElementById('goalsGrid');
  if (!grid) return;

  const goalCards = goals.map(g => {
    const pct     = Math.min(100, g.target > 0 ? (g.saved / g.target * 100) : 0);
    const pctFmt  = pct.toFixed(1);
    const color   = pct >= 100 ? '#22c55e' : pct >= 60 ? '#7c5cfc' : pct >= 30 ? '#f59e0b' : '#e879f9';
    return ` + "`" + `<div class="goal-card" data-id="${g.id}">
      <div style="display:flex;justify-content:space-between;align-items:flex-start;">
        <span class="goal-icon">${esc(g.icon||'🎯')}</span>
        <div style="display:flex;gap:6px;">
          <button onclick="openEditGoal('${g.id}')" style="background:none;border:none;cursor:pointer;color:var(--text-2);padding:2px;" title="Edit">
            <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          </button>
          <button onclick="deleteGoal('${g.id}')" style="background:none;border:none;cursor:pointer;color:var(--red-text);padding:2px;" title="Delete">
            <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/></svg>
          </button>
        </div>
      </div>
      <div class="goal-name">${esc(g.name)}</div>
      <div class="goal-target">Target: ${fmtCurrency(g.target)}</div>
      <div class="prog-t" style="margin:10px 0 6px;">
        <div class="prog-f" style="width:${pctFmt}%;background:${color};transition:width .8s cubic-bezier(.4,0,.2,1);"></div>
      </div>
      <div class="goal-pct">${pctFmt}% — <span class="goal-amount">${fmtCurrency(g.saved)} saved</span></div>
      <div style="margin-top:8px;display:flex;gap:6px;">
        <input type="number" placeholder="Add amount" min="0" step="0.01"
          style="flex:1;padding:5px 8px;border:1px solid var(--border);border-radius:8px;background:var(--input-bg);color:var(--text-1);font-size:12px;"
          id="goalAdd_${g.id}">
        <button onclick="addToGoal('${g.id}')"
          style="padding:5px 10px;background:var(--accent);color:#fff;border:none;border-radius:8px;font-size:12px;cursor:pointer;">
          + Add
        </button>
      </div>
    </div>` + "`" + `;
  }).join('');

  const addCard = ` + "`" + `<div class="goal-card" id="addGoalCard"
    style="display:flex;flex-direction:column;align-items:center;justify-content:center;cursor:pointer;border:2px dashed var(--border);background:transparent;min-height:180px;"
    onclick="openAddGoal()">
    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="var(--text-3)" stroke-width="1.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
    <div style="font-size:13px;color:var(--text-2);margin-top:8px;font-weight:500;">Add new goal</div>
  </div>` + "`" + `;

  grid.innerHTML = goalCards + addCard;
  updateGoalsBadge(goals.length);
}

function updateGoalsBadge(n) {
  const el = document.getElementById('goalsBadge');
  if (el) el.textContent = n;
}

function openAddGoal() {
  document.getElementById('goalModalTitle').textContent = 'Add goal';
  document.getElementById('gIcon').value   = '🎯';
  document.getElementById('gName').value   = '';
  document.getElementById('gTarget').value = '';
  document.getElementById('gSaved').value  = '0';
  document.getElementById('gEditId').value = '';
  document.getElementById('gNameErr').classList.remove('show');
  document.getElementById('gTargetErr').classList.remove('show');
  document.getElementById('goalModal').classList.add('open');
}

function openEditGoal(id) {
  const goals = loadGoals();
  const g = goals.find(x=>x.id===id);
  if (!g) return;
  document.getElementById('goalModalTitle').textContent = 'Edit goal';
  document.getElementById('gIcon').value   = g.icon  || '🎯';
  document.getElementById('gName').value   = g.name  || '';
  document.getElementById('gTarget').value = g.target|| '';
  document.getElementById('gSaved').value  = g.saved || '0';
  document.getElementById('gEditId').value = id;
  document.getElementById('gNameErr').classList.remove('show');
  document.getElementById('gTargetErr').classList.remove('show');
  document.getElementById('goalModal').classList.add('open');
}

function closeGoalModal() {
  document.getElementById('goalModal').classList.remove('open');
}

function saveGoalModal() {
  const name   = document.getElementById('gName').value.trim();
  const target = parseFloat(document.getElementById('gTarget').value);
  const saved  = parseFloat(document.getElementById('gSaved').value)||0;
  const icon   = document.getElementById('gIcon').value.trim()||'🎯';
  const editId = document.getElementById('gEditId').value;

  let valid = true;
  if (!name) { document.getElementById('gNameErr').classList.add('show'); valid=false; }
  else document.getElementById('gNameErr').classList.remove('show');
  if (!target || target <= 0) { document.getElementById('gTargetErr').classList.add('show'); valid=false; }
  else document.getElementById('gTargetErr').classList.remove('show');
  if (!valid) return;

  const goals = loadGoals();
  if (editId) {
    const idx = goals.findIndex(x=>x.id===editId);
    if (idx>=0) goals[idx] = {...goals[idx], icon, name, target, saved};
  } else {
    goals.push({id: crypto.randomUUID(), icon, name, target, saved});
  }
  saveGoals(goals);
  closeGoalModal();
  renderGoalsPage();
  showToast(editId ? 'Goal updated!' : 'Goal added!');
}

function deleteGoal(id) {
  if (!confirm('Delete this goal?')) return;
  const goals = loadGoals().filter(g=>g.id!==id);
  saveGoals(goals);
  renderGoalsPage();
  showToast('Goal deleted.');
}

function addToGoal(id) {
  const input = document.getElementById('goalAdd_'+id);
  const amt   = parseFloat(input?.value);
  if (!amt || amt <= 0) { showToast('Enter a valid amount.'); return; }
  const goals = loadGoals();
  const idx   = goals.findIndex(g=>g.id===id);
  if (idx < 0) return;
  goals[idx].saved = (goals[idx].saved||0) + amt;
  saveGoals(goals);
  renderGoalsPage();
  showToast('+' + fmtCurrency(amt) + ' added to goal!');
}

/* ══════════════════════════════════
   TX ROW HTML HELPER
══════════════════════════════════ */
function txRowHtml(t, showDelete) {
  const isInc = t.type==='income';
  const iconBg= isInc?'var(--green-bg)':'var(--red-bg)';
  const iconSt= isInc?'var(--green-text)':'var(--red-text)';
  const iconSvg= isInc
    ? ` + "`" + `<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="${iconSt}" stroke-width="2.2"><polyline points="17 11 12 6 7 11"/><line x1="12" y1="6" x2="12" y2="18"/></svg>` + "`" + `
    : ` + "`" + `<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="${iconSt}" stroke-width="2.2"><polyline points="7 13 12 18 17 13"/><line x1="12" y1="18" x2="12" y2="6"/></svg>` + "`" + `;
  const delBtn = showDelete
    ? ` + "`" + `<td><button class="del-btn" data-id="${t.id}" title="Delete">
        <svg viewBox="0 0 24 24"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4h6v2"/></svg>
      </button></td>` + "`" + ` : '';
  const typeTd = showDelete
    ? ` + "`" + `<td><span class="cchip ${isInc?'g':'r'}">${isInc?'Income':'Expense'}</span></td>` + "`" + ` : '';
  return ` + "`" + `<tr>
    <td class="td-d" style="padding-left:20px;">${fmtDate(t.date)}</td>
    <td class="${isInc?'pos':'neg'}">${isInc?'+':'−'}${fmtCurrency(t.amount)}</td>
    <td><span class="txi" style="background:${iconBg};">${iconSvg}</span><span class="txn">${esc(t.note||t.category)}</span></td>
    <td><span class="mchip"><svg viewBox="0 0 24 24" width="11" height="11" stroke="currentColor" fill="none" stroke-width="1.8"><rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/></svg>${esc(t.method)}</span></td>
    <td><span class="cchip ${isInc?'g':''}">${esc(t.category)}</span></td>
    ${typeTd}${delBtn}
  </tr>` + "`" + `;
}

/* ══════════════════════════════════
   MODAL SYSTEM
══════════════════════════════════ */
function openModal() {
  _selectedType = 'income';
  document.getElementById('typeBtnIncome').className ='type-btn sel-income';
  document.getElementById('typeBtnExpense').className='type-btn';
  document.getElementById('fAmount').value='';
  document.getElementById('fAmount').classList.remove('err');
  document.getElementById('fAmountErr').classList.remove('show');
  document.getElementById('fCategory').value='';
  document.getElementById('fCategory').classList.remove('err');
  document.getElementById('fCategoryErr').classList.remove('show');
  document.getElementById('fCustomCat').value='';
  document.getElementById('fCustomCatRow').style.display='none';
  document.getElementById('fMethod').value='Cash';
  document.getElementById('fDate').value=new Date().toISOString().split('T')[0];
  document.getElementById('fNote').value='';
  document.getElementById('txModal').classList.add('open');
}
function closeModal() {
  document.getElementById('txModal').classList.remove('open');
}

async function submitTransaction() {
  let valid = true;
  const amtRaw = parseFloat(document.getElementById('fAmount').value);
  if (isNaN(amtRaw)||amtRaw<=0) {
    document.getElementById('fAmount').classList.add('err');
    document.getElementById('fAmountErr').classList.add('show');
    valid=false;
  } else {
    document.getElementById('fAmount').classList.remove('err');
    document.getElementById('fAmountErr').classList.remove('show');
  }
  let cat = document.getElementById('fCategory').value;
  if (!cat) {
    document.getElementById('fCategory').classList.add('err');
    document.getElementById('fCategoryErr').classList.add('show');
    valid=false;
  } else {
    document.getElementById('fCategory').classList.remove('err');
    document.getElementById('fCategoryErr').classList.remove('show');
    if (cat==='Other') { const c=document.getElementById('fCustomCat').value.trim(); if(c) cat=c; }
  }
  if (!valid) return;

  const btn = document.getElementById('modalSubmitBtn');
  btn.textContent='Saving…'; btn.disabled=true;
  try {
    await API.post('/api/transactions', {
      type: _selectedType,
      amount: amtRaw,
      category: cat,
      method: document.getElementById('fMethod').value,
      date:   document.getElementById('fDate').value,
      note:   document.getElementById('fNote').value.trim(),
    });
    invalidateCache();
    closeModal();
    await renderDashboard();
    const txPage=document.querySelector('.page[data-page="transactions"]');
    if (txPage&&txPage.classList.contains('active')) await renderTransactionsPage();
    showToast(` + "`" + `${_selectedType==='income'?'Income':'Expense'} saved! 🎉` + "`" + `);
  } catch(e) {
    showToast('Error: '+e.message);
  } finally {
    btn.textContent='Save transaction'; btn.disabled=false;
  }
}

/* ══════════════════════════════════
   SETTINGS: EXPORT / IMPORT / RESET
══════════════════════════════════ */
async function exportData() {
  const txs  = await getTx();
  const blob = new Blob([JSON.stringify({transactions:txs, exportedAt:new Date().toISOString()},null,2)],{type:'application/json'});
  const a    = document.createElement('a');
  a.href     = URL.createObjectURL(blob); a.download='finset-transactions.json'; a.click();
  showToast('Export downloaded!');
}

async function importData(file) {
  if (!file) return;
  const reader = new FileReader();
  reader.onload = async e => {
    try {
      const obj = JSON.parse(e.target.result);
      const txs = obj.transactions||(Array.isArray(obj)?obj:null);
      if (!txs||!Array.isArray(txs)) throw new Error('Expected {transactions:[…]}');
      const res = await API.post('/api/import', {transactions:txs});
      invalidateCache();
      await renderDashboard();
      showToast(` + "`" + `Imported ${res.inserted} new, skipped ${res.skipped} duplicates.` + "`" + `);
    } catch(err) { showToast('Import failed: '+err.message); }
  };
  reader.readAsText(file);
}

/* ══════════════════════════════════
   THEME
══════════════════════════════════ */
function applyTheme(t) {
  document.documentElement.setAttribute('data-theme',t);
  const l=t==='dark'?'Dark mode':'Light mode';
  ['themeLabel','themeLabel2'].forEach(id=>{const el=document.getElementById(id);if(el)el.textContent=l;});
  localStorage.setItem(LS_THEME,t);
}
function toggleTheme() {
  applyTheme(document.documentElement.getAttribute('data-theme')==='light'?'dark':'light');
}

/* ══════════════════════════════════
   SIDEBAR COLLAPSE
══════════════════════════════════ */
function setSidebarCollapsed(c) {
  document.getElementById('sidebar').classList.toggle('collapsed',c);
  localStorage.setItem(LS_SIDEBAR,c?'1':'0');
}

/* ══════════════════════════════════
   TOAST
══════════════════════════════════ */
function showToast(msg,dur=3000) {
  const el=document.createElement('div'); el.className='toast'; el.textContent=msg;
  document.getElementById('toastContainer').appendChild(el);
  setTimeout(()=>{el.classList.add('out');setTimeout(()=>el.remove(),300);},dur);
}

/* ══════════════════════════════════
   HELPERS
══════════════════════════════════ */
function fmtCurrency(n) {
  return '$'+Math.abs(n).toLocaleString('en',{minimumFractionDigits:2,maximumFractionDigits:2});
}
function fmtK(n){ return n>=1000?(n/1000).toFixed(0)+'k':Math.round(n).toString(); }
function fmtDate(s){
  if(!s) return '';
  const d=new Date(s+'T12:00:00');
  return d.toLocaleDateString('en',{day:'2-digit',month:'short',year:'numeric'});
}
function esc(s){
  return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}
function animateProgressBars(){
  document.querySelectorAll('.prog-f[data-w]').forEach(b=>{
    if(b.style.width==='0%'||!b.style.width) requestAnimationFrame(()=>{b.style.width=b.dataset.w+'%';});
  });
}
function updateTxBadge(n){
  const el=document.getElementById('txBadge'); if(el) el.textContent=n;
}
function showPageSpinner(containerId) {
  const el=document.getElementById(containerId);
  if(el) el.innerHTML=` + "`" + `<div style="grid-column:1/-1;text-align:center;padding:20px;color:var(--text-2);font-size:12px;">Loading…</div>` + "`" + `;
}
function showError(containerId, msg) {
  const el=document.getElementById(containerId);
  if(el) el.innerHTML=` + "`" + `<div style="grid-column:1/-1;text-align:center;padding:20px;color:var(--red-text);font-size:12px;">⚠ ${esc(msg)}</div>` + "`" + `;
}

/* ══════════════════════════════════
   BOOT — bind all events once
══════════════════════════════════ */
document.addEventListener('DOMContentLoaded', () => {

  // Sidebar collapse
  const sidebar     = document.getElementById('sidebar');
  const collapseBtn = document.getElementById('collapseBtn');
  if (localStorage.getItem(LS_SIDEBAR)==='1') sidebar.classList.add('collapsed');
  collapseBtn.addEventListener('click', ()=>setSidebarCollapsed(!sidebar.classList.contains('collapsed')));

  // Nav routing
  document.getElementById('mainNav').addEventListener('click', e=>{
    const item=e.target.closest('[data-page]');
    if(item) setActivePage(item.dataset.page);
  });

  // Period pills
  document.querySelectorAll('.pill[data-period]').forEach(p=>p.addEventListener('click',()=>{
    document.querySelectorAll('.pill[data-period]').forEach(x=>x.classList.remove('active'));
    p.classList.add('active');
    _activePeriod=p.dataset.period;
    getTx().then(txs=>renderDashStats(txs));
  }));

  // Goal modal
  document.getElementById('addGoalBtn').addEventListener('click', openAddGoal);
  document.getElementById('goalModalClose').addEventListener('click', closeGoalModal);
  document.getElementById('goalModalCancel').addEventListener('click', closeGoalModal);
  document.getElementById('goalModalSave').addEventListener('click', saveGoalModal);
  document.getElementById('goalModal').addEventListener('click', e=>{
    if(e.target===document.getElementById('goalModal')) closeGoalModal();
  });

  // Init goals badge
  updateGoalsBadge(loadGoals().length);

  // Add transaction buttons
  document.getElementById('addTxBtnDash').addEventListener('click', openModal);
  document.getElementById('addTxBtnTx').addEventListener('click', openModal);

  // See all
  document.getElementById('seeAllBtn').addEventListener('click', ()=>setActivePage('transactions'));

  // Modal
  document.getElementById('modalClose').addEventListener('click', closeModal);
  document.getElementById('modalCancelBtn').addEventListener('click', closeModal);
  document.getElementById('txModal').addEventListener('click', e=>{
    if(e.target===document.getElementById('txModal')) closeModal();
  });
  document.getElementById('modalSubmitBtn').addEventListener('click', submitTransaction);
  document.addEventListener('keydown', e=>{ if(e.key==='Escape') closeModal(); });

  // Type toggle
  document.getElementById('typeBtnIncome').addEventListener('click',()=>{
    _selectedType='income';
    document.getElementById('typeBtnIncome').className ='type-btn sel-income';
    document.getElementById('typeBtnExpense').className='type-btn';
  });
  document.getElementById('typeBtnExpense').addEventListener('click',()=>{
    _selectedType='expense';
    document.getElementById('typeBtnIncome').className ='type-btn';
    document.getElementById('typeBtnExpense').className='type-btn sel-expense';
  });

  // Custom category
  document.getElementById('fCategory').addEventListener('change', e=>{
    document.getElementById('fCustomCatRow').style.display=e.target.value==='Other'?'block':'none';
  });

  // Tx page filters
  ['txSearch','txTypeFilter','txDateFrom','txDateTo'].forEach(id=>{
    const el=document.getElementById(id);
    if(el){
      el.addEventListener('input',  applyTxFilters);
      el.addEventListener('change', applyTxFilters);
    }
  });

  // Analytics period toggle
  document.querySelectorAll('[data-anperiod]').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('[data-anperiod]').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      _anPeriod = btn.dataset.anperiod;
      const anPage = document.querySelector('.page[data-page="analytics"]');
      if (anPage && anPage.classList.contains('active')) renderAnalyticsPage();
    });
  });

  // Analytics CSV export
  const csvBtn = document.getElementById('anExportCsv');
  if (csvBtn) csvBtn.addEventListener('click', async () => {
    const txs = await getTx();
    const header = 'id,type,amount,category,method,date,note';
    const rows = txs.map(t =>
      [t.id,t.type,t.amount,t.category,t.method,t.date,(t.note||'').replace(/,/g,';')].join(',')
    );
    const blob = new Blob([header+'\n'+rows.join('\n')], {type:'text/csv'});
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob); a.download='finset-transactions.csv'; a.click();
    showToast('CSV exported!');
  });

  // Theme
  document.getElementById('themeToggle').addEventListener('click',  toggleTheme);
  document.getElementById('themeToggle2').addEventListener('click', toggleTheme);

  // Settings
  document.getElementById('exportBtn').addEventListener('click', exportData);
  document.getElementById('importFile').addEventListener('change', e=>importData(e.target.files[0]));
  document.getElementById('resetDemoBtn').addEventListener('click', async ()=>{
    if(!confirm('This will wipe all data and reload demo transactions. Continue?')) return;
    try {
      // Clear all transactions via dedicated endpoint
      await API.del('/api/transactions/all');
      // Re-import demo data (UUIDs are now proper RFC-4122 format)
      const demo = getDemoTransactions();
      await API.post('/api/import', {transactions:demo});
      invalidateCache();
      await renderDashboard();
      showToast('Demo data loaded!');
    } catch(e) { showToast('Reset failed: '+e.message); }
  });

  // Apply saved theme
  applyTheme(localStorage.getItem(LS_THEME)||'light');

  // Init page from hash or storage
  const hashPage   = location.hash.slice(1);
  const storedPage = localStorage.getItem(LS_PAGE)||'dashboard';
  const initPage   = PAGES.includes(hashPage)?hashPage:storedPage;
  history.replaceState({page:initPage},'','#'+initPage);
  setActivePage(initPage);
});

/* ══════════════════════════════════
   DEMO DATA (used by "Reset demo" button)
══════════════════════════════════ */
function getDemoTransactions() {
  const now = new Date();
  const d = offset => { const dt=new Date(now); dt.setDate(dt.getDate()-offset); return dt.toISOString().split('T')[0]; };
  // Generate RFC-4122 compliant UUIDs (required by PostgreSQL uuid column type)
  const uid = () => ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16));
  return [
    {id:uid(),type:'income', amount:3200, category:'Salary',       method:'Transfer',date:d(1), note:'Monthly salary',       created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:80,   category:'Food',         method:'Card',    date:d(2), note:'Yaposhka restaurant',  created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:150,  category:'Shopping',     method:'Card',    date:d(3), note:'Reserved item',        created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:13,   category:'Entertainment',method:'Card',    date:d(3), note:'YouTube Premium',      created_at:new Date().toISOString()},
    {id:uid(),type:'income', amount:800,  category:'Freelance',    method:'Transfer',date:d(5), note:'Design project',       created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:120,  category:'Bills',        method:'Transfer',date:d(6), note:'Electricity',          created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:55,   category:'Transport',    method:'Cash',    date:d(8), note:'Taxi',                 created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:200,  category:'Health',       method:'Card',    date:d(10),note:'Pharmacy',             created_at:new Date().toISOString()},
    {id:uid(),type:'income', amount:500,  category:'Investment',   method:'Transfer',date:d(12),note:'Dividends',            created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:90,   category:'Food',         method:'Cash',    date:d(14),note:'Groceries',            created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:35,   category:'Education',    method:'Card',    date:d(16),note:'Online course',        created_at:new Date().toISOString()},
    {id:uid(),type:'income', amount:1200, category:'Freelance',    method:'Transfer',date:d(20),note:'Web project',          created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:75,   category:'Entertainment',method:'Card',    date:d(22),note:'Netflix + Spotify',    created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:300,  category:'Shopping',     method:'Card',    date:d(25),note:'Clothes',              created_at:new Date().toISOString()},
    {id:uid(),type:'expense',amount:60,   category:'Transport',    method:'Card',    date:d(28),note:'Fuel',                 created_at:new Date().toISOString()},
  ];
}
</script>
</body>
</html>
`
